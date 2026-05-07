# 13 One Click Deployment Plan

## Status

completed

## Objective

提供服务器一键部署入口，基于 Docker Compose 启动 PostgreSQL、Redis、后端、Agent 和前端。

## Scope

- `deploy/scripts/deploy.sh`
- `.env.example` 部署变量说明。
- Compose 健康检查。
- README 部署命令。

## Out Of Scope

- 不实现 Kubernetes。
- 不自动购买云服务器。

## Dependencies

- [09 Hardening Release](</Users/youngking/Documents/jijin/plans/09-hardening-release.plan.md>)
- [10 Configuration Management](</Users/youngking/Documents/jijin/plans/10-configuration-management.plan.md>)

## Required Specs

- One Click Deployment Spec

## Tasks

- [x] 生成部署 spec。
- [x] 添加部署脚本。
- [x] 校验 compose。
- [x] 更新文档。

## Testing Gate

- Shell 语法检查通过。
- Docker Compose config 通过。

## Completion Criteria

- 用户可以用一条命令启动部署流程。
- 测试通过并记录。

## Delivery Notes

- Implementation files: `deploy/scripts/deploy.sh`, `docs/release.md`, `README.md`.
- Test files: none.
- Test commands: `sh -n deploy/scripts/deploy.sh`; `docker compose -f deploy/docker-compose.yml config`.
- Test result: passed.
- Remaining risks: containers were not started or image-built in this turn; server deployment still requires Docker runtime and network access for images.
