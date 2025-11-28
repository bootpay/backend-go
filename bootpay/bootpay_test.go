package bootpay

import (
	"fmt"
	"testing"
	"time"
)

// 테스트용 API 인스턴스 생성
func getTestApi() *Api {
	return Api{}.New("5b8f6a4d396fa665fdc2b5ea", "rm6EYECr6aroQVG2ntW0A6LpWnkTgP4uQ3H18sDDUYw=", nil, "")
}

// =====================================================
// 메인 테스트 함수
// =====================================================

func TestFunctions(t *testing.T) {
	bootpay := getTestApi()

	// 토큰 발급 (필수)
	GetToken(bootpay)

	// 추가 테스트 (읽기 전용 API)
	GetUserWallets(bootpay)
	RequestUserToken(bootpay)

	// 아래 테스트들은 필요에 따라 주석 해제하여 실행
	// ReceiptCancel(bootpay)
	// GetReceipt(bootpay)
	// GetBillingKey(bootpay)
	// RequestSubscribe(bootpay)
	// LookupBillingKey(bootpay)
	// LookupBillingKeyByKey(bootpay)
	// ReserveSubscribe(bootpay)
	// ReserveSubscribeLookup(bootpay)
	// ReserveCancel(bootpay)
	// DestroyBillingKey(bootpay)
	// GetUserToken(bootpay)
	// GetVerify(bootpay)
	// ServerConfirm(bootpay)
	// Certificate(bootpay)
	// ShippingStart(bootpay)
	// RequestCashReceiptByBootpay(bootpay)
	// RequestCashReceiptCancelByBootpay(bootpay)
	// RequestCashReceipt(bootpay)
	// RequestCashReceiptCancel(bootpay)
	// RequestAuthentication(bootpay)
	// ConfirmAuthentication(bootpay)
	// RealarmAuthentication(bootpay)
	// RequestSubscribeAutomaticTransferBillingKey(bootpay)
	// PublishAutomaticTransferBillingKey(bootpay)
	// RequestWalletPayment(bootpay)
}

// =====================================================
// Token API
// =====================================================

func GetToken(api *Api) {
	fmt.Println("--------------- GetToken() Start ---------------")
	token, err := api.GetToken()
	fmt.Println(token)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- GetToken() End ---------------")
}

// =====================================================
// Billing Key API
// =====================================================

func GetBillingKey(api *Api) {
	fmt.Println("--------------- GetBillingKey() Start ---------------")
	subscriptId := fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond)))
	payload := BillingKeyPayload{
		SubscriptionId:  subscriptId,
		Pg:              "nicepay",
		OrderName:       "정기결제 테스트 아이템",
		CardNo:          "5570********1074",
		CardPw:          "**",
		CardExpireYear:  "**",
		CardExpireMonth: "**",
		CardIdentityNo:  "",
	}
	billingKey, err := api.GetBillingKey(payload)

	fmt.Println(billingKey)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- GetBillingKey() End ---------------")
}

func LookupBillingKey(api *Api) {
	receiptId := "62afccb3cf9f6d001b7d101d"
	fmt.Println("--------------- LookupBillingKey() Start ---------------")
	verify, err := api.LookupBillingKey(receiptId)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}

	fmt.Println(verify)
	fmt.Println("--------------- LookupBillingKey() End ---------------")
}

func LookupBillingKeyByKey(api *Api) {
	billingKey := "66542dfb4d18d5fc7b43e1b6"
	fmt.Println("--------------- LookupBillingKeyByKey() Start ---------------")
	verify, err := api.LookupBillingKeyByKey(billingKey)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}

	fmt.Println(verify)
	fmt.Println("--------------- LookupBillingKeyByKey() End ---------------")
}

func DestroyBillingKey(api *Api) {
	billingKey := "62afc52dcf9f6d001d7d1035"
	fmt.Println("--------------- DestroyBillingKey() Start ---------------")
	res, err := api.DestroyBillingKey(billingKey)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- DestroyBillingKey() End ---------------")
}

// =====================================================
// Subscribe Payment API
// =====================================================

func RequestSubscribe(api *Api) {
	payload := SubscribePayload{
		BillingKey: "62afc52dcf9f6d001d7d1035",
		OrderName:  "정기결제 테스트",
		OrderId:    fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		Price:      1000,
		Items: []Item{
			{
				Name:  "테스트 상품",
				Qty:   1,
				Id:    "item_1",
				Price: 1000,
			},
		},
	}

	fmt.Println("--------------- RequestSubscribe() Start ---------------")
	res, err := api.RequestSubscribe(payload)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- RequestSubscribe() End ---------------")
}

func ReserveSubscribe(api *Api) {
	s10 := time.Now().Add(time.Second * 100).Format("2006-01-02T15:04:05-07:00")
	payload := SubscribePayload{
		BillingKey:       "62aff193cf9f6d001a7d10be",
		OrderName:        "정기결제 테스트",
		OrderId:          fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		ReserveExecuteAt: s10,
		Price:            1000,
	}

	fmt.Println("--------------- ReserveSubscribe() Start ---------------")
	res, err := api.ReserveSubscribe(payload)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ReserveSubscribe() End ---------------")
}

func ReserveSubscribeLookup(api *Api) {
	reserveId := "6490149ca575b40024f0b70d"

	fmt.Println("--------------- ReserveSubscribeLookup() Start ---------------")
	res, err := api.ReserveSubscribeLookup(reserveId)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ReserveSubscribeLookup() End ---------------")
}

func ReserveCancel(api *Api) {
	reserveId := "62aff2a0cf9f6d001a7d10c4"
	fmt.Println("--------------- ReserveCancel() Start ---------------")
	res, err := api.ReserveCancelSubscribe(reserveId)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ReserveCancel() End ---------------")
}

// =====================================================
// Receipt/Payment API
// =====================================================

func GetReceipt(api *Api) {
	receiptId := "62afc194e38c300021b345d4"
	fmt.Println("--------------- GetReceipt() Start ---------------")
	verify, err := api.GetReceipt(receiptId)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}

	fmt.Println(verify)
	fmt.Println("--------------- GetReceipt() End ---------------")
}

func GetVerify(api *Api) {
	receiptId := "62afc3c5cf9f6d001b7d101a"
	fmt.Println("--------------- GetVerify() Start ---------------")
	verify, err := api.GetReceipt(receiptId)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}

	fmt.Println(verify)
	fmt.Println("--------------- GetVerify() End ---------------")
}

func ReceiptCancel(api *Api) {
	payload := CancelData{
		ReceiptId:      "62afc3c5cf9f6d001b7d101a",
		CancelId:       fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		CancelUsername: "관리자",
		CancelMessage:  "테스트 결제 취소를 테스트",
	}
	fmt.Println("--------------- ReceiptCancel() Start ---------------")
	res, err := api.ReceiptCancel(payload)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ReceiptCancel() End ---------------")
}

func ServerConfirm(api *Api) {
	receiptId := "62afda41cf9f6d001f7d105f"
	fmt.Println("--------------- ServerConfirm() Start ---------------")
	res, err := api.ServerConfirm(receiptId)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ServerConfirm() End ---------------")
}

func Certificate(api *Api) {
	receiptId := "6285ffa6cf9f6d0022c4346b"
	fmt.Println("--------------- Certificate() Start ---------------")
	res, err := api.Certificate(receiptId)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- Certificate() End ---------------")
}

// =====================================================
// User Token API
// =====================================================

func GetUserToken(api *Api) {
	payload := EasyUserTokenPayload{
		UserId: "user_1234",
		Email:  "test1234@gmail.com",
		Name:   "홍길동",
		Gender: 0,
		Birth:  "19861014",
		Phone:  "01012345678",
	}

	fmt.Println("--------------- GetUserToken() Start ---------------")
	res, err := api.GetUserToken(payload)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- GetUserToken() End ---------------")
}

func RequestUserToken(api *Api) {
	payload := UserTokenRequest{
		UserId: "gosomi1",
		Phone:  "01012345678",
	}

	fmt.Println("--------------- RequestUserToken() Start ---------------")
	res, err := api.RequestUserToken(payload)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- RequestUserToken() End ---------------")
}

// =====================================================
// Escrow/Shipping API
// =====================================================

func ShippingStart(api *Api) {
	shipping := Shipping{
		ReceiptId:      "628ae7ffd01c7e001e9b6066",
		ReceiptUrl:     "https://example.com/receipt",
		TrackingNumber: "123456",
		DeliveryCorp:   "CJ대한통운",
		User: ShippingUser{
			Username: "홍길동",
			Phone:    "01000000000",
			Address:  "서울특별시 종로구",
			Zipcode:  "08490",
		},
	}

	fmt.Println("--------------- ShippingStart() Start ---------------")
	res, err := api.PutShippingStart(shipping)
	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ShippingStart() End ---------------")
}

// =====================================================
// Cash Receipt API
// =====================================================

func RequestCashReceiptByBootpay(api *Api) {
	cashReceipt := CashReceiptData{
		ReceiptId:       "62e0f11f1fc192036b1b3c92",
		Username:        "테스트",
		Email:           "test@bootpay.co.kr",
		Phone:           "01000000000",
		IdentityNo:      "01000000000",
		CashReceiptType: "소득공제",
	}

	fmt.Println("--------------- RequestCashReceiptByBootpay() Start ---------------")
	res, err := api.RequestCashReceiptByBootpay(cashReceipt)
	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- RequestCashReceiptByBootpay() End ---------------")
}

func RequestCashReceiptCancelByBootpay(api *Api) {
	cancelData := CancelData{
		ReceiptId:      "62e0f11f1fc192036b1b3c92",
		CancelUsername: "테스트",
		CancelMessage:  "테스트 관리자",
	}

	fmt.Println("--------------- RequestCashReceiptCancelByBootpay() Start ---------------")
	res, err := api.RequestCashReceiptCancelByBootpay(cancelData)
	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- RequestCashReceiptCancelByBootpay() End ---------------")
}

func RequestCashReceipt(api *Api) {
	purchasedAt := time.Now().Format("2006-01-02T15:04:05-07:00")

	cashReceipt := CashReceiptData{
		Pg:              "토스",
		Price:           1000,
		OrderName:       "테스트",
		CashReceiptType: "소득공제",
		IdentityNo:      "01000000000",
		OrderId:         fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		PurchasedAt:     purchasedAt,
		Currency:        "KRW",
	}

	fmt.Println("--------------- RequestCashReceipt() Start ---------------")
	res, err := api.RequestCashReceipt(cashReceipt)
	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- RequestCashReceipt() End ---------------")
}

func RequestCashReceiptCancel(api *Api) {
	cancelData := CancelData{
		ReceiptId:      "62f4be7f1fc192036f9f4bc6",
		CancelUsername: "테스트",
		CancelMessage:  "테스트 관리자",
	}

	fmt.Println("--------------- RequestCashReceiptCancel() Start ---------------")
	res, err := api.RequestCashReceiptCancel(cancelData)
	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- RequestCashReceiptCancel() End ---------------")
}

// =====================================================
// Authentication API
// =====================================================

func RequestAuthentication(api *Api) {
	authData := Authentication{
		Pg:               "다날",
		Method:           "본인인증",
		Username:         "사용자명",
		IdentityNo:       "0000000",
		Carrier:          "SKT",
		Phone:            "01010002000",
		SiteUrl:          "https://www.bootpay.co.kr",
		OrderName:        "회원 본인인증",
		AuthenticationId: fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		ClientIp:         "127.0.0.1",
		AuthenticateType: "sms",
	}

	fmt.Println("--------------- RequestAuthentication() Start ---------------")
	res, err := api.RequestAuthentication(authData)
	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- RequestAuthentication() End ---------------")
}

func ConfirmAuthentication(api *Api) {
	authParams := AuthenticationParams{
		ReceiptId: "636a020d1fc1920373e6d8ff",
		Otp:       "613026",
	}

	fmt.Println("--------------- ConfirmAuthentication() Start ---------------")
	res, err := api.ConfirmAuthentication(authParams)
	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ConfirmAuthentication() End ---------------")
}

func RealarmAuthentication(api *Api) {
	authParams := AuthenticationParams{
		ReceiptId: "636a020d1fc1920373e6d8ff",
	}

	fmt.Println("--------------- RealarmAuthentication() Start ---------------")
	res, err := api.RealarmAuthentication(authParams)
	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- RealarmAuthentication() End ---------------")
}

// =====================================================
// Automatic Transfer Billing API
// =====================================================

func RequestSubscribeAutomaticTransferBillingKey(api *Api) {
	fmt.Println("--------------- RequestSubscribeAutomaticTransferBillingKey() Start ---------------")

	subscriptId := fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond)))

	payload := BillingKeyPayload{
		SubscriptionId:        subscriptId,
		Pg:                    "nicepay",
		OrderName:             "정기결제 테스트 아이템",
		Username:              "홍길동",
		AuthType:              "ARS",
		BankName:              "국민",
		BankAccount:           "6756123412342472",
		IdentityNo:            "901014",
		CashReceiptType:       "소득공제",
		CashReceiptIdentityNo: "01012341234",
		Phone:                 "01012341234",
	}
	res, err := api.RequestSubscribeAutomaticTransferBillingKey(payload)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- RequestSubscribeAutomaticTransferBillingKey() End ---------------")
}

func PublishAutomaticTransferBillingKey(api *Api) {
	fmt.Println("--------------- PublishAutomaticTransferBillingKey() Start ---------------")

	res, err := api.PublishAutomaticTransferBillingKey("6655069ca691573f1bb9c28a")

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- PublishAutomaticTransferBillingKey() End ---------------")
}

// =====================================================
// Wallet API (신규 추가)
// =====================================================

func GetUserWallets(api *Api) {
	fmt.Println("--------------- GetUserWallets() Start ---------------")

	res, err := api.GetUserWallets("bootpay", true)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- GetUserWallets() End ---------------")
}

func RequestWalletPayment(api *Api) {
	fmt.Println("--------------- RequestWalletPayment() Start ---------------")

	payload := WalletRequest{
		UserId:    "bootpay",
		OrderName: "테스트 결제",
		OrderId:   fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		Price:     100,
		Sandbox:   true,
		User: User{
			Phone:    "01012341234",
			Username: "홍길동",
			Email:    "test@bootpay.co.kr",
		},
	}

	res, err := api.RequestWalletPayment(payload)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- RequestWalletPayment() End ---------------")
}
