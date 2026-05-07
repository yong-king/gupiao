import unittest

from agent_core import health_payload


class HealthTests(unittest.TestCase):
    def test_health_payload_shape(self):
        self.assertEqual(
            health_payload(),
            {
                "status": "ok",
                "service": "agent",
            },
        )


if __name__ == "__main__":
    unittest.main()
