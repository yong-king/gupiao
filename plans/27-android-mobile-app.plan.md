# 27 安卓移动端应用计划

状态：spec_ready

## 目标

为项目新增一个面向安卓手机的可安装移动应用，复用现有后端 API、认证和股票助手能力，优先覆盖登录、持仓、自选、提醒、股票详情和股票助手等高频核心场景。

## 范围

- 新增 `mobile/` 目录与移动端工程结构。
- 引入 Capacitor Android 壳。
- 设计并实现移动端基础导航与页面框架。
- 适配登录、持仓、自选、提醒、股票详情、股票助手和设置页面。
- 复用现有 API 与鉴权机制。
- 补充移动端开发、调试、构建文档。
- 为后续推送、生物识别和原生能力扩展预留接口。

## 非目标

- 不在本计划内实现 iOS。
- 不在本计划内实现完整桌面端所有功能。
- 不在本计划内正式接入推送服务。
- 不在本计划内实现离线数据库同步。
- 不实现自动交易相关能力。

## 依赖

- [08 前端控制台计划](</Users/youngking/Documents/jijin/plans/08-frontend-console.plan.md>)
- [09 发布与加固计划](</Users/youngking/Documents/jijin/plans/09-hardening-release.plan.md>)
- [11 注册登录计划](</Users/youngking/Documents/jijin/plans/11-auth-registration-login.plan.md>)
- [24 LangGraph Agent 聊天体验与界面优化计划](</Users/youngking/Documents/jijin/plans/24-langgraph-agent-chat-ux.plan.md>)
- [26 股票智能体多 Agent 计划](</Users/youngking/Documents/jijin/plans/26-stock-intelligence-multi-agent.plan.md>)

## 对应规范

- [27 安卓移动端应用规范](</Users/youngking/Documents/jijin/docs/specs/27-android-mobile-app.spec.md>)

## 任务

- [ ] 确认移动端技术路线与目录结构。
- [ ] 初始化 `mobile/` 工程与 Capacitor Android 壳。
- [ ] 设计移动端导航、主题和基础状态管理。
- [ ] 实现登录与会话持久化。
- [ ] 实现首页、持仓、自选、提醒页面。
- [ ] 实现股票详情页与图表区域移动端布局。
- [ ] 实现移动端股票助手聊天体验与流式输出。
- [ ] 实现设置页与 API 地址配置。
- [ ] 补充移动端 README、开发文档与构建说明。
- [ ] 运行移动端相关测试与安卓构建验证。
- [ ] 完成后提交代码并记录剩余风险。

## 测试门槛

- 移动端页面渲染测试通过。
- 移动端状态管理与 API 调用测试通过。
- 现有前端测试与类型检查继续通过。
- 如修改后端接口契约，相关 Go 测试同步通过。
- Android debug 构建成功。
- 真机或模拟器完成一次登录、查看持仓、查看详情、使用股票助手的冒烟验证。

## 完成标准

- 仓库中存在可运行的 `mobile/` 工程。
- 安卓 App 可以安装、登录并连接现有后端。
- 用户可以完成移动端高频核心操作。
- 文档包含移动端开发与构建说明。
- 测试与构建验证结果记录完整。

## 交付说明

- 第一阶段优先保证“可安装、可登录、可查看、可问答”，而不是一次性覆盖桌面端全部功能。
- 移动端布局必须以单列竖屏为主。
- 后续如需要原生推送、生物识别或更强性能图表，可在本计划完成后继续拆分子计划。
