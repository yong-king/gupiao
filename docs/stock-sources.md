# Stock Sources

Preferred order:

1. Official or licensed market data API.
2. Exchange/public issuer pages.
3. Public HTML parsing only when allowed by the site terms.

Crawler rules:

- Do not bypass login, captcha, paywalls, or anti-bot controls.
- Do not scrape broker account pages.
- Keep source URL, source name, request time, and data time.
- Apply rate limits per source.

MVP parser support:

- JSON payloads with `market`, `symbol`, `price`, `previous_close`, `change_percent`, `volume`, `source`, `data_time`.
- Simple HTML elements using `data-*` attributes.
