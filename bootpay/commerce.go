package bootpay

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

// getBasicAuthHeader returns Basic Auth header value
func (api *CommerceApi) getBasicAuthHeader() string {
	credentials := fmt.Sprintf("%s:%s", api.clientKey, api.secretKey)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return fmt.Sprintf("Basic %s", encoded)
}

// newRequest creates a new HTTP request with common headers
func (api *CommerceApi) newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, api.baseUrl+"/"+url, body)
	if err != nil {
		return nil, errors.New("cannot create Commerce API request: " + err.Error())
	}

	if api.token != "" {
		req.Header.Set("Authorization", "Bearer "+api.token)
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

// GetAccessToken obtains an access token using client_key and secret_key
func (api *CommerceApi) GetAccessToken() (map[string]interface{}, error) {
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

	result := make(map[string]interface{})
	json.NewDecoder(res.Body).Decode(&result)

	if accessToken, ok := result["access_token"].(string); ok {
		api.token = accessToken
	}

	return result, nil
}

// doRequest performs an HTTP request and returns the response
func (api *CommerceApi) doRequest(method string, url string, data interface{}) (map[string]interface{}, error) {
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

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result := make(map[string]interface{})
	json.NewDecoder(res.Body).Decode(&result)

	return result, nil
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
