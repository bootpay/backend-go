package bootpay

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
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

// Bootpay-Role scope 헤더를 구성한다 (Idempotency-Key 미지정시 자동 생성)
func commerceHeaders(role string, idempotencyKey string) map[string]string {
	if idempotencyKey == "" {
		idempotencyKey = randomIdempotencyKey()
	}
	return map[string]string{
		"Idempotency-Key": idempotencyKey,
		"Bootpay-Role":    role,
	}
}

// supervisor scope 전용 헤더를 구성한다 (Idempotency-Key 미지정시 자동 생성)
func supervisorHeaders(idempotencyKey string) map[string]string {
	return commerceHeaders("supervisor", idempotencyKey)
}

// user scope 전용 헤더를 구성한다 (Idempotency-Key 미지정시 자동 생성)
func userHeaders(idempotencyKey string) map[string]string {
	return commerceHeaders("user", idempotencyKey)
}

// manager scope 전용 헤더를 구성한다 (Idempotency-Key 미지정시 자동 생성)
func managerHeaders(idempotencyKey string) map[string]string {
	return commerceHeaders("manager", idempotencyKey)
}

// project_id가 있으면 supervisor(프로젝트 전체 검색), 없으면 user scope로 조회한다
func projectScopeRole(projectId interface{}) string {
	if projectId == nil {
		return "user"
	}
	if value, ok := projectId.(string); ok && value == "" {
		return "user"
	}
	return "supervisor"
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

// 호출자가 지정하지 않은(nil인) 항목만 기본값으로 채운다
func withDefaults(params APIResponse, defaults APIResponse) APIResponse {
	merged := APIResponse{}
	for name, value := range defaults {
		merged[name] = value
	}
	for name, value := range params {
		if value == nil {
			continue
		}
		merged[name] = value
	}
	return merged
}

// nil이거나 빈 문자열인 값은 제외하고 query string을 구성한다
// 배열은 Rails 규약에 맞춰 key[]=v1&key[]=v2 형식으로 전송한다
func buildQuery(params APIResponse) string {
	query := url.Values{}
	for name, value := range params {
		if value == nil {
			continue
		}
		switch typed := value.(type) {
		case string:
			if typed == "" {
				continue
			}
			query.Set(name, typed)
		default:
			reflected := reflect.ValueOf(value)
			if reflected.Kind() == reflect.Slice || reflected.Kind() == reflect.Array {
				for index := 0; index < reflected.Len(); index++ {
					query.Add(name+"[]", fmt.Sprintf("%v", reflected.Index(index).Interface()))
				}
				continue
			}
			query.Set(name, fmt.Sprintf("%v", value))
		}
	}
	return query.Encode()
}

// query string이 있는 경우에만 path에 덧붙인다
func commercePath(path string, params APIResponse) string {
	if encoded := buildQuery(params); encoded != "" {
		return path + "?" + encoded
	}
	return path
}

// scope 헤더가 필요한 쓰기 요청을 전송한다
func (api *Api) commerceWriteWithHeaders(method string, path string, payload interface{}, headers map[string]string) (APIResponse, error) {
	postBody, _ := json.Marshal(payload)
	return api.commerceRequest(method, path, bytes.NewBuffer(postBody), headers)
}

// multipart/form-data로 전송할 파일 (Field는 images[0] 처럼 서버가 읽는 form 필드명)
type multipartFile struct {
	Field string
	Path  string
}

// multipart form 값을 정규화한다 (배열·해시는 JSON, 나머지는 문자열)
func multipartValue(value interface{}) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case bool:
		return strconv.FormatBool(typed)
	}
	switch reflect.ValueOf(value).Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
		if encoded, err := json.Marshal(value); err == nil {
			return string(encoded)
		}
	}
	return fmt.Sprintf("%v", value)
}

// multipart/form-data 전송 (파일 업로드용)
// ⚠️ Content-Type은 boundary가 포함된 값으로 덮어써야 한다. JSON용 헤더를 그대로 두면 본문이 깨진다
func (api *Api) commerceMultipart(path string, fields APIResponse, files []multipartFile, headers map[string]string) (APIResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for name, value := range fields {
		if value == nil {
			continue
		}
		if err := writer.WriteField(name, multipartValue(value)); err != nil {
			return APIResponse{}, err
		}
	}
	for _, upload := range files {
		if err := writeMultipartFile(writer, upload); err != nil {
			return APIResponse{}, err
		}
	}
	if err := writer.Close(); err != nil {
		return APIResponse{}, err
	}

	merged := map[string]string{"Content-Type": writer.FormDataContentType()}
	for name, value := range headers {
		merged[name] = value
	}
	return api.commerceRequest(http.MethodPost, path, body, merged)
}

func writeMultipartFile(writer *multipart.Writer, upload multipartFile) error {
	file, err := os.Open(upload.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile(upload.Field, filepath.Base(upload.Path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	return err
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
// 상품 목록을 조회한다 (GET /products)
// ⚠️ keyword는 서버가 읽지 않는다. 컨트롤러는 page·limit·category_id·ex_uid·sort만 사용하므로
// keyword를 보내도 조용히 무시된다 (하위호환 때문에 인자는 유지한다)
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

// Invoice
// 청구서 목록을 조회한다 (GET /invoices), user scope
// params: page, limit, keyword, cs_type, user_id, product_type, css_at, cse_at
// 응답은 { list: [...], count: N } 구조다 ({ items, total }이 아니다)
// 서버 기본 limit은 24다
func (api *Api) InvoiceList(params APIResponse, idempotencyKey string) (APIResponse, error) {
	query := withDefaults(params, APIResponse{"page": 1, "limit": 24})
	return api.commerceRequest(http.MethodGet, commercePath("/invoices", query), nil, userHeaders(idempotencyKey))
}

// 청구서 상세를 조회한다 (GET /invoices/:id), user scope
func (api *Api) InvoiceDetail(invoiceId string, idempotencyKey string) (APIResponse, error) {
	return api.commerceRequest(http.MethodGet, "/invoices/"+invoiceId, nil, userHeaders(idempotencyKey))
}

// 청구서를 재안내한다 (POST /invoices/:id/notify), user scope
// sendTypes가 비어있으면 서버가 빈 배열로 처리한다
// ⚠️ 실제 고객에게 알림이 발송되므로 테스트 호출에 주의한다
func (api *Api) InvoiceNotify(invoiceId string, sendTypes []string, idempotencyKey string) (APIResponse, error) {
	payload := APIResponse{}
	if len(sendTypes) > 0 {
		payload["send_types"] = sendTypes
	}
	return api.commerceWriteWithHeaders(http.MethodPost, "/invoices/"+invoiceId+"/notify", payload, userHeaders(idempotencyKey))
}

// Order Cancel
// 주문 취소 요청 내역을 조회한다 (GET /order/cancel), user scope
// orderNumber 또는 orderId로 필터하며 둘 다 비어있으면 전체를 조회한다
// 승인·반려·철회에 넘길 order_cancellation_request_id를 여기서 얻는다
func (api *Api) OrderCancelList(orderNumber string, orderId string, idempotencyKey string) (APIResponse, error) {
	path := commercePath("/order/cancel", APIResponse{
		"order_number": orderNumber,
		"order_id":     orderId,
	})
	return api.commerceRequest(http.MethodGet, path, nil, userHeaders(idempotencyKey))
}

// (구매자) 주문 취소 요청을 철회한다 (PUT /order/cancel/:id/withdraw), user scope
// ⚠️ DELETE /order/cancel/:id (destroy)와는 다른 라우트다. 매뉴얼이 문서화한 withdraw를 사용한다
func (api *Api) OrderCancelWithdraw(orderCancellationRequestId string, idempotencyKey string) (APIResponse, error) {
	return api.commerceRequest(http.MethodPut, "/order/cancel/"+orderCancellationRequestId+"/withdraw", nil, userHeaders(idempotencyKey))
}

// Order Subscription
// 구독 계약 내용을 변경한다 (PUT /order_subscriptions/:id), supervisor scope
// 바뀐 값만 보내면 되며 나머지는 서버가 그대로 유지한다
// payload: product_id, product_option_id, order_name, total_subscription_duration, quantity,
// address_id, username, phone, email, use_free_trial, free_trial_day, service_start_at, service_end_at
func (api *Api) OrderSubscriptionUpdate(orderSubscriptionId string, payload APIResponse, idempotencyKey string) (APIResponse, error) {
	return api.commerceWriteWithHeaders(http.MethodPut, "/order_subscriptions/"+orderSubscriptionId, compactPayload(payload), supervisorHeaders(idempotencyKey))
}

// 가감산 조정항목을 추가한다 (POST /order_subscriptions/:id/adjustments), supervisor scope
// ⚠️ /adjustments 한 경로에 POST·PUT·DELETE 세 동사가 걸려있으므로 메서드를 혼동하지 않는다
// type 미전달시 서버가 price > 0 이면 SETUP_PRICE, 아니면 PERIOD_DISCOUNT로 자동 판정한다
func (api *Api) OrderSubscriptionAdjustmentCreate(orderSubscriptionId string, adjustment APIResponse, idempotencyKey string) (APIResponse, error) {
	payload := withDefaults(adjustment, APIResponse{
		"price":          0,
		"duration":       1,
		"tax_free_price": 0,
	})
	return api.commerceWriteWithHeaders(http.MethodPost, "/order_subscriptions/"+orderSubscriptionId+"/adjustments", payload, supervisorHeaders(idempotencyKey))
}

// 특정 회차의 조정항목을 통째로 교체한다 (PUT /order_subscriptions/:id/adjustments), supervisor scope
// 서버는 duration(회차) 단위로 갈아끼운다
func (api *Api) OrderSubscriptionAdjustmentUpdate(orderSubscriptionId string, duration int, adjustments []APIResponse, idempotencyKey string) (APIResponse, error) {
	if adjustments == nil {
		adjustments = []APIResponse{}
	}
	payload := APIResponse{
		"duration":    duration,
		"adjustments": adjustments,
	}
	return api.commerceWriteWithHeaders(http.MethodPut, "/order_subscriptions/"+orderSubscriptionId+"/adjustments", payload, supervisorHeaders(idempotencyKey))
}

// 조정항목을 삭제한다 (DELETE /order_subscriptions/:id/adjustments), supervisor scope
func (api *Api) OrderSubscriptionAdjustmentDelete(orderSubscriptionId string, orderSubscriptionAdjustmentId string, idempotencyKey string) (APIResponse, error) {
	payload := APIResponse{"order_subscription_adjustment_id": orderSubscriptionAdjustmentId}
	return api.commerceWriteWithHeaders(http.MethodDelete, "/order_subscriptions/"+orderSubscriptionId+"/adjustments", payload, supervisorHeaders(idempotencyKey))
}

// 구독 빌(회차) 목록을 조회한다 (GET /order_subscription_bills), user scope
// params: order_subscription_id, page, limit, status
// ⚠️ 경로가 order_subscription_bills로 언더스코어다 (하이픈이 아니다)
func (api *Api) OrderSubscriptionBillList(params APIResponse, idempotencyKey string) (APIResponse, error) {
	query := withDefaults(params, APIResponse{"page": 1, "limit": 20})
	return api.commerceRequest(http.MethodGet, commercePath("/order_subscription_bills", query), nil, userHeaders(idempotencyKey))
}

// 구독 진행중 요청 (requests/ing)
// ⚠️ 동사가 제각각이다. pause·purchase·termination·transfer는 POST, resume만 PUT,
// calculate_termination_fee는 GET이다. resume을 POST로 바꾸지 않는다
// ⚠️ SupervisorRequestOrderSubscription* (관리자가 즉시 실행)와는 다른 라우트로,
// 이쪽은 구매자가 요청을 올리는 면이다

// order_subscription_id를 payload에 병합한다 (requests/ing 계열 공통)
func orderSubscriptionRequestPayload(orderSubscriptionId string, payload APIResponse) APIResponse {
	merged := compactPayload(payload)
	merged["order_subscription_id"] = orderSubscriptionId
	return merged
}

// 구독 일시중지 요청 (POST /order_subscriptions/requests/ing/pause), user scope
// payload: reason, paused_at, expected_resume_at
func (api *Api) OrderSubscriptionRequestsIngPause(orderSubscriptionId string, payload APIResponse, idempotencyKey string) (APIResponse, error) {
	return api.commerceWriteWithHeaders(http.MethodPost, "/order_subscriptions/requests/ing/pause",
		orderSubscriptionRequestPayload(orderSubscriptionId, payload), userHeaders(idempotencyKey))
}

// 구독 재개 요청 (PUT /order_subscriptions/requests/ing/resume), user scope
// ⚠️ requests/ing 계열 중 유일하게 PUT이다
// payload: reason
func (api *Api) OrderSubscriptionRequestsIngResume(orderSubscriptionId string, payload APIResponse, idempotencyKey string) (APIResponse, error) {
	return api.commerceWriteWithHeaders(http.MethodPut, "/order_subscriptions/requests/ing/resume",
		orderSubscriptionRequestPayload(orderSubscriptionId, payload), userHeaders(idempotencyKey))
}

// 중도인수 요청 (POST /order_subscriptions/requests/ing/purchase), user scope
// payload: price, tax_free_price, reason
func (api *Api) OrderSubscriptionRequestsIngPurchase(orderSubscriptionId string, payload APIResponse, idempotencyKey string) (APIResponse, error) {
	return api.commerceWriteWithHeaders(http.MethodPost, "/order_subscriptions/requests/ing/purchase",
		orderSubscriptionRequestPayload(orderSubscriptionId, payload), userHeaders(idempotencyKey))
}

// 중도해지 요청 (POST /order_subscriptions/requests/ing/termination), user scope
// payload: order_number, reason, termination_fee, last_bill_refund_price, final_fee, service_end_at
func (api *Api) OrderSubscriptionRequestsIngTermination(orderSubscriptionId string, payload APIResponse, idempotencyKey string) (APIResponse, error) {
	return api.commerceWriteWithHeaders(http.MethodPost, "/order_subscriptions/requests/ing/termination",
		orderSubscriptionRequestPayload(orderSubscriptionId, payload), userHeaders(idempotencyKey))
}

// 구독 이전·승계 요청 (POST /order_subscriptions/requests/ing/transfer), user scope
// payload: new_user_id, new_username, new_user_email, new_user_phone, new_user_address, wallet_id, reason
func (api *Api) OrderSubscriptionRequestsIngTransfer(orderSubscriptionId string, payload APIResponse, idempotencyKey string) (APIResponse, error) {
	return api.commerceWriteWithHeaders(http.MethodPost, "/order_subscriptions/requests/ing/transfer",
		orderSubscriptionRequestPayload(orderSubscriptionId, payload), userHeaders(idempotencyKey))
}

// 중도해지 수수료 사전계산 (GET /order_subscriptions/requests/ing/calculate_termination_fee), user scope
// 해지 요청 전에 얼마가 나오는지 미리 보여줄 때 사용한다
func (api *Api) OrderSubscriptionCalculateTerminationFee(orderSubscriptionId string, orderNumber string, idempotencyKey string) (APIResponse, error) {
	path := commercePath("/order_subscriptions/requests/ing/calculate_termination_fee", APIResponse{
		"order_subscription_id": orderSubscriptionId,
		"order_number":          orderNumber,
	})
	return api.commerceRequest(http.MethodGet, path, nil, userHeaders(idempotencyKey))
}

// 구독 요청 리소스 (order-subscription-requests)
// ⚠️ 하이픈 경로는 order-subscription-requests와 user-groups 둘 뿐이다.
// order_subscriptions · order_subscription_bills는 언더스코어다

// 구독 변경요청 목록 (GET /order-subscription-requests)
// params: project_id, order_subscription_id, page, limit, keyword, s_at, e_at, status, request_type, user_id, user_group_id
// project_id를 주면 supervisor 모드(프로젝트 전체 검색), 없으면 본인 요청만 조회한다
func (api *Api) OrderSubscriptionRequestList(params APIResponse, idempotencyKey string) (APIResponse, error) {
	query := withDefaults(params, APIResponse{"page": 1, "limit": 20})
	headers := commerceHeaders(projectScopeRole(query["project_id"]), idempotencyKey)
	return api.commerceRequest(http.MethodGet, commercePath("/order-subscription-requests", query), nil, headers)
}

// 구독 변경요청 상세 (GET /order-subscription-requests/:id)
func (api *Api) OrderSubscriptionRequestDetail(requestHistoryId string, projectId string, idempotencyKey string) (APIResponse, error) {
	path := commercePath("/order-subscription-requests/"+requestHistoryId, APIResponse{"project_id": projectId})
	return api.commerceRequest(http.MethodGet, path, nil, commerceHeaders(projectScopeRole(projectId), idempotencyKey))
}

// 구독 변경요청 승인·반려 (PUT /order-subscription-requests/:id), supervisor scope
// ⚠️ 승인과 반려는 별도 액션이 아니라 approval("approve" 또는 "reject") 값으로 갈린다
// (서버가 params[:action]을 Rails 예약어로 쓰기 때문에 키 이름이 approval이다)
// payload: reason, price, tax_free_price, termination_fee, last_bill_refund_price, final_fee, service_end_at
func (api *Api) OrderSubscriptionRequestUpdate(requestHistoryId string, approval string, payload APIResponse, idempotencyKey string) (APIResponse, error) {
	body := compactPayload(payload)
	body["approval"] = approval
	return api.commerceWriteWithHeaders(http.MethodPut, "/order-subscription-requests/"+requestHistoryId, body, supervisorHeaders(idempotencyKey))
}

// Product 쓰기
// 상품을 등록한다 (POST /products), manager scope
// product: name, display_price, desc, content, category_id, type, stock, status_sale, status_display,
// use_subscription, subscription_setting_id, save_by 등 (여기 없는 값도 그대로 전달된다)
// images에 업로드할 파일 경로를 넘기면 multipart/form-data로, 없으면 JSON으로 전송한다
func (api *Api) ProductCreate(product APIResponse, images []string, idempotencyKey string) (APIResponse, error) {
	payload := compactPayload(product)
	headers := managerHeaders(idempotencyKey)
	if len(images) == 0 {
		return api.commerceWriteWithHeaders(http.MethodPost, "/products", payload, headers)
	}

	files := make([]multipartFile, 0, len(images))
	for index, image := range images {
		files = append(files, multipartFile{Field: fmt.Sprintf("images[%d]", index), Path: image})
	}
	return api.commerceMultipart("/products", payload, files, headers)
}

// 상품을 수정한다 (PUT /products/:id), manager scope
// 바뀐 값만 보내면 된다. category_id는 키 존재 여부로 해제 의사를 판별하므로 주의한다
func (api *Api) ProductUpdate(productId string, product APIResponse, idempotencyKey string) (APIResponse, error) {
	return api.commerceWriteWithHeaders(http.MethodPut, "/products/"+productId, compactPayload(product), managerHeaders(idempotencyKey))
}

// 상품을 삭제한다 (DELETE /products/:id), manager scope
func (api *Api) ProductDelete(productId string, idempotencyKey string) (APIResponse, error) {
	return api.commerceRequest(http.MethodDelete, "/products/"+productId, nil, managerHeaders(idempotencyKey))
}

// 상품 판매·노출 상태를 변경한다 (PUT /products/:id/status), manager scope
// status: status_sale, status_display, status_frozen, status_review, use_display_period,
// display_start_at, display_end_at, use_sale_period, sale_start_at, sale_end_at
// 재고(stock)는 여기가 아니라 ProductUpdate로 변경한다
func (api *Api) ProductStatus(productId string, status APIResponse, idempotencyKey string) (APIResponse, error) {
	return api.commerceWriteWithHeaders(http.MethodPut, "/products/"+productId+"/status", compactPayload(status), managerHeaders(idempotencyKey))
}

// User 쓰기
// 외부 uid(ex_uid) 중복검사 (GET /users/join/uid-exist), user scope
func (api *Api) UidExist(uid string, idempotencyKey string) (APIResponse, error) {
	path := commercePath("/users/join/uid-exist", APIResponse{"pk": uid})
	return api.commerceRequest(http.MethodGet, path, nil, userHeaders(idempotencyKey))
}

// 회원 정보를 수정한다 (PUT /users/:id), user scope
// user: login_id, login_pw, name, phone, email, tel, nickname, bank_username, bank_account,
// bank_code, comment, gender, birth, group 등 (여기 없는 값도 그대로 전달된다)
// 사업자 정보는 group: { company_name, business_number, registration_number } 로 중첩 전달한다
func (api *Api) UserUpdate(userId string, user APIResponse, idempotencyKey string) (APIResponse, error) {
	return api.commerceWriteWithHeaders(http.MethodPut, "/users/"+userId, compactPayload(user), userHeaders(idempotencyKey))
}

// User Group
// 그룹 구매한도를 설정한다 (PUT /user-groups/:id/limit), manager scope
// limit: use_limit, limit_month_purchase, limit_week_purchase, limit_message
// ⚠️ 그룹 수정 API로는 한도가 반영되지 않는다. 서버가 use_limit·limit_* 를 제거하므로 이 전용 라우트로만 바뀐다
func (api *Api) UserGroupLimit(userGroupId string, limit APIResponse, idempotencyKey string) (APIResponse, error) {
	return api.commerceWriteWithHeaders(http.MethodPut, "/user-groups/"+userGroupId+"/limit", compactPayload(limit), managerHeaders(idempotencyKey))
}

// 그룹 구독 합산청구(정산주기) 설정을 변경한다 (PUT /user-groups/:id/aggregate-transaction), manager scope
// setting: use_subscription_aggregate_transaction, subscription_month_day, subscription_week_day
func (api *Api) UserGroupAggregateTransaction(userGroupId string, setting APIResponse, idempotencyKey string) (APIResponse, error) {
	return api.commerceWriteWithHeaders(http.MethodPut, "/user-groups/"+userGroupId+"/aggregate-transaction", compactPayload(setting), managerHeaders(idempotencyKey))
}

// Webhook
// 테스트 웹훅을 발송한다 (POST /webhook/test)
// headerContentType이 비어있으면 서버 기본값으로 발송된다
func (api *Api) SendTestWebhook(headerContentType string) (APIResponse, error) {
	payload := APIResponse{}
	if headerContentType != "" {
		payload["header_content_type"] = headerContentType
	}
	return api.commerceWrite(http.MethodPost, "/webhook/test", payload)
}
