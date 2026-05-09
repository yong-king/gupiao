import json
import os
import tempfile
import unittest

from agent_core.config import load_config


class ConfigTests(unittest.TestCase):
    def test_load_config_file(self):
        with tempfile.NamedTemporaryFile("w", delete=False) as file:
            json.dump(
                {
                    "host": "0.0.0.0",
                    "port": 9000,
                    "llm": {"model": "test"},
                    "mcp": {
                        "stock_research": {
                            "enabled": True,
                            "name": "yfinance-mcp-server",
                            "repository": "https://github.com/barvhaim/yfinance-mcp-server",
                        }
                    },
                },
                file,
            )
            path = file.name
        try:
            cfg = load_config(path)
            self.assertEqual(cfg.host, "0.0.0.0")
            self.assertEqual(cfg.port, 9000)
            self.assertEqual(cfg.llm.model, "test")
            self.assertTrue(cfg.mcp.stock_research.enabled)
            self.assertEqual(cfg.mcp.stock_research.name, "yfinance-mcp-server")
        finally:
            os.unlink(path)


if __name__ == "__main__":
    unittest.main()
