# Redis Key Plan

Redis is used for short-lived coordination and rate control. Persistent business data belongs in PostgreSQL.

## Keys

| Key | Type | TTL | Fields / Value | Purpose |
| --- | --- | --- | --- | --- |
| `session:{token_hash}` | String | session TTL | `user_id` | Fast session lookup without storing plaintext token. |
| `rate:stock_source:{source}:{market}:{symbol}` | String counter | 60s | integer count | Prevent frequent provider calls that could trigger account/source limits. |
| `lock:refresh:{user_id}:{watchlist_id}` | String lock | 5m | job id | Avoid duplicate manual/automatic refresh jobs. |
| `cache:quote:{market}:{symbol}` | Hash | 30-120s | `price`, `open`, `high`, `low`, `volume`, `source`, `data_time` | Low-frequency quote cache for UI refresh. |
| `alert:cooldown:{rule_id}` | String | rule cooldown | last event id | Prevent repeated alert spam. |
| `job:refresh:{job_id}` | Hash | 24h | `status`, `started_at`, `finished_at`, `error` | Runtime status for refresh task polling. |

## Local Check

```bash
docker compose -f deploy/docker-compose.yml exec redis redis-cli KEYS '*'
```
