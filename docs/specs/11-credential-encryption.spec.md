# Credential Encryption Spec

## 1. Background

系统涉及用户密码、session token、LLM key、数据库密码和潜在平台凭证，需要明确保护方式。

## 2. Rules

- 密码永不明文保存。
- Session token 服务端只保存 hash。
- API key 和数据库密码不提交到仓库。
- 券商交易密码不支持保存。
- 后续如需平台只读凭证，必须使用服务器侧加密或密钥管理服务。

## 3. Testing

- Password hash 不包含原始密码。
- Token repository 不以明文 token 作为 key。

## 4. Definition Of Done

- 认证包测试覆盖这些规则。
