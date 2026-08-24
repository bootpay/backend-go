### 2.5.0

#### Commerce scope(BOOTPAY-ROLE) 정합성 (동작 변경)

서버(commerce-api)가 `scope_invalid!` 로 supervisor / manager scope 를 요구하는 10개 엔드포인트가 `BOOTPAY-ROLE: user` 로 나가고 있었다. 요청 단위로 올바른 scope 를 붙인다. Java SDK 3.3.0 · Ruby SDK 와 같은 규약이다.

- `OrderSubscription` — `SupervisorApprove` / `SupervisorReject` / `SupervisorTerminate` / `SupervisorPause` / `SupervisorResume` → **supervisor**
- `Category` — `Create` / `Update` / `Destroy` → **supervisor**
- `UserGroup` — `UserCreate` / `UserDelete` → **manager**

부수 효과로 이 10개 호출에 `Idempotency-Key` 가 자동 부착된다 (다른 supervisor 메서드·Ruby SDK 와 동일). 요청 경로·바디는 변경 없다.
⚠️ 그동안 이 API 들은 올바른 키로도 scope 오류로 거절됐다. 우회하려고 role 을 직접 조작하던 코드가 있다면 제거해도 된다.

- 파라미터 struct 에 `IdempotencyKey` 필드(`json:"-"`)를 추가했다 (순수 추가 — 기존 필드·JSON 바디 불변). `Category.Destroy(categoryId, idempotencyKey ...string)` / `UserGroup.UserCreate(userGroupId, userId, idempotencyKey ...string)` / `UserGroup.UserDelete(...)` 는 가변 인자로 받아 기존 호출부가 그대로 컴파일된다.

#### Commerce 최상위 배열 응답 파싱 (버그 수정)

- 일부 엔드포인트(`GET categories`, `GET coupon`, `GET coupon/available`)는 객체가 아니라 **최상위 JSON 배열**을 내려준다. 이전에는 `map[string]interface{}` 로만 디코드를 시도하고 그 에러를 버려서 **호출자가 빈 map 을 error=nil 로 받았다.**
- 배열 응답은 `{"data": [...]}` 로 감싸 돌려준다 (Java SDK `BootpayStoreObject` 폴백과 동일). 객체 응답의 모양은 그대로다.
- 빈 본문은 빈 map, JSON 이 아닌 본문(HTML 5xx 등)은 이제 **에러**를 낸다 — 조용한 빈 map 이 성공으로 보이던 문제를 없앤다.

#### PG

- `Api.GetAccessToken()` 추가 — `Api.GetToken()` 의 별칭. 같은 패키지에서 `Api.GetToken()` 은 토큰을 *발급*하고 `CommerceApi.GetToken()` 은 저장된 값을 *읽어서* 이름과 의미가 반대였다. 기존 `GetToken()` 은 그대로 유지된다.

### 2.4.0
- NodeJS SDK 2.9.0 parity (2.7.0 ~ 2.9.0 변경분 반영). 기존 함수 시그니처·struct 필드는 모두 유지 — 신규 기능은 신규 메서드/optional 필드로만 추가.
- PG: 우선순위(순차) 결제 빌링키 조회 `LookupSequentialBillingKey(widgetKey, billingKey, userId)` 추가 — `GET subscribe/sequential_billing_key/{billing_key}?widget_key=&user_id=` (QueryEscape 적용).
- Commerce: 공통 요청 계층에 요청별 헤더 지원 신설.
  - 요청별로 지정한 `BOOTPAY-ROLE` 을 공통 계층이 덮어쓰지 않는다 (supervisor 전용 endpoint 대응).
  - `Idempotency-Key` 헤더 자동 생성(UUID v4, 외부 의존성 없이 crypto/rand) — 파라미터/인자로 직접 지정 가능.
  - DELETE body 전송 지원 (charge 해지·조정항목 삭제), multipart POST 시 form-data boundary 를 보존.
- Commerce: 신규 모듈
  - `MallSetting` — `GetMallSetting`/`Detail` (`GET mall-setting`), `UpdateMallSetting`/`Update` (`PUT mall-setting`, flatten 바디에 지정된 값만 전송 — `*bool`/`*int` 로 명시적 false/0 전송 가능). supervisor 전용.
  - `Webhook.SendTest` — `POST webhook/test` (`header_content_type` 파라미터, 테스트 웹훅 발송).
- Commerce: 수시결제(온디맨드) charge_key (supervisor 전용)
  - `OrderSubscription.SupervisorCharge` — `POST order_subscriptions/charge`. charge_key 는 body 로만 전송 (URL/query 금지).
  - `OrderSubscription.SupervisorChargeRevoke` — `DELETE order_subscriptions/charge` (body 전송). 해지 후 해당 키로 재결제 불가.
- Commerce: V1 Mall API
  - User — `UserLogin`, `UserSession`, `UserLogout`, `UserJoin`, `UserJoinCheck`, `UidExist`. 복수형 `users/...` 경로 사용 (v1 에 단수 라우트 없음). 회원 JWT 는 `Bootpay-User-JWT` 헤더로 값이 있을 때만 부착. `MALL_USER_JOIN_CHECK_*` 상수 추가.
  - Product — `Products` (page/limit 기본 1/20, category_id/sort/user_jwt), `ProductDetail(productId, userJwt, idempotencyKey)`.
- Commerce: 정정·확장
  - OrderSubscription.RequestIng 에 `Purchase`(`POST .../requests/ing/purchase`) / `Transfer`(`POST .../requests/ing/transfer`) 추가. Resume 파라미터에 `Reason` 추가.
  - Product.Create — 이미지가 있으면 multipart/form-data (`images[0]`, `images[1]` … 인덱싱), 없으면 JSON. 기존 구현의 이미지 없음 시 이중 요청 버그 수정.
  - Invoice — `ListWithParams(*InvoiceListParams)` 추가 (limit 기본 24, cs_type/user_id/product_type/css_at/cse_at). 기존 `List` 도 page/limit 기본값을 전송하도록 정렬. 응답 data 는 `{ list, count }` 구조 (`{ items, total }` 아님 — Go 는 map 반환이라 메서드 주석으로 문서화). `Notify` 의 sendTypes 가 nil 이면 `send_types` 키 미전송.
  - OrderCancel — `OrderCancelActionParams` 에 정식 필드 `OrderCancellationRequestId` 추가 (구 `OrderCancelRequestHistoryId` 도 계속 동작, 둘 다 세팅 시 신 필드 우선). 서버가 읽는 `Message` 필드 추가.
  - OrderSubscriptionAdjustment — `Delete` 는 대상 ID 를 query 가 아니라 **body** 로 전송. `Update` 에 `Adjustments` 배열 지원 (서버는 duration 회차 단위로 교체, duration 기본 1). `Create` 는 price/duration/tax_free_price 기본값 0/1/0 명시 전송.
  - 파라미터 확장 — UserGroupLimitParams 에 `LimitMonthPurchase`/`LimitWeekPurchase` (*int, 서버 정식 인자명 — Update 로는 한도 미반영); OrderListParams 에 `SearchDateFrom`/`SearchDateTo` (CssAt/CseAt 는 서버 별칭으로 유지); OrderSubscriptionListParams 에 `SearchDateFrom`/`SearchDateTo`/`Status`(*int); OrderSubscriptionRequestListParams 에 `OrderSubscriptionId`/`UserId`/`UserGroupId`; OrderSubscriptionUpdateParams 에 계약변경 필드 일괄 추가 (ProductId/OrderName/Quantity/UseFreeTrial 등).
- Commerce: endpoint 별 `BOOTPAY-ROLE` scope 명시
  - manager — 상품 쓰기(Create/CreateSimple/Update/Status/Delete), UserGroup Limit/AggregateTransaction.
  - supervisor — 구독 계약변경(OrderSubscription.Update), 조정항목(Create/Update/Delete), 요청 승인(OrderSubscriptionRequest.Update), charge/chargeRevoke, MallSetting, OrderCancel Approve/Reject.
  - user — requests/ing 계열, invoice, orderCancel List/Withdraw, uid-exist.
  - OrderSubscriptionRequest List/Detail 은 project_id 유무로 supervisor/user 분기.
  - Store 조회에 `Idempotency-Key` 자동 부착 (`InfoWithIdempotencyKey`/`DetailWithIdempotencyKey` 로 직접 지정 가능).
- Commerce: OrderSubscriptionBill.List parity — page/limit 기본 1/20 상시 전송, `BOOTPAY-ROLE: user` + `Idempotency-Key` 자동 부착. `OrderSubscriptionBillListParams` 에 `IdempotencyKey` 필드 추가.
- Commerce: `ProductStatusParams` 에 `UseSalePeriod`/`SaleStartAt`/`SaleEndAt` 추가.
- Commerce: OrderSubscriptionRequest Update 의 Extra 유실 수정 — 기존에는 `Extra` 가 body 에서 조용히 누락됐다. `OrderSubscriptionRequestUpdateBody` 에 `MarshalJSON` 을 구현해 Extra 를 body 에 병합하고, 해지 승인용 타입 필드 `Price`/`TaxFreePrice`/`TerminationFee`/`LastBillRefundPrice`/`FinalFee`/`ServiceEndAt` (*int — 명시적 0 전송 가능) 를 Params/Body 에 추가 (같은 키 충돌 시 타입 필드 우선, nil Extra 값 미전송).
- 테스트: 기존 라이브 통합 테스트 전체에 `BOOTPAY_ENV=development` 게이트 추가 — 기본(production) 환경에서 `go test ./...` 가 실서버를 호출하던 문제 수정 (네트워크 없는 테스트 2종은 게이트 제외). Example 함수들은 `Output:` 지시자를 비활성화해 go test 실행에서 제외 (라이브 예제 문서용으로 유지).
- 테스트: wire-format(URL·헤더·바디) mock 검증 테스트 추가 — 신규 기능 (`commerce_parity_mock_test.go`, `pg_sequential_billing_mock_test.go`) 및 변경 기능 회귀 (`commerce_parity_changed_mock_test.go`, httptest 실서버 왕복 검증 포함: bill.List·request.Detail project_id 분기·mallSetting alias·Extra 병합·sale period 필드). 신규 endpoint 라이브 테스트는 `BOOTPAY_ENV=development` 전용 게이트 (`commerce_parity_live_test.go`).

### 2.3.0
- 모듈 레이아웃 복원 — `bootpay/` 서브디렉토리에 있던 `go.mod` 와 모든 `.go` 파일을 repo root 로 이동. v2.0.9 시점 레이아웃 복원. 이전 v2.1.0 ~ v2.2.0 은 module path 와 디렉토리 구조 불일치로 Go module proxy 에 publish 되지 않음 (실제 사용자는 v2.0.9 에 묶여 있었음). v2.3.0 부터 다시 정상 publish.
- 인증: client_key/secret_key Basic Auth 지원 (PG + Commerce 공통).
  - 기존 application_id/private_key Bearer 방식 하위 호환 유지.
  - PG: `NewAPIWithClientKey(clientKey, secretKey, client, mode)` 팩토리 추가.
    - ck/sk 모드에서는 매 요청 자동 Basic Auth 헤더 부착 — `GetToken()` 은 합성 응답을 반환하며 `request/token` 호출이 발생하지 않음.
    - `RestConfig` 에 `ClientKey` / `SecretKey` 필드 추가 (`omitempty`).
  - Commerce: `CommerceApi` 의 모든 호출이 ck/sk 로 Basic Auth 사용. ck/sk 누락 시 Authorization 헤더 미부착.
- `WalletRequest`, `RequestWalletPayment`, `WalletPaymentResponse`, `WalletDataResp`, `GetUserWallets` `Deprecated` doc 추가 — 다음 메이저 버전에서 제거 예정.
- `HttpStatus` 응답 struct 필드 (10개 response 타입) 와 `result["http_status"]` map key 를 deprecated 처리.
  - 동작 변경 없음 — 필드/map key 모두 그대로 유지하며 `go vet` / IDE 에서 사용 시 경고만 표시.
  - 성공 여부는 함수 반환 `error` 값 또는 응답 `Status` 필드로 판단 권장.
  - 다음 메이저 버전(v3) 에서 제거 예정. 기존 가맹점 코드는 그대로 컴파일·동작.
- `GetToken` 응답 디코딩 견고성 — JSON 파싱 실패는 에러로 반환, ck/sk 모드는 빈 access_token + http_status:200 합성 응답.
- 테스트 인프라: `.env` / `BOOTPAY_AUTH_MODE=new|legacy` / `BOOTPAY_ENV` 토글로 ck/sk · legacy 양쪽 검증. `pg_token_test.go` 의 ck/sk 테스트는 AUTH_MODE 와 무관하게 ck/sk 인스턴스 직접 생성.

### 2.2.0
- Commerce API 추가
  - User 모듈 (Token, Join, CheckExist, Login, List, Detail, Update, Delete, AuthenticationData)
  - UserGroup 모듈 (Create, List, Detail, Update, UserCreate, UserDelete, Limit, AggregateTransaction)
  - Product 모듈 (List, Create, CreateSimple, Detail, Update, Status, Delete)
  - Invoice 모듈 (List, Create, Detail, Notify)
  - Order 모듈 (List, Detail, Month)
  - OrderCancel 모듈 (List, Request, Withdraw, Approve, Reject)
  - OrderSubscription 모듈 (List, Detail, Update, RequestIng.Pause/Resume/CalculateTerminationFee/Termination)
  - OrderSubscriptionBill 모듈 (List, Detail, Update)
  - OrderSubscriptionAdjustment 모듈 (Create, Update, Delete)
- NewAPI() 함수 추가 (Go 표준 팩토리 패턴)
- NewCommerceAPI() 함수 추가
- 기존 Api{}.New() 및 NewCommerceApi() 하위 호환성 유지

### 2.1.5
- 배송등록 api 필드 추가 

### 2.1.4
- 예약 조회 API 추가 

### 2.1.3
- 계좌 자동이체 API 추가 

### 2.0.9
- 본인인증 REST API 추가 

### 2.0.8
- 현금영수증 API 추가 
- return 타입에 http_status 추가 
- getToken Api 호출시 http_status_code -> http_status 로 키 변경 

### 2.0.7
-  putShippingStart -> PutShippingStart renamed function 

### 2.0.6
-  go.sum added

### 2.0.5
-  package name v2 added

### 2.0.4
-  shipping api http method change get -> put

### 2.0.3
-  shipping model field added

### 2.0.2
-  republish

### 2.0.1
-  escrow api added

### 2.0.0 (2-x-development)
-  bootpay api v1 -> v2 upgrade

### 1.0.7
- verify response date update 

### 1.0.6
- reserve api schedyleType 값이 비어있을때 보완처리  

### 1.0.5
- rename interface

### 1.0.4
- rename interface

### 1.0.3
- package name update

### 1.0.2
- response data type update

### 1.0.1
- response data update

### 1.0.0
- first release  