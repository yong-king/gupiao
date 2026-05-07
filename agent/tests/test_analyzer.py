import unittest

from agent_core.analyzer import analyze
from agent_core.indicators import change_percent, moving_average
from agent_core.text_analysis import extract_risks


class AnalyzerTests(unittest.TestCase):
    def test_missing_data_returns_data_issue(self):
        result = analyze({"symbol": "AAPL", "market": "US", "prices": []})

        self.assertEqual(result["signal"], "data_issue")
        self.assertIn("prices", result["missing_data"])

    def test_analyze_returns_structured_result(self):
        result = analyze({
            "symbol": "AAPL",
            "market": "US",
            "prices": [100, 101, 102, 103, 104],
            "triggered_rules": ["price_above"],
            "texts": ["Company faces investigation"],
            "source_refs": [{"type": "price", "source": "mock"}],
        })

        self.assertEqual(result["signal"], "risk_warning")
        self.assertEqual(result["risk_level"], "high")
        self.assertEqual(result["triggered_rules"], ["price_above"])
        self.assertEqual(result["missing_data"], [])
        self.assertIn("sma_5", result["indicators"])

    def test_indicators(self):
        self.assertEqual(moving_average([1, 2, 3, 4, 5], 5), 3)
        self.assertEqual(change_percent([100, 110]), 10)

    def test_extract_risks(self):
        self.assertEqual(extract_risks(["重大诉讼 and downgrade"]), ["downgrade", "lawsuit"])


if __name__ == "__main__":
    unittest.main()
