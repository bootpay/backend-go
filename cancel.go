package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type CancelData struct {
	RestConfig
	ReceiptId string     `json:"receipt_id"`
	Price     float64    `json:"price,omitempty"`
	Name      string     `json:"name"`
	Reason    string     `json:"reason"`
	Refund    RefundData `json:"refund,omitempty"`
}
//type CancelPartData struct {
//	CancelData
//	Price     float64    `json:"price"`
//}
//type CancelRefundData struct {
//	CancelData
//	Refund    RefundData `json:"refund"`
//}

type RefundData struct {
	Account       string `json:"account"`
	Accountholder string `json:"accountholder"`
	Bankcode      string `json:"bankcode"`
}
type Cancel struct {
	APIResponse
	Data    struct {
		ReceiptId          string 		`json:"receipt_id"`
		RequestCancelPrice float64    	`json:"request_cancel_price"`
		RemainPrice        float64    	`json:"remain_price"`
		RemainTaxFree      float64    	`json:"remain_tax_free"`
		CancelledPrice     float64    	`json:"cancelled_price"`
		CancelledTaxFree   float64    	`json:"cancelled_tax_free"`
		RevokedAt          string 		`json:"revoked_at"`
		Tid                string 		`json:"tid"`
	} `json:"data"`
}
func (api *Api) ReceiptCancel(receiptId string, price float64, name string, reason string, refund RefundData) (Cancel, error) {
	cancel := make(map[string]interface{})
	cancel["application_id"] = api.applicationId
	cancel["private_key"] = api.privateKey
	cancel["receipt_id"] = receiptId
	if price != 0 {  cancel["price"] = price }
	cancel["name"] = name
	cancel["reason"] = reason
	if refund.Bankcode != "" { cancel["refund"] = refund }

	postBody, _ := json.Marshal(cancel)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/cancel", body)
	if err != nil {
		errors.New("bootpay: ReserveCancelSubscribe error: " + err.Error())
		return Cancel{}, err
	}
	req.Header.Set("Authorization", api.token)
	res, err := api.client.Do(req)

	defer res.Body.Close()

	result := Cancel{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}
