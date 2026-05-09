from dataclasses import asdict, dataclass, field
from datetime import datetime, timezone
from typing import Any, Dict, List

from ..config import Config


@dataclass
class ResearchItem:
    source: str
    title: str
    summary: str
    url: str = ""
    published_at: str = ""
    metadata: Dict[str, Any] = field(default_factory=dict)


@dataclass
class StockResearchResult:
    provider: str
    mcp_server: str
    items: List[ResearchItem]
    source_refs: List[Dict[str, Any]]
    warnings: List[str] = field(default_factory=list)

    def to_dict(self) -> Dict[str, Any]:
        return {
            "provider": self.provider,
            "mcp_server": self.mcp_server,
            "items": [asdict(item) for item in self.items],
            "source_refs": self.source_refs,
            "warnings": self.warnings,
        }


class StockResearchSkill:
    """股票信息抓取 Skill。

    当前实现先把开源 MCP 的配置纳入工作流上下文，并使用后端传入的 profile
    生成可审计公开信息条目。后续接入 MCP stdio/http 客户端时，只需要在
    collect_public_info 中替换 `_profile_items`，工作流节点和 RAG schema 不变。
    """

    def __init__(self, cfg: Config):
        self.cfg = cfg

    def collect_public_info(self, market: str, symbol: str, profile: Dict[str, Any]) -> StockResearchResult:
        mcp = self.cfg.mcp.stock_research
        items = self._profile_items(market, symbol, profile)
        refs = [
            {
                "type": "mcp_server",
                "source": mcp.name,
                "repository": mcp.repository,
                "tool_hint": ",".join(mcp.tools),
                "time": datetime.now(timezone.utc).isoformat(),
            }
        ]
        warnings = []
        if not mcp.enabled:
            warnings.append("stock_research MCP 未启用，使用后端 profile 公开信息作为降级输入。")
        return StockResearchResult(
            provider=mcp.provider,
            mcp_server=mcp.name,
            items=items,
            source_refs=refs,
            warnings=warnings,
        )

    def _profile_items(self, market: str, symbol: str, profile: Dict[str, Any]) -> List[ResearchItem]:
        name = str(profile.get("name") or profile.get("Name") or symbol)
        business = str(profile.get("business") or profile.get("Business") or "暂无业务描述")
        analysis = str(profile.get("analysis") or profile.get("Analysis") or "暂无公开分析")
        products = profile.get("products") or profile.get("Products") or []
        if isinstance(products, list):
            product_text = "、".join(str(item) for item in products if str(item).strip())
        else:
            product_text = str(products)
        summary = f"{name} 公开业务：{business}；主要产品/业务：{product_text or '暂无'}；已有归纳：{analysis}"
        return [
            ResearchItem(
                source="backend_profile",
                title=f"{market}:{symbol} 公司/产品公开信息",
                summary=summary,
                metadata={"market": market, "symbol": symbol, "name": name},
            )
        ]
