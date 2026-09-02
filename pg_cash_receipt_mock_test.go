package bootpay

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// captureCashReceiptRequest 현금영수증 요청의 method/URL/바디를 가로채는 mock 클라이언트를 만든다.
func captureCashReceiptRequest(t *testing.T, method *string, url *string, body *map[string]interface{}) *http.Client {
	t.Helper()
	return &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			*method = req.Method
			*url = req.URL.String()
			if req.Body != nil {
				if err := json.NewDecoder(req.Body).Decode(body); err != nil {
					t.Fatalf("decode payload failed: %v", err)
				}
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"receipt_id":"cr_1"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}
}

// TestRequestCashReceiptOmitsPgWhenBlank mock 검증 —
// Ruby SDK c716a1f 로 request_cash_receipt 의 pg 가 필수에서 선택(pg: nil)으로 바뀌었다.
// Go 는 pg 가 비면 키 자체를 싣지 않아 서버가 기본 PG 로 발행한다.
func TestRequestCashReceiptOmitsPgWhenBlank(t *testing.T) {
	var capturedMethod, capturedURL string
	capturedBody := map[string]interface{}{}
	client := captureCashReceiptRequest(t, &capturedMethod, &capturedURL, &capturedBody)
	api := NewAPIWithClientKey("ck", "sk", client, "development")

	result, err := api.RequestCashReceipt(CashReceiptData{
		Price:           1000,
		OrderName:       "cash receipt without pg",
		CashReceiptType: "소득공제",
		IdentityNo:      "01000000000",
		OrderId:         "order_no_pg",
		PurchasedAt:     "2026-08-31T14:50:00+09:00",
	})
	if err != nil {
		t.Fatal(err)
	}

	if capturedMethod != http.MethodPost {
		t.Fatalf("method mismatch: %s", capturedMethod)
	}
	if expected := DEVELOPMENT + "/request/cash/receipt"; capturedURL != expected {
		t.Fatalf("URL mismatch:\n got %s\nwant %s", capturedURL, expected)
	}
	if _, exists := capturedBody["pg"]; exists {
		t.Fatalf("pg must be omitted when blank: %+v", capturedBody)
	}
	// pg 만 빠지고 나머지 필드는 그대로 실려야 한다.
	if capturedBody["order_id"] != "order_no_pg" || capturedBody["cash_receipt_type"] != "소득공제" {
		t.Fatalf("payload mismatch: %+v", capturedBody)
	}
	if result["receipt_id"] != "cr_1" {
		t.Fatalf("response decode mismatch: %+v", result)
	}
}

// TestRequestCashReceiptKeepsPgWhenSet pg 를 지정하면 그대로 실린다 (기존 동작 유지).
func TestRequestCashReceiptKeepsPgWhenSet(t *testing.T) {
	var capturedMethod, capturedURL string
	capturedBody := map[string]interface{}{}
	client := captureCashReceiptRequest(t, &capturedMethod, &capturedURL, &capturedBody)
	api := NewAPIWithClientKey("ck", "sk", client, "development")

	if _, err := api.RequestCashReceipt(CashReceiptData{
		Pg:              "toss",
		Price:           1000,
		OrderName:       "cash receipt with pg",
		CashReceiptType: "소득공제",
		IdentityNo:      "01000000000",
		OrderId:         "order_with_pg",
		PurchasedAt:     "2026-08-31T14:50:00+09:00",
	}); err != nil {
		t.Fatal(err)
	}

	if capturedBody["pg"] != "toss" {
		t.Fatalf("pg must be sent when set: %+v", capturedBody)
	}
}

// TestRequestCashReceiptByBootpayOmitsPgWhenBlank 결제건 현금영수증 발행(cash_receipt_publish_on_receipt)도
// Ruby 와 동일하게 pg 가 선택값이다.
func TestRequestCashReceiptByBootpayOmitsPgWhenBlank(t *testing.T) {
	var capturedMethod, capturedURL string
	capturedBody := map[string]interface{}{}
	client := captureCashReceiptRequest(t, &capturedMethod, &capturedURL, &capturedBody)
	api := NewAPIWithClientKey("ck", "sk", client, "development")

	if _, err := api.RequestCashReceiptByBootpay(CashReceiptData{
		ReceiptId:       "62e0f11f1fc192036b1b3c92",
		Username:        "홍길동",
		Email:           "test@bootpay.co.kr",
		Phone:           "01000000000",
		IdentityNo:      "01000000000",
		CashReceiptType: "소득공제",
	}); err != nil {
		t.Fatal(err)
	}

	if capturedMethod != http.MethodPost {
		t.Fatalf("method mismatch: %s", capturedMethod)
	}
	if expected := DEVELOPMENT + "/request/receipt/cash/publish"; capturedURL != expected {
		t.Fatalf("URL mismatch:\n got %s\nwant %s", capturedURL, expected)
	}
	if _, exists := capturedBody["pg"]; exists {
		t.Fatalf("pg must be omitted when blank: %+v", capturedBody)
	}
	if capturedBody["receipt_id"] != "62e0f11f1fc192036b1b3c92" {
		t.Fatalf("payload mismatch: %+v", capturedBody)
	}
}
