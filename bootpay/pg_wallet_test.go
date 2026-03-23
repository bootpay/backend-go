package bootpay

import (
	"fmt"
	"testing"
	"time"
)

func TestPgGetUserWallets(t *testing.T) {
	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	result, err := api.GetUserWallets(TestUserId, true)
	if err != nil {
		t.Fatalf("GetUserWallets failed: %v", err)
	}

	t.Logf("GetUserWallets response: %+v", result)
}

func TestPgRequestWalletPayment(t *testing.T) {
	t.Skip("Skipping: requires a valid user with wallet setup")

	api := CreatePgApi()
	_, err := api.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	request := WalletRequest{
		UserId:    TestUserId,
		OrderName: "wallet payment test",
		Price:     1000,
		OrderId:   fmt.Sprintf("%+8d", time.Now().UnixNano()/int64(time.Millisecond)),
		Sandbox:   true,
		Items: []Item{
			{
				Name:  "test item",
				Qty:   1,
				Id:    "item_1",
				Price: 1000,
			},
		},
	}

	result, err := api.RequestWalletPayment(request)
	if err != nil {
		t.Fatalf("RequestWalletPayment failed: %v", err)
	}

	t.Logf("RequestWalletPayment response: %+v", result)
}
