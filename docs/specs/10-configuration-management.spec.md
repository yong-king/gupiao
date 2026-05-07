# Configuration Management Spec

## 1. Background

部署和后续接入数据库、Redis、大模型、股票数据源时，需要明确配置文件位置和环境变量覆盖规则。

## 2. Goals

- 后端配置放在 `config/backend.example.json`。
- Python Agent 配置放在 `agent/config/agent.example.json`。
- 支持环境变量覆盖敏感项。
- 文档说明哪些配置必须加密或只走环境变量。

## 3. Non-Goals

- 不提交真实 API key。
- 不实现密钥管理服务。

## 4. Functional Scope

### Must Have

- 数据库、Redis、Agent URL、LLM、股票源配置字段。
- Go JSON 配置加载。
- Python JSON 配置加载。

## 5. Security

需要保护：

- 用户密码：只存 salted password hash。
- Session token：只存 token hash。
- LLM API key：只走环境变量或服务器密钥。
- 数据库密码：部署环境变量注入。
- 券商/股票平台凭证：MVP 不保存。

## 6. Testing

- Go 配置加载测试。
- Python 配置加载测试。

## 7. Acceptance Criteria

- Given 配置文件存在 When 加载配置 Then 数据库、Redis、Agent、LLM 和股票源字段可读取。
- Given 环境变量存在 When 加载配置 Then 环境变量覆盖配置文件敏感字段。

## 8. Definition Of Done

- 配置文件、加载代码、测试和文档都存在。
