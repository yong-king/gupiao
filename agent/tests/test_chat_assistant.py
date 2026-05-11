import unittest

from agent_core.chat_assistant import chat
from agent_core.config import Config, LLMConfig


class ChatAssistantTest(unittest.TestCase):
    def test_uses_history_and_rag_context(self):
        result = chat(
            {
                "market": "CN",
                "symbol": "000821",
                "question": "还能继续关注吗？",
                "context_summary": "持仓数量 1000，成本 8.00",
                "history": [{"question": "之前风险是什么？", "answer": "注意波动。"}],
                "rag_documents": [{"content": "订单改善，但原材料价格波动。"}],
                "web_research": [{"source": "public", "summary": "公开信息显示新产品订单增加。"}],
                "data_gaps": ["缺少最近行情样本"],
                "profile": {"analysis": "公司产品以钢铁和装备为主。"},
            },
            Config(llm=LLMConfig(chat_model="deepseek-chat")),
        )

        self.assertEqual(result["model"], "deepseek-chat")
        self.assertEqual(result["provider"], "deepseek")
        self.assertEqual(result["llm_status"], "missing-api-key")
        self.assertIn("先给结论", result["answer"])
        self.assertIn("继续研究观察", result["answer"])
        self.assertIn("缺少最近行情样本", result["answer"])
        self.assertIn("订单改善", result["answer"])
        self.assertIn("不构成买卖指令", result["answer"])

    def test_answers_model_question_directly(self):
        result = chat(
            {
                "market": "CN",
                "symbol": "002555",
                "question": "测试，你是什么模型",
                "context_summary": "持仓数量 1000，成本 23.00",
                "history": [],
                "rag_documents": [],
                "web_research": [],
                "data_gaps": [],
                "profile": {"analysis": "主营业务为移动游戏与内容运营。"},
            },
            Config(llm=LLMConfig(chat_model="deepseek-chat")),
        )

        self.assertEqual(result["llm_status"], "model-info")
        self.assertIn("deepseek-chat", result["answer"])
        self.assertNotIn("股票助手分析", result["answer"])


if __name__ == "__main__":
    unittest.main()
