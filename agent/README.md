# Jijin Python Agent

This service contains analysis and Agent capabilities. The foundation plan provides a minimal health entry point and tests. Business analysis, FastAPI wiring, technical indicators, and LLM-backed summaries are introduced in later plans.

## Stock Intelligence Workflow

The research workflow is a five-node LangGraph multi-agent flow:

- `stock_info_collect`: stock information Skill Agent. It is configured for the open-source `yfinance-mcp-server` MCP server and falls back to backend profile context when MCP is disabled.
- `trade_market_collect`: trading market and K-line context Agent.
- `information_summarize`: summary Agent.
- `investment_analysis`: analysis and risk framing Agent.
- `rag_vector_write`: RAG/vector payload Agent.

MCP configuration lives in `agent/config/agent.example.json` under `mcp.stock_research`. Keep `enabled=false` until the MCP server is installed and tested locally.
