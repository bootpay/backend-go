package bootpay

import "testing"

func TestCommerceStoreInfo(t *testing.T) {
	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.Store.Info()
	if err != nil {
		t.Fatalf("Commerce Store.Info failed: %v", err)
	}

	t.Logf("Commerce Store.Info response: %+v", result)
}

func TestCommerceStoreDetail(t *testing.T) {
	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.Store.Detail()
	if err != nil {
		t.Fatalf("Commerce Store.Detail failed: %v", err)
	}

	t.Logf("Commerce Store.Detail response: %+v", result)
}
