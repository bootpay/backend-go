package bootpay

import "testing"

func TestPgGetToken(t *testing.T) {
	// AUTH_MODE=new 는 합성 응답(no-network)이지만 legacy 토글 시 실서버를 치므로 게이트
	skipUnlessDevelopment(t)
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

func TestPgGetTokenCkSkLeavesInternalTokenEmpty(t *testing.T) {
	// ck/sk 모드에서는 매 요청 Basic Auth 로 인증하므로 GetToken 은 no-op 합성 응답을 반환하고
	// 내부 token 은 비어 있는 상태를 유지한다. AUTH_MODE 토글과 무관하게 항상 ck/sk 인스턴스를 만든다.
	clientKey, secretKey := GetPgClientKeys()
	mode := ""
	if getEnv() == "development" {
		mode = "development"
	}
	api := NewAPIWithClientKey(clientKey, secretKey, nil, mode)

	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	if api.token != "" {
		t.Errorf("ck/sk 모드의 GetToken 이후 internal token 은 비어 있어야 한다, got %q", api.token)
	}
}

func TestPgGetTokenLegacyIssuesRealAccessToken(t *testing.T) {
	skipUnlessDevelopment(t)
	// legacy application_id/private_key 모드는 실제 request/token 호출 후 토큰을 발급받는다.
	appId, privateKey := GetPgKeys()
	mode := ""
	if getEnv() == "development" {
		mode = "development"
	}
	api := NewAPI(appId, privateKey, nil, mode)

	result, err := api.GetToken()
	if err != nil {
		t.Fatalf("legacy GetToken failed: %v", err)
	}

	token, ok := result["access_token"].(string)
	if !ok || token == "" {
		t.Fatalf("expected non-empty access_token, got %#v", result["access_token"])
	}
	if api.token != token {
		t.Errorf("internal token 이 발급된 access_token 과 일치하지 않음: api.token=%q, result=%q", api.token, token)
	}
	if expire, ok := result["expire_in"].(float64); !ok || expire <= 0 {
		t.Errorf("expected positive expire_in, got %#v", result["expire_in"])
	}
}
