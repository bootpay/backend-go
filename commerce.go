package bootpay

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (api *Api) commerceGet(path string) (APIResponse, error) {
	return api.commerceRequest(http.MethodGet, path, nil, nil)
}

func (api *Api) commerceWrite(method string, path string, payload interface{}) (APIResponse, error) {
	postBody, _ := json.Marshal(payload)
	return api.commerceRequest(method, path, bytes.NewBuffer(postBody), nil)
}

func (api *Api) commerceRequest(method string, path string, body io.Reader, headers map[string]string) (APIResponse, error) {
	req, err := api.NewRequest(method, path, body)
	if err != nil {
		return APIResponse{}, err
	}
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	res, err := api.client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	result["http_status_code"] = res.StatusCode
	return result, nil
}

// supervisor scope 전용 헤더를 구성한다 (Idempotency-Key 미지정시 자동 생성)
func supervisorHeaders(idempotencyKey string) map[string]string {
	if idempotencyKey == "" {
		idempotencyKey = randomIdempotencyKey()
	}
	return map[string]string{
		"Idempotency-Key": idempotencyKey,
		"Bootpay-Role":    "supervisor",
	}
}

func randomIdempotencyKey() string {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 16)
	}
	buffer[6] = (buffer[6] & 0x0f) | 0x40
	buffer[8] = (buffer[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", buffer[0:4], buffer[4:6], buffer[6:8], buffer[8:10], buffer[10:16])
}

// 값이 nil인 항목은 서버로 전송하지 않는다
func compactPayload(payload APIResponse) APIResponse {
	compacted := APIResponse{}
	for name, value := range payload {
		if value == nil {
			continue
		}
		compacted[name] = value
	}
	return compacted
}

// Store
func (api *Api) GetStore() (APIResponse, error) {
	return api.commerceGet("/store")
}

func (api *Api) StoreInfo() (APIResponse, error) {
	return api.GetStore()
}

func (api *Api) GetStoreDetail() (APIResponse, error) {
	return api.commerceGet("/store/detail")
}

func (api *Api) StoreDetail() (APIResponse, error) {
	return api.GetStoreDetail()
}

// User
func (api *Api) UserLogin(loginId string, loginPw string) (APIResponse, error) {
	return api.commerceWrite(http.MethodPost, "/users/login", APIResponse{
		"login_id": loginId,
		"login_pw": loginPw,
	})
}

func (api *Api) UserJoin(user APIResponse) (APIResponse, error) {
	return api.commerceWrite(http.MethodPost, "/users/join", user)
}

func (api *Api) UserJoinCheck(checkType string, pk string) (APIResponse, error) {
	encoded := url.QueryEscape(pk)
	return api.commerceGet(fmt.Sprintf("/users/join/%s?pk=%s", checkType, encoded))
}

// Mall aliases
func (api *Api) MallUserLogin(loginId string, loginPw string) (APIResponse, error) {
	return api.UserLogin(loginId, loginPw)
}

func (api *Api) MallUserJoin(user APIResponse) (APIResponse, error) {
	return api.UserJoin(user)
}

func (api *Api) MallUserJoinCheck(checkType string, pk string) (APIResponse, error) {
	return api.UserJoinCheck(checkType, pk)
}

// Product
func (api *Api) ProductList(params APIResponse) (APIResponse, error) {
	q := url.Values{}
	if params != nil {
		if v, ok := params["page"]; ok {
			q.Set("page", fmt.Sprintf("%v", v))
		}
		if v, ok := params["limit"]; ok {
			q.Set("limit", fmt.Sprintf("%v", v))
		}
		if v, ok := params["keyword"]; ok {
			q.Set("keyword", fmt.Sprintf("%v", v))
		}
		if v, ok := params["type"]; ok {
			switch t := v.(type) {
			case int:
				q.Set("type", strconv.Itoa(t))
			default:
				q.Set("type", fmt.Sprintf("%v", v))
			}
		}
	}
	path := "/products"
	if encoded := q.Encode(); encoded != "" {
		path = path + "?" + encoded
	}
	return api.commerceGet(path)
}

func (api *Api) ProductDetail(productId string) (APIResponse, error) {
	return api.commerceGet("/products/" + productId)
}

// Mall aliases
func (api *Api) Products(params APIResponse) (APIResponse, error) {
	return api.ProductList(params)
}

func (api *Api) MallProductDetail(productId string) (APIResponse, error) {
	return api.ProductDetail(productId)
}

// Mall Setting
// 몰 설정 조회 (GET /mall-setting), supervisor scope 전용
func (api *Api) GetMallSetting(idempotencyKey string) (APIResponse, error) {
	return api.commerceRequest(http.MethodGet, "/mall-setting", nil, supervisorHeaders(idempotencyKey))
}

// 몰 설정 수정 (PUT /mall-setting), supervisor scope 전용
// 요청 바디는 flatten 형식이며 전달된 값(non-nil)만 서버로 전송된다.
func (api *Api) UpdateMallSetting(setting APIResponse, idempotencyKey string) (APIResponse, error) {
	postBody, _ := json.Marshal(compactPayload(setting))
	return api.commerceRequest(http.MethodPut, "/mall-setting", bytes.NewBuffer(postBody), supervisorHeaders(idempotencyKey))
}

// Supervisor
type OrderSubscriptionChargePayload struct {
	IdempotencyKey string      `json:"-"`
	ChargeKey      string      `json:"charge_key"`
	Price          float64     `json:"price"`
	TaxFreePrice   float64     `json:"tax_free_price,omitempty"`
	User           APIResponse `json:"user,omitempty"`
	Metadata       APIResponse `json:"metadata,omitempty"`
}

type OrderSubscriptionChargeRevokePayload struct {
	IdempotencyKey string      `json:"-"`
	ChargeKey      string      `json:"charge_key"`
	User           APIResponse `json:"user,omitempty"`
}

// 수시결제(온디맨드) charge_key 즉시 결제
// charge_key는 body로만 전송한다 (URL/query 금지 - 액세스 로그 노출 방지)
func (api *Api) SupervisorRequestOrderSubscriptionCharge(payload OrderSubscriptionChargePayload) (APIResponse, error) {
	postBody, _ := json.Marshal(payload)
	return api.commerceRequest(http.MethodPost, "/order_subscriptions/charge", bytes.NewBuffer(postBody), supervisorHeaders(payload.IdempotencyKey))
}

// 수시결제(온디맨드) charge_key 해지
// 해지 이후 해당 키로의 재결제는 불가능하다
func (api *Api) SupervisorRequestOrderSubscriptionChargeRevoke(payload OrderSubscriptionChargeRevokePayload) (APIResponse, error) {
	postBody, _ := json.Marshal(payload)
	return api.commerceRequest(http.MethodDelete, "/order_subscriptions/charge", bytes.NewBuffer(postBody), supervisorHeaders(payload.IdempotencyKey))
}
