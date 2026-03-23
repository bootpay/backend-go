package bootpay

import "testing"

func TestCommerceProductList(t *testing.T) {
	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.Product.List(&ProductListParams{
		ListParams: ListParams{
			Page:  1,
			Limit: 10,
		},
	})
	if err != nil {
		t.Fatalf("Commerce Product.List failed: %v", err)
	}

	t.Logf("Commerce Product.List response: %+v", result)
}

func TestCommerceProductDetail(t *testing.T) {
	t.Skip("Skipping: requires a valid product_id in the Commerce system")

	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.Product.Detail("test_product_123")
	if err != nil {
		t.Fatalf("Commerce Product.Detail failed: %v", err)
	}

	t.Logf("Commerce Product.Detail response: %+v", result)
}

func TestCommerceProductCreateSimpleAndDelete(t *testing.T) {
	t.Skip("Skipping: creates and deletes a real product in the Commerce system")

	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	newProduct := CommerceProduct{
		Name:         "Integration Test Product",
		DisplayPrice: 10000,
		TaxFreePrice: 0,
		Type:         1,
	}

	result, err := commerce.Product.CreateSimple(newProduct)
	if err != nil {
		t.Fatalf("Commerce Product.CreateSimple failed: %v", err)
	}

	t.Logf("Commerce Product.CreateSimple response: %+v", result)

	// Cleanup: delete the created product if we have a product_id
	if productId, ok := result["product_id"].(string); ok && productId != "" {
		deleteResult, err := commerce.Product.Delete(productId)
		if err != nil {
			t.Logf("Commerce Product.Delete cleanup failed: %v", err)
		} else {
			t.Logf("Commerce Product.Delete cleanup response: %+v", deleteResult)
		}
	}
}
