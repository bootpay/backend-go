package bootpay

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestPgClientKeyCredentialsUseBasicAuth(t *testing.T) {
	api := NewAPIWithClientKey("client_key", "secret_key", nil, "production")
	api.token = "ignored_token"

	req, err := api.NewRequest(http.MethodGet, "/receipt/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	want := "Basic " + base64.StdEncoding.EncodeToString([]byte("client_key:secret_key"))
	if got := req.Header.Get("Authorization"); got != want {
		t.Fatalf("Authorization header mismatch: got %q want %q", got, want)
	}
}

func TestPgPartialCredentialsFailBeforeNetworkRequest(t *testing.T) {
	tests := []struct {
		name string
		api  *Api
	}{
		{name: "client_key only", api: NewAPIWithClientKey("client_key", "", nil, "production")},
		{name: "secret_key only", api: NewAPIWithClientKey("", "secret_key", nil, "production")},
		{name: "application_id only", api: NewAPI("application_id", "", nil, "production")},
		{name: "private_key only", api: NewAPI("", "private_key", nil, "production")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			tt.api.client = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				called = true
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{}`)),
					Header:     make(http.Header),
				}, nil
			})}

			if _, err := tt.api.GetToken(); err == nil {
				t.Fatal("expected incomplete credentials error")
			}
			if called {
				t.Fatal("network request must not run with incomplete credentials")
			}
		})
	}
}

func TestCommerceAlwaysUsesBasicAuthEvenWhenTokenIsStored(t *testing.T) {
	api := NewCommerceAPI("client_key", "secret_key", nil, "production")
	api.SetToken("stored_token")

	req, err := api.newRequest(http.MethodGet, "users", nil)
	if err != nil {
		t.Fatal(err)
	}
	want := "Basic " + base64.StdEncoding.EncodeToString([]byte("client_key:secret_key"))
	if got := req.Header.Get("Authorization"); got != want {
		t.Fatalf("Commerce Authorization must remain Basic: got %q want %q", got, want)
	}
}

func TestCommercePartialCredentialsFailBeforeNetworkRequest(t *testing.T) {
	tests := []struct {
		name      string
		clientKey string
		secretKey string
	}{
		{name: "client_key only", clientKey: "client_key"},
		{name: "secret_key only", secretKey: "secret_key"},
		{name: "both missing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				called = true
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{}`)),
					Header:     make(http.Header),
				}, nil
			})}
			api := NewCommerceAPI(tt.clientKey, tt.secretKey, client, "production")

			if _, err := api.Get("users"); err == nil {
				t.Fatal("expected incomplete Commerce credentials error")
			}
			if called {
				t.Fatal("network request must not run with incomplete Commerce credentials")
			}
		})
	}
}
