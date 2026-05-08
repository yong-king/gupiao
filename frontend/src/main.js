import { formatRefreshStatus, getViewCopy, isKnownView, layoutClassForAuthState, navItems, refreshModes } from "./app.js?v=25";
import { API_BASE, deleteJSON, getJSON, postJSON, shouldInvalidateSession } from "./api.js?v=25";
import { formatDailyChange, monitorText, renderChangeCalendar, renderPriceChart, summarizeMarketNumbers, summarizeProfile, valueOf } from "./market.js?v=25";

const root = document.querySelector("#app");
const STOCK_DETAIL_COOLDOWN_MS = 5 * 60 * 1000;
const assistantSessionID = window.localStorage.getItem("jijin_assistant_session") || `chat-${Date.now().toString(36)}`;
window.localStorage.setItem("jijin_assistant_session", assistantSessionID);
const state = {
  token: window.localStorage.getItem("jijin_token") || "",
  email: window.localStorage.getItem("jijin_email") || "",
  userID: window.localStorage.getItem("jijin_user_id") || "user-demo",
  activeView: window.localStorage.getItem("jijin_active_view") || "watchlists",
  selectedMarket: window.localStorage.getItem("jijin_selected_market") || "CN",
  selectedSymbol: window.localStorage.getItem("jijin_selected_symbol") || "000821",
  selectedWatchlistID: window.localStorage.getItem("jijin_selected_watchlist") || "",
  theme: window.localStorage.getItem("jijin_theme") || "dark",
  displayName: window.localStorage.getItem("jijin_display_name") || "",
  message: "",
  alerts: [],
  notifications: [],
  watchlists: [],
  watchlist: null,
  quoteByKey: {},
  profileByKey: {},
  holdings: [],
  rules: [],
  job: null,
  collected: null,
  snapshots: [],
  dailyChanges: [],
  profile: null,
  dependencies: null,
  accountConfigs: [],
  researchResult: null,
  workflows: [],
  workflowResult: null,
  assistantAnswer: null,
  assistantMessages: [],
  assistantSessionID,
  assistantStreaming: false,
  quoteFetchedAt: {},
};

applyTheme();

window.addEventListener("jijin-auth-expired", () => {
  state.token = "";
  state.email = "";
  state.userID = "user-demo";
  state.message = "登录已过期，请重新登录。";
  render();
});

function render() {
  applyTheme();
  if (!state.token) {
    renderAuth();
    return;
  }
  root.className = layoutClassForAuthState(true);
  if (!isKnownView(state.activeView)) {
    state.activeView = "watchlists";
  }
  const nav = navItems
    .map((item) => `<button class="${item.id === state.activeView ? "active" : ""}" data-view="${item.id}" type="button">${item.label}</button>`)
    .join("");
  root.innerHTML = `
    <aside class="sidebar">
      <div class="brand">
        <strong>股票监控</strong>
        <span>人工决策工作台</span>
      </div>
      <nav>${nav}</nav>
    </aside>
    <main class="content">
      <section class="topbar">
        <div>
          <h1>投资监控工作台</h1>
          <p class="muted">当前股票：${state.selectedMarket}:${state.selectedSymbol}</p>
        </div>
        <div class="userbar">
          <span class="muted">${state.displayName || state.email}</span>
          <button id="logout" type="button">退出登录</button>
        </div>
      </section>
      ${renderActiveView()}
      <p class="status-line">${state.message || " "}</p>
    </main>
  `;
  document.querySelectorAll("[data-view]").forEach((button) => {
    button.addEventListener("click", () => {
      state.activeView = button.dataset.view;
      window.localStorage.setItem("jijin_active_view", state.activeView);
      state.message = "";
      render();
    });
  });
  document.querySelector("#logout").addEventListener("click", () => {
    ["jijin_token", "jijin_email", "jijin_user_id"].forEach((key) => window.localStorage.removeItem(key));
    state.token = "";
    state.userID = "user-demo";
    render();
  });
  bindViewActions();
}

function renderActiveView() {
  const view = getViewCopy(state.activeView);
  return `
    <section class="view-header">
      <div>
        <h2>${view.title}</h2>
        <p class="muted">${view.description}</p>
      </div>
      <span class="safety-badge">只提醒，不自动交易</span>
    </section>
    ${renderViewBody(view)}
  `;
}

function renderViewBody(view) {
  if (state.activeView === "watchlists") return renderWatchlistsView(view);
  if (state.activeView === "holdings") return renderHoldingsView(view);
  if (state.activeView === "rules") return renderRulesView(view);
  if (state.activeView === "refresh") return renderRefreshView();
  if (state.activeView === "alerts") return renderAlertsView(view);
  if (state.activeView === "reports") return renderReportsView(view);
  if (state.activeView === "assistant") return renderAssistantView(view);
  if (state.activeView === "accounts") return renderAccountsView();
  return renderSettingsView(view);
}

function renderWatchlistsView(view) {
  const symbols = currentSymbols();
  const pools = state.watchlists.map((item) => {
    const id = itemID(item);
    const active = id === state.selectedWatchlistID ? "active" : "";
    return `<button class="pool-pill ${active}" data-select-watchlist="${id}" type="button">${escapeHTML(valueOf(item, "Name", "name") || id)}</button>`;
  }).join("");
  const symbolRows = symbols.map((item) => {
    const market = valueOf(item, "Market", "market");
    const symbol = valueOf(item, "Symbol", "symbol");
    const key = stockKey(market, symbol);
    const quote = state.quoteByKey[key];
    return `<tr>
      <td><button class="link-button" data-load-stock="${key}" type="button">${market}:${symbol}</button></td>
      <td>${formatPrice(valueOf(item, "BuyPrice", "buy_price"))}</td>
      <td>${formatPrice(valueOf(item, "SellPrice", "sell_price"))}</td>
      <td>${quote ? formatPrice(valueOf(quote, "Price", "price")) : "-"}</td>
      <td>${monitorText(quote, item)}</td>
    </tr>`;
  }).join("");
  const selectedProfile = state.profileByKey[stockKey(state.selectedMarket, state.selectedSymbol)] || state.profile;
  const selectedQuote = state.quoteByKey[stockKey(state.selectedMarket, state.selectedSymbol)] || state.collected;
  return `
    <section class="grid">
      <article>
        <h3>新建股票池</h3>
        <label>股票池名称 <input id="watchlist-name" placeholder="例如：短线观察、长期配置" /></label>
        <button id="create-watchlist" type="button">新建股票池</button>
      </article>
      <article>
        <h3>我的股票池</h3>
        <div class="pill-row">${pools || `<p>${view.empty}</p>`}</div>
        <div class="actions">
          <button id="load-watchlists" type="button">刷新股票池</button>
          <button id="delete-watchlist" type="button">删除当前股票池</button>
        </div>
        <label>目标股票池
          <select id="target-watchlist">
            ${state.watchlists.filter((item) => itemID(item) !== state.selectedWatchlistID).map((item) => `<option value="${itemID(item)}">${escapeHTML(valueOf(item, "Name", "name") || itemID(item))}</option>`).join("")}
          </select>
        </label>
        <label>合并来源股票池
          <select id="merge-watchlist">
            ${state.watchlists.filter((item) => itemID(item) !== state.selectedWatchlistID).map((item) => `<option value="${itemID(item)}">${escapeHTML(valueOf(item, "Name", "name") || itemID(item))}</option>`).join("")}
          </select>
        </label>
        <div class="actions">
          <button id="copy-stock" type="button">复制当前股票到目标池</button>
          <button id="move-stock" type="button">移动当前股票到目标池</button>
          <button id="merge-watchlists" type="button">合并来源池到当前池</button>
        </div>
      </article>
      <article>
        <h3>加入当前股票池</h3>
        <p class="muted">当前池：${currentWatchlistName()}</p>
        ${renderStockPicker("stock", true)}
        <div class="two-col">
          <label>买入关注价 <input id="buy-price" type="number" min="0" step="0.01" placeholder="低于该价格提醒" /></label>
          <label>卖出关注价 <input id="sell-price" type="number" min="0" step="0.01" placeholder="高于该价格提醒" /></label>
        </div>
        <div class="actions">
          <button id="lookup-stock" type="button">获取股票信息</button>
          <button id="save-monitor-stock" type="button">加入/更新监控</button>
        </div>
      </article>
      <article>
        <h3>股票信息</h3>
        ${selectedQuote ? `<p class="market-number">${state.selectedMarket}:${state.selectedSymbol} 当前 ${formatPrice(valueOf(selectedQuote, "Price", "price"))}</p>` : "<p>输入股票代码后点击获取股票信息。</p>"}
        <p>${selectedProfile ? summarizeProfile(selectedProfile) : "暂无公司和产品信息。"}</p>
        <p class="muted">${selectedProfile ? valueOf(selectedProfile, "analysis", "Analysis") : "采集后会展示业务、产品和走势观察。"}</p>
      </article>
      <article class="wide">
        <h3>当前池股票</h3>
        ${symbols.length ? `<table class="stock-table"><thead><tr><th>代码</th><th>买入关注</th><th>卖出关注</th><th>当前价</th><th>状态</th></tr></thead><tbody>${symbolRows}</tbody></table>` : `<p>${view.empty}</p>`}
      </article>
    </section>
  `;
}

function renderHoldingsView(view) {
  const rows = state.holdings.map((item) => {
    const market = valueOf(item, "Market", "market");
    const symbol = valueOf(item, "Symbol", "symbol");
    const quantity = Number(valueOf(item, "Quantity", "quantity") || 0);
    const cost = Number(valueOf(item, "CostBasis", "cost_basis") || 0);
    const quote = state.quoteByKey[stockKey(market, symbol)];
    const latest = Number(valueOf(quote, "Price", "price") || 0);
    const pnl = latest > 0 ? (latest - cost) * quantity : 0;
    const attention = valueOf(item, "AttentionLevel", "attention_level") || "medium";
    return `<tr>
      <td><button class="link-button" data-load-stock="${stockKey(market, symbol)}" type="button">${market}:${symbol}</button></td>
      <td>${quantity}</td>
      <td>${cost.toFixed(2)}</td>
      <td><span class="risk-badge risk-${attentionRisk(attention)}">${attentionLabel(attention)}</span></td>
      <td>${latest ? latest.toFixed(2) : "-"}</td>
      <td class="${pnl >= 0 ? "gain" : "loss"}">${latest ? pnl.toFixed(2) : "-"}</td>
      <td class="row-actions">
        <button data-edit-holding="${stockKey(market, symbol)}" type="button">修改</button>
        <button data-delete-holding="${stockKey(market, symbol)}" type="button">删除</button>
      </td>
    </tr>`;
  }).join("");
  return `
    <section class="grid">
      <article>
        <h3>持仓配置</h3>
        <label>从股票池选择
          <select id="holding-source">
            <option value="">手动输入</option>
            ${stockOptions(allPoolSymbols())}
          </select>
        </label>
        ${renderStockPicker("holding", false)}
        <div class="two-col">
          <label>数量 <input id="holding-quantity" type="number" min="0" step="0.0001" placeholder="例如 1000" /></label>
          <label>成本价 <input id="holding-cost" type="number" min="0" step="0.01" placeholder="例如 8.20" /></label>
        </div>
        <label>关注等级
          <select id="holding-attention">
            <option value="high">高：4 小时采集一次</option>
            <option value="medium" selected>中：6 小时采集一次</option>
            <option value="low">低：1 天采集一次</option>
          </select>
        </label>
        <div class="actions">
          <button id="save-holding" type="button">保存/更新持仓</button>
          <button id="load-holdings" type="button">刷新持仓</button>
          <button id="collect-holding-research" type="button">采集并总结当前持仓信息</button>
        </div>
      </article>
      <article>
        <h3>持仓概览</h3>
        <p>${state.holdings.length ? `当前 ${state.holdings.length} 条持仓。保存后会自动获取行情并加入分析目标。` : view.empty}</p>
        <button id="analyze-holdings" type="button">刷新持仓行情并分析</button>
      </article>
      <article class="wide">
        <h3>当前持仓</h3>
          ${state.holdings.length ? `<table class="stock-table"><thead><tr><th>代码</th><th>数量</th><th>成本价</th><th>关注等级</th><th>最新价</th><th>浮动盈亏</th><th>操作</th></tr></thead><tbody>${rows}</tbody></table>` : "<p>暂无持仓，先保存一条持仓。</p>"}
      </article>
    </section>
  `;
}

function renderRulesView(view) {
  const rules = state.rules.map((item) => {
    const risk = valueOf(item, "RiskLevel", "risk_level") || "low";
    return `<tr>
      <td>${valueOf(item, "Market", "market")}:${valueOf(item, "Symbol", "symbol")}</td>
      <td>${ruleTypeLabel(valueOf(item, "Type", "type"))}</td>
      <td>${formatPrice(valueOf(item, "Threshold", "threshold"))}</td>
      <td>${signalLabel(valueOf(item, "Signal", "signal"))}</td>
      <td><span class="risk-badge risk-${risk}">${riskLabel(risk)}</span></td>
      <td>${valueOf(item, "Enabled", "enabled") === false ? "停用" : "启用"}</td>
      <td><button data-delete-rule="${valueOf(item, "ID", "id")}" type="button">删除</button></td>
    </tr>`;
  }).join("");
  return `
    <section class="grid">
      <article>
        <h3>提醒规则配置</h3>
        <label>规则股票
          <select id="rule-source">
            <option value="">手动输入</option>
            ${stockOptions(analysisTargets())}
          </select>
        </label>
        ${renderStockPicker("rule", false)}
        <label>触发条件
          <select id="rule-type">
            <option value="price_below">价格低于</option>
            <option value="price_above">价格高于</option>
            <option value="change_percent_below">跌幅低于</option>
            <option value="change_percent_above">涨幅高于</option>
          </select>
        </label>
        <div class="two-col">
          <label>阈值 <input id="rule-threshold" type="number" step="0.01" placeholder="例如 8.00 或 -3" /></label>
          <label>提醒信号
            <select id="rule-signal">
              <option value="buy_watch">买入关注</option>
              <option value="sell_watch">卖出关注</option>
              <option value="risk_warning">风险提示</option>
              <option value="hold_watch">继续观察</option>
            </select>
          </label>
        </div>
        <label>提醒等级
          <select id="rule-risk">
            <option value="low">低</option>
            <option value="medium" selected>中</option>
            <option value="high">高</option>
            <option value="critical">严重</option>
          </select>
        </label>
        <div class="actions">
          <button id="save-rule" type="button">保存提醒规则</button>
          <button id="load-rules" type="button">刷新规则</button>
        </div>
      </article>
      <article>
        <h3>规则说明</h3>
        <p class="muted">提醒规则只写入日志和通知中心，不会下单。等级越高，提醒中心颜色越醒目。</p>
      </article>
      <article class="wide">
        <h3>当前规则</h3>
        ${state.rules.length ? `<table class="stock-table"><thead><tr><th>代码</th><th>条件</th><th>阈值</th><th>信号</th><th>等级</th><th>状态</th><th>操作</th></tr></thead><tbody>${rules}</tbody></table>` : `<p>${view.empty}</p>`}
      </article>
    </section>
  `;
}

function renderRefreshView() {
  const modes = refreshModes.map((mode) => `<option value="${mode.id}">${mode.label}</option>`).join("");
  return `
    <section class="grid">
      <article>
        <h3>刷新对象</h3>
        <p>当前股票池：${currentWatchlistName()}</p>
        <p>${formatRefreshStatus(state.job)}</p>
        <div class="actions">
          <button id="refresh-current-pool" type="button">刷新当前股票池</button>
          <button id="refresh-holdings" type="button">刷新全部持仓</button>
        </div>
      </article>
      <article>
        <h3>刷新模式</h3>
        <select id="refresh-mode">${modes}</select>
        <p class="muted">后端默认自动刷新间隔为 30 分钟；手动刷新用于立即验证股票池或持仓数据。</p>
      </article>
      <article>
        <h3>最近保存行情</h3>
        <p>${state.collected ? `${valueOf(state.collected, "Name", "name") || `${state.selectedMarket}:${state.selectedSymbol}`} 最新 ${formatPrice(valueOf(state.collected, "Price", "price"))}` : "暂无保存记录。"}</p>
      </article>
    </section>
  `;
}

function renderAlertsView(view) {
  const alerts = state.alerts.map((item) => {
    const risk = valueOf(item, "RiskLevel", "risk_level") || "low";
    return `<li class="alert-item risk-card-${risk}">
      <strong>${signalLabel(valueOf(item, "Signal", "signal"))}</strong>
      <span>${valueOf(item, "Market", "market")}:${valueOf(item, "Symbol", "symbol")}</span>
      <span class="risk-badge risk-${risk}">${riskLabel(risk)}</span>
      <p>${valueOf(item, "Summary", "summary")}</p>
    </li>`;
  }).join("");
  const notifications = state.notifications.map((item) => {
    const id = valueOf(item, "ID", "id");
    const risk = valueOf(item, "RiskLevel", "risk_level") || "low";
    return `<li class="alert-item risk-card-${risk}">
      <strong>${valueOf(item, "Title", "title")}</strong>
      <span>${valueOf(item, "Read", "read") ? "已读" : "未读"}</span>
      <p>${valueOf(item, "Summary", "summary")}</p>
      <button data-read-notification="${id}" type="button">标记已读</button>
    </li>`;
  }).join("");
  return `
    <section class="grid">
      <article>
        <h3>提醒检查</h3>
        <p>${state.job ? formatRefreshStatus(state.job) : "运行后会刷新当前股票池并根据规则写入提醒日志。"}</p>
        <div class="actions">
          <button id="run-alert-check" type="button">运行当前池监控检查</button>
          <button id="load-alerts" type="button">刷新提醒中心</button>
        </div>
      </article>
      <article>
        <h3>提醒日志</h3>
        ${state.alerts.length ? `<ul class="alert-list">${alerts}</ul>` : `<p>${view.empty}</p>`}
      </article>
      <article class="wide">
        <h3>通知记录</h3>
        ${state.notifications.length ? `<ul class="alert-list">${notifications}</ul>` : "<p>暂无通知。规则触发后会出现在这里。</p>"}
      </article>
    </section>
  `;
}

function renderReportsView(view) {
  const selectedPoolSymbols = currentSymbols().map((item) => ({ market: valueOf(item, "Market", "market"), symbol: valueOf(item, "Symbol", "symbol"), source: "股票池" }));
  const holdingRows = renderAnalysisRows(state.holdings.map((item) => ({
    market: valueOf(item, "Market", "market"),
    symbol: valueOf(item, "Symbol", "symbol"),
    source: "持仓",
    holding: item,
  })));
  const poolRows = renderAnalysisRows(selectedPoolSymbols);
  const changes = state.dailyChanges.map((item) => `<li>${formatDailyChange(item)}</li>`).join("");
  const latestSnapshot = state.snapshots[state.snapshots.length - 1] || state.collected;
  const latestChange = state.dailyChanges[state.dailyChanges.length - 1];
  const holding = findHolding(state.selectedMarket, state.selectedSymbol);
  return `
    <section class="grid">
      <article class="wide">
        <h3>持仓股票今日涨跌</h3>
        ${holdingRows || `<p>${view.empty}</p>`}
        <button id="analyze-holdings" type="button">刷新持仓今日涨跌</button>
      </article>
      <article class="wide">
        <h3>股票池分析</h3>
        <p class="muted">当前股票池：${currentWatchlistName()}</p>
        <div class="pill-row">${state.watchlists.map((item) => `<button class="pool-pill ${itemID(item) === state.selectedWatchlistID ? "active" : ""}" data-select-watchlist="${itemID(item)}" type="button">${escapeHTML(valueOf(item, "Name", "name") || itemID(item))}</button>`).join("")}</div>
        ${poolRows || "<p>当前股票池还没有股票。</p>"}
        <div class="inline-form">
          ${renderStockPicker("analysis", false)}
          <button id="collect-report-market" type="button">采集并分析当前股票</button>
          <button id="load-market-analysis" type="button">加载已保存分析</button>
          <button id="collect-research" type="button">采集产品信息并生成 RAG 总结</button>
        </div>
      </article>
      <article>
        <h3>个人持仓情况</h3>
        <p>${holding ? holdingSummary(holding) : "当前选中股票不在持仓中，可继续作为自选股票分析。"}</p>
      </article>
      <article>
        <h3>关注等级 AI 工作流</h3>
        <p class="muted">按持仓关注等级批量运行：高 4 小时、中 6 小时、低 1 天。每次会记录 Agent 步骤，写入 RAG 和本地向量。</p>
        <div class="actions">
          <button data-run-workflow="high" type="button">运行高关注</button>
          <button data-run-workflow="medium" type="button">运行中关注</button>
          <button data-run-workflow="low" type="button">运行低关注</button>
          <button id="load-workflows" type="button">刷新工作流记录</button>
        </div>
        ${state.workflowResult ? `<p class="research-summary">${valueOf(state.workflowResult, "job")?.summary || valueOf(state.workflowResult, "Job")?.Summary || "工作流已执行。"}</p>` : ""}
      </article>
      <article>
        <h3>详细曲线</h3>
        <p class="market-number">${summarizeMarketNumbers(latestSnapshot, latestChange)}</p>
        ${renderPriceChart(state.snapshots)}
      </article>
      <article>
        <h3>涨跌日历</h3>
        ${renderChangeCalendar(state.dailyChanges)}
      </article>
      <article>
        <h3>每日涨跌记录</h3>
        ${state.dailyChanges.length ? `<ul>${changes}</ul>` : "<p>暂无每日记录，请先采集行情。</p>"}
        <p class="muted">每日记录已包含 RAG 文本字段，可用于后续问答分析。</p>
      </article>
      <article>
        <h3>公司和产品情况</h3>
        <p>${state.profile ? summarizeProfile(state.profile) : "尚未加载公司/产品分析。"}</p>
        <p>${state.profile ? valueOf(state.profile, "analysis", "Analysis") : "点击采集后展示业务、产品和趋势归纳。"}</p>
        ${state.researchResult ? `<p class="research-summary">${state.researchResult.summary || state.researchResult.Summary}</p>` : ""}
        <p class="muted">${state.profile ? valueOf(state.profile, "disclaimer", "Disclaimer") : "仅用于监控和研究，不构成投资建议。"}</p>
      </article>
    </section>
  `;
}

function renderAssistantView(view) {
  const targets = analysisTargets();
  const messages = state.assistantMessages.map((message) => `
    <div class="chat-message ${message.role}">
      <span>${message.role === "user" ? "你" : "AI"}</span>
      <p>${escapeHTML(message.content || "")}</p>
    </div>
  `).join("");
  const workflowRows = state.workflows.map((job) => {
    const steps = asArray(valueOf(job, "Steps", "steps"));
    return `<tr>
      <td>${valueOf(job, "AttentionLevel", "attention_level")}</td>
      <td>${valueOf(job, "Status", "status")}</td>
      <td>${valueOf(job, "TargetCount", "target_count")}</td>
      <td>${escapeHTML(valueOf(job, "Summary", "summary") || "")}</td>
      <td>${steps.length}</td>
    </tr>`;
  }).join("");
  return `
    <section class="assistant-layout">
      <article class="assistant-config">
        <div class="panel-title">
          <h3>股票助手</h3>
          <span class="safety-badge">多轮上下文</span>
        </div>
        <label>从已有股票选择
          <select id="assistant-source">
            <option value="">手动输入</option>
            ${stockOptions(targets)}
          </select>
        </label>
        ${renderStockPicker("assistant", false)}
        <p class="muted">会话：${escapeHTML(state.assistantSessionID)}。回答会结合持仓、股票池、行情、产品信息和 RAG 历史。</p>
        <div class="actions">
          <button id="load-workflows-assistant" type="button">刷新工作流记录</button>
          <button id="clear-assistant-chat" type="button">清空本页对话</button>
        </div>
      </article>
      <article class="chat-panel">
        <div class="panel-title">
          <h3>${state.selectedMarket}:${state.selectedSymbol}</h3>
          <span class="muted">${state.assistantStreaming ? "正在流式输出..." : "DeepSeek 路由分析"}</span>
        </div>
        <div class="chat-window" id="assistant-chat-window">
          ${messages || `<div class="empty-chat">${view.empty}</div>`}
        </div>
        <label class="chat-composer">你的问题
          <textarea id="assistant-question" rows="3" placeholder="例如：结合我的持仓成本、最近产品信息和历史 RAG，总结这只股票下一步重点关注什么。"></textarea>
        </label>
        <div class="actions">
          <button id="ask-assistant" type="button" ${state.assistantStreaming ? "disabled" : ""}>发送并分析</button>
        </div>
        ${state.assistantAnswer ? `<p class="muted">${escapeHTML(valueOf(state.assistantAnswer, "ContextSummary", "context_summary") || "")}</p>` : ""}
      </article>
      <article class="wide">
        <h3>最近 AI 工作流</h3>
        ${state.workflows.length ? `<table class="stock-table"><thead><tr><th>关注等级</th><th>状态</th><th>股票数</th><th>摘要</th><th>步骤数</th></tr></thead><tbody>${workflowRows}</tbody></table>` : "<p>暂无工作流记录。</p>"}
      </article>
    </section>
  `;
}

function renderAccountsView() {
  const accounts = state.accountConfigs.map((item) => {
    const metadata = valueOf(item, "Metadata", "metadata") || {};
    return `<tr>
      <td>${valueOf(item, "ID", "id")}</td>
      <td>${valueOf(item, "Alias", "alias")}</td>
      <td>${providerLabel(metadata.provider)}</td>
      <td>${modeLabel(valueOf(item, "RefreshMode", "refresh_mode", "refreshMode"))}</td>
      <td>${valueOf(item, "ReadOnly", "read_only", "readOnly") ? "只读" : "未限制"}</td>
    </tr>`;
  }).join("");
  return `
    <section class="grid">
      <article>
        <h3>外部工具只读接入</h3>
        <label>接入来源
          <select id="account-provider">
            <option value="ths">同花顺</option>
            <option value="eastmoney">东方财富</option>
            <option value="xueqiu">雪球</option>
            <option value="csv">CSV/手工导入</option>
            <option value="other">其他只读来源</option>
          </select>
        </label>
        <label>账户 ID <input id="account-id" placeholder="例如 ths-readonly-1" /></label>
        <label>账户别名 <input id="account-alias" placeholder="例如 我的同花顺只读账户" /></label>
        <label>刷新模式
          <select id="account-refresh-mode">
            <option value="manual">手动</option>
            <option value="conservative">保守自动</option>
            <option value="standard">标准自动</option>
          </select>
        </label>
        <label class="inline-check"><input id="account-readonly" type="checkbox" checked /> 只读接入</label>
        <div class="actions">
          <button id="save-account" type="button">保存账户配置</button>
          <button id="load-accounts" type="button">刷新账户配置</button>
        </div>
      </article>
      <article>
        <h3>安全边界</h3>
        <p class="muted">这里只保存来源、别名、刷新模式等只读配置；不保存交易密码、token、secret，也不会执行买入或卖出。</p>
      </article>
      <article class="wide">
        <h3>已配置账户</h3>
        ${state.accountConfigs.length ? `<table class="stock-table"><thead><tr><th>ID</th><th>别名</th><th>来源</th><th>刷新模式</th><th>权限</th></tr></thead><tbody>${accounts}</tbody></table>` : "<p>暂无账户配置。</p>"}
      </article>
    </section>
  `;
}

function renderSettingsView(view) {
  return `
    <section class="grid">
      <article>
        <h3>配置文件</h3>
        <p class="code">config/backend.example.json</p>
        <p class="code">agent/config/agent.example.json</p>
        <p class="code">deploy/docker-compose.yml</p>
        <p class="muted">${view.empty}</p>
      </article>
      <article>
        <h3>系统检测</h3>
        <button id="load-dependencies" type="button">测试后端、数据库、Redis、大模型和行情源</button>
        ${renderDependencies()}
      </article>
      <article>
        <h3>主题设置</h3>
        <label>显示主题
          <select id="theme-select">
            <option value="dark" ${state.theme === "dark" ? "selected" : ""}>深色专业</option>
            <option value="light" ${state.theme === "light" ? "selected" : ""}>浅色清爽</option>
          </select>
        </label>
      </article>
      <article>
        <h3>用户信息</h3>
        <p>邮箱：${state.email}</p>
        <p>用户 ID：${state.userID}</p>
        <label>显示名称 <input id="display-name" value="${escapeHTML(state.displayName)}" placeholder="用于页面右上角展示" /></label>
        <button id="save-user-settings" type="button">保存用户信息</button>
      </article>
    </section>
  `;
}

function bindViewActions() {
  document.querySelector("#create-watchlist")?.addEventListener("click", createWatchlistFromForm);
  document.querySelector("#delete-watchlist")?.addEventListener("click", deleteCurrentWatchlist);
  document.querySelector("#copy-stock")?.addEventListener("click", () => transferCurrentStock(false));
  document.querySelector("#move-stock")?.addEventListener("click", () => transferCurrentStock(true));
  document.querySelector("#merge-watchlists")?.addEventListener("click", mergeWatchlists);
  document.querySelectorAll("[data-select-watchlist]").forEach((button) => button.addEventListener("click", () => selectWatchlist(button.dataset.selectWatchlist)));
  document.querySelector("#lookup-stock")?.addEventListener("click", collectStockInfoFromForm);
  document.querySelector("#save-monitor-stock")?.addEventListener("click", saveMonitorStock);
  document.querySelectorAll("[data-load-stock]").forEach((button) => button.addEventListener("click", () => {
    const [market, symbol] = button.dataset.loadStock.split(":");
    collectStockInfo(market, symbol);
  }));
  document.querySelector("#load-watchlists")?.addEventListener("click", () => loadWatchlists(true));
  document.querySelector("#holding-source")?.addEventListener("change", fillHoldingFromSource);
  document.querySelector("#save-holding")?.addEventListener("click", saveHolding);
  document.querySelector("#load-holdings")?.addEventListener("click", () => loadHoldings(true));
  document.querySelector("#collect-holding-research")?.addEventListener("click", collectResearchForCurrentStock);
  document.querySelectorAll("[data-edit-holding]").forEach((button) => button.addEventListener("click", () => fillHoldingForm(button.dataset.editHolding)));
  document.querySelectorAll("[data-delete-holding]").forEach((button) => button.addEventListener("click", () => deleteHolding(button.dataset.deleteHolding)));
  document.querySelector("#rule-source")?.addEventListener("change", fillRuleFromSource);
  document.querySelector("#save-rule")?.addEventListener("click", saveRuleFromForm);
  document.querySelector("#load-rules")?.addEventListener("click", () => loadRules(true));
  document.querySelectorAll("[data-delete-rule]").forEach((button) => button.addEventListener("click", () => deleteRule(button.dataset.deleteRule)));
  document.querySelector("#refresh-current-pool")?.addEventListener("click", refreshCurrentPool);
  document.querySelector("#refresh-holdings")?.addEventListener("click", refreshHoldings);
  document.querySelector("#run-alert-check")?.addEventListener("click", runAlertCheck);
  document.querySelector("#load-alerts")?.addEventListener("click", () => loadAlerts(true));
  document.querySelectorAll("[data-read-notification]").forEach((button) => button.addEventListener("click", () => markNotificationRead(button.dataset.readNotification)));
  document.querySelector("#collect-report-market")?.addEventListener("click", collectReportMarket);
  document.querySelector("#load-market-analysis")?.addEventListener("click", loadMarketAnalysis);
  document.querySelector("#collect-research")?.addEventListener("click", collectResearchForCurrentStock);
  document.querySelectorAll("[data-run-workflow]").forEach((button) => button.addEventListener("click", () => runAttentionWorkflow(button.dataset.runWorkflow)));
  document.querySelector("#load-workflows")?.addEventListener("click", () => loadWorkflows(true));
  document.querySelector("#analyze-holdings")?.addEventListener("click", refreshHoldings);
  document.querySelector("#assistant-source")?.addEventListener("change", fillAssistantFromSource);
  document.querySelector("#ask-assistant")?.addEventListener("click", askAssistant);
  document.querySelector("#clear-assistant-chat")?.addEventListener("click", () => {
    state.assistantMessages = [];
    state.assistantAnswer = null;
    render();
  });
  document.querySelector("#load-workflows-assistant")?.addEventListener("click", () => loadWorkflows(true));
  document.querySelector("#save-account")?.addEventListener("click", saveAccountConfig);
  document.querySelector("#load-accounts")?.addEventListener("click", () => loadAccounts(true));
  document.querySelector("#load-dependencies")?.addEventListener("click", loadDependencies);
  document.querySelector("#theme-select")?.addEventListener("change", (event) => {
    state.theme = event.target.value;
    window.localStorage.setItem("jijin_theme", state.theme);
    render();
  });
  document.querySelector("#save-user-settings")?.addEventListener("click", saveLocalUserSettings);
}

function renderAuth() {
  root.className = layoutClassForAuthState(false);
  root.innerHTML = `
    <main class="auth">
      <section class="auth-card">
        <h1>股票监控操作台</h1>
        <p class="muted">请先注册或登录。密码在后端加盐哈希保存，不保存明文。</p>
        <label>邮箱 <input id="email" value="${escapeHTML(state.email || "demo@example.com")}" /></label>
        <label>密码 <input id="password" type="password" value="password123" /></label>
        <div class="actions">
          <button id="register">注册并登录</button>
          <button id="login">登录</button>
        </div>
        <p class="muted">${state.message}</p>
      </section>
    </main>
  `;
  document.querySelector("#register").addEventListener("click", () => authenticate("register"));
  document.querySelector("#login").addEventListener("click", () => authenticate("login"));
}

async function authenticate(mode) {
  const email = document.querySelector("#email").value;
  const password = document.querySelector("#password").value;
  try {
    const payload = mode === "register" ? { id: `user-${Date.now().toString(36)}`, email, password } : { email, password };
    const data = await postJSON(`/api/auth/${mode}`, payload);
    state.token = data.token;
    state.email = data.email;
    state.userID = data.user_id || data.userID || state.userID;
    window.localStorage.setItem("jijin_token", state.token);
    window.localStorage.setItem("jijin_email", state.email);
    window.localStorage.setItem("jijin_user_id", state.userID);
    state.message = "登录成功，正在加载你的数据。";
    await bootstrapUserData();
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function bootstrapUserData() {
  await Promise.allSettled([loadWatchlists(false), loadHoldings(false), loadRules(false), loadAlerts(false), loadAccounts(false), loadWorkflows(false)]);
  await ensureDefaultWatchlist();
}

async function createWatchlistFromForm() {
  const name = document.querySelector("#watchlist-name")?.value?.trim() || "";
  if (!name) {
    state.message = "请填写股票池名称。";
    render();
    return;
  }
  try {
    const id = `wl-${state.userID}-${Date.now().toString(36)}`.replace(/[^A-Za-z0-9_-]/g, "-");
    await postJSON("/api/watchlists", { id, user_id: state.userID, name }, state.token);
    state.selectedWatchlistID = id;
    window.localStorage.setItem("jijin_selected_watchlist", id);
    await loadWatchlists(false);
    state.watchlist = await getJSON(`/api/watchlists/${encodeURIComponent(id)}`, state.token);
    state.message = `股票池「${name}」已创建。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function selectWatchlist(id) {
  try {
    state.selectedWatchlistID = id;
    window.localStorage.setItem("jijin_selected_watchlist", id);
    state.watchlist = await getJSON(`/api/watchlists/${encodeURIComponent(id)}`, state.token);
    state.message = `已切换到股票池「${currentWatchlistName()}」。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function deleteCurrentWatchlist() {
  if (!state.selectedWatchlistID) {
    state.message = "请先选择要删除的股票池。";
    render();
    return;
  }
  try {
    await deleteJSON(`/api/watchlists/${encodeURIComponent(state.selectedWatchlistID)}?user_id=${encodeURIComponent(state.userID)}`, state.token);
    state.selectedWatchlistID = "";
    window.localStorage.removeItem("jijin_selected_watchlist");
    await loadWatchlists(false);
    await ensureDefaultWatchlist();
    state.message = "股票池已删除。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function transferCurrentStock(move) {
  const targetID = document.querySelector("#target-watchlist")?.value || "";
  if (!targetID) {
    state.message = "请先选择目标股票池。";
    render();
    return;
  }
  try {
    const symbol = currentSymbols().find((item) => stockKey(valueOf(item, "Market", "market"), valueOf(item, "Symbol", "symbol")) === stockKey(state.selectedMarket, state.selectedSymbol)) || {
      market: state.selectedMarket,
      symbol: state.selectedSymbol,
    };
    await postJSON(`/api/watchlists/${encodeURIComponent(targetID)}/symbols`, {
      market: valueOf(symbol, "Market", "market"),
      symbol: valueOf(symbol, "Symbol", "symbol"),
      buy_price: valueOf(symbol, "BuyPrice", "buy_price") || 0,
      sell_price: valueOf(symbol, "SellPrice", "sell_price") || 0,
    }, state.token);
    if (move) {
      await deleteJSON(`/api/watchlists/${encodeURIComponent(state.selectedWatchlistID)}/symbols?market=${encodeURIComponent(state.selectedMarket)}&symbol=${encodeURIComponent(state.selectedSymbol)}`, state.token).catch(() => {});
    }
    await loadWatchlists(false);
    state.watchlist = await getJSON(`/api/watchlists/${encodeURIComponent(state.selectedWatchlistID)}`, state.token).catch(() => state.watchlist);
    state.message = `${state.selectedMarket}:${state.selectedSymbol} 已${move ? "移动" : "复制"}到目标股票池。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function mergeWatchlists() {
  const sourceID = document.querySelector("#merge-watchlist")?.value || "";
  if (!sourceID || !state.selectedWatchlistID) {
    state.message = "请选择要合并的来源股票池。";
    render();
    return;
  }
  try {
    const source = await getJSON(`/api/watchlists/${encodeURIComponent(sourceID)}`, state.token);
    const symbols = source.Symbols || source.symbols || [];
    for (const item of symbols) {
      await postJSON(`/api/watchlists/${encodeURIComponent(state.selectedWatchlistID)}/symbols`, {
        market: valueOf(item, "Market", "market"),
        symbol: valueOf(item, "Symbol", "symbol"),
        buy_price: valueOf(item, "BuyPrice", "buy_price") || 0,
        sell_price: valueOf(item, "SellPrice", "sell_price") || 0,
      }, state.token);
    }
    await loadWatchlists(false);
    state.watchlist = await getJSON(`/api/watchlists/${encodeURIComponent(state.selectedWatchlistID)}`, state.token);
    state.message = `已合并 ${symbols.length} 只股票到当前股票池。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function ensureDefaultWatchlist() {
  if (state.watchlists.length && state.selectedWatchlistID) {
    state.watchlist = await getJSON(`/api/watchlists/${encodeURIComponent(state.selectedWatchlistID)}`, state.token).catch(() => state.watchlist);
    return;
  }
  if (state.watchlists.length) {
    state.selectedWatchlistID = itemID(state.watchlists[0]);
    window.localStorage.setItem("jijin_selected_watchlist", state.selectedWatchlistID);
    state.watchlist = await getJSON(`/api/watchlists/${encodeURIComponent(state.selectedWatchlistID)}`, state.token).catch(() => state.watchlists[0]);
    return;
  }
  const id = `wl-${state.userID}-default`.replace(/[^A-Za-z0-9_-]/g, "-");
  await postJSON("/api/watchlists", { id, user_id: state.userID, name: "我的股票池" }, state.token).catch(() => {});
  state.selectedWatchlistID = id;
  window.localStorage.setItem("jijin_selected_watchlist", id);
  await loadWatchlists(false);
  state.watchlist = await getJSON(`/api/watchlists/${encodeURIComponent(id)}`, state.token).catch(() => null);
}

async function loadWatchlists(showMessage) {
  try {
    state.watchlists = asArray(await getJSON(`/api/watchlists?user_id=${encodeURIComponent(state.userID)}`, state.token));
    if (state.selectedWatchlistID) {
      state.watchlist = await getJSON(`/api/watchlists/${encodeURIComponent(state.selectedWatchlistID)}`, state.token).catch(() => state.watchlist);
    }
    if (showMessage) state.message = `已加载 ${state.watchlists.length} 个股票池。`;
  } catch (error) {
    state.message = error.message;
  }
  if (showMessage) render();
}

async function saveMonitorStock() {
  const form = readStockForm("stock");
  if (!form) return;
  try {
    await ensureDefaultWatchlist();
    await collectStockInfoData(form.market, form.symbol, { force: true });
    state.watchlist = await postJSON(`/api/watchlists/${encodeURIComponent(state.selectedWatchlistID)}/symbols`, {
      market: form.market,
      symbol: form.symbol,
      buy_price: form.buyPrice,
      sell_price: form.sellPrice,
    }, state.token);
    await createPriceRules(form);
    await loadWatchlists(false);
    await loadRules(false);
    state.message = `${form.market}:${form.symbol} 已加入「${currentWatchlistName()}」。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function createPriceRules(form) {
  const base = `${state.userID}-${form.market}-${form.symbol}`.replace(/[^A-Za-z0-9_-]/g, "-");
  if (form.buyPrice > 0) {
    await postJSON("/api/alert-rules", {
      id: `rule-${base}-buy`,
      user_id: state.userID,
      market: form.market,
      symbol: form.symbol,
      type: "price_below",
      threshold: form.buyPrice,
      signal: "buy_watch",
      risk_level: "medium",
      enabled: true,
      cooldown_seconds: 1800,
    }, state.token).catch(() => {});
  }
  if (form.sellPrice > 0) {
    await postJSON("/api/alert-rules", {
      id: `rule-${base}-sell`,
      user_id: state.userID,
      market: form.market,
      symbol: form.symbol,
      type: "price_above",
      threshold: form.sellPrice,
      signal: "sell_watch",
      risk_level: "medium",
      enabled: true,
      cooldown_seconds: 1800,
    }, state.token).catch(() => {});
  }
}

async function loadHoldings(showMessage) {
  try {
    state.holdings = asArray(await getJSON(`/api/holdings?user_id=${encodeURIComponent(state.userID)}`, state.token));
    if (showMessage) state.message = `已加载 ${state.holdings.length} 条持仓。`;
  } catch (error) {
    state.message = error.message;
  }
  if (showMessage) render();
}

async function saveHolding() {
  const market = document.querySelector("#holding-market")?.value || state.selectedMarket;
  const symbol = (document.querySelector("#holding-symbol")?.value || "").trim().toUpperCase();
  const quantity = Number(document.querySelector("#holding-quantity")?.value || 0);
  const costBasis = Number(document.querySelector("#holding-cost")?.value || 0);
  const attentionLevel = document.querySelector("#holding-attention")?.value || "medium";
  if (!symbol || quantity <= 0 || costBasis < 0) {
    state.message = "请填写有效的股票代码、数量和成本价。";
    render();
    return;
  }
  try {
    await postJSON("/api/holdings", { user_id: state.userID, market, symbol, quantity, cost_basis: costBasis, attention_level: attentionLevel }, state.token);
    await ensureDefaultWatchlist();
    await postJSON(`/api/watchlists/${encodeURIComponent(state.selectedWatchlistID)}/symbols`, { market, symbol }, state.token).catch(() => {});
    await collectStockInfoData(market, symbol, { force: true });
    await Promise.all([loadHoldings(false), loadWatchlists(false)]);
    state.message = `${market}:${symbol} 持仓已保存，并已加入当前股票池。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function deleteHolding(key) {
  const [market, symbol] = key.split(":");
  try {
    await deleteJSON(`/api/holdings?user_id=${encodeURIComponent(state.userID)}&market=${encodeURIComponent(market)}&symbol=${encodeURIComponent(symbol)}`, state.token);
    await loadHoldings(false);
    state.message = `${key} 持仓已删除。分析页仍可手动选择该股票继续观察。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

function fillHoldingForm(key) {
  const [market, symbol] = key.split(":");
  const item = state.holdings.find((holding) => stockKey(valueOf(holding, "Market", "market"), valueOf(holding, "Symbol", "symbol")) === key);
  setInput("#holding-market", market);
  setInput("#holding-symbol", symbol);
  setInput("#holding-quantity", valueOf(item, "Quantity", "quantity") || "");
  setInput("#holding-cost", valueOf(item, "CostBasis", "cost_basis") || "");
  setInput("#holding-attention", valueOf(item, "AttentionLevel", "attention_level") || "medium");
}

function fillHoldingFromSource() {
  const key = document.querySelector("#holding-source")?.value || "";
  if (key) fillHoldingFormFromKey(key, "holding");
}

function fillRuleFromSource() {
  const key = document.querySelector("#rule-source")?.value || "";
  if (key) fillHoldingFormFromKey(key, "rule");
}

function fillAssistantFromSource() {
  const key = document.querySelector("#assistant-source")?.value || "";
  if (key) fillHoldingFormFromKey(key, "assistant");
}

function fillHoldingFormFromKey(key, prefix) {
  const [market, symbol] = key.split(":");
  setInput(`#${prefix}-market`, market);
  setInput(`#${prefix}-symbol`, symbol);
}

async function saveRuleFromForm() {
  const market = document.querySelector("#rule-market")?.value || state.selectedMarket;
  const symbol = (document.querySelector("#rule-symbol")?.value || "").trim().toUpperCase();
  const type = document.querySelector("#rule-type")?.value || "price_below";
  const signal = document.querySelector("#rule-signal")?.value || "buy_watch";
  const riskLevel = document.querySelector("#rule-risk")?.value || "medium";
  const threshold = Number(document.querySelector("#rule-threshold")?.value || 0);
  if (!symbol || threshold === 0) {
    state.message = "请填写股票代码和阈值。";
    render();
    return;
  }
  try {
    const id = `rule-${state.userID}-${market}-${symbol}-${type}-${Date.now()}`.replace(/[^A-Za-z0-9_-]/g, "-");
    await postJSON("/api/alert-rules", {
      id,
      user_id: state.userID,
      market,
      symbol,
      type,
      threshold,
      signal,
      risk_level: riskLevel,
      enabled: true,
      cooldown_seconds: 1800,
    }, state.token);
    await loadRules(false);
    state.message = `${market}:${symbol} ${riskLabel(riskLevel)}提醒规则已保存。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function deleteRule(id) {
  try {
    await deleteJSON(`/api/alert-rules?user_id=${encodeURIComponent(state.userID)}&id=${encodeURIComponent(id)}`, state.token);
    await loadRules(false);
    state.message = "提醒规则已删除。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function loadRules(showMessage) {
  try {
    state.rules = asArray(await getJSON(`/api/alert-rules?user_id=${encodeURIComponent(state.userID)}`, state.token));
    if (showMessage) state.message = `已加载 ${state.rules.length} 条提醒规则。`;
  } catch (error) {
    state.message = error.message;
  }
  if (showMessage) render();
}

async function refreshCurrentPool() {
  try {
    await ensureDefaultWatchlist();
    state.job = await postJSON("/api/refresh/manual", { user_id: state.userID, watchlist_id: state.selectedWatchlistID, job_id: `pool-${Date.now()}` }, state.token);
    state.snapshots = state.job.Snapshots || state.job.snapshots || [];
    for (const snapshot of state.snapshots) {
      state.quoteByKey[stockKey(valueOf(snapshot, "Market", "market"), valueOf(snapshot, "Symbol", "symbol"))] = snapshot;
    }
    await Promise.all([loadAlerts(false), loadMarketAnalysisData().catch(() => {})]);
    state.message = `当前股票池刷新完成，共 ${state.snapshots.length} 条行情。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function refreshHoldings() {
  try {
    if (!state.holdings.length) await loadHoldings(false);
    for (const item of state.holdings) {
      await collectStockInfoData(valueOf(item, "Market", "market"), valueOf(item, "Symbol", "symbol"), { force: true });
    }
    state.message = state.holdings.length ? `已刷新 ${state.holdings.length} 条持仓行情。` : "暂无持仓可刷新。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function runAlertCheck() {
  try {
    await refreshCurrentPool();
    await loadAlerts(false);
    state.message = state.alerts.length ? `监控检查完成，已有 ${state.alerts.length} 条提醒日志。` : "监控检查完成，当前没有规则被触发。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function loadAlerts(showMessage) {
  try {
    state.alerts = asArray(await getJSON(`/api/alerts?user_id=${encodeURIComponent(state.userID)}`, state.token));
    state.notifications = asArray(await getJSON(`/api/notifications?user_id=${encodeURIComponent(state.userID)}`, state.token));
    if (showMessage) state.message = "提醒中心已刷新。";
  } catch (error) {
    state.message = error.message;
  }
  if (showMessage) render();
}

async function markNotificationRead(id) {
  try {
    await postJSON("/api/notifications/read", { id }, state.token);
    await loadAlerts(false);
    state.message = "通知已标记为已读。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function collectStockInfoFromForm() {
  const form = readStockForm("stock");
  if (form) await collectStockInfo(form.market, form.symbol, { force: true });
}

async function collectReportMarket() {
  const form = readStockForm("analysis");
  if (!form) return;
  await collectStockInfo(form.market, form.symbol, { force: true });
}

async function collectResearchForCurrentStock() {
  const holding = findHolding(state.selectedMarket, state.selectedSymbol);
  const attentionLevel = document.querySelector("#holding-attention")?.value || valueOf(holding, "AttentionLevel", "attention_level") || "medium";
  try {
    state.researchResult = await postJSON("/api/research/collect", {
      user_id: state.userID,
      market: state.selectedMarket,
      symbol: state.selectedSymbol,
      attention_level: attentionLevel,
    }, state.token);
    state.message = `${state.selectedMarket}:${state.selectedSymbol} 产品信息总结已保存为 RAG 文档。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function runAttentionWorkflow(level) {
  try {
    state.workflowResult = await postJSON("/api/workflows/research/run", {
      user_id: state.userID,
      attention_level: level,
    }, state.token);
    await loadWorkflows(false);
    const job = valueOf(state.workflowResult, "job", "Job");
    state.message = valueOf(job, "summary", "Summary") || `${attentionLabel(level)}关注 AI 工作流已完成。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function loadWorkflows(showMessage) {
  try {
    state.workflows = asArray(await getJSON(`/api/workflows?user_id=${encodeURIComponent(state.userID)}`, state.token));
    if (showMessage) state.message = `已加载 ${state.workflows.length} 条工作流记录。`;
  } catch (error) {
    state.message = error.message;
  }
  if (showMessage) render();
}

async function askAssistant() {
  const market = document.querySelector("#assistant-market")?.value || state.selectedMarket;
  const symbol = (document.querySelector("#assistant-symbol")?.value || "").trim().toUpperCase();
  const question = document.querySelector("#assistant-question")?.value?.trim() || "";
  if (!symbol || !question) {
    state.message = "请填写股票代码和问题。";
    render();
    return;
  }
  state.assistantMessages.push({ role: "user", content: question });
  state.assistantMessages.push({ role: "assistant", content: "" });
  state.assistantStreaming = true;
  state.message = "股票助手正在流式分析。";
  render();
  try {
    state.assistantAnswer = await streamAssistantChat({
      user_id: state.userID,
      session_id: state.assistantSessionID,
      market,
      symbol,
      question,
    });
    state.selectedMarket = valueOf(state.assistantAnswer, "Market", "market") || market;
    state.selectedSymbol = valueOf(state.assistantAnswer, "Symbol", "symbol") || symbol;
    state.message = "股票助手已完成分析。";
  } catch (error) {
    state.message = error.message;
    const last = state.assistantMessages[state.assistantMessages.length - 1];
    if (last?.role === "assistant" && !last.content) last.content = `分析失败：${error.message}`;
  } finally {
    state.assistantStreaming = false;
  }
  render();
}

async function streamAssistantChat(payload) {
  const response = await fetch(`${API_BASE}/api/assistant/chat/stream`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Request-ID": `web-${Date.now()}`,
      ...(state.token ? { Authorization: `Bearer ${state.token}` } : {}),
    },
    body: JSON.stringify(payload),
  });
  if (!response.ok) {
    const data = await response.json().catch(() => ({}));
    if (shouldInvalidateSession(data, response.status)) {
      window.localStorage.removeItem("jijin_token");
      throw new Error("登录已过期，请重新登录。");
    }
    throw new Error(data?.error?.message || `HTTP ${response.status}`);
  }
  if (!response.body) {
    return await postJSON("/api/assistant/chat", payload, state.token);
  }
  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  let buffer = "";
  let finalResponse = null;
  while (true) {
    const { value, done } = await reader.read();
    if (done) break;
    buffer += decoder.decode(value, { stream: true });
    const events = buffer.split("\n\n");
    buffer = events.pop() || "";
    for (const event of events) {
      const line = event.split("\n").find((item) => item.startsWith("data:"));
      if (!line) continue;
      const data = JSON.parse(line.slice(5).trim());
      if (data.delta) {
        const last = state.assistantMessages[state.assistantMessages.length - 1];
        if (last?.role === "assistant") last.content += data.delta;
        render();
      }
      if (data.done) {
        finalResponse = data.response;
      }
    }
  }
  return finalResponse || { market: payload.market, symbol: payload.symbol, answer: state.assistantMessages.at(-1)?.content || "" };
}

async function collectStockInfo(market, symbol, options = {}) {
  try {
    const result = await collectStockInfoData(market, symbol, options);
    state.message = result.cached
      ? `${market}:${symbol} 已使用 ${Math.ceil(STOCK_DETAIL_COOLDOWN_MS / 60000)} 分钟内缓存，避免频繁请求数据源。`
      : `已获取 ${market}:${symbol} 行情和公司信息。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function collectStockInfoData(market, symbol, options = {}) {
  const normalizedMarket = String(market || "CN").trim().toUpperCase();
  const normalizedSymbol = String(symbol || "").trim().toUpperCase();
  if (!normalizedSymbol) throw new Error("请输入股票代码。");
  state.selectedMarket = normalizedMarket;
  state.selectedSymbol = normalizedSymbol;
  window.localStorage.setItem("jijin_selected_market", normalizedMarket);
  window.localStorage.setItem("jijin_selected_symbol", normalizedSymbol);
  const key = stockKey(normalizedMarket, normalizedSymbol);
  const fetchedAt = state.quoteFetchedAt[key] || 0;
  if (!options.force && state.quoteByKey[key] && Date.now() - fetchedAt < STOCK_DETAIL_COOLDOWN_MS) {
    await loadMarketAnalysisData();
    return { cached: true };
  }
  const data = await postJSON("/api/market/collect", { market: normalizedMarket, symbol: normalizedSymbol }, state.token);
  state.collected = data.snapshot || data.Snapshot;
  state.quoteByKey[key] = state.collected;
  state.quoteFetchedAt[key] = Date.now();
  state.snapshots = asArray(await getJSON(`/api/market/snapshots?market=${encodeURIComponent(normalizedMarket)}&symbol=${encodeURIComponent(normalizedSymbol)}`, state.token));
  state.dailyChanges = asArray(data.daily_changes || data.DailyChanges);
  state.profile = data.profile || data.Profile || null;
  state.profileByKey[key] = state.profile;
  return { cached: false };
}

async function loadMarketAnalysis() {
  try {
    await loadMarketAnalysisData();
    state.message = `${state.selectedMarket}:${state.selectedSymbol} 已加载保存的行情、涨跌和公司分析。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function loadMarketAnalysisData() {
  const market = encodeURIComponent(state.selectedMarket);
  const symbol = encodeURIComponent(state.selectedSymbol);
  state.snapshots = asArray(await getJSON(`/api/market/snapshots?market=${market}&symbol=${symbol}`, state.token));
  state.dailyChanges = asArray(await getJSON(`/api/market/daily-changes?market=${market}&symbol=${symbol}`, state.token));
  state.profile = await getJSON(`/api/stocks/profile?market=${market}&symbol=${symbol}`, state.token);
  state.profileByKey[stockKey(state.selectedMarket, state.selectedSymbol)] = state.profile;
}

async function saveAccountConfig() {
  const provider = document.querySelector("#account-provider")?.value || "csv";
  const id = document.querySelector("#account-id")?.value?.trim() || `${provider}-${Date.now().toString(36)}`;
  const alias = document.querySelector("#account-alias")?.value?.trim() || providerLabel(provider);
  const refreshMode = document.querySelector("#account-refresh-mode")?.value || "manual";
  const readOnly = document.querySelector("#account-readonly")?.checked === true;
  try {
    await postJSON("/api/accounts", {
      id,
      user_id: state.userID,
      alias,
      refresh_mode: refreshMode,
      read_only: readOnly,
      metadata: { provider },
    }, state.token);
    await loadAccounts(false);
    state.message = `${providerLabel(provider)}只读账户配置已保存。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function loadAccounts(showMessage) {
  try {
    state.accountConfigs = asArray(await getJSON(`/api/accounts?user_id=${encodeURIComponent(state.userID)}`, state.token));
    if (showMessage) state.message = `已加载 ${state.accountConfigs.length} 个账户配置。`;
  } catch (error) {
    state.message = error.message;
  }
  if (showMessage) render();
}

async function loadDependencies() {
  try {
    state.dependencies = await getJSON("/api/system/dependencies", state.token);
    state.message = "依赖状态已刷新。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

function saveLocalUserSettings() {
  state.displayName = document.querySelector("#display-name")?.value?.trim() || "";
  window.localStorage.setItem("jijin_display_name", state.displayName);
  state.message = "用户显示信息已保存到本机。";
  render();
}

function renderStockPicker(prefix, includePrices) {
  return `
    <label>市场
      <select id="${prefix}-market">
        ${["CN", "HK", "US"].map((market) => `<option value="${market}" ${market === state.selectedMarket ? "selected" : ""}>${marketLabel(market)}</option>`).join("")}
      </select>
    </label>
    <label>股票代码 <input id="${prefix}-symbol" value="${escapeHTML(state.selectedSymbol)}" placeholder="例如 000821" /></label>
    ${includePrices ? "" : ""}
  `;
}

function renderDependencies() {
  if (!state.dependencies) return "<p class=\"muted\">点击按钮后检查连接状态。</p>";
  const labels = { database: "数据库", redis: "Redis", llm: "大模型", stock_source: "行情源" };
  return `<ul class="dependency-list">${Object.keys(labels).map((key) => {
    const item = state.dependencies[key];
    return `<li><span class="${item?.reachable ? "dot ok" : "dot bad"}"></span>${labels[key]}：${item?.reachable ? "可用" : "未就绪"} <span class="muted">${item?.message || ""}</span></li>`;
  }).join("")}</ul>`;
}

function readStockForm(prefix) {
  const market = document.querySelector(`#${prefix}-market`)?.value || state.selectedMarket || "CN";
  const symbol = (document.querySelector(`#${prefix}-symbol`)?.value || "").trim().toUpperCase();
  const buyPrice = Number(document.querySelector("#buy-price")?.value || 0);
  const sellPrice = Number(document.querySelector("#sell-price")?.value || 0);
  if (!symbol) {
    state.message = "请输入股票代码。";
    render();
    return null;
  }
  if (buyPrice < 0 || sellPrice < 0) {
    state.message = "买入关注价和卖出关注价不能为负数。";
    render();
    return null;
  }
  return { market, symbol, buyPrice, sellPrice };
}

function currentSymbols() {
  return asArray(state.watchlist?.symbols || state.watchlist?.Symbols);
}

function allPoolSymbols() {
  const out = [];
  for (const pool of state.watchlists) {
    for (const item of (pool.Symbols || pool.symbols || [])) {
      out.push({ market: valueOf(item, "Market", "market"), symbol: valueOf(item, "Symbol", "symbol"), source: "股票池" });
    }
  }
  return uniqueStocks(out);
}

function analysisTargets() {
  const holdingTargets = state.holdings.map((item) => ({
    market: valueOf(item, "Market", "market"),
    symbol: valueOf(item, "Symbol", "symbol"),
    source: "持仓",
  }));
  return uniqueStocks([...holdingTargets, ...allPoolSymbols(), { market: state.selectedMarket, symbol: state.selectedSymbol, source: "当前" }]);
}

function renderAnalysisRows(items) {
  const rows = uniqueStocks(items).map((item) => {
    const key = stockKey(item.market, item.symbol);
    const quote = state.quoteByKey[key];
    const percent = Number(valueOf(quote, "ChangePercent", "change_percent") || 0);
    const cls = percent > 0 ? "cn-up" : percent < 0 ? "cn-down" : "muted";
    const sign = percent > 0 ? "+" : "";
    return `<tr>
      <td><button class="link-button" data-load-stock="${key}" type="button">${key}</button></td>
      <td>${item.source || "自选"}</td>
      <td>${quote ? formatPrice(valueOf(quote, "Price", "price")) : "未刷新"}</td>
      <td class="${cls}">${quote ? `${sign}${percent.toFixed(2)}%` : "-"}</td>
      <td>${findHolding(item.market, item.symbol) ? holdingSummary(findHolding(item.market, item.symbol)) : "未持仓"}</td>
    </tr>`;
  }).join("");
  if (!rows) return "";
  return `<table class="stock-table"><thead><tr><th>代码</th><th>来源</th><th>最新价</th><th>今日涨跌</th><th>持仓情况</th></tr></thead><tbody>${rows}</tbody></table>`;
}

function findHolding(market, symbol) {
  const key = stockKey(market, symbol);
  return state.holdings.find((item) => stockKey(valueOf(item, "Market", "market"), valueOf(item, "Symbol", "symbol")) === key);
}

function holdingSummary(item) {
  const market = valueOf(item, "Market", "market");
  const symbol = valueOf(item, "Symbol", "symbol");
  const quantity = Number(valueOf(item, "Quantity", "quantity") || 0);
  const cost = Number(valueOf(item, "CostBasis", "cost_basis") || 0);
  const attention = valueOf(item, "AttentionLevel", "attention_level") || "medium";
  const quote = state.quoteByKey[stockKey(market, symbol)];
  const latest = Number(valueOf(quote, "Price", "price") || cost);
  const pnl = (latest - cost) * quantity;
  return `${quantity} 股，成本 ${cost.toFixed(2)}，关注${attentionLabel(attention)}，估算盈亏 ${pnl.toFixed(2)}`;
}

function uniqueStocks(items) {
  const map = new Map();
  for (const item of items) {
    if (!item.market || !item.symbol) continue;
    const key = stockKey(item.market, item.symbol);
    if (!map.has(key)) map.set(key, { market: String(item.market).toUpperCase(), symbol: String(item.symbol).toUpperCase(), source: item.source || "自选" });
  }
  return [...map.values()];
}

function stockOptions(items) {
  return items.map((item) => `<option value="${stockKey(item.market, item.symbol)}">${item.source || "自选"} ${stockKey(item.market, item.symbol)}</option>`).join("");
}

function asArray(value) {
  return Array.isArray(value) ? value : [];
}

function portfolioSummary() {
  if (!state.holdings.length) return "暂无持仓，添加持仓后会自动进入组合分析。";
  let costValue = 0;
  let latestValue = 0;
  for (const item of state.holdings) {
    const market = valueOf(item, "Market", "market");
    const symbol = valueOf(item, "Symbol", "symbol");
    const quantity = Number(valueOf(item, "Quantity", "quantity") || 0);
    const cost = Number(valueOf(item, "CostBasis", "cost_basis") || 0);
    const quote = state.quoteByKey[stockKey(market, symbol)];
    const latest = Number(valueOf(quote, "Price", "price") || cost);
    costValue += cost * quantity;
    latestValue += latest * quantity;
  }
  const pnl = latestValue - costValue;
  return `持仓 ${state.holdings.length} 条，成本 ${costValue.toFixed(2)}，估算市值 ${latestValue.toFixed(2)}，浮动盈亏 ${pnl.toFixed(2)}。`;
}

function currentWatchlistName() {
  const item = state.watchlists.find((pool) => itemID(pool) === state.selectedWatchlistID) || state.watchlist;
  return valueOf(item, "Name", "name") || "我的股票池";
}

function itemID(item) {
  return valueOf(item, "ID", "id") || "";
}

function stockKey(market, symbol) {
  return `${String(market || "").toUpperCase()}:${String(symbol || "").toUpperCase()}`;
}

function formatPrice(value) {
  const price = Number(value || 0);
  return price > 0 ? price.toFixed(2) : "-";
}

function setInput(selector, value) {
  const input = document.querySelector(selector);
  if (input) input.value = value;
}

function marketLabel(market) {
  return { CN: "A 股 / CN", HK: "港股 / HK", US: "美股 / US" }[market] || market;
}

function modeLabel(mode) {
  return { manual: "手动", conservative: "保守自动", standard: "标准自动" }[mode] || mode || "-";
}

function providerLabel(provider) {
  return { ths: "同花顺", eastmoney: "东方财富", xueqiu: "雪球", csv: "CSV/手工导入", other: "其他只读来源" }[provider] || provider || "-";
}

function ruleTypeLabel(type) {
  return { price_below: "价格低于", price_above: "价格高于", change_percent_below: "跌幅低于", change_percent_above: "涨幅高于" }[type] || type || "-";
}

function signalLabel(signal) {
  return { buy_watch: "买入关注", sell_watch: "卖出关注", risk_warning: "风险提示", hold_watch: "继续观察" }[signal] || signal || "-";
}

function riskLabel(risk) {
  return { low: "低", medium: "中", high: "高", critical: "严重" }[risk] || risk || "低";
}

function attentionLabel(level) {
  return { high: "高", medium: "中", low: "低" }[level] || "中";
}

function attentionRisk(level) {
  return { high: "high", medium: "medium", low: "low" }[level] || "medium";
}

function applyTheme() {
  document.documentElement.dataset.theme = state.theme;
}

function escapeHTML(value) {
  return String(value ?? "").replace(/[&<>"']/g, (char) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", "\"": "&quot;", "'": "&#39;" }[char]));
}

if (state.token) {
  bootstrapUserData().finally(render);
} else {
  render();
}
