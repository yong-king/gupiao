package api

import (
	"encoding/json"
	"net/http"
)

func handleSwaggerOpenAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	writeJSON(w, http.StatusOK, swaggerOpenAPISpec())
}

func handleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html><html lang="zh-CN"><head><meta charset="utf-8"><title>Jijin Swagger</title></head><body><h1>Jijin API Swagger</h1><p>OpenAPI JSON: <a href="/swagger/openapi.json">/swagger/openapi.json</a></p></body></html>`))
}

func swaggerOpenAPISpec() map[string]interface{} {
	// 用代码生成最小 Swagger/OpenAPI 契约，接口测试直接校验这份契约覆盖关键股票智能助手与多 Agent API。
	return map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":       "Jijin Stock Monitoring API",
			"version":     "0.1.0",
			"description": "股票监控、Multi-Agent 研究工作流和股票智能助手 API。所有投资相关输出仅用于提醒、观察和风险提示。",
		},
		"paths": map[string]interface{}{
			"/api/workflows/research/run": swaggerPath("post", "运行多 Agent 股票研究工作流", "RunWorkflowRequest", "RunWorkflowResponse"),
			"/api/workflows":              swaggerPath("get", "查询多 Agent 工作流记录", "", "WorkflowJobList"),
			"/api/assistant/chat":         swaggerPath("post", "股票智能助手非流式对话", "AssistantChatRequest", "AssistantChatResponse"),
			"/api/assistant/chat/stream":  swaggerPath("post", "股票智能助手 SSE 流式对话", "AssistantChatRequest", "AssistantChatStream"),
			"/api/operation-logs":         swaggerPath("get", "查询 Agent、RAG 和行情采集操作日志", "", "OperationLogList"),
			"/api/delivery/settings":      swaggerPath("post", "保存提醒投递配置", "DeliverySettingsRequest", "DeliverySettings"),
			"/api/delivery/test":          swaggerPath("post", "发送提醒投递测试消息", "DeliveryTestRequest", "DeliveryTestResponse"),
			"/api/delivery/logs":          swaggerPath("get", "查询提醒投递日志", "", "DeliveryLogList"),
		},
		"components": map[string]interface{}{
			"securitySchemes": map[string]interface{}{
				"bearerAuth": map[string]interface{}{"type": "http", "scheme": "bearer"},
			},
			"schemas": swaggerSchemas(),
		},
		"security": []map[string][]string{{"bearerAuth": []string{}}},
	}
}

func swaggerPath(method string, summary string, requestSchema string, responseSchema string) map[string]interface{} {
	operation := map[string]interface{}{
		"summary": summary,
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"description": "OK",
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{"$ref": "#/components/schemas/" + responseSchema},
					},
				},
			},
		},
	}
	if requestSchema != "" {
		operation["requestBody"] = map[string]interface{}{
			"required": true,
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{"$ref": "#/components/schemas/" + requestSchema},
				},
			},
		}
	}
	return map[string]interface{}{method: operation}
}

func swaggerSchemas() map[string]interface{} {
	return map[string]interface{}{
		"RunWorkflowRequest": objectSchema(map[string]interface{}{
			"user_id":         stringSchema("用户 ID"),
			"attention_level": stringSchema("关注等级 high/medium/low"),
			"market":          stringSchema("市场，可选"),
			"symbol":          stringSchema("股票代码，可选"),
		}, []string{"user_id", "attention_level"}),
		"RunWorkflowResponse": objectSchema(map[string]interface{}{
			"job":       objectSchema(map[string]interface{}{}, nil),
			"targets":   arraySchema(objectSchema(map[string]interface{}{}, nil)),
			"documents": arraySchema(objectSchema(map[string]interface{}{}, nil)),
		}, nil),
		"WorkflowJobList": arraySchema(objectSchema(map[string]interface{}{}, nil)),
		"AssistantChatRequest": objectSchema(map[string]interface{}{
			"user_id":    stringSchema("用户 ID"),
			"session_id": stringSchema("会话 ID"),
			"market":     stringSchema("市场"),
			"symbol":     stringSchema("股票代码"),
			"question":   stringSchema("用户问题"),
		}, []string{"user_id", "symbol", "question"}),
		"AssistantChatResponse": objectSchema(map[string]interface{}{
			"market":           stringSchema("市场"),
			"symbol":           stringSchema("股票代码"),
			"session_id":       stringSchema("会话 ID"),
			"answer":           stringSchema("中文研究提醒回答"),
			"context_summary":  stringSchema("上下文摘要"),
			"rag_document_ids": arraySchema(stringSchema("RAG 文档 ID")),
			"model":            stringSchema("模型名称"),
			"provider":         stringSchema("模型提供方或降级来源"),
			"llm_status":       stringSchema("模型调用状态，例如 deepseek-api、missing-api-key、local-fallback"),
		}, nil),
		"AssistantChatStream": objectSchema(map[string]interface{}{
			"delta": stringSchema("SSE 增量片段"),
			"done":  map[string]interface{}{"type": "boolean"},
		}, nil),
		"OperationLogList": arraySchema(objectSchema(map[string]interface{}{}, nil)),
		"DeliverySettingsRequest": objectSchema(map[string]interface{}{
			"user_id":            stringSchema("用户 ID"),
			"enabled":            boolSchema("是否启用外部高提醒投递"),
			"default_channel":    stringSchema("默认投递通道 wecom_webhook/wechat_placeholder"),
			"wecom_webhook_url":  stringSchema("企业微信机器人 Webhook"),
			"wecom_mention":      stringSchema("企业微信提醒对象"),
			"wechat_target":      stringSchema("微信目标占位"),
			"wechat_note":        stringSchema("微信说明或 MCP 备注"),
			"min_risk_level":     stringSchema("最小外发风险等级"),
			"cooldown_seconds":   integerSchema("同类提醒外发冷却秒数"),
			"mcp_enabled":        boolSchema("是否启用后续 MCP 通道"),
		}, []string{"user_id"}),
		"DeliverySettings": objectSchema(map[string]interface{}{
			"user_id":            stringSchema("用户 ID"),
			"enabled":            boolSchema("是否启用外部高提醒投递"),
			"default_channel":    stringSchema("默认投递通道"),
			"wecom_webhook_url":  stringSchema("企业微信 Webhook，已配置时会脱敏显示"),
			"wecom_mention":      stringSchema("企业微信提醒对象"),
			"wechat_target":      stringSchema("微信目标占位"),
			"wechat_note":        stringSchema("微信说明或 MCP 备注"),
			"min_risk_level":     stringSchema("最小外发风险等级"),
			"cooldown_seconds":   integerSchema("同类提醒外发冷却秒数"),
			"mcp_enabled":        boolSchema("是否启用后续 MCP 通道"),
			"last_test_status":   stringSchema("最近一次测试状态"),
			"last_test_message":  stringSchema("最近一次测试结果说明"),
			"last_tested_at":     stringSchema("最近一次测试时间"),
		}, []string{"user_id"}),
		"DeliveryTestRequest": objectSchema(map[string]interface{}{
			"user_id": stringSchema("用户 ID"),
		}, []string{"user_id"}),
		"DeliveryTestResponse": objectSchema(map[string]interface{}{
			"status":  stringSchema("发送状态"),
			"message": stringSchema("发送结果说明"),
		}, nil),
		"DeliveryLog": objectSchema(map[string]interface{}{
			"id":               stringSchema("日志 ID"),
			"user_id":          stringSchema("用户 ID"),
			"alert_event_id":   stringSchema("关联提醒事件 ID"),
			"channel_type":     stringSchema("通道类型"),
			"channel_target":   stringSchema("通道目标摘要"),
			"request_summary":  stringSchema("请求摘要"),
			"response_summary": stringSchema("响应摘要"),
			"status":           stringSchema("sent/failed/skipped"),
			"error_message":    stringSchema("错误信息"),
			"created_at":       stringSchema("创建时间"),
		}, nil),
		"DeliveryLogList": arraySchema(map[string]interface{}{"$ref": "#/components/schemas/DeliveryLog"}),
	}
}

func objectSchema(properties map[string]interface{}, required []string) map[string]interface{} {
	schema := map[string]interface{}{"type": "object", "properties": properties}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func arraySchema(item interface{}) map[string]interface{} {
	return map[string]interface{}{"type": "array", "items": item}
}

func stringSchema(description string) map[string]interface{} {
	return map[string]interface{}{"type": "string", "description": description}
}

func boolSchema(description string) map[string]interface{} {
	return map[string]interface{}{"type": "boolean", "description": description}
}

func integerSchema(description string) map[string]interface{} {
	return map[string]interface{}{"type": "integer", "description": description}
}

func swaggerOpenAPIJSON() ([]byte, error) {
	return json.Marshal(swaggerOpenAPISpec())
}
