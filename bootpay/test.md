# Go SDK 테스트 실행 가이드

## 환경 설정

`bootpay/config_test.go` 파일에서 환경을 설정합니다:

```go
// "production" 또는 "development"로 설정
const CurrentEnv = "production"
```

## 테스트 실행

### 전체 테스트 실행
```bash
cd /Users/taesupyoon/bootpay/server/sdk/go
go test ./bootpay/... -v
```

### 개별 테스트 활성화

`bootpay/bootpay_test.go`의 `TestFunctions` 함수에서 원하는 테스트의 주석을 해제합니다:

```go
func TestFunctions(t *testing.T) {
    bootpay := CreatePgApi()

    // 토큰 발급 (필수)
    GetToken(bootpay)

    // 아래 테스트들은 필요에 따라 주석 해제하여 실행
    // ReceiptCancel(bootpay)           // 결제 취소
    // GetReceipt(bootpay)              // 결제 조회
    // GetBillingKey(bootpay)           // 빌링키 발급
    // RequestSubscribe(bootpay)        // 정기결제 실행
    // LookupBillingKey(bootpay)        // 빌링키 조회 (receipt_id)
    // LookupBillingKeyByKey(bootpay)   // 빌링키 조회 (billing_key)
    // ReserveSubscribe(bootpay)        // 예약 결제
    // ReserveSubscribeLookup(bootpay)  // 예약 결제 조회
    // ReserveCancel(bootpay)           // 예약 결제 취소
    // DestroyBillingKey(bootpay)       // 빌링키 삭제
    // GetUserToken(bootpay)            // 사용자 토큰 발급
    // ServerConfirm(bootpay)           // 결제 승인
    // Certificate(bootpay)             // 본인인증 조회
    // ShippingStart(bootpay)           // 에스크로 배송시작
    // RequestCashReceiptByBootpay(bootpay)       // 결제건 현금영수증 발행
    // RequestCashReceiptCancelByBootpay(bootpay) // 결제건 현금영수증 취소
    // RequestCashReceipt(bootpay)                // 현금영수증 발행
    // RequestCashReceiptCancel(bootpay)          // 현금영수증 취소
}
```

## 테스트 데이터

`bootpay/config_test.go`에서 테스트 데이터를 관리합니다:

```go
const (
    TestReceiptId          = "628b2206d01c7e00209b6087"
    TestReceiptIdConfirm   = "62876963d01c7e00209b6028"
    TestReceiptIdCash      = "62e0f11f1fc192036b1b3c92"
    TestReceiptIdEscrow    = "628ae7ffd01c7e001e9b6066"
    TestReceiptIdBilling   = "62c7ccebcf9f6d001b3adcd4"
    TestReceiptIdTransfer  = "66541bc4ca4517e69343e24c"
    TestBillingKey         = "628b2644d01c7e00209b6092"
    TestBillingKey2        = "66542dfb4d18d5fc7b43e1b6"
    TestReserveId          = "6490149ca575b40024f0b70d"
    TestReserveId2         = "628b316cd01c7e00219b6081"
    TestUserId             = "1234"
    TestCertificateReceiptId = "61b009aaec81b4057e7f6ecd"
)
```

## 헬퍼 함수

`config_test.go`에서 제공하는 헬퍼 함수들:

```go
// PG API 키 가져오기
appId, privateKey := GetPgKeys()

// Commerce API 키 가져오기
clientKey, secretKey := GetCommerceKeys()

// PG API 인스턴스 생성 (환경 설정 자동 반영)
bootpay := CreatePgApi()

// Commerce API 인스턴스 생성 (환경 설정 자동 반영)
commerce := CreateCommerceApi()
```

## 폴더 구조

```
bootpay/
├── config_test.go    # 환경 설정 및 테스트 데이터
├── bootpay_test.go   # PG API 테스트
├── commerce_test.go  # Commerce API 테스트
└── test.md           # 테스트 가이드
```
