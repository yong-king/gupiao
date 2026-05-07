# Auth Registration Login Spec

## 1. Background

系统需要注册和登录，否则业务 API 对未登录用户开放。

## 2. Goals

- 注册用户。
- 登录返回 session token。
- 业务 API 需要 Bearer token。
- 密码和 token 不能明文保存。

## 3. Non-Goals

- 不实现 OAuth。
- 不实现短信验证码。
- 不实现复杂 RBAC。

## 4. Functional Scope

### Must Have

- `POST /api/auth/register`
- `POST /api/auth/login`
- Password PBKDF2-HMAC-SHA256 salted hash。
- Session token hash。
- Auth middleware。

## 5. Encryption And Hashing

- 用户密码：PBKDF2-HMAC-SHA256 + random salt。
- Session token：返回明文 token 给客户端，服务端只保存 SHA-256 hash。
- LLM key、数据库密码：不进入代码仓库，通过环境变量或部署密钥注入。

## 6. Testing

- 注册成功。
- 密码验证。
- 登录失败。
- 未登录访问业务 API 返回 401。
- 登录后访问业务 API 成功。

## 7. Acceptance Criteria

- Given 未登录用户 When 访问业务 API Then 返回 401。
- Given 用户登录成功 When 带 Bearer token 访问业务 API Then 允许访问。

## 8. Definition Of Done

- Go 和前端认证测试通过。
