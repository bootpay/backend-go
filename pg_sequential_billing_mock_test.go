package bootpay

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// TestLookupSequentialBillingKeyRequestFormat mock 검증 —
// GET subscribe/sequential_billing_key/{billing_key}?widget_key=&user_id= (QueryEscape 포함)
func TestLookupSequentialBillingKeyRequestFormat(t *testing.T) {
	var capturedURL string
	var capturedMethod string
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			capturedURL = req.URL.String()
			capturedMethod = req.Method
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"billing_key":"bk_1"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}
	api := NewAPIWithClientKey("ck", "sk", client, "development")

	result, err := api.LookupSequentialBillingKey("widget key/1", "bk_1", "user id&2")
	if err != nil {
		t.Fatal(err)
	}

	if capturedMethod != http.MethodGet {
		t.Fatalf("method mismatch: %s", capturedMethod)
	}
	expected := DEVELOPMENT + "/subscribe/sequential_billing_key/bk_1?widget_key=widget+key%2F1&user_id=user+id%262"
	if capturedURL != expected {
		t.Fatalf("URL mismatch:\n got %s\nwant %s", capturedURL, expected)
	}
	if result["billing_key"] != "bk_1" {
		t.Fatalf("response decode mismatch: %+v", result)
	}
}
