export const navItems = [
  { id: "watchlists", label: "股票池" },
  { id: "holdings", label: "持仓" },
  { id: "rules", label: "提醒规则" },
  { id: "refresh", label: "刷新任务" },
  { id: "alerts", label: "提醒中心" },
  { id: "reports", label: "分析报告" },
  { id: "accounts", label: "账户监控" },
  { id: "settings", label: "系统设置" },
];

export const viewCopy = {
  watchlists: {
    title: "股票池",
    description: "创建多个自定义股票池，维护需要持续观察的股票代码和关注价。",
    empty: "还没有股票池。先新建一个股票池，再加入要监控的代码。",
  },
  holdings: {
    title: "持仓",
    description: "维护真实或手工录入的持仓，用于风险提醒和组合分析。",
    empty: "还没有持仓。可以从股票池选择代码，也可以手动新增。",
  },
  rules: {
    title: "提醒规则",
    description: "配置观察型买入、卖出和风险提醒阈值，不包含自动下单。",
    empty: "还没有提醒规则。请选择持仓或股票池代码创建规则。",
  },
  refresh: {
    title: "刷新任务",
    description: "手动或低频自动刷新行情，避免高频请求导致账号或数据源受限。",
    empty: "暂无刷新任务。",
  },
  alerts: {
    title: "提醒中心",
    description: "查看规则触发后的提醒、原因和风险等级。",
    empty: "等待监控规则触发。",
  },
  reports: {
    title: "分析报告",
    description: "汇总行情、规则触发和组合风险，辅助人工决策。",
    empty: "盘后生成每日复盘。",
  },
  accounts: {
    title: "账户监控",
    description: "管理同花顺、东方财富、雪球等外部工具的只读接入配置。",
    empty: "当前只保存只读接入配置，不保存交易密码，也不做自动交易。",
  },
  settings: {
    title: "系统设置",
    description: "查看配置文件位置、刷新策略和外部服务连接方式。",
    empty: "后端、Agent、大模型、数据库和 Redis 配置见项目配置文件。",
  },
};

export function isKnownView(viewID) {
  return navItems.some((item) => item.id === viewID);
}

export function getViewCopy(viewID) {
  return viewCopy[isKnownView(viewID) ? viewID : navItems[0].id];
}

export function layoutClassForAuthState(hasToken) {
  return hasToken ? "app-shell" : "auth-layout";
}

export const refreshModes = [
  { id: "manual", label: "手动模式" },
  { id: "conservative", label: "保守模式" },
  { id: "standard", label: "标准模式" },
];

export function formatRefreshStatus(job) {
  if (!job) return "暂无刷新任务";
  if (job.status === "rate_limited") return `刷新冷却中：${job.error || "请稍后再试"}`;
  if (job.status === "failed") return `刷新失败：${job.error || "未知错误"}`;
  if (job.status === "succeeded") return "刷新成功，数据已更新";
  return `刷新状态：${job.status}`;
}

export function formatAlertMessage(message) {
  return {
    title: message?.title || `${message?.signal || "alert"} ${message?.market || ""}:${message?.symbol || ""}`,
    meta: `${message?.riskLevel || message?.risk_level || "low"} · ${message?.dataTime || message?.data_time || ""}`,
    summary: message?.summary || "",
  };
}

export function hasTradingControls(items = navItems) {
  return items.some((item) => /下单|买入下单|卖出下单|交易/.test(item.label));
}
