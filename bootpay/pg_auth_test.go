package bootpay

import (
	"fmt"
	"testing"
	"time"
)

func TestPgRequestAuthentication(t *testing.T) {
	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	authData := Authentication{
		Pg:               "danal",
		Method:           "identity",
		Username:         "test_user",
		IdentityNo:       "0000000",
		Carrier:          "SKT",
		Phone:            "01010002000",
		SiteUrl:          "https://www.bootpay.co.kr",
		OrderName:        "identity verification test",
		AuthenticationId: fmt.Sprintf("%+8d", time.Now().UnixNano()/int64(time.Millisecond)),
		ClientIp:         "127.0.0.1",
		AuthenticateType: "sms",
	}

	result, err := api.RequestAuthentication(authData)
	if err != nil {
		t.Fatalf("RequestAuthentication failed: %v", err)
	}

	t.Logf("RequestAuthentication response: %+v", result)
}

func TestPgConfirmAuthentication(t *testing.T) {
	t.Skip("Skipping: requires a valid receipt_id and OTP from a real authentication request")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	params := AuthenticationParams{
		ReceiptId: "636a0bc4d01c7e00331cd25e",
		Otp:       "123456",
	}

	result, err := api.ConfirmAuthentication(params)
	if err != nil {
		t.Fatalf("ConfirmAuthentication failed: %v", err)
	}

	t.Logf("ConfirmAuthentication response: %+v", result)
}

func TestPgRealarmAuthentication(t *testing.T) {
	t.Skip("Skipping: requires a valid receipt_id from a real authentication request")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	params := AuthenticationParams{
		ReceiptId: "6369dc33d01c7e00271cccad",
	}

	result, err := api.RealarmAuthentication(params)
	if err != nil {
		t.Fatalf("RealarmAuthentication failed: %v", err)
	}

	t.Logf("RealarmAuthentication response: %+v", result)
}
