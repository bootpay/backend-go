# Go SDK 테스트 가이드

## 테스트 실행 방법

> **중요**: 모든 테스트는 `bootpay` 폴더 안에서 실행해야 합니다.

```bash
cd bootpay
```

### macOS 환경 (Go 1.22+)

macOS에서 Go 1.22 이상 버전 사용 시 `CGO_ENABLED=0` 옵션이 필요합니다:

```bash
CGO_ENABLED=0 go test -v
```

### 전체 테스트 실행

```bash
CGO_ENABLED=0 go test -v
```

### 특정 테스트만 실행

```bash
# PG API 테스트
CGO_ENABLED=0 go test -v -run TestFunctions

# Commerce API 테스트
CGO_ENABLED=0 go test -v -run TestCommerce
```

### Example 테스트 실행

```bash
# 모든 Example 함수 실행
CGO_ENABLED=0 go test -v -run Example

# 특정 Example만 실행
CGO_ENABLED=0 go test -v -run ExampleNewCommerceApi
CGO_ENABLED=0 go test -v -run ExampleCommerceApi_User
CGO_ENABLED=0 go test -v -run ExampleCommerceApi_Product
```

---

## PG API 테스트

### 테스트 파일
- `bootpay/bootpay_test.go`

### API 키 설정
`bootpay_test.go` 파일에서 테스트용 API 키를 설정합니다:

```go
func getTestApi() *Api {
    return NewAPI("YOUR_APPLICATION_ID", "YOUR_PRIVATE_KEY", nil, "development")
}
```

### 사용 예시

```go
package main

import (
    "fmt"
    "github.com/bootpay/backend-go/bootpay"
)

func main() {
    // 새로운 방식 (권장)
    api := bootpay.NewAPI("application_id", "private_key", nil, "development")

    // 기존 방식 (하위 호환성)
    // api := bootpay.Api{}.New("application_id", "private_key", nil, "development")

    // 토큰 발급
    token, err := api.GetToken()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Token:", token)

    // 결제 조회
    receipt, _ := api.GetReceipt("receipt_id")
    fmt.Println("Receipt:", receipt)
}
```

---

## Commerce API 테스트

### 테스트 파일
- `bootpay/commerce_example_test.go`

### API 키 설정
`commerce_example_test.go` 파일에서 테스트용 키를 수정합니다:

```go
const (
    testClientKey = "your_client_key"
    testSecretKey = "your_secret_key"
)
```

### 사용 예시

```go
package main

import (
    "fmt"
    "github.com/bootpay/backend-go/bootpay"
)

func main() {
    // 새로운 방식 (권장)
    commerce := bootpay.NewCommerceAPI("client_key", "secret_key", nil, "development")

    // 기존 방식 (하위 호환성)
    // commerce := bootpay.NewCommerceApi("client_key", "secret_key", nil, "development")

    // 토큰 발급
    result, err := commerce.GetAccessToken()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Token:", result)

    // Role 설정 (method chaining)
    commerce.AsManager()

    // User 모듈 사용
    users, _ := commerce.User.List(&bootpay.UserListParams{
        ListParams: bootpay.ListParams{
            Page:  1,
            Limit: 10,
        },
    })
    fmt.Println("Users:", users)
}
```

---

## 모듈별 테스트 예시

### User 모듈

```go
commerce := bootpay.NewCommerceAPI("client_key", "secret_key", nil, "development")
commerce.GetAccessToken()

// 사용자 토큰 발급
commerce.User.Token("user_123")

// 회원가입
commerce.User.Join(bootpay.CommerceUser{
    LoginId: "test@example.com",
    LoginPw: "password123",
    Name:    "Test User",
})

// 사용자 목록 조회
commerce.User.List(&bootpay.UserListParams{
    ListParams: bootpay.ListParams{Page: 1, Limit: 10},
})

// 사용자 상세 조회
commerce.User.Detail("user_id")

// 사용자 수정
commerce.User.Update(bootpay.CommerceUser{
    UserId: "user_id",
    Name:   "Updated Name",
})

// 사용자 삭제
commerce.User.Delete("user_id")
```

### UserGroup 모듈

```go
// 그룹 생성
commerce.UserGroup.Create(bootpay.CommerceUserGroup{
    CompanyName:   "Test Company",
    CorporateType: bootpay.CORPORATE_TYPE_CORPORATE,
})

// 그룹 목록 조회
commerce.UserGroup.List(&bootpay.UserGroupListParams{
    ListParams: bootpay.ListParams{Page: 1, Limit: 10},
})

// 그룹에 사용자 추가
commerce.UserGroup.UserCreate("group_id", "user_id")

// 그룹에서 사용자 제거
commerce.UserGroup.UserDelete("group_id", "user_id")

// 그룹 한도 설정
commerce.UserGroup.Limit(bootpay.UserGroupLimitParams{
    UserGroupId:   "group_id",
    UseLimit:      true,
    PurchaseLimit: 1000000,
})
```

### Product 모듈

```go
// 상품 목록 조회
commerce.Product.List(&bootpay.ProductListParams{
    ListParams: bootpay.ListParams{Page: 1, Limit: 10},
    Type:       1,
})

// 상품 생성 (이미지 없이)
commerce.Product.CreateSimple(bootpay.CommerceProduct{
    Name:         "Test Product",
    DisplayPrice: 10000,
    Type:         1,
})

// 상품 생성 (이미지 포함)
commerce.Product.Create(bootpay.CommerceProduct{
    Name:         "Test Product",
    DisplayPrice: 10000,
}, []string{"/path/to/image.jpg"})

// 상품 상세 조회
commerce.Product.Detail("product_id")

// 상품 상태 변경
commerce.Product.Status(bootpay.ProductStatusParams{
    ProductId: "product_id",
    Status:    1,
})

// 상품 삭제
commerce.Product.Delete("product_id")
```

### Invoice 모듈

```go
// 청구서 목록 조회
commerce.Invoice.List(&bootpay.ListParams{Page: 1, Limit: 10})

// 청구서 생성
commerce.Invoice.Create(bootpay.CommerceInvoice{
    Title:  "Test Invoice",
    Price:  50000,
    UserId: "user_123",
    InvoiceItems: []bootpay.CommerceInvoiceItem{
        {Name: "Item 1", Price: 30000, Qty: 1},
        {Name: "Item 2", Price: 20000, Qty: 1},
    },
})

// 청구서 알림 발송
commerce.Invoice.Notify("invoice_id", []int{
    bootpay.INVOICE_SEND_TYPE_SMS,
    bootpay.INVOICE_SEND_TYPE_EMAIL,
})
```

### Order 모듈

```go
// 주문 목록 조회
commerce.Order.List(&bootpay.OrderListParams{
    ListParams:    bootpay.ListParams{Page: 1, Limit: 10},
    UserId:        "user_123",
    Status:        []int{1, 2, 3},
    PaymentStatus: []int{1},
})

// 주문 상세 조회
commerce.Order.Detail("order_id")

// 월별 주문 조회
commerce.Order.Month("group_id", "2024-01")
```

### OrderCancel 모듈

```go
// 취소 요청 목록 조회
commerce.OrderCancel.List(&bootpay.OrderCancelListParams{
    OrderId: "order_123",
})

// 취소 요청
commerce.OrderCancel.Request(bootpay.OrderCancelParams{
    OrderNumber: "ORDER-2024-001",
    RequestCancelParameters: &bootpay.RequestCancelParameter{
        CancelReason: "Customer request",
        CancelType:   1,
    },
})

// 취소 요청 철회
commerce.OrderCancel.Withdraw("cancel_request_id")

// 취소 승인
commerce.OrderCancel.Approve(bootpay.OrderCancelActionParams{
    OrderCancelRequestHistoryId: "cancel_request_id",
    CancelReason:                "Approved",
})

// 취소 거절
commerce.OrderCancel.Reject(bootpay.OrderCancelActionParams{
    OrderCancelRequestHistoryId: "cancel_request_id",
    CancelReason:                "Cannot cancel",
})
```

### OrderSubscription 모듈

```go
// 구독 목록 조회
commerce.OrderSubscription.List(&bootpay.OrderSubscriptionListParams{
    ListParams: bootpay.ListParams{Page: 1, Limit: 10},
    UserId:     "user_123",
})

// 구독 상세 조회
commerce.OrderSubscription.Detail("subscription_id")

// 구독 일시정지
commerce.OrderSubscription.RequestIng.Pause(bootpay.OrderSubscriptionPauseParams{
    OrderSubscriptionId: "subscription_id",
    Reason:              "Temporary pause",
})

// 구독 재개
commerce.OrderSubscription.RequestIng.Resume(bootpay.OrderSubscriptionResumeParams{
    OrderSubscriptionId: "subscription_id",
})

// 해지 위약금 계산
commerce.OrderSubscription.RequestIng.CalculateTerminationFee("subscription_id", "")

// 구독 해지
commerce.OrderSubscription.RequestIng.Termination(bootpay.OrderSubscriptionTerminationParams{
    OrderSubscriptionId: "subscription_id",
    Reason:              "Customer request",
    TerminationFee:      10000,
})
```

### OrderSubscriptionBill 모듈

```go
// 청구 목록 조회
commerce.OrderSubscriptionBill.List(&bootpay.OrderSubscriptionBillListParams{
    ListParams:          bootpay.ListParams{Page: 1, Limit: 10},
    OrderSubscriptionId: "subscription_id",
    Status:              []int{1, 2},
})

// 청구 상세 조회
commerce.OrderSubscriptionBill.Detail("bill_id")

// 청구 수정
commerce.OrderSubscriptionBill.Update(bootpay.CommerceOrderSubscriptionBill{
    OrderSubscriptionBillId: "bill_id",
    Status:                  2,
})
```

### OrderSubscriptionAdjustment 모듈

```go
// 조정 생성 (할인 등)
commerce.OrderSubscriptionAdjustment.Create("subscription_id", bootpay.CommerceOrderSubscriptionAdjustment{
    Name:     "Discount",
    Price:    -5000,
    Duration: 3,
    Type:     bootpay.SUBSCRIPTION_ADJUSTMENT_TYPE_PERIOD_DISCOUNT,
})

// 조정 수정
commerce.OrderSubscriptionAdjustment.Update(bootpay.OrderSubscriptionAdjustmentUpdateParams{
    OrderSubscriptionId:           "subscription_id",
    OrderSubscriptionAdjustmentId: "adjustment_id",
    Price:                         -10000,
})

// 조정 삭제
commerce.OrderSubscriptionAdjustment.Delete("subscription_id", "adjustment_id")
```

---

## Role 설정

Commerce API는 역할(Role)에 따라 접근 권한이 다릅니다:

```go
commerce := bootpay.NewCommerceAPI("client_key", "secret_key", nil, "development")
commerce.GetAccessToken()

// Method chaining으로 Role 설정
commerce.AsUser().User.List(nil)       // 일반 사용자 권한
commerce.AsManager().Order.List(nil)   // 관리자 권한
commerce.AsPartner().Product.List(nil) // 파트너 권한
commerce.AsVendor().Invoice.List(nil)  // 벤더 권한
commerce.AsSupervisor().UserGroup.List(nil) // 슈퍼바이저 권한

// 현재 Role 확인
fmt.Println(commerce.GetRole())
```

---

## 환경별 설정

### PG API

```go
// Development
api := bootpay.NewAPI(appId, privateKey, nil, "development")

// Test
api := bootpay.NewAPI(appId, privateKey, nil, "test")

// Stage
api := bootpay.NewAPI(appId, privateKey, nil, "stage")

// Production (기본값)
api := bootpay.NewAPI(appId, privateKey, nil, "production")
api := bootpay.NewAPI(appId, privateKey, nil, "")
```

### Commerce API

```go
// Development
commerce := bootpay.NewCommerceAPI(clientKey, secretKey, nil, "development")

// Stage
commerce := bootpay.NewCommerceAPI(clientKey, secretKey, nil, "stage")

// Production (기본값)
commerce := bootpay.NewCommerceAPI(clientKey, secretKey, nil, "production")
commerce := bootpay.NewCommerceAPI(clientKey, secretKey, nil, "")
```
