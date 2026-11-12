package bootpay

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestCommercePayoutEndpointContract(t *testing.T) {
	cases := []struct {
		name, method, path string
		call               func(*CommerceApi) (map[string]interface{}, error)
	}{
		{"Create", http.MethodPost, "sdk/payouts", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Payout.Create(PayoutCreateParams{SubSellerId: "ss-1", Amount: 150000, MerchantPayoutId: "mp-1"})
		}},
		{"Detail", http.MethodGet, "sdk/payouts/mp-1", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Payout.Detail("mp-1")
		}},
		{"List", http.MethodGet, "sdk/payouts", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Payout.List(nil)
		}},
		{"DetailEscape", http.MethodGet, "sdk/payouts/a%2Fb%20c%3Fd%23e%ED%95%9C", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Payout.Detail("a/b c?d#e한")
		}},
	}
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
				if key == "" || keys[key] {
					t.Fatalf("expected fresh Idempotency-Key, got %q", key)
				}
				keys[key] = true
			}
		})
	}
}

func TestCommercePayoutExplicitIdempotencyKey(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)
	calls := []func() (map[string]interface{}, error){
		func() (map[string]interface{}, error) {
			return api.Payout.Create(PayoutCreateParams{SubSellerId: "ss-1", Amount: 1, MerchantPayoutId: "mp-1", IdempotencyKey: "idem-payout"})
		},
		func() (map[string]interface{}, error) { return api.Payout.Detail("mp-1", "idem-payout") },
		func() (map[string]interface{}, error) {
			return api.Payout.List(&PayoutListParams{IdempotencyKey: "idem-payout"})
		},
	}
	for _, call := range calls {
		if _, err := call(); err != nil {
			t.Fatal(err)
		}
		if got := lastRequest(t, captured).Header.Get("Idempotency-Key"); got != "idem-payout" {
			t.Fatalf("Idempotency-Key=%q", got)
		}
	}
}

func TestCommercePayoutCreateCompactBody(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)
	params := PayoutCreateParams{SubSellerId: "ss-1", MerchantPayoutId: "mp-1", IdempotencyKey: "header-only"}
	if _, err := api.Payout.Create(params); err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{"sub_seller_id": "ss-1", "amount": float64(0), "merchant_payout_id": "mp-1"}
	if got := decodeBody(t, lastRequest(t, captured)); !reflect.DeepEqual(got, want) {
		t.Fatalf("compact body=%+v, want %+v", got, want)
	}

	params.Amount = 150000
	params.OrderNumber = "ORDER-1"
	params.PayoutDate = "2026-11-20"
	if _, err := api.Payout.Create(params); err != nil {
		t.Fatal(err)
	}
	want["amount"] = float64(150000)
	want["order_number"] = "ORDER-1"
	want["payout_date"] = "2026-11-20"
	if got := decodeBody(t, lastRequest(t, captured)); !reflect.DeepEqual(got, want) {
		t.Fatalf("optional body=%+v, want %+v", got, want)
	}
}

func TestCommercePayoutListQuery(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)
	if _, err := api.Payout.List(nil); err != nil {
		t.Fatal(err)
	}
	if got := lastRequest(t, captured).URL; got != COMMERCE_DEVELOPMENT+"/sdk/payouts" {
		t.Fatalf("nil params URL=%s", got)
	}
	if _, err := api.Payout.List(&PayoutListParams{}); err != nil {
		t.Fatal(err)
	}
	if got := lastRequest(t, captured).URL; got != COMMERCE_DEVELOPMENT+"/sdk/payouts" {
		t.Fatalf("empty params URL=%s", got)
	}

	params := &PayoutListParams{Status: PAYOUT_STATUS_READY, SubSellerId: "ss-1", From: "2026-11-01", To: "2026-11-30", Page: 2, Limit: 50}
	if _, err := api.Payout.List(params); err != nil {
		t.Fatal(err)
	}
	query := lastRequest(t, captured).URL
	for _, item := range []string{"status=ready", "sub_seller_id=ss-1", "from=2026-11-01", "to=2026-11-30", "page=2", "limit=50"} {
		if !strings.Contains(query, item) {
			t.Fatalf("query %q missing from %s", item, query)
		}
	}
}

func TestCommercePayoutBareArrayAndResponseTypes(t *testing.T) {
	response := `[{"id":"p-1","merchant_payout_id":"mp-1","sub_seller_id":"ss-1","amount":150000,"status":"ready","payout_date":"2026-11-20","order_number":null,"last_code":null,"last_message":null,"created_at":"2026-11-12T10:00:00+09:00","updated_at":"2026-11-12T10:00:00+09:00"},{"id":"p-2","merchant_payout_id":"mp-2","sub_seller_id":"ss-2","amount":200000,"status":"succeeded","payout_date":"2026-11-21","order_number":"ORDER-2","last_code":"0000","last_message":"완료","created_at":"2026-11-12T11:00:00+09:00","updated_at":"2026-11-12T12:00:00+09:00"}]`
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(response)), Header: make(http.Header)}, nil
	})}
	api := NewCommerceAPI("ck", "sk", client, "development")
	result, err := api.Payout.List(nil)
	if err != nil {
		t.Fatal(err)
	}
	data, ok := result["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Fatalf("bare array envelope=%#v", result)
	}
	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	var payouts []CommercePayout
	if err := json.Unmarshal(raw, &payouts); err != nil {
		t.Fatal(err)
	}
	if payouts[0].Status != PAYOUT_STATUS_READY || payouts[0].Amount != 150000 || payouts[0].OrderNumber != nil || payouts[0].LastCode != nil || payouts[0].LastMessage != nil {
		t.Fatalf("first payout=%+v", payouts[0])
	}
	if payouts[1].Status != PAYOUT_STATUS_SUCCEEDED || payouts[1].OrderNumber == nil || *payouts[1].OrderNumber != "ORDER-2" || payouts[1].LastCode == nil || *payouts[1].LastCode != "0000" || payouts[1].LastMessage == nil || *payouts[1].LastMessage != "완료" {
		t.Fatalf("second payout=%+v", payouts[1])
	}
}

func TestCommercePayoutRequiredIdValidation(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)
	if _, err := api.Payout.Detail(""); err == nil {
		t.Fatal("Detail: empty id must return an error")
	}
	if len(captured) != 0 {
		t.Fatalf("no request must be sent for an empty id, got %d", len(captured))
	}
}
