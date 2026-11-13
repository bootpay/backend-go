package bootpay

import (
	"encoding/json"
	"net/http"
	"reflect"
	"regexp"
	"testing"
)

// 입금대기는 알림톡 메시지 호스트가 아니라 커머스 /v1 로 요청한다.
func TestCommerceDepositEndpointContract(t *testing.T) {
	cases := []struct {
		name, method, path string
		call               func(*CommerceApi) (map[string]interface{}, error)
	}{
		{"Accounts", http.MethodGet, "sdk/deposits/accounts", func(api *CommerceApi) (map[string]interface{}, error) { return api.Deposit.Accounts() }},
		{"WebhookSecret", http.MethodPost, "sdk/deposits/webhook_secret", func(api *CommerceApi) (map[string]interface{}, error) { return api.Deposit.WebhookSecret() }},
		{"Create", http.MethodPost, "sdk/deposits", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Deposit.Create(DepositCreateParams{MerchantOrderId: "m-1", Price: 15000, Depositor: "홍길동"})
		}},
		{"Detail", http.MethodGet, "sdk/deposits/m-1", func(api *CommerceApi) (map[string]interface{}, error) { return api.Deposit.Detail("m-1") }},
		{"Cancel", http.MethodDelete, "sdk/deposits/m-1", func(api *CommerceApi) (map[string]interface{}, error) { return api.Deposit.Cancel("m-1") }},
		{"DetailEscape", http.MethodGet, "sdk/deposits/a%2Fb", func(api *CommerceApi) (map[string]interface{}, error) { return api.Deposit.Detail("a/b") }},
		{"CancelEscape", http.MethodDelete, "sdk/deposits/a%2Fb", func(api *CommerceApi) (map[string]interface{}, error) { return api.Deposit.Cancel("a/b") }},
		{"DetailSpecialCharacters", http.MethodGet, "sdk/deposits/%ED%95%9C%20%3F%23%25", func(api *CommerceApi) (map[string]interface{}, error) { return api.Deposit.Detail("한 ?#%") }},
	}
	uuid := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	keys := map[string]bool{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var captured []capturedCommerceRequest
			api := newMockCommerceApi(&captured)
			api.AsManager()
			for i := 0; i < 2; i++ {
				if _, err := tc.call(api); err != nil {
					t.Fatal(err)
				}
				req := lastRequest(t, captured)
				if req.Method != tc.method || req.URL != COMMERCE_DEVELOPMENT+"/"+tc.path {
					t.Fatalf("got %s %s", req.Method, req.URL)
				}
				if req.Header.Get("BOOTPAY-ROLE") != "supervisor" || api.role != "manager" {
					t.Fatalf("request role=%q, instance role=%q", req.Header.Get("BOOTPAY-ROLE"), api.role)
				}
				key := req.Header.Get("Idempotency-Key")
				if !uuid.MatchString(key) || keys[key] {
					t.Fatalf("expected fresh UUID, got %q", key)
				}
				keys[key] = true
				if tc.name == "WebhookSecret" && string(req.Body) != "{}" {
					t.Fatalf("secret body=%s", req.Body)
				}
			}
		})
	}
}

func TestCommerceDepositExplicitIdempotencyKey(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)
	calls := []func() (map[string]interface{}, error){
		func() (map[string]interface{}, error) {
			return api.Deposit.Create(DepositCreateParams{MerchantOrderId: "m-1", Price: 1, Depositor: "홍", IdempotencyKey: "idem-deposit"})
		},
		func() (map[string]interface{}, error) { return api.Deposit.Accounts("idem-deposit") },
		func() (map[string]interface{}, error) { return api.Deposit.WebhookSecret("idem-deposit") },
		func() (map[string]interface{}, error) { return api.Deposit.Detail("a/b", "idem-deposit") },
		func() (map[string]interface{}, error) { return api.Deposit.Cancel("a/b", "idem-deposit") },
	}
	for _, call := range calls {
		if _, err := call(); err != nil {
			t.Fatal(err)
		}
		if got := lastRequest(t, captured).Header.Get("Idempotency-Key"); got != "idem-deposit" {
			t.Fatalf("Idempotency-Key=%q", got)
		}
	}
}

func TestCommerceDepositCreateCompactBody(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)
	params := DepositCreateParams{MerchantOrderId: "m-1", Price: 15000, Depositor: "홍길동", AccountId: "acc1", IdempotencyKey: "idem-deposit",
		Notify: &DepositNotifyParams{DepositGuide: BoolPtr(true), DepositPaid: BoolPtr(false), Link: "https://example.com/o/m-1"}}
	if _, err := api.Deposit.Create(params); err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{"merchant_order_id": "m-1", "price": float64(15000), "depositor": "홍길동", "account_id": "acc1",
		"notify": map[string]interface{}{"deposit_guide": true, "deposit_paid": false, "link": "https://example.com/o/m-1"}}
	if got := decodeBody(t, lastRequest(t, captured)); !reflect.DeepEqual(got, want) {
		t.Fatalf("body=%+v, want %+v", got, want)
	}

	// nil notify 와 미지정 필드는 생략하고 필수 price 의 0 은 보존한다.
	params = DepositCreateParams{MerchantOrderId: "m-2", Depositor: "홍"}
	if _, err := api.Deposit.Create(params); err != nil {
		t.Fatal(err)
	}
	want = map[string]interface{}{"merchant_order_id": "m-2", "price": float64(0), "depositor": "홍"}
	if got := decodeBody(t, lastRequest(t, captured)); !reflect.DeepEqual(got, want) {
		t.Fatalf("compact body=%+v", got)
	}

	params.OrderName, params.Memo, params.BuyerName = "주문", "메모", "구매자"
	params.BuyerPhone, params.BuyerEmail, params.WebhookUrl = "01012345678", "buyer@example.com", "https://example.com/hook"
	params.Notify = &DepositNotifyParams{DepositGuide: BoolPtr(false)}
	if _, err := api.Deposit.Create(params); err != nil {
		t.Fatal(err)
	}
	want["order_name"], want["memo"], want["buyer_name"] = params.OrderName, params.Memo, params.BuyerName
	want["buyer_phone"], want["buyer_email"], want["webhook_url"] = params.BuyerPhone, params.BuyerEmail, params.WebhookUrl
	want["notify"] = map[string]interface{}{"deposit_guide": false}
	if got := decodeBody(t, lastRequest(t, captured)); !reflect.DeepEqual(got, want) {
		t.Fatalf("optional body=%+v", got)
	}
}

func TestCommerceDepositResponseTypes(t *testing.T) {
	// 웹훅 본문 형태 — 주문의 account 는 등록 때 스냅샷이라 primary 가 없고(계좌 목록에만 실린다), deposit 은 웹훅에만 실린다
	var deposit CommerceDeposit
	body := `{"merchant_order_id":"m-1","order_number":"o-1","order_name":"무통장 입금","amount":15000,"status":"paid","buyer_name":"구매자","depositor":"홍길동","account":{"id":"acc1","bank_name":"은행","account_no":"1234","holder":"판매자"},"deadline":"2026-10-24","deadline_label":"내일","paid_at":"2026-10-23","notify":{"deposit_guide":true,"deposit_paid":false},"deposit":{"key":"key1","bank_name":"은행","account_tail":"1234","depositor":"홍길동","amount":15000,"received_at":"2026-10-23"}}`
	if err := json.Unmarshal([]byte(body), &deposit); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(deposit)
	if err != nil {
		t.Fatal(err)
	}
	var got, want map[string]interface{}
	_ = json.Unmarshal(encoded, &got)
	_ = json.Unmarshal([]byte(body), &want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("response lost fields: %s", encoded)
	}
	if err := json.Unmarshal([]byte(`{"account":null,"deadline":null,"deadline_label":null,"paid_at":null,"notify":null,"deposit":null}`), &deposit); err != nil {
		t.Fatal(err)
	}
	if deposit.Account != nil || deposit.Deadline != nil || deposit.DeadlineLabel != nil || deposit.PaidAt != nil || deposit.Notify != nil || deposit.Deposit != nil {
		t.Fatal("nullable fields must decode to nil")
	}
	var accounts DepositAccountsResponse
	if err := json.Unmarshal([]byte(`{"accounts":[{"id":"acc1","primary":false}]}`), &accounts); err != nil {
		t.Fatal(err)
	}
	if len(accounts.Accounts) != 1 || accounts.Accounts[0].Primary == nil || *accounts.Accounts[0].Primary {
		t.Fatal("account primary false lost")
	}
	var secret DepositWebhookSecretResponse
	if err := json.Unmarshal([]byte(`{"secret":"secret"}`), &secret); err != nil || secret.Secret != "secret" {
		t.Fatalf("secret=%+v, err=%v", secret, err)
	}
}

// 식별자가 비면 네트워크로 나가기 전에 막는다 — 빈 id 는 상위 컬렉션 경로(sdk/deposits/)를 때린다.
func TestCommerceDepositRequiredIdValidation(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.Deposit.Detail(""); err == nil {
		t.Fatalf("Detail: empty id must return an error")
	}
	if _, err := api.Deposit.Cancel(""); err == nil {
		t.Fatalf("Cancel: empty id must return an error")
	}
	if len(captured) != 0 {
		t.Fatalf("no request must be sent for an empty id, got %d", len(captured))
	}
}
