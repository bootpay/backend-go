### Unreleased
- `HttpStatus` 응답 struct 필드 (10개 response 타입) 와 `result["http_status"]` map key 를 deprecated 처리.
  - 동작 변경 없음 — 필드/map key 모두 그대로 유지하며 `go vet` / IDE 에서 사용 시 경고만 표시.
  - 성공 여부는 함수 반환 `error` 값 또는 응답 `Status` 필드로 판단 권장.
  - 다음 메이저 버전(v3) 에서 제거 예정. 기존 가맹점 코드는 그대로 컴파일·동작.

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