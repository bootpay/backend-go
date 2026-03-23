package bootpay

import "testing"

func TestCommerceGetAccessTokenIntegration(t *testing.T) {
	commerce := CreateCommerceApi()

	result, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	t.Logf("Commerce GetAccessToken response: %+v", result)

	if result["access_token"] == nil {
		t.Error("Expected access_token in Commerce response, got nil")
	}
}

func TestCommerceGetAccessTokenSetsInternalToken(t *testing.T) {
	commerce := CreateCommerceApi()

	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	if commerce.GetToken() == "" {
		t.Error("Expected internal token to be set after GetAccessToken, but it is empty")
	}
}

func TestCommerceRoleChaining(t *testing.T) {
	commerce := CreateCommerceApi()

	commerce.AsUser()
	if commerce.GetRole() != "user" {
		t.Errorf("Expected role 'user', got '%s'", commerce.GetRole())
	}

	commerce.AsManager()
	if commerce.GetRole() != "manager" {
		t.Errorf("Expected role 'manager', got '%s'", commerce.GetRole())
	}

	commerce.AsPartner()
	if commerce.GetRole() != "partner" {
		t.Errorf("Expected role 'partner', got '%s'", commerce.GetRole())
	}

	commerce.AsVendor()
	if commerce.GetRole() != "vendor" {
		t.Errorf("Expected role 'vendor', got '%s'", commerce.GetRole())
	}

	commerce.AsSupervisor()
	if commerce.GetRole() != "supervisor" {
		t.Errorf("Expected role 'supervisor', got '%s'", commerce.GetRole())
	}
}
