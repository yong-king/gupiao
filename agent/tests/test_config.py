import json
import os
import tempfile
import unittest

from agent_core.config import load_config


class ConfigTests(unittest.TestCase):
    def test_load_config_file(self):
        with tempfile.NamedTemporaryFile("w", delete=False) as file:
            json.dump({"host": "0.0.0.0", "port": 9000, "llm": {"model": "test"}}, file)
            path = file.name
        try:
            cfg = load_config(path)
            self.assertEqual(cfg.host, "0.0.0.0")
            self.assertEqual(cfg.port, 9000)
            self.assertEqual(cfg.llm.model, "test")
        finally:
            os.unlink(path)


if __name__ == "__main__":
    unittest.main()
