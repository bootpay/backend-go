# Bootpay Go Server Side Library

부트페이 공식 Go 라이브러리입니다.

[![Go Reference](https://pkg.go.dev/badge/github.com/bootpay/backend-go/v2.svg)](https://pkg.go.dev/github.com/bootpay/backend-go/v2)

## 목차

- [설치하기](#설치하기)
- [PG API 사용하기](#사용하기)
- [API 목록](#api-목록)
  - [1. 토큰 발급](#1-토큰-발급)
  - [2. 결제 단건 조회](#2-결제-단건-조회)
  - [3. 결제 취소](#3-결제-취소-전액-취소--부분-취소)
  - [4. 자동/빌링/정기 결제](#4-자동빌링정기-결제)
    - [4-1. 카드 빌링키 발급](#4-1-카드-빌링키-발급)
    - [4-2. 계좌 빌링키 발급](#4-2-계좌-빌링키-발급)
    - [4-3. 결제 요청하기](#4-3-결제-요청하기)
    - [4-4. 결제 예약하기](#4-4-결제-예약하기)
    - [4-5. 예약 조회하기](#4-5-예약-조회하기)
    - [4-6. 예약 취소하기](#4-6-예약-취소하기)
    - [4-7. 빌링키 삭제하기](#4-7-빌링키-삭제하기)
    - [4-8. 빌링키 조회하기](#4-8-빌링키-조회하기)
  - [5. 회원 토큰 발급요청](#5-회원-토큰-발급요청)
  - [6. 서버 승인 요청](#6-서버-승인-요청)
  - [7. 본인 인증](#7-본인-인증)
    - [7-1. 본인인증 결과 조회](#7-1-본인인증-결과-조회)
    - [7-2. 본인인증 REST API 요청](#7-2-본인인증-rest-api-요청)
    - [7-3. 본인인증 승인](#7-3-본인인증-승인)
    - [7-4. 본인인증 SMS 재전송](#7-4-본인인증-sms-재전송)
  - [8. 에스크로](#8-에스크로-이용시-pg사로-배송정보-보내기)
  - [9. 현금영수증](#9-현금영수증)
    - [9-1. 현금영수증 발행하기](#9-1-현금영수증-발행하기)
    - [9-2. 현금영수증 발행 취소](#9-2-현금영수증-발행-취소)
    - [9-3. 별건 현금영수증 발행](#9-3-별건-현금영수증-발행)
    - [9-4. 별건 현금영수증 발행 취소](#9-4-별건-현금영수증-발행-취소)
- [Commerce API 사용하기](#10-commerce-api)
  - [10-1. Commerce API 초기화](#10-1-commerce-api-초기화)
  - [10-2. 사용자 관리](#10-2-사용자-관리)
  - [10-3. 상품 관리](#10-3-상품-관리)
  - [10-4. 주문 관리](#10-4-주문-관리)
  - [10-5. 정기구독 관리](#10-5-정기구독-관리)
  - [10-6. 청구서 관리](#10-6-청구서-관리)
  - [10-7. 알림톡](#10-7-알림톡)
- [Example 프로젝트](#example-프로젝트)
- [Documentation](#documentation)
- [기술문의](#기술문의)
- [License](#license)

---

## 설치하기

```bash
go get -u github.com/bootpay/backend-go/v2
```

---

#
## 환경변수 설정

예제와 테스트는 각 SDK 루트의 `.env` 파일을 우선 읽습니다. 먼저 `.env.example`을 복사한 뒤 필요한 키만 변경하세요. `.env`는 gitignore 처리되어 커밋되지 않습니다.

```bash
cp .env.example .env
# BOOTPAY_ENV=production 또는 development
```

주요 변수:

```env
BOOTPAY_ENV=production
BOOTPAY_PG_CLIENT_KEY_PROD=...
BOOTPAY_PG_SECRET_KEY_PROD=...
BOOTPAY_PG_CLIENT_KEY_DEV=...
BOOTPAY_PG_SECRET_KEY_DEV=...
BOOTPAY_COMMERCE_CLIENT_KEY_PROD=...
BOOTPAY_COMMERCE_SECRET_KEY_PROD=...
BOOTPAY_COMMERCE_CLIENT_KEY_PROD=...
BOOTPAY_COMMERCE_SECRET_KEY_PROD=...
```

변수가 없으면 SDK 테스트용 기본값(NodeJS 기준 ck/sk)으로 fallback 합니다.

# 사용하기

```go
import (
    bootpay "github.com/bootpay/backend-go/v2"
)

func main() {
    // 권장: client_key/secret_key 방식
    api := bootpay.NewAPIWithClientKey(
        os.Getenv("BOOTPAY_PG_CLIENT_KEY_PROD"),
        os.Getenv("BOOTPAY_PG_SECRET_KEY_PROD"),
        nil,  // http.Client (nil이면 기본값 사용)
        "",   // mode: "", "development", "test", "stage"
    )
    // legacy application_id/private_key 방식도 NewAPI로 계속 지원됩니다.

    // 토큰 발급
    token, err := api.GetToken()
    if err != nil {
        panic(err)
    }
    fmt.Println(token)
}
```

### 환경 설정

| Mode | Base URL |
|------|----------|
| `""` (기본값) | `https://api.bootpay.co.kr/v2` |
| `"development"` | `https://dev-api.bootpay.co.kr/v2` |
| `"test"` | `https://test-api.bootpay.co.kr/v2` |
| `"stage"` | `https://stage-api.bootpay.co.kr/v2` |

---

## API 목록

### 1. 토큰 발급

부트페이와 서버간 통신을 하기 위해서는 부트페이 서버로부터 토큰을 발급받아야 합니다.
발급된 토큰은 **30분간 유효**하며, 30분이 지날 경우 토큰 발급 함수를 재호출 해주셔야 합니다.

```go
func GetToken(api *bootpay.Api) {
    token, err := api.GetToken()
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(token)
}
```

---

### 2. 결제 단건 조회

결제창 및 정기결제에서 승인/취소된 결제건에 대하여 올바른 결제건인지 서버간 통신으로 결제검증을 합니다.

```go
func GetReceipt(api *bootpay.Api) {
    receiptId := "62afc194e38c300021b345d4"
    receipt, err := api.GetReceipt(receiptId)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(receipt)
}
```

---

### 3. 결제 취소 (전액 취소 / 부분 취소)

`CancelPrice`를 지정하지 않으면 전액취소 됩니다.

**주의사항:**
- 휴대폰 결제의 경우 이월될 경우 이통사 정책상 취소되지 않습니다
- 정산받으실 금액보다 취소금액이 클 경우 PG사 정책상 취소되지 않을 수 있습니다
- 가상계좌의 경우 CMS 특약이 되어있지 않으면 취소되지 않습니다

**부분취소 지원 PG사:** 이니시스, KCP, 다날, 페이레터, 나이스페이, 카카오페이, 페이코

```go
func ReceiptCancel(api *bootpay.Api) {
    payload := bootpay.CancelData{
        ReceiptId:      "62afc3c5cf9f6d001b7d101a",
        CancelId:       fmt.Sprintf("%d", time.Now().UnixMilli()), // 중복 취소 방지용
        CancelUsername: "관리자",
        CancelMessage:  "테스트 결제 취소",
        // CancelPrice:  500,  // 부분 취소시 금액 지정
        // CancelTaxFree: 0,   // 부분 취소시 면세 금액
    }

    result, err := api.ReceiptCancel(payload)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

---

### 4. 자동/빌링/정기 결제

#### 4-1. 카드 빌링키 발급

REST API 방식으로 고객으로부터 카드 정보를 전달하여 PG사에게 빌링키를 발급받습니다.

```go
func GetBillingKey(api *bootpay.Api) {
    payload := bootpay.BillingKeyPayload{
        SubscriptionId:  fmt.Sprintf("%d", time.Now().UnixMilli()),
        Pg:              "nicepay",
        OrderName:       "정기결제 테스트 아이템",
        CardNo:          "5570********1074",
        CardPw:          "**",
        CardExpireYear:  "**",
        CardExpireMonth: "**",
        CardIdentityNo:  "",
    }

    billingKey, err := api.GetBillingKey(payload)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(billingKey)
}
```

#### 4-2. 계좌 빌링키 발급

고객의 계좌 정보로 빌링키 발급을 요청합니다. 출금동의 확인 절차까지 진행해야 빌링키가 발급됩니다.

**Step 1: 빌링키 발급 요청**
```go
func RequestSubscribeAutomaticTransferBillingKey(api *bootpay.Api) {
    payload := bootpay.BillingKeyPayload{
        SubscriptionId:        fmt.Sprintf("%d", time.Now().UnixMilli()),
        Pg:                    "nicepay",
        OrderName:             "정기결제 테스트 아이템",
        Username:              "홍길동",
        AuthType:              "ARS",  // "ARS" 또는 "간편인증"
        BankName:              "국민",
        BankAccount:           "6756123412342472",
        IdentityNo:            "901014",
        CashReceiptType:       "소득공제",  // "소득공제" 또는 "지출증빙"
        CashReceiptIdentityNo: "01012341234",
        Phone:                 "01012341234",
    }

    result, err := api.RequestSubscribeAutomaticTransferBillingKey(payload)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

**Step 2: 출금 동의 확인**
```go
func PublishAutomaticTransferBillingKey(api *bootpay.Api) {
    receiptId := "6655069ca691573f1bb9c28a"

    result, err := api.PublishAutomaticTransferBillingKey(receiptId)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)  // 빌링키 발급됨
}
```

#### 4-3. 결제 요청하기

발급된 빌링키로 원하는 시점에 원하는 금액으로 결제 승인 요청을 합니다.

```go
func RequestSubscribe(api *bootpay.Api) {
    payload := bootpay.SubscribePayload{
        BillingKey: "62afc52dcf9f6d001d7d1035",
        OrderName:  "정기결제 테스트",
        OrderId:    fmt.Sprintf("%d", time.Now().UnixMilli()),
        Price:      1000,
        Items: []bootpay.Item{
            {
                Name:  "테스트 상품",
                Qty:   1,
                Id:    "item_1",
                Price: 1000,
            },
        },
    }

    result, err := api.RequestSubscribe(payload)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

#### 4-4. 결제 예약하기

빌링키 발급 이후 원하는 시점에 결제를 예약할 수 있습니다. (빌링키당 최대 10건)

```go
func ReserveSubscribe(api *bootpay.Api) {
    reserveTime := time.Now().Add(time.Minute * 10).Format("2006-01-02T15:04:05-07:00")

    payload := bootpay.SubscribePayload{
        BillingKey:       "62aff193cf9f6d001a7d10be",
        OrderName:        "정기결제 테스트",
        OrderId:          fmt.Sprintf("%d", time.Now().UnixMilli()),
        ReserveExecuteAt: reserveTime,
        Price:            1000,
    }

    result, err := api.ReserveSubscribe(payload)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

#### 4-5. 예약 조회하기

```go
func ReserveSubscribeLookup(api *bootpay.Api) {
    reserveId := "6490149ca575b40024f0b70d"

    result, err := api.ReserveSubscribeLookup(reserveId)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

#### 4-6. 예약 취소하기

```go
func ReserveCancel(api *bootpay.Api) {
    reserveId := "62aff2a0cf9f6d001a7d10c4"

    result, err := api.ReserveCancelSubscribe(reserveId)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

#### 4-7. 빌링키 삭제하기

발급된 빌링키를 삭제합니다. 삭제하더라도 예약된 결제건은 취소되지 않습니다.

```go
func DestroyBillingKey(api *bootpay.Api) {
    billingKey := "62afc52dcf9f6d001d7d1035"

    result, err := api.DestroyBillingKey(billingKey)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

#### 4-8. 빌링키 조회하기

클라이언트에서 빌링키 발급시, 보안상 클라이언트 이벤트에 빌링키를 전달해주지 않습니다.

**receiptId로 조회:**
```go
func LookupBillingKey(api *bootpay.Api) {
    receiptId := "62afccb3cf9f6d001b7d101d"

    result, err := api.LookupBillingKey(receiptId)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

**billingKey로 조회:**
```go
func LookupBillingKeyByKey(api *bootpay.Api) {
    billingKey := "66542dfb4d18d5fc7b43e1b6"

    result, err := api.LookupBillingKeyByKey(billingKey)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

---

### 5. 회원 토큰 발급요청

ㅇㅇ페이 사용을 위해 가맹점 회원의 토큰을 발급합니다.

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

    result, err := api.GetUserToken(payload)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

또는 간단하게:

```go
func RequestUserToken(api *bootpay.Api) {
    payload := bootpay.UserTokenRequest{
        UserId: "gosomi1",
        Phone:  "01012345678",
    }

    result, err := api.RequestUserToken(payload)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

---

### 6. 서버 승인 요청

결제승인 방식은 클라이언트 승인 방식과 서버 승인 방식 2가지가 있습니다.

**서버 승인이 필요한 경우:**
1. 100% 안정적인 결제 후 고객 안내가 필요한 경우
2. 단일 트랜잭션 개념이 필요한 경우 (재고 파악이 중요한 커머스)

```go
func ServerConfirm(api *bootpay.Api) {
    receiptId := "62afda41cf9f6d001f7d105f"

    result, err := api.ServerConfirm(receiptId)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

---

### 7. 본인 인증

#### 7-1. 본인인증 결과 조회

다날 본인인증 후 결과값을 조회합니다.

```go
func Certificate(api *bootpay.Api) {
    receiptId := "6285ffa6cf9f6d0022c4346b"

    result, err := api.Certificate(receiptId)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

#### 7-2. 본인인증 REST API 요청

```go
func RequestAuthentication(api *bootpay.Api) {
    authData := bootpay.Authentication{
        Pg:               "다날",
        Method:           "본인인증",
        Username:         "홍길동",
        IdentityNo:       "9001011",  // 생년월일 + 주민번호 뒷 1자리
        Carrier:          "SKT",      // SKT, KT, LGU+, SKT_MVNO, KT_MVNO, LGU+_MVNO
        Phone:            "01012345678",
        SiteUrl:          "https://www.example.com",
        OrderName:        "회원 본인인증",
        AuthenticationId: fmt.Sprintf("%d", time.Now().UnixMilli()),
        ClientIp:         "127.0.0.1",
        AuthenticateType: "sms",  // "sms" 또는 "pass"
    }

    result, err := api.RequestAuthentication(authData)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

#### 7-3. 본인인증 승인

```go
func ConfirmAuthentication(api *bootpay.Api) {
    authParams := bootpay.AuthenticationParams{
        ReceiptId: "636a020d1fc1920373e6d8ff",
        Otp:       "613026",
    }

    result, err := api.ConfirmAuthentication(authParams)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

#### 7-4. 본인인증 SMS 재전송

```go
func RealarmAuthentication(api *bootpay.Api) {
    authParams := bootpay.AuthenticationParams{
        ReceiptId: "636a020d1fc1920373e6d8ff",
    }

    result, err := api.RealarmAuthentication(authParams)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

---

### 8. 에스크로 이용시 PG사로 배송정보 보내기

현금 거래에 한해 구매자의 안전거래를 보장하는 매매보호서비스입니다.

**지원 PG사:** 이니시스, KCP

```go
func ShippingStart(api *bootpay.Api) {
    shipping := bootpay.Shipping{
        ReceiptId:      "628ae7ffd01c7e001e9b6066",
        ReceiptUrl:     "https://example.com/receipt",
        TrackingNumber: "123456789",
        DeliveryCorp:   "CJ대한통운",
        ShippingPrepayment: true,  // 선불 여부
        ShippingDay:    3,         // 배송일
        User: bootpay.ShippingUser{
            Username: "홍길동",
            Phone:    "01012345678",
            Address:  "서울특별시 종로구",
            Zipcode:  "03000",
        },
        Company: bootpay.ShippingCompany{
            Name:    "부트페이",
            Phone:   "02-1234-5678",
            Zipcode: "06000",
            Addr1:   "서울특별시 강남구",
            Addr2:   "테헤란로 123",
        },
    }

    result, err := api.PutShippingStart(shipping)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

---

### 9. 현금영수증

#### 9-1. 현금영수증 발행하기

부트페이 API를 통해 결제된 건에 대하여 현금영수증을 발행합니다.

```go
func RequestCashReceiptByBootpay(api *bootpay.Api) {
    cashReceipt := bootpay.CashReceiptData{
        ReceiptId:       "62e0f11f1fc192036b1b3c92",
        Username:        "홍길동",
        Email:           "test@bootpay.co.kr",
        Phone:           "01012345678",
        IdentityNo:      "01012345678",
        CashReceiptType: "소득공제",  // "소득공제" 또는 "지출증빙"
    }

    result, err := api.RequestCashReceiptByBootpay(cashReceipt)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

#### 9-2. 현금영수증 발행 취소

```go
func RequestCashReceiptCancelByBootpay(api *bootpay.Api) {
    cancelData := bootpay.CancelData{
        ReceiptId:      "62e0f11f1fc192036b1b3c92",
        CancelUsername: "관리자",
        CancelMessage:  "고객 요청으로 취소",
    }

    result, err := api.RequestCashReceiptCancelByBootpay(cancelData)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

#### 9-3. 별건 현금영수증 발행

부트페이 결제와 상관없이 현금영수증을 발행합니다.

```go
func RequestCashReceipt(api *bootpay.Api) {
    purchasedAt := time.Now().Format("2006-01-02T15:04:05-07:00")

    cashReceipt := bootpay.CashReceiptData{
        Pg:              "토스",
        Price:           1000,
        OrderName:       "테스트 상품",
        CashReceiptType: "소득공제",
        IdentityNo:      "01012345678",
        OrderId:         fmt.Sprintf("%d", time.Now().UnixMilli()),
        PurchasedAt:     purchasedAt,
        Currency:        "KRW",
    }

    result, err := api.RequestCashReceipt(cashReceipt)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

#### 9-4. 별건 현금영수증 발행 취소

```go
func RequestCashReceiptCancel(api *bootpay.Api) {
    cancelData := bootpay.CancelData{
        ReceiptId:      "62f4be7f1fc192036f9f4bc6",
        CancelUsername: "관리자",
        CancelMessage:  "고객 요청으로 취소",
    }

    result, err := api.RequestCashReceiptCancel(cancelData)
    if err != nil {
        fmt.Println("error:", err.Error())
        return
    }
    fmt.Println(result)
}
```

---

## 응답 타입

Go SDK는 상세 응답 타입을 제공합니다. `types.go` 파일에서 확인할 수 있습니다:

- `AccessTokenResponse` - 토큰 발급 응답
- `ReceiptResponse` - 결제 정보 응답
- `CertificateResponse` - 본인인증 응답
- `SubscriptionBillingResponse` - 빌링키 응답
- `UserTokenResponse` - 사용자 토큰 응답
- 등...

---

## 10. Commerce API

부트페이 Commerce API를 사용하여 사용자, 상품, 주문, 정기구독 등을 관리할 수 있습니다.

### 10-1. Commerce API 초기화

```go
import (
    commerce "github.com/bootpay/backend-go/v2/commerce"
)

func main() {
    api := commerce.ApiCommerce{}.New(
        os.Getenv("BOOTPAY_COMMERCE_CLIENT_KEY_PROD"),
        os.Getenv("BOOTPAY_COMMERCE_SECRET_KEY_PROD"),
        nil,  // http.Client (nil이면 기본값 사용)
        "production", // mode: "production", "development", "stage"
    )

    // 토큰 발급
    token, err := api.GetToken()
    if err != nil {
        panic(err)
    }
}
```

### 10-2. 사용자 관리

```go
// 사용자 목록 조회
users, err := api.UserList(map[string]interface{}{
    "page":  1,
    "limit": 10,
})

// 사용자 상세 조회
user, err := api.UserDetail("USER_ID")

// 회원가입
newUser, err := api.UserJoin(map[string]interface{}{
    "login_id": "test@example.com",
    "login_pw": "password123",
    "name":     "홍길동",
    "email":    "test@example.com",
    "phone":    "010-1234-5678",
})

// 사용자 정보 수정
updatedUser, err := api.UserUpdate(map[string]interface{}{
    "user_id": "USER_ID",
    "name":    "수정된 이름",
})
```

### 10-3. 상품 관리

```go
// 상품 목록 조회
products, err := api.ProductList(map[string]interface{}{
    "page":  1,
    "limit": 10,
})

// 상품 생성
product, err := api.ProductCreate(map[string]interface{}{
    "name":        "테스트 상품",
    "price":       10000,
    "description": "상품 설명",
})

// 상품 상세 조회
productDetail, err := api.ProductDetail("PRODUCT_ID")

// 상품 수정
updatedProduct, err := api.ProductUpdate(map[string]interface{}{
    "product_id": "PRODUCT_ID",
    "name":       "수정된 상품명",
    "price":      15000,
})
```

### 10-4. 주문 관리

```go
// 주문 목록 조회
orders, err := api.OrderList(map[string]interface{}{
    "page":  1,
    "limit": 10,
})

// 주문 상세 조회
order, err := api.OrderDetail("ORDER_ID")

// 월별 주문 조회
monthOrders, err := api.OrderMonth("USER_GROUP_ID", "2024-12")
```

### 10-5. 정기구독 관리

```go
// 정기구독 목록 조회
subscriptions, err := api.OrderSubscriptionList(nil)

// 정기구독 상세 조회
subscription, err := api.OrderSubscriptionDetail("ORDER_SUBSCRIPTION_ID")

// 정기구독 일시정지
err = api.OrderSubscriptionPause(map[string]interface{}{
    "order_subscription_id": "ORDER_SUBSCRIPTION_ID",
    "pause_days":            30,
    "reason":                "일시정지 사유",
})

// 정기구독 재개
err = api.OrderSubscriptionResume(map[string]interface{}{
    "order_subscription_id": "ORDER_SUBSCRIPTION_ID",
})

// 정기구독 해지
err = api.OrderSubscriptionTermination(map[string]interface{}{
    "order_subscription_id": "ORDER_SUBSCRIPTION_ID",
    "reason":                "해지 사유",
})
```

### 10-6. 청구서 관리

```go
// 청구서 목록 조회
invoices, err := api.InvoiceList(nil)

// 청구서 생성
invoice, err := api.InvoiceCreate(map[string]interface{}{
    "user_id": "USER_ID",
    "amount":  50000,
    "title":   "청구서 제목",
})

// 청구서 알림 전송
err = api.InvoiceNotify("INVOICE_ID", []int{1, 2}) // 1: SMS, 2: Email
```

### 10-7. 알림톡

Commerce API 인스턴스(`NewCommerceAPI`)의 알림톡 모듈로 발신프로필 등록부터 템플릿 검수, 발송, 결과 조회까지 처리합니다.

> ⚠️ 알림톡은 **샌드박스가 없습니다.** 발송·OTP·발신프로필 등록·템플릿 등록은 모두 실제로 나가고 과금됩니다.
> 알림톡 API 는 `Idempotency-Key` 헤더를 읽지 않습니다 — 멱등은 발송의 `ref_id` 로만 성립합니다.

```go
api := bootpay.NewCommerceAPI(clientKey, secretKey, nil, "production")

// 발신프로필(카카오채널) 등록 — 카테고리 조회 → OTP → 등록
categories, err := api.AlimtalkSender.Categories()
_, err = api.AlimtalkSender.Otp(bootpay.AlimtalkSenderOtpParams{
    YellowId: "@부트페이", Phone: "01012345678",
})
sender, err := api.AlimtalkSender.Create(bootpay.AlimtalkSenderCreateParams{
    Otp: "123456", YellowId: "@부트페이", Phone: "01012345678", CategoryCode: "001001",
})

// 공식 템플릿 (부트페이가 미리 승인받아 둔 카탈로그 — 검수 없이 즉시 발송 가능)
official, err := api.AlimtalkOfficial.List(&bootpay.AlimtalkOfficialListParams{Keyword: "주문"})

// 자체 템플릿 — 초안 생성 후 확인하고 등록/검수하는 것을 권장
draft, err := api.AlimtalkTemplate.Create(bootpay.AlimtalkTemplateCreateParams{
    KspId:    "KSP_ID",
    Name:     "주문완료 안내",
    Content:  "#{user_name}님, 주문이 완료되었습니다.",
    Register: bootpay.BoolPtr(false), // ⚠️ 생략하면 생성 즉시 카카오에 실제 등록됩니다
})
_, err = api.AlimtalkTemplate.Register("TEMPLATE_ID")
_, err = api.AlimtalkTemplate.Inspect("TEMPLATE_ID") // ⚠️ 검수 요청은 취소할 수 없습니다

// 발송 — ref_id 가 멱등 키입니다
receipt, err := api.AlimtalkSend.Send(bootpay.AlimtalkSendParams{
    TemplateCode: "TEMPLATE_CODE",
    To:           "01012345678",
    Variables:    map[string]interface{}{"user_name": "홍길동"},
    RefId:        "order-20260827-0001",
    Fallback:     bootpay.BoolPtr(true), // 미지정(nil)이면 프로젝트 기본값을 따릅니다
})

// 벌크 발송 (1요청 = N수신자)
_, err = api.AlimtalkSend.Bulk(bootpay.AlimtalkSendBulkParams{
    TemplateCode: "TEMPLATE_CODE",
    Recipients: []bootpay.AlimtalkSendRecipient{
        {To: "01011112222", RefId: "bulk-0001", Variables: map[string]interface{}{"user_name": "홍길동"}},
    },
})

// 발송 결과 · 집계
messages, err := api.AlimtalkMessage.List(&bootpay.AlimtalkMessageListParams{Status: "success"})
stats, err := api.AlimtalkMessage.Stats("2026-08-01", "2026-08-27")
detail, err := api.AlimtalkMessage.Detail("RECEIPT_ID")

// 수신거부 (발송 판정과 같은 기준 — 전역 차단은 조회만 되고 해제되지 않습니다)
_, err = api.AlimtalkOptout.Check(bootpay.AlimtalkOptoutCheckParams{Phones: []string{"01012345678"}})
_, err = api.AlimtalkOptout.Create(bootpay.AlimtalkOptoutCreateParams{Phone: "01012345678"})

// 알림톡 전용 웹훅 (⚠️ 주문·구독 통합 웹훅과 별개입니다)
_, err = api.AlimtalkWebhook.Update(bootpay.AlimtalkWebhookUpdateParams{
    Url:    "https://example.com/alimtalk-hook", // https 만 허용
    Events: []int{301, 302, 310, 311},
})
```

---

더 자세한 Commerce API 사용 예제는 [commerce_test](./commerce_test) 디렉토리를 참고해주세요.

---

## Example 프로젝트

[적용 샘플 프로젝트](https://github.com/bootpay/backend-go-example)를 참조해주세요.

---

## Documentation

[부트페이 개발매뉴얼](https://developer.bootpay.co.kr/)을 참조해주세요.

---

## 기술문의

[부트페이 홈페이지](https://www.bootpay.co.kr) 우측 하단 채팅을 통해 기술문의 주세요!

---

## License

[MIT License](https://opensource.org/licenses/MIT)
