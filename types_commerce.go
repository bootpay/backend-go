package bootpay

import "encoding/json"

// ============================================
// Common Types
// ============================================

// ListParams represents common list query parameters
type ListParams struct {
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Keyword string `json:"keyword,omitempty"`
}

// CommerceAddress represents an address
type CommerceAddress struct {
	AddressId string `json:"address_id,omitempty"`
	Zipcode   string `json:"zipcode,omitempty"`
	Addr1     string `json:"addr1,omitempty"`
	Addr2     string `json:"addr2,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Name      string `json:"name,omitempty"`
	Memo      string `json:"memo,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
}

// ============================================
// User Types
// ============================================

// CommerceUserGroupRef represents a reference to a user group
type CommerceUserGroupRef struct {
	UserGroupId string `json:"user_group_id,omitempty"`
	Name        string `json:"name,omitempty"`
}

// CommerceUser represents a commerce user
type CommerceUser struct {
	UserId    string `json:"user_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`

	// 고객 유형
	MembershipType int `json:"membership_type,omitempty"`

	// 고객 정보
	Name         string `json:"name,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Email        string `json:"email,omitempty"`
	Tel          string `json:"tel,omitempty"`
	Nickname     string `json:"nickname,omitempty"`
	BankUsername string `json:"bank_username,omitempty"`
	BankAccount  string `json:"bank_account,omitempty"`
	BankCode     string `json:"bank_code,omitempty"`
	Comment      string `json:"comment,omitempty"`

	// 최종상태
	Count  int `json:"count,omitempty"`
	Status int `json:"status,omitempty"`

	// 개인 고객
	Gender              int                    `json:"gender,omitempty"`
	Birth               string                 `json:"birth,omitempty"`
	IndividualExtension map[string]interface{} `json:"individual_extension,omitempty"`

	// 쇼핑몰 회원
	LoginId   string `json:"login_id,omitempty"`
	LoginPw   string `json:"login_pw,omitempty"`
	LoginType int    `json:"login_type,omitempty"`

	GroupTags []string               `json:"group_tags,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`

	// 인증정보
	AuthSms   bool   `json:"auth_sms,omitempty"`
	AuthPhone bool   `json:"auth_phone,omitempty"`
	AuthEmail bool   `json:"auth_email,omitempty"`
	Ci        string `json:"ci,omitempty"`
	Cd        string `json:"cd,omitempty"`

	JoinAt          string `json:"join_at,omitempty"`
	JoinConfirmType int    `json:"join_confirm_type,omitempty"`
	LastedAt        string `json:"lasted_at,omitempty"`

	// 약관 동의
	MarketingAcceptType     int      `json:"marketing_accept_type,omitempty"`
	MarketingAcceptCreateAt string   `json:"marketing_accept_create_at,omitempty"`
	MarketingAcceptUpdateAt string   `json:"marketing_accept_update_at,omitempty"`
	TermIds                 []string `json:"term_ids,omitempty"`

	Group *CommerceUserGroupRef `json:"group,omitempty"`

	ExternalUid string `json:"external_uid,omitempty"`
	IsExternal  string `json:"is_external,omitempty"`
	UserGroupId string `json:"user_group_id,omitempty"`
}

// UserListParams represents user list query parameters
// ⚠️ The membership filter key read by the server (v1/users_controller#index) is membership_type.
// member_type was silently ignored (no error — the full list came back), so it is kept only as a
// legacy alias and is mapped onto membership_type when MembershipType is unset.
type UserListParams struct {
	ListParams
	MembershipType int `json:"membership_type,omitempty"`
	// MemberType is the legacy alias of MembershipType (kept for backward compatibility)
	MemberType int    `json:"member_type,omitempty"`
	Type       string `json:"type,omitempty"`
}

// CommerceUserTokenResponse represents commerce user token response
type CommerceUserTokenResponse struct {
	AccessToken string        `json:"access_token,omitempty"`
	ExpiredAt   string        `json:"expired_at,omitempty"`
	User        *CommerceUser `json:"user,omitempty"`
}

// CommerceUserLoginResponse represents commerce user login response
type CommerceUserLoginResponse struct {
	AccessToken string        `json:"access_token,omitempty"`
	ExpiredAt   string        `json:"expired_at,omitempty"`
	User        *CommerceUser `json:"user,omitempty"`
}

// Constants for V1 Mall API join-check types (GET users/join/{type})
const (
	MALL_USER_JOIN_CHECK_EMAIL_EXIST                 = "email-exist"
	MALL_USER_JOIN_CHECK_ID_EXIST                    = "id-exist"
	MALL_USER_JOIN_CHECK_PHONE_EXIST                 = "phone-exist"
	MALL_USER_JOIN_CHECK_UID_EXIST                   = "uid-exist"
	MALL_USER_JOIN_CHECK_GROUP_BUSINESS_NUMBER_EXIST = "group-business-number-exist"
)

// MallUserLoginParams represents V1 Mall API member login parameters (POST users/login)
// CorporateType: 0 = individual, 1 = corporate (always sent; defaults to 0)
type MallUserLoginParams struct {
	LoginId       string `json:"login_id"`
	Password      string `json:"password"`
	CorporateType int    `json:"corporate_type"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// MallUserJoinParams represents V1 Mall API member join parameters (POST users/join)
type MallUserJoinParams struct {
	LoginId  string `json:"login_id"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	// Gender: explicit 0 (female) must be sendable, so pointer type
	Gender *int   `json:"gender,omitempty"`
	Birth  string `json:"birth,omitempty"`
	// CorporateType: 0 = individual, 1 = corporate (always sent; defaults to 0)
	CorporateType int                    `json:"corporate_type"`
	Group         map[string]interface{} `json:"group,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// ============================================
// UserGroup Types
// ============================================

// Constants for corporate type
const (
	CORPORATE_TYPE_INDIVIDUAL = 1
	CORPORATE_TYPE_CORPORATE  = 2
)

// CommerceUserGroup represents a user group
type CommerceUserGroup struct {
	UserGroupId   string `json:"user_group_id,omitempty"`
	SellerId      string `json:"seller_id,omitempty"`
	ProjectId     string `json:"project_id,omitempty"`
	CorporateType int    `json:"corporate_type,omitempty"`

	Bank     string `json:"bank,omitempty"`
	BankCode string `json:"bank_code,omitempty"`

	Count         int    `json:"count,omitempty"`
	LastUpdatedAt string `json:"last_updated_at,omitempty"`
	Status        int    `json:"status,omitempty"`

	Phone              string                 `json:"phone,omitempty"`
	Email              string                 `json:"email,omitempty"`
	Zipcode            string                 `json:"zipcode,omitempty"`
	Address            string                 `json:"address,omitempty"`
	AddressDetail      string                 `json:"address_detail,omitempty"`
	CorporateExtension map[string]interface{} `json:"corporate_extension,omitempty"`
	AuthBank           bool                   `json:"auth_bank,omitempty"`

	CompanyName          string `json:"company_name,omitempty"`
	BusinessNumber       string `json:"business_number,omitempty"`
	RegistrationNumber   string `json:"registration_number,omitempty"`
	CorporateEstablished string `json:"corporate_established,omitempty"`
	BusinessType         string `json:"business_type,omitempty"`
	BusinessCategory     string `json:"business_category,omitempty"`
	CeoName              string `json:"ceo_name,omitempty"`
	AuthCompany          bool   `json:"auth_company,omitempty"`

	ManagerName  string `json:"manager_name,omitempty"`
	ManagerPhone string `json:"manager_phone,omitempty"`
	ManagerEmail string `json:"manager_email,omitempty"`

	PersonalCustomsClearanceCode string `json:"personal_customs_clearance_code,omitempty"`

	Point                   int    `json:"point,omitempty"`
	Accumulation            int    `json:"accumulation,omitempty"`
	MarketingAcceptType     int    `json:"marketing_accept_type,omitempty"`
	MarketingAcceptCreateAt string `json:"marketing_accept_create_at,omitempty"`
	MarketingAcceptUpdateAt string `json:"marketing_accept_update_at,omitempty"`

	UseSubscriptionAggregateTransaction bool `json:"use_subscription_aggregate_transaction,omitempty"`
	SubscriptionMonthDay                int  `json:"subscription_month_day,omitempty"`
	SubscriptionWeekDay                 int  `json:"subscription_week_day,omitempty"`

	UseLimit        bool   `json:"use_limit,omitempty"`
	PurchaseLimit   int    `json:"purchase_limit,omitempty"`
	SubscribedLimit int    `json:"subscribed_limit,omitempty"`
	LimitMessage    string `json:"limit_message,omitempty"`
	ExternalUid     string `json:"external_uid,omitempty"`
	IsExternal      string `json:"is_external,omitempty"`
}

// UserGroupListParams represents user group list query parameters
type UserGroupListParams struct {
	ListParams
	CorporateType int `json:"corporate_type,omitempty"`
}

// UserGroupLimitParams represents user group limit parameters
// LimitMonthPurchase / LimitWeekPurchase are the official server parameter names —
// limits are never applied through Update, only through this dedicated route (server scope: manager).
type UserGroupLimitParams struct {
	UserGroupId     string `json:"user_group_id"`
	UseLimit        bool   `json:"use_limit,omitempty"`
	PurchaseLimit   int    `json:"purchase_limit,omitempty"`
	SubscribedLimit int    `json:"subscribed_limit,omitempty"`
	LimitMessage    string `json:"limit_message,omitempty"`
	// explicit 0 must be sendable, so pointer types
	LimitMonthPurchase *int `json:"limit_month_purchase,omitempty"`
	LimitWeekPurchase  *int `json:"limit_week_purchase,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// UserGroupAggregateTransactionParams represents aggregate transaction parameters
type UserGroupAggregateTransactionParams struct {
	UserGroupId                         string `json:"user_group_id"`
	UseSubscriptionAggregateTransaction bool   `json:"use_subscription_aggregate_transaction,omitempty"`
	SubscriptionMonthDay                int    `json:"subscription_month_day,omitempty"`
	SubscriptionWeekDay                 int    `json:"subscription_week_day,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// ============================================
// Product Types
// ============================================

// CommerceProductOption represents a product option
type CommerceProductOption struct {
	OptionId string `json:"option_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Price    int    `json:"price,omitempty"`
	Stock    int    `json:"stock,omitempty"`
}

// CommerceSubscriptionSetting represents subscription settings
type CommerceSubscriptionSetting struct {
	SubscriptionSettingId string `json:"subscription_setting_id,omitempty"`
	PeriodType            string `json:"period_type,omitempty"`
	PeriodValue           int    `json:"period_value,omitempty"`
	BillingDay            int    `json:"billing_day,omitempty"`
	BillingCount          int    `json:"billing_count,omitempty"`
}

// CommerceProduct represents a product
type CommerceProduct struct {
	ProductId             string `json:"product_id,omitempty"`
	CategoryId            string `json:"category_id,omitempty"`
	ProjectId             string `json:"project_id,omitempty"`
	SellerId              string `json:"seller_id,omitempty"`
	SubscriptionSettingId string `json:"subscription_setting_id,omitempty"`
	DeliveryShippingId    string `json:"delivery_shipping_id,omitempty"`
	BrandId               string `json:"brand_id,omitempty"`
	ManufacturerId        string `json:"manufacturer_id,omitempty"`

	ExUid string `json:"ex_uid,omitempty"`

	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	Images         []string `json:"images,omitempty"`
	Type           int      `json:"type,omitempty"`
	TaxType        int      `json:"tax_type,omitempty"`
	UseStock       bool     `json:"use_stock,omitempty"`
	Stock          int      `json:"stock,omitempty"`
	UseOptionStock bool     `json:"use_option_stock,omitempty"`
	UseStockSafe   bool     `json:"use_stock_safe,omitempty"`
	StockSafe      int      `json:"stock_safe,omitempty"`

	DisplayPrice      int    `json:"display_price,omitempty"`
	TaxFreePrice      int    `json:"tax_free_price,omitempty"`
	UseDiscount       bool   `json:"use_discount,omitempty"`
	DiscountPrice     int    `json:"discount_price,omitempty"`
	DiscountPriceType int    `json:"discount_price_type,omitempty"`
	UseDiscountPeriod bool   `json:"use_discount_period,omitempty"`
	DiscountStartAt   string `json:"discount_start_at,omitempty"`
	DiscountEndAt     string `json:"discount_end_at,omitempty"`

	UseAccumulation       bool `json:"use_accumulation,omitempty"`
	AccumulationPoint     int  `json:"accumulation_point,omitempty"`
	AccumulationPointType int  `json:"accumulation_point_type,omitempty"`

	StatusDisplay    bool   `json:"status_display,omitempty"`
	UseDisplayPeriod bool   `json:"use_display_period,omitempty"`
	DisplayStartAt   string `json:"display_start_at,omitempty"`
	DisplayEndAt     string `json:"display_end_at,omitempty"`
	StatusSale       bool   `json:"status_sale,omitempty"`
	UseSalePeriod    bool   `json:"use_sale_period,omitempty"`
	SaleStartAt      string `json:"sale_start_at,omitempty"`
	SaleEndAt        string `json:"sale_end_at,omitempty"`

	CountSale   int `json:"count_sale,omitempty"`
	CountQna    int `json:"count_qna,omitempty"`
	CountLike   int `json:"count_like,omitempty"`
	CountReview int `json:"count_review,omitempty"`

	Barcode        string   `json:"barcode,omitempty"`
	Sku            string   `json:"sku,omitempty"`
	SearchTags     []string `json:"search_tags,omitempty"`
	EventTags      []string `json:"event_tags,omitempty"`
	TargetUserTags []string `json:"target_user_tags,omitempty"`
	DeliveryTags   []string `json:"delivery_tags,omitempty"`
	EmotionTags    []string `json:"emotion_tags,omitempty"`

	UseCoupon   bool   `json:"use_coupon,omitempty"`
	UseMinor    bool   `json:"use_minor,omitempty"`
	UseFreeGift bool   `json:"use_free_gift,omitempty"`
	FreeGift    string `json:"free_gift,omitempty"`

	UseBulkPurchaseDiscount bool                   `json:"use_bulk_purchase_discount,omitempty"`
	BulkPurchaseDiscount    map[string]interface{} `json:"bulk_purchase_discount,omitempty"`

	UseReviewPoint bool                   `json:"use_review_point,omitempty"`
	ReviewPoint    map[string]interface{} `json:"review_point,omitempty"`

	UseSeo             bool     `json:"use_seo,omitempty"`
	SeoPageTitle       string   `json:"seo_page_title,omitempty"`
	SeoMetaDescription string   `json:"seo_meta_description,omitempty"`
	SeoMetaTags        []string `json:"seo_meta_tags,omitempty"`

	ModelId          string `json:"model_id,omitempty"`
	ModelName        string `json:"model_name,omitempty"`
	ManufacturerName string `json:"manufacturer_name,omitempty"`
	BrandName        string `json:"brand_name,omitempty"`
	OriginCode       string `json:"origin_code,omitempty"`
	OriginName       string `json:"origin_name,omitempty"`
	Importer         string `json:"importer,omitempty"`

	Used           bool   `json:"used,omitempty"`
	ExpiredAt      string `json:"expired_at,omitempty"`
	ManufacturedAt string `json:"manufactured_at,omitempty"`

	UseSetupFee   bool   `json:"use_setup_fee,omitempty"`
	SetupFeeValue int    `json:"setup_fee_value,omitempty"`
	SetupFeeType  int    `json:"setup_fee_type,omitempty"`
	SetupFeeName  string `json:"setup_fee_name,omitempty"`
	SetupFeeText  string `json:"setup_fee_text,omitempty"`

	UseDeliveryShipping       bool   `json:"use_delivery_shipping,omitempty"`
	DeliveryShippingFeeType   int    `json:"delivery_shipping_fee_type,omitempty"`
	UseOverseasShipping       bool   `json:"use_overseas_shipping,omitempty"`
	UseDeliveryShippingBundle bool   `json:"use_delivery_shipping_bundle,omitempty"`
	DeliveryShippingBundleId  string `json:"delivery_shipping_bundle_id,omitempty"`

	UseSubscription      bool `json:"use_subscription,omitempty"`
	UseSubscriptionTimes bool `json:"use_subscription_times,omitempty"`
	UseProductPrice      bool `json:"use_product_price,omitempty"`

	UseCancel     bool `json:"use_cancel,omitempty"`
	UseAbleRefund bool `json:"use_able_refund,omitempty"`
	UseAbleCart   bool `json:"use_able_cart,omitempty"`

	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`

	Options             []CommerceProductOption      `json:"options,omitempty"`
	SubscriptionSetting *CommerceSubscriptionSetting `json:"subscription_setting,omitempty"`
}

// ProductListParams represents product list query parameters (GET /v1/products)
//
// ⚠️ The server (v1/products_controller#index) reads only
// page / limit / keyword / category_id / ex_uid / sort.
// Type, PeriodType, SAt, EAt and CategoryCode are still sent for backward compatibility
// but are silently ignored — the full list comes back, so do not mistake them for filters.
// (Keyword takes effect from the 26-08-26 server change; older deployments ignore it.)
type ProductListParams struct {
	ListParams
	// CategoryId filters by category (child categories included)
	CategoryId string `json:"category_id,omitempty"`
	// ExUid looks a product up by its external UID
	ExUid string `json:"ex_uid,omitempty"`
	// Sort key — position | created_at | -created_at | price | -price | -sold
	Sort string `json:"sort,omitempty"`

	// ── The server does not read the fields below (kept for backward compatibility) ──
	// Type is numeric here, while the server's product type filter takes a string
	// (subscription | discount | normal) — the value systems do not match.
	Type         int    `json:"type,omitempty"`
	PeriodType   string `json:"period_type,omitempty"`
	SAt          string `json:"s_at,omitempty"`
	EAt          string `json:"e_at,omitempty"`
	CategoryCode string `json:"category_code,omitempty"`
}

// MallProductListParams represents V1 Mall API product list query parameters (GET products)
// page/limit default to 1/20 when unset.
// ⚠️ The server (v1/products_controller#index) reads only
// page / limit / keyword / category_id / ex_uid / sort — Type, PeriodType, SAt, EAt and
// CategoryCode are silently ignored. Keyword takes effect from the 26-08-26 server change.
type MallProductListParams struct {
	ProductListParams
	// CategoryId / ExUid / Sort are declared here as well as on the embedded
	// ProductListParams. Removing them would break existing composite literals
	// (Go forbids setting a promoted field in a struct literal), so they shadow the
	// embedded ones. Products() prefers these and falls back to the embedded values.
	CategoryId string `json:"category_id,omitempty"`
	// ExUid looks a product up by its external UID (read by the controller as params[:ex_uid])
	ExUid string `json:"ex_uid,omitempty"`
	Sort  string `json:"sort,omitempty"`
	// UserJwt is sent as the Bootpay-User-JWT header (attached only when present)
	UserJwt string `json:"-"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// ProductStatusParams represents product status change parameters
// ⚠️ stock is changed through Update, not here.
type ProductStatusParams struct {
	ProductId     string `json:"product_id"`
	Status        int    `json:"status"`
	StatusDisplay bool   `json:"status_display,omitempty"`
	StatusSale    bool   `json:"status_sale,omitempty"`
	// explicit false must be sendable, so pointer types
	StatusFrozen     *bool  `json:"status_frozen,omitempty"`
	StatusReview     *bool  `json:"status_review,omitempty"`
	UseDisplayPeriod *bool  `json:"use_display_period,omitempty"`
	DisplayStartAt   string `json:"display_start_at,omitempty"`
	DisplayEndAt     string `json:"display_end_at,omitempty"`
	UseSalePeriod    *bool  `json:"use_sale_period,omitempty"`
	SaleStartAt      string `json:"sale_start_at,omitempty"`
	SaleEndAt        string `json:"sale_end_at,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// ============================================
// Invoice Types
// ============================================

// Constants for invoice send type
const (
	INVOICE_SEND_TYPE_SMS   = 1
	INVOICE_SEND_TYPE_KAKAO = 2
	INVOICE_SEND_TYPE_EMAIL = 3
	INVOICE_SEND_TYPE_PUSH  = 4
)

// CommerceInvoiceItem represents an invoice item
type CommerceInvoiceItem struct {
	InvoiceItemId string `json:"invoice_item_id,omitempty"`
	Name          string `json:"name,omitempty"`
	Price         int    `json:"price,omitempty"`
	Qty           int    `json:"qty,omitempty"`
	TaxFreePrice  int    `json:"tax_free_price,omitempty"`
}

// CommerceInvoice represents an invoice
type CommerceInvoice struct {
	InvoiceId string `json:"invoice_id,omitempty"`
	ProjectId string `json:"project_id,omitempty"`
	SellerId  string `json:"seller_id,omitempty"`

	Name        string `json:"name,omitempty"`
	Title       string `json:"title,omitempty"`
	Memo        string `json:"memo,omitempty"`
	ProductName string `json:"product_name,omitempty"`

	CreatedOwnerId   string `json:"created_owner_id,omitempty"`
	CreatedOwnerType int    `json:"created_owner_type,omitempty"`

	Unit     int                    `json:"unit,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	RequestId string `json:"request_id,omitempty"`
	Sku       string `json:"sku,omitempty"`

	UseRedirect bool   `json:"use_redirect,omitempty"`
	RedirectUrl string `json:"redirect_url,omitempty"`

	Type     int    `json:"type,omitempty"`
	ParentId string `json:"parent_id,omitempty"`

	SubscriptionType    int    `json:"subscription_type,omitempty"`
	SubscriptionStartAt string `json:"subscription_start_at,omitempty"`
	SubscriptionEndAt   string `json:"subscription_end_at,omitempty"`

	ExpiredAt string `json:"expired_at,omitempty"`
	Status    int    `json:"status,omitempty"`
	Deleted   bool   `json:"deleted,omitempty"`

	UserCollectionType int  `json:"user_collection_type,omitempty"`
	UseLinkRedirect    bool `json:"use_link_redirect,omitempty"`

	UserId string `json:"user_id,omitempty"`

	SendStatus int   `json:"send_status,omitempty"`
	SendTypes  []int `json:"send_types,omitempty"`

	MessageTemplateId string `json:"message_template_id,omitempty"`
	MessageId         string `json:"message_id,omitempty"`
	MessageFrom       string `json:"message_from,omitempty"`
	MessageType       int    `json:"message_type,omitempty"`
	MessageResponse   string `json:"message_response,omitempty"`

	SentAt string `json:"sent_at,omitempty"`
	PayAt  string `json:"pay_at,omitempty"`

	Price        int `json:"price,omitempty"`
	TaxFreePrice int `json:"tax_free_price,omitempty"`

	UseEditableUsername bool `json:"use_editable_username,omitempty"`
	UseEditablePhone    bool `json:"use_editable_phone,omitempty"`
	UseEditableEmail    bool `json:"use_editable_email,omitempty"`
	UseMemo             bool `json:"use_memo,omitempty"`

	ProductIds       []string `json:"product_ids,omitempty"`
	ProductOptionIds []string `json:"product_option_ids,omitempty"`

	Tags []string `json:"tags,omitempty"`

	Password string `json:"password,omitempty"`
	OrderId  string `json:"order_id,omitempty"`
	Uuid     string `json:"uuid,omitempty"`

	WebhookUrl        string `json:"webhook_url,omitempty"`
	HeaderContentType int    `json:"header_content_type,omitempty"`
	WebhookRetryCount int    `json:"webhook_retry_count,omitempty"`

	ProductType int  `json:"product_type,omitempty"`
	IsOpenLink  bool `json:"is_open_link,omitempty"`

	InvoiceItems  []CommerceInvoiceItem `json:"invoice_items,omitempty"`
	SelectedUsers []string              `json:"selected_users,omitempty"`
}

// InvoiceListParams represents invoice list query parameters (GET invoices)
// page/limit default to 1/24 when unset (24 matches the server default).
// ⚠️ The response data is { list: [...], count: N } — not { items, total }.
type InvoiceListParams struct {
	ListParams
	CsType string `json:"cs_type,omitempty"`
	UserId string `json:"user_id,omitempty"`
	// explicit 0 must be sendable, so pointer type
	ProductType *int   `json:"product_type,omitempty"`
	CssAt       string `json:"css_at,omitempty"`
	CseAt       string `json:"cse_at,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// ============================================
// Order Types
// ============================================

// Constants for subscription billing type
const (
	SUBSCRIPTION_BILLING_TYPE_NONE  = 0
	SUBSCRIPTION_BILLING_TYPE_EACH  = 1
	SUBSCRIPTION_BILLING_TYPE_GROUP = 2
)

// CommerceChosenProductOption represents a chosen product option
type CommerceChosenProductOption struct {
	ChosenProductOptionId string `json:"chosen_product_option_id,omitempty"`
	ProductId             string `json:"product_id,omitempty"`
	ProductOptionId       string `json:"product_option_id,omitempty"`
	ProductName           string `json:"product_name,omitempty"`
	OptionName            string `json:"option_name,omitempty"`
	Price                 int    `json:"price,omitempty"`
	TaxFreePrice          int    `json:"tax_free_price,omitempty"`
	Qty                   int    `json:"qty,omitempty"`
}

// CommerceOrderCancellationRequestHistory represents order cancellation request history
type CommerceOrderCancellationRequestHistory struct {
	OrderCancellationRequestHistoryId string `json:"order_cancellation_request_history_id,omitempty"`
	OrderId                           string `json:"order_id,omitempty"`
	Status                            int    `json:"status,omitempty"`
	CancelReason                      string `json:"cancel_reason,omitempty"`
	CancelType                        int    `json:"cancel_type,omitempty"`
	RequestedAt                       string `json:"requested_at,omitempty"`
	ProcessedAt                       string `json:"processed_at,omitempty"`
}

// CommerceOrder represents an order
type CommerceOrder struct {
	OrderId              string                        `json:"order_id,omitempty"`
	OrderPreId           string                        `json:"order_pre_id,omitempty"`
	ChosenProductOptions []CommerceChosenProductOption `json:"chosen_product_options,omitempty"`

	ParentOrderId  string `json:"parent_order_id,omitempty"`
	UserId         string `json:"user_id,omitempty"`
	SellerId       string `json:"seller_id,omitempty"`
	ProjectId      string `json:"project_id,omitempty"`
	Status         int    `json:"status,omitempty"`
	Currency       int    `json:"currency,omitempty"`
	IsSubscription bool   `json:"is_subscription,omitempty"`
	IsLeaf         bool   `json:"is_leaf,omitempty"`
	TotalPrice     int    `json:"total_price,omitempty"`
	TaxFreePrice   int    `json:"tax_free_price,omitempty"`
	DiscountAmount int    `json:"discount_amount,omitempty"`
	DeliveryPrice  int    `json:"delivery_price,omitempty"`
	PaymentMethod  string `json:"payment_method,omitempty"`
	ReceiptId      string `json:"receipt_id,omitempty"`
	WebhookUrl     string `json:"webhook_url,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`

	CancelledRequestHistory []CommerceOrderCancellationRequestHistory `json:"cancelled_request_history,omitempty"`
}

// OrderListParams represents order list query parameters
// SearchDateFrom / SearchDateTo are the official date keys (CssAt / CseAt remain as server aliases).
type OrderListParams struct {
	ListParams
	UserId                  string   `json:"user_id,omitempty"`
	UserGroupId             string   `json:"user_group_id,omitempty"`
	Status                  []int    `json:"status,omitempty"`
	PaymentStatus           []int    `json:"payment_status,omitempty"`
	CsType                  string   `json:"cs_type,omitempty"`
	SearchDateFrom          string   `json:"search_date_from,omitempty"`
	SearchDateTo            string   `json:"search_date_to,omitempty"`
	CssAt                   string   `json:"css_at,omitempty"`
	CseAt                   string   `json:"cse_at,omitempty"`
	SubscriptionBillingType int      `json:"subscription_billing_type,omitempty"`
	OrderSubscriptionIds    []string `json:"order_subscription_ids,omitempty"`
}

// ============================================
// OrderCancel Types
// ============================================

// OrderCancelListParams represents order cancel list query parameters
type OrderCancelListParams struct {
	OrderId     string `json:"order_id,omitempty"`
	OrderNumber string `json:"order_number,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// CancelProduct represents a product to cancel
type CancelProduct struct {
	OrderProductId string `json:"order_product_id,omitempty"`
	ProductId      string `json:"product_id,omitempty"`
	Qty            int    `json:"qty,omitempty"`
	CancelPrice    int    `json:"cancel_price,omitempty"`
}

// CancelOrderSubscriptionBill represents a subscription bill to cancel
type CancelOrderSubscriptionBill struct {
	OrderSubscriptionBillId string `json:"order_subscription_bill_id,omitempty"`
	CancelPrice             int    `json:"cancel_price,omitempty"`
}

// RequestCancelParameter represents cancel request parameters
type RequestCancelParameter struct {
	CancelProducts               []CancelProduct               `json:"cancel_products,omitempty"`
	CancelOrderSubscriptionBills []CancelOrderSubscriptionBill `json:"cancel_order_subscription_bills,omitempty"`
	CancelReason                 string                        `json:"cancel_reason,omitempty"`
	CancelType                   int                           `json:"cancel_type,omitempty"`
	RefundPrice                  int                           `json:"refund_price,omitempty"`
}

// OrderCancelParams represents order cancel parameters
type OrderCancelParams struct {
	OrderNumber             string                  `json:"order_number,omitempty"`
	RequestCancelParameters *RequestCancelParameter `json:"request_cancel_parameters,omitempty"`
	IsSupervisor            bool                    `json:"is_supervisor,omitempty"`
}

// OrderCancelActionParams represents order cancel action parameters (approve / reject)
// The official id name is OrderCancellationRequestId; the old OrderCancelRequestHistoryId
// keeps working (when both are set, the official one wins).
// The value the server reads for approve/reject is Message.
type OrderCancelActionParams struct {
	OrderCancellationRequestId  string `json:"order_cancellation_request_id,omitempty"`
	OrderCancelRequestHistoryId string `json:"order_cancel_request_history_id"`
	Message                     string `json:"message,omitempty"`
	CancelReason                string `json:"cancel_reason,omitempty"`
	RefundPrice                 int    `json:"refund_price,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// CommerceOrderCancelRequestHistory represents order cancel request history
type CommerceOrderCancelRequestHistory struct {
	OrderCancelRequestHistoryId string `json:"order_cancel_request_history_id,omitempty"`
	OrderId                     string `json:"order_id,omitempty"`
	OrderNumber                 string `json:"order_number,omitempty"`
	Status                      int    `json:"status,omitempty"`
	CancelReason                string `json:"cancel_reason,omitempty"`
	CancelType                  int    `json:"cancel_type,omitempty"`
	RequestedAt                 string `json:"requested_at,omitempty"`
	ProcessedAt                 string `json:"processed_at,omitempty"`
	RefundPrice                 int    `json:"refund_price,omitempty"`
}

// ============================================
// OrderSubscription Types
// ============================================

// CommerceOrderSubscription represents an order subscription
type CommerceOrderSubscription struct {
	OrderSubscriptionId string `json:"order_subscription_id,omitempty"`
	SellerId            string `json:"seller_id,omitempty"`
	ProjectId           string `json:"project_id,omitempty"`
	OrderId             string `json:"order_id,omitempty"`
	OrderPreId          string `json:"order_pre_id,omitempty"`
	UserId              string `json:"user_id,omitempty"`
	UserGroupId         string `json:"user_group_id,omitempty"`
	WalletId            string `json:"wallet_id,omitempty"`

	SubscriptionBillingType      int `json:"subscription_billing_type,omitempty"`
	SubscriptionPaymentCycleType int `json:"subscription_payment_cycle_type,omitempty"`
	SubscriptionPaymentDate      int `json:"subscription_payment_date,omitempty"`
	SubscriptionBillingBaseDay   int `json:"subscription_billing_base_day,omitempty"`

	Quantity       int  `json:"quantity,omitempty"`
	IsFirstPrepaid bool `json:"is_first_prepaid,omitempty"`

	OneUnitPrice        int `json:"one_unit_price,omitempty"`
	OneUnitTaxFreePrice int `json:"one_unit_tax_free_price,omitempty"`
	Price               int `json:"price,omitempty"`
	TaxFreePrice        int `json:"tax_free_price,omitempty"`
	SetupPrice          int `json:"setup_price,omitempty"`

	Unit        int      `json:"unit,omitempty"`
	OrderName   string   `json:"order_name,omitempty"`
	ProductName string   `json:"product_name,omitempty"`
	OptionNames []string `json:"option_names,omitempty"`

	ServiceStartAt string `json:"service_start_at,omitempty"`
	ServiceEndAt   string `json:"service_end_at,omitempty"`

	LastBillingCreatedAt string `json:"last_billing_created_at,omitempty"`
	LatestPurchasedAt    string `json:"latest_purchased_at,omitempty"`
	LatestFailedAt       string `json:"latest_failed_at,omitempty"`
	PaymentNextAt        string `json:"payment_next_at,omitempty"`

	CurrentDuration           int `json:"current_duration,omitempty"`
	CreatedLastDuration       int `json:"created_last_duration,omitempty"`
	PaymentLastDuration       int `json:"payment_last_duration,omitempty"`
	TotalSubscriptionDuration int `json:"total_subscription_duration,omitempty"`

	MembershipType       int  `json:"membership_type,omitempty"`
	UseSubscriptionTimes bool `json:"use_subscription_times,omitempty"`

	RenewalStatus int    `json:"renewal_status,omitempty"`
	CancelStatus  int    `json:"cancel_status,omitempty"`
	Status        int    `json:"status,omitempty"`
	CancelAt      string `json:"cancel_at,omitempty"`
}

// OrderSubscriptionListParams represents order subscription list query parameters
// ⚠️ Date keys are SearchDateFrom / SearchDateTo (or SAt / EAt) — different from orders' CssAt / CseAt.
type OrderSubscriptionListParams struct {
	ListParams
	SearchDateFrom string `json:"search_date_from,omitempty"`
	SearchDateTo   string `json:"search_date_to,omitempty"`
	SAt            string `json:"s_at,omitempty"`
	EAt            string `json:"e_at,omitempty"`
	RequestType    string `json:"request_type,omitempty"`
	UserGroupId    string `json:"user_group_id,omitempty"`
	UserId         string `json:"user_id,omitempty"`
	// OrderNumber looks a subscription up in reverse from an order number (server #index reads params[:order_number])
	OrderNumber string `json:"order_number,omitempty"`
	// Status: explicit 0 must be sendable, so pointer type
	Status *int `json:"status,omitempty"`
}

// OrderSubscriptionUpdateParams represents order subscription contract update parameters
// Only changed values need to be sent. The server requires supervisor scope.
type OrderSubscriptionUpdateParams struct {
	OrderSubscriptionId       string `json:"order_subscription_id"`
	ProductId                 string `json:"product_id,omitempty"`
	ProductOptionId           string `json:"product_option_id,omitempty"`
	OrderName                 string `json:"order_name,omitempty"`
	TotalSubscriptionDuration int    `json:"total_subscription_duration,omitempty"`
	Quantity                  int    `json:"quantity,omitempty"`
	AddressId                 string `json:"address_id,omitempty"`
	Username                  string `json:"username,omitempty"`
	Phone                     string `json:"phone,omitempty"`
	Email                     string `json:"email,omitempty"`
	UseFreeTrial              *bool  `json:"use_free_trial,omitempty"`
	FreeTrialDay              int    `json:"free_trial_day,omitempty"`
	ServiceStartAt            string `json:"service_start_at,omitempty"`
	NextBillingAt             string `json:"next_billing_at,omitempty"`
	BillingKey                string `json:"billing_key,omitempty"`
	Status                    int    `json:"status,omitempty"`
	PaymentNextAt             string `json:"payment_next_at,omitempty"`
	ServiceEndAt              string `json:"service_end_at,omitempty"`
	// Price is the base charge amount per billing cycle.
	// Changing it immediately recalculates the amount of the READY (scheduled) cycle and
	// every cycle created afterwards. Already-paid cycles are untouched.
	// Values of 0 or less are rejected by the server (0 is therefore not sent).
	// To add/subtract on specific cycles only, use OrderSubscriptionAdjustment.Create instead.
	Price int `json:"price,omitempty"`
	// Memo is the reason recorded on the change history (SUBSCRIPTION_ACTION_UPDATE)
	Memo string `json:"memo,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// OrderSubscriptionPauseParams represents subscription pause parameters
type OrderSubscriptionPauseParams struct {
	OrderSubscriptionId string `json:"order_subscription_id,omitempty"`
	OrderNumber         string `json:"order_number,omitempty"`
	Reason              string `json:"reason,omitempty"`
	PausedAt            string `json:"paused_at,omitempty"`
	ExpectedResumeAt    string `json:"expected_resume_at,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// OrderSubscriptionResumeParams represents subscription resume parameters
type OrderSubscriptionResumeParams struct {
	OrderSubscriptionId string `json:"order_subscription_id,omitempty"`
	OrderNumber         string `json:"order_number,omitempty"`
	Reason              string `json:"reason,omitempty"`
	ResumeAt            string `json:"resume_at,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// OrderSubscriptionPurchaseParams represents mid-contract purchase request parameters
// (POST order_subscriptions/requests/ing/purchase)
type OrderSubscriptionPurchaseParams struct {
	OrderSubscriptionId string `json:"order_subscription_id,omitempty"`
	OrderNumber         string `json:"order_number,omitempty"`
	Price               int    `json:"price,omitempty"`
	TaxFreePrice        int    `json:"tax_free_price,omitempty"`
	Reason              string `json:"reason,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// OrderSubscriptionTransferParams represents subscription transfer request parameters
// (POST order_subscriptions/requests/ing/transfer)
type OrderSubscriptionTransferParams struct {
	OrderSubscriptionId string `json:"order_subscription_id,omitempty"`
	NewUserId           string `json:"new_user_id,omitempty"`
	NewUsername         string `json:"new_username,omitempty"`
	NewUserEmail        string `json:"new_user_email,omitempty"`
	NewUserPhone        string `json:"new_user_phone,omitempty"`
	NewUserAddress      string `json:"new_user_address,omitempty"`
	WalletId            string `json:"wallet_id,omitempty"`
	Reason              string `json:"reason,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// OrderSubscriptionTerminationParams represents subscription termination parameters
type OrderSubscriptionTerminationParams struct {
	OrderSubscriptionId string `json:"order_subscription_id,omitempty"`
	OrderNumber         string `json:"order_number,omitempty"`
	TerminationFee      int    `json:"termination_fee,omitempty"`
	LastBillRefundPrice int    `json:"last_bill_refund_price,omitempty"`
	FinalFee            int    `json:"final_fee,omitempty"`
	ServiceEndAt        string `json:"service_end_at,omitempty"`
	Reason              string `json:"reason,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// CalcTerminateFeeResponse represents terminate fee calculation response
type CalcTerminateFeeResponse struct {
	TerminationFee      int `json:"termination_fee,omitempty"`
	RefundAmount        int `json:"refund_amount,omitempty"`
	LastBillRefundPrice int `json:"last_bill_refund_price,omitempty"`
	FinalFee            int `json:"final_fee,omitempty"`
}

// Supervisor OrderSubscription Types

// SupervisorOrderSubscriptionApproveParams represents supervisor approve parameters
type SupervisorOrderSubscriptionApproveParams struct {
	Reason string `json:"reason,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// SupervisorOrderSubscriptionRejectParams represents supervisor reject parameters
type SupervisorOrderSubscriptionRejectParams struct {
	Reason string `json:"reason,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// SupervisorOrderSubscriptionTerminateParams represents supervisor terminate parameters
type SupervisorOrderSubscriptionTerminateParams struct {
	Reason              string `json:"reason,omitempty"`
	TerminationFee      int    `json:"termination_fee,omitempty"`
	LastBillRefundPrice int    `json:"last_bill_refund_price,omitempty"`
	FinalFee            int    `json:"final_fee,omitempty"`
	ServiceEndAt        string `json:"service_end_at,omitempty"`
	CancelDate          string `json:"cancel_date,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// SupervisorOrderSubscriptionPauseParams represents supervisor pause parameters
type SupervisorOrderSubscriptionPauseParams struct {
	Reason           string `json:"reason,omitempty"`
	PausedAt         string `json:"paused_at"`
	ExpectedResumeAt string `json:"expected_resume_at,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// SupervisorOrderSubscriptionResumeParams represents supervisor resume parameters
type SupervisorOrderSubscriptionResumeParams struct {
	Reason string `json:"reason,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// SupervisorOrderSubscriptionChargeParams represents on-demand charge_key payment parameters
// (POST order_subscriptions/charge, supervisor only)
// ⚠️ charge_key is sent only in the body — never in the URL/query (access log exposure).
type SupervisorOrderSubscriptionChargeParams struct {
	ChargeKey    string                 `json:"charge_key"`
	Price        int                    `json:"price"`
	TaxFreePrice int                    `json:"tax_free_price,omitempty"`
	User         map[string]interface{} `json:"user,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// SupervisorOrderSubscriptionChargeRevokeParams represents on-demand charge_key revoke parameters
// (DELETE order_subscriptions/charge, supervisor only)
// After revoking, the key can never be charged again.
type SupervisorOrderSubscriptionChargeRevokeParams struct {
	ChargeKey string                 `json:"charge_key"`
	User      map[string]interface{} `json:"user,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// ============================================
// OrderSubscriptionBill Types
// ============================================

// CommerceOrderSubscriptionBill represents an order subscription bill
type CommerceOrderSubscriptionBill struct {
	OrderSubscriptionBillId string `json:"order_subscription_bill_id,omitempty"`
	OrderSubscriptionId     string `json:"order_subscription_id,omitempty"`
	UserId                  string `json:"user_id,omitempty"`
	UserGroupId             string `json:"user_group_id,omitempty"`

	SubscriptionBillingType int    `json:"subscription_billing_type,omitempty"`
	OrderName               string `json:"order_name,omitempty"`
	PaidWalletId            string `json:"paid_wallet_id,omitempty"`
	ReservedWalletId        string `json:"reserved_wallet_id,omitempty"`

	OrderNumber               string `json:"order_number,omitempty"`
	OrderPreId                string `json:"order_pre_id,omitempty"`
	OrderId                   string `json:"order_id,omitempty"`
	Duration                  int    `json:"duration,omitempty"`
	TotalSubscriptionDuration int    `json:"total_subscription_duration,omitempty"`

	OneUnitPrice        int `json:"one_unit_price,omitempty"`
	OneUnitTaxFreePrice int `json:"one_unit_tax_free_price,omitempty"`
	SetupPrice          int `json:"setup_price,omitempty"`

	Price        int `json:"price,omitempty"`
	TaxFreePrice int `json:"tax_free_price,omitempty"`
	Unit         int `json:"unit,omitempty"`

	PurchasePrice        int `json:"purchase_price,omitempty"`
	PurchaseTaxFreePrice int `json:"purchase_tax_free_price,omitempty"`

	CancelledPrice        int `json:"cancelled_price,omitempty"`
	CancelledTaxFreePrice int `json:"cancelled_tax_free_price,omitempty"`
	CancelledFee          int `json:"cancelled_fee,omitempty"`

	MembershipType int `json:"membership_type,omitempty"`

	AddressId          string `json:"address_id,omitempty"`
	UserAddress        string `json:"user_address,omitempty"`
	Username           string `json:"username,omitempty"`
	UserPhone          string `json:"user_phone,omitempty"`
	UserEmail          string `json:"user_email,omitempty"`
	UserCompanyName    string `json:"user_company_name,omitempty"`
	UserBusinessNumber string `json:"user_business_number,omitempty"`

	ProductIds               []string `json:"product_ids,omitempty"`
	ProductOptionIds         []string `json:"product_option_ids,omitempty"`
	ProductSnapshotIds       []string `json:"product_snapshot_ids,omitempty"`
	ProductOptionSnapshotIds []string `json:"product_option_snapshot_ids,omitempty"`
	ProductType              int      `json:"product_type,omitempty"`
	Quantity                 int      `json:"quantity,omitempty"`

	ReservePaymentAt string `json:"reserve_payment_at,omitempty"`
	PurchasedAt      string `json:"purchased_at,omitempty"`
	RevokedAt        string `json:"revoked_at,omitempty"`
	LastErrorAt      string `json:"last_error_at,omitempty"`

	Status       int    `json:"status,omitempty"`
	CancelStatus int    `json:"cancel_status,omitempty"`
	TestCode     string `json:"test_code,omitempty"`

	ServiceStartAt string `json:"service_start_at,omitempty"`
	ServiceEndAt   string `json:"service_end_at,omitempty"`
}

// OrderSubscriptionBillListParams represents order subscription bill list query parameters
// page/limit default to 1/20 when unset.
type OrderSubscriptionBillListParams struct {
	ListParams
	OrderSubscriptionId string `json:"order_subscription_id,omitempty"`
	Status              []int  `json:"status,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// ============================================
// OrderSubscriptionAdjustment Types
// ============================================

// Constants for subscription adjustment type
const (
	SUBSCRIPTION_ADJUSTMENT_TYPE_PERIOD_DISCOUNT = 1
)

// CommerceOrderSubscriptionAdjustment represents an order subscription adjustment
//
// Three ways to target the billing cycles (widest last) — see OrderSubscriptionAdjustmentModule.Create:
//   - Duration: 5                       → the 5th cycle only
//   - DurationFrom: 3, DurationTo: 7    → cycles 3~7, one record each (5 records)
//   - DurationFrom: 3, IsUnlimited:true → from cycle 3 to the end of the contract
//     (a single record, DurationTo is ignored)
type CommerceOrderSubscriptionAdjustment struct {
	OrderSubscriptionAdjustmentId string `json:"order_subscription_adjustment_id,omitempty"`
	Duration                      int    `json:"duration,omitempty"`
	Price                         int    `json:"price,omitempty"`
	TaxFreePrice                  int    `json:"tax_free_price,omitempty"`
	Name                          string `json:"name,omitempty"`
	Type                          int    `json:"type,omitempty"`
	// DurationFrom is the first cycle of the range (1-based, inclusive)
	DurationFrom int `json:"duration_from,omitempty"`
	// DurationTo is the last cycle of the range (inclusive, ignored when IsUnlimited is true)
	DurationTo int `json:"duration_to,omitempty"`
	// IsUnlimited applies the adjustment from DurationFrom to the end of the contract
	IsUnlimited *bool  `json:"is_unlimited,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// OrderSubscriptionAdjustmentUpdateParams represents subscription adjustment update parameters
// The server replaces the adjustments of the given Duration (defaults to 1 when unset) as a whole —
// pass the full Adjustments array for that duration.
type OrderSubscriptionAdjustmentUpdateParams struct {
	OrderSubscriptionId           string                                `json:"order_subscription_id"`
	OrderSubscriptionAdjustmentId string                                `json:"order_subscription_adjustment_id,omitempty"`
	Duration                      int                                   `json:"duration,omitempty"`
	Adjustments                   []CommerceOrderSubscriptionAdjustment `json:"adjustments,omitempty"`
	Price                         int                                   `json:"price,omitempty"`
	TaxFreePrice                  int                                   `json:"tax_free_price,omitempty"`
	Name                          string                                `json:"name,omitempty"`
	Type                          int                                   `json:"type,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// ============================================
// Category Types
// ============================================

// CommerceCategory represents a category
type CommerceCategory struct {
	CategoryId       string   `json:"category_id,omitempty"`
	SellerId         string   `json:"seller_id,omitempty"`
	ProjectId        string   `json:"project_id,omitempty"`
	Name             string   `json:"name,omitempty"`
	ParentCategoryId string   `json:"parent_category_id,omitempty"`
	ParentCategories []string `json:"parent_categories,omitempty"`
	StatusDisplay    *bool    `json:"status_display,omitempty"`
	StatusBest       *bool    `json:"status_best,omitempty"`
	FilterColor      int      `json:"filter_color,omitempty"`
	FilterSize       int      `json:"filter_size,omitempty"`
	Idx              int      `json:"idx,omitempty"`
	CreatedAt        string   `json:"created_at,omitempty"`
	UpdatedAt        string   `json:"updated_at,omitempty"`
}

// CategoryCreateParams represents category creation parameters
type CategoryCreateParams struct {
	Name             string `json:"name"`
	ParentCategoryId string `json:"parent_category_id,omitempty"`
	StatusDisplay    *bool  `json:"status_display,omitempty"`
	StatusBest       *bool  `json:"status_best,omitempty"`
	FilterColor      int    `json:"filter_color,omitempty"`
	FilterSize       int    `json:"filter_size,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// CategoryUpdateParams represents category update parameters
type CategoryUpdateParams struct {
	CategoryId       string `json:"-"`
	Name             string `json:"name,omitempty"`
	ParentCategoryId string `json:"parent_category_id,omitempty"`
	StatusDisplay    *bool  `json:"status_display,omitempty"`
	StatusBest       *bool  `json:"status_best,omitempty"`
	FilterColor      int    `json:"filter_color,omitempty"`
	FilterSize       int    `json:"filter_size,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// CategoryUpdateBody represents the body sent to PUT categories/{id}
type CategoryUpdateBody struct {
	Name             string `json:"name,omitempty"`
	ParentCategoryId string `json:"parent_category_id,omitempty"`
	StatusDisplay    *bool  `json:"status_display,omitempty"`
	StatusBest       *bool  `json:"status_best,omitempty"`
	FilterColor      int    `json:"filter_color,omitempty"`
	FilterSize       int    `json:"filter_size,omitempty"`
}

// ============================================
// Coupon Types
// ============================================

// CommerceCoupon represents a coupon
type CommerceCoupon struct {
	CouponId          string `json:"coupon_id,omitempty"`
	CouponTemplateId  string `json:"coupon_template_id,omitempty"`
	UserId            string `json:"user_id,omitempty"`
	ProjectId         string `json:"project_id,omitempty"`
	Name              string `json:"name,omitempty"`
	DiscountType      int    `json:"discount_type,omitempty"`
	DiscountValue     int    `json:"discount_value,omitempty"`
	MinOrderAmount    int    `json:"min_order_amount,omitempty"`
	MaxDiscountAmount int    `json:"max_discount_amount,omitempty"`
	Status            int    `json:"status,omitempty"`
	IssuedAt          string `json:"issued_at,omitempty"`
	UsedAt            string `json:"used_at,omitempty"`
	ExpiresAt         string `json:"expires_at,omitempty"`
	CreatedAt         string `json:"created_at,omitempty"`
}

// CouponListParams represents coupon list query parameters
type CouponListParams struct {
	Status string `json:"status,omitempty"`
	Page   int    `json:"page,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// CouponDownloadParams represents coupon download parameters
type CouponDownloadParams struct {
	CouponTemplateId string `json:"coupon_template_id"`
}

// ============================================
// Point Types
// ============================================

// PointBalance represents point balance
type PointBalance struct {
	AvailableBalance int  `json:"available_balance,omitempty"`
	TotalEarned      int  `json:"total_earned,omitempty"`
	TotalUsed        int  `json:"total_used,omitempty"`
	IsNegative       bool `json:"is_negative,omitempty"`
}

// PointTransaction represents a point transaction
type PointTransaction struct {
	TransactionId    string `json:"transaction_id,omitempty"`
	TransactionType  int    `json:"transaction_type,omitempty"`
	Amount           int    `json:"amount,omitempty"`
	BalanceAfter     int    `json:"balance_after,omitempty"`
	Reason           string `json:"reason,omitempty"`
	Type             int    `json:"type,omitempty"`
	OrderId          string `json:"order_id,omitempty"`
	ReviewId         string `json:"review_id,omitempty"`
	EarnedAt         string `json:"earned_at,omitempty"`
	ExpiresAt        string `json:"expires_at,omitempty"`
	Expired          bool   `json:"expired,omitempty"`
	RemainingBalance int    `json:"remaining_balance,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
}

// PointTransactionsResponse represents the response from point transactions
type PointTransactionsResponse struct {
	Transactions []PointTransaction `json:"transactions,omitempty"`
	TotalCount   int                `json:"total_count,omitempty"`
	Page         int                `json:"page,omitempty"`
	Limit        int                `json:"limit,omitempty"`
	TotalPages   int                `json:"total_pages,omitempty"`
}

// PointTransactionsParams represents point transactions query parameters
type PointTransactionsParams struct {
	Page            int `json:"page,omitempty"`
	Limit           int `json:"limit,omitempty"`
	TransactionType int `json:"transaction_type,omitempty"`
}

// ============================================
// Cart Types
// ============================================

// CartItemPayload represents a cart item payload
type CartItemPayload struct {
	ProductId            string `json:"product_id"`
	ProductOptionId      string `json:"product_option_id,omitempty"`
	Quantity             int    `json:"quantity,omitempty"`
	IsSubscription       bool   `json:"is_subscription,omitempty"`
	SubscriptionPeriodId string `json:"subscription_period_id,omitempty"`
}

// ShippingAddressPayload represents a shipping address payload
type ShippingAddressPayload struct {
	Zipcode string `json:"zipcode,omitempty"`
}

// OrderPreviewParams represents order preview parameters
type OrderPreviewParams struct {
	MemberMode      string                  `json:"member_mode,omitempty"`
	CartItems       []CartItemPayload       `json:"cart_items,omitempty"`
	ShippingAddress *ShippingAddressPayload `json:"shipping_address,omitempty"`
	CouponIds       []string                `json:"coupon_ids,omitempty"`
	PointAmount     int                     `json:"point_amount,omitempty"`
	UserGroupId     string                  `json:"user_group_id,omitempty"`
}

// DeliveryGroupItem represents a delivery group item
type DeliveryGroupItem struct {
	CartItemId      string `json:"cart_item_id,omitempty"`
	ProductId       string `json:"product_id"`
	ProductOptionId string `json:"product_option_id,omitempty"`
	ProductName     string `json:"product_name,omitempty"`
	Quantity        int    `json:"quantity"`
	Price           int    `json:"price"`
	Subtotal        int    `json:"subtotal,omitempty"`
}

// DeliveryGroup represents a delivery group
type DeliveryGroup struct {
	GroupKey                 string              `json:"group_key,omitempty"`
	SellerId                 string              `json:"seller_id,omitempty"`
	DeliveryShippingId       string              `json:"delivery_shipping_id,omitempty"`
	DeliveryShippingBundleId string              `json:"delivery_shipping_bundle_id,omitempty"`
	BundleId                 string              `json:"bundle_id,omitempty"`
	Items                    []DeliveryGroupItem `json:"items"`
	TotalPrice               int                 `json:"total_price"`
	TotalQuantity            int                 `json:"total_quantity"`
	DeliveryFee              int                 `json:"delivery_fee"`
	DeliveryExtraFeeJeju     int                 `json:"delivery_extra_fee_jeju,omitempty"`
	DeliveryExtraFeeRemote   int                 `json:"delivery_extra_fee_remote,omitempty"`
	ShippingAvailable        *bool               `json:"shipping_available,omitempty"`
}

// AppliedCouponSnapshot represents an applied coupon snapshot
type AppliedCouponSnapshot struct {
	CouponId             string                 `json:"coupon_id,omitempty"`
	CouponTemplateId     string                 `json:"coupon_template_id,omitempty"`
	Name                 string                 `json:"name,omitempty"`
	DiscountType         int                    `json:"discount_type,omitempty"`
	DiscountValue        int                    `json:"discount_value,omitempty"`
	ActualDiscountAmount int                    `json:"actual_discount_amount,omitempty"`
	Extra                map[string]interface{} `json:"-"`
}

// OrderPreviewSummary represents an order preview summary
type OrderPreviewSummary struct {
	TotalItems            int                     `json:"total_items"`
	TotalQuantity         int                     `json:"total_quantity"`
	TotalProductPrice     int                     `json:"total_product_price"`
	TotalDeliveryFee      int                     `json:"total_delivery_fee"`
	TotalDeliveryExtraFee int                     `json:"total_delivery_extra_fee"`
	CouponDiscountAmount  int                     `json:"coupon_discount_amount"`
	AppliedCoupons        []AppliedCouponSnapshot `json:"applied_coupons"`
	PointUseAmount        int                     `json:"point_use_amount"`
	PointMaxUsable        int                     `json:"point_max_usable"`
	PointBalanceAfter     int                     `json:"point_balance_after"`
	TotalOrderPrice       int                     `json:"total_order_price"`
}

// OrderPreviewUnavailableItem represents an unavailable item in order preview
type OrderPreviewUnavailableItem struct {
	CartItemId  string `json:"cart_item_id,omitempty"`
	ProductId   string `json:"product_id"`
	ProductName string `json:"product_name,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

// OrderPreviewResponse represents the response from order preview
type OrderPreviewResponse struct {
	CartId           string                        `json:"cart_id,omitempty"`
	UserId           string                        `json:"user_id,omitempty"`
	DeliveryGroups   []DeliveryGroup               `json:"delivery_groups"`
	Summary          OrderPreviewSummary           `json:"summary"`
	UnavailableItems []OrderPreviewUnavailableItem `json:"unavailable_items,omitempty"`
}

// ============================================
// OrderSubscriptionRequest Types
// ============================================

// OrderSubscriptionRequest represents an order subscription request
type OrderSubscriptionRequest struct {
	OrderSubscriptionRequestHistoryId string `json:"order_subscription_request_history_id,omitempty"`
	OrderSubscriptionId               string `json:"order_subscription_id,omitempty"`
	ProjectId                         string `json:"project_id,omitempty"`
	UserId                            string `json:"user_id,omitempty"`
	RequestType                       int    `json:"request_type,omitempty"`
	Status                            int    `json:"status,omitempty"`
	Reason                            string `json:"reason,omitempty"`
	RequestedAt                       string `json:"requested_at,omitempty"`
	ProcessedAt                       string `json:"processed_at,omitempty"`
	CreatedAt                         string `json:"created_at,omitempty"`
	UpdatedAt                         string `json:"updated_at,omitempty"`
}

// OrderSubscriptionRequestListParams represents request list query parameters
// page/limit default to 1/20 when unset.
type OrderSubscriptionRequestListParams struct {
	ProjectId           string `json:"project_id,omitempty"`
	OrderSubscriptionId string `json:"order_subscription_id,omitempty"`
	Page                int    `json:"page,omitempty"`
	Limit               int    `json:"limit,omitempty"`
	RequestType         int    `json:"request_type,omitempty"`
	Status              int    `json:"status,omitempty"`
	SAt                 string `json:"s_at,omitempty"`
	EAt                 string `json:"e_at,omitempty"`
	Keyword             string `json:"keyword,omitempty"`
	UserId              string `json:"user_id,omitempty"`
	UserGroupId         string `json:"user_group_id,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// OrderSubscriptionRequestUpdateParams represents request update parameters (supervisor)
// Approval values: "approve" | "reject"
// Termination approval accepts fee overrides (Price/TerminationFee/... — pointer types so
// explicit 0 is sendable); arbitrary additional keys go through Extra.
type OrderSubscriptionRequestUpdateParams struct {
	OrderSubscriptionRequestHistoryId string                 `json:"-"`
	Approval                          string                 `json:"approval"`
	Reason                            string                 `json:"reason,omitempty"`
	Price                             *int                   `json:"price,omitempty"`
	TaxFreePrice                      *int                   `json:"tax_free_price,omitempty"`
	TerminationFee                    *int                   `json:"termination_fee,omitempty"`
	LastBillRefundPrice               *int                   `json:"last_bill_refund_price,omitempty"`
	FinalFee                          *int                   `json:"final_fee,omitempty"`
	ServiceEndAt                      string                 `json:"service_end_at,omitempty"`
	Extra                             map[string]interface{} `json:"-"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// OrderSubscriptionRequestUpdateBody represents the request body sent to PUT order-subscription-requests/{id}
type OrderSubscriptionRequestUpdateBody struct {
	Approval            string                 `json:"approval"`
	Reason              string                 `json:"reason,omitempty"`
	Price               *int                   `json:"price,omitempty"`
	TaxFreePrice        *int                   `json:"tax_free_price,omitempty"`
	TerminationFee      *int                   `json:"termination_fee,omitempty"`
	LastBillRefundPrice *int                   `json:"last_bill_refund_price,omitempty"`
	FinalFee            *int                   `json:"final_fee,omitempty"`
	ServiceEndAt        string                 `json:"service_end_at,omitempty"`
	Extra               map[string]interface{} `json:"-"`
}

// MarshalJSON merges Extra into the serialized body — with the plain struct tags Extra was
// silently dropped (json:"-"). Typed fields win over same-named Extra keys; nil Extra values
// are not sent (compact semantics).
func (b OrderSubscriptionRequestUpdateBody) MarshalJSON() ([]byte, error) {
	body := map[string]interface{}{}
	for key, value := range b.Extra {
		if value != nil {
			body[key] = value
		}
	}
	body["approval"] = b.Approval
	if b.Reason != "" {
		body["reason"] = b.Reason
	}
	if b.Price != nil {
		body["price"] = *b.Price
	}
	if b.TaxFreePrice != nil {
		body["tax_free_price"] = *b.TaxFreePrice
	}
	if b.TerminationFee != nil {
		body["termination_fee"] = *b.TerminationFee
	}
	if b.LastBillRefundPrice != nil {
		body["last_bill_refund_price"] = *b.LastBillRefundPrice
	}
	if b.FinalFee != nil {
		body["final_fee"] = *b.FinalFee
	}
	if b.ServiceEndAt != "" {
		body["service_end_at"] = b.ServiceEndAt
	}
	return json.Marshal(body)
}

// ============================================
// MallSetting Types
// ============================================

// MallSettingUpdateParams represents mall setting update parameters (PUT mall-setting)
// The body is flat and only provided (non-nil / non-zero) values are sent —
// pointer types (*bool / *int) allow sending explicit false / 0.
// supervisor scope only.
type MallSettingUpdateParams struct {
	// 위젯
	NormalWidgetKey       string `json:"normal_widget_key,omitempty"`
	SubscriptionWidgetKey string `json:"subscription_widget_key,omitempty"`

	// 사업자 정보
	SellerName           string `json:"seller_name,omitempty"`
	SellerNameEn         string `json:"seller_name_en,omitempty"`
	BizEmail             string `json:"biz_email,omitempty"`
	BizTel               string `json:"biz_tel,omitempty"`
	BizFax               string `json:"biz_fax,omitempty"`
	RegistrationNo       string `json:"registration_no,omitempty"`
	CorpRegNo            string `json:"corp_reg_no,omitempty"`
	MailOrderSalesNumber string `json:"mail_order_sales_number,omitempty"`
	OwnerName            string `json:"owner_name,omitempty"`
	Zip                  string `json:"zip,omitempty"`
	Addr1                string `json:"addr_1,omitempty"`
	Addr2                string `json:"addr_2,omitempty"`
	PrivacyName          string `json:"privacy_name,omitempty"`
	PrivacyEmail         string `json:"privacy_email,omitempty"`

	// 몰 기본 정보
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Status       *int   `json:"status,omitempty"`
	InvoiceTitle string `json:"invoice_title,omitempty"`

	// 브랜딩
	UseLogo      *bool  `json:"use_logo,omitempty"`
	Logo         string `json:"logo,omitempty"`
	UseFavicon   *bool  `json:"use_favicon,omitempty"`
	Favicon      string `json:"favicon,omitempty"`
	UseOpenGraph *bool  `json:"use_open_graph,omitempty"`
	OgImage      string `json:"og_image,omitempty"`
	UseSignature *bool  `json:"use_signature,omitempty"`
	Signature    string `json:"signature,omitempty"`

	// 고객센터 운영시간
	UseOperationTime                   *bool  `json:"use_operation_time,omitempty"`
	CustomerServiceCenterOperationTime string `json:"customer_service_center_operation_time,omitempty"`
	RestStartHour                      *int   `json:"rest_start_hour,omitempty"`
	RestStartMinute                    *int   `json:"rest_start_minute,omitempty"`
	RestEndHour                        *int   `json:"rest_end_hour,omitempty"`
	RestEndMinute                      *int   `json:"rest_end_minute,omitempty"`
	// 휴무일 (요일 코드 배열 또는 서버 정의 문자열)
	RestDay        interface{} `json:"rest_day,omitempty"`
	HostingService string      `json:"hosting_service,omitempty"`

	// 주문/연령 정책
	UseNonMemberOrder       *bool `json:"use_non_member_order,omitempty"`
	UseAgeAccept19          *bool `json:"use_age_accept_19,omitempty"`
	UseAgeAccept14          *bool `json:"use_age_accept_14,omitempty"`
	UseAgeAcceptParentName  *bool `json:"use_age_accept_parent_name,omitempty"`
	UseAgeAcceptParentBirth *bool `json:"use_age_accept_parent_birth,omitempty"`
	UseAgeAcceptParentEmail *bool `json:"use_age_accept_parent_email,omitempty"`

	// 회원가입 수집 항목
	UseMembershipCollectPhone       *bool `json:"use_membership_collect_phone,omitempty"`
	UseMembershipCollectTel         *bool `json:"use_membership_collect_tel,omitempty"`
	UseMembershipCollectEmail       *bool `json:"use_membership_collect_email,omitempty"`
	UseMembershipCollectAddress     *bool `json:"use_membership_collect_address,omitempty"`
	UseMembershipCollectBank        *bool `json:"use_membership_collect_bank,omitempty"`
	UseMembershipCollectBirth       *bool `json:"use_membership_collect_birth,omitempty"`
	UseMembershipCollectGender      *bool `json:"use_membership_collect_gender,omitempty"`
	UseMembershipCollectInterest    *bool `json:"use_membership_collect_interest,omitempty"`
	MembershipCollectInterestNumber *int  `json:"membership_collect_interest_number,omitempty"`
	UseMembershipCollectCustoms     *bool `json:"use_membership_collect_customs,omitempty"`
	UseMembershipCollectNickname    *bool `json:"use_membership_collect_nickname,omitempty"`
	UseMembershipCollectRecommendId *bool `json:"use_membership_collect_recommend_id,omitempty"`
	RecommendIdPointTo              *int  `json:"recommend_id_point_to,omitempty"`
	RecommendIdPointFrom            *int  `json:"recommend_id_point_from,omitempty"`
	UseMembershipCollectBusiness    *bool `json:"use_membership_collect_business,omitempty"`
	UseMembershipCollectRegister    *bool `json:"use_membership_collect_register,omitempty"`
	MembershipOnlyBusiness          *bool `json:"membership_only_business,omitempty"`

	// 기업(그룹) 회원
	UseCorporateDepartment     *bool `json:"use_corporate_department,omitempty"`
	SubGroupType               *int  `json:"sub_group_type,omitempty"`
	UseCorporateSignupApproval *bool `json:"use_corporate_signup_approval,omitempty"`
	// 기업 회원 허용 이메일 도메인 목록
	CorporateEmailDomains   interface{} `json:"corporate_email_domains,omitempty"`
	UseCorporateAutoApprove *bool       `json:"use_corporate_auto_approve,omitempty"`
	UseCorporateInviteOnly  *bool       `json:"use_corporate_invite_only,omitempty"`

	// 회원 정보 노출 항목
	UseMemberInfoPhone    *bool `json:"use_member_info_phone,omitempty"`
	UseMemberInfoTel      *bool `json:"use_member_info_tel,omitempty"`
	UseMemberInfoEmail    *bool `json:"use_member_info_email,omitempty"`
	UseMemberInfoAddress  *bool `json:"use_member_info_address,omitempty"`
	UseMemberInfoBank     *bool `json:"use_member_info_bank,omitempty"`
	UseMemberInfoBirth    *bool `json:"use_member_info_birth,omitempty"`
	UseMemberInfoGender   *bool `json:"use_member_info_gender,omitempty"`
	UseMemberInfoCustoms  *bool `json:"use_member_info_customs,omitempty"`
	UseMemberInfoNickname *bool `json:"use_member_info_nickname,omitempty"`
	UseMemberInfoRegister *bool `json:"use_member_info_register,omitempty"`

	// 주문자 수집 항목
	OrdererCollectPhone *bool `json:"orderer_collect_phone,omitempty"`
	OrdererCollectTel   *bool `json:"orderer_collect_tel,omitempty"`
	OrdererCollectEmail *bool `json:"orderer_collect_email,omitempty"`

	// 주문/취소 정책
	OrderPrefix    string `json:"order_prefix,omitempty"`
	UseOrderCancel *bool  `json:"use_order_cancel,omitempty"`
	// 취소 승인 사용 여부 (서버 필드명 오타 그대로 유지)
	UseOderCancelApproval *bool `json:"use_oder_cancel_approval,omitempty"`
	// 취소 사유 목록
	OrderCancelReasons            interface{} `json:"order_cancel_reasons,omitempty"`
	OrderCancelReasonRequiredType *int        `json:"order_cancel_reason_required_type,omitempty"`
	OrderCancelRequestMessage     string      `json:"order_cancel_request_message,omitempty"`
	OrderCancelDoneMessage        string      `json:"order_cancel_done_message,omitempty"`

	// 회원 가입/인증 방식
	UseGeneralMembership          *bool `json:"use_general_membership,omitempty"`
	GeneralMembershipDuplication  *int  `json:"general_membership_duplication,omitempty"`
	UseCertification              *bool `json:"use_certification,omitempty"`
	CertificationType             *int  `json:"certification_type,omitempty"`
	GeneralMembershipIdType       *int  `json:"general_membership_id_type,omitempty"`
	UseMembershipDuplicationEmail *bool `json:"use_membership_duplication_email,omitempty"`
	UseMembershipDuplicationPhone *bool `json:"use_membership_duplication_phone,omitempty"`
	UseSocialMembership           *bool `json:"use_social_membership,omitempty"`
	// 사용 소셜 로그인 타입
	SocialMembershipType interface{} `json:"social_membership_type,omitempty"`

	// 적립금
	UsePoint            *bool  `json:"use_point,omitempty"`
	UsePointTransaction *bool  `json:"use_point_transaction,omitempty"`
	PointDisplayName    string `json:"point_display_name,omitempty"`
	PointMinBalance     *int   `json:"point_min_balance,omitempty"`
	// 적립 제외 조건
	PointNotCondition interface{} `json:"point_not_condition,omitempty"`
	// 적립 조건
	PointCondition           interface{} `json:"point_condition,omitempty"`
	UsePointMaxRate          *bool       `json:"use_point_max_rate,omitempty"`
	PointMaxRate             *int        `json:"point_max_rate,omitempty"`
	UsePointMaxAmount        *bool       `json:"use_point_max_amount,omitempty"`
	PointMaxAmount           *int        `json:"point_max_amount,omitempty"`
	PointRate                *int        `json:"point_rate,omitempty"`
	PointCalcType1           *int        `json:"point_calc_type1,omitempty"`
	PointCalcType2           *int        `json:"point_calc_type2,omitempty"`
	UsePointAdvanceDiscount  *bool       `json:"use_point_advance_discount,omitempty"`
	PointAdvanceDiscountRate *int        `json:"point_advance_discount_rate,omitempty"`
	UsePointExpire           *bool       `json:"use_point_expire,omitempty"`
	PointExpireType          *int        `json:"point_expire_type,omitempty"`
	PointIssueEventType      *int        `json:"point_issue_event_type,omitempty"`
	PointIssueDelayDays      *int        `json:"point_issue_delay_days,omitempty"`

	// 오픈마켓 / 상품
	UseOpenMarket                 *bool  `json:"use_open_market,omitempty"`
	UseProductApproval            *bool  `json:"use_product_approval,omitempty"`
	UseProductReview              *bool  `json:"use_product_review,omitempty"`
	UseProductReviewPoint         *bool  `json:"use_product_review_point,omitempty"`
	ProductReviewPoint            *int   `json:"product_review_point,omitempty"`
	ProductReviewPhotoPoint       *int   `json:"product_review_photo_point,omitempty"`
	UseProductReviewAnswer        *bool  `json:"use_product_review_answer,omitempty"`
	UseProductReviewAutoAnswer    *bool  `json:"use_product_review_auto_answer,omitempty"`
	ProductReviewAutoAnswerMinute *int   `json:"product_review_auto_answer_minute,omitempty"`
	ProductReviewAutoAnswerText   string `json:"product_review_auto_answer_text,omitempty"`
	UseProductQna                 *bool  `json:"use_product_qna,omitempty"`
	ProductQnaMemberAuth          *int   `json:"product_qna_member_auth,omitempty"`
	UseProductQnaAnswerOption     *bool  `json:"use_product_qna_answer_option,omitempty"`

	// 게시판 / 상담
	UseNotice       *bool  `json:"use_notice,omitempty"`
	UseQna          *bool  `json:"use_qna,omitempty"`
	UseFaq          *bool  `json:"use_faq,omitempty"`
	UseChatSupport  *bool  `json:"use_chat_support,omitempty"`
	ChatSupportType *int   `json:"chat_support_type,omitempty"`
	ChatSupportKey  string `json:"chat_support_key,omitempty"`

	// 휴면 / 탈퇴
	UseDormant                     *bool  `json:"use_dormant,omitempty"`
	DormantYear                    *int   `json:"dormant_year,omitempty"`
	DormantRestore                 *int   `json:"dormant_restore,omitempty"`
	UseWithdrawal                  *bool  `json:"use_withdrawal,omitempty"`
	UseWithdrawalGuideMessage      *bool  `json:"use_withdrawal_guide_message,omitempty"`
	UseWithdrawalGuideMessageAfter *bool  `json:"use_withdrawal_guide_message_after,omitempty"`
	WithdrawalGuideMessageAfter    string `json:"withdrawal_guide_message_after,omitempty"`
	UseWithdrawalAuto              *bool  `json:"use_withdrawal_auto,omitempty"`
	WithdrawalAutoYear             *int   `json:"withdrawal_auto_year,omitempty"`

	// 정기구독 정산
	UseSubscriptionAggregateTransaction *bool `json:"use_subscription_aggregate_transaction,omitempty"`
	SubscriptionMonthDay                *int  `json:"subscription_month_day,omitempty"`
	SubscriptionWeekDay                 *int  `json:"subscription_week_day,omitempty"`

	// 구매 한도
	UseLimit           *bool `json:"use_limit,omitempty"`
	LimitMonthPurchase *int  `json:"limit_month_purchase,omitempty"`
	LimitWeekPurchase  *int  `json:"limit_week_purchase,omitempty"`
	UseLimitPayment    *bool `json:"use_limit_payment,omitempty"`
	UseLimitMessage    *bool `json:"use_limit_message,omitempty"`

	// 약관
	TermsOfService        string `json:"terms_of_service,omitempty"`
	TermsOfPrivacyPolicy  string `json:"terms_of_privacy_policy,omitempty"`
	TermsOfPrivacyCollect string `json:"terms_of_privacy_collect,omitempty"`
	TermsOfPrivacyThird   string `json:"terms_of_privacy_third,omitempty"`

	// 결제 / 노출
	PaymentTimeout         *int   `json:"payment_timeout,omitempty"`
	ProductSortType        *int   `json:"product_sort_type,omitempty"`
	MallThemeType          *int   `json:"mall_theme_type,omitempty"`
	CatalogDisplayType     *int   `json:"catalog_display_type,omitempty"`
	CatalogHeadline        string `json:"catalog_headline,omitempty"`
	CatalogBgColor         string `json:"catalog_bg_color,omitempty"`
	CatalogViewTypePc      *int   `json:"catalog_view_type_pc,omitempty"`
	CatalogViewTypeMobile  *int   `json:"catalog_view_type_mobile,omitempty"`
	CatalogProductSortType *int   `json:"catalog_product_sort_type,omitempty"`

	// 장바구니 / 위시리스트
	UseCart             *bool `json:"use_cart,omitempty"`
	CartStoragePeriod   *int  `json:"cart_storage_period,omitempty"`
	CartMaxLimit        *int  `json:"cart_max_limit,omitempty"`
	CartAddAction       *int  `json:"cart_add_action,omitempty"`
	CartDirectPurchase  *bool `json:"cart_direct_purchase,omitempty"`
	CartOptionChange    *bool `json:"cart_option_change,omitempty"`
	CartDiscountDisplay *bool `json:"cart_discount_display,omitempty"`
	UseWishlist         *bool `json:"use_wishlist,omitempty"`
	WishlistMaxLimit    *int  `json:"wishlist_max_limit,omitempty"`
	CartWishlistDisplay *bool `json:"cart_wishlist_display,omitempty"`
}

// ============================================
// Webhook Types
// ============================================

// SendTestWebhookParams represents test webhook send parameters (POST webhook/test)
type SendTestWebhookParams struct {
	// 웹훅 본문 Content-Type (미지정시 서버 기본값; explicit 0 must be sendable, so pointer type)
	HeaderContentType *int `json:"header_content_type,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header (auto-generated when empty)
	IdempotencyKey string `json:"-"`
}

// ============================================
// Alimtalk Types
// ============================================

// AlimtalkMessageListParams represents alimtalk send-history query parameters
// (GET /v1/alimtalk/messages)
//
// ⚠️ 유료 알림톡만 조회된다 (무료 커머스 알림톡은 포함되지 않는다).
// ⚠️ 기간 기본값은 최근 30일이고 최대 조회 폭은 92일이다 — 초과분은 거부하지 않고
//
//	시작일을 당겨 잘라낸다. 실제 적용된 구간은 응답의 period 로 확인한다.
type AlimtalkMessageListParams struct {
	TemplateCode string `json:"template_code,omitempty"`
	// requested · success · failed · canceled
	Status string `json:"status,omitempty"`
	// 발송 시 넘긴 멱등키
	RefId string `json:"ref_id,omitempty"`
	// 수신번호 (하이픈 무관, 정확 매칭)
	To   string `json:"to,omitempty"`
	SAt  string `json:"s_at,omitempty"`
	EAt  string `json:"e_at,omitempty"`
	Page int    `json:"page,omitempty"`
	// 서버 기본 20, 최대 100
	Limit int `json:"limit,omitempty"`
}

// AlimtalkOfficialListParams represents official template catalog search parameters
// (GET /v1/alimtalk/official)
type AlimtalkOfficialListParams struct {
	// 본문·이름·분류를 부분일치(대소문자 무시)로 훑는다.
	// 서버는 q 를 먼저 보고 없으면 keyword 를 본다 — 정본 키인 q 로 보낸다.
	Keyword  string `json:"q,omitempty"`
	Category string `json:"category,omitempty"`
	// BA(기본형) · EX(부가정보형)만 존재한다 — 그룹 템플릿이라 AD/MI 는 쓸 수 없다.
	MsgType string `json:"msg_type,omitempty"`
	Page    int    `json:"page,omitempty"`
	// 서버 기본 20, 최대 100 으로 clamp
	Per int `json:"per,omitempty"`
	// 주면 그 채널의 변수 예문 사전으로 variable_examples 를 채워 준다(표시용)
	KspId string `json:"ksp_id,omitempty"`
}

// AlimtalkOfficialRecommendParams represents official template recommendation parameters
// (POST /v1/alimtalk/official/recommend)
type AlimtalkOfficialRecommendParams struct {
	Text     string `json:"text"`
	Category string `json:"category,omitempty"`
	// 서버 기본 5
	Limit int    `json:"limit,omitempty"`
	KspId string `json:"ksp_id,omitempty"`
}

// AlimtalkOptoutListParams represents optout list query parameters
// (GET /v1/alimtalk/optouts)
// phone 은 숫자만 남겨 **부분일치**로 찾는다(정확 매칭이 아니다). 50건 단위로 페이징된다.
type AlimtalkOptoutListParams struct {
	Phone string `json:"phone,omitempty"`
	Page  int    `json:"page,omitempty"`
}

// AlimtalkOptoutCreateParams represents optout registration parameters
// (POST /v1/alimtalk/optouts)
// 내 프로젝트 스코프로 등록된다(source: api). 같은 번호를 다시 등록해도 멱등이다.
type AlimtalkOptoutCreateParams struct {
	Phone  string `json:"phone"`
	Reason string `json:"reason,omitempty"`
}

// AlimtalkOptoutCheckParams represents optout pre-check parameters
// (POST /v1/alimtalk/optouts/check)
// 단건(Phone)·다건(Phones) 모두 받는다.
// ⚠️ 1회 최대 1,000건이고 넘으면 -48 이다 (중복은 서버가 제거).
type AlimtalkOptoutCheckParams struct {
	Phones []string `json:"phones,omitempty"`
	Phone  string   `json:"phone,omitempty"`
}

// AlimtalkSendParams represents single alimtalk send parameters (POST /v1/alimtalk/send)
//
// ⚠️ 실제로 카카오톡이 발송되고 과금된다. 샌드박스가 없다.
type AlimtalkSendParams struct {
	TemplateCode string `json:"template_code"`
	To           string `json:"to"`
	// { "company_name": "부트페이몰", "user_name": "홍길동" } 형태의 치환값.
	// 템플릿 응답의 required_variables 를 모두 채워야 한다 — 하나라도 비면 3017 로 거부된다.
	Variables map[string]interface{} `json:"variables,omitempty"`
	// 가맹점 발송 식별자 — **멱등 키**로 쓰인다
	RefId string `json:"ref_id,omitempty"`
	// 알림톡 실패 시 문자(LMS) 대체발송 여부.
	// ⚠️ **미지정(nil)과 false 는 다르다** — nil 이면 프로젝트 기본값을 따르고, false 는 명시적으로 끈다.
	// 그래서 pointer type 이다.
	Fallback *bool `json:"fallback,omitempty"`
	// 예약 발송 시각(ISO8601). 생략하면 즉시 발송
	ReservedAt string `json:"reserved_at,omitempty"`
	// 채널 공개키. 생략하면 프로젝트 연동 채널로 해석하며, 연동 채널이 둘 이상일 때만 필수다.
	// (ksp_id 는 내부 문서 id 라 발송 API 에 쓰지 않는다)
	SenderKey string `json:"sender_key,omitempty"`
	UserId    string `json:"user_id,omitempty"`
}

// AlimtalkSendRecipient represents one recipient of a bulk alimtalk send
type AlimtalkSendRecipient struct {
	To        string                 `json:"to"`
	RefId     string                 `json:"ref_id,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// AlimtalkSendBulkParams represents bulk alimtalk send parameters (POST /v1/alimtalk/send/bulk)
//
// ⚠️ 수신자 수만큼 실제 발송되고 과금된다.
//   - 쿼터를 넘으면 요청 시점에 **전체 거부**된다(3022) — 일부만 나가지 않는다.
//   - 수신거부 번호는 skipped 이며 과금되지 않고 발송 기록도 만들지 않는다.
//   - Fallback 은 요청 단위로 한 번만 판정한다 — 발신번호가 없으면 요청 전체가 3030 으로 거부된다.
type AlimtalkSendBulkParams struct {
	TemplateCode string                  `json:"template_code"`
	Recipients   []AlimtalkSendRecipient `json:"recipients"`
	// 미지정(nil)과 false 는 다르다 — AlimtalkSendParams.Fallback 주석 참고
	Fallback   *bool  `json:"fallback,omitempty"`
	ReservedAt string `json:"reserved_at,omitempty"`
	SenderKey  string `json:"sender_key,omitempty"`
	UserId     string `json:"user_id,omitempty"`
}

// AlimtalkSenderOtpParams represents channel-admin OTP request parameters
// (POST /v1/alimtalk/senders/otp)
// ⚠️ 실제로 문자가 나간다.
type AlimtalkSenderOtpParams struct {
	YellowId string `json:"yellow_id"`
	Phone    string `json:"phone"`
}

// AlimtalkSenderCreateParams represents sender profile registration parameters
// (POST /v1/alimtalk/senders)
// ⚠️ 카카오에 발신프로필이 실제 등록된다. 같은 yellow_id 를 다시 등록하면 기존 프로필을 재사용한다(dedup).
type AlimtalkSenderCreateParams struct {
	Otp          string `json:"otp"`
	YellowId     string `json:"yellow_id"`
	Phone        string `json:"phone"`
	CategoryCode string `json:"category_code"`
}

// AlimtalkTemplateListParams represents own-template list query parameters
// (GET /v1/alimtalk/templates)
// ⚠️ 페이지네이션이 없다 — 필터에 걸린 템플릿을 한 번에 모두 돌려준다.
type AlimtalkTemplateListParams struct {
	// 검수상태 필터 — 1 REG(등록) / 2 REQ(검수요청) / 3 APR(승인) / 4 KRR(등록거절) / 5 REJ(승인반려).
	// 숫자문자열·벤더 문자열('APR' 등)을 모두 받는다. 해석 못 하는 값은 필터 없음으로 떨어진다.
	Ins string `json:"ins,omitempty"`
	// latest(기본) · oldest · code
	Sort string `json:"sort,omitempty"`
	// 코드·이름·본문·분류 부분일치
	Keyword string `json:"keyword,omitempty"`
}

// AlimtalkTemplateButton represents a template button
type AlimtalkTemplateButton struct {
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	UrlMobile     string `json:"url_mobile,omitempty"`
	UrlPc         string `json:"url_pc,omitempty"`
	SchemeIos     string `json:"scheme_ios,omitempty"`
	SchemeAndroid string `json:"scheme_android,omitempty"`
	// 서버가 읽는 그 밖의 키는 Extra 로 넘긴다
	Extra map[string]interface{} `json:"-"`
}

// MarshalJSON merges Extra into the serialized button — typed fields win over
// same-named Extra keys; nil Extra values are not sent (compact semantics).
func (b AlimtalkTemplateButton) MarshalJSON() ([]byte, error) {
	type alias AlimtalkTemplateButton
	return mergeExtraJSON(alias(b), b.Extra)
}

// AlimtalkTemplateCreateParams represents own-template creation parameters
// (POST /v1/alimtalk/templates)
//
// ⚠️ Register 를 false 로 주지 않으면 대행사·카카오에 **실제 등록**된다(되돌리려면 삭제해야 한다).
// ⚠️ 본문 변수는 `#{변수명}` 형식이고 템플릿 전체에서 최대 40개다.
//
// EmphasizeType: NONE · TEXT(강조표기형) · IMAGE(이미지형) · ITEM_LIST(아이템리스트형)
//   - TEXT 는 EmphasizeTitle·EmphasizeSubtitle 둘 다 필수(각 50자·40자)
//   - IMAGE 는 이미지 필수 — AlimtalkTemplate.Image 로 올린 URL 을 StorageImageUrl 로 넘긴다
//   - ITEM_LIST 는 TemplateItem.list(2~10개) 필수 + TemplateHeader·ItemHighlight·이미지 중 하나 이상
//
// MsgType: BA(기본형) · EX(부가정보형, TemplateExtra 필수) · AD(채널추가형) · MI(복합형)
//   - AD·MI 는 채널추가(AC) 버튼이 필수다
type AlimtalkTemplateCreateParams struct {
	KspId string `json:"ksp_id"`
	// ⚠️ 미지정(nil)이면 즉시 등록된다 — 초안만 만들려면 명시적으로 false 를 넘긴다.
	// 그래서 pointer type 이다.
	Register          *bool                    `json:"register,omitempty"`
	Name              string                   `json:"name,omitempty"`
	Content           string                   `json:"content,omitempty"`
	Buttons           []AlimtalkTemplateButton `json:"buttons,omitempty"`
	MsgType           string                   `json:"msg_type,omitempty"`
	EmphasizeType     string                   `json:"emphasize_type,omitempty"`
	EmphasizeTitle    string                   `json:"emphasize_title,omitempty"`
	EmphasizeSubtitle string                   `json:"emphasize_subtitle,omitempty"`
	TemplateExtra     string                   `json:"template_extra,omitempty"`
	TemplateHeader    string                   `json:"template_header,omitempty"`
	ItemHighlight     map[string]interface{}   `json:"item_highlight,omitempty"`
	TemplateItem      map[string]interface{}   `json:"template_item,omitempty"`
	ImageUrl          string                   `json:"image_url,omitempty"`
	StorageImageUrl   string                   `json:"storage_image_url,omitempty"`
	SecurityFlag      *bool                    `json:"security_flag,omitempty"`
	Category          string                   `json:"category,omitempty"`
	Tags              []string                 `json:"tags,omitempty"`
	// 변수 예문(표시용). 주면 **모든 변수에 예문이 있어야** 한다(없으면 3017).
	Examples     map[string]interface{} `json:"examples,omitempty"`
	TemplateCode string                 `json:"template_code,omitempty"`
	// Ruby `**attrs` 대응 — 타입 필드에 없는 키를 그대로 실어 보낸다
	Extra map[string]interface{} `json:"-"`
}

// MarshalJSON merges Extra into the serialized body — typed fields win over
// same-named Extra keys; nil Extra values are not sent (compact semantics).
func (p AlimtalkTemplateCreateParams) MarshalJSON() ([]byte, error) {
	type alias AlimtalkTemplateCreateParams
	return mergeExtraJSON(alias(p), p.Extra)
}

// AlimtalkTemplateUpdateParams represents own-template update parameters
// (PUT /v1/alimtalk/templates/{template_id})
//
// ⚠️ **부분 수정이 아니다.** 보내지 않은 필드는 서버에서 nil 로 덮어써지므로 항상 전체 필드를 보낸다.
// ⚠️ 등록된 템플릿을 수정하면 벤더에도 수정 요청이 나간다. 수정 가능 상태는
//
//	초안 / REG(등록) / REJ(승인반려) / KRR(등록거절) 뿐이다 — APR·REQ 는 거부된다.
//
// storage_image_url 을 **빈 값으로** 보내면 이미지 삭제로 처리되어 벤더에도 전달된다.
// omitempty 때문에 빈 문자열은 실리지 않으므로, 삭제 의도는
// `Extra: map[string]interface{}{"storage_image_url": ""}` 로 명시한다.
type AlimtalkTemplateUpdateParams struct {
	Name              string                   `json:"name,omitempty"`
	Content           string                   `json:"content,omitempty"`
	Buttons           []AlimtalkTemplateButton `json:"buttons,omitempty"`
	MsgType           string                   `json:"msg_type,omitempty"`
	EmphasizeType     string                   `json:"emphasize_type,omitempty"`
	EmphasizeTitle    string                   `json:"emphasize_title,omitempty"`
	EmphasizeSubtitle string                   `json:"emphasize_subtitle,omitempty"`
	TemplateExtra     string                   `json:"template_extra,omitempty"`
	TemplateHeader    string                   `json:"template_header,omitempty"`
	ItemHighlight     map[string]interface{}   `json:"item_highlight,omitempty"`
	TemplateItem      map[string]interface{}   `json:"template_item,omitempty"`
	ImageUrl          string                   `json:"image_url,omitempty"`
	StorageImageUrl   string                   `json:"storage_image_url,omitempty"`
	SecurityFlag      *bool                    `json:"security_flag,omitempty"`
	Category          string                   `json:"category,omitempty"`
	Tags              []string                 `json:"tags,omitempty"`
	Examples          map[string]interface{}   `json:"examples,omitempty"`
	TemplateCode      string                   `json:"template_code,omitempty"`
	// Ruby `**attrs` 대응 — 타입 필드에 없는 키(및 빈 값 명시 전송)를 그대로 실어 보낸다
	Extra map[string]interface{} `json:"-"`
}

// MarshalJSON merges Extra into the serialized body — typed fields win over
// same-named Extra keys; nil Extra values are not sent (compact semantics).
func (p AlimtalkTemplateUpdateParams) MarshalJSON() ([]byte, error) {
	type alias AlimtalkTemplateUpdateParams
	return mergeExtraJSON(alias(p), p.Extra)
}

// AlimtalkTemplateExportParams represents template export parameters
// (GET /v1/alimtalk/templates/export)
// ⚠️ 1회 5,000건을 넘으면 3031 로 거부되므로 채널·상태 필터로 좁힌다.
type AlimtalkTemplateExportParams struct {
	// json(SDK 기본) · csv. 서버 기본은 csv 지만 csv 본문은 JSON 이 아니라서
	// 공용 파서를 통과하지 못한다 — SDK 는 json 을 기본으로 두고, csv 는 원문 그대로 돌려준다.
	Format string `json:"format,omitempty"`
	// private(기본, 내 채널 자체 템플릿) · official(공식 카탈로그) · all
	Scope          string `json:"scope,omitempty"`
	KspId          string `json:"ksp_id,omitempty"`
	Status         string `json:"status,omitempty"`
	IncludeContent *bool  `json:"include_content,omitempty"`
}

// AlimtalkWebhookUpdateParams represents alimtalk webhook configuration parameters
// (PUT /v1/alimtalk/webhook)
//
// ⚠️ 주문·구독 통합 웹훅과 완전히 별개다 (Webhook.SendTest 는 주문 웹훅용이다).
type AlimtalkWebhookUpdateParams struct {
	// **https 만** 허용한다(아니면 3028). 최초 저장 시 서명 시크릿이 자동 발급된다.
	Url string `json:"url,omitempty"`
	// 구독할 이벤트 코드. 목록에 없는 값은 저장 시 조용히 버려진다(유령 구독 방지).
	//   300 발송 접수(기본 미구독) / 301 전달 성공 / 302 전달 실패 / 303 예약 취소 /
	//   304 문자(LMS) 대체발송 전환 / 310 검수 승인 / 311 검수 반려 / 320 수신거부 등록(기본 미구독)
	// 비우면 기본 구독셋(301·302·303·304·310·311)이 적용된다.
	Events []int `json:"events,omitempty"`
	// explicit false 를 보내야 하므로 pointer type
	Enabled *bool `json:"enabled,omitempty"`
}

// AlimtalkWebhookDeliveriesParams represents webhook delivery history query parameters
// (GET /v1/alimtalk/webhook/deliveries)
type AlimtalkWebhookDeliveriesParams struct {
	Page int `json:"page,omitempty"`
	// 서버 기본 20, 최대 100
	Limit int `json:"limit,omitempty"`
}

// mergeExtraJSON serializes value and merges extra into the resulting object.
// Typed fields win over same-named extra keys, and nil extra values are dropped
// (Ruby `.merge(attrs).compact` semantics).
func mergeExtraJSON(value interface{}, extra map[string]interface{}) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if len(extra) == 0 {
		return raw, nil
	}
	body := map[string]interface{}{}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	for key, item := range extra {
		if item == nil {
			continue
		}
		if _, exists := body[key]; exists {
			continue
		}
		body[key] = item
	}
	return json.Marshal(body)
}
