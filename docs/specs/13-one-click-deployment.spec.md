# One Click Deployment Spec

## 1. Background

项目需要后续能一键部署到服务器，启动 PostgreSQL、Redis、后端、Agent 和前端。

## 2. Goals

- 提供部署脚本。
- 校验 Docker 和 Compose。
- 使用 `.env` 注入部署配置。

## 3. Non-Goals

- 不自动购买服务器。
- 不实现 Kubernetes。

## 4. Functional Scope

### Must Have

- `deploy/scripts/deploy.sh`
- `docker compose config` 校验。
- `docker compose up -d --build`

## 5. Testing

- Shell 语法检查。
- Compose config 校验。

## 6. Acceptance Criteria

- Given Docker 可用 When 执行部署脚本 Then Compose 服务启动。

## 7. Definition Of Done

- 脚本和文档存在，静态测试通过。
