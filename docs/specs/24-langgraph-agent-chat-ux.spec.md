# 24 LangGraph Agent Chat UX Spec

## Purpose

The project must expose a real, inspectable agent workflow instead of only hard-coded backend summaries. Stock research should be coordinated by a Python LangGraph workflow, routed to configured DeepSeek model tiers by task type, then surfaced through a multi-turn streaming stock assistant and a cleaner financial console UI.

## Requirements

- Python agent service must provide visible LangGraph/LangChain-style workflow code:
  - market/context collector.
  - company/product information collector.
  - summarizer.
  - risk reviewer.
  - RAG payload writer.
- The agent must run even when optional LangGraph packages are not installed, using the same node functions in deterministic sequential fallback mode.
- Model routing must be explicit and configurable:
  - lightweight crawling/collection uses `deepseek-v4-flash`.
  - synthesis/chat uses `deepseek-chat`.
  - risk review and higher-stakes reasoning uses `deepseek-v4-pro`.
- Backend research workflow must call the Python agent workflow first and only use the local Go fallback when the agent is unavailable or returns invalid output.
- Workflow step metadata must show which engine and model were used.
- Stock assistant must support multi-turn conversation:
  - per-user session id.
  - recent conversation history.
  - holdings, stock-pool, profile, snapshots, and RAG context.
  - persisted messages.
- Stock assistant must support streaming output for the frontend chat box.
- Frontend stock assistant must be a real chat interface with message bubbles, target selector, and streaming answer updates.
- Stock clicks from holdings, stock pools, and reports must load detail data while respecting a local refresh cooldown to reduce pressure on stock sources.
- Frontend layout must be cleaned up into a denser, modern finance console:
  - stable cards/panels and tables.
  - clearer toolbar spacing.
  - consistent dark/light theme behavior.
  - no broken or dead clickable navigation.

## Non-Goals

- No automatic trading.
- No broker write actions.
- No anti-bot bypass or aggressive scraping.
- No requirement that DeepSeek API calls be made during tests; deterministic fallback remains required.

## Acceptance Criteria

- Python tests can import and run the research workflow and see LangGraph/fallback engine plus model routing metadata.
- Backend workflow run stores agent-returned step summaries and model metadata when the Python agent service is reachable.
- Backend assistant saves and loads chat history by session.
- Streaming assistant endpoint emits incremental chunks and a final payload.
- Frontend chat appends the user message, streams the assistant answer into the current bubble, and preserves context for the session.
- Clicking a stock repeatedly within the cooldown uses cached detail data and shows a clear status message.
- Frontend, backend, and agent test gates pass.
