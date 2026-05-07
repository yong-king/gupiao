from datetime import datetime, timezone

from .indicators import change_percent, moving_average, volatility
from .schema import AnalyzeRequest, AnalyzeResult
from .text_analysis import extract_risks


def analyze(payload):
    request = AnalyzeRequest(
        symbol=payload.get("symbol", ""),
        market=payload.get("market", ""),
        prices=payload.get("prices", []) or [],
        triggered_rules=payload.get("triggered_rules", []) or [],
        texts=payload.get("texts", []) or [],
        source_refs=payload.get("source_refs", []) or [],
    )

    missing = []
    if not request.symbol:
        missing.append("symbol")
    if not request.market:
        missing.append("market")
    if not request.prices:
        missing.append("prices")

    now = datetime.now(timezone.utc).isoformat()
    if missing:
        return AnalyzeResult(
            signal="data_issue",
            confidence=0.0,
            risk_level="low",
            triggered_rules=request.triggered_rules,
            summary="数据不足，无法生成买入或卖出观察。",
            reasoning=["missing_data"],
            data_time=now,
            source_refs=request.source_refs,
            missing_data=missing,
            recommended_action="补充数据后继续观察",
        ).to_dict()

    risks = extract_risks(request.texts)
    latest_change = change_percent(request.prices)
    indicators = {
        "sma_5": moving_average(request.prices, 5),
        "latest_change_percent": latest_change,
        "volatility": volatility(request.prices),
    }

    signal = "hold_watch"
    risk_level = "low"
    confidence = 0.55
    reasoning = []

    if request.triggered_rules:
        signal = "risk_warning"
        risk_level = "medium"
        confidence = 0.72
        reasoning.append("deterministic_rules_triggered")
    if risks:
        risk_level = "high"
        confidence = max(confidence, 0.78)
        reasoning.append("text_risk_markers:" + ",".join(risks))
    if latest_change is not None:
        reasoning.append(f"latest_change_percent:{latest_change:.2f}")

    return AnalyzeResult(
        signal=signal,
        confidence=confidence,
        risk_level=risk_level,
        triggered_rules=request.triggered_rules,
        summary=f"{request.market}:{request.symbol} 当前为{signal}，风险等级 {risk_level}。",
        reasoning=reasoning,
        data_time=now,
        source_refs=request.source_refs,
        missing_data=[],
        recommended_action="继续观察并由人工确认",
        indicators=indicators,
    ).to_dict()
