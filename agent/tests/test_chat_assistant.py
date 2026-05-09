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
                "profile": {"analysis": "公司产品以钢铁和装备为主。"},
            },
            Config(llm=LLMConfig(chat_model="deepseek-chat")),
        )

        self.assertEqual(result["model"], "deepseek-chat")
        self.assertIn("历史对话", result["answer"])
        self.assertIn("网络公开信息", result["answer"])
        self.assertIn("订单改善", result["answer"])
        self.assertIn("不构成买卖指令", result["answer"])


if __name__ == "__main__":
    unittest.main()
