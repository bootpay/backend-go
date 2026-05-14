package bootpay

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestLegacyNewAPICompatibility(t *testing.T) {
	api := NewAPI("legacy_application_id", "legacy_private_key", nil, "development")

	if api.applicationId != "legacy_application_id" {
		t.Fatalf("applicationId mismatch: %s", api.applicationId)
	}
	if api.privateKey != "legacy_private_key" {
		t.Fatalf("privateKey mismatch: %s", api.privateKey)
	}
	if api.clientKey != "" || api.secretKey != "" {
		t.Fatalf("legacy API should not set client/secret keys")
	}
	if !strings.HasPrefix(api.baseUrl, DEVELOPMENT) {
		t.Fatalf("development base URL mismatch: %s", api.baseUrl)
	}
}

func TestLegacyNewRequestUsesBearerTokenWhenPresent(t *testing.T) {
	api := NewAPI("legacy_application_id", "legacy_private_key", nil, "production")
	api.token = "legacy_access_token"

	req, err := api.NewRequest(http.MethodGet, "receipt/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	if got := req.Header.Get("Authorization"); got != "Bearer legacy_access_token" {
		t.Fatalf("Authorization header mismatch: %s", got)
	}
}

func TestLegacyGetTokenOmitsAuthorizationAndUsesLegacyPayload(t *testing.T) {
	var capturedAuth string
	var capturedPayload RestConfig
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			capturedAuth = req.Header.Get("Authorization")
			if err := json.NewDecoder(req.Body).Decode(&capturedPayload); err != nil {
				t.Fatalf("decode payload failed: %v", err)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"legacy_access_token"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}
	api := NewAPI("legacy_application_id", "legacy_private_key", client, "development")

	result, err := api.GetToken()
	if err != nil {
		t.Fatal(err)
	}

	if capturedAuth != "" {
		t.Fatalf("legacy token request should not include Authorization, got %q", capturedAuth)
	}
	if capturedPayload.ApplicationId != "legacy_application_id" || capturedPayload.PrivateKey != "legacy_private_key" {
		t.Fatalf("legacy payload mismatch: %+v", capturedPayload)
	}
	if capturedPayload.ClientKey != "" || capturedPayload.SecretKey != "" {
		t.Fatalf("legacy payload should not include client/secret keys: %+v", capturedPayload)
	}
	if api.token != "legacy_access_token" {
		t.Fatalf("legacy token was not stored: %q", api.token)
	}
	// http_status 는 deprecated 이지만 backward-compat 보장 차원에서 응답 map 에 계속 포함되어야 한다.
	if result["http_status"] != http.StatusOK {
		t.Fatalf("http_status should be preserved for compatibility: %+v", result)
	}
}

func TestClientKeyGetTokenSkipsHTTPCallAndReturnsSyntheticResponse(t *testing.T) {
	// client_key/secret_key 모드에서는 매 요청에 Basic Auth 헤더가 자동 부착되므로
	// request/token 호출이 불필요하다. GetToken() 은 네트워크 호출 없이 합성 응답을 돌려준다.
	called := false
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			called = true
			return nil, nil
		}),
	}
	api := NewAPIWithClientKey("ck", "sk", client, "development")

	result, err := api.GetToken()
	if err != nil {
		t.Fatal(err)
	}
	if called {
		t.Fatal("ck/sk 경로에서는 request/token HTTP 호출이 발생하면 안 된다")
	}
	if result["access_token"] != "" {
		t.Fatalf("synthetic access_token should be empty string, got %+v", result["access_token"])
	}
	if result["http_status"] != http.StatusOK {
		t.Fatalf("synthetic http_status should be 200, got %+v", result["http_status"])
	}
	if api.token != "" {
		t.Fatalf("ck/sk 경로의 token 값은 사용되지 않으므로 비어 있어야 한다: %q", api.token)
	}
}

func TestGetTokenIgnoresNonStringAccessTokenWithoutPanic(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":123}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}
	api := NewAPI("legacy_application_id", "legacy_private_key", client, "development")
	api.token = "stale_access_token"

	if _, err := api.GetToken(); err != nil {
		t.Fatal(err)
	}
	if api.token != "" {
		t.Fatalf("non-string access_token should clear stale token: %q", api.token)
	}
}

func TestGetTokenReturnsInvalidJSONError(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`not-json`)),
				Header:     make(http.Header),
			}, nil
		}),
	}
	api := NewAPI("legacy_application_id", "legacy_private_key", client, "development")

	if _, err := api.GetToken(); err == nil {
		t.Fatal("expected invalid JSON error")
	}
}
