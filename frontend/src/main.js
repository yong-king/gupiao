import { formatRefreshStatus, getViewCopy, isKnownView, layoutClassForAuthState, navItems, refreshModes } from "./app.js?v=19";
import { getJSON, postJSON } from "./api.js?v=19";
import { formatDailyChange, monitorText, renderChangeCalendar, renderPriceChart, summarizeProfile, valueOf } from "./market.js?v=19";

const root = document.querySelector("#app");
const state = {
  token: window.localStorage.getItem("jijin_token") || "",
  email: window.localStorage.getItem("jijin_email") || "",
  userID: window.localStorage.getItem("jijin_user_id") || "user-demo",
  activeView: window.localStorage.getItem("jijin_active_view") || "watchlists",
  message: "",
  alerts: [],
  notifications: [],
  watchlists: [],
  watchlist: null,
  selectedMarket: window.localStorage.getItem("jijin_selected_market") || "US",
  selectedSymbol: window.localStorage.getItem("jijin_selected_symbol") || "AAPL",
  quoteByKey: {},
  profileByKey: {},
  holdings: [],
  rules: [],
  rule: null,
  holdingsImport: null,
  job: null,
  collected: null,
  snapshots: [],
  dailyChanges: [],
  profile: null,
  dependencies: null,
};

function render() {
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
    <aside class="sidebar">${nav}</aside>
    <main class="content">
      <section class="topbar">
        <h1>股票监控操作台</h1>
        <div class="userbar">
          <span class="muted">当前用户：${state.email}</span>
          <button id="logout" type="button">退出登录</button>
        </div>
      </section>
      ${renderActiveView()}
      <p class="muted">${state.message}</p>
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
    window.localStorage.removeItem("jijin_token");
    window.localStorage.removeItem("jijin_email");
    window.localStorage.removeItem("jijin_user_id");
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
      <span class="safety-badge">仅提醒，不自动交易</span>
    </section>
    ${renderViewBody(view)}
  `;
}

function renderViewBody(view) {
  if (state.activeView === "watchlists") {
    const symbols = state.watchlist?.symbols || state.watchlist?.Symbols || [];
    const watchlists = state.watchlists.map((item) => `<li>${item.Name || item.name} (${item.ID || item.id})</li>`).join("");
    const symbolRows = symbols.map((item) => {
      const market = valueOf(item, "Market", "market");
      const symbol = valueOf(item, "Symbol", "symbol");
      const key = stockKey(market, symbol);
      const quote = state.quoteByKey[key];
      return `<tr>
        <td>${market}:${symbol}</td>
        <td>${formatPrice(valueOf(item, "BuyPrice", "buy_price"))}</td>
        <td>${formatPrice(valueOf(item, "SellPrice", "sell_price"))}</td>
        <td>${quote ? formatPrice(valueOf(quote, "Price", "price")) : "-"}</td>
        <td>${monitorText(quote, item)}</td>
        <td><button data-load-stock="${key}" type="button">获取</button></td>
      </tr>`;
    }).join("");
    const selectedProfile = state.profileByKey[stockKey(state.selectedMarket, state.selectedSymbol)] || state.profile;
    const selectedQuote = state.quoteByKey[stockKey(state.selectedMarket, state.selectedSymbol)] || state.collected;
    return `
      <section class="grid">
        <article>
          <h3>添加股票</h3>
          <p>${state.watchlist ? `名称：${state.watchlist.name || state.watchlist.Name}` : view.empty}</p>
          <label>市场
            <select id="stock-market">
              ${["US", "HK", "CN"].map((market) => `<option value="${market}" ${market === state.selectedMarket ? "selected" : ""}>${market}</option>`).join("")}
            </select>
          </label>
          <label>股票代码 <input id="stock-symbol" value="${state.selectedSymbol}" placeholder="AAPL" /></label>
          <div class="two-col">
            <label>买入关注价 <input id="buy-price" type="number" min="0" step="0.01" placeholder="例如 160" /></label>
            <label>卖出关注价 <input id="sell-price" type="number" min="0" step="0.01" placeholder="例如 230" /></label>
          </div>
          <div class="actions">
            <button id="lookup-stock" type="button">获取股票信息</button>
            <button id="save-monitor-stock" type="button">加入/更新监控</button>
          </div>
          <button id="load-watchlists" type="button">刷新股票池列表</button>
        </article>
        <article>
          <h3>股票信息</h3>
          ${selectedQuote ? `<p>${state.selectedMarket}:${state.selectedSymbol} 当前 ${formatPrice(valueOf(selectedQuote, "Price", "price"))}</p>` : "<p>输入股票代码后点击获取股票信息。</p>"}
          <p>${selectedProfile ? summarizeProfile(selectedProfile) : "暂无公司信息。"}</p>
          <p class="muted">${selectedProfile ? valueOf(selectedProfile, "analysis", "Analysis") : "获取后会展示公司、产品和观察建议。"}</p>
        </article>
        <article class="wide">
          <h3>已监控代码</h3>
          ${symbols.length ? `<table class="stock-table"><thead><tr><th>代码</th><th>买入关注</th><th>卖出关注</th><th>当前价</th><th>状态</th><th></th></tr></thead><tbody>${symbolRows}</tbody></table>` : `<p>${view.empty}</p>`}
        </article>
        <article>
          <h3>我的股票池</h3>
          ${state.watchlists.length ? `<ul>${watchlists}</ul>` : "<p>暂无列表，点击刷新股票池列表。</p>"}
        </article>
      </section>
    `;
  }
  if (state.activeView === "holdings") {
    const rows = state.holdings.map((item) => `<li>${item.Market || item.market}:${item.Symbol || item.symbol} 数量 ${item.Quantity || item.quantity} 成本 ${item.CostBasis || item.cost_basis}</li>`).join("");
    return `
      <section class="grid">
        <article>
          <h3>持仓导入</h3>
          <p>${state.holdingsImport ? `已导入 ${state.holdingsImport.imported} 条持仓` : view.empty}</p>
          <button id="import-holdings" type="button">导入 Demo 持仓</button>
          <button id="load-holdings" type="button">刷新持仓列表</button>
        </article>
        <article>
          <h3>导入格式</h3>
          <p class="code">market,symbol,quantity,cost_basis</p>
          <p class="muted">只读数据用于提醒和风险分析。</p>
        </article>
        <article>
          <h3>当前持仓</h3>
          ${state.holdings.length ? `<ul>${rows}</ul>` : "<p>暂无持仓，点击导入 Demo 持仓。</p>"}
        </article>
      </section>
    `;
  }
  if (state.activeView === "rules") {
    const rules = state.rules.map((item) => `<li>${item.Market || item.market}:${item.Symbol || item.symbol} ${item.Type || item.type} ${item.Threshold || item.threshold} -> ${item.Signal || item.signal}</li>`).join("");
    return `
      <section class="grid">
        <article>
          <h3>AAPL 价格提醒</h3>
          <p>${state.rule ? `${state.rule.Signal || state.rule.signal}，阈值 ${state.rule.Threshold || state.rule.threshold}` : view.empty}</p>
          <button id="create-rule" type="button">创建 Demo 提醒规则</button>
          <button id="load-rules" type="button">刷新规则列表</button>
        </article>
        <article>
          <h3>规则原则</h3>
          <p class="muted">规则只产生观察、买入关注、卖出关注和风险提示，不触发交易。</p>
        </article>
        <article>
          <h3>当前规则</h3>
          ${state.rules.length ? `<ul>${rules}</ul>` : "<p>暂无规则，点击创建 Demo 提醒规则。</p>"}
        </article>
      </section>
    `;
  }
  if (state.activeView === "refresh") {
    const modes = refreshModes.map((mode) => `<option value="${mode.id}">${mode.label}</option>`).join("");
    return `
      <section class="grid">
        <article>
          <h3>刷新状态</h3>
          <p>${formatRefreshStatus(state.job)}</p>
          <button id="run-demo" type="button">运行真实 AAPL 刷新</button>
          <button id="collect-real" type="button">测试并保存 AAPL 行情</button>
        </article>
        <article>
          <h3>刷新模式</h3>
          <select>${modes}</select>
          <p class="muted">自动刷新应保持低频，并遵守数据源限制。</p>
        </article>
        <article>
          <h3>最近保存行情</h3>
          <p>${state.collected ? `${valueOf(state.collected, "Name", "name") || "AAPL"} 收盘 ${valueOf(state.collected, "Price", "price")}` : "暂无保存记录。"}</p>
        </article>
      </section>
    `;
  }
  if (state.activeView === "alerts") {
    const alerts = state.alerts.map((item) => `<li>${item.Signal || item.signal} ${item.Market || item.market}:${item.Symbol || item.symbol} - ${item.Summary || item.summary}</li>`).join("");
    return `
      <section class="grid">
        <article>
          <h3>已触发提醒</h3>
          ${state.alerts.length ? `<ul>${alerts}</ul>` : `<p>${view.empty}</p>`}
          <button id="load-alerts" type="button">刷新提醒列表</button>
        </article>
        <article>
          <h3>通知状态</h3>
          <p>${state.notifications.length ? `通知 ${state.notifications.length} 条` : "暂无通知。"}</p>
        </article>
      </section>
    `;
  }
  if (state.activeView === "reports") {
    const profileSummary = state.profile ? summarizeProfile(state.profile) : "尚未加载公司/产品分析。";
    const changes = state.dailyChanges.map((item) => `<li>${formatDailyChange(item)}</li>`).join("");
    return `
      <section class="grid">
        <article>
          <h3>${state.selectedSymbol} 真实行情曲线</h3>
          ${renderPriceChart(state.snapshots)}
          <button id="collect-report-market" type="button">测试并保存当前股票行情</button>
          <button id="load-market-analysis" type="button">加载已保存行情与分析</button>
        </article>
        <article>
          <h3>涨跌日历</h3>
          ${renderChangeCalendar(state.dailyChanges)}
        </article>
        <article>
          <h3>每日涨跌记录</h3>
          ${state.dailyChanges.length ? `<ul>${changes}</ul>` : `<p>${view.empty}</p>`}
          <p class="muted">这些记录已生成 RAG 文本字段，后续可入库向量化。</p>
        </article>
        <article>
          <h3>公司和产品分析</h3>
          <p>${profileSummary}</p>
          <p>${state.profile ? valueOf(state.profile, "analysis", "Analysis") : "加载后展示业务、产品和走势分析。"}</p>
          <p class="muted">${state.profile ? valueOf(state.profile, "disclaimer", "Disclaimer") : "仅用于监控和研究，不构成投资建议。"}</p>
        </article>
      </section>
    `;
  }
  if (state.activeView === "accounts") {
    return `
      <section class="grid">
        <article>
          <h3>账户接入</h3>
          <p>${view.empty}</p>
        </article>
        <article>
          <h3>安全边界</h3>
          <p class="muted">只允许只读数据接入，不实现买入、卖出或自动下单。</p>
        </article>
      </section>
    `;
  }
  return `
    <section class="grid">
      <article>
        <h3>配置文件</h3>
        <p class="code">config/backend.example.json</p>
        <p class="code">agent/config/agent.example.json</p>
        <button id="load-dependencies" type="button">检查后端依赖状态</button>
      </article>
      <article>
        <h3>部署入口</h3>
        <p class="code">sh deploy/scripts/deploy.sh</p>
        <p class="muted">${view.empty}</p>
      </article>
      <article>
        <h3>依赖状态</h3>
        ${renderDependencies()}
      </article>
    </section>
  `;
}

function bindViewActions() {
  document.querySelector("#ensure-watchlist")?.addEventListener("click", ensureWatchlist);
  document.querySelector("#add-aapl")?.addEventListener("click", addDemoSymbol);
  document.querySelector("#lookup-stock")?.addEventListener("click", collectStockInfoFromForm);
  document.querySelector("#save-monitor-stock")?.addEventListener("click", saveMonitorStock);
  document.querySelectorAll("[data-load-stock]").forEach((button) => {
    button.addEventListener("click", () => {
      const [market, symbol] = button.dataset.loadStock.split(":");
      collectStockInfo(market, symbol);
    });
  });
  document.querySelector("#load-watchlists")?.addEventListener("click", loadWatchlists);
  document.querySelector("#import-holdings")?.addEventListener("click", importDemoHoldings);
  document.querySelector("#load-holdings")?.addEventListener("click", loadHoldings);
  document.querySelector("#create-rule")?.addEventListener("click", createDemoRule);
  document.querySelector("#load-rules")?.addEventListener("click", loadRules);
  document.querySelector("#run-demo")?.addEventListener("click", runDemo);
  document.querySelector("#collect-real")?.addEventListener("click", collectRealMarket);
  document.querySelector("#load-alerts")?.addEventListener("click", loadAlerts);
  document.querySelector("#collect-report-market")?.addEventListener("click", collectRealMarket);
  document.querySelector("#load-market-analysis")?.addEventListener("click", loadMarketAnalysis);
  document.querySelector("#load-dependencies")?.addEventListener("click", loadDependencies);
}

function renderAuth() {
  root.className = layoutClassForAuthState(false);
  root.innerHTML = `
    <main class="auth">
      <section class="auth-card">
        <h1>股票监控操作台</h1>
        <p class="muted">请先注册或登录。密码会在后端做加盐哈希，系统不保存明文密码。</p>
        <label>邮箱 <input id="email" value="demo@example.com" /></label>
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
    const payload = mode === "register"
      ? { id: `user-${Date.now().toString(36)}`, email, password }
      : { email, password };
    const data = await postJSON(`/api/auth/${mode}`, payload);
    state.token = data.token;
    state.email = data.email;
    state.userID = data.user_id || data.userID || state.userID;
    window.localStorage.setItem("jijin_token", state.token);
    window.localStorage.setItem("jijin_email", state.email);
    window.localStorage.setItem("jijin_user_id", state.userID);
    state.message = "登录成功。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function ensureWatchlist() {
  try {
    await postJSON("/api/watchlists", { id: "wl-demo", user_id: state.userID, name: "Demo Watchlist" }, state.token).catch(() => {});
    state.watchlist = await getJSON("/api/watchlists/wl-demo", state.token);
    state.watchlists = await getJSON(`/api/watchlists?user_id=${encodeURIComponent(state.userID)}`, state.token);
    state.message = "Demo 股票池已准备。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function addDemoSymbol() {
  try {
    await ensureWatchlistData();
    state.watchlist = await postJSON("/api/watchlists/wl-demo/symbols", { market: "US", symbol: "AAPL", buy_price: 160, sell_price: 230 }, state.token).catch(async () => getJSON("/api/watchlists/wl-demo", state.token));
    state.watchlists = await getJSON(`/api/watchlists?user_id=${encodeURIComponent(state.userID)}`, state.token);
    state.message = "AAPL 已加入股票池。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function importDemoHoldings() {
  try {
    const csv = "market,symbol,quantity,cost_basis\nUS,AAPL,10,145\nUS,MSFT,5,390\n";
    state.holdingsImport = await postJSON("/api/holdings/import", { user_id: state.userID, csv }, state.token);
    state.holdings = state.holdingsImport.holdings || state.holdingsImport.Holdings || await getJSON(`/api/holdings?user_id=${encodeURIComponent(state.userID)}`, state.token);
    state.message = "Demo 持仓已导入。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function collectStockInfoFromForm() {
  const form = readStockForm();
  if (!form) return;
  await collectStockInfo(form.market, form.symbol);
}

async function saveMonitorStock() {
  const form = readStockForm();
  if (!form) return;
  try {
    await ensureWatchlistData();
    await collectStockInfoData(form.market, form.symbol);
    state.watchlist = await postJSON("/api/watchlists/wl-demo/symbols", {
      market: form.market,
      symbol: form.symbol,
      buy_price: form.buyPrice,
      sell_price: form.sellPrice,
    }, state.token);
    await createPriceRules(form);
    state.watchlists = await getJSON(`/api/watchlists?user_id=${encodeURIComponent(state.userID)}`, state.token);
    state.rules = await getJSON(`/api/alert-rules?user_id=${encodeURIComponent(state.userID)}`, state.token);
    state.message = `${form.market}:${form.symbol} 已加入监控。`;
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

async function createDemoRule() {
  try {
    state.rule = await postJSON("/api/alert-rules", {
      id: "rule-demo",
      user_id: state.userID,
      market: "US",
      symbol: "AAPL",
      type: "price_below",
      threshold: 160,
      signal: "buy_watch",
      risk_level: "medium",
      enabled: true,
      cooldown_seconds: 1800,
    }, state.token).catch(async () => {
      const rules = await getJSON(`/api/alert-rules?user_id=${encodeURIComponent(state.userID)}`, state.token);
      return rules.find((item) => (item.ID || item.id) === "rule-demo") || null;
    });
    state.rules = await getJSON(`/api/alert-rules?user_id=${encodeURIComponent(state.userID)}`, state.token);
    state.message = "Demo 提醒规则已准备。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function loadWatchlists() {
  try {
    state.watchlists = await getJSON(`/api/watchlists?user_id=${encodeURIComponent(state.userID)}`, state.token);
    state.message = `已加载 ${state.watchlists.length} 个股票池。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function loadHoldings() {
  try {
    state.holdings = await getJSON(`/api/holdings?user_id=${encodeURIComponent(state.userID)}`, state.token);
    state.message = `已加载 ${state.holdings.length} 条持仓。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function loadRules() {
  try {
    state.rules = await getJSON(`/api/alert-rules?user_id=${encodeURIComponent(state.userID)}`, state.token);
    state.message = `已加载 ${state.rules.length} 条提醒规则。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function loadAlerts() {
  try {
    state.alerts = await getJSON(`/api/alerts?user_id=${encodeURIComponent(state.userID)}`, state.token);
    state.notifications = await getJSON(`/api/notifications?user_id=${encodeURIComponent(state.userID)}`, state.token);
    state.message = "提醒列表已刷新。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function ensureWatchlistData() {
  await postJSON("/api/watchlists", { id: "wl-demo", user_id: state.userID, name: "Demo Watchlist" }, state.token).catch(() => {});
  state.watchlist = await getJSON("/api/watchlists/wl-demo", state.token);
}

async function runDemo() {
  try {
    await ensureWatchlistData();
    await postJSON("/api/watchlists/wl-demo/symbols", { market: "US", symbol: "AAPL", buy_price: 160, sell_price: 230 }, state.token).catch(() => {});
    await createDemoRuleData();
    state.job = await postJSON("/api/refresh/manual", { user_id: state.userID, watchlist_id: "wl-demo", job_id: `job-${Date.now()}` }, state.token);
    state.snapshots = state.job.Snapshots || state.job.snapshots || [];
    state.alerts = await getJSON(`/api/alerts?user_id=${encodeURIComponent(state.userID)}`, state.token);
    await loadMarketAnalysisData();
    state.message = "Demo 链路已执行。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function collectRealMarket() {
  try {
    await collectStockInfoData(state.selectedMarket, state.selectedSymbol);
    state.message = `已测试真实行情源并保存 ${state.selectedMarket}:${state.selectedSymbol} 快照。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function collectStockInfo(market, symbol) {
  try {
    await collectStockInfoData(market, symbol);
    state.message = `已获取 ${market}:${symbol} 股票信息。`;
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function collectStockInfoData(market, symbol) {
  const normalizedMarket = String(market || "US").trim().toUpperCase();
  const normalizedSymbol = String(symbol || "").trim().toUpperCase();
  if (!normalizedSymbol) throw new Error("请输入股票代码。");
  state.selectedMarket = normalizedMarket;
  state.selectedSymbol = normalizedSymbol;
  window.localStorage.setItem("jijin_selected_market", normalizedMarket);
  window.localStorage.setItem("jijin_selected_symbol", normalizedSymbol);
  const data = await postJSON("/api/market/collect", { market: normalizedMarket, symbol: normalizedSymbol }, state.token);
  state.collected = data.snapshot || data.Snapshot;
  state.quoteByKey[stockKey(normalizedMarket, normalizedSymbol)] = state.collected;
  state.snapshots = await getJSON(`/api/market/snapshots?market=${encodeURIComponent(normalizedMarket)}&symbol=${encodeURIComponent(normalizedSymbol)}`, state.token);
  state.dailyChanges = data.daily_changes || data.DailyChanges || [];
  state.profile = data.profile || data.Profile || null;
  state.profileByKey[stockKey(normalizedMarket, normalizedSymbol)] = state.profile;
}

async function createDemoRuleData() {
  state.rule = await postJSON("/api/alert-rules", {
    id: "rule-demo",
    user_id: state.userID,
    market: "US",
    symbol: "AAPL",
    type: "price_below",
    threshold: 160,
    signal: "buy_watch",
    risk_level: "medium",
    enabled: true,
    cooldown_seconds: 1800,
  }, state.token).catch(async () => {
    const rules = await getJSON(`/api/alert-rules?user_id=${encodeURIComponent(state.userID)}`, state.token);
    return rules.find((item) => (item.ID || item.id) === "rule-demo") || null;
  });
}

async function loadMarketAnalysis() {
  try {
    await loadMarketAnalysisData();
    state.message = "真实行情、每日涨跌和公司分析已加载。";
  } catch (error) {
    state.message = error.message;
  }
  render();
}

async function loadMarketAnalysisData() {
  const market = encodeURIComponent(state.selectedMarket);
  const symbol = encodeURIComponent(state.selectedSymbol);
  state.snapshots = await getJSON(`/api/market/snapshots?market=${market}&symbol=${symbol}`, state.token);
  state.dailyChanges = await getJSON(`/api/market/daily-changes?market=${market}&symbol=${symbol}`, state.token);
  state.profile = await getJSON(`/api/stocks/profile?market=${market}&symbol=${symbol}`, state.token);
  state.profileByKey[stockKey(state.selectedMarket, state.selectedSymbol)] = state.profile;
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

function renderDependencies() {
  if (!state.dependencies) return "<p>点击检查后端依赖状态。</p>";
  return `<ul>${["database", "redis", "llm", "stock_source"].map((key) => {
    const item = state.dependencies[key];
    return `<li>${key}: ${item?.reachable ? "可用" : "未就绪"} - ${item?.message || ""}</li>`;
  }).join("")}</ul>`;
}

function readStockForm() {
  const market = document.querySelector("#stock-market")?.value || state.selectedMarket || "US";
  const symbol = (document.querySelector("#stock-symbol")?.value || "").trim().toUpperCase();
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

function stockKey(market, symbol) {
  return `${String(market || "").toUpperCase()}:${String(symbol || "").toUpperCase()}`;
}

function formatPrice(value) {
  const price = Number(value || 0);
  return price > 0 ? price.toFixed(2) : "-";
}

render();
