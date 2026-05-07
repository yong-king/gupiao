# Shared Contracts

## Time

- Store timestamps as UTC.
- API payloads use RFC 3339 strings.
- UI must display the source data time separately from refresh time.

## Market Codes

Initial market code set:

- `US`
- `HK`
- `CN`

## Symbol Format

- Symbol is required.
- Market is required.
- Symbols are normalized to uppercase ASCII.
- The pair `(market, symbol)` identifies a listed security in user-facing APIs.

## Error Response

```json
{
  "error": {
    "code": "validation_error",
    "message": "Invalid request.",
    "request_id": "req_..."
  }
}
```

Initial error codes:

- `validation_error`
- `not_found`
- `rate_limited`
- `data_source_error`
- `agent_error`
- `insufficient_data`
- `conflict`
- `internal_error`

## Audit Fields

Models that represent user actions or automated decisions should include:

- `created_at`
- `updated_at`
- `created_by`
- `source`
- `data_time`
- `request_id`

## Alert Language

Alerts use observation language:

- `buy_watch`
- `sell_watch`
- `hold_watch`
- `take_profit_watch`
- `stop_loss_watch`
- `risk_warning`
- `abnormal_movement`
- `data_issue`

No API or UI copy may instruct automatic trade execution.
