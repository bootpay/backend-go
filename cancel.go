package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type CancelData struct {
	//RestConfig
	ReceiptId       string     `json:"receipt_id"`
	CancelPrice     float64    `json:"cancel_price,omitempty"`
	CancelTaxFree   float64    `json:"cancel_tax_free,omitempty"`
	CancelId        string     `json:"cancel_id,omitempty"`
	CancelUsername  string     `json:"cancel_username"`
	CancelMessage   string     `json:"cancel_message"`
	//Name      		string     `json:"name"`
	//Reason    		string     `json:"reason"`
	Refund    		RefundData `json:"refund,omitempty"`
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
	BankAccount       string `json:"bank_account"`
	BankUsername      string `json:"bank_username"`
	Bankcode          string `json:"bankcode"`
}
//type Cancel struct {
//	APIResponse
//	Data    struct {
//		ReceiptId          string 		`json:"receipt_id"`
//		RequestCancelPrice float64    	`json:"request_cancel_price"`
//		RemainPrice        float64    	`json:"remain_price"`
//		RemainTaxFree      float64    	`json:"remain_tax_free"`
//		CancelledPrice     float64    	`json:"cancelled_price"`
//		CancelledTaxFree   float64    	`json:"cancelled_tax_free"`
//		RevokedAt          string 		`json:"revoked_at"`
//		Tid                string 		`json:"tid"`
//	} `json:"data"`
//}
//func (api *Api) ReceiptCancel(receiptId string, price float64, name string, reason string, refund RefundData) (APIResponse, error) {
func (api *Api) ReceiptCancel(cancelData CancelData) (APIResponse, error) {
	//cancel := make(map[string]interface{})
	//cancel["application_id"] = api.applicationId
	//cancel["private_key"] = api.privateKey
	//cancel["receipt_id"] = receiptId
	//if price != 0 {  cancel["price"] = price }
	//cancel["name"] = name
	//cancel["reason"] = reason
	//if refund.Bankcode != "" { cancel["refund"] = refund }

	postBody, _ := json.Marshal(cancelData)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/cancel", body)
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
