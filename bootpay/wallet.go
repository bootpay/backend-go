package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// WalletRequest represents the request parameters for wallet payment
type WalletRequest struct {
	UserId      string                 `json:"user_id"`
	OrderName   string                 `json:"order_name"`
	Price       float64                `json:"price"`
	TaxFree     float64                `json:"tax_free,omitempty"`
	OrderId     string                 `json:"order_id"`
	WebhookUrl  string                 `json:"webhook_url,omitempty"`
	ContentType string                 `json:"content_type,omitempty"` // "application/json" or "application/x-www-form-urlencoded"
	Items       []Item                 `json:"items,omitempty"`
	User        User                   `json:"user,omitempty"`
	Extra       SubscribeExtra         `json:"extra,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Sandbox     bool                   `json:"sandbox"`
}

// WalletDataPart represents a single wallet entry
type WalletDataPart struct {
	WalletId          string    `json:"wallet_id"`
	Type              int       `json:"type"`
	Sandbox           int       `json:"sandbox"`
	Order             int       `json:"order"`
	PaymentStatus     int       `json:"payment_status"`
	BatchData         BatchData `json:"batch_data"`
	CardCode          string    `json:"card_code"`
	ExpiredAt         string    `json:"expired_at"`
	LatestPurchasedAt string    `json:"latest_purchased_at"`
}

// BatchData represents card batch data
type BatchData struct {
	CardNo          string `json:"card_no"`
	CardCompany     string `json:"card_company"`
	CardCompanyCode string `json:"card_company_code"`
	CardType        int    `json:"card_type"`
	CardHash        string `json:"card_hash"`
}

// GetUserWallets retrieves the list of wallets for a user
func (api *Api) GetUserWallets(userId string, sandbox bool) (APIResponse, error) {
	sandboxStr := "false"
	if sandbox {
		sandboxStr = "true"
	}

	url := fmt.Sprintf("/wallet?user_id=%s&sandbox=%s", userId, sandboxStr)
	req, err := api.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		errors.New("bootpay: GetUserWallets error: " + err.Error())
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil {
		result = map[string]interface{}{}
	}
	result["http_status"] = res.StatusCode
	return result, nil
}

// RequestWalletPayment requests a payment using wallet
func (api *Api) RequestWalletPayment(request WalletRequest) (APIResponse, error) {
	postBody, _ := json.Marshal(request)
	body := bytes.NewBuffer(postBody)

	req, err := api.NewRequest(http.MethodPost, "/wallet/payment", body)
	if err != nil {
		errors.New("bootpay: RequestWalletPayment error: " + err.Error())
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil {
		result = map[string]interface{}{}
	}
	result["http_status"] = res.StatusCode
	return result, nil
}
