package contracts

var MarketCodes = []string{"US", "HK", "CN"}

var ErrorCodes = []string{
	"validation_error",
	"not_found",
	"rate_limited",
	"data_source_error",
	"agent_error",
	"insufficient_data",
	"conflict",
	"internal_error",
}

func IsMarketCode(code string) bool {
	for _, marketCode := range MarketCodes {
		if marketCode == code {
			return true
		}
	}
	return false
}
