package bootpay

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// capturedCommerceRequest 는 mock transport 가 잡아낸 wire-format 스냅샷이다.
type capturedCommerceRequest struct {
	Method string
	URL    string
	Header http.Header
	Body   []byte
}

// newMockCommerceApi 는 실제 네트워크 호출 없이 요청을 캡처하는 Commerce API 인스턴스를 만든다.
func newMockCommerceApi(captured *[]capturedCommerceRequest) *CommerceApi {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			var body []byte
			if req.Body != nil {
				body, _ = io.ReadAll(req.Body)
			}
			*captured = append(*captured, capturedCommerceRequest{
				Method: req.Method,
				URL:    req.URL.String(),
				Header: req.Header.Clone(),
				Body:   body,
			})
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"success":true}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}
	return NewCommerceAPI("ck", "sk", client, "development")
}

func lastRequest(t *testing.T, captured []capturedCommerceRequest) capturedCommerceRequest {
	t.Helper()
	if len(captured) == 0 {
		t.Fatal("no request captured")
	}
	return captured[len(captured)-1]
}

func decodeBody(t *testing.T, req capturedCommerceRequest) map[string]interface{} {
	t.Helper()
	body := map[string]interface{}{}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		t.Fatalf("body is not JSON: %v (%s)", err, string(req.Body))
	}
	return body
}

func TestCommerceMallSettingRequests(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	// 인스턴스 role 이 user 여도 supervisor 전용 endpoint 는 supervisor 로 나가야 한다 (role 덮어쓰기 금지 검증)
	commerce.SetRole("user")

	if _, err := commerce.MallSetting.GetMallSetting(""); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodGet || req.URL != COMMERCE_DEVELOPMENT+"/mall-setting" {
		t.Fatalf("get mall-setting mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" {
		t.Fatalf("mall-setting must use supervisor role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
	if req.Header.Get("Idempotency-Key") == "" {
		t.Fatal("Idempotency-Key must be auto-generated")
	}

	useLogo := false
	restStartHour := 0
	if _, err := commerce.MallSetting.UpdateMallSetting(MallSettingUpdateParams{
		Name:          "테스트몰",
		UseLogo:       &useLogo,
		RestStartHour: &restStartHour,
	}, "fixed-key"); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodPut || req.URL != COMMERCE_DEVELOPMENT+"/mall-setting" {
		t.Fatalf("update mall-setting mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("Idempotency-Key") != "fixed-key" {
		t.Fatalf("explicit Idempotency-Key must win: %q", req.Header.Get("Idempotency-Key"))
	}
	body := decodeBody(t, req)
	if body["name"] != "테스트몰" {
		t.Fatalf("name missing in body: %+v", body)
	}
	if body["use_logo"] != false {
		t.Fatalf("explicit false must be sent: %+v", body)
	}
	if body["rest_start_hour"] != float64(0) {
		t.Fatalf("explicit 0 must be sent: %+v", body)
	}
	if _, exists := body["use_favicon"]; exists {
		t.Fatalf("unset values must not be sent: %+v", body)
	}
	if _, exists := body["seller_name"]; exists {
		t.Fatalf("unset values must not be sent: %+v", body)
	}
}

func TestCommerceWebhookSendTest(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	headerContentType := 0
	if _, err := commerce.Webhook.SendTest(&SendTestWebhookParams{HeaderContentType: &headerContentType}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodPost || req.URL != COMMERCE_DEVELOPMENT+"/webhook/test" {
		t.Fatalf("webhook/test mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("Idempotency-Key") == "" {
		t.Fatal("Idempotency-Key must be auto-generated")
	}
	body := decodeBody(t, req)
	if body["header_content_type"] != float64(0) {
		t.Fatalf("explicit header_content_type=0 must be sent: %+v", body)
	}

	// params nil 이면 빈 바디로 전송된다 (header_content_type 미전송)
	if _, err := commerce.Webhook.SendTest(nil); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	body = decodeBody(t, req)
	if _, exists := body["header_content_type"]; exists {
		t.Fatalf("unset header_content_type must not be sent: %+v", body)
	}
}

func TestCommerceSupervisorCharge(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	if _, err := commerce.OrderSubscription.SupervisorCharge(SupervisorOrderSubscriptionChargeParams{}); err == nil {
		t.Fatal("empty charge_key must be rejected")
	}
	if len(captured) != 0 {
		t.Fatal("no request must be sent when charge_key is missing")
	}

	if _, err := commerce.OrderSubscription.SupervisorCharge(SupervisorOrderSubscriptionChargeParams{
		ChargeKey:      "charge_key_1",
		Price:          1000,
		IdempotencyKey: "charge-key-fixed",
	}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodPost || req.URL != COMMERCE_DEVELOPMENT+"/order_subscriptions/charge" {
		t.Fatalf("charge mismatch: %s %s", req.Method, req.URL)
	}
	if strings.Contains(req.URL, "charge_key_1") {
		t.Fatal("charge_key must never appear in the URL")
	}
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" {
		t.Fatalf("charge must use supervisor role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
	if req.Header.Get("Idempotency-Key") != "charge-key-fixed" {
		t.Fatalf("explicit Idempotency-Key must win: %q", req.Header.Get("Idempotency-Key"))
	}
	body := decodeBody(t, req)
	if body["charge_key"] != "charge_key_1" || body["price"] != float64(1000) {
		t.Fatalf("charge body mismatch: %+v", body)
	}

	if _, err := commerce.OrderSubscription.SupervisorChargeRevoke(SupervisorOrderSubscriptionChargeRevokeParams{
		ChargeKey: "charge_key_1",
	}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodDelete || req.URL != COMMERCE_DEVELOPMENT+"/order_subscriptions/charge" {
		t.Fatalf("charge revoke mismatch: %s %s", req.Method, req.URL)
	}
	body = decodeBody(t, req)
	if body["charge_key"] != "charge_key_1" {
		t.Fatalf("revoke must send charge_key in the DELETE body: %+v", body)
	}
	if req.Header.Get("Idempotency-Key") == "" {
		t.Fatal("Idempotency-Key must be auto-generated")
	}
}

func TestCommerceMallUserEndpoints(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	if _, err := commerce.User.UserLogin(MallUserLoginParams{LoginId: "tester", Password: "pw"}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodPost || req.URL != COMMERCE_DEVELOPMENT+"/users/login" {
		t.Fatalf("users/login mismatch: %s %s", req.Method, req.URL)
	}
	body := decodeBody(t, req)
	if body["login_id"] != "tester" || body["password"] != "pw" || body["corporate_type"] != float64(0) {
		t.Fatalf("login body mismatch (corporate_type must default to 0): %+v", body)
	}
	if req.Header.Get("Bootpay-User-JWT") != "" {
		t.Fatal("login must not send Bootpay-User-JWT")
	}

	if _, err := commerce.User.UserSession("jwt-token", ""); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodGet || req.URL != COMMERCE_DEVELOPMENT+"/users/session" {
		t.Fatalf("users/session mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("Bootpay-User-JWT") != "jwt-token" {
		t.Fatalf("session must send Bootpay-User-JWT: %q", req.Header.Get("Bootpay-User-JWT"))
	}

	if _, err := commerce.User.UserSession("", ""); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if _, exists := req.Header["Bootpay-User-Jwt"]; exists {
		t.Fatal("empty JWT must not attach the Bootpay-User-JWT header")
	}

	if _, err := commerce.User.UserLogout("jwt-token", ""); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodDelete || req.URL != COMMERCE_DEVELOPMENT+"/users/session" {
		t.Fatalf("logout mismatch: %s %s", req.Method, req.URL)
	}

	if _, err := commerce.User.UserJoinCheck(MALL_USER_JOIN_CHECK_ID_EXIST, "some id&", ""); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/users/join/id-exist?pk=some+id%26" {
		t.Fatalf("join check URL mismatch: %s", req.URL)
	}

	if _, err := commerce.User.UidExist("ex_uid_1", ""); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/users/join/uid-exist?pk=ex_uid_1" {
		t.Fatalf("uid-exist URL mismatch: %s", req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("uid-exist must use user role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
}

func TestCommerceMallProducts(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	// page/limit 미지정시 1/20 기본값이 항상 전송된다
	if _, err := commerce.Product.Products(nil); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/products?limit=20&page=1" {
		t.Fatalf("mall products default URL mismatch: %s", req.URL)
	}

	params := &MallProductListParams{CategoryId: "cat_1", Sort: "recent", UserJwt: "jwt-token"}
	params.Page = 2
	params.Limit = 5
	if _, err := commerce.Product.Products(params); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/products?category_id=cat_1&limit=5&page=2&sort=recent" {
		t.Fatalf("mall products URL mismatch: %s", req.URL)
	}
	if req.Header.Get("Bootpay-User-JWT") != "jwt-token" {
		t.Fatal("mall products must send Bootpay-User-JWT when present")
	}

	if _, err := commerce.Product.ProductDetail("prod_1", "jwt-token", ""); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodGet || req.URL != COMMERCE_DEVELOPMENT+"/products/prod_1" {
		t.Fatalf("mall product detail mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("Bootpay-User-JWT") != "jwt-token" {
		t.Fatal("mall product detail must send Bootpay-User-JWT when present")
	}
}

func TestCommerceProductCreateMultipart(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	// 이미지가 없으면 JSON 으로 전송된다
	if _, err := commerce.Product.Create(CommerceProduct{Name: "상품A", DisplayPrice: 1000}, nil); err != nil {
		t.Fatal(err)
	}
	if len(captured) != 1 {
		t.Fatalf("JSON create must send exactly one request, got %d", len(captured))
	}
	req := lastRequest(t, captured)
	if !strings.HasPrefix(req.Header.Get("Content-Type"), "application/json") {
		t.Fatalf("no-image create must be JSON: %q", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("BOOTPAY-ROLE") != "manager" {
		t.Fatalf("product create must use manager role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
	body := decodeBody(t, req)
	if body["name"] != "상품A" || body["display_price"] != float64(1000) {
		t.Fatalf("JSON create body mismatch: %+v", body)
	}

	// 이미지가 있으면 multipart/form-data + images[0], images[1] 인덱싱
	dir := t.TempDir()
	image1 := filepath.Join(dir, "one.png")
	image2 := filepath.Join(dir, "two.png")
	if err := os.WriteFile(image1, []byte("png-1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(image2, []byte("png-2"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := commerce.Product.Create(CommerceProduct{
		Name:    "상품B",
		Options: []CommerceProductOption{{Name: "옵션1", Price: 500}},
	}, []string{image1, image2}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodPost || req.URL != COMMERCE_DEVELOPMENT+"/products" {
		t.Fatalf("multipart create mismatch: %s %s", req.Method, req.URL)
	}

	mediaType, mediaParams, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		t.Fatal(err)
	}
	if mediaType != "multipart/form-data" || mediaParams["boundary"] == "" {
		t.Fatalf("Content-Type must be multipart/form-data with boundary: %q", req.Header.Get("Content-Type"))
	}

	reader := multipart.NewReader(bytes.NewReader(req.Body), mediaParams["boundary"])
	fields := map[string]string{}
	files := map[string]string{}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		content, _ := io.ReadAll(part)
		if part.FileName() != "" {
			files[part.FormName()] = part.FileName()
		} else {
			fields[part.FormName()] = string(content)
		}
	}
	if fields["name"] != "상품B" {
		t.Fatalf("multipart name field mismatch: %+v", fields)
	}
	var options []map[string]interface{}
	if err := json.Unmarshal([]byte(fields["options"]), &options); err != nil || len(options) != 1 {
		t.Fatalf("options must be a JSON string field: %q", fields["options"])
	}
	if files["images[0]"] != "one.png" || files["images[1]"] != "two.png" {
		t.Fatalf("images must be indexed as images[0], images[1]: %+v", files)
	}
}

func TestCommerceInvoiceRequests(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	// page/limit 미지정시 1/24 기본값이 항상 전송된다 (구 List 시그니처 포함)
	if _, err := commerce.Invoice.List(nil); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/invoices?limit=24&page=1" {
		t.Fatalf("invoice list default URL mismatch: %s", req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("invoice list must use user role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
	if req.Header.Get("Idempotency-Key") == "" {
		t.Fatal("Idempotency-Key must be auto-generated")
	}

	productType := 0
	params := &InvoiceListParams{CsType: "cs", UserId: "user_1", ProductType: &productType, CssAt: "2026-01-01", CseAt: "2026-01-31"}
	if _, err := commerce.Invoice.ListWithParams(params); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	expected := COMMERCE_DEVELOPMENT + "/invoices?cs_type=cs&cse_at=2026-01-31&css_at=2026-01-01&limit=24&page=1&product_type=0&user_id=user_1"
	if req.URL != expected {
		t.Fatalf("invoice list URL mismatch:\n got %s\nwant %s", req.URL, expected)
	}

	// sendTypes nil → send_types 미전송, 빈 배열 → 빈 배열 전송
	if _, err := commerce.Invoice.Notify("inv_1", nil); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/invoices/inv_1/notify" {
		t.Fatalf("notify URL mismatch: %s", req.URL)
	}
	body := decodeBody(t, req)
	if _, exists := body["send_types"]; exists {
		t.Fatalf("nil sendTypes must omit send_types: %+v", body)
	}

	if _, err := commerce.Invoice.Notify("inv_1", []int{}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	body = decodeBody(t, req)
	if _, exists := body["send_types"]; !exists {
		t.Fatalf("empty sendTypes must send send_types: %+v", body)
	}

	if _, err := commerce.Invoice.Detail("inv_1"); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Header.Get("Idempotency-Key") == "" {
		t.Fatal("invoice detail must attach Idempotency-Key")
	}
}

func TestCommerceOrderCancelIdResolution(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	if _, err := commerce.OrderCancel.Approve(OrderCancelActionParams{}); err == nil {
		t.Fatal("missing id must be rejected")
	}

	// 구 필드만 지정 — 계속 동작해야 한다
	if _, err := commerce.OrderCancel.Approve(OrderCancelActionParams{OrderCancelRequestHistoryId: "old_id"}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/order/cancel/old_id/approve" {
		t.Fatalf("old id field must keep working: %s", req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" {
		t.Fatalf("approve must use supervisor role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}

	// 둘 다 지정시 신 필드 우선
	if _, err := commerce.OrderCancel.Reject(OrderCancelActionParams{
		OrderCancellationRequestId:  "new_id",
		OrderCancelRequestHistoryId: "old_id",
		Message:                     "반려 사유",
	}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/order/cancel/new_id/reject" {
		t.Fatalf("official id must win when both are set: %s", req.URL)
	}
	body := decodeBody(t, req)
	if body["message"] != "반려 사유" {
		t.Fatalf("message must be sent in body: %+v", body)
	}

	if _, err := commerce.OrderCancel.Withdraw("cancel_id"); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodPut || req.URL != COMMERCE_DEVELOPMENT+"/order/cancel/cancel_id/withdraw" {
		t.Fatalf("withdraw mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("withdraw must use user role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
}

func TestCommerceAdjustmentRequests(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	// Create — price/duration/tax_free_price 는 0 / 1 / 0 기본값을 명시 전송한다
	if _, err := commerce.OrderSubscriptionAdjustment.Create("sub_1", CommerceOrderSubscriptionAdjustment{}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodPost || req.URL != COMMERCE_DEVELOPMENT+"/order_subscriptions/sub_1/adjustments" {
		t.Fatalf("adjustment create mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" {
		t.Fatalf("adjustment must use supervisor role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
	body := decodeBody(t, req)
	if body["price"] != float64(0) || body["duration"] != float64(1) || body["tax_free_price"] != float64(0) {
		t.Fatalf("create defaults mismatch: %+v", body)
	}

	// Update — adjustments 배열 지원 + duration 기본 1
	if _, err := commerce.OrderSubscriptionAdjustment.Update(OrderSubscriptionAdjustmentUpdateParams{
		OrderSubscriptionId: "sub_1",
		Adjustments: []CommerceOrderSubscriptionAdjustment{
			{Price: 1000, Duration: 2, Name: "할인"},
		},
	}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodPut {
		t.Fatalf("adjustment update method mismatch: %s", req.Method)
	}
	body = decodeBody(t, req)
	if body["duration"] != float64(1) {
		t.Fatalf("update duration must default to 1: %+v", body)
	}
	adjustments, ok := body["adjustments"].([]interface{})
	if !ok || len(adjustments) != 1 {
		t.Fatalf("adjustments array must be sent: %+v", body)
	}

	// Delete — 대상 ID 는 query 가 아니라 body 로 전송한다
	if _, err := commerce.OrderSubscriptionAdjustment.Delete("sub_1", "adj_1"); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodDelete || req.URL != COMMERCE_DEVELOPMENT+"/order_subscriptions/sub_1/adjustments" {
		t.Fatalf("adjustment delete mismatch (id must not be in query): %s %s", req.Method, req.URL)
	}
	body = decodeBody(t, req)
	if body["order_subscription_adjustment_id"] != "adj_1" {
		t.Fatalf("delete must send target id in body: %+v", body)
	}
}

func TestCommerceRequestIngPurchaseTransfer(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	if _, err := commerce.OrderSubscription.RequestIng.Purchase(OrderSubscriptionPurchaseParams{
		OrderSubscriptionId: "sub_1",
		Price:               10000,
	}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodPost || req.URL != COMMERCE_DEVELOPMENT+"/order_subscriptions/requests/ing/purchase" {
		t.Fatalf("purchase mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("requests/ing must use user role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}

	if _, err := commerce.OrderSubscription.RequestIng.Transfer(OrderSubscriptionTransferParams{
		OrderSubscriptionId: "sub_1",
		NewUserId:           "user_2",
	}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodPost || req.URL != COMMERCE_DEVELOPMENT+"/order_subscriptions/requests/ing/transfer" {
		t.Fatalf("transfer mismatch: %s %s", req.Method, req.URL)
	}
	body := decodeBody(t, req)
	if body["new_user_id"] != "user_2" {
		t.Fatalf("transfer body mismatch: %+v", body)
	}
}

func TestCommerceOrderSubscriptionRequestScope(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	// project_id 없음 → user scope, page/limit 기본 1/20
	if _, err := commerce.OrderSubscriptionRequest.List(nil); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/order-subscription-requests?limit=20&page=1" {
		t.Fatalf("request list default URL mismatch: %s", req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("without project_id role must be user, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}

	// project_id 있음 → supervisor scope + 신규 파라미터
	if _, err := commerce.OrderSubscriptionRequest.List(&OrderSubscriptionRequestListParams{
		ProjectId:           "proj_1",
		OrderSubscriptionId: "sub_1",
		UserId:              "user_1",
		UserGroupId:         "group_1",
	}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	expected := COMMERCE_DEVELOPMENT + "/order-subscription-requests?limit=20&order_subscription_id=sub_1&page=1&project_id=proj_1&user_group_id=group_1&user_id=user_1"
	if req.URL != expected {
		t.Fatalf("request list URL mismatch:\n got %s\nwant %s", req.URL, expected)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" {
		t.Fatalf("with project_id role must be supervisor, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}

	if _, err := commerce.OrderSubscriptionRequest.Detail("req_1", ""); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("detail without project_id must be user, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}

	if _, err := commerce.OrderSubscriptionRequest.Update(OrderSubscriptionRequestUpdateParams{
		OrderSubscriptionRequestHistoryId: "req_1",
		Approval:                          "approve",
	}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" {
		t.Fatalf("request update must be supervisor, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
}

func TestCommerceUserGroupLimitParams(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	limitMonth := 0
	limitWeek := 50000
	if _, err := commerce.UserGroup.Limit(UserGroupLimitParams{
		UserGroupId:        "group_1",
		LimitMonthPurchase: &limitMonth,
		LimitWeekPurchase:  &limitWeek,
	}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/user-groups/group_1/limit" {
		t.Fatalf("limit URL mismatch: %s", req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "manager" {
		t.Fatalf("limit must use manager role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
	body := decodeBody(t, req)
	if body["limit_month_purchase"] != float64(0) || body["limit_week_purchase"] != float64(50000) {
		t.Fatalf("limit body mismatch (explicit 0 must be sent): %+v", body)
	}
}

func TestCommerceOrderListSearchDates(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	if _, err := commerce.Order.List(&OrderListParams{
		SearchDateFrom: "2026-01-01",
		SearchDateTo:   "2026-01-31",
		CssAt:          "2026-02-01",
	}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	expected := COMMERCE_DEVELOPMENT + "/orders?css_at=2026-02-01&search_date_from=2026-01-01&search_date_to=2026-01-31"
	if req.URL != expected {
		t.Fatalf("order list URL mismatch:\n got %s\nwant %s", req.URL, expected)
	}
}

func TestCommerceOrderSubscriptionListNewParams(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	status := 0
	if _, err := commerce.OrderSubscription.List(&OrderSubscriptionListParams{
		SearchDateFrom: "2026-01-01",
		SearchDateTo:   "2026-01-31",
		Status:         &status,
	}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	expected := COMMERCE_DEVELOPMENT + "/order_subscriptions?search_date_from=2026-01-01&search_date_to=2026-01-31&status=0"
	if req.URL != expected {
		t.Fatalf("subscription list URL mismatch:\n got %s\nwant %s", req.URL, expected)
	}

	if _, err := commerce.OrderSubscription.Update(OrderSubscriptionUpdateParams{
		OrderSubscriptionId: "sub_1",
		OrderName:           "구독 변경",
	}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Method != http.MethodPut || req.URL != COMMERCE_DEVELOPMENT+"/order_subscriptions/sub_1" {
		t.Fatalf("subscription update mismatch: %s %s", req.Method, req.URL)
	}
	if req.Header.Get("BOOTPAY-ROLE") != "supervisor" {
		t.Fatalf("subscription update must use supervisor role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
}

func TestCommerceStoreIdempotency(t *testing.T) {
	var captured []capturedCommerceRequest
	commerce := newMockCommerceApi(&captured)

	if _, err := commerce.Store.Info(); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Header.Get("Idempotency-Key") == "" {
		t.Fatal("store info must auto-attach Idempotency-Key")
	}

	if _, err := commerce.Store.DetailWithIdempotencyKey("store-key"); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.Header.Get("Idempotency-Key") != "store-key" {
		t.Fatalf("explicit Idempotency-Key must win: %q", req.Header.Get("Idempotency-Key"))
	}
	// store 조회는 role 을 덮어쓰지 않는다 (인스턴스 기본 role 유지)
	if req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("store lookup must keep the instance role, got %q", req.Header.Get("BOOTPAY-ROLE"))
	}
}
