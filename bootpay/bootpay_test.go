package bootpay

import ( 
	"fmt"
	"testing"
	"time"
)

func TestFunctions(t *testing.T) {
	bootpay := Api{}.New("5b8f6a4d396fa665fdc2b5ea", "rm6EYECr6aroQVG2ntW0A6LpWnkTgP4uQ3H18sDDUYw=", nil, "")
	//bootpay := Api{}.New("59bfc738e13f337dbd6ca48a", "pDc0NwlkEX3aSaHTp/PPL/i8vn5E/CqRChgyEp/gHD0=", nil, "development")

	GetToken(bootpay)
	//ReceiptCancel(bootpay)
	//GetReceipt(bootpay)
	// GetBillingKey(bootpay)
	//RequestSubscribe(bootpay)
	// LookupBillingKey(bootpay)
	// LookupBillingKeyByKey(bootpay)
// 	LookupSubscribeBillingKey(bootpay)
	//ReserveSubscribe(bootpay)
	//ReserveCancel(bootpay)
	//DestroyBillingKey(bootpay)
	//GetUserToken(bootpay)
	//GetVerify(bootpay)
	//ServerConfirm(bootpay)
	//Certificate(bootpay)
	//ShoppingStart(bootpay)
	//
	//RequestCashReceiptByBootpay(bootpay)
	//RequestCashReceiptCancelByBootpay(bootpay)
	//RequestCashReceipt(bootpay)
	//RequestCashReceiptCancel(bootpay)

// 	RequestAuthentication(bootpay)
	//ConfirmAuthentication(bootpay)
	//RealarmAuthentication(bootpay)
	requestSubscribeAutomaticTransferBillingKey(bootpay)
	// publishAutomaticTransferBillingKey(bootpay)
}

func GetToken(api *Api) {
	fmt.Println("--------------- GetToken() Start ---------------")
	token , err := api.GetToken()
	fmt.Println(token)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- GetToken() End ---------------")
}

func GetBillingKey(api *Api) {
	fmt.Println("--------------- GetBillingKey() Start ---------------")
	subscriptId := fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond)))
	payload := BillingKeyPayload{
		SubscriptionId: subscriptId,
		Pg: "nicepay",
		OrderName: "정기결제 테스트 아이템",
		CardNo: "5570********1074",
		CardPw: "**",
		CardExpireYear: "**",
		CardExpireMonth: "**",
		CardIdentityNo: "",
	}
	billingKey, err := api.GetBillingKey(payload)

	fmt.Println(billingKey)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- GetBillingKey() End ---------------")
}


func GetReceipt(api *Api) {
	receiptId := "62afc194e38c300021b345d4"
	fmt.Println("--------------- getReceipt() Start ---------------")
	verify, err := api.GetReceipt(receiptId)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}

	fmt.Println(verify)
	fmt.Println("--------------- GetVerify() End ---------------")
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


func ReceiptCancel(api *Api) {
	payload := CancelData{
		ReceiptId: "62afc3c5cf9f6d001b7d101a",
		CancelId:  fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		CancelUsername: "관리자",
		CancelMessage: "테스트 결제 취소를 테스트",
	}
	//receiptId := "610cc0cb7b5ba40044b04530"
	//name := "관리자"
	//reason := "테스트 결제 취소를 테스트"
	fmt.Println("--------------- ReceiptCancel() Start ---------------")
	cancel, err := api.ReceiptCancel(payload)

	fmt.Println(cancel)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ReceiptCancel() End ---------------")
}

func RequestSubscribe(api *Api) {
	payload := SubscribePayload{
		BillingKey: "62afc52dcf9f6d001d7d1035",
		OrderName: "정기결제 테스트",
		OrderId:  fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		Price: 1000,
	}

	fmt.Println("--------------- requestSubscribe() Start ---------------")
	cancel, err := api.RequestSubscribe(payload)

	fmt.Println(cancel)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- requestSubscribe() End ---------------")
}

func ReserveSubscribe(api *Api) {
	s10 := time.Now().Add(time.Second * 100).Format("2006-01-02T15:04:05-07:00")
	payload := SubscribePayload{
		BillingKey: "62aff193cf9f6d001a7d10be",
		OrderName: "정기결제 테스트",
		OrderId:  fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		ReserveExecuteAt: s10,
		Price: 1000,
	}

	fmt.Println("--------------- ReserveSubscribe() Start ---------------")
	cancel, err := api.ReserveSubscribe(payload)

	fmt.Println(cancel)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ReserveSubscribe() End ---------------")
}


func ReserveCancel(api *Api) {
	reserveId := "62aff2a0cf9f6d001a7d10c4"
	fmt.Println("--------------- ReserveCancel() Start ---------------")
	cancel, err := api.ReserveCancelSubscribe(reserveId)

	fmt.Println(cancel)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ReserveCancel() End ---------------")
}


func RequestLink(api *Api) {
	//payload := Payload{
	//	Pg: "kcp",
	//	Method: "card",
	//	Price: 1000,
	//	OrderId: fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
	//	Name: "테스트 결제 상품",
	//}
	//fmt.Println("--------------- RequestLink() End ---------------")
	//res, err := api.RequestLink(payload)
	//
	//fmt.Println(res)
	//if err != nil {
	//	fmt.Println("error: " + err.Error())
	//}
	//fmt.Println("--------------- RequestLink() End ---------------")
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


func GetUserToken(api *Api) {
	payload := EasyUserTokenPayload{
		UserId: "user_1234",
		Email: "test1234@gmail.com",
		Name: "홍길동",
		Gender: 0,
		Birth: "19861014",
		Phone: "01012345678",
	}

	fmt.Println("--------------- GetUserToken() Start ---------------")
	cancel, err := api.GetUserToken(payload)

	fmt.Println(cancel)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- GetUserToken() End ---------------")
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

func ShoppingStart(api *Api) {
	shipping := Shipping{
		ReceiptId: "628ae7ffd01c7e001e9b6066",
		TrackingNumber: "123456",
		DeliveryCorp: "CJ대한통운",
		User: ShippingUser{
			Username: "홍길동",
			Phone: "01000000000",
			Address: "서울특별시 종로구",
			Zipcode: "08490",
		},
	}

	fmt.Println("--------------- ShoppingStart() Start ---------------")
	res, err := api.PutShippingStart(shipping)
	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- ShoppingStart() End ---------------")
}


func RequestCashReceiptByBootpay(api *Api) {
	cashReceipt := CashReceiptData{
		ReceiptId: "62e0f11f1fc192036b1b3c92",
		Username: "테스트",
		Email: "test@bootpay.co.kr",
		Phone: "01000000000",
		IdentityNo: "01000000000",
		CashReceiptType: "01000000000",
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
		ReceiptId: "62e0f11f1fc192036b1b3c92",
		CancelUsername: "테스트",
		CancelMessage: "테스트 관리자",
	}

	fmt.Println("--------------- RequestCashReceiptCancelByBootpay() Start ---------------")
	fmt.Println("2135554")
	res, err := api.RequestCashReceiptCancelByBootpay(cancelData)
	fmt.Println("2134")
	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- RequestCashReceiptCancelByBootpay() End ---------------")
}

func RequestCashReceipt(api *Api) {
	purchasedAt := time.Now().Format("2006-01-02T15:04:05-07:00")

	cashReceipt := CashReceiptData{
		Pg: "토스",
		Price: 1000,
		OrderName: "테스트",
		CashReceiptType: "소득공제",
		IdentityNo: "01000000000",
		OrderId:  fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		PurchasedAt:  purchasedAt,
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
		ReceiptId: "62f4be7f1fc192036f9f4bc6",
		CancelUsername: "테스트",
		CancelMessage: "테스트 관리자",
	}

	fmt.Println("--------------- RequestCashReceiptCancel() Start ---------------")
	res, err := api.RequestCashReceiptCancel(cancelData)
	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- RequestCashReceiptCancel() End ---------------")
}


func RequestAuthentication(api *Api) {
	authData := Authentication{
		Pg: "다날",
		Method: "본인인증",
		Username: "사용자명",
		IdentityNo: "0000000", //생년월일 + 주민번호 뒷 1자리
		Carrier: "SKT",
		Phone:  "01010002000", //사용자 전화번호
		SiteUrl: "https://www.bootpay.co.kr",
		OrderName: "회원 본인인증",
		AuthenticationId: fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
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
		Otp: "613026",
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

func requestSubscribeAutomaticTransferBillingKey(api *Api) {
    fmt.Println("--------------- requestSubscribeAutomaticTransferBillingKey() Start ---------------")
	
	subscriptId := fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond)))

	payload := BillingKeyPayload{
		SubscriptionId: subscriptId,
		Pg: "nicepay",
		OrderName: "정기결제 테스트 아이템",
		Username: "홍길동",
		AuthType: "ARS",
		BankName: "국민",
		BankAccount: "6756123412342472",
		IdentityNo: "901014",
		CashReceiptType: "소득공제",
		CashReceiptIdentityNo: "01012341234",
		Phone: "01012341234",
	}
	res, err := api.requestSubscribeAutomaticTransferBillingKey(payload)

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- requestSubscribeAutomaticTransferBillingKey() End ---------------")
}

func publishAutomaticTransferBillingKey(api *Api) {
    fmt.Println("--------------- publishAutomaticTransferBillingKey() Start ---------------")

	res, err := api.publishAutomaticTransferBillingKey("6655069ca691573f1bb9c28a")

	fmt.Println(res)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("--------------- publishAutomaticTransferBillingKey() End ---------------")
}