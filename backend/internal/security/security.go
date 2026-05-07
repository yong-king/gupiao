package security

import "strings"

var sensitiveTerms = []string{"password", "secret", "token", "private_key"}
var tradingTerms = []string{"order", "place_order", "buy_order", "sell_order", "auto_trade", "自动下单", "买入下单", "卖出下单"}

func ContainsSensitiveKey(keys []string) bool {
	for _, key := range keys {
		normalized := strings.ToLower(key)
		for _, term := range sensitiveTerms {
			if strings.Contains(normalized, term) {
				return true
			}
		}
	}
	return false
}

func ContainsTradingAction(labels []string) bool {
	for _, label := range labels {
		normalized := strings.ToLower(label)
		for _, term := range tradingTerms {
			if strings.Contains(normalized, term) {
				return true
			}
		}
	}
	return false
}
