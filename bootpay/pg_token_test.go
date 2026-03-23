package bootpay

import "testing"

func TestPgGetToken(t *testing.T) {
	api := CreatePgApi()

	result, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	t.Logf("GetToken response: %+v", result)

	if result["access_token"] == nil {
		t.Error("Expected access_token in response, got nil")
	}
}

func TestPgGetTokenSetsInternalToken(t *testing.T) {
	api := CreatePgApi()

	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	if api.token == "" {
		t.Error("Expected internal token to be set after GetToken, but it is empty")
	}
}
