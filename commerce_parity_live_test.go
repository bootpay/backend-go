package bootpay

import "testing"

// 신규 endpoint 라이브 테스트 — development 전용 (BOOTPAY_ENV=development 일 때만 실행).
// production 호출 금지 정책에 따라 기본(production) 환경에서는 skip 된다.
// 조회성(read-only) endpoint 만 호출한다.
// (skipUnlessDevelopment 가드는 config_test.go — 모든 라이브 테스트 공용)

func TestLiveCommerceMallSettingGet(t *testing.T) {
	skipUnlessDevelopment(t)
	commerce := CreateCommerceApi()
	if _, err := commerce.GetAccessToken(); err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.MallSetting.GetMallSetting("")
	if err != nil {
		t.Fatalf("Commerce MallSetting.GetMallSetting failed: %v", err)
	}
	t.Logf("Commerce MallSetting.GetMallSetting response: %+v", result)
}

func TestLiveCommerceUidExist(t *testing.T) {
	skipUnlessDevelopment(t)
	commerce := CreateCommerceApi()
	if _, err := commerce.GetAccessToken(); err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.User.UidExist("go-sdk-live-test-uid", "")
	if err != nil {
		t.Fatalf("Commerce User.UidExist failed: %v", err)
	}
	t.Logf("Commerce User.UidExist response: %+v", result)
}

func TestLiveCommerceMallProducts(t *testing.T) {
	skipUnlessDevelopment(t)
	commerce := CreateCommerceApi()
	if _, err := commerce.GetAccessToken(); err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.Product.Products(nil)
	if err != nil {
		t.Fatalf("Commerce Product.Products failed: %v", err)
	}
	t.Logf("Commerce Product.Products response: %+v", result)
}

func TestLiveCommerceInvoiceList(t *testing.T) {
	skipUnlessDevelopment(t)
	commerce := CreateCommerceApi()
	if _, err := commerce.GetAccessToken(); err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.Invoice.ListWithParams(nil)
	if err != nil {
		t.Fatalf("Commerce Invoice.ListWithParams failed: %v", err)
	}
	t.Logf("Commerce Invoice.ListWithParams response: %+v", result)
}

func TestLiveCommerceOrderSubscriptionRequestList(t *testing.T) {
	skipUnlessDevelopment(t)
	commerce := CreateCommerceApi()
	if _, err := commerce.GetAccessToken(); err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.OrderSubscriptionRequest.List(nil)
	if err != nil {
		t.Fatalf("Commerce OrderSubscriptionRequest.List failed: %v", err)
	}
	t.Logf("Commerce OrderSubscriptionRequest.List response: %+v", result)
}
