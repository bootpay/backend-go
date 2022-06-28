## Bootpay Go Server Side Library
부트페이 공식 Go 라이브러리 입니다 (서버사이드 용)
* PG 결제창 연동은 클라이언트 라이브러리에서 수행됩니다. (Javascript, Android, iOS, React Native, Flutter 등)
* 결제 검증 및 취소, 빌링키 발급, 본인인증 등의 수행은 서버사이드에서 진행됩니다. (Java, PHP, Python, Ruby, Node.js, Go, ASP.NET 등)

 
## 기능
1. (부트페이 통신을 위한) 토큰 발급
2. 결제 단건 조회
3. 결제 취소 (전액 취소 / 부분 취소)
4. 신용카드 자동결제 (빌링결제)

   4-1. 빌링키 발급

   4-2. 발급된 빌링키로 결제 승인 요청

   4-3. 발급된 빌링키로 결제 예약 요청

   4-4. 발급된 빌링키로 결제 예약 - 취소 요청

   4-5. 빌링키 삭제

   4-6. 빌링키 조회

5. (생체인증, 비밀번호 결제를 위한) 구매자 토큰 발급
6. 서버 승인 요청
7. 본인 인증 결과 조회
8. (에스크로 이용시) PG사로 배송정보 보내기

## 설치하기

```curl
  go get -u  github.com/bootpay/backend-go/v2
``` 

# 사용하기  
```go 
import (
	bootpay "github.com/bootpay/backend-go/v2"
)

bootpay := bootpay.Api{}.New("5b8f6a4d396fa665fdc2b5ea", "rm6EYECr6aroQVG2ntW0A6LpWnkTgP4uQ3H18sDDUYw=", nil, "")
bootpay.GetToken()
```

## 1. 토큰 발급

부트페이와 서버간 통신을 하기 위해서는 부트페이 서버로부터 토큰을 발급받아야 합니다.  
발급된 토큰은 30분간 유효하며, 최초 발급일로부터 30분이 지날 경우 토큰 발급 함수를 재호출 해주셔야 합니다.
```go  
func GetToken(api *bootpay.Api) {
	fmt.Println("--------------- GetToken() Start ---------------")
	token, err := api.GetToken()
	fmt.Println(token)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- GetToken() End ---------------")
}
```

## 2. 결제 단건 조회
결제창 및 정기결제에서 승인/취소된 결제건에 대하여 올바른 결제건인지 서버간 통신으로 결제검증을 합니다.
```go  
func GetReceipt(api *bootpay.Api) {
	receiptId := "62afc194e38c300021b345d4"
	fmt.Println("--------------- getReceipt() Start ---------------")
	verify, err := api.GetReceipt(receiptId)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}

	fmt.Println(verify)
	fmt.Println("--------------- GetVerify() End ---------------")
}
```

## 3. 결제 취소 (전액 취소 / 부분 취소)
price를 지정하지 않으면 전액취소 됩니다.
* 휴대폰 결제의 경우 이월될 경우 이통사 정책상 취소되지 않습니다
* 정산받으실 금액보다 취소금액이 클 경우 PG사 정책상 취소되지 않을 수 있습니다. 이때 PG사에 문의하시면 되겠습니다.
* 가상계좌의 경우 CMS 특약이 되어있지 않으면 취소되지 않습니다. 그러므로 결제 테스트시에는 가상계좌로 테스트 하지 않길 추천합니다.

부분취는 카드로 결제된 건만 가능하며, 일부 PG사만 지원합니다. 요청시 price에 금액을 지정하시면 되겠습니다.
* (지원가능 PG사: 이니시스, kcp, 다날, 페이레터, 나이스페이, 카카오페이, 페이코)

간혹 개발사에서 실수로 여러번 부분취소를 보내서 여러번 취소되는 경우가 있기때문에, 부트페이에서는 부분취소 중복 요청을 막기 위해 cancel_id 라는 필드를 추가했습니다. cancel_id를 지정하시면, 해당 건에 대해 중복 요청방지가 가능합니다.
```go  

func ReceiptCancel(api *bootpay.Api) {
	payload := bootpay.CancelData{
		ReceiptId:      "62afc3c5cf9f6d001b7d101a",
		CancelId:       fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		CancelUsername: "관리자",
		CancelMessage:  "테스트 결제 취소를 테스트",
	}
	//receiptId := "610cc0cb7b5ba40044b04530"
	//name := "관리자"
	//reason := "테스트 결제 취소를 테스트"
	fmt.Println("--------------- ReceiptCancel() Start ---------------")
	cancel, err := api.ReceiptCancel(payload)

	fmt.Println(cancel)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- ReceiptCancel() End ---------------")
}

```

## 4-1. 빌링키 발급
REST API 방식으로 고객으로부터 카드 정보를 전달하여, PG사에게 빌링키를 발급받을 수 있습니다.
발급받은 빌링키를 저장하고 있다가, 원하는 시점, 원하는 금액에 결제 승인 요청하여 좀 더 자유로운 결제시나리오에 적용이 가능합니다.
* 비인증 정기결제(REST API) 방식을 지원하는 PG사만 사용 가능합니다.
```go  

func GetBillingKey(api *bootpay.Api) {
	fmt.Println("--------------- GetBillingKey() Start ---------------")
	payload := bootpay.BillingKeyPayload{
		SubscriptionId:  fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
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
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- GetBillingKey() End ---------------")
}

```

## 4-2. 발급된 빌링키로 결제 승인 요청
발급된 빌링키로 원하는 시점에 원하는 금액으로 결제 승인 요청을 할 수 있습니다. 잔액이 부족하거나 도난 카드 등의 특별한 건이 아니면 PG사에서 결제를 바로 승인합니다.
```go  
func RequestSubscribe(api *bootpay.Api) {
	payload := bootpay.SubscribePayload{
		BillingKey: "62afc52dcf9f6d001d7d1035",
		OrderName:  "정기결제 테스트",
		OrderId:    fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		Price:      1000,
	}

	fmt.Println("--------------- requestSubscribe() Start ---------------")
	cancel, err := api.RequestSubscribe(payload)

	fmt.Println(cancel)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- requestSubscribe() End ---------------")
}
```
## 4-3. 발급된 빌링키로 결제 예약 요청
원하는 시점에 4-1로 결제 승인 요청을 보내도 되지만, 빌링키 발급 이후에 바로 결제 예약 할 수 있습니다. (빌링키당 최대 10건)
```go  
func ReserveSubscribe(api *bootpay.Api) {
	s10 := time.Now().Add(time.Second * 100).Format("2006-01-02T15:04:05-07:00")
	payload := bootpay.SubscribePayload{
		BillingKey:       "62aff193cf9f6d001a7d10be",
		OrderName:        "정기결제 테스트",
		OrderId:          fmt.Sprintf("%+8d", (time.Now().UnixNano() / int64(time.Millisecond))),
		ReserveExecuteAt: s10,
		Price:            1000,
	}

	fmt.Println("--------------- ReserveSubscribe() Start ---------------")
	cancel, err := api.ReserveSubscribe(payload)

	fmt.Println(cancel)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- ReserveSubscribe() End ---------------")
}
```
## 4-4. 발급된 빌링키로 결제 예약 - 취소 요청
빌링키로 예약된 결제건을 취소합니다.
```go  
func ReserveCancel(api *bootpay.Api) {
	reserveId := "62aff2a0cf9f6d001a7d10c4"
	fmt.Println("--------------- ReserveCancel() Start ---------------")
	cancel, err := api.ReserveCancelSubscribe(reserveId)

	fmt.Println(cancel)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- ReserveCancel() End ---------------")
}
```
## 4-5. 빌링키 삭제
발급된 빌링키로 더 이상 사용되지 않도록, 삭제 요청합니다.
```go  
func DestroyBillingKey(api *bootpay.Api) {
	billingKey := "62afc52dcf9f6d001d7d1035"
	fmt.Println("--------------- DestroyBillingKey() Start ---------------")
	res, err := api.DestroyBillingKey(billingKey)

	fmt.Println(res)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- DestroyBillingKey() End ---------------")
}
```
## 5. 사용자 토큰 발급
(부트페이 단독) 부트페이에서 제공하는 간편결제창, 생체인증 기반의 결제 사용을 위해서는 개발사에서 회원 고유번호를 관리해야하며, 해당 회원에 대한 사용자 토큰을 발급합니다.
이 토큰값을 기반으로 클라이언트에서 결제요청 하시면 되겠습니다.
```go  
func GetUserToken(api *bootpay.Api) {
	payload := bootpay.EasyUserTokenPayload{
		UserId: "user_1234",
		Email:  "test1234@gmail.com",
		Name:   "홍길동",
		Gender: 0,
		Birth:  "19861014",
		Phone:  "01012345678",
	}

	fmt.Println("--------------- GetUserToken() Start ---------------")
	cancel, err := api.GetUserToken(payload)

	fmt.Println(cancel)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- GetUserToken() End ---------------")
}
``` 

## 6. 서버 승인 요청
결제승인 방식은 클라이언트 승인 방식과, 서버 승인 방식으로 총 2가지가 있습니다.

클라이언트 승인 방식은 pythonscript나 native 등에서 confirm 함수에서 진행하는 일반적인 방법입니다만, 경우에 따라 서버 승인 방식이 필요할 수 있습니다.

필요한 이유
1. 100% 안정적인 결제 후 고객 안내를 위해 - 클라이언트에서 PG결제 진행 후 승인 완료될 때 onDone이 수행되지 않아 (인터넷 환경 등), 결제 이후 고객에게 안내하지 못할 수 있습니다
2. 단일 트랜잭션의 개념이 필요할 경우 - 재고파악이 중요한 커머스를 운영할 경우 트랜잭션 개념이 필요할 수 있겠으며, 이를 위해서는 서버 승인을 사용해야 합니다.
```go  
func ServerConfirm(api *bootpay.Api) {
	receiptId := "62afda41cf9f6d001f7d105f"
	fmt.Println("--------------- ServerConfirm() Start ---------------")
	res, err := api.ServerConfirm(receiptId)

	fmt.Println(res)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- ServerConfirm() End ---------------")
}
```

## 7. 본인 인증 결과 조회
다날 본인인증 후 결과값을 조회합니다.
다날 본인인증에서 통신사, 외국인여부, 전화번호 이 3가지 정보는 다날에 추가로 요청하셔야 받으실 수 있습니다.
```go 
func Certificate(api *bootpay.Api) {
	receiptId := "6285ffa6cf9f6d0022c4346b"
	fmt.Println("--------------- Certificate() Start ---------------")
	res, err := api.Certificate(receiptId)

	fmt.Println(res)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- Certificate() End ---------------")
}
```


8. (에스크로 이용시) PG사로 배송정보 보내기
   현금 거래에 한해 구매자의 안전거래를 보장하는 방법으로, 판매자와 구매자의 온라인 전자상거래가 원활하게 이루어질 수 있도록 중계해주는 매매보호서비스입니다. 국내법에 따라 전자상거래에서 반드시 적용이 되어 있어야합니다. PG에서도 에스크로 결제를 지원하며, 에스크로 결제 사용을 원하시면 PG사 가맹시에 에스크로결제를 미리 얘기하고나서 진행을 하시는 것이 수월합니다.

PG사로 배송정보( 이니시스, KCP만 지원 )를 보내서 에스크로 상태를 변경하는 API 입니다.
```go 
func ShoppingStart(api *bootpay.Api) {
	shipping := bootpay.Shipping{
		ReceiptId:      "628ae7ffd01c7e001e9b6066",
		TrackingNumber: "123456",
		DeliveryCorp:   "CJ대한통운",
		User: bootpay.ShippingUser{
			Username: "홍길동",
			Phone:    "01000000000",
			Address:  "서울특별시 종로구",
			Zipcode:  "08490",
		},
	}

	fmt.Println("--------------- ShoppingStart() Start ---------------")
	res, err := api.PutShippingStart(shipping)

	fmt.Println(res)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- ShoppingStart() End ---------------")
}
```

## Example 프로젝트

[적용한 샘플 프로젝트](https://github.com/bootpay/backend-go-example)을 참조해주세요

## Documentation

[부트페이 개발매뉴얼](https://docs.bootpay.co.kr/next/)을 참조해주세요

## 기술문의

[부트페이 홈페이지](https://www.bootpay.co.kr) 우측 하단 채팅을 통해 기술문의 주세요!

## License

[MIT License](https://opensource.org/licenses/MIT).

