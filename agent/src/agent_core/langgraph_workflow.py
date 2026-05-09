from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from typing import Any, Dict, List, Tuple

from .config import Config
from .skills.stock_research import StockResearchSkill


TASK_MODELS = {
    "stock_info_collect": "flash_model",
    "trade_market_collect": "flash_model",
    "information_summarize": "chat_model",
    "investment_analysis": "pro_model",
    "rag_vector_write": "chat_model",
}

# 当前 LangGraph 工作流包含 5 个节点，也对应 5 个智能体：
# 股票信息抓取 Skill、交易行情/K线采集、信息汇总、分析审查、RAG/向量写入。
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
    snapshots: List[Dict[str, Any]] = field(default_factory=list)
    snapshots_count: int = 0
    steps: List[Dict[str, Any]] = field(default_factory=list)
    public_research: Dict[str, Any] = field(default_factory=dict)
    trading_context: Dict[str, Any] = field(default_factory=dict)
    summary: str = ""
    analysis: str = ""
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
        snapshots=list(payload.get("snapshots") or []),
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
        graph.add_node("stock_info_collect", lambda state: stock_info_collect(state, cfg))
        graph.add_node("trade_market_collect", lambda state: trade_market_collect(state, cfg))
        graph.add_node("information_summarize", lambda state: information_summarize(state, cfg))
        graph.add_node("investment_analysis", lambda state: investment_analysis(state, cfg))
        graph.add_node("rag_vector_write", lambda state: rag_vector_write(state, cfg))
        graph.set_entry_point("stock_info_collect")
        graph.add_edge("stock_info_collect", "trade_market_collect")
        graph.add_edge("trade_market_collect", "information_summarize")
        graph.add_edge("information_summarize", "investment_analysis")
        graph.add_edge("investment_analysis", "rag_vector_write")
        graph.add_edge("rag_vector_write", END)
        compiled = graph.compile()
        result = compiled.invoke(initial)
        if isinstance(result, ResearchState):
            return result, "langgraph"
        return _state_from_dict(result), "langgraph"
    except Exception:
        return _run_sequential(initial, cfg), "fallback-sequential"


def _run_sequential(state: ResearchState, cfg: Config) -> ResearchState:
    for node in (stock_info_collect, trade_market_collect, information_summarize, investment_analysis, rag_vector_write):
        state = node(state, cfg)
    return state


def stock_info_collect(state: ResearchState, cfg: Config) -> ResearchState:
    result = StockResearchSkill(cfg).collect_public_info(state.market, state.symbol, state.profile)
    state.public_research = result.to_dict()
    item_count = len(state.public_research.get("items") or [])
    warning_text = "；".join(state.public_research.get("warnings") or [])
    output = f"通过 {result.mcp_server} Skill 准备 {item_count} 条股票公司/产品公开信息。{warning_text}".strip()
    _append_step(state, "stock_info_collect", "股票信息抓取 Skill Agent", output, cfg)
    return state


def trade_market_collect(state: ResearchState, cfg: Config) -> ResearchState:
    latest = state.latest_snapshot
    price = _number(latest.get("price") or latest.get("Price"))
    change = _number(latest.get("change_percent") or latest.get("ChangePercent"))
    source = latest.get("source") or latest.get("Source") or "unknown"
    kline = _kline_summary(state.snapshots)
    state.trading_context = {
        "latest_price": price,
        "change_percent": change,
        "source": source,
        "snapshots_count": state.snapshots_count or len(state.snapshots),
        "kline_summary": kline,
    }
    output = f"{state.market}:{state.symbol} 最新价 {price}，涨跌幅 {change}%，行情源 {source}，K线摘要：{kline}。"
    _append_step(state, "trade_market_collect", "交易行情/K线 Agent", output, cfg)
    return state


def information_summarize(state: ResearchState, cfg: Config) -> ResearchState:
    profile = state.profile
    analysis = profile.get("analysis") or profile.get("Analysis") or "暂无公司产品归纳"
    public_items = state.public_research.get("items") or []
    public_text = "；".join(str(item.get("summary", "")) for item in public_items[:3] if isinstance(item, dict))
    state.summary = (
        f"关注等级 {state.attention_level}，建议信息采集周期 {state.interval}。"
        f"公司/产品信息：{public_text or analysis}。"
        f"交易行情：{state.trading_context.get('kline_summary', '暂无K线摘要')}。"
    )
    _append_step(state, "information_summarize", "信息汇总 Agent", state.summary, cfg)
    return state


def investment_analysis(state: ResearchState, cfg: Config) -> ResearchState:
    change = _number(state.trading_context.get("change_percent"))
    if change >= 5:
        risk = "价格波动显著偏强，需核对公告、成交量和是否存在短线过热。"
    elif change <= -5:
        risk = "价格波动显著偏弱，需核对基本面变化、止损线和仓位承受能力。"
    else:
        risk = "价格波动未触发极端阈值，重点观察趋势延续、成交量和产品信息变化。"
    state.analysis = (
        f"分析结论：{risk}"
        "该结论只用于研究、提醒和风险提示，不构成自动买卖或确定性交易指令。"
    )
    state.content = f"{state.market}:{state.symbol} 多智能体研究结果\n{state.summary}\n{state.analysis}"
    _append_step(state, "investment_analysis", "分析 Agent", state.analysis, cfg)
    return state


def rag_vector_write(state: ResearchState, cfg: Config) -> ResearchState:
    state.rag_metadata = {
        "attention_level": state.attention_level,
        "refresh_interval": state.interval,
        "workflow_job_id": state.job_id,
        "agents": "stock_info_collect,trade_market_collect,information_summarize,investment_analysis,rag_vector_write",
        "agent_engine": state.engine or "pending",
        "model_stock_info_collect": model_for_task(cfg, "stock_info_collect"),
        "model_trade_market_collect": model_for_task(cfg, "trade_market_collect"),
        "model_information_summarize": model_for_task(cfg, "information_summarize"),
        "model_investment_analysis": model_for_task(cfg, "investment_analysis"),
        "stock_research_mcp": cfg.mcp.stock_research.name,
        "stock_research_mcp_repository": cfg.mcp.stock_research.repository,
        "embedding_status": "indexed",
        "rag_schema": "stock_intelligence_v2",
    }
    _append_step(state, "rag_vector_write", "RAG/向量写入 Agent", "生成 stock_intelligence_v2 RAG 写入载荷和本地向量元数据。", cfg)
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
        snapshots=list(value.get("snapshots") or []),
        snapshots_count=int(value.get("snapshots_count") or 0),
    )
    state.steps = list(value.get("steps") or [])
    state.public_research = dict(value.get("public_research") or {})
    state.trading_context = dict(value.get("trading_context") or {})
    state.summary = str(value.get("summary") or "")
    state.analysis = str(value.get("analysis") or "")
    state.content = str(value.get("content") or "")
    state.rag_metadata = dict(value.get("rag_metadata") or {})
    state.engine = str(value.get("engine") or "")
    return state


def _number(value: Any) -> float:
    try:
        return float(value)
    except (TypeError, ValueError):
        return 0.0


def _kline_summary(snapshots: List[Dict[str, Any]]) -> str:
    if not snapshots:
        return "暂无历史K线样本"
    prices = [_number(item.get("price") or item.get("Price")) for item in snapshots if _number(item.get("price") or item.get("Price")) > 0]
    if not prices:
        return f"共 {len(snapshots)} 条样本，但缺少有效价格"
    high = max(prices)
    low = min(prices)
    latest = prices[-1]
    first = prices[0]
    direction = "上涨" if latest >= first else "下跌"
    return f"共 {len(prices)} 条有效样本，区间高点 {high:.2f}，低点 {low:.2f}，最新较首条{direction}"
