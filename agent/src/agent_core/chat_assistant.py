from datetime import datetime, timezone
from typing import Any, Dict, List

from .config import Config


def chat(payload: Dict[str, Any], cfg: Config) -> Dict[str, Any]:
    market = str(payload.get("market", "CN")).upper()
    symbol = str(payload.get("symbol", "")).upper()
    question = str(payload.get("question", "")).strip()
    history = _history_lines(payload.get("history") or [])
    context = str(payload.get("context_summary") or "")
    rag_docs = payload.get("rag_documents") or []
    rag = str(rag_docs[0].get("content") or rag_docs[0].get("Content")) if rag_docs else "暂无 RAG 历史总结"
    web_research = _web_research_lines(payload.get("web_research") or [])
    profile = payload.get("profile") or {}
    profile_text = profile.get("analysis") or profile.get("Analysis") or profile.get("business") or profile.get("Business") or "暂无公司产品信息"
    answer = (
        f"{market}:{symbol} 多轮分析：{profile_text}。"
        f" 当前上下文：{context or '暂无持仓/股票池上下文'}。"
        f" 历史对话：{history or '本轮是该会话的首个问题'}。"
        f" 网络公开信息：{web_research or '暂无新增公开信息'}。"
        f" RAG 参考：{rag}。"
        f" 针对你的问题「{question}」，建议把价格波动、关注等级、产品变化和仓位成本一起核对；这里只做研究提醒，不构成买卖指令。"
    )
    return {
        "market": market,
        "symbol": symbol,
        "answer": answer,
        "model": cfg.llm.chat_model,
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
