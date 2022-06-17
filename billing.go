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
	taxFree    		float64 	   			 `json:"tax_free"`
	User            User           			 `json:"user"`
	Extra           SubscribeExtra 			 `json:"extra"`
	Metadata        map[string]interface{}   `json:"metadata"`
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
	Name 	 string  `json:"name"`
	Qty      int     `json:"qty"`
	Id   	 string  `json:"id"`
	Price    float64 `json:"price"`
	Cat1     string  `json:"cat1"`
	Cat2     string  `json:"cat2"`
	Cat3     string  `json:"cat3"`
}
type SubscribePayload struct {
	BillingKey          string         `json:"billing_key"`
	OrderName           string         `json:"order_name"`
	OrderId             string         `json:"order_id"`
	Price               float64        `json:"price"`
	TaxFree             float64        `json:"tax_free"`

	CardQuota           string         `json:"card_quota"`
	CardInterest        string         `json:"card_interest"`
	User            	User           `json:"user"`
	FeedbackUrl         string         `json:"feedback_url"` // webhook 통지시 받으실 url 주소 (localhost 사용 불가)
	ContentType 		string         `json:"content_type"` // webhook 통지시 받으실 데이터 타입 (application/json 또는 application/x-www-form-urlencoded 선택)
	Extra               SubscribeExtra `json:"extra"`
	//SchedulerType       string         `json:"scheduler_type"`
	ReserveExecuteAt    string         `json:"reserve_execute_at"` //ex. "2022-04-21T17:05:36+09:00"
	Metadata            map[string]interface{}   `json:"metadata"`
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
	if payload.ApplicationId == "" {
		payload.ApplicationId = api.applicationId
	}
	if payload.PrivateKey == "" {
		payload.PrivateKey = api.privateKey
	}

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
	return result, nil
}

func (api *Api) DestroyBillingKey(billingKey string) (APIResponse, error) {
	req, err := api.NewRequest(http.MethodDelete, "/subscribe/billing/" + billingKey, nil)
	if err != nil {
		errors.New("bootpay: DestroyBillingKey error: " + err.Error())
		return APIResponse{}, err
	} 
	res, err := api.client.Do(req)

	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}

func (api *Api) RequestSubscribe(payload SubscribePayload) (APIResponse, error) {
	if payload.ApplicationId == "" {
		payload.ApplicationId = api.applicationId
	}
	if payload.PrivateKey == "" {
		payload.PrivateKey = api.privateKey
	}
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
	return result, nil
}

func (api *Api) ReserveSubscribe(payload SubscribePayload) (APIResponse, error) {
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
	req, err := api.NewRequest(http.MethodPost, "/subscribe/payment/reserve", body)
	if err != nil {
		errors.New("bootpay: ReserveSubscribe error: " + err.Error())
		return APIResponse{}, err
	} 
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
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
	return result, nil
}
