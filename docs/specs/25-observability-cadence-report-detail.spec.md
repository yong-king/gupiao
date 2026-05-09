# 25 Observability Cadence Report Detail Spec

## Purpose

The console needs clearer operational boundaries: report pages should summarize today's movement and reminders only, stock-specific analysis should live on a detail page, and AI/crawler activity must be observable through logs. Research and quote refresh cadence must be configurable by holding attention level.

## Requirements

- LangGraph workflow must document its node/agent count and model routing.
- If LangGraph is not available in the runtime, the project must provide dependency installation through agent requirements and Docker build.
- Holding product/company information collection cadence:
  - high: 1 hour.
  - medium: 2 hours.
  - low: 4 hours.
- Realtime stock quote/detail collection cadence:
  - high: 2 minutes.
  - medium: 5 minutes.
  - low: 10 minutes.
- Cadence values must be configurable in backend JSON config and documented.
- Analysis report page must only show:
  - holding stocks and current stock-pool stocks.
  - latest price when available.
  - today change using Chinese market convention: rising red, falling green.
  - buy/sell/reminder status.
- Stock-specific charts, K-line-like price curve, daily change calendar, daily records, holding details, RAG/product summaries, and AI workflow controls must move to a stock detail page opened by clicking a stock code.
- Operation log system must record:
  - AI model call component, model, input summary, output summary, and metadata.
  - crawler/market-data collection target, source/provider, input summary, output summary, and metadata.
  - RAG writes and workflow step outputs.
- Frontend must provide a log monitoring page to inspect recent operation logs.
- Stock assistant must use RAG documents, profile/product context, recent market data, and public-information collection summaries when forming the answer.
- Code changes should include succinct Chinese comments where they clarify workflow, logging, or cadence behavior.

## Non-Goals

- No automatic trading or broker write action.
- No anti-bot bypassing or aggressive scraping.
- No production scheduler daemon in this plan; this plan configures cadence and uses it in manual/API paths.

## Acceptance Criteria

- `/api/operation-logs?user_id=...` returns recent logs for AI, crawler, and workflow actions.
- Running market collect writes a crawler/quote log.
- Running workflow writes AI/agent step logs with model names.
- Asking stock assistant writes a model log with RAG document ids in metadata.
- Analysis report lists latest price, today's change, and reminder status but does not show charts or product analysis.
- Clicking a stock code opens stock detail with charts, daily records, profile/RAG summary, and workflow controls.
- Backend config exposes product-research and quote refresh cadence by attention level.
- Tests pass for backend, frontend, and agent.
