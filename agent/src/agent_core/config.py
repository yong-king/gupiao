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
class Config:
    host: str = "127.0.0.1"
    port: int = 8090
    llm: LLMConfig = field(default_factory=LLMConfig)


def load_config(path=None):
    cfg = Config()
    path = path or os.getenv("JIJIN_AGENT_CONFIG")
    if path:
        with open(path, "r", encoding="utf-8") as file:
            raw = json.load(file)
        llm = raw.get("llm", {})
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
        )
    cfg.host = os.getenv("AGENT_HOST", cfg.host)
    cfg.port = int(os.getenv("AGENT_PORT", cfg.port))
    cfg.llm.model = os.getenv("DEEPSEEK_MODEL", cfg.llm.model)
    cfg.llm.chat_model = os.getenv("DEEPSEEK_CHAT_MODEL", cfg.llm.chat_model)
    cfg.llm.flash_model = os.getenv("DEEPSEEK_FLASH_MODEL", cfg.llm.flash_model)
    cfg.llm.pro_model = os.getenv("DEEPSEEK_PRO_MODEL", cfg.llm.pro_model)
    return cfg
