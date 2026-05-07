from dataclasses import dataclass, field, asdict
from typing import Any, Dict, List


@dataclass
class SourceRef:
    type: str
    source: str
    time: str = ""


@dataclass
class AnalyzeRequest:
    symbol: str
    market: str
    prices: List[float] = field(default_factory=list)
    triggered_rules: List[str] = field(default_factory=list)
    texts: List[str] = field(default_factory=list)
    source_refs: List[Dict[str, Any]] = field(default_factory=list)


@dataclass
class AnalyzeResult:
    signal: str
    confidence: float
    risk_level: str
    triggered_rules: List[str]
    summary: str
    reasoning: List[str]
    data_time: str
    source_refs: List[Dict[str, Any]]
    missing_data: List[str]
    recommended_action: str
    indicators: Dict[str, Any] = field(default_factory=dict)

    def to_dict(self):
        return asdict(self)
