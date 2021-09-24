package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type BillingKeyCardData struct {
	CardCode string `json:"card_code"`
	CardName string `json:"card_name"`
	CardNo   string `json:"card_no"`
	CardCl   string `json:"card_cl"`
}
type BillingKeyData struct {
	BillingKey string             `json:"billing_key"`
	PgName     string             `json:"pg_name"`
	MethodName string             `json:"method_name"`
	method     string             `json:"method"`
	Data       BillingKeyCardData `json:"data"`
	EndAt      string             `json:"e_at"`
	CreateAt   string             `json:"c_at"`
}
type BillingKey struct {
	APIResponse
	Data BillingKeyData `json:"data"`
}
type BillingKeyPayload struct {
	RestConfig
	OrderId        string         `json:"order_id"`
	Pg             string         `json:"pg"`
	ItemName       string         `json:"item_name"`
	CardNo         string         `json:"card_no"`
	CardPw         string         `json:"card_pw"`
	ExpireYear     string         `json:"expire_year"`
	ExpireMonth    string         `json:"expire_month"`
	IdentifyNumber string         `json:"identify_number"`
	UserInfo       User           `json:"user_info"`
	Extra          SubscribeExtra `json:"extra"`
}
type SubscribeExtra struct {
	SubscribeTestPayment int `json:"subscribe_test_payment"`
	RawData              int `json:"raw_data"`
}
type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Gender   int    `json:"gender"`
	Area     string `json:"area"`
	Birth    string `json:"birth"`
}
type Item struct {
	ItemName string  `json:"item_name"`
	Qty      int     `json:"qty"`
	Unique   string  `json:"unique"`
	Price    float64 `json:"price"`
	Cat1     string  `json:"cat1"`
	Cat2     string  `json:"cat2"`
	Cat3     string  `json:"cat3"`
}
type SubscribePayload struct {
	RestConfig
	BillingKey          string         `json:"billing_key"`
	ItemName            string         `json:"item_name"`
	Price               float64        `json:"price"`
	TaxFree             float64        `json:"tax_free"`
	OrderId             string         `json:"order_id"`
	Quota               int            `json:"quota"`
	Interest            int            `json:"interest"`
	UserInfo            User           `json:"user_info"`
	FeedbackUrl         string         `json:"feedback_url"`
	FeedbackContentType string         `json:"feedback_content_type"`
	Extra               SubscribeExtra `json:"extra"`
	SchedulerType       string         `json:"scheduler_type"`
	ExecuteAt           int64          `json:"execute_at"`
}

type SubscribeBilling struct {
	APIResponse
	Data    struct {
		ReceiptId          string 					`json:"receipt_id"`
		Price     		   float64    				`json:"price"`
		CardNo     		   string    				`json:"card_no"`
		CardCode     	   string    				`json:"card_code"`
		CardName     	   string    				`json:"card_name"`
		CardQuota     	   string    				`json:"card_quota"`
		Params             map[string]interface{}   `json:"params"`
		ItemName     	   string    				`json:"item_name"`
		OrderId     	   string    				`json:"order_id"`
		Url     		   string    				`json:"url"`
		PaymentName        string    				`json:"payment_name"`
		PgName     	  	   string    				`json:"pg_name"`
		Pg     			   string    				`json:"pg"`
		Method     		   string    				`json:"method"`
		MethodName     	   string    				`json:"method_name"`
		RequestedAt        string                   `json:"requested_at"`
		PurchasedAt        string                   `json:"purchased_at"`
		Status             int                      `json:"status"`
	} `json:"data"`
}

type SubscribeBillingReserve struct {
	APIResponse
	Data    struct {
		ReserveId          string 			`json:"reserve_id"` 
		BillingKey         string 			`json:"billing_key"` 
		ItemName           string 			`json:"item_name"` 
		Price              float64 			`json:"price"` 
		PurchaseCount      int 				`json:"purchase_count"` 
		PurchaseLimit      int 				`json:"purchase_limit"` 
		ExecuteAt          int64          	`json:"execute_at"`
		Status             int              `json:"status"`
	} `json:"data"`
}

func (api *Api) GetBillingKey(payload BillingKeyPayload) (BillingKey, error) {
	if payload.ApplicationId == "" {
		payload.ApplicationId = api.applicationId
	}
	if payload.PrivateKey == "" {
		payload.PrivateKey = api.privateKey
	}

	postBody, _ := json.Marshal(payload)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/request/card_rebill", body)
	if err != nil {
		errors.New("bootpay: GetBillingKey error: " + err.Error())
		return BillingKey{}, err
	}
	req.Header.Set("Authorization", api.token)
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := BillingKey{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}
func (api *Api) DestroyBillingKey(billingKey string) (APIResponse, error) {
	req, err := api.NewRequest(http.MethodDelete, "/subscribe/billing/"+billingKey, nil)
	if err != nil {
		errors.New("bootpay: DestroyBillingKey error: " + err.Error())
		return APIResponse{}, err
	}
	req.Header.Set("Authorization", api.token)
	res, err := api.client.Do(req)

	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}
func (api *Api) RequestSubscribe(payload SubscribePayload) (SubscribeBilling, error) {
	if payload.ApplicationId == "" {
		payload.ApplicationId = api.applicationId
	}
	if payload.PrivateKey == "" {
		payload.PrivateKey = api.privateKey
	}
	postBody, _ := json.Marshal(payload)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/subscribe/billing", body)
	if err != nil {
		errors.New("bootpay: RequestSubscribe error: " + err.Error())
		return SubscribeBilling{}, err
	}
	req.Header.Set("Authorization", api.token)
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := SubscribeBilling{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}

func (api *Api) ReserveSubscribe(payload SubscribePayload) (SubscribeBillingReserve, error) {
	if payload.ApplicationId == "" {
		payload.ApplicationId = api.applicationId
	}
	if payload.PrivateKey == "" {
		payload.PrivateKey = api.privateKey
	}
	if payload.SchedulerType == "" {
		payload.SchedulerType = "oneshot"
	}

	postBody, _ := json.Marshal(payload)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/subscribe/billing/reserve", body)
	if err != nil {
		errors.New("bootpay: ReserveSubscribe error: " + err.Error())
		return SubscribeBillingReserve{}, err
	}
	req.Header.Set("Authorization", api.token)
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := SubscribeBillingReserve{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}

func (api *Api) ReserveCancelSubscribe(reserveId string) (APIResponse, error) {
	req, err := api.NewRequest(http.MethodDelete, "/subscribe/billing/reserve/"+reserveId, nil)
	if err != nil {
		errors.New("bootpay: ReserveCancelSubscribe error: " + err.Error())
		return APIResponse{}, err
	}
	req.Header.Set("Authorization", api.token)
	res, err := api.client.Do(req)

	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}
