export const supportedMarkets = ["US", "HK", "CN"];

export const alertSignals = [
  "buy_watch",
  "sell_watch",
  "hold_watch",
  "take_profit_watch",
  "stop_loss_watch",
  "risk_warning",
  "abnormal_movement",
  "data_issue",
];

export function normalizeSymbol(symbol) {
  return String(symbol || "").trim().toUpperCase();
}
