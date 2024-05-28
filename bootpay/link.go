package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Payload struct {
	RestConfig
	Pg        string   `json:"pg"`
	Method    string   `json:"method"`
	Methods   []string `json:"methods"`
	Price     float64  `json:"price"`
	OrderId   string   `json:"order_id"`
	Params    string   `json:"params"`
	TaxFree   float64  `json:"tax_free"`
	OrderName string   `json:"order_name"`
	User      User     `json:"user"`
	Items     []Item   `json:"items"`
	//ReturnUrl string   `json:"return_url"`
	Extra     Extra    `json:"extra""`
	//UserInfo            User           `json:"user_info"`
}

type Extra struct {
	CardQuota            string `json:"card_quota"`
	SellerName           string `json:"seller_name"`
	DeliveryDay          int    `json:"delivery_day"`
	Locale          	 string `json:"locale"`
	OfferPeriod          string `json:"offer_period"`
	DisplayCashReceipt   bool   `json:"display_cash_receipt"`
	DepositExpiration    string `json:"deposit_expiration"`
	AppScheme    		 string `json:"app_scheme"` //서버에선 사용되지 않음
	UseCardPoint         bool   `json:"use_card_point"`
	DirectCard           bool   `json:"direct_card"`
	UseOrderId           bool   `json:"use_order_id"`
	InternationalCardOnly bool  `json:"international_card_only"`
	DirectAppCard        bool   `json:"direct_app_card"`
	DirectSamsungpay     bool   `json:"direct_samsungpay"`
	EnableErrorWebhook   bool   `json:"enable_error_webhook"`
	SeparatelyConfirmed  bool   `json:"separately_confirmed"`
	ConfirmOnlyRestApi   bool   `json:"confirm_only_rest_api"`
	OpenType   			 string `json:"open_type"`
	UseBootpayInappSdk   bool   `json:"use_bootpay_inapp_sdk"`
	RedirectUrl   		 string `json:"redirect_url"`
	DisplaySuccessResult bool   `json:"display_success_result"`
	DisplayErrorResult   bool   `json:"display_error_result"`
	IsposableCupDeposit  int    `json:"disposable_cup_deposit"`

	CardEasyOption       BootExtraCardEasyOption `json:"card_easy_option"`
	BrowserOpenType      []BrowserOpenType   `json:"browser_open_type"`

	UseWelcomepayment    bool   `json:"use_welcomepayment"`
}

type BootExtraCardEasyOption struct {
	Title        string   `json:"title"`
}


type BrowserOpenType struct {
	Browser        string    `json:"browser"`
	OpenType       string   `json:"open_type"`
}


func (api *Api) RequestLink(payload Payload) (APIResponse, error) {
	if payload.ApplicationId == "" {
		payload.ApplicationId = api.applicationId
	}
	if payload.PrivateKey == "" {
		payload.PrivateKey = api.privateKey
	}
	postBody, _ := json.Marshal(payload)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/request/payment", body)
	if err != nil {
		errors.New("bootpay: RequestLink error: " + err.Error())
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
