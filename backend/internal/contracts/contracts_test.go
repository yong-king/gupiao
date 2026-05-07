package contracts

import "testing"

func TestIsMarketCode(t *testing.T) {
	if !IsMarketCode("US") {
		t.Fatal("expected US to be a supported market")
	}
	if IsMarketCode("AUTO_TRADE") {
		t.Fatal("unexpected unsupported market accepted")
	}
}

func TestInitialErrorCodesIncludeSafetyCodes(t *testing.T) {
	required := map[string]bool{
		"rate_limited":      false,
		"agent_error":       false,
		"insufficient_data": false,
	}

	for _, code := range ErrorCodes {
		if _, ok := required[code]; ok {
			required[code] = true
		}
	}

	for code, found := range required {
		if !found {
			t.Fatalf("missing required error code %q", code)
		}
	}
}
