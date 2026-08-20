package bootpay

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 변경된(기존) 기능의 wire-format 회귀 테스트 —
// 신규 기능 테스트(commerce_parity_mock_test.go)와 분리해 기존 동작 정렬분을 검증한다.

func TestCommerceMallUserJoin(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	gender := 0
	if _, err := commerce.User.UserJoin(MallUserJoinParams{
		LoginId:  "tester",
		Password: "pw",
		Name:     "홍길동",
		Gender:   &gender,
	}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodPost || req.URL != COMMERCE_DEVELOPMENT+"/users/join" {
		t.Fatalf("users/join mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("Idempotency-Key") == "" {
		t.Fatal("Idempotency-Key must be auto-generated")
	}
	body := decodeBody(t, req)
	if body["login_id"] != "tester" || body["name"] != "홍길동" {
		t.Fatalf("join body mismatch: %+v", body)
	}
	if body["corporate_type"] != float64(0) {
		t.Fatalf("corporate_type must default to 0: %+v", body)
	}
	if body["gender"] != float64(0) {
		t.Fatalf("explicit gender=0 must be sent: %+v", body)
	}
	if _, exists := body["email"]; exists {
		t.Fatalf("unset values must not be sent: %+v", body)
	}
	if _, exists := body["group"]; exists {
		t.Fatalf("unset values must not be sent: %+v", body)
	}
}

// 기존 RequestIng 계열(Pause/Resume/Termination/CalculateTerminationFee)이
// user role + Idempotency-Key 헤더로 정렬됐는지 검증한다 (변경 기능 회귀).
func TestCommerceRequestIngChangedBehavior(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	// supervisor 인스턴스여도 requests/ing 는 user 로 나가야 한다
	commerce.AsSupervisor()

	if _, err := commerce.OrderSubscription.RequestIng.Pause(OrderSubscriptionPauseParams{
		OrderSubscriptionId: "sub_1",
		IdempotencyKey:      "pause-key",
	}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodPost || req.URL != COMMERCE_DEVELOPMENT+"/order_subscriptions/requests/ing/pause" {
		t.Fatalf("pause mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("pause must use user role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
	if req.Header.Get("Idempotency-Key") != "pause-key" {
		t.Fatalf("explicit Idempotency-Key must win: %q", req.Header.Get("Idempotency-Key"))
	}

	// Resume 은 requests/ing 계열 중 유일한 PUT — POST 로 바뀌면 안 된다. Reason 필드 신규 지원.
	if _, err := commerce.OrderSubscription.RequestIng.Resume(OrderSubscriptionResumeParams{
		OrderSubscriptionId: "sub_1",
		Reason:              "재개 사유",
	}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodPut || req.URL != COMMERCE_DEVELOPMENT+"/order_subscriptions/requests/ing/resume" {
		t.Fatalf("resume must stay PUT: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("resume must use user role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
	body := decodeBody(t, req)
	if body["reason"] != "재개 사유" {
		t.Fatalf("resume reason must be sent: %+v", body)
	}

	if _, err := commerce.OrderSubscription.RequestIng.Termination(OrderSubscriptionTerminationParams{
		OrderNumber: "order_1",
	}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodPost || req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("termination mismatch: %s role=%q", req.Method, req.Header.Get("BOOTPAY-ROLE"))
	}

	// CalculateTerminationFee — 둘 다 지정시 둘 다 전송 (nodejs HEAD 정렬), user role + Idempotency-Key
	if _, err := commerce.OrderSubscription.RequestIng.CalculateTerminationFee("sub_1", "order_1"); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	expected := COMMERCE_DEVELOPMENT + "/order_subscriptions/requests/ing/calculate_termination_fee?order_number=order_1&order_subscription_id=sub_1"
	if req.URL != expected {
		t.Fatalf("calculate fee URL mismatch:\n got %s\nwant %s", req.URL, expected)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "user" || req.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("calculate fee headers mismatch: role=%q", req.Header.Get("BOOTPAY-ROLE"))
	}

	// 둘 다 비어있으면 요청 없이 에러
	before := len(captured)
	if _, err := commerce.OrderSubscription.RequestIng.CalculateTerminationFee("", ""); err == nil {
		t.Fatal("missing both ids must be rejected")
	}
	if len(captured) != before {
		t.Fatal("no request must be sent when both ids are missing")
	}
}

// 상품 쓰기(Update/Status/Delete)가 manager 헤더로 정렬됐는지 검증한다 (변경 기능 회귀).
func TestCommerceProductWriteHeaders(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	if _, err := commerce.Product.Update(CommerceProduct{}); err == nil {
		t.Fatal("missing product_id must be rejected")
	}

	if _, err := commerce.Product.Update(CommerceProduct{ProductId: "prod_1", Name: "수정"}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodPut || req.URL != COMMERCE_DEVELOPMENT+"/products/prod_1" {
		t.Fatalf("product update mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "manager" || req.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("product update must use manager role + Idempotency-Key: role=%q", req.Header.Get("BOOTPAY-ROLE"))
	}

	statusFrozen := false
	if _, err := commerce.Product.Status(ProductStatusParams{
		ProductId:      "prod_1",
		Status:         2,
		StatusFrozen:   &statusFrozen,
		IdempotencyKey: "status-key",
	}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/products/prod_1/status" {
		t.Fatalf("product status URL mismatch: %s", req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "manager" || req.Header.Get("Idempotency-Key") != "status-key" {
		t.Fatalf("product status headers mismatch: role=%q key=%q", req.Header.Get("BOOTPAY-ROLE"), req.Header.Get("Idempotency-Key"))
	}
	body := decodeBody(t, req)
	if body["status"] != float64(2) {
		t.Fatalf("status body mismatch: %+v", body)
	}
	if body["status_frozen"] != false {
		t.Fatalf("explicit status_frozen=false must be sent: %+v", body)
	}

	if _, err := commerce.Product.Delete("prod_1"); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodDelete || req.Header.Get("BOOTPAY-ROLE") != "manager" {
		t.Fatalf("product delete mismatch: %s role=%q", req.Method, req.Header.Get("BOOTPAY-ROLE"))
	}

	if _, err := commerce.Product.CreateSimple(CommerceProduct{Name: "단순생성"}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Header.Get("BOOTPAY-ROLE") != "manager" {
		t.Fatalf("CreateSimple must use manager role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
}

// OrderCancel List 가 user 헤더 + query 필터로 정렬됐는지 검증한다 (변경 기능 회귀).
func TestCommerceOrderCancelListHeaders(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	if _, err := commerce.OrderCancel.List(&OrderCancelListParams{
		OrderNumber:    "order_1",
		IdempotencyKey: "list-key",
	}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/order/cancel?order_number=order_1" {
		t.Fatalf("cancel list URL mismatch: %s", req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "user" || req.Header.Get("Idempotency-Key") != "list-key" {
		t.Fatalf("cancel list headers mismatch: role=%q key=%q", req.Header.Get("BOOTPAY-ROLE"), req.Header.Get("Idempotency-Key"))
	}

	// params nil 도 기존처럼 동작 + 헤더 자동 부착
	if _, err := commerce.OrderCancel.List(nil); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/order/cancel" {
		t.Fatalf("cancel list nil params URL mismatch: %s", req.URL)
	}
	if req.Header.Get("Idempotency-Key") == "" {
		t.Fatal("Idempotency-Key must be auto-generated")
	}
}

// UserGroup AggregateTransaction 이 manager 헤더로 정렬됐는지 검증한다 (변경 기능 회귀).
func TestCommerceUserGroupAggregateTransactionHeaders(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	if _, err := commerce.UserGroup.AggregateTransaction(UserGroupAggregateTransactionParams{}); err == nil {
		t.Fatal("missing user_group_id must be rejected")
	}

	if _, err := commerce.UserGroup.AggregateTransaction(UserGroupAggregateTransactionParams{
		UserGroupId:          "group_1",
		SubscriptionMonthDay: 15,
	}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodPut || req.URL != COMMERCE_DEVELOPMENT+"/user-groups/group_1/aggregate-transaction" {
		t.Fatalf("aggregate-transaction mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "manager" || req.Header.Get("Idempotency-Key") == "" {
		t.Fatalf("aggregate-transaction must use manager role + Idempotency-Key: role=%q", req.Header.Get("BOOTPAY-ROLE"))
	}
}

// role 헤더 우선순위 양방향 검증 —
// (1) 요청별 role 이 지정된 endpoint 는 인스턴스 role 을 무시하고,
// (2) role 미지정 endpoint 는 인스턴스 role 을 그대로 쓴다.
func TestCommerceRoleHeaderNotOverwritten(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	// supervisor 인스턴스 → user 전용 endpoint 는 user 로
	commerce.AsSupervisor()
	if _, err := commerce.Invoice.Detail("inv_1"); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("user-scoped endpoint must send user even on supervisor instance, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}

	// user 인스턴스 → supervisor 전용 endpoint 는 supervisor 로
	commerce.AsUser()
	if _, err := commerce.OrderSubscriptionAdjustment.Delete("sub_1", "adj_1"); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" {
		t.Fatalf("supervisor-scoped endpoint must send supervisor even on user instance, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}

	// role 미지정 endpoint (Order.List, mall 계열) 는 인스턴스 role 유지
	commerce.AsSupervisor()
	if _, err := commerce.Order.List(nil); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" {
		t.Fatalf("role-neutral endpoint must keep instance role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
	if _, err := commerce.User.UserSession("jwt-token", ""); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" {
		t.Fatalf("mall endpoint must keep instance role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
}

// httptest 실서버 왕복 검증 — mock transport 가 아닌 실제 HTTP 스택을 통과한 요청이
// 서버 수신 시점에도 method/path/헤더/바디를 유지하는지 확인한다.
func TestCommerceHttptestServerRoundTrip(t *testing.T) {
	type received struct {
		Method string
		Path   string
		Role   string
		Idem   string
		Body   map[string]interface{}
	}
	var got []received

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := map[string]interface{}{}
		if r.Body != nil {
			raw, _ := io.ReadAll(r.Body)
			if len(raw) > 0 {
				_ = json.Unmarshal(raw, &body)
			}
		}
		got = append(got, received{
			Method: r.Method,
			Path:   r.URL.RequestURI(),
			Role:   r.Header.Get("BOOTPAY-ROLE"),
			Idem:   r.Header.Get("Idempotency-Key"),
			Body:   body,
		})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	commerce := NewCommerceAPI("ck", "sk", nil, "development")
	commerce.baseUrl = server.URL

	if _, err := commerce.MallSetting.GetMallSetting("mall-key"); err != nil {
		t.Fatal(err)
	}
	if _, err := commerce.OrderSubscription.SupervisorCharge(SupervisorOrderSubscriptionChargeParams{
		ChargeKey: "charge_key_1",
		Price:     1000,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := commerce.OrderSubscription.SupervisorChargeRevoke(SupervisorOrderSubscriptionChargeRevokeParams{
		ChargeKey: "charge_key_1",
	}); err != nil {
		t.Fatal(err)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 received requests, got %d", len(got))
	}
	if got[0].Method != http.MethodGet || got[0].Path != "/mall-setting" || got[0].Role != "supervisor" || got[0].Idem != "mall-key" {
		t.Fatalf("mall-setting received mismatch: %+v", got[0])
	}
	if got[1].Method != http.MethodPost || got[1].Path != "/order_subscriptions/charge" || got[1].Role != "supervisor" {
		t.Fatalf("charge received mismatch: %+v", got[1])
	}
	if got[1].Body["charge_key"] != "charge_key_1" || got[1].Body["price"] != float64(1000) {
		t.Fatalf("charge received body mismatch: %+v", got[1].Body)
	}
	// DELETE body 가 실서버 수신 시점에도 살아있는지 (charge_key 는 URL 이 아니라 body)
	if got[2].Method != http.MethodDelete || got[2].Path != "/order_subscriptions/charge" {
		t.Fatalf("charge revoke received mismatch: %+v", got[2])
	}
	if got[2].Body["charge_key"] != "charge_key_1" {
		t.Fatalf("charge revoke DELETE body must survive the real HTTP stack: %+v", got[2].Body)
	}
	if strings.Contains(got[2].Path, "charge_key_1") {
		t.Fatal("charge_key must never appear in the URL")
	}
}

// Bill.List parity — page/limit 기본 1/20 상시 전송 + user role + Idempotency-Key (P1-1)
func TestCommerceBillListParity(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	if _, err := commerce.OrderSubscriptionBill.List(nil); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/order_subscription_bills?limit=20&page=1" {
		t.Fatalf("bill list default URL mismatch: %s", req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("bill list must use user role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
	if req.Header.Get("Idempotency-Key") == "" {
		t.Fatal("Idempotency-Key must be auto-generated")
	}

	params := &OrderSubscriptionBillListParams{
		OrderSubscriptionId: "sub_1",
		Status:              []int{1, 2},
		IdempotencyKey:      "bill-key",
	}
	params.Page = 3
	params.Limit = 5
	if _, err := commerce.OrderSubscriptionBill.List(params); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	expected := COMMERCE_DEVELOPMENT + "/order_subscription_bills?limit=5&order_subscription_id=sub_1&page=3&status=1%2C2"
	if req.URL != expected {
		t.Fatalf("bill list URL mismatch:\n got %s\nwant %s", req.URL, expected)
	}
	if req.Header.Get("Idempotency-Key") != "bill-key" {
		t.Fatalf("explicit Idempotency-Key must win: %q", req.Header.Get("Idempotency-Key"))
	}
}

// Request.Detail 의 project_id 분기 — supervisor 모드 (P1 테스트 공백 보강)
func TestCommerceRequestDetailProjectScope(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	if _, err := commerce.OrderSubscriptionRequest.Detail("req_1", "proj_1"); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/order-subscription-requests/req_1?project_id=proj_1" {
		t.Fatalf("detail URL mismatch: %s", req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" {
		t.Fatalf("detail with project_id must be supervisor, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
}

// MallSetting 의 Detail/Update alias 가 본체와 동일 wire 인지 검증
func TestCommerceMallSettingAliases(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	if _, err := commerce.MallSetting.Detail("alias-key"); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodGet || req.URL != COMMERCE_DEVELOPMENT+"/mall-setting" {
		t.Fatalf("Detail alias mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" || req.Header.Get("Idempotency-Key") != "alias-key" {
		t.Fatalf("Detail alias headers mismatch: role=%q key=%q", req.Header.Get("BOOTPAY-ROLE"), req.Header.Get("Idempotency-Key"))
	}

	if _, err := commerce.MallSetting.Update(MallSettingUpdateParams{Name: "별칭몰"}, "alias-key-2"); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodPut || req.URL != COMMERCE_DEVELOPMENT+"/mall-setting" {
		t.Fatalf("Update alias mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" || req.Header.Get("Idempotency-Key") != "alias-key-2" {
		t.Fatalf("Update alias headers mismatch: role=%q key=%q", req.Header.Get("BOOTPAY-ROLE"), req.Header.Get("Idempotency-Key"))
	}
	body := decodeBody(t, req)
	if body["name"] != "별칭몰" {
		t.Fatalf("Update alias body mismatch: %+v", body)
	}
}

// Request Update body — Extra 병합 유실 수정 검증 (P1-3): 타입 필드 + Extra 가 body 에 병합되고,
// 같은 키는 타입 필드가 우선하며, nil Extra 값은 미전송된다.
func TestCommerceRequestUpdateBodyMerge(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	price := 0
	terminationFee := 5000
	if _, err := commerce.OrderSubscriptionRequest.Update(OrderSubscriptionRequestUpdateParams{
		OrderSubscriptionRequestHistoryId: "req_1",
		Approval:                          "approve",
		Reason:                            "승인 사유",
		Price:                             &price,
		TerminationFee:                    &terminationFee,
		ServiceEndAt:                      "2026-12-31",
		Extra: map[string]interface{}{
			"custom_key": "custom_value",
			"price":      999, // 타입 필드와 충돌 — 타입 필드(0)가 이겨야 한다
			"nil_key":    nil, // nil 은 미전송
		},
		IdempotencyKey: "req-key",
	}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodPut || req.URL != COMMERCE_DEVELOPMENT+"/order-subscription-requests/req_1" {
		t.Fatalf("request update mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("Idempotency-Key") != "req-key" {
		t.Fatalf("explicit Idempotency-Key must win: %q", req.Header.Get("Idempotency-Key"))
	}
	body := decodeBody(t, req)
	if body["approval"] != "approve" || body["reason"] != "승인 사유" {
		t.Fatalf("approval/reason mismatch: %+v", body)
	}
	if body["price"] != float64(0) {
		t.Fatalf("typed price=0 must win over Extra: %+v", body)
	}
	if body["termination_fee"] != float64(5000) || body["service_end_at"] != "2026-12-31" {
		t.Fatalf("typed fields must be sent: %+v", body)
	}
	if body["custom_key"] != "custom_value" {
		t.Fatalf("Extra must be merged into the body (was silently dropped): %+v", body)
	}
	if _, exists := body["nil_key"]; exists {
		t.Fatalf("nil Extra values must not be sent: %+v", body)
	}
}

// ProductStatusParams 의 sale period 3필드 (P1-2)
func TestCommerceProductStatusSalePeriodFields(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	useSalePeriod := false
	if _, err := commerce.Product.Status(ProductStatusParams{
		ProductId:     "prod_1",
		UseSalePeriod: &useSalePeriod,
		SaleStartAt:   "2026-09-01",
		SaleEndAt:     "2026-09-30",
	}); err != nil {
		t.Fatal(err)
	}
	body := decodeBody(t, lastRequest(t, captured))
	if body["use_sale_period"] != false {
		t.Fatalf("explicit use_sale_period=false must be sent: %+v", body)
	}
	if body["sale_start_at"] != "2026-09-01" || body["sale_end_at"] != "2026-09-30" {
		t.Fatalf("sale period fields must be sent: %+v", body)
	}
}
