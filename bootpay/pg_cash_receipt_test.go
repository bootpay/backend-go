package bootpay

import (
	"fmt"
	"testing"
	"time"
)

func TestPgRequestCashReceipt(t *testing.T) {
	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	purchasedAt := time.Now().Format("2006-01-02T15:04:05-07:00")
	cashReceipt := CashReceiptData{
		Pg:              "toss",
		Price:           1000,
		OrderName:       "cash receipt test",
		CashReceiptType: "income_deduction",
		IdentityNo:      "01000000000",
		OrderId:         fmt.Sprintf("%+8d", time.Now().UnixNano()/int64(time.Millisecond)),
		PurchasedAt:     purchasedAt,
	}

	result, err := api.RequestCashReceipt(cashReceipt)
	if err != nil {
		t.Fatalf("RequestCashReceipt failed: %v", err)
	}

	t.Logf("RequestCashReceipt response: %+v", result)
}

func TestPgCashReceiptPublishOnReceipt(t *testing.T) {
	t.Skip("Skipping: requires a valid receipt_id from a real payment")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	cashReceipt := CashReceiptData{
		ReceiptId:       TestReceiptIdCash,
		Username:        "test_user",
		Email:           "test@bootpay.co.kr",
		Phone:           "01000000000",
		IdentityNo:      "01000000000",
		CashReceiptType: "income_deduction",
	}

	result, err := api.RequestCashReceiptByBootpay(cashReceipt)
	if err != nil {
		t.Fatalf("RequestCashReceiptByBootpay failed: %v", err)
	}

	t.Logf("RequestCashReceiptByBootpay response: %+v", result)
}

func TestPgCashReceiptCancelOnReceipt(t *testing.T) {
	t.Skip("Skipping: requires a valid receipt_id with a published cash receipt")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	cancelData := CancelData{
		ReceiptId:      TestReceiptIdCash,
		CancelUsername: "test_admin",
		CancelMessage:  "cash receipt cancel test",
	}

	result, err := api.RequestCashReceiptCancelByBootpay(cancelData)
	if err != nil {
		t.Fatalf("RequestCashReceiptCancelByBootpay failed: %v", err)
	}

	t.Logf("RequestCashReceiptCancelByBootpay response: %+v", result)
}

func TestPgCashReceiptCancel(t *testing.T) {
	t.Skip("Skipping: requires a valid receipt_id from a standalone cash receipt")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	cancelData := CancelData{
		ReceiptId:      TestReceiptIdCash,
		CancelUsername: "test_admin",
		CancelMessage:  "standalone cash receipt cancel test",
	}

	result, err := api.RequestCashReceiptCancel(cancelData)
	if err != nil {
		t.Fatalf("RequestCashReceiptCancel failed: %v", err)
	}

	t.Logf("RequestCashReceiptCancel response: %+v", result)
}
