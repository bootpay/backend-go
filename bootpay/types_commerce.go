package bootpay

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
type UserListParams struct {
	ListParams
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

	Phone                string                 `json:"phone,omitempty"`
	Email                string                 `json:"email,omitempty"`
	Zipcode              string                 `json:"zipcode,omitempty"`
	Address              string                 `json:"address,omitempty"`
	AddressDetail        string                 `json:"address_detail,omitempty"`
	CorporateExtension   map[string]interface{} `json:"corporate_extension,omitempty"`
	AuthBank             bool                   `json:"auth_bank,omitempty"`

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
type UserGroupLimitParams struct {
	UserGroupId     string `json:"user_group_id"`
	UseLimit        bool   `json:"use_limit,omitempty"`
	PurchaseLimit   int    `json:"purchase_limit,omitempty"`
	SubscribedLimit int    `json:"subscribed_limit,omitempty"`
	LimitMessage    string `json:"limit_message,omitempty"`
}

// UserGroupAggregateTransactionParams represents aggregate transaction parameters
type UserGroupAggregateTransactionParams struct {
	UserGroupId                         string `json:"user_group_id"`
	UseSubscriptionAggregateTransaction bool   `json:"use_subscription_aggregate_transaction,omitempty"`
	SubscriptionMonthDay                int    `json:"subscription_month_day,omitempty"`
	SubscriptionWeekDay                 int    `json:"subscription_week_day,omitempty"`
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

	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Images       []string `json:"images,omitempty"`
	Type         int      `json:"type,omitempty"`
	TaxType      int      `json:"tax_type,omitempty"`
	UseStock     bool     `json:"use_stock,omitempty"`
	Stock        int      `json:"stock,omitempty"`
	UseOptionStock bool   `json:"use_option_stock,omitempty"`
	UseStockSafe bool     `json:"use_stock_safe,omitempty"`
	StockSafe    int      `json:"stock_safe,omitempty"`

	DisplayPrice      int  `json:"display_price,omitempty"`
	TaxFreePrice      int  `json:"tax_free_price,omitempty"`
	UseDiscount       bool `json:"use_discount,omitempty"`
	DiscountPrice     int  `json:"discount_price,omitempty"`
	DiscountPriceType int  `json:"discount_price_type,omitempty"`
	UseDiscountPeriod bool `json:"use_discount_period,omitempty"`
	DiscountStartAt   string `json:"discount_start_at,omitempty"`
	DiscountEndAt     string `json:"discount_end_at,omitempty"`

	UseAccumulation       bool `json:"use_accumulation,omitempty"`
	AccumulationPoint     int  `json:"accumulation_point,omitempty"`
	AccumulationPointType int  `json:"accumulation_point_type,omitempty"`

	StatusDisplay      bool   `json:"status_display,omitempty"`
	UseDisplayPeriod   bool   `json:"use_display_period,omitempty"`
	DisplayStartAt     string `json:"display_start_at,omitempty"`
	DisplayEndAt       string `json:"display_end_at,omitempty"`
	StatusSale         bool   `json:"status_sale,omitempty"`
	UseSalePeriod      bool   `json:"use_sale_period,omitempty"`
	SaleStartAt        string `json:"sale_start_at,omitempty"`
	SaleEndAt          string `json:"sale_end_at,omitempty"`

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

	UseSetupFee    bool   `json:"use_setup_fee,omitempty"`
	SetupFeeValue  int    `json:"setup_fee_value,omitempty"`
	SetupFeeType   int    `json:"setup_fee_type,omitempty"`
	SetupFeeName   string `json:"setup_fee_name,omitempty"`
	SetupFeeText   string `json:"setup_fee_text,omitempty"`

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

// ProductListParams represents product list query parameters
type ProductListParams struct {
	ListParams
	Type         int    `json:"type,omitempty"`
	PeriodType   string `json:"period_type,omitempty"`
	SAt          string `json:"s_at,omitempty"`
	EAt          string `json:"e_at,omitempty"`
	CategoryCode string `json:"category_code,omitempty"`
}

// ProductStatusParams represents product status change parameters
type ProductStatusParams struct {
	ProductId     string `json:"product_id"`
	Status        int    `json:"status"`
	StatusDisplay bool   `json:"status_display,omitempty"`
	StatusSale    bool   `json:"status_sale,omitempty"`
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
	InvoiceId   string `json:"invoice_id,omitempty"`
	ProjectId   string `json:"project_id,omitempty"`
	SellerId    string `json:"seller_id,omitempty"`

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
type OrderListParams struct {
	ListParams
	UserId                  string   `json:"user_id,omitempty"`
	UserGroupId             string   `json:"user_group_id,omitempty"`
	Status                  []int    `json:"status,omitempty"`
	PaymentStatus           []int    `json:"payment_status,omitempty"`
	CsType                  string   `json:"cs_type,omitempty"`
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
	CancelProducts              []CancelProduct               `json:"cancel_products,omitempty"`
	CancelOrderSubscriptionBills []CancelOrderSubscriptionBill `json:"cancel_order_subscription_bills,omitempty"`
	CancelReason                string                        `json:"cancel_reason,omitempty"`
	CancelType                  int                           `json:"cancel_type,omitempty"`
	RefundPrice                 int                           `json:"refund_price,omitempty"`
}

// OrderCancelParams represents order cancel parameters
type OrderCancelParams struct {
	OrderNumber             string                  `json:"order_number,omitempty"`
	RequestCancelParameters *RequestCancelParameter `json:"request_cancel_parameters,omitempty"`
	IsSupervisor            bool                    `json:"is_supervisor,omitempty"`
}

// OrderCancelActionParams represents order cancel action parameters
type OrderCancelActionParams struct {
	OrderCancelRequestHistoryId string `json:"order_cancel_request_history_id"`
	CancelReason                string `json:"cancel_reason,omitempty"`
	RefundPrice                 int    `json:"refund_price,omitempty"`
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
type OrderSubscriptionListParams struct {
	ListParams
	SAt         string `json:"s_at,omitempty"`
	EAt         string `json:"e_at,omitempty"`
	RequestType string `json:"request_type,omitempty"`
	UserGroupId string `json:"user_group_id,omitempty"`
	UserId      string `json:"user_id,omitempty"`
}

// OrderSubscriptionUpdateParams represents order subscription update parameters
type OrderSubscriptionUpdateParams struct {
	OrderSubscriptionId string `json:"order_subscription_id"`
	NextBillingAt       string `json:"next_billing_at,omitempty"`
	BillingKey          string `json:"billing_key,omitempty"`
	Status              int    `json:"status,omitempty"`
	PaymentNextAt       string `json:"payment_next_at,omitempty"`
	ServiceEndAt        string `json:"service_end_at,omitempty"`
}

// OrderSubscriptionPauseParams represents subscription pause parameters
type OrderSubscriptionPauseParams struct {
	OrderSubscriptionId string `json:"order_subscription_id,omitempty"`
	OrderNumber         string `json:"order_number,omitempty"`
	Reason              string `json:"reason,omitempty"`
	PausedAt            string `json:"paused_at,omitempty"`
	ExpectedResumeAt    string `json:"expected_resume_at,omitempty"`
}

// OrderSubscriptionResumeParams represents subscription resume parameters
type OrderSubscriptionResumeParams struct {
	OrderSubscriptionId string `json:"order_subscription_id,omitempty"`
	OrderNumber         string `json:"order_number,omitempty"`
	ResumeAt            string `json:"resume_at,omitempty"`
}

// OrderSubscriptionTerminationParams represents subscription termination parameters
type OrderSubscriptionTerminationParams struct {
	OrderSubscriptionId  string `json:"order_subscription_id,omitempty"`
	OrderNumber          string `json:"order_number,omitempty"`
	TerminationFee       int    `json:"termination_fee,omitempty"`
	LastBillRefundPrice  int    `json:"last_bill_refund_price,omitempty"`
	FinalFee             int    `json:"final_fee,omitempty"`
	ServiceEndAt         string `json:"service_end_at,omitempty"`
	Reason               string `json:"reason,omitempty"`
}

// CalcTerminateFeeResponse represents terminate fee calculation response
type CalcTerminateFeeResponse struct {
	TerminationFee      int `json:"termination_fee,omitempty"`
	RefundAmount        int `json:"refund_amount,omitempty"`
	LastBillRefundPrice int `json:"last_bill_refund_price,omitempty"`
	FinalFee            int `json:"final_fee,omitempty"`
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

	ProductIds              []string `json:"product_ids,omitempty"`
	ProductOptionIds        []string `json:"product_option_ids,omitempty"`
	ProductSnapshotIds      []string `json:"product_snapshot_ids,omitempty"`
	ProductOptionSnapshotIds []string `json:"product_option_snapshot_ids,omitempty"`
	ProductType             int      `json:"product_type,omitempty"`
	Quantity                int      `json:"quantity,omitempty"`

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
type OrderSubscriptionBillListParams struct {
	ListParams
	OrderSubscriptionId string `json:"order_subscription_id,omitempty"`
	Status              []int  `json:"status,omitempty"`
}

// ============================================
// OrderSubscriptionAdjustment Types
// ============================================

// Constants for subscription adjustment type
const (
	SUBSCRIPTION_ADJUSTMENT_TYPE_PERIOD_DISCOUNT = 1
)

// CommerceOrderSubscriptionAdjustment represents an order subscription adjustment
type CommerceOrderSubscriptionAdjustment struct {
	OrderSubscriptionAdjustmentId string `json:"order_subscription_adjustment_id,omitempty"`
	Duration                      int    `json:"duration,omitempty"`
	Price                         int    `json:"price,omitempty"`
	TaxFreePrice                  int    `json:"tax_free_price,omitempty"`
	Name                          string `json:"name,omitempty"`
	Type                          int    `json:"type,omitempty"`
	CreatedAt                     string `json:"created_at,omitempty"`
}

// OrderSubscriptionAdjustmentUpdateParams represents subscription adjustment update parameters
type OrderSubscriptionAdjustmentUpdateParams struct {
	OrderSubscriptionId           string `json:"order_subscription_id"`
	OrderSubscriptionAdjustmentId string `json:"order_subscription_adjustment_id,omitempty"`
	Duration                      int    `json:"duration,omitempty"`
	Price                         int    `json:"price,omitempty"`
	TaxFreePrice                  int    `json:"tax_free_price,omitempty"`
	Name                          string `json:"name,omitempty"`
	Type                          int    `json:"type,omitempty"`
}
