import unittest

from agent_core.config import Config, LLMConfig
from agent_core.langgraph_workflow import model_for_task, run_research_workflow


class LangGraphWorkflowTest(unittest.TestCase):
    def test_runs_workflow_and_records_model_routing(self):
        cfg = Config(
            llm=LLMConfig(
                chat_model="deepseek-chat",
                flash_model="deepseek-v4-flash",
                pro_model="deepseek-v4-pro",
            )
        )
        result = run_research_workflow(
            {
                "user_id": "u1",
                "job_id": "wf-1",
                "market": "CN",
                "symbol": "000821",
                "attention_level": "high",
                "interval": "4h",
                "snapshots_count": 2,
                "latest_snapshot": {"price": 8.21, "change_percent": 1.2},
                "profile": {"analysis": "钢铁与装备业务需要跟踪订单和原材料价格。", "products": ["钢材", "装备"]},
            },
            cfg,
        )

        self.assertIn(result["engine"], ["langgraph", "fallback-sequential"])
        self.assertEqual(result["metadata"]["model_stock_info_collect"], "deepseek-v4-flash")
        self.assertEqual(result["metadata"]["model_information_summarize"], "deepseek-chat")
        self.assertEqual(result["metadata"]["model_investment_analysis"], "deepseek-v4-pro")
        self.assertEqual(result["metadata"]["rag_schema"], "stock_intelligence_v2")
        self.assertEqual(len(result["steps"]), 5)
        self.assertIn("000821", result["content"])
        self.assertEqual(result["steps"][0]["step_name"], "stock_info_collect")

    def test_model_for_task_defaults_to_chat_model(self):
        cfg = Config(llm=LLMConfig(chat_model="chat", flash_model="flash", pro_model="pro"))
        self.assertEqual(model_for_task(cfg, "stock_info_collect"), "flash")
        self.assertEqual(model_for_task(cfg, "investment_analysis"), "pro")
        self.assertEqual(model_for_task(cfg, "unknown"), "chat")


if __name__ == "__main__":
    unittest.main()
