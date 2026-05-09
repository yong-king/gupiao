export function valueOf(item, ...keys) {
  for (const key of keys) {
    if (item?.[key] !== undefined && item?.[key] !== null) return item[key];
  }
  return undefined;
}

export function formatDailyChange(change) {
  const date = valueOf(change, "Date", "date") || "";
  const close = Number(valueOf(change, "Close", "close") || 0);
  const diff = Number(valueOf(change, "Change", "change") || 0);
  const percent = Number(valueOf(change, "ChangePercent", "change_percent") || 0);
  const sign = diff > 0 ? "+" : "";
  return `${date} 收盘 ${close.toFixed(2)}，涨跌 ${sign}${diff.toFixed(2)} (${sign}${percent.toFixed(2)}%)`;
}

export function buildSparklinePath(points, width = 520, height = 180) {
  if (!points.length) return "";
  const prices = points.map((item) => Number(valueOf(item, "Price", "price", "Close", "close") || 0));
  const min = Math.min(...prices);
  const max = Math.max(...prices);
  const span = max - min || 1;
  return prices.map((price, index) => {
    const x = points.length === 1 ? width / 2 : (index / (points.length - 1)) * width;
    const y = height - ((price - min) / span) * height;
    return `${index === 0 ? "M" : "L"} ${x.toFixed(1)} ${y.toFixed(1)}`;
  }).join(" ");
}

export function renderPriceChart(points) {
  const path = buildSparklinePath(points);
  if (!path) return "<p>暂无曲线数据，请先刷新真实行情。</p>";
  const latest = points[points.length - 1] || {};
  const latestPrice = Number(valueOf(latest, "Price", "price", "Close", "close") || 0);
  return `
    <svg class="price-chart" viewBox="0 0 520 180" role="img" aria-label="价格曲线">
      <path class="chart-grid" d="M 0 45 H 520 M 0 90 H 520 M 0 135 H 520"></path>
      <path class="chart-line" d="${path}"></path>
      ${points.length === 1 ? `<circle class="chart-point" cx="260" cy="90" r="5"></circle>` : ""}
    </svg>
    <p class="chart-value">最新价 ${latestPrice.toFixed(2)}，样本 ${points.length} 条</p>
  `;
}

export function renderRealtimeQuote(snapshot) {
  if (!snapshot) return "<p>暂无实时行情，请点击刷新。</p>";
  const price = Number(valueOf(snapshot, "Price", "price") || 0);
  const open = Number(valueOf(snapshot, "Open", "open") || 0);
  const high = Number(valueOf(snapshot, "High", "high") || 0);
  const low = Number(valueOf(snapshot, "Low", "low") || 0);
  const previousClose = Number(valueOf(snapshot, "PreviousClose", "previous_close") || 0);
  const percent = Number(valueOf(snapshot, "ChangePercent", "change_percent") || 0);
  const volume = Number(valueOf(snapshot, "Volume", "volume") || 0);
  const source = valueOf(snapshot, "Source", "source") || "-";
  const dataTime = valueOf(snapshot, "DataTime", "data_time") || "";
  const sign = percent > 0 ? "+" : "";
  const cls = percent > 0 ? "cn-up" : percent < 0 ? "cn-down" : "muted";
  return `
    <div class="quote-board">
      <div><span>最新价</span><strong class="${cls}">${price.toFixed(2)}</strong></div>
      <div><span>涨跌幅</span><strong class="${cls}">${sign}${percent.toFixed(2)}%</strong></div>
      <div><span>今开</span><strong>${open.toFixed(2)}</strong></div>
      <div><span>最高</span><strong>${high.toFixed(2)}</strong></div>
      <div><span>最低</span><strong>${low.toFixed(2)}</strong></div>
      <div><span>昨收</span><strong>${previousClose.toFixed(2)}</strong></div>
      <div><span>成交量</span><strong>${formatVolume(volume)}</strong></div>
      <div><span>来源</span><strong>${source}</strong></div>
    </div>
    <p class="muted">数据时间：${formatDateTime(dataTime)}</p>
  `;
}

export function renderCandlestickChart(items) {
  if (!items.length) return "<p>暂无K线数据，请先刷新真实行情。</p>";
  const width = 560;
  const height = 220;
  const padding = 18;
  const visible = items.slice(-60);
  const highs = visible.map((item) => Number(valueOf(item, "High", "high") || 0)).filter(Boolean);
  const lows = visible.map((item) => Number(valueOf(item, "Low", "low") || 0)).filter(Boolean);
  if (!highs.length || !lows.length) return "<p>暂无有效K线价格。</p>";
  const max = Math.max(...highs);
  const min = Math.min(...lows);
  const span = max - min || 1;
  const step = (width - padding * 2) / Math.max(visible.length, 1);
  const candleWidth = Math.max(4, Math.min(12, step * 0.55));
  const y = (price) => padding + (height - padding * 2) * (1 - (price - min) / span);
  const candles = visible.map((item, index) => {
    const open = Number(valueOf(item, "Open", "open") || 0);
    const close = Number(valueOf(item, "Close", "close") || valueOf(item, "Price", "price") || 0);
    const high = Number(valueOf(item, "High", "high") || Math.max(open, close));
    const low = Number(valueOf(item, "Low", "low") || Math.min(open, close));
    const x = padding + index * step + step / 2;
    const top = Math.min(y(open), y(close));
    const bodyHeight = Math.max(2, Math.abs(y(open) - y(close)));
    const cls = close >= open ? "up" : "down";
    return `<g class="candle ${cls}">
      <line x1="${x.toFixed(1)}" y1="${y(high).toFixed(1)}" x2="${x.toFixed(1)}" y2="${y(low).toFixed(1)}"></line>
      <rect x="${(x - candleWidth / 2).toFixed(1)}" y="${top.toFixed(1)}" width="${candleWidth.toFixed(1)}" height="${bodyHeight.toFixed(1)}"></rect>
    </g>`;
  }).join("");
  const latest = visible[visible.length - 1];
  const latestClose = Number(valueOf(latest, "Close", "close") || 0);
  const latestDate = valueOf(latest, "Date", "date") || "";
  return `
    <svg class="kline-chart" viewBox="0 0 ${width} ${height}" role="img" aria-label="K线图">
      <path class="chart-grid" d="M 0 55 H ${width} M 0 110 H ${width} M 0 165 H ${width}"></path>
      ${candles}
    </svg>
    <p class="chart-value">K线 ${visible.length} 条，最新 ${latestDate} 收盘 ${latestClose.toFixed(2)}，区间 ${min.toFixed(2)}-${max.toFixed(2)}</p>
  `;
}

export function summarizeMarketNumbers(snapshot, change) {
  if (!snapshot) return "暂无行情数值。";
  const price = Number(valueOf(snapshot, "Price", "price") || 0);
  const open = Number(valueOf(snapshot, "Open", "open") || 0);
  const high = Number(valueOf(snapshot, "High", "high") || 0);
  const low = Number(valueOf(snapshot, "Low", "low") || 0);
  const percent = Number(valueOf(change, "ChangePercent", "change_percent") || valueOf(snapshot, "ChangePercent", "change_percent") || 0);
  const sign = percent > 0 ? "+" : "";
  return `现价 ${price.toFixed(2)}，开 ${open.toFixed(2)}，高 ${high.toFixed(2)}，低 ${low.toFixed(2)}，涨跌幅 ${sign}${percent.toFixed(2)}%`;
}

function formatVolume(value) {
  if (value >= 100000000) return `${(value / 100000000).toFixed(2)}亿`;
  if (value >= 10000) return `${(value / 10000).toFixed(2)}万`;
  return String(Math.round(value || 0));
}

function formatDateTime(value) {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return date.toLocaleString("zh-CN", { hour12: false });
}

export function monitorText(snapshot, symbol) {
  const price = Number(valueOf(snapshot, "Price", "price") || 0);
  const buyPrice = Number(valueOf(symbol, "BuyPrice", "buy_price") || 0);
  const sellPrice = Number(valueOf(symbol, "SellPrice", "sell_price") || 0);
  if (!price) return "等待行情";
  if (buyPrice > 0 && price <= buyPrice) return `到达买入关注价 ${buyPrice.toFixed(2)}`;
  if (sellPrice > 0 && price >= sellPrice) return `到达卖出关注价 ${sellPrice.toFixed(2)}`;
  return `观察中，当前 ${price.toFixed(2)}`;
}

export function renderChangeCalendar(changes) {
  if (!changes.length) return "<p>暂无日历数据，请先保存行情。</p>";
  const cells = changes.slice(-31).map((item) => {
    const date = valueOf(item, "Date", "date") || "";
    const percent = Number(valueOf(item, "ChangePercent", "change_percent") || 0);
    const close = Number(valueOf(item, "Close", "close") || 0);
    const className = percent > 0 ? "up" : percent < 0 ? "down" : "flat";
    const day = String(date).slice(-2) || "--";
    const sign = percent > 0 ? "+" : "";
    return `<div class="calendar-cell ${className}" title="${date} 收盘 ${close.toFixed(2)}"><strong>${day}</strong><span>${sign}${percent.toFixed(2)}%</span></div>`;
  }).join("");
  return `<div class="change-calendar">${cells}</div>`;
}

export function summarizeProfile(profile) {
  const name = valueOf(profile, "name", "Name") || "";
  const sector = valueOf(profile, "sector", "Sector") || "";
  const recommendation = valueOf(profile, "recommendation", "Recommendation") || "observe";
  return `${name} ${sector} ${recommendation}`.trim();
}
