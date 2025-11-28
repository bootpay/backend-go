package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

//type BillingKeyCardData struct {
//	CardCode string `json:"card_code"`
//	CardName string `json:"card_name"`
//	CardNo   string `json:"card_no"`
//	CardCl   string `json:"card_cl"`
//}
//type BillingKeyData struct {
//	BillingKey string             `json:"billing_key"`
//	PgName     string             `json:"pg_name"`
//	MethodName string             `json:"method_name"`
//	method     string             `json:"method"`
//	Data       BillingKeyCardData `json:"data"`
//	EndAt      string             `json:"e_at"`
//	CreateAt   string             `json:"c_at"`
//}
//type BillingKey struct {
//	APIResponse
//	Data BillingKeyData `json:"data"`
//}
type BillingKeyPayload struct {
	//RestConfig
	Pg             string         			 `json:"pg"`
	Method         string         			 `json:"method,omitempty"`
	OrderName      string         			 `json:"order_name"`
	SubscriptionId string         			 `json:"subscription_id"`

	CardNo          string         			 `json:"card_no"`
	CardPw          string         			 `json:"card_pw"`
	CardIdentityNo  string         			 `json:"card_identity_no"`
	CardExpireYear  string         			 `json:"card_expire_year"`
	CardExpireMonth string         			 `json:"card_expire_month"`
	Price    		float64 	   			 `json:"price"`

	Username                string           `json:"username"`
	AuthType                string           `json:"auth_type"`
	BankName                string           `json:"bank_name"`
	BankAccount             string           `json:"bank_account"`
	IdentityNo              string           `json:"identity_no"`
	CashReceiptType         string           `json:"cash_receipt_type"`
	CashReceiptIdentityNo   string           `json:"cash_receipt_identity_no"`
	Phone                   string           `json:"phone"`

	TaxFree    		float64 	   			 `json:"tax_free"`
	User            User           			 `json:"user"`
	Extra           SubscribeExtra 			 `json:"extra"`
	Metadata        map[string]interface{}   `json:"metadata"`
}
type SubscribeExtra struct {
    CardQuota           string  `json:"card_quota"`
	SubscribeTestPayment int    `json:"subscribe_test_payment"`
	RawData              int    `json:"raw_data"`
}
type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Gender   int    `json:"gender"`
	Area     string `json:"area"`
	Birth    string `json:"birth"`
	Addr     string `json:"addr"`
}
type Item struct {
	Name 	         string  `json:"name"`
	Qty              int     `json:"qty"`
	Id   	         string  `json:"id"`
	Price            float64 `json:"price"`
	Cat1             string  `json:"cat1"`
	Cat2             string  `json:"cat2"`
	Cat3             string  `json:"cat3"`
	CategoryType     string  `json:"category_type"`
	CategoryCode     string  `json:"category_code"`
	StartDate        string  `json:"start_date"`
	EndDate          string  `json:"end_date"`
}
type SubscribePayload struct {
	BillingKey       string                 `json:"billing_key"`
	OrderName        string                 `json:"order_name"`
	OrderId          string                 `json:"order_id"`
	Price            float64                `json:"price"`
	TaxFree          float64                `json:"tax_free,omitempty"`
	CardQuota        string                 `json:"card_quota,omitempty"`
	CardInterest     string                 `json:"card_interest,omitempty"`
	User             User                   `json:"user,omitempty"`
	Items            []Item                 `json:"items,omitempty"`
	FeedbackUrl      string                 `json:"feedback_url,omitempty"`  // webhook 통지시 받으실 url 주소 (localhost 사용 불가)
	ContentType      string                 `json:"content_type,omitempty"`  // webhook 통지시 받으실 데이터 타입 (application/json 또는 application/x-www-form-urlencoded 선택)
	Extra            SubscribeExtra         `json:"extra,omitempty"`
	ReserveExecuteAt string                 `json:"reserve_execute_at,omitempty"` //ex. "2022-04-21T17:05:36+09:00"
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

type RequestPayload struct {
	ReceiptID string `json:"receipt_id"`
}

//type SubscribeBilling struct {
//	APIResponse
//	Data    struct {
//		ReceiptId          string 					`json:"receipt_id"`
//		Price     		   float64    				`json:"price"`
//		CardNo     		   string    				`json:"card_no"`
//		CardCode     	   string    				`json:"card_code"`
//		CardName     	   string    				`json:"card_name"`
//		CardQuota     	   string    				`json:"card_quota"`
//		Params             map[string]interface{}   `json:"params"`
//		ItemName     	   string    				`json:"item_name"`
//		OrderId     	   string    				`json:"order_id"`
//		Url     		   string    				`json:"url"`
//		PaymentName        string    				`json:"payment_name"`
//		PgName     	  	   string    				`json:"pg_name"`
//		Pg     			   string    				`json:"pg"`
//		Method     		   string    				`json:"method"`
//		MethodName     	   string    				`json:"method_name"`
//		RequestedAt        string                   `json:"requested_at"`
//		PurchasedAt        string                   `json:"purchased_at"`
//		Status             int                      `json:"status"`
//	} `json:"data"`
//}

//type SubscribeBillingReserve struct {
//	APIResponse
//	Data    struct {
//		ReserveId          string 			`json:"reserve_id"`
//		BillingKey         string 			`json:"billing_key"`
//		ItemName           string 			`json:"item_name"`
//		Price              float64 			`json:"price"`
//		PurchaseCount      int 				`json:"purchase_count"`
//		PurchaseLimit      int 				`json:"purchase_limit"`
//		ExecuteAt          int64          	`json:"execute_at"`
//		Status             int              `json:"status"`
//	} `json:"data"`
//}

func (api *Api) GetBillingKey(payload BillingKeyPayload) (APIResponse, error) {
	//if payload.ApplicationId == "" {
	//	payload.ApplicationId = api.applicationId
	//}
	//if payload.PrivateKey == "" {
	//	payload.PrivateKey = api.privateKey
	//}

	postBody, _ := json.Marshal(payload)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/request/subscribe", body)
	if err != nil {
		errors.New("bootpay: GetBillingKey error: " + err.Error())
		return APIResponse{}, err
	} 
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode

	return result, nil
}


func (api *Api) LookupBillingKey(receiptId string) (APIResponse, error) {

	req, err := api.NewRequest(http.MethodGet, "/subscribe/billing_key/" + receiptId, nil)
	if err != nil {
		errors.New("bootpay: LookupBillingKey error: " + err.Error())
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode
	return result, nil
}

func (api *Api) LookupBillingKeyByKey(billingKey string) (APIResponse, error) {

	req, err := api.NewRequest(http.MethodGet, "/billing_key/" + billingKey, nil)
	if err != nil {
		errors.New("bootpay: LookupBillingKeyByKey error: " + err.Error())
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode
	return result, nil
}


func (api *Api) DestroyBillingKey(billingKey string) (APIResponse, error) {
	req, err := api.NewRequest(http.MethodDelete, "/subscribe/billing_key/" + billingKey, nil)
	if err != nil {
		errors.New("bootpay: DestroyBillingKey error: " + err.Error())
		return APIResponse{}, err
	} 
	res, err := api.client.Do(req)

	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode
	return result, nil
}

func (api *Api) RequestSubscribe(payload SubscribePayload) (APIResponse, error) {
	//if payload.ApplicationId == "" {
	//	payload.ApplicationId = api.applicationId
	//}
	//if payload.PrivateKey == "" {
	//	payload.PrivateKey = api.privateKey
	//}
	postBody, _ := json.Marshal(payload)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/subscribe/payment", body)
	if err != nil {
		errors.New("bootpay: RequestSubscribe error: " + err.Error())
		return APIResponse{}, err
	} 
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode
	return result, nil
}

func (api *Api) ReserveSubscribe(payload SubscribePayload) (APIResponse, error) {
	postBody, _ := json.Marshal(payload)
	body := bytes.NewBuffer(postBody)
	
	req, err := api.NewRequest(http.MethodPost, "/subscribe/payment/reserve", body)
	if err != nil {
		errors.New("bootpay: ReserveSubscribe error: " + err.Error())
		return APIResponse{}, err
	} 
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode
	return result, nil
}


func (api *Api) ReserveSubscribeLookup(reserveId string) (APIResponse, error) {	 

	req, err := api.NewRequest(http.MethodGet, "/subscribe/payment/reserve/" + reserveId, nil)
	if err != nil {
		errors.New("bootpay: ReserveSubscribe error: " + err.Error())
		return APIResponse{}, err
	} 
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode
	return result, nil
}

func (api *Api) ReserveCancelSubscribe(reserveId string) (APIResponse, error) {
	req, err := api.NewRequest(http.MethodDelete, "/subscribe/payment/reserve/" + reserveId, nil)
	if err != nil {
		errors.New("bootpay: ReserveCancelSubscribe error: " + err.Error())
		return APIResponse{}, err
	} 
	res, err := api.client.Do(req)

	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode
	return result, nil
}

func (api *Api) RequestSubscribeAutomaticTransferBillingKey(payload BillingKeyPayload) (APIResponse, error) {
    postBody, _ := json.Marshal(payload)
	body := bytes.NewBuffer(postBody)

    req, err := api.NewRequest(http.MethodPost, "/request/subscribe/automatic-transfer", body)
	if err != nil {
		errors.New("bootpay: requestSubscribeAutomaticTransferBillingKey error: " + err.Error())
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode

	return result, nil
}

func (api *Api) PublishAutomaticTransferBillingKey(receiptId string) (APIResponse, error) {
    payload := RequestPayload{
		ReceiptID: receiptId,
	}

    postBody, _ := json.Marshal(payload)
	body := bytes.NewBuffer(postBody)

    req, err := api.NewRequest(http.MethodPost, "/request/subscribe/automatic-transfer/publish", body)
	if err != nil {
		errors.New("bootpay: publishAutomaticTransferBillingKey error: " + err.Error())
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode

	return result, nil
}