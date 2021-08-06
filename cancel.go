package bootpay_go

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
		ReceiptId          string `json:"receipt_id"`
		RequestCancelPrice int    `json:"request_cancel_price"`
		RemainPrice        int    `json:"remain_price"`
		RemainTaxFree      int    `json:"remain_tax_free"`
		CancelledPrice     int    `json:"cancelled_price"`
		CancelledTaxFree   int    `json:"cancelled_tax_free"`
		RevokedAt          string `json:"revoked_at"`
		Tid                string `json:"tid"`
	} `json:"data"`
}
func (c *Client) ReceiptCancel(receiptId string, price float64, name string, reason string, refund RefundData) (Cancel, error) {
	cancel := make(map[string]interface{})
	cancel["application_id"] = c.applicationId
	cancel["private_key"] = c.privateKey
	cancel["receipt_id"] = receiptId
	if price != 0 {  cancel["price"] = price }
	cancel["name"] = name
	cancel["reason"] = reason
	if refund.Bankcode != "" { cancel["refund"] = refund }

	postBody, _ := json.Marshal(cancel)
	body := bytes.NewBuffer(postBody)
	req, err := c.NewRequest(http.MethodPost, "/cancel", body)
	if err != nil {
		errors.New("bootpay: ReserveCancelSubscribe error: " + err.Error())
		return Cancel{}, err
	}
	req.Header.Set("Authorization", c.token)
	res, err := c.httpClient.Do(req)

	defer res.Body.Close()

	result := Cancel{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}
