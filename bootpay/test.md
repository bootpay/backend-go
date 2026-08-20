# Go SDK 테스트 실행 가이드

## 테스트 파일 구성

Go 는 동일 디렉토리 내 모든 테스트가 같은 package 를 공유하므로 PG/Commerce 를 폴더로 분리하지 않고 파일명 prefix 로 구분한다.

- `pg_*_test.go` — PG API 테스트 (auth, billing, cash receipt, payment, token)
- `commerce_*_test.go` — Commerce API 테스트 (order, product, store, token, user)
- `legacy_compatibility_test.go` — application_id/private_key 레거시 인증 호환성
- `bootpay_test.go`, `config_test.go` — 공통 setup/helper

특정 그룹만 실행하려면 `-run` 으로 prefix 매칭한다:

```bash
go test ./bootpay/... -run "^TestPg" -v       # PG 테스트만
go test ./bootpay/... -run "^TestCommerce" -v # Commerce 테스트만
```

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
// PG API 키 가져오기 (권장)
clientKey, secretKey := GetPgClientKeys()

// Legacy fallback
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

## PG 인증 방식 토글 (BOOTPAY_AUTH_MODE)

PG 테스트는 기본적으로 신규 `client_key/secret_key` 방식으로 동작한다. 매 실행 시 환경변수로 레거시 `application_id/private_key` 방식으로 전환할 수 있다.

### 토글 contract

| `BOOTPAY_AUTH_MODE` | 동작 |
|---|---|
| `new` (기본, 미설정 시 동일) | `NewAPIWithClientKey(client_key, secret_key, mode)` 로 인스턴스 생성. Basic Auth 헤더 자동 부착. |
| `legacy` | `NewAPI(application_id, private_key, mode)` 로 인스턴스 생성. `GetAccessToken()` 호출 후 `Bearer` 헤더 사용. |

키 값은 모두 `.env` (또는 환경변수) 로 주입한다 — `.env.example` 참고.

### 사용법

```bash
# (1) 기본 — env var 생략 (= new)
CGO_ENABLED=0 go test -v -run TestPg

# (2) 한 번만 legacy 로 전환
BOOTPAY_AUTH_MODE=legacy CGO_ENABLED=0 go test -v -run TestPg

# (3) 셸 세션 동안 legacy 고정
export BOOTPAY_AUTH_MODE=legacy
CGO_ENABLED=0 go test -v ./bootpay/...
unset BOOTPAY_AUTH_MODE

# (4) 영구 전환 — .env 의 BOOTPAY_AUTH_MODE 값을 legacy 로 바꾸면 셸 export 없이도 동작
```

### 진입 헬퍼 — 어디서 토글이 흡수되는가

`bootpay/config_test.go` 의 `CreatePgApi()` 가 `BOOTPAY_AUTH_MODE` 값에 따라 `NewAPIWithClientKey(...)` 또는 `NewAPI(...)` 를 반환한다. 모든 PG 테스트는 한 줄로 두 모드를 모두 지원한다:

```go
api := CreatePgApi()
```

### 실행 시 인증 모드 표시

`CreatePgApi()` 호출 시마다 stdout 에 한 줄로 어떤 모드가 활성화됐는지 표시된다. `go test -v` 출력에서 즉시 확인 가능:

```
[BOOTPAY_AUTH_MODE=new] PG: client_key/secret_key (Basic Auth) | env=production
[BOOTPAY_AUTH_MODE=legacy] PG: application_id/private_key (Bearer) | env=production
```

### 토글의 영향을 받지 않는 파일

다음 테스트는 한 함수 안에서 두 모드를 모두 검증하므로 환경변수에 무관하게 동일한 동작을 한다:

- `pg_token_test.go`
- `legacy_compatibility_test.go`
