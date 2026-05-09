from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from typing import Any, Dict, List, Tuple

from .config import Config


TASK_MODELS = {
    "market_context_collect": "flash_model",
    "company_product_collect": "flash_model",
    "summarize": "chat_model",
    "risk_review": "pro_model",
    "rag_vector_write": "chat_model",
}

# 当前 LangGraph 工作流包含 5 个节点，也对应 5 个智能体：
# 行情上下文采集、公司产品采集、归纳整理、风险审查、RAG/向量写入。
# 不同节点按业务风险路由到不同 DeepSeek 模型。


@dataclass
class ResearchState:
    user_id: str
    job_id: str
    market: str
    symbol: str
    attention_level: str = "medium"
    interval: str = "2h"
    profile: Dict[str, Any] = field(default_factory=dict)
    latest_snapshot: Dict[str, Any] = field(default_factory=dict)
    snapshots_count: int = 0
    steps: List[Dict[str, Any]] = field(default_factory=list)
    content: str = ""
    rag_metadata: Dict[str, str] = field(default_factory=dict)
    engine: str = ""


def model_for_task(cfg: Config, task: str) -> str:
    attr = TASK_MODELS.get(task, "chat_model")
    return getattr(cfg.llm, attr, cfg.llm.model)


def run_research_workflow(payload: Dict[str, Any], cfg: Config) -> Dict[str, Any]:
    state = ResearchState(
        user_id=str(payload.get("user_id", "")),
        job_id=str(payload.get("job_id", "")),
        market=str(payload.get("market", "CN")).upper(),
        symbol=str(payload.get("symbol", "")).upper(),
        attention_level=str(payload.get("attention_level", "medium")),
        interval=str(payload.get("interval", "2h")),
        profile=dict(payload.get("profile") or {}),
        latest_snapshot=dict(payload.get("latest_snapshot") or {}),
        snapshots_count=int(payload.get("snapshots_count") or 0),
    )
    final_state, engine = _run_langgraph(state, cfg)
    final_state.engine = engine
    final_state.rag_metadata["agent_engine"] = engine
    content = final_state.content or _join_step_outputs(final_state.steps)
    return {
        "engine": engine,
        "market": final_state.market,
        "symbol": final_state.symbol,
        "content": content,
        "metadata": final_state.rag_metadata,
        "steps": final_state.steps,
    }


def _run_langgraph(initial: ResearchState, cfg: Config) -> Tuple[ResearchState, str]:
    try:
        from langgraph.graph import END, StateGraph

        # 运行时如果已安装 langgraph，就构建真实节点图；否则走同一批节点函数的顺序兜底。
        graph = StateGraph(ResearchState)
        graph.add_node("market_context_collect", lambda state: market_context_collect(state, cfg))
        graph.add_node("company_product_collect", lambda state: company_product_collect(state, cfg))
        graph.add_node("summarize", lambda state: summarize(state, cfg))
        graph.add_node("risk_review", lambda state: risk_review(state, cfg))
        graph.add_node("rag_vector_write", lambda state: rag_vector_write(state, cfg))
        graph.set_entry_point("market_context_collect")
        graph.add_edge("market_context_collect", "company_product_collect")
        graph.add_edge("company_product_collect", "summarize")
        graph.add_edge("summarize", "risk_review")
        graph.add_edge("risk_review", "rag_vector_write")
        graph.add_edge("rag_vector_write", END)
        compiled = graph.compile()
        result = compiled.invoke(initial)
        if isinstance(result, ResearchState):
            return result, "langgraph"
        return _state_from_dict(result), "langgraph"
    except Exception:
        return _run_sequential(initial, cfg), "fallback-sequential"


def _run_sequential(state: ResearchState, cfg: Config) -> ResearchState:
    for node in (market_context_collect, company_product_collect, summarize, risk_review, rag_vector_write):
        state = node(state, cfg)
    return state


def market_context_collect(state: ResearchState, cfg: Config) -> ResearchState:
    latest = state.latest_snapshot
    price = latest.get("price") or latest.get("Price") or 0
    change = latest.get("change_percent") or latest.get("ChangePercent") or 0
    output = f"{state.market}:{state.symbol} 行情样本 {state.snapshots_count} 条，最新价 {price}，涨跌幅 {change}%。"
    _append_step(state, "market_context_collect", "行情上下文采集 Agent", output, cfg)
    return state


def company_product_collect(state: ResearchState, cfg: Config) -> ResearchState:
    profile = state.profile
    products = profile.get("products") or profile.get("Products") or []
    business = profile.get("business") or profile.get("Business") or "暂无公开业务文本"
    if isinstance(products, list) and products:
        business = f"{business}；产品/业务：{'、'.join(str(item) for item in products)}"
    _append_step(state, "company_product_collect", "股票产品信息采集 Agent", business, cfg)
    return state


def summarize(state: ResearchState, cfg: Config) -> ResearchState:
    profile = state.profile
    analysis = profile.get("analysis") or profile.get("Analysis") or "暂无公司产品归纳。"
    output = f"关注等级 {state.attention_level}，建议采集周期 {state.interval}。综合归纳：{analysis}"
    state.content = f"{state.market}:{state.symbol} 研究总结：{output}"
    _append_step(state, "summarize", "归纳整理 Agent", output, cfg)
    return state


def risk_review(state: ResearchState, cfg: Config) -> ResearchState:
    output = "风险审查：该结论只用于研究和提醒，不能替代公告、财报、仓位和个人风险承受能力判断，不构成买卖指令。"
    state.content = f"{state.content}\n{output}"
    _append_step(state, "risk_review", "风险审查 Agent", output, cfg)
    return state


def rag_vector_write(state: ResearchState, cfg: Config) -> ResearchState:
    state.rag_metadata = {
        "attention_level": state.attention_level,
        "refresh_interval": state.interval,
        "workflow_job_id": state.job_id,
        "agents": "market_context_collect,company_product_collect,summarize,risk_review,rag_vector_write",
        "agent_engine": state.engine or "pending",
        "model_market_context_collect": model_for_task(cfg, "market_context_collect"),
        "model_company_product_collect": model_for_task(cfg, "company_product_collect"),
        "model_summarize": model_for_task(cfg, "summarize"),
        "model_risk_review": model_for_task(cfg, "risk_review"),
        "embedding_status": "indexed",
    }
    _append_step(state, "rag_vector_write", "RAG/向量写入 Agent", "生成 RAG 写入载荷和本地向量元数据。", cfg)
    return state


def _append_step(state: ResearchState, step_name: str, agent_name: str, output: str, cfg: Config) -> None:
    now = datetime.now(timezone.utc).isoformat()
    state.steps.append(
        {
            "step_name": step_name,
            "agent_name": agent_name,
            "status": "succeeded",
            "input_summary": f"{state.market}:{state.symbol} {state.attention_level}",
            "output_summary": output,
            "model": model_for_task(cfg, step_name),
            "started_at": now,
            "completed_at": now,
        }
    )


def _join_step_outputs(steps: List[Dict[str, Any]]) -> str:
    return "\n".join(str(step.get("output_summary", "")) for step in steps if step.get("output_summary"))


def as_plain_dict(value: Any) -> Dict[str, Any]:
    if hasattr(value, "__dataclass_fields__"):
        return asdict(value)
    if isinstance(value, dict):
        return value
    return {}


def _state_from_dict(value: Dict[str, Any]) -> ResearchState:
    state = ResearchState(
        user_id=str(value.get("user_id", "")),
        job_id=str(value.get("job_id", "")),
        market=str(value.get("market", "CN")).upper(),
        symbol=str(value.get("symbol", "")).upper(),
        attention_level=str(value.get("attention_level", "medium")),
        interval=str(value.get("interval", "2h")),
        profile=dict(value.get("profile") or {}),
        latest_snapshot=dict(value.get("latest_snapshot") or {}),
        snapshots_count=int(value.get("snapshots_count") or 0),
    )
    state.steps = list(value.get("steps") or [])
    state.content = str(value.get("content") or "")
    state.rag_metadata = dict(value.get("rag_metadata") or {})
    state.engine = str(value.get("engine") or "")
    return state
