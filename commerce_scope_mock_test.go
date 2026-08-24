package bootpay

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// commerce-api 가 supervisor / manager scope 를 요구하는 엔드포인트들.
// 헤더를 붙이지 않으면 인스턴스 기본값 user 로 조용히 나가고 서버가 scope_invalid! 로 거절한다.
func TestCommerceScopeRequiredEndpointsSendExpectedRole(t *testing.T) {
	cases := []struct {
		name   string
		call   func(api *CommerceApi) (map[string]interface{}, error)
		method string
		url    string
		role   string
	}{
		{"SupervisorApprove", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.OrderSubscription.SupervisorApprove("s1", &SupervisorOrderSubscriptionApproveParams{Reason: "승인"})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/order_subscriptions/s1/approve", "supervisor"},
		{"SupervisorReject", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.OrderSubscription.SupervisorReject("s1", &SupervisorOrderSubscriptionRejectParams{Reason: "반려"})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/order_subscriptions/s1/reject", "supervisor"},
		{"SupervisorTerminate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.OrderSubscription.SupervisorTerminate("s1", &SupervisorOrderSubscriptionTerminateParams{Reason: "해지"})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/order_subscriptions/s1/terminate", "supervisor"},
		{"SupervisorPause", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.OrderSubscription.SupervisorPause("s1", SupervisorOrderSubscriptionPauseParams{PausedAt: "2026-01-01"})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/order_subscriptions/s1/pause", "supervisor"},
		{"SupervisorResume", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.OrderSubscription.SupervisorResume("s1", nil)
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/order_subscriptions/s1/resume", "supervisor"},
		{"CategoryCreate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Category.Create(CategoryCreateParams{Name: "카테고리"})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/categories", "supervisor"},
		{"CategoryUpdate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Category.Update(CategoryUpdateParams{CategoryId: "c1", Name: "변경"})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/categories/c1", "supervisor"},
		{"CategoryDestroy", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Category.Destroy("c1")
		}, http.MethodDelete, COMMERCE_DEVELOPMENT + "/categories/c1", "supervisor"},
		{"UserGroupUserCreate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.UserGroup.UserCreate("g1", "u1")
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/user-groups/g1/user", "manager"},
		{"UserGroupUserDelete", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.UserGroup.UserDelete("g1", "u1")
		}, http.MethodDelete, COMMERCE_DEVELOPMENT + "/user-groups/g1/user/u1", "manager"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var captured []capturedCommerceRequest
			api := newMockCommerceApi(&captured)
			if _, err := tc.call(api); err != nil {
				t.Fatalf("%s: %v", tc.name, err)
			}
			req := lastRequest(t, captured)
			if req.Method != tc.method || req.URL != tc.url {
				t.Fatalf("%s: got %s %s, want %s %s", tc.name, req.Method, req.URL, tc.method, tc.url)
			}
			if got := req.Header.Get("BOOTPAY-ROLE"); got != tc.role {
				t.Fatalf("%s: BOOTPAY-ROLE = %q, want %q", tc.name, got, tc.role)
			}
			if req.Header.Get("Idempotency-Key") == "" {
				t.Fatalf("%s: Idempotency-Key is not attached", tc.name)
			}
		})
	}
}

// 명시한 Idempotency-Key 는 그대로 전달되고 바디에는 실리지 않는다.
func TestCommerceScopeExplicitIdempotencyKey(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.Category.Create(CategoryCreateParams{Name: "카테고리", IdempotencyKey: "fixed-key"}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if got := req.Header.Get("Idempotency-Key"); got != "fixed-key" {
		t.Fatalf("Idempotency-Key = %q, want fixed-key", got)
	}
	if strings.Contains(string(req.Body), "idempotency") {
		t.Fatalf("idempotency key leaked into body: %s", string(req.Body))
	}

	if _, err := api.UserGroup.UserCreate("g1", "u1", "member-key"); err != nil {
		t.Fatal(err)
	}
	if got := lastRequest(t, captured).Header.Get("Idempotency-Key"); got != "member-key" {
		t.Fatalf("Idempotency-Key = %q, want member-key", got)
	}
}

// 최상위가 배열인 응답(GET /v1/categories 등)도 받아야 한다.
// 이전에는 map 파싱 실패 에러를 버려서 호출자가 빈 map 을 성공으로 받았다.
func TestCommerceTopLevelArrayResponse(t *testing.T) {
	respond := func(body string) *CommerceApi {
		client := &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}),
		}
		return NewCommerceAPI("ck", "sk", client, "development")
	}

	api := respond(`[{"category_id":"c1"},{"category_id":"c2"}]`)
	result, err := api.Category.List()
	if err != nil {
		t.Fatalf("array response must not error: %v", err)
	}
	items, ok := result["data"].([]interface{})
	if !ok {
		t.Fatalf("array response must be exposed under data, got %#v", result)
	}
	if len(items) != 2 {
		t.Fatalf("data length = %d, want 2", len(items))
	}

	// 객체 응답은 기존과 동일한 모양을 유지한다 (data 로 감싸지 않는다).
	objectApi := respond(`{"success":true,"total":2}`)
	objectResult, err := objectApi.Category.List()
	if err != nil {
		t.Fatal(err)
	}
	if objectResult["success"] != true {
		t.Fatalf("object response changed shape: %#v", objectResult)
	}

	// 빈 본문은 빈 map.
	emptyApi := respond(``)
	emptyResult, err := emptyApi.Category.List()
	if err != nil {
		t.Fatalf("empty body must not error: %v", err)
	}
	if len(emptyResult) != 0 {
		t.Fatalf("empty body must decode to an empty map, got %#v", emptyResult)
	}

	// JSON 이 아닌 본문(HTML 5xx 등)은 조용한 빈 map 대신 에러를 낸다.
	htmlApi := respond(`<html><body>500</body></html>`)
	if _, err := htmlApi.Category.List(); err == nil {
		t.Fatal("non-JSON body must return an error instead of an empty map")
	}
}

// PG 의 GetAccessToken 은 GetToken 과 같은 요청을 만든다 (별칭).
func TestPgGetAccessTokenAliasMatchesGetToken(t *testing.T) {
	var captured []capturedCommerceRequest
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			var body []byte
			if req.Body != nil {
				body, _ = io.ReadAll(req.Body)
			}
			captured = append(captured, capturedCommerceRequest{
				Method: req.Method, URL: req.URL.String(), Header: req.Header.Clone(), Body: body,
			})
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"t","expire_in":1800}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	api := NewAPI("app_id", "private_key", client, "development")
	if _, err := api.GetToken(); err != nil {
		t.Fatal(err)
	}
	if _, err := api.GetAccessToken(); err != nil {
		t.Fatal(err)
	}
	if len(captured) != 2 {
		t.Fatalf("captured %d requests, want 2", len(captured))
	}
	if captured[0].Method != captured[1].Method || captured[0].URL != captured[1].URL {
		t.Fatalf("alias diverged: %v vs %v", captured[0], captured[1])
	}
	if string(captured[0].Body) != string(captured[1].Body) {
		t.Fatalf("alias body diverged: %s vs %s", captured[0].Body, captured[1].Body)
	}
}
