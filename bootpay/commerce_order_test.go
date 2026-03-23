package bootpay

import "testing"

func TestCommerceOrderList(t *testing.T) {
	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.Order.List(&OrderListParams{
		ListParams: ListParams{
			Page:  1,
			Limit: 10,
		},
	})
	if err != nil {
		t.Fatalf("Commerce Order.List failed: %v", err)
	}

	t.Logf("Commerce Order.List response: %+v", result)
}

func TestCommerceOrderDetail(t *testing.T) {
	t.Skip("Skipping: requires a valid order_id in the Commerce system")

	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.Order.Detail("test_order_123")
	if err != nil {
		t.Fatalf("Commerce Order.Detail failed: %v", err)
	}

	t.Logf("Commerce Order.Detail response: %+v", result)
}

func TestCommerceOrderMonth(t *testing.T) {
	t.Skip("Skipping: requires a valid user_group_id in the Commerce system")

	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.Order.Month("test_group_123", "2024-01")
	if err != nil {
		t.Fatalf("Commerce Order.Month failed: %v", err)
	}

	t.Logf("Commerce Order.Month response: %+v", result)
}
