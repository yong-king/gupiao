package security

import "testing"

func TestContainsSensitiveKey(t *testing.T) {
	if !ContainsSensitiveKey([]string{"api_token"}) {
		t.Fatal("expected sensitive key")
	}
	if ContainsSensitiveKey([]string{"alias"}) {
		t.Fatal("did not expect sensitive key")
	}
}

func TestContainsTradingAction(t *testing.T) {
	if !ContainsTradingAction([]string{"place_order"}) {
		t.Fatal("expected trading action")
	}
	if ContainsTradingAction([]string{"buy_watch"}) {
		t.Fatal("watch signal should not be trading action")
	}
}
