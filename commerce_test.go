package bootpay

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 테스트용 API 서버로 요청이 전달되도록 Api를 구성한다
func newTestApi(handler http.HandlerFunc) (*Api, *httptest.Server) {
	server := httptest.NewServer(handler)
	api := &Api{
		clientKey: "test_client_key",
		secretKey: "test_secret_key",
		baseUrl:   server.URL,
		client:    server.Client(),
	}
	return api, server
}

func decodeRequestBody(t *testing.T, req *http.Request) map[string]interface{} {
	t.Helper()
	raw, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("request body read error: %s", err.Error())
	}
	payload := map[string]interface{}{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("request body decode error: %s (body: %s)", err.Error(), string(raw))
	}
	return payload
}

func TestLookupSequentialBillingKey(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"billing_key":"62afc52dcf9f6d001d7d1035"}`))
	})
	defer server.Close()

	result, err := api.LookupSequentialBillingKey("widget_1234", "62afc52dcf9f6d001d7d1035")
	if err != nil {
		t.Fatalf("LookupSequentialBillingKey error: %s", err.Error())
	}
	if request.Method != http.MethodGet {
		t.Errorf("method = %s, want %s", request.Method, http.MethodGet)
	}
	if request.URL.Path != "/subscribe/sequential_billing_key/62afc52dcf9f6d001d7d1035" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if widgetKey := request.URL.Query().Get("widget_key"); widgetKey != "widget_1234" {
		t.Errorf("widget_key = %s, want widget_1234", widgetKey)
	}
	if result["billing_key"] != "62afc52dcf9f6d001d7d1035" {
		t.Errorf("billing_key = %v", result["billing_key"])
	}
}

func TestGetMallSetting(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"name":"부트페이 몰"}`))
	})
	defer server.Close()

	result, err := api.GetMallSetting("")
	if err != nil {
		t.Fatalf("GetMallSetting error: %s", err.Error())
	}
	if request.Method != http.MethodGet {
		t.Errorf("method = %s, want %s", request.Method, http.MethodGet)
	}
	if request.URL.Path != "/mall-setting" {
		t.Errorf("path = %s, want /mall-setting", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "supervisor" {
		t.Errorf("Bootpay-Role = %s, want supervisor", role)
	}
	if idempotencyKey := request.Header.Get("Idempotency-Key"); idempotencyKey == "" {
		t.Error("Idempotency-Key header is empty")
	}
	if result["name"] != "부트페이 몰" {
		t.Errorf("name = %v", result["name"])
	}
}

func TestUpdateMallSetting(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.UpdateMallSetting(APIResponse{
		"name":        "부트페이 몰",
		"use_logo":    false,
		"point_rate":  0,
		"description": nil,
	}, "0d3d4bd0-3f8a-4d94-93d8-1c1e1f9f0c11")
	if err != nil {
		t.Fatalf("UpdateMallSetting error: %s", err.Error())
	}
	if request.Method != http.MethodPut {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPut)
	}
	if request.URL.Path != "/mall-setting" {
		t.Errorf("path = %s, want /mall-setting", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "supervisor" {
		t.Errorf("Bootpay-Role = %s, want supervisor", role)
	}
	if idempotencyKey := request.Header.Get("Idempotency-Key"); idempotencyKey != "0d3d4bd0-3f8a-4d94-93d8-1c1e1f9f0c11" {
		t.Errorf("Idempotency-Key = %s", idempotencyKey)
	}
	if payload["name"] != "부트페이 몰" {
		t.Errorf("name = %v", payload["name"])
	}
	if payload["use_logo"] != false {
		t.Errorf("use_logo = %v, want false", payload["use_logo"])
	}
	if payload["point_rate"] != float64(0) {
		t.Errorf("point_rate = %v, want 0", payload["point_rate"])
	}
	if _, ok := payload["description"]; ok {
		t.Error("description is nil, so it must not be sent")
	}
}

func TestSupervisorRequestOrderSubscriptionCharge(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"receipt_id":"62afc3c5cf9f6d001b7d101a"}`))
	})
	defer server.Close()

	result, err := api.SupervisorRequestOrderSubscriptionCharge(OrderSubscriptionChargePayload{
		ChargeKey: "charge_1234",
		Price:     1000,
		User:      APIResponse{"id": "user_1234"},
	})
	if err != nil {
		t.Fatalf("SupervisorRequestOrderSubscriptionCharge error: %s", err.Error())
	}
	if request.Method != http.MethodPost {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPost)
	}
	if request.URL.Path != "/order_subscriptions/charge" {
		t.Errorf("path = %s, want /order_subscriptions/charge", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "supervisor" {
		t.Errorf("Bootpay-Role = %s, want supervisor", role)
	}
	if idempotencyKey := request.Header.Get("Idempotency-Key"); idempotencyKey == "" {
		t.Error("Idempotency-Key header is empty")
	}
	// charge_key는 body로만 전송한다 (URL/query 금지)
	if strings.Contains(request.URL.RequestURI(), "charge_1234") {
		t.Errorf("charge_key must not be sent by url: %s", request.URL.RequestURI())
	}
	if payload["charge_key"] != "charge_1234" {
		t.Errorf("charge_key = %v", payload["charge_key"])
	}
	if payload["price"] != float64(1000) {
		t.Errorf("price = %v, want 1000", payload["price"])
	}
	if _, ok := payload["tax_free_price"]; ok {
		t.Error("tax_free_price is empty, so it must not be sent")
	}
	if _, ok := payload["metadata"]; ok {
		t.Error("metadata is empty, so it must not be sent")
	}
	if result["receipt_id"] != "62afc3c5cf9f6d001b7d101a" {
		t.Errorf("receipt_id = %v", result["receipt_id"])
	}
}

func TestSupervisorRequestOrderSubscriptionChargeRevoke(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.SupervisorRequestOrderSubscriptionChargeRevoke(OrderSubscriptionChargeRevokePayload{
		IdempotencyKey: "5f0f0cd6-1f45-4a5f-9d0f-2f1a1b0c8d20",
		ChargeKey:      "charge_1234",
	})
	if err != nil {
		t.Fatalf("SupervisorRequestOrderSubscriptionChargeRevoke error: %s", err.Error())
	}
	if request.Method != http.MethodDelete {
		t.Errorf("method = %s, want %s", request.Method, http.MethodDelete)
	}
	if request.URL.Path != "/order_subscriptions/charge" {
		t.Errorf("path = %s, want /order_subscriptions/charge", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "supervisor" {
		t.Errorf("Bootpay-Role = %s, want supervisor", role)
	}
	if idempotencyKey := request.Header.Get("Idempotency-Key"); idempotencyKey != "5f0f0cd6-1f45-4a5f-9d0f-2f1a1b0c8d20" {
		t.Errorf("Idempotency-Key = %s", idempotencyKey)
	}
	if payload["charge_key"] != "charge_1234" {
		t.Errorf("charge_key = %v", payload["charge_key"])
	}
	if _, ok := payload["user"]; ok {
		t.Error("user is empty, so it must not be sent")
	}
}
