# Stock Source Collection Spec

## 1. Background

股票平台信息可以来自官方 API、授权数据源或公开页面解析。MVP 先定义可替换接口和解析器。

## 2. Goals

- 定义股票平台数据源接口。
- 支持 JSON 和简单 HTML 解析。
- 保留 source 和 data_time。
- 明确爬虫边界。

## 3. Non-Goals

- 不绕过登录、验证码或反爬。
- 不抓取需要授权的账户页面。
- 不高频访问。

## 4. Functional Scope

### Must Have

- `Collector` 接口。
- JSON parser。
- HTML parser。
- Mock collector。

## 5. Testing

- JSON 解析。
- HTML 解析。
- 缺失字段失败。

## 6. Acceptance Criteria

- Given 平台 JSON When 解析 Then 生成行情 snapshot。
- Given 公开 HTML 包含字段 When 解析 Then 生成行情 snapshot。

## 7. Definition Of Done

- Go 测试通过。
