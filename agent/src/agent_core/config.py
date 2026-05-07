import json
import os
from dataclasses import dataclass


@dataclass
class LLMConfig:
    provider: str = "deepseek"
    model: str = "deepseek-chat"
    api_key_env: str = "DEEPSEEK_API_KEY"


@dataclass
class Config:
    host: str = "127.0.0.1"
    port: int = 8090
    llm: LLMConfig = LLMConfig()


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
                api_key_env=llm.get("api_key_env", cfg.llm.api_key_env),
            ),
        )
    cfg.host = os.getenv("AGENT_HOST", cfg.host)
    cfg.port = int(os.getenv("AGENT_PORT", cfg.port))
    return cfg
