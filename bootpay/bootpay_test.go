package bootpay

import (
	"fmt"
	"testing"
	"time"
)

// =====================================================
// 메인 테스트 함수
// =====================================================

func TestFunctions(t *testing.T) {
	bootpay := CreatePgApi()

	// 토큰 발급 (필수)
	GetToken(bootpay)

	// 추가 테스트 (읽기 전용 API)
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
		Pg:              "라이트페이",
		OrderName:       "정기결제 테스트 아이템",
		CardNo:          "5570**********1074",
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
	fmt.Println("--------------- LookupBillingKey() Start ---------------")
	verify, err := api.LookupBillingKey(TestReceiptIdBilling)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}

	fmt.Println(verify)
	fmt.Println("--------------- LookupBillingKey() End ---------------")
}

func LookupBillingKeyByKey(api *Api) {
	fmt.Println("--------------- LookupBillingKeyByKey() Start ---------------")
	verify, err := api.LookupBillingKeyByKey(TestBillingKey2)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}

	fmt.Println(verify)
	fmt.Println("--------------- LookupBillingKeyByKey() End ---------------")
}

func DestroyBillingKey(api *Api) {
	fmt.Println("--------------- DestroyBillingKey() Start ---------------")
	res, err := api.DestroyBillingKey(TestBillingKey)

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
		BillingKey: TestBillingKey,
		OrderName:  "아이템01",
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
	s10 := time.Now().Add(time.Second * 10).Format("2006-01-02T15:04:05-07:00")
	payload := SubscribePayload{
		BillingKey:       TestBillingKey,
		OrderName:        "아이템01",
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
	fmt.Println("--------------- ReserveSubscribeLookup() Start ---------------")
	res, err := api.ReserveSubscribeLookup(TestReserveId)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ReserveSubscribeLookup() End ---------------")
}

func ReserveCancel(api *Api) {
	fmt.Println("--------------- ReserveCancel() Start ---------------")
	res, err := api.ReserveCancelSubscribe(TestReserveId)

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
	fmt.Println("--------------- GetReceipt() Start ---------------")
	verify, err := api.GetReceipt(TestReceiptId)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}

	fmt.Println(verify)
	fmt.Println("--------------- GetReceipt() End ---------------")
}

func GetVerify(api *Api) {
	fmt.Println("--------------- GetVerify() Start ---------------")
	verify, err := api.GetReceipt(TestReceiptId)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}

	fmt.Println(verify)
	fmt.Println("--------------- GetVerify() End ---------------")
}

func ReceiptCancel(api *Api) {
	payload := CancelData{
		ReceiptId:      TestReceiptId,
		CancelId:       fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		CancelUsername: "관리자",
		CancelMessage:  "테스트 결제 취소",
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
	fmt.Println("--------------- ServerConfirm() Start ---------------")
	res, err := api.ServerConfirm(TestReceiptIdConfirm)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ServerConfirm() End ---------------")
}

func Certificate(api *Api) {
	fmt.Println("--------------- Certificate() Start ---------------")
	res, err := api.Certificate(TestCertificateReceiptId)

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
		UserId: TestUserId,
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
		UserId: TestUserId,
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
		ReceiptId:      TestReceiptIdEscrow,
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
		ReceiptId:       TestReceiptIdCash,
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
		ReceiptId:      TestReceiptIdCash,
		CancelUsername: "테스트 관리자",
		CancelMessage:  "테스트 결제",
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
		ReceiptId:      TestReceiptIdCash,
		CancelUsername: "테스트 관리자",
		CancelMessage:  "테스트 결제",
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
		ReceiptId: "636a0bc4d01c7e00331cd25e",
		Otp:       "556659",
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
		ReceiptId: "6369dc33d01c7e00271cccad",
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
		Pg:                    "나이스페이",
		OrderName:             "테스트 결제",
		Username:              "홍길동",
		AuthType:              "ARS",
		BankName:              "국민",
		BankAccount:           "67512341234472",
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

	res, err := api.PublishAutomaticTransferBillingKey(TestReceiptIdTransfer)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- PublishAutomaticTransferBillingKey() End ---------------")
}
