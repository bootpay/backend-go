package bootpay_go

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
	Name      string   `json:"name"`
	UserInfo  User     `json:"user_info"`
	Items     []Item   `json:"items"`
	ReturnUrl string   `json:"return_url"`
	Extra     Extra    `json:"extra""`
	//UserInfo            User           `json:"user_info"`
}

type Extra struct {
	Escrow               bool   `json:"escrow"`
	ExpireMonth          int    `json:"expire_month"`
	Quota                []int  `json:"quota"`
	SubscribeTestPayment bool   `json:"subscribe_test_payment"`
	DispCashResult       bool   `json:"disp_cash_result"`
	OfferPeriod          string `json:"offer_period"`
	SellerName           string `json:"seller_name"`
	Theme                string `json:"theme"`
	CustomBackground     string `json:"custom_background"`
	CustomFontColor      string `json:"custom_font_color"`
}

func (c *Client) RequestLink(payload Payload) (APIResponse, error) {
	if payload.ApplicationId == "" {
		payload.ApplicationId = c.applicationId
	}
	if payload.PrivateKey == "" {
		payload.PrivateKey = c.privateKey
	}
	postBody, _ := json.Marshal(payload)
	body := bytes.NewBuffer(postBody)
	req, err := c.NewRequest(http.MethodPost, "/request/payment", body)
	if err != nil {
		errors.New("bootpay: RequestLink error: " + err.Error())
		return APIResponse{}, err
	}
	req.Header.Set("Authorization", c.token)
	res, err := c.httpClient.Do(req)

	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}
