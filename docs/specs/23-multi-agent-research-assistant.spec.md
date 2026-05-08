# 23 Multi-Agent Research Assistant Spec

## Purpose

The stock monitor needs a closed-loop AI research workflow and a stock chat assistant. Users should not have to manually click every stock. Holdings attention level controls research cadence and workflow scope, while the assistant can analyze holdings, stock-pool stocks, or arbitrary stock codes using RAG context.

## Requirements

- Research workflow runs by attention level instead of a single manual stock button:
  - high: every 4 hours.
  - medium: every 6 hours.
  - low: every 24 hours.
- The workflow is represented as coordinated Agent steps:
  - market/context collector.
  - product/company information collector.
  - summarizer.
  - risk reviewer.
  - RAG/vector writer.
- Each workflow run stores a job record and step records with status, input summary, output summary, and errors.
- Workflow run can be triggered manually for one attention level from the frontend.
- Workflow run can optionally target one explicit stock for smoke/debug use.
- RAG documents must be saved per user/stock and include attention level, collection cadence, price movement, company/product information, and risk framing.
- Vector database support must be implemented locally through `rag_vectors` JSONB embeddings, with deterministic local embeddings when no external embedding provider is configured.
- Stock chat assistant accepts user input plus stock identity.
- Chat assistant can analyze:
  - a holding stock.
  - a stock-pool stock.
  - an arbitrary stock code.
- Chat assistant uses RAG documents, latest snapshots, company/product profile, holding context, and stock-pool context when available.
- Chat answers are saved for audit/history and must include a safety boundary that they are research only, not trading instructions.

## Non-Goals

- No automatic buy/sell order placement.
- No storage of broker credentials.
- No unsupported web crawling or anti-bot bypass.
- No requirement for a production vector extension; local vector-ready JSONB storage is acceptable for this phase.

## Acceptance Criteria

- Given holdings with different attention levels, running the high workflow only processes high-attention holdings.
- Given workflow runs, `/api/workflows` returns jobs and step summaries.
- Given a workflow completes, `rag_documents` and `rag_vectors` contain a document and indexed local vector for each processed stock.
- Given a user asks the assistant about a holding, the answer includes holding context and RAG context when present.
- Given a user asks the assistant about a stock-pool stock, the answer includes stock-pool context.
- Given a user asks by arbitrary stock code, the assistant returns a research-only answer using available market/profile data.
- Frontend provides controls to run high/medium/low workflows and a stock assistant panel.
