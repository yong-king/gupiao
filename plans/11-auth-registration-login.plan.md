# 11 Auth Registration Login Plan

## Status

completed

## Objective

实现用户注册、登录和受保护 API，避免未登录用户直接访问业务功能。

## Scope

- 注册 API。
- 登录 API。
- 密码哈希和盐。
- Session token。
- 业务 API Bearer token 校验。
- 前端登录状态模型。

## Out Of Scope

- 不实现 OAuth。
- 不实现短信验证码。
- 不实现多租户权限系统。

## Dependencies

- [10 Configuration Management](</Users/youngking/Documents/jijin/plans/10-configuration-management.plan.md>)

## Required Specs

- Auth Registration Login Spec
- Credential Encryption Spec

## Tasks

- [x] 生成认证 spec。
- [x] 实现密码哈希。
- [x] 实现注册和登录。
- [x] 保护业务 API。
- [x] 前端增加登录状态模型。
- [x] 补测试。

## Testing Gate

- Go auth 单元测试通过。
- API 未登录拒绝测试通过。
- 注册登录后业务 API 测试通过。
- 前端认证状态测试通过。

## Completion Criteria

- 未登录用户不能访问业务 API。
- 密码不明文保存。
- Token 不明文落库，仓库保存 token hash。

## Delivery Notes

- Implementation files: `backend/internal/auth/`, `backend/internal/api/auth_handlers.go`, `backend/internal/api/router.go`, `frontend/src/auth.js`, `frontend/src/main.js`.
- Test files: `backend/internal/auth/auth_test.go`, `backend/internal/api/auth_handlers_test.go`, `frontend/tests/auth.test.js`.
- Test commands: `env GOCACHE=/Users/youngking/Documents/jijin/backend/.gocache go test ./...`; `npm test`; `npm run typecheck`.
- Test result: passed.
- Remaining risks: auth persistence is in-memory until database-backed repositories are wired; OAuth and password reset are not implemented.
