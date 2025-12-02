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