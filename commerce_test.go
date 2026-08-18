package bootpay

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

	result, err := api.LookupSequentialBillingKey("widget_1234", "62afc52dcf9f6d001d7d1035", "user_1234")
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
	if userId := request.URL.Query().Get("user_id"); userId != "user_1234" {
		t.Errorf("user_id = %s, want user_1234", userId)
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

func TestInvoiceList(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"list":[{"id":"invoice_1234"}],"count":1}`))
	})
	defer server.Close()

	result, err := api.InvoiceList(APIResponse{
		"keyword": "부트페이",
		"user_id": nil,
	}, "")
	if err != nil {
		t.Fatalf("InvoiceList error: %s", err.Error())
	}
	if request.Method != http.MethodGet {
		t.Errorf("method = %s, want %s", request.Method, http.MethodGet)
	}
	if request.URL.Path != "/invoices" {
		t.Errorf("path = %s, want /invoices", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "user" {
		t.Errorf("Bootpay-Role = %s, want user", role)
	}
	query := request.URL.Query()
	if page := query.Get("page"); page != "1" {
		t.Errorf("page = %s, want 1", page)
	}
	// 서버 기본 limit은 24다
	if limit := query.Get("limit"); limit != "24" {
		t.Errorf("limit = %s, want 24", limit)
	}
	if keyword := query.Get("keyword"); keyword != "부트페이" {
		t.Errorf("keyword = %s", keyword)
	}
	if _, ok := query["user_id"]; ok {
		t.Error("user_id is nil, so it must not be sent")
	}
	if result["count"] != float64(1) {
		t.Errorf("count = %v, want 1", result["count"])
	}
}

func TestInvoiceDetail(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"id":"invoice_1234"}`))
	})
	defer server.Close()

	_, err := api.InvoiceDetail("invoice_1234", "")
	if err != nil {
		t.Fatalf("InvoiceDetail error: %s", err.Error())
	}
	if request.Method != http.MethodGet {
		t.Errorf("method = %s, want %s", request.Method, http.MethodGet)
	}
	if request.URL.Path != "/invoices/invoice_1234" {
		t.Errorf("path = %s, want /invoices/invoice_1234", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "user" {
		t.Errorf("Bootpay-Role = %s, want user", role)
	}
}

func TestInvoiceNotify(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.InvoiceNotify("invoice_1234", []string{"sms", "email"}, "")
	if err != nil {
		t.Fatalf("InvoiceNotify error: %s", err.Error())
	}
	if request.Method != http.MethodPost {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPost)
	}
	if request.URL.Path != "/invoices/invoice_1234/notify" {
		t.Errorf("path = %s, want /invoices/invoice_1234/notify", request.URL.Path)
	}
	sendTypes, ok := payload["send_types"].([]interface{})
	if !ok || len(sendTypes) != 2 || sendTypes[0] != "sms" || sendTypes[1] != "email" {
		t.Errorf("send_types = %v", payload["send_types"])
	}
}

func TestOrderCancelList(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"list":[]}`))
	})
	defer server.Close()

	_, err := api.OrderCancelList("order_1234", "", "")
	if err != nil {
		t.Fatalf("OrderCancelList error: %s", err.Error())
	}
	if request.Method != http.MethodGet {
		t.Errorf("method = %s, want %s", request.Method, http.MethodGet)
	}
	if request.URL.Path != "/order/cancel" {
		t.Errorf("path = %s, want /order/cancel", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "user" {
		t.Errorf("Bootpay-Role = %s, want user", role)
	}
	if orderNumber := request.URL.Query().Get("order_number"); orderNumber != "order_1234" {
		t.Errorf("order_number = %s, want order_1234", orderNumber)
	}
	if _, ok := request.URL.Query()["order_id"]; ok {
		t.Error("order_id is empty, so it must not be sent")
	}
}

func TestOrderCancelWithdraw(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.OrderCancelWithdraw("cancellation_1234", "")
	if err != nil {
		t.Fatalf("OrderCancelWithdraw error: %s", err.Error())
	}
	if request.Method != http.MethodPut {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPut)
	}
	if request.URL.Path != "/order/cancel/cancellation_1234/withdraw" {
		t.Errorf("path = %s, want /order/cancel/cancellation_1234/withdraw", request.URL.Path)
	}
}

func TestOrderSubscriptionUpdate(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionUpdate("subscription_1234", APIResponse{
		"order_name": "구독 상품",
		"quantity":   2,
		"phone":      nil,
	}, "")
	if err != nil {
		t.Fatalf("OrderSubscriptionUpdate error: %s", err.Error())
	}
	if request.Method != http.MethodPut {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPut)
	}
	if request.URL.Path != "/order_subscriptions/subscription_1234" {
		t.Errorf("path = %s, want /order_subscriptions/subscription_1234", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "supervisor" {
		t.Errorf("Bootpay-Role = %s, want supervisor", role)
	}
	if payload["order_name"] != "구독 상품" {
		t.Errorf("order_name = %v", payload["order_name"])
	}
	if _, ok := payload["phone"]; ok {
		t.Error("phone is nil, so it must not be sent")
	}
}

func TestOrderSubscriptionAdjustmentCreate(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionAdjustmentCreate("subscription_1234", APIResponse{
		"name":  "설치비",
		"price": 10000,
	}, "")
	if err != nil {
		t.Fatalf("OrderSubscriptionAdjustmentCreate error: %s", err.Error())
	}
	if request.Method != http.MethodPost {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPost)
	}
	if request.URL.Path != "/order_subscriptions/subscription_1234/adjustments" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if payload["price"] != float64(10000) {
		t.Errorf("price = %v, want 10000", payload["price"])
	}
	if payload["duration"] != float64(1) {
		t.Errorf("duration = %v, want 1", payload["duration"])
	}
	if payload["tax_free_price"] != float64(0) {
		t.Errorf("tax_free_price = %v, want 0", payload["tax_free_price"])
	}
	// type 미전달시 서버가 자동 판정한다
	if _, ok := payload["type"]; ok {
		t.Error("type is not given, so it must not be sent")
	}
}

func TestOrderSubscriptionAdjustmentUpdate(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionAdjustmentUpdate("subscription_1234", 3, []APIResponse{
		{"name": "할인", "price": -1000},
	}, "")
	if err != nil {
		t.Fatalf("OrderSubscriptionAdjustmentUpdate error: %s", err.Error())
	}
	if request.Method != http.MethodPut {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPut)
	}
	if request.URL.Path != "/order_subscriptions/subscription_1234/adjustments" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if payload["duration"] != float64(3) {
		t.Errorf("duration = %v, want 3", payload["duration"])
	}
	adjustments, ok := payload["adjustments"].([]interface{})
	if !ok || len(adjustments) != 1 {
		t.Errorf("adjustments = %v", payload["adjustments"])
	}
}

func TestOrderSubscriptionAdjustmentDelete(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionAdjustmentDelete("subscription_1234", "adjustment_1234", "")
	if err != nil {
		t.Fatalf("OrderSubscriptionAdjustmentDelete error: %s", err.Error())
	}
	if request.Method != http.MethodDelete {
		t.Errorf("method = %s, want %s", request.Method, http.MethodDelete)
	}
	if request.URL.Path != "/order_subscriptions/subscription_1234/adjustments" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if payload["order_subscription_adjustment_id"] != "adjustment_1234" {
		t.Errorf("order_subscription_adjustment_id = %v", payload["order_subscription_adjustment_id"])
	}
}

func TestOrderSubscriptionBillList(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"list":[]}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionBillList(APIResponse{
		"order_subscription_id": "subscription_1234",
	}, "")
	if err != nil {
		t.Fatalf("OrderSubscriptionBillList error: %s", err.Error())
	}
	// ⚠️ 경로는 언더스코어다 (하이픈이 아니다)
	if request.URL.Path != "/order_subscription_bills" {
		t.Errorf("path = %s, want /order_subscription_bills", request.URL.Path)
	}
	query := request.URL.Query()
	if id := query.Get("order_subscription_id"); id != "subscription_1234" {
		t.Errorf("order_subscription_id = %s", id)
	}
	if limit := query.Get("limit"); limit != "20" {
		t.Errorf("limit = %s, want 20", limit)
	}
}

func TestOrderSubscriptionRequestsIngPause(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionRequestsIngPause("subscription_1234", APIResponse{
		"reason": "이사",
	}, "")
	if err != nil {
		t.Fatalf("OrderSubscriptionRequestsIngPause error: %s", err.Error())
	}
	if request.Method != http.MethodPost {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPost)
	}
	if request.URL.Path != "/order_subscriptions/requests/ing/pause" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "user" {
		t.Errorf("Bootpay-Role = %s, want user", role)
	}
	if payload["order_subscription_id"] != "subscription_1234" {
		t.Errorf("order_subscription_id = %v", payload["order_subscription_id"])
	}
	if payload["reason"] != "이사" {
		t.Errorf("reason = %v", payload["reason"])
	}
}

// ⚠️ requests/ing 계열 중 resume만 PUT이다
func TestOrderSubscriptionRequestsIngResume(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionRequestsIngResume("subscription_1234", nil, "")
	if err != nil {
		t.Fatalf("OrderSubscriptionRequestsIngResume error: %s", err.Error())
	}
	if request.Method != http.MethodPut {
		t.Errorf("method = %s, want %s (resume only uses PUT)", request.Method, http.MethodPut)
	}
	if request.URL.Path != "/order_subscriptions/requests/ing/resume" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if payload["order_subscription_id"] != "subscription_1234" {
		t.Errorf("order_subscription_id = %v", payload["order_subscription_id"])
	}
}

func TestOrderSubscriptionRequestsIngPurchase(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionRequestsIngPurchase("subscription_1234", APIResponse{"price": 5000}, "")
	if err != nil {
		t.Fatalf("OrderSubscriptionRequestsIngPurchase error: %s", err.Error())
	}
	if request.Method != http.MethodPost {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPost)
	}
	if request.URL.Path != "/order_subscriptions/requests/ing/purchase" {
		t.Errorf("path = %s", request.URL.Path)
	}
}

func TestOrderSubscriptionRequestsIngTermination(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionRequestsIngTermination("subscription_1234", APIResponse{
		"termination_fee": 3000,
		"order_number":    "order_1234",
	}, "")
	if err != nil {
		t.Fatalf("OrderSubscriptionRequestsIngTermination error: %s", err.Error())
	}
	if request.Method != http.MethodPost {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPost)
	}
	if request.URL.Path != "/order_subscriptions/requests/ing/termination" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if payload["termination_fee"] != float64(3000) {
		t.Errorf("termination_fee = %v", payload["termination_fee"])
	}
}

func TestOrderSubscriptionRequestsIngTransfer(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionRequestsIngTransfer("subscription_1234", APIResponse{
		"new_user_id": "user_5678",
	}, "")
	if err != nil {
		t.Fatalf("OrderSubscriptionRequestsIngTransfer error: %s", err.Error())
	}
	if request.Method != http.MethodPost {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPost)
	}
	if request.URL.Path != "/order_subscriptions/requests/ing/transfer" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if payload["new_user_id"] != "user_5678" {
		t.Errorf("new_user_id = %v", payload["new_user_id"])
	}
}

func TestOrderSubscriptionCalculateTerminationFee(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"termination_fee":3000}`))
	})
	defer server.Close()

	result, err := api.OrderSubscriptionCalculateTerminationFee("subscription_1234", "", "")
	if err != nil {
		t.Fatalf("OrderSubscriptionCalculateTerminationFee error: %s", err.Error())
	}
	if request.Method != http.MethodGet {
		t.Errorf("method = %s, want %s", request.Method, http.MethodGet)
	}
	if request.URL.Path != "/order_subscriptions/requests/ing/calculate_termination_fee" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if id := request.URL.Query().Get("order_subscription_id"); id != "subscription_1234" {
		t.Errorf("order_subscription_id = %s", id)
	}
	if _, ok := request.URL.Query()["order_number"]; ok {
		t.Error("order_number is empty, so it must not be sent")
	}
	if result["termination_fee"] != float64(3000) {
		t.Errorf("termination_fee = %v", result["termination_fee"])
	}
}

func TestOrderSubscriptionRequestList(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"list":[]}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionRequestList(APIResponse{}, "")
	if err != nil {
		t.Fatalf("OrderSubscriptionRequestList error: %s", err.Error())
	}
	// ⚠️ 하이픈 경로다
	if request.URL.Path != "/order-subscription-requests" {
		t.Errorf("path = %s, want /order-subscription-requests", request.URL.Path)
	}
	// project_id가 없으면 본인 요청만 조회한다
	if role := request.Header.Get("Bootpay-Role"); role != "user" {
		t.Errorf("Bootpay-Role = %s, want user", role)
	}
	if page := request.URL.Query().Get("page"); page != "1" {
		t.Errorf("page = %s, want 1", page)
	}
}

func TestOrderSubscriptionRequestListWithProject(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"list":[]}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionRequestList(APIResponse{"project_id": "project_1234"}, "")
	if err != nil {
		t.Fatalf("OrderSubscriptionRequestList error: %s", err.Error())
	}
	// project_id를 주면 supervisor 모드(프로젝트 전체 검색)다
	if role := request.Header.Get("Bootpay-Role"); role != "supervisor" {
		t.Errorf("Bootpay-Role = %s, want supervisor", role)
	}
	if projectId := request.URL.Query().Get("project_id"); projectId != "project_1234" {
		t.Errorf("project_id = %s", projectId)
	}
}

func TestOrderSubscriptionRequestDetail(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"id":"request_1234"}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionRequestDetail("request_1234", "", "")
	if err != nil {
		t.Fatalf("OrderSubscriptionRequestDetail error: %s", err.Error())
	}
	if request.URL.Path != "/order-subscription-requests/request_1234" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "user" {
		t.Errorf("Bootpay-Role = %s, want user", role)
	}
	if _, ok := request.URL.Query()["project_id"]; ok {
		t.Error("project_id is empty, so it must not be sent")
	}
}

func TestOrderSubscriptionRequestUpdate(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.OrderSubscriptionRequestUpdate("request_1234", "approve", APIResponse{
		"reason": "승인합니다",
	}, "")
	if err != nil {
		t.Fatalf("OrderSubscriptionRequestUpdate error: %s", err.Error())
	}
	if request.Method != http.MethodPut {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPut)
	}
	if request.URL.Path != "/order-subscription-requests/request_1234" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "supervisor" {
		t.Errorf("Bootpay-Role = %s, want supervisor", role)
	}
	// 승인·반려는 approval 값으로 갈린다 (action이 아니다)
	if payload["approval"] != "approve" {
		t.Errorf("approval = %v, want approve", payload["approval"])
	}
	if _, ok := payload["action"]; ok {
		t.Error("action must not be sent")
	}
}

func TestProductCreate(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"id":"product_1234"}`))
	})
	defer server.Close()

	_, err := api.ProductCreate(APIResponse{
		"name":          "부트페이 상품",
		"display_price": 10000,
		"desc":          nil,
	}, nil, "")
	if err != nil {
		t.Fatalf("ProductCreate error: %s", err.Error())
	}
	if request.Method != http.MethodPost {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPost)
	}
	if request.URL.Path != "/products" {
		t.Errorf("path = %s, want /products", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "manager" {
		t.Errorf("Bootpay-Role = %s, want manager", role)
	}
	// images가 없으면 JSON으로 전송한다
	if contentType := request.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type = %s, want application/json", contentType)
	}
	if payload["name"] != "부트페이 상품" {
		t.Errorf("name = %v", payload["name"])
	}
	if _, ok := payload["desc"]; ok {
		t.Error("desc is nil, so it must not be sent")
	}
}

func TestProductCreateWithImages(t *testing.T) {
	var request *http.Request
	var form *multipart.Form
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		if err := req.ParseMultipartForm(1 << 20); err != nil {
			t.Errorf("multipart parse error: %s", err.Error())
		}
		form = req.MultipartForm
		res.Write([]byte(`{"id":"product_1234"}`))
	})
	defer server.Close()

	imagePath := filepath.Join(t.TempDir(), "product.png")
	if err := os.WriteFile(imagePath, []byte("fake-png"), 0644); err != nil {
		t.Fatalf("temp image write error: %s", err.Error())
	}

	_, err := api.ProductCreate(APIResponse{
		"name":             "부트페이 상품",
		"display_price":    10000,
		"use_subscription": true,
		"metadata":         APIResponse{"color": "red"},
	}, []string{imagePath}, "")
	if err != nil {
		t.Fatalf("ProductCreate error: %s", err.Error())
	}
	// ⚠️ boundary가 포함된 multipart Content-Type으로 전송되어야 한다
	if contentType := request.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "multipart/form-data; boundary=") {
		t.Errorf("Content-Type = %s, want multipart/form-data with boundary", contentType)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "manager" {
		t.Errorf("Bootpay-Role = %s, want manager", role)
	}
	if values := form.Value["name"]; len(values) != 1 || values[0] != "부트페이 상품" {
		t.Errorf("name = %v", form.Value["name"])
	}
	if values := form.Value["use_subscription"]; len(values) != 1 || values[0] != "true" {
		t.Errorf("use_subscription = %v, want true", form.Value["use_subscription"])
	}
	// 해시·배열은 JSON 문자열로 전송한다
	if values := form.Value["metadata"]; len(values) != 1 || values[0] != `{"color":"red"}` {
		t.Errorf("metadata = %v", form.Value["metadata"])
	}
	files := form.File["images[0]"]
	if len(files) != 1 || files[0].Filename != "product.png" {
		t.Errorf("images[0] = %v", files)
	}
}

func TestProductUpdate(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.ProductUpdate("product_1234", APIResponse{"stock": 10}, "")
	if err != nil {
		t.Fatalf("ProductUpdate error: %s", err.Error())
	}
	if request.Method != http.MethodPut {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPut)
	}
	if request.URL.Path != "/products/product_1234" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "manager" {
		t.Errorf("Bootpay-Role = %s, want manager", role)
	}
	if payload["stock"] != float64(10) {
		t.Errorf("stock = %v, want 10", payload["stock"])
	}
}

func TestProductDelete(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.ProductDelete("product_1234", "")
	if err != nil {
		t.Fatalf("ProductDelete error: %s", err.Error())
	}
	if request.Method != http.MethodDelete {
		t.Errorf("method = %s, want %s", request.Method, http.MethodDelete)
	}
	if request.URL.Path != "/products/product_1234" {
		t.Errorf("path = %s", request.URL.Path)
	}
}

func TestProductStatus(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.ProductStatus("product_1234", APIResponse{
		"status_sale":    1,
		"status_display": 0,
	}, "")
	if err != nil {
		t.Fatalf("ProductStatus error: %s", err.Error())
	}
	if request.Method != http.MethodPut {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPut)
	}
	if request.URL.Path != "/products/product_1234/status" {
		t.Errorf("path = %s", request.URL.Path)
	}
	if payload["status_sale"] != float64(1) {
		t.Errorf("status_sale = %v, want 1", payload["status_sale"])
	}
	if payload["status_display"] != float64(0) {
		t.Errorf("status_display = %v, want 0", payload["status_display"])
	}
}

func TestUidExist(t *testing.T) {
	var request *http.Request
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		res.Write([]byte(`{"exist":false}`))
	})
	defer server.Close()

	_, err := api.UidExist("uid_1234", "")
	if err != nil {
		t.Fatalf("UidExist error: %s", err.Error())
	}
	if request.Method != http.MethodGet {
		t.Errorf("method = %s, want %s", request.Method, http.MethodGet)
	}
	if request.URL.Path != "/users/join/uid-exist" {
		t.Errorf("path = %s, want /users/join/uid-exist", request.URL.Path)
	}
	if pk := request.URL.Query().Get("pk"); pk != "uid_1234" {
		t.Errorf("pk = %s, want uid_1234", pk)
	}
}

func TestUserUpdate(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.UserUpdate("user_1234", APIResponse{
		"name": "홍길동",
		"group": APIResponse{
			"company_name":    "부트페이",
			"business_number": "1234567890",
		},
	}, "")
	if err != nil {
		t.Fatalf("UserUpdate error: %s", err.Error())
	}
	if request.Method != http.MethodPut {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPut)
	}
	if request.URL.Path != "/users/user_1234" {
		t.Errorf("path = %s, want /users/user_1234", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "user" {
		t.Errorf("Bootpay-Role = %s, want user", role)
	}
	group, ok := payload["group"].(map[string]interface{})
	if !ok || group["company_name"] != "부트페이" {
		t.Errorf("group = %v", payload["group"])
	}
}

func TestUserGroupLimit(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.UserGroupLimit("group_1234", APIResponse{
		"use_limit":            true,
		"limit_month_purchase": 100000,
	}, "")
	if err != nil {
		t.Fatalf("UserGroupLimit error: %s", err.Error())
	}
	if request.Method != http.MethodPut {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPut)
	}
	if request.URL.Path != "/user-groups/group_1234/limit" {
		t.Errorf("path = %s, want /user-groups/group_1234/limit", request.URL.Path)
	}
	if role := request.Header.Get("Bootpay-Role"); role != "manager" {
		t.Errorf("Bootpay-Role = %s, want manager", role)
	}
	if payload["use_limit"] != true {
		t.Errorf("use_limit = %v, want true", payload["use_limit"])
	}
}

func TestUserGroupAggregateTransaction(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.UserGroupAggregateTransaction("group_1234", APIResponse{
		"use_subscription_aggregate_transaction": true,
		"subscription_month_day":                 25,
	}, "")
	if err != nil {
		t.Fatalf("UserGroupAggregateTransaction error: %s", err.Error())
	}
	if request.Method != http.MethodPut {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPut)
	}
	if request.URL.Path != "/user-groups/group_1234/aggregate-transaction" {
		t.Errorf("path = %s, want /user-groups/group_1234/aggregate-transaction", request.URL.Path)
	}
	if payload["subscription_month_day"] != float64(25) {
		t.Errorf("subscription_month_day = %v, want 25", payload["subscription_month_day"])
	}
}

func TestSendTestWebhook(t *testing.T) {
	var request *http.Request
	var payload map[string]interface{}
	api, server := newTestApi(func(res http.ResponseWriter, req *http.Request) {
		request = req
		payload = decodeRequestBody(t, req)
		res.Write([]byte(`{"status":1}`))
	})
	defer server.Close()

	_, err := api.SendTestWebhook("application/json")
	if err != nil {
		t.Fatalf("SendTestWebhook error: %s", err.Error())
	}
	if request.Method != http.MethodPost {
		t.Errorf("method = %s, want %s", request.Method, http.MethodPost)
	}
	if request.URL.Path != "/webhook/test" {
		t.Errorf("path = %s, want /webhook/test", request.URL.Path)
	}
	if payload["header_content_type"] != "application/json" {
		t.Errorf("header_content_type = %v", payload["header_content_type"])
	}
}
