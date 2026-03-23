package bootpay

import (
	"fmt"
	"testing"
	"time"
)

func TestPgGetReceipt(t *testing.T) {
	t.Skip("Skipping: requires a valid receipt_id from a real payment")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	result, err := api.GetReceipt(TestReceiptId)
	if err != nil {
		t.Fatalf("GetReceipt failed: %v", err)
	}

	t.Logf("GetReceipt response: %+v", result)
}

func TestPgGetReceiptWithUserData(t *testing.T) {
	t.Skip("Skipping: requires a valid receipt_id from a real payment")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	result, err := api.GetReceiptWithUserData(TestReceiptId, true)
	if err != nil {
		t.Fatalf("GetReceiptWithUserData failed: %v", err)
	}

	t.Logf("GetReceiptWithUserData response: %+v", result)
}

func TestPgCancelPayment(t *testing.T) {
	t.Skip("Skipping: requires a valid receipt_id from a real payment to cancel")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	cancelData := CancelData{
		ReceiptId:      TestReceiptId,
		CancelId:       fmt.Sprintf("%+8d", time.Now().UnixNano()/int64(time.Millisecond)),
		CancelUsername: "test_admin",
		CancelMessage:  "integration test cancel",
	}

	result, err := api.ReceiptCancel(cancelData)
	if err != nil {
		t.Fatalf("ReceiptCancel failed: %v", err)
	}

	t.Logf("ReceiptCancel response: %+v", result)
}

func TestPgConfirmPayment(t *testing.T) {
	t.Skip("Skipping: requires a valid receipt_id pending confirmation")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	result, err := api.ServerConfirm(TestReceiptIdConfirm)
	if err != nil {
		t.Fatalf("ServerConfirm failed: %v", err)
	}

	t.Logf("ServerConfirm response: %+v", result)
}

func TestPgCertificate(t *testing.T) {
	t.Skip("Skipping: requires a valid certificate receipt_id")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	result, err := api.Certificate(TestCertificateReceiptId)
	if err != nil {
		t.Fatalf("Certificate failed: %v", err)
	}

	t.Logf("Certificate response: %+v", result)
}

func TestPgRequestUserToken(t *testing.T) {
	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	payload := UserTokenRequest{
		UserId: TestUserId,
		Phone:  "01012345678",
	}

	result, err := api.RequestUserToken(payload)
	if err != nil {
		t.Fatalf("RequestUserToken failed: %v", err)
	}

	t.Logf("RequestUserToken response: %+v", result)
}
