package bootpay

import (
	"fmt"
	"testing"
	"time"
)

func TestPgRequestSubscribeBillingKey(t *testing.T) {
	t.Skip("Skipping: requires real card information")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	subscriptionId := fmt.Sprintf("%+8d", time.Now().UnixNano()/int64(time.Millisecond))
	payload := BillingKeyPayload{
		SubscriptionId:  subscriptionId,
		Pg:              "nicepay",
		OrderName:       "billing key test item",
		CardNo:          "5570**********1074",
		CardPw:          "**",
		CardExpireYear:  "**",
		CardExpireMonth: "**",
		CardIdentityNo:  "******",
	}

	result, err := api.GetBillingKey(payload)
	if err != nil {
		t.Fatalf("GetBillingKey failed: %v", err)
	}

	t.Logf("GetBillingKey response: %+v", result)
}

func TestPgLookupBillingKey(t *testing.T) {
	t.Skip("Skipping: requires a valid billing receipt_id")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	result, err := api.LookupBillingKey(TestReceiptIdBilling)
	if err != nil {
		t.Fatalf("LookupBillingKey failed: %v", err)
	}

	t.Logf("LookupBillingKey response: %+v", result)
}

func TestPgLookupBillingKeyByKey(t *testing.T) {
	t.Skip("Skipping: requires a valid billing_key")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	result, err := api.LookupBillingKeyByKey(TestBillingKey2)
	if err != nil {
		t.Fatalf("LookupBillingKeyByKey failed: %v", err)
	}

	t.Logf("LookupBillingKeyByKey response: %+v", result)
}

func TestPgDestroyBillingKey(t *testing.T) {
	t.Skip("Skipping: requires a valid billing_key to destroy")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	result, err := api.DestroyBillingKey(TestBillingKey)
	if err != nil {
		t.Fatalf("DestroyBillingKey failed: %v", err)
	}

	t.Logf("DestroyBillingKey response: %+v", result)
}

func TestPgRequestSubscribePayment(t *testing.T) {
	t.Skip("Skipping: requires a valid billing_key")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	payload := SubscribePayload{
		BillingKey: TestBillingKey,
		OrderName:  "subscribe payment test",
		OrderId:    fmt.Sprintf("%+8d", time.Now().UnixNano()/int64(time.Millisecond)),
		Price:      1000,
		Items: []Item{
			{
				Name:  "test item",
				Qty:   1,
				Id:    "item_1",
				Price: 1000,
			},
		},
	}

	result, err := api.RequestSubscribe(payload)
	if err != nil {
		t.Fatalf("RequestSubscribe failed: %v", err)
	}

	t.Logf("RequestSubscribe response: %+v", result)
}

func TestPgReserveSubscribePayment(t *testing.T) {
	t.Skip("Skipping: requires a valid billing_key")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	executeAt := time.Now().Add(10 * time.Second).Format("2006-01-02T15:04:05-07:00")
	payload := SubscribePayload{
		BillingKey:       TestBillingKey,
		OrderName:        "reserve subscribe test",
		OrderId:          fmt.Sprintf("%+8d", time.Now().UnixNano()/int64(time.Millisecond)),
		ReserveExecuteAt: executeAt,
		Price:            1000,
	}

	result, err := api.ReserveSubscribe(payload)
	if err != nil {
		t.Fatalf("ReserveSubscribe failed: %v", err)
	}

	t.Logf("ReserveSubscribe response: %+v", result)
}

func TestPgReserveSubscribeLookup(t *testing.T) {
	t.Skip("Skipping: requires a valid reserve_id")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	result, err := api.ReserveSubscribeLookup(TestReserveId)
	if err != nil {
		t.Fatalf("ReserveSubscribeLookup failed: %v", err)
	}

	t.Logf("ReserveSubscribeLookup response: %+v", result)
}

func TestPgReserveCancelSubscribe(t *testing.T) {
	t.Skip("Skipping: requires a valid reserve_id to cancel")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	result, err := api.ReserveCancelSubscribe(TestReserveId)
	if err != nil {
		t.Fatalf("ReserveCancelSubscribe failed: %v", err)
	}

	t.Logf("ReserveCancelSubscribe response: %+v", result)
}

func TestPgRequestSubscribeAutomaticTransferBillingKey(t *testing.T) {
	t.Skip("Skipping: requires real bank account information")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	subscriptionId := fmt.Sprintf("%+8d", time.Now().UnixNano()/int64(time.Millisecond))
	payload := BillingKeyPayload{
		SubscriptionId:        subscriptionId,
		Pg:                    "nicepay",
		OrderName:             "auto transfer billing test",
		Username:              "test_user",
		AuthType:              "ARS",
		BankName:              "kookmin",
		BankAccount:           "67512341234472",
		IdentityNo:            "901014",
		CashReceiptType:       "income_deduction",
		CashReceiptIdentityNo: "01012341234",
		Phone:                 "01012341234",
	}

	result, err := api.RequestSubscribeAutomaticTransferBillingKey(payload)
	if err != nil {
		t.Fatalf("RequestSubscribeAutomaticTransferBillingKey failed: %v", err)
	}

	t.Logf("RequestSubscribeAutomaticTransferBillingKey response: %+v", result)
}

func TestPgPublishAutomaticTransferBillingKey(t *testing.T) {
	t.Skip("Skipping: requires a valid receipt_id from automatic transfer request")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	result, err := api.PublishAutomaticTransferBillingKey(TestReceiptIdTransfer)
	if err != nil {
		t.Fatalf("PublishAutomaticTransferBillingKey failed: %v", err)
	}

	t.Logf("PublishAutomaticTransferBillingKey response: %+v", result)
}
