# Jijin Stock Monitoring Agent

Jijin 是一个股票监控、分析和告警系统。它围绕用户提供的股票代码、自选列表、手动导入持仓或只读账户数据进行监控，并生成可解释的分析结果与告警，供人工判断。

系统不下单，也不应实现自动买卖或自动提交订单功能。

## 目录结构

```text
backend/   Go API、任务、规则、存储、通知、审计日志
agent/     Python 分析服务
frontend/  前端控制台
deploy/    本地部署与 Docker 相关文件
docs/      规格、契约、开发说明
plans/     各模块执行计划
```

## 环境要求

本地启动前请先准备：

- Docker Desktop（或可用的 Docker Engine，且支持 `docker compose`）
- Go
- Python 3
- Node.js 和 npm

如果你只是想先把整套项目跑起来，优先使用 Docker 方式即可。

## 推荐启动方式：Docker Compose

在仓库根目录执行：

```bash
docker compose -f deploy/docker-compose.yml up -d --build
```

首次启动后，执行数据库迁移：

```bash
sh deploy/scripts/apply-migrations.sh
```

启动完成后可访问：

- 前端: `http://127.0.0.1:8081`
- 后端: `http://127.0.0.1:8080`
- Agent: `http://127.0.0.1:8090`

健康检查：

```bash
curl http://127.0.0.1:8080/healthz
```

查看服务状态：

```bash
docker compose -f deploy/docker-compose.yml ps
```

查看日志：

```bash
docker compose -f deploy/docker-compose.yml logs -f
```

停止服务：

```bash
docker compose -f deploy/docker-compose.yml down
```

如果希望连数据库数据卷一起清理：

```bash
docker compose -f deploy/docker-compose.yml down -v
```

## 分服务本地启动

如果你不想用 Docker 跑应用服务，也可以分别启动后端、Agent 和前端检查命令。

### 1. 启动依赖中间件

先启动 PostgreSQL 和 Redis：

```bash
docker compose -f deploy/docker-compose.yml up -d postgres redis
```

然后执行数据库迁移：

```bash
sh deploy/scripts/apply-migrations.sh
```

### 2. 启动后端

```bash
cd backend
env GOCACHE=$PWD/.gocache go run ./cmd/server
```

默认监听 `:8080`。如需修改，可使用环境变量 `GO_BACKEND_ADDR`。

后端默认配置样例见：

- `config/backend.example.json`

### 3. 启动 Python Agent

```bash
cd agent
PYTHONPATH=src python3 -m agent_core.server
```

默认监听 `127.0.0.1:8090`。如需修改，可使用环境变量 `AGENT_HOST` 和 `AGENT_PORT`。

Agent 配置样例见：

- `agent/config/agent.example.json`

### 4. 启动前端

当前前端是静态页面，可以直接打开：

```text
/Users/youngking/Documents/jijin/frontend/index.html
```

如果你已经通过 Docker Compose 启动了整套服务，也可以直接访问：

```text
http://127.0.0.1:8081
```

## 常用检查命令

后端测试：

```bash
cd backend
env GOCACHE=$PWD/.gocache go test ./...
```

Agent 测试：

```bash
cd agent
PYTHONPATH=src python3 -m unittest discover -s tests
```

前端测试：

```bash
cd frontend
npm test
npm run typecheck
```

完整回归检查：

```bash
cd backend
env GOCACHE=$PWD/.gocache go test ./...
cd ../agent
PYTHONPATH=src python3 -m unittest discover -s tests
cd ../frontend
npm test
npm run typecheck
cd ..
docker compose -f deploy/docker-compose.yml config
```

## 相关文件

- 全局 Agent 说明：`agent.md`
- 规范说明：`spec.md`
- 计划索引：`plan.md`
