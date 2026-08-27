package bootpay

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	COMMERCE_DEVELOPMENT string = "https://dev-api.bootapi.com/v1"
	COMMERCE_STAGE       string = "https://stage-api.bootapi.com/v1"
	COMMERCE_PRODUCTION  string = "https://api.bootapi.com/v1"

	COMMERCE_API_VERSION string = "1.0.0"
	COMMERCE_SDK_VERSION string = "1.0.0"
)

// CommerceApi is the main struct for Commerce API
type CommerceApi struct {
	token     string
	clientKey string
	secretKey string
	baseUrl   string
	role      string
	client    *http.Client

	// Modules
	User                        *UserModule
	UserGroup                   *UserGroupModule
	Product                     *ProductModule
	Invoice                     *InvoiceModule
	Order                       *OrderModule
	OrderCancel                 *OrderCancelModule
	OrderSubscription           *OrderSubscriptionModule
	OrderSubscriptionBill       *OrderSubscriptionBillModule
	OrderSubscriptionAdjustment *OrderSubscriptionAdjustmentModule
	OrderSubscriptionRequest    *OrderSubscriptionRequestModule
	Store                       *StoreModule
	Category                    *CategoryModule
	Coupon                      *CouponModule
	Point                       *PointModule
	Cart                        *CartModule
	MallSetting                 *MallSettingModule
	Webhook                     *WebhookModule

	// 알림톡 v1 API
	AlimtalkMessage  *AlimtalkMessageModule
	AlimtalkOfficial *AlimtalkOfficialModule
	AlimtalkOptout   *AlimtalkOptoutModule
	AlimtalkSend     *AlimtalkSendModule
	AlimtalkSender   *AlimtalkSenderModule
	AlimtalkTemplate *AlimtalkTemplateModule
	AlimtalkWebhook  *AlimtalkWebhookModule
}

// CommerceResponse is the common response structure for Commerce API
type CommerceResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	ErrorCode int         `json:"error_code,omitempty"`
	Message   string      `json:"message,omitempty"`
}

// CommerceListResponse is the common response structure for list APIs
type CommerceListResponse[T any] struct {
	Success   bool   `json:"success"`
	Data      T      `json:"data,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
	Message   string `json:"message,omitempty"`
}

// CommerceTokenResponse represents the token response
type CommerceTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiredAt   string `json:"expired_at,omitempty"`
}

// NewCommerceAPI creates a new Commerce API instance (recommended)
func NewCommerceAPI(clientKey string, secretKey string, client *http.Client, mode string) *CommerceApi {
	if client == nil {
		client = &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
			},
		}
	}

	baseUrl := COMMERCE_PRODUCTION
	if mode == "development" {
		baseUrl = COMMERCE_DEVELOPMENT
	} else if mode == "stage" {
		baseUrl = COMMERCE_STAGE
	}

	api := &CommerceApi{
		clientKey: clientKey,
		secretKey: secretKey,
		baseUrl:   baseUrl,
		role:      "user",
		client:    client,
	}

	// Initialize modules
	api.User = &UserModule{api: api}
	api.UserGroup = &UserGroupModule{api: api}
	api.Product = &ProductModule{api: api}
	api.Invoice = &InvoiceModule{api: api}
	api.Order = &OrderModule{api: api}
	api.OrderCancel = &OrderCancelModule{api: api}
	api.OrderSubscription = &OrderSubscriptionModule{
		api:        api,
		RequestIng: &OrderSubscriptionRequestIngModule{api: api},
	}
	api.OrderSubscriptionBill = &OrderSubscriptionBillModule{api: api}
	api.OrderSubscriptionAdjustment = &OrderSubscriptionAdjustmentModule{api: api}
	api.OrderSubscriptionRequest = &OrderSubscriptionRequestModule{api: api}
	api.Store = &StoreModule{api: api}
	api.Category = &CategoryModule{api: api}
	api.Coupon = &CouponModule{api: api}
	api.Point = &PointModule{api: api}
	api.Cart = &CartModule{api: api}
	api.MallSetting = &MallSettingModule{api: api}
	api.Webhook = &WebhookModule{api: api}
	api.AlimtalkMessage = &AlimtalkMessageModule{api: api}
	api.AlimtalkOfficial = &AlimtalkOfficialModule{api: api}
	api.AlimtalkOptout = &AlimtalkOptoutModule{api: api}
	api.AlimtalkSend = &AlimtalkSendModule{api: api}
	api.AlimtalkSender = &AlimtalkSenderModule{api: api}
	api.AlimtalkTemplate = &AlimtalkTemplateModule{api: api}
	api.AlimtalkWebhook = &AlimtalkWebhookModule{api: api}

	return api
}

// NewCommerceApi creates a new Commerce API instance (deprecated: use NewCommerceAPI instead)
func NewCommerceApi(clientKey string, secretKey string, client *http.Client, mode string) *CommerceApi {
	return NewCommerceAPI(clientKey, secretKey, client, mode)
}

// SetRole sets the role for API requests
func (api *CommerceApi) SetRole(role string) *CommerceApi {
	api.role = role
	return api
}

// AsUser sets role to "user"
func (api *CommerceApi) AsUser() *CommerceApi {
	return api.SetRole("user")
}

// AsManager sets role to "manager"
func (api *CommerceApi) AsManager() *CommerceApi {
	return api.SetRole("manager")
}

// AsPartner sets role to "partner"
func (api *CommerceApi) AsPartner() *CommerceApi {
	return api.SetRole("partner")
}

// AsVendor sets role to "vendor"
func (api *CommerceApi) AsVendor() *CommerceApi {
	return api.SetRole("vendor")
}

// AsSupervisor sets role to "supervisor"
func (api *CommerceApi) AsSupervisor() *CommerceApi {
	return api.SetRole("supervisor")
}

// GetRole returns the current role
func (api *CommerceApi) GetRole() string {
	return api.role
}

// GetToken returns the current token
func (api *CommerceApi) GetToken() string {
	return api.token
}

// SetToken sets the access token
func (api *CommerceApi) SetToken(token string) {
	api.token = token
}

func (api *CommerceApi) validateCredentials() error {
	if api.clientKey == "" || api.secretKey == "" {
		return errors.New("commerce: client_key and secret_key must be provided together")
	}
	return nil
}

// getBasicAuthHeader returns Basic Auth header value
func (api *CommerceApi) getBasicAuthHeader() string {
	if api.clientKey == "" || api.secretKey == "" {
		return ""
	}
	credentials := fmt.Sprintf("%s:%s", api.clientKey, api.secretKey)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return fmt.Sprintf("Basic %s", encoded)
}

// newRequest creates a new HTTP request with common headers
func (api *CommerceApi) newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	if err := api.validateCredentials(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, api.baseUrl+"/"+url, body)
	if err != nil {
		return nil, errors.New("cannot create Commerce API request: " + err.Error())
	}

	if basic := api.getBasicAuthHeader(); basic != "" {
		req.Header.Set("Authorization", basic)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("BOOTPAY-SDK-VERSION", COMMERCE_SDK_VERSION)
	req.Header.Set("BOOTPAY-API-VERSION", COMMERCE_API_VERSION)
	req.Header.Set("BOOTPAY-SDK-TYPE", "305")
	req.Header.Set("BOOTPAY-ROLE", api.role)

	return req, nil
}

// decodeCommerceBody parses a Commerce API response body into map[string]interface{}.
//
// 26-08-24: some endpoints (GET /v1/categories, /v1/coupon, /v1/coupon/available)
// answer with a top-level JSON array. Decoding those straight into a map fails, and
// the previous code discarded that error — callers silently received an empty map.
// A top-level array is now wrapped as {"data": [...]}, mirroring the Java SDK's
// BootpayStoreObject fallback. Object responses keep their exact previous shape.
// An empty body stays an empty map; anything else that is not JSON returns an error
// instead of a silent empty result (an HTML 5xx page used to look like success).
func decodeCommerceBody(body io.Reader) (map[string]interface{}, error) {
	raw, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(raw)) == 0 {
		return map[string]interface{}{}, nil
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(raw, &result); err == nil {
		return result, nil
	}

	var list []interface{}
	if listErr := json.Unmarshal(raw, &list); listErr == nil {
		if list == nil {
			list = []interface{}{}
		}
		return map[string]interface{}{"data": list}, nil
	}

	return nil, fmt.Errorf("commerce: cannot parse response body as JSON: %s", strings.TrimSpace(string(raw)))
}

// GetAccessToken obtains an access token using client_key and secret_key
func (api *CommerceApi) GetAccessToken() (map[string]interface{}, error) {
	if err := api.validateCredentials(); err != nil {
		return nil, err
	}
	data := map[string]string{
		"client_key": api.clientKey,
		"secret_key": api.secretKey,
	}

	postBody, _ := json.Marshal(data)
	body := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodPost, api.baseUrl+"/request/token", body)
	if err != nil {
		return nil, errors.New("commerce: getAccessToken error: " + err.Error())
	}

	req.Header.Set("Authorization", api.getBasicAuthHeader())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("BOOTPAY-SDK-VERSION", COMMERCE_SDK_VERSION)
	req.Header.Set("BOOTPAY-API-VERSION", COMMERCE_API_VERSION)
	req.Header.Set("BOOTPAY-SDK-TYPE", "305")
	req.Header.Set("BOOTPAY-ROLE", api.role)

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result, err := decodeCommerceBody(res.Body)
	if err != nil {
		return nil, err
	}

	if accessToken, ok := result["access_token"].(string); ok {
		api.token = accessToken
	}

	return result, nil
}

// doRequest performs an HTTP request and returns the response
func (api *CommerceApi) doRequest(method string, url string, data interface{}) (map[string]interface{}, error) {
	return api.doRequestWithHeaders(method, url, data, nil)
}

// doRequestWithHeaders performs an HTTP request with per-request headers.
// Per-request headers (BOOTPAY-ROLE, Idempotency-Key, Bootpay-User-JWT, ...) are applied
// after the common headers, so a request-specific BOOTPAY-ROLE is never overwritten by
// the instance default (supervisor-only endpoint support).
func (api *CommerceApi) doRequestWithHeaders(method string, url string, data interface{}, headers map[string]string) (map[string]interface{}, error) {
	var body io.Reader
	if data != nil {
		postBody, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(postBody)
	}

	req, err := api.newRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return decodeCommerceBody(res.Body)
}

// getWithHeaders performs a GET request with per-request headers
func (api *CommerceApi) getWithHeaders(url string, headers map[string]string) (map[string]interface{}, error) {
	return api.doRequestWithHeaders(http.MethodGet, url, nil, headers)
}

// postWithHeaders performs a POST request with per-request headers
func (api *CommerceApi) postWithHeaders(url string, data interface{}, headers map[string]string) (map[string]interface{}, error) {
	return api.doRequestWithHeaders(http.MethodPost, url, data, headers)
}

// putWithHeaders performs a PUT request with per-request headers
func (api *CommerceApi) putWithHeaders(url string, data interface{}, headers map[string]string) (map[string]interface{}, error) {
	return api.doRequestWithHeaders(http.MethodPut, url, data, headers)
}

// deleteWithHeaders performs a DELETE request with per-request headers.
// data (when non-nil) is sent as the request body — some endpoints
// (order_subscriptions/charge, adjustments) require the target id in the body, not the query.
func (api *CommerceApi) deleteWithHeaders(url string, data interface{}, headers map[string]string) (map[string]interface{}, error) {
	return api.doRequestWithHeaders(http.MethodDelete, url, data, headers)
}

// postMultipart performs a multipart/form-data POST.
// Content-Type must keep the boundary generated by the multipart writer —
// overwriting it makes the server parse the body as null.
func (api *CommerceApi) postMultipart(url string, body io.Reader, contentType string, headers map[string]string) (map[string]interface{}, error) {
	req, err := api.newRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", contentType)

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return decodeCommerceBody(res.Body)
}

// getRaw performs a GET request and returns the body without parsing it.
//
// ⚠️ decodeCommerceBody always parses the body as JSON, so an endpoint answering with a
// non-JSON body (alimtalk template export with format=csv) fails there and the caller sees
// a parse error for a request that actually succeeded. This path keeps the original text.
// The result is {"body": "<raw text>", "content_type": "..."} — mirroring the Ruby SDK.
func (api *CommerceApi) getRaw(url string, headers map[string]string) (map[string]interface{}, error) {
	req, err := api.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Accept", "*/*")

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"body":         string(raw),
		"content_type": res.Header.Get("Content-Type"),
	}, nil
}

// newIdempotencyKey generates a UUID v4 string for the Idempotency-Key header
func newIdempotencyKey() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("idem-%d", time.Now().UnixNano())
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// commerceIdempotencyHeaders returns headers with an Idempotency-Key (auto-generated when empty)
func commerceIdempotencyHeaders(idempotencyKey string) map[string]string {
	if idempotencyKey == "" {
		idempotencyKey = newIdempotencyKey()
	}
	return map[string]string{"Idempotency-Key": idempotencyKey}
}

// commerceRoleHeaders returns headers with a per-request BOOTPAY-ROLE and Idempotency-Key
func commerceRoleHeaders(role string, idempotencyKey string) map[string]string {
	headers := commerceIdempotencyHeaders(idempotencyKey)
	headers["BOOTPAY-ROLE"] = role
	return headers
}

// alimtalkHeaders returns the header set for the Alimtalk v1 API.
//
// ⚠️ No Idempotency-Key here, unlike commerceRoleHeaders. The alimtalk API does not read
// that header — idempotency is established only by ref_id on send. Attaching it anyway
// would advertise a guarantee the server never gives.
// BOOTPAY-ROLE is always user: every alimtalk scope key is user:alimtalk_*.
func alimtalkHeaders() map[string]string {
	return map[string]string{"BOOTPAY-ROLE": "user"}
}

// commerceMallHeaders returns V1 Mall API headers.
// Bootpay-User-JWT is attached only when a member JWT is present.
func commerceMallHeaders(userJwt string, idempotencyKey string) map[string]string {
	headers := commerceIdempotencyHeaders(idempotencyKey)
	if userJwt != "" {
		headers["Bootpay-User-JWT"] = userJwt
	}
	return headers
}

// Get performs a GET request
func (api *CommerceApi) Get(url string) (map[string]interface{}, error) {
	return api.doRequest(http.MethodGet, url, nil)
}

// Post performs a POST request
func (api *CommerceApi) Post(url string, data interface{}) (map[string]interface{}, error) {
	return api.doRequest(http.MethodPost, url, data)
}

// Put performs a PUT request
func (api *CommerceApi) Put(url string, data interface{}) (map[string]interface{}, error) {
	return api.doRequest(http.MethodPut, url, data)
}

// Delete performs a DELETE request
func (api *CommerceApi) Delete(url string) (map[string]interface{}, error) {
	return api.doRequest(http.MethodDelete, url, nil)
}

// BoolPtr returns a pointer to v.
//
// Several Commerce/Alimtalk fields are tri-state: nil follows the server/project default,
// while an explicit false means "turn this off" (AlimtalkSendParams.Fallback,
// AlimtalkTemplateCreateParams.Register, ...). A plain bool would lose that distinction
// because omitempty drops false — hence the pointer, and this helper to build one inline.
func BoolPtr(v bool) *bool {
	return &v
}

// firstOrEmpty returns the first element of an optional variadic string argument.
// Used by Commerce methods that accept an optional Idempotency-Key without breaking
// the existing call signature.
func firstOrEmpty(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

// withQuery appends an encoded query string to a uri, omitting the "?" when there is
// nothing to send (Ruby's params `.compact` leaves no trailing "?" either).
func withQuery(uri string, query url.Values) string {
	if len(query) == 0 {
		return uri
	}
	return uri + "?" + query.Encode()
}
