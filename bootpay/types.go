package bootpay

import "time"

// =====================================================
// Response Types - 응답 타입 정의
// =====================================================

// AccessTokenResponse represents the access token response
type AccessTokenResponse struct {
	ExpireIn    int    `json:"expire_in"`
	AccessToken string `json:"access_token"`
	HttpStatus  int    `json:"http_status"`
}

// ReceiptResponse represents the receipt/payment response
type ReceiptResponse struct {
	ReceiptId          string                 `json:"receipt_id"`
	OrderId            string                 `json:"order_id"`
	Price              float64                `json:"price"`
	TaxFree            float64                `json:"tax_free"`
	CancelledPrice     float64                `json:"cancelled_price"`
	CancelledTaxFree   float64                `json:"cancelled_tax_free"`
	OrderName          string                 `json:"order_name"`
	CompanyName        string                 `json:"company_name"`
	GatewayUrl         string                 `json:"gateway_url"`
	Metadata           string                 `json:"metadata"`
	Sandbox            bool                   `json:"sandbox"`
	Pg                 string                 `json:"pg"`
	Method             string                 `json:"method"`
	MethodSymbol       string                 `json:"method_symbol"`
	MethodOrigin       string                 `json:"method_origin"`
	MethodOriginSymbol string                 `json:"method_origin_symbol"`
	PurchasedAt        *time.Time             `json:"purchased_at,omitempty"`
	RequestedAt        time.Time              `json:"requested_at"`
	CancelledAt        *time.Time             `json:"cancelled_at,omitempty"`
	Status             int                    `json:"status"`
	StatusLocale       string                 `json:"status_locale"`
	ReceiptUrl         string                 `json:"receipt_url,omitempty"`
	CardData           *CardDataResponse      `json:"card_data,omitempty"`
	PhoneData          *PhoneDataResponse     `json:"phone_data,omitempty"`
	BankData           *BankDataResponse      `json:"bank_data,omitempty"`
	VbankData          *BankDataResponse      `json:"vbank_data,omitempty"`
	EscrowData         *EscrowDataResponse    `json:"escrow_data,omitempty"`
	CashReceiptData    *CashReceiptResponse   `json:"cash_receipt_data,omitempty"`
	NaverPointData     *PointDataResponse     `json:"naver_point_data,omitempty"`
	KakaoMoneyData     *PointDataResponse     `json:"kakao_money_data,omitempty"`
	PaycoPointData     *PointDataResponse     `json:"payco_point_data,omitempty"`
	TossPointData      *PointDataResponse     `json:"toss_point_data,omitempty"`
	Currency           string                 `json:"currency,omitempty"`
	HttpStatus         int                    `json:"http_status"`
}

// CardDataResponse represents card payment data
type CardDataResponse struct {
	Tid             string `json:"tid"`
	CardApproveNo   string `json:"card_approve_no"`
	CardNo          string `json:"card_no"`
	CardQuota       string `json:"card_quota"`
	CardCompanyCode string `json:"card_company_code"`
	CardCompany     string `json:"card_company"`
	CardInterest    string `json:"card_interest"`
	ReceiptUrl      string `json:"receipt_url,omitempty"`
	CardType        string `json:"card_type,omitempty"`
	CardOwnerType   string `json:"card_owner_type,omitempty"`
	Point           int    `json:"point,omitempty"`
}

// PhoneDataResponse represents phone payment data
type PhoneDataResponse struct {
	Tid        string `json:"tid"`
	AuthNo     string `json:"auth_no,omitempty"`
	Phone      string `json:"phone,omitempty"`
	ReceiptUrl string `json:"receipt_url,omitempty"`
}

// BankDataResponse represents bank/virtual bank payment data
type BankDataResponse struct {
	Tid            string     `json:"tid"`
	BankCode       string     `json:"bank_code"`
	BankName       string     `json:"bank_name"`
	BankUsername   string     `json:"bank_username"`
	BankAccount    string     `json:"bank_account,omitempty"`
	SenderName     string     `json:"sender_name,omitempty"`
	ExpiredAt      *time.Time `json:"expired_at,omitempty"`
	CashReceiptTid string     `json:"cash_receipt_tid,omitempty"`
	CashReceiptType string    `json:"cash_receipt_type,omitempty"`
	CashReceiptNo  string     `json:"cash_receipt_no,omitempty"`
	ReceiptUrl     string     `json:"receipt_url,omitempty"`
}

// EscrowDataResponse represents escrow data
type EscrowDataResponse struct {
	Status             int        `json:"status"`
	StatusLocale       string     `json:"status_locale"`
	ShippingStartedAt  time.Time  `json:"shipping_started_at"`
	ReceiptConfirmedAt *time.Time `json:"receipt_confirmed_at,omitempty"`
}

// CashReceiptResponse represents cash receipt data
type CashReceiptResponse struct {
	Tid             string `json:"tid,omitempty"`
	CashReceiptType int    `json:"cash_receipt_type,omitempty"`
	CashReceiptNo   string `json:"cash_receipt_no,omitempty"`
	ReceiptUrl      string `json:"receipt_url,omitempty"`
}

// PointDataResponse represents point payment data (Naver, Kakao, Payco, Toss)
type PointDataResponse struct {
	Tid string `json:"tid,omitempty"`
}

// CertificateResponse represents certificate/authentication response
type CertificateResponse struct {
	ReceiptId          string               `json:"receipt_id"`
	AuthenticateId     string               `json:"authenticate_id"`
	Pg                 string               `json:"pg"`
	Method             string               `json:"method"`
	MethodOrigin       string               `json:"method_origin"`
	MethodOriginSymbol string               `json:"method_origin_symbol"`
	AuthenticatedAt    time.Time            `json:"authenticated_at"`
	RequestedAt        time.Time            `json:"requested_at"`
	Status             int                  `json:"status"`
	StatusLocale       string               `json:"status_locale"`
	AuthenticateData   AuthenticateDataResp `json:"authenticate_data"`
	HttpStatus         int                  `json:"http_status"`
}

// AuthenticateDataResp represents authentication data in response
type AuthenticateDataResp struct {
	Name             string     `json:"name,omitempty"`
	Phone            string     `json:"phone,omitempty"`
	Unique           string     `json:"unique,omitempty"`
	Birth            *time.Time `json:"birth,omitempty"`
	Gender           int        `json:"gender,omitempty"`
	Foreigner        int        `json:"foreigner,omitempty"`
	Carrier          string     `json:"carrier,omitempty"`
	NumberOfRealarms int        `json:"number_of_realarms,omitempty"`
	Tid              string     `json:"tid"`
}

// SubscriptionBillingResponse represents subscription billing key response
type SubscriptionBillingResponse struct {
	BillingKey         string           `json:"billing_key"`
	BillingData        BillingDataResp  `json:"billing_data"`
	ReceiptId          string           `json:"receipt_id"`
	SubscriptionId     string           `json:"subscription_id"`
	GatewayUrl         string           `json:"gateway_url,omitempty"`
	Metadata           interface{}      `json:"metadata"`
	Pg                 string           `json:"pg"`
	Method             string           `json:"method"`
	MethodOrigin       string           `json:"method_origin,omitempty"`
	MethodOriginSymbol string           `json:"method_origin_symbol,omitempty"`
	MethodSymbol       string           `json:"method_symbol,omitempty"`
	PublishedAt        time.Time        `json:"published_at"`
	RequestedAt        time.Time        `json:"requested_at"`
	ReceiptData        *ReceiptResponse `json:"receipt_data,omitempty"`
	BillingExpireAt    time.Time        `json:"billing_expire_at"`
	Status             int              `json:"status"`
	StatusLocale       string           `json:"status_locale,omitempty"`
	HttpStatus         int              `json:"http_status"`
}

// BillingDataResp represents billing data in response
type BillingDataResp struct {
	CardNo          string `json:"card_no"`
	CardCompany     string `json:"card_company"`
	CardCompanyCode string `json:"card_company_code"`
	CardType        int    `json:"card_type"`
	CardHash        string `json:"card_hash,omitempty"`
	RtnKeyInfo      string `json:"rtn_key_info,omitempty"` // KCP 전용
}

// DestroySubscribeResponse represents the response when destroying a billing key
type DestroySubscribeResponse struct {
	BillingKey string `json:"billing_key"`
	HttpStatus int    `json:"http_status"`
}

// UserTokenResponse represents user token response
type UserTokenResponse struct {
	UserToken  string    `json:"user_token"`
	ExpiredAt  time.Time `json:"expired_at"`
	HttpStatus int       `json:"http_status"`
}

// SubscribePaymentReserveResponse represents reserve payment response
type SubscribePaymentReserveResponse struct {
	ReserveId        string    `json:"reserve_id"`
	ReserveExecuteAt time.Time `json:"reserve_execute_at"`
	HttpStatus       int       `json:"http_status"`
}

// CancelSubscribeReserveResponse represents cancel reserve response
type CancelSubscribeReserveResponse struct {
	ReserveId  string `json:"reserve_id"`
	Success    bool   `json:"success"`
	HttpStatus int    `json:"http_status"`
}

// SubscribePaymentLookupResponse represents reserve payment lookup response
type SubscribePaymentLookupResponse struct {
	ReserveId          string      `json:"reserve_id"`
	ReceiptId          string      `json:"receipt_id,omitempty"`
	OrderId            string      `json:"order_id"`
	Price              float64     `json:"price"`
	TaxFree            float64     `json:"tax_free"`
	OrderName          string      `json:"order_name"`
	User               User        `json:"user"`
	FeedbackUrl        string      `json:"feedback_url"`
	Metadata           interface{} `json:"metadata,omitempty"`
	ContentType        string      `json:"content_type"`
	Version            int         `json:"version"`
	Extra              interface{} `json:"extra"`
	ReserveRequestedAt string      `json:"reserve_requested_at"`
	ReserveExecuteAt   string      `json:"reserve_execute_at"`
	ReserveStartedAt   string      `json:"reserve_started_at"`
	ReserveFinishedAt  string      `json:"reserve_finished_at"`
	ReserveRevokedAt   string      `json:"reserve_revoked_at"`
	Status             int         `json:"status"`
	HttpStatus         int         `json:"http_status"`
}

// WalletPaymentResponse represents wallet payment response
type WalletPaymentResponse struct {
	CancelledPrice     float64                `json:"cancelled_price"`
	WalletData         WalletDataResp         `json:"wallet_data"`
	Metadata           map[string]interface{} `json:"metadata"`
	CancelledTaxFree   float64                `json:"cancelled_tax_free"`
	Method             string                 `json:"method"`
	CardData           *CardDataResponse      `json:"card_data,omitempty"`
	Sandbox            bool                   `json:"sandbox"`
	ReceiptId          string                 `json:"receipt_id"`
	MethodOrigin       string                 `json:"method_origin"`
	OrderName          string                 `json:"order_name"`
	MethodOriginSymbol string                 `json:"method_origin_symbol"`
	ReceiptUrl         string                 `json:"receipt_url"`
	MethodSymbol       string                 `json:"method_symbol"`
	PurchasedAt        string                 `json:"purchased_at"`
	TaxFree            float64                `json:"tax_free"`
	Price              float64                `json:"price"`
	CompanyName        string                 `json:"company_name"`
	Pg                 string                 `json:"pg"`
	StatusLocale       string                 `json:"status_locale"`
	Currency           string                 `json:"currency"`
	HttpStatus         int                    `json:"http_status"`
	OrderId            string                 `json:"order_id"`
	RequestedAt        string                 `json:"requested_at"`
	Status             int                    `json:"status"`
}

// WalletDataResp represents wallet data in response
type WalletDataResp struct {
	Success WalletDataPart `json:"success"`
	Failure []interface{}  `json:"failure"`
}

// =====================================================
// Request Parameter Types - 요청 파라미터 타입 정의
// =====================================================

// Refund represents refund information for cancellation
type Refund struct {
	BankAccount  string `json:"bank_account"`
	BankUsername string `json:"bank_username"`
	BankCode     string `json:"bank_code"`
}

// ExtraModel represents extra options
type ExtraModel struct {
	SubscribeTestPayment bool `json:"subscribe_test_payment,omitempty"`
}

// UserModel represents user information
type UserModel struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
}

// CompanyModel represents company information
type CompanyModel struct {
	Name    string `json:"name,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Zipcode string `json:"zipcode,omitempty"`
	Addr1   string `json:"addr1,omitempty"`
	Addr2   string `json:"addr2,omitempty"`
}

// ItemModel represents item information
type ItemModel struct {
	Id    string  `json:"id,omitempty"`
	Name  string  `json:"name,omitempty"`
	Qty   int     `json:"qty,omitempty"`
	Price float64 `json:"price,omitempty"`
}
