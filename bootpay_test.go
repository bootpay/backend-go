package bootpay_go

import (
	"fmt"
	"testing"
	"time"
)

func TestGetBillingKey(t *testing.T) {
	bc := Client{}.New("5b8f6a4d396fa665fdc2b5ea", "rm6EYECr6aroQVG2ntW0A6LpWnkTgP4uQ3H18sDDUYw=", nil, "")
	GetToken(bc)
	GetBillingKey(bc)
	GetVerify(bc)
	ReceiptCancel(bc)
	RequestLink(bc)
	ServerSubmit(bc)
}

func GetToken(bc *Client) {
	fmt.Println("--------------- GetToken() Start ---------------")
	token, err := bc.GetToken()
	fmt.Println("token : " + token.Data.Token)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- GetToken() End ---------------")
}

func GetBillingKey(bc *Client) {
	fmt.Println("--------------- GetBillingKey() Start ---------------")
	payload := BillingKeyPayload{
		OrderId: fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		Pg: "nicepay",
		ItemName: "정기결제 테스트 아이템",
		CardNo: "5570********1074",
		CardPw: "**",
		ExpireYear: "**",
		ExpireMonth: "**",
		IdentifyNumber: "",
	}
	billingKey, err := bc.GetBillingKey(payload)

	fmt.Println(billingKey)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- GetBillingKey() End ---------------")
}

func GetVerify(bc *Client) {
	receiptId := "610c96352386840036db8bef"
	fmt.Println("--------------- GetVerify() Start ---------------")
	verify, err := bc.Verify(receiptId)

	fmt.Println(verify)
	fmt.Println(verify.Data.PaymentData)
	fmt.Println(verify.Data.PaymentData["o_id"])
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- GetVerify() End ---------------")
}

func ReceiptCancel(bc *Client) {
	receiptId := "610cc0cb7b5ba40044b04530"
	name := "관리자"
	reason := "테스트 결제 취소를 테스트"
	fmt.Println("--------------- ReceiptCancel() Start ---------------")
	cancel, err := bc.ReceiptCancel(receiptId, 0, name, reason, RefundData{})

	fmt.Println(cancel)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- ReceiptCancel() End ---------------")
}

func RequestLink(bc *Client) {
	payload := Payload{
		Pg: "kcp",
		Method: "card",
		Price: 1000,
		OrderId: fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		Name: "테스트 결제 상품",
	}
	fmt.Println("--------------- RequestLink() End ---------------")
	res, err := bc.RequestLink(payload)

	fmt.Println(res)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- RequestLink() End ---------------")
}

func ServerSubmit(bc *Client) {
	receiptId := "610cc01b238684002adb904e"
	fmt.Println("--------------- ServerSubmit() Start ---------------")
	res, err := bc.ServerSubmit(receiptId)

	fmt.Println(res)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- ServerSubmit() End ---------------")
}
