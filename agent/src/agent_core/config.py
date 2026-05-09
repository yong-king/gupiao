import json
import os
from dataclasses import dataclass, field


@dataclass
class LLMConfig:
    provider: str = "deepseek"
    model: str = "deepseek-chat"
    chat_model: str = "deepseek-chat"
    flash_model: str = "deepseek-v4-flash"
    pro_model: str = "deepseek-v4-pro"
    api_key_env: str = "DEEPSEEK_API_KEY"


@dataclass
class MCPServerConfig:
    enabled: bool = False
    provider: str = "yfinance"
    name: str = "yfinance-mcp-server"
    repository: str = "https://github.com/barvhaim/yfinance-mcp-server"
    command: str = "uvx"
    args: list = field(default_factory=lambda: ["yfinance-mcp-server"])
    tools: list = field(default_factory=lambda: ["get_stock_info", "get_news", "get_history"])


@dataclass
class MCPConfig:
    stock_research: MCPServerConfig = field(default_factory=MCPServerConfig)


@dataclass
class Config:
    host: str = "127.0.0.1"
    port: int = 8090
    llm: LLMConfig = field(default_factory=LLMConfig)
    mcp: MCPConfig = field(default_factory=MCPConfig)


def load_config(path=None):
    cfg = Config()
    path = path or os.getenv("JIJIN_AGENT_CONFIG")
    if path:
        with open(path, "r", encoding="utf-8") as file:
            raw = json.load(file)
        llm = raw.get("llm", {})
        mcp = raw.get("mcp", {})
        stock_research = mcp.get("stock_research", {})
        cfg = Config(
            host=raw.get("host", cfg.host),
            port=int(raw.get("port", cfg.port)),
            llm=LLMConfig(
                provider=llm.get("provider", cfg.llm.provider),
                model=llm.get("model", cfg.llm.model),
                chat_model=llm.get("chat_model", llm.get("model", cfg.llm.chat_model)),
                flash_model=llm.get("flash_model", cfg.llm.flash_model),
                pro_model=llm.get("pro_model", cfg.llm.pro_model),
                api_key_env=llm.get("api_key_env", cfg.llm.api_key_env),
            ),
            mcp=MCPConfig(
                stock_research=MCPServerConfig(
                    enabled=bool(stock_research.get("enabled", cfg.mcp.stock_research.enabled)),
                    provider=stock_research.get("provider", cfg.mcp.stock_research.provider),
                    name=stock_research.get("name", cfg.mcp.stock_research.name),
                    repository=stock_research.get("repository", cfg.mcp.stock_research.repository),
                    command=stock_research.get("command", cfg.mcp.stock_research.command),
                    args=list(stock_research.get("args", cfg.mcp.stock_research.args)),
                    tools=list(stock_research.get("tools", cfg.mcp.stock_research.tools)),
                )
            ),
        )
    cfg.host = os.getenv("AGENT_HOST", cfg.host)
    cfg.port = int(os.getenv("AGENT_PORT", cfg.port))
    cfg.llm.model = os.getenv("DEEPSEEK_MODEL", cfg.llm.model)
    cfg.llm.chat_model = os.getenv("DEEPSEEK_CHAT_MODEL", cfg.llm.chat_model)
    cfg.llm.flash_model = os.getenv("DEEPSEEK_FLASH_MODEL", cfg.llm.flash_model)
    cfg.llm.pro_model = os.getenv("DEEPSEEK_PRO_MODEL", cfg.llm.pro_model)
    cfg.mcp.stock_research.enabled = os.getenv("STOCK_RESEARCH_MCP_ENABLED", str(cfg.mcp.stock_research.enabled)).lower() == "true"
    return cfg
