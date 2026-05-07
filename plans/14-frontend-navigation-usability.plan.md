# 14 Frontend Navigation Usability Plan

## Status

completed

## Objective

修复登录页不居中和左侧导航不可用的问题，让前端操作台达到可点击、可切换、可验证的 MVP 使用状态。

## Scope

- 登录页居中布局。
- 已登录应用布局与侧边栏 active 状态。
- 左侧导航点击切换视图。
- 各导航项的基础内容和安全 Demo 操作。
- 认证成功后保存并复用 `user_id`。

## Out Of Scope

- 不实现自动买卖或下单能力。
- 不引入完整 Vue 构建链。
- 不为本 UI 修复强制接入 PostgreSQL 或 Redis。

## Dependencies

- [08 Frontend Console](</Users/youngking/Documents/jijin/plans/08-frontend-console.plan.md>)
- [11 Auth Registration Login](</Users/youngking/Documents/jijin/plans/11-auth-registration-login.plan.md>)

## Required Specs

- Frontend Navigation Usability Spec

## Tasks

- [x] 更新 spec，明确登录居中和导航可用性验收标准。
- [x] 更新 plan，限定本次 UI 可用性修复范围。
- [x] 修复登录页布局类和 CSS。
- [x] 实现侧边栏点击切换 active view。
- [x] 为每个导航项渲染独立页面内容。
- [x] Demo API 操作使用登录返回的 `user_id`。
- [x] 增加导航视图模型测试。
- [x] 运行相关测试并记录结果。

## Testing Gate

- `cd frontend && npm test`
- `cd frontend && npm run typecheck`
- 浏览器验证：登录页居中，登录后点击每个侧边栏按钮都能切换主内容。

## Completion Criteria

- 登录页在浏览器中居中。
- 左侧导航不再是无效按钮。
- 每个导航项都有对应主内容。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `frontend/index.html`, `frontend/src/app.js`, `frontend/src/main.js`.
- Test files: `frontend/tests/app.test.js`.
- Test commands: `npm test`; `npm run typecheck`.
- Test result: passed, 19 frontend tests.
- Browser verification: passed in Chrome at `http://127.0.0.1:5173/index.html?v=16`; auth card is centered and sidebar clicks switch the main view.
- Remaining risks: this plan keeps the static ES module MVP shell; full Vue SFC routing remains a future implementation step.
