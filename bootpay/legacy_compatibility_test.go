package bootpay

import (
	"encoding/base64"
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
	if api.clientKey != "" || api.serverKey != "" {
		t.Fatalf("legacy API should not set client/server keys")
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
	if capturedPayload.ClientKey != "" || capturedPayload.ServerKey != "" {
		t.Fatalf("legacy payload should not include client/server keys: %+v", capturedPayload)
	}
	if api.token != "legacy_access_token" {
		t.Fatalf("legacy token was not stored: %q", api.token)
	}
	// http_status 는 deprecated 이지만 backward-compat 보장 차원에서 응답 map 에 계속 포함되어야 한다.
	if result["http_status"] != http.StatusOK {
		t.Fatalf("http_status should be preserved for compatibility: %+v", result)
	}
}

func TestClientKeyGetTokenUsesBasicAuthorizationAndClientPayload(t *testing.T) {
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
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"client_key_access_token"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}
	api := NewAPIWithClientKey("ck", "sk", client, "development")

	if _, err := api.GetToken(); err != nil {
		t.Fatal(err)
	}

	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("ck:sk"))
	if capturedAuth != expectedAuth {
		t.Fatalf("Basic Authorization mismatch: got %q want %q", capturedAuth, expectedAuth)
	}
	if capturedPayload.ClientKey != "ck" || capturedPayload.ServerKey != "sk" {
		t.Fatalf("client payload mismatch: %+v", capturedPayload)
	}
	if capturedPayload.ApplicationId != "" || capturedPayload.PrivateKey != "" {
		t.Fatalf("client payload should not include legacy keys: %+v", capturedPayload)
	}
	if api.token != "client_key_access_token" {
		t.Fatalf("client key token was not stored: %q", api.token)
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
