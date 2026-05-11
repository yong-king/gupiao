from datetime import datetime, timezone
from typing import Any, Dict, List

from .config import Config
from .llm import complete_text, first_text


def chat(payload: Dict[str, Any], cfg: Config) -> Dict[str, Any]:
    market = str(payload.get("market", "CN")).upper()
    symbol = str(payload.get("symbol", "")).upper()
    question = str(payload.get("question", "")).strip()
    if _is_model_question(question):
        return {
            "market": market,
            "symbol": symbol,
            "answer": (
                f"当前股票助手对话默认使用 {cfg.llm.chat_model}。"
                " 它负责研究问答、风险提示和总结，不负责自动交易。"
                " 仅供研究提醒，不构成买卖指令。"
            ),
            "model": cfg.llm.chat_model,
            "provider": cfg.llm.provider,
            "llm_status": "model-info",
            "created_at": datetime.now(timezone.utc).isoformat(),
        }
    history = _history_lines(payload.get("history") or [])
    context = str(payload.get("context_summary") or "")
    rag_docs = payload.get("rag_documents") or []
    rag = str(rag_docs[0].get("content") or rag_docs[0].get("Content")) if rag_docs else "暂无 RAG 历史总结"
    web_research = _web_research_lines(payload.get("web_research") or [])
    data_gaps = _data_gap_lines(payload.get("data_gaps") or [])
    profile = payload.get("profile") or {}
    profile_name = _trim_sentence_end(profile.get("name") or profile.get("Name") or "")
    profile_text = _trim_sentence_end(profile.get("analysis") or profile.get("Analysis") or profile.get("business") or profile.get("Business") or "暂无公司产品信息")
    # 缺少行情、持仓或 RAG 时，兜底回答优先直接回答问题，不复读整段原始上下文。
    fallback_answer = _fallback_answer(
        market=market,
        symbol=symbol,
        question=question,
        context=context,
        data_gaps=data_gaps,
        rag=rag,
        profile_name=profile_name,
        profile_text=profile_text,
        model=cfg.llm.chat_model,
    )
    system_prompt = (
        "你是股票研究助手。"
        "只能输出研究、提醒和风险提示，不得输出自动下单指令。"
        "优先直接回答用户问题，不要把上下文、RAG、公开信息逐段复读给用户。"
        "如果信息不足，要明确指出缺口，并给出下一步观察建议。"
        "如果用户在问你使用的模型或能力边界，要直接、简短、明确回答。"
        "回答请使用中文，结构清晰，控制在 3 到 6 个短段落或要点内。"
    )
    user_prompt = (
        f"股票：{market}:{symbol}\n"
        f"用户问题：{question}\n"
        f"当前模型：{cfg.llm.chat_model}\n"
        f"上下文：{context}\n"
        f"历史对话：{history}\n"
        f"公开信息：{web_research}\n"
        f"数据缺口：{data_gaps}\n"
        f"RAG：{rag}\n"
        f"公司信息：{profile_text}\n"
        "请先判断用户是在问模型信息、结论分析、风险、还是下一步动作。"
        "回答时优先给结论，再给依据，再给风险或数据缺口。"
        "除非用户明确要求，不要把“当前上下文/历史对话/RAG/公开信息”这些标签原样写进回答。"
        "最后保留一句“仅供研究提醒，不构成买卖指令”。"
    )
    llm = complete_text(cfg, cfg.llm.chat_model, system_prompt, user_prompt)
    answer = first_text(llm, fallback_answer)
    return {
        "market": market,
        "symbol": symbol,
        "answer": answer,
        "model": cfg.llm.chat_model,
        "provider": llm.provider,
        "llm_status": llm.status,
        "created_at": datetime.now(timezone.utc).isoformat(),
    }


def _history_lines(history: List[Dict[str, Any]]) -> str:
    lines = []
    for item in history[-6:]:
        question = str(item.get("question") or item.get("Question") or "").strip()
        answer = str(item.get("answer") or item.get("Answer") or "").strip()
        if question:
            lines.append("用户：" + question)
        if answer:
            lines.append("助手：" + answer[:120])
    return "；".join(lines)


def _web_research_lines(items: List[Dict[str, Any]]) -> str:
    lines = []
    for item in items[:4]:
        source = str(item.get("source") or "public")
        summary = str(item.get("summary") or "").strip()
        if summary:
            lines.append(f"{source}: {summary[:220]}")
    return "；".join(lines)


def _data_gap_lines(items: List[Any]) -> str:
    gaps = [str(item).strip() for item in items if str(item).strip()]
    return "；".join(gaps[:6])


def _trim_sentence_end(value: Any) -> str:
    return str(value).strip().rstrip("。.!！")


def _is_model_question(question: str) -> bool:
    lowered = question.lower()
    return any(token in lowered for token in ["什么模型", "哪个模型", "你是什么模型", "llm"])


def _fallback_answer(
    market: str,
    symbol: str,
    question: str,
    context: str,
    data_gaps: str,
    rag: str,
    profile_name: str,
    profile_text: str,
    model: str,
) -> str:
    if _is_model_question(question):
        return (
            f"当前股票助手优先使用 {model} 处理对话分析。"
            " 如果大模型不可用，会降级为本地研究提醒模式。"
            " 仅供研究提醒，不构成买卖指令。"
        )

    company = profile_name or f"{market}:{symbol}"
    basis = profile_text if profile_text and profile_text != "暂无公司产品信息" else f"{company} 的公司与产品资料还不完整"
    gap_line = f"当前主要数据缺口：{data_gaps}。" if data_gaps and data_gaps != "暂无明显缺口" else ""
    rag_hint = ""
    if rag and rag != "暂无 RAG 历史总结":
        rag_hint = f" 历史研究里提到：{rag[:120]}。"
    context_hint = f" 当前持仓/关注上下文：{context}。" if context else ""
    return (
        f"先给结论：{company} 目前更适合继续研究观察，暂时不适合只凭现有信息下确定性判断。"
        f" 依据主要有两点：{basis}。{context_hint}{rag_hint}"
        f" {gap_line}"
        "下一步建议优先补齐最新行情、业务变化和公告信息，再结合你的持仓成本看风险收益是否匹配。"
        " 仅供研究提醒，不构成买卖指令。"
    ).replace("  ", " ").strip()
