package bootpay

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type CashReceiptData struct {
	ReceiptId       string                 `json:"receipt_id"`
	OrderId         string                 `json:"order_id,omitempty"`
	OrderName       string                 `json:"order_name,omitempty"`
	IdentityNo      string                 `json:"identity_no,omitempty"`
	PurchasedAt     string                 `json:"purchased_at,omitempty"`
	CashReceiptType string                 `json:"cash_receipt_type,omitempty"` // "소득공제" or "지출증빙"
	Price           float64                `json:"price"`
	TaxFree         float64                `json:"tax_free,omitempty"`
	Currency        string                 `json:"currency,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	Extra           map[string]interface{} `json:"extra,omitempty"`
	/* 구매자 정보 */
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	/* 별건 요청 파라미터 */
	// Pg 는 선택값이다. 비워두면 키를 싣지 않고, 서버가 프로젝트 기본 PG 로 발행한다.
	// (Ruby SDK c716a1f — request_cash_receipt 의 pg 가 필수에서 선택으로 바뀐 것과 동일한 계약)
	Pg   string `json:"pg,omitempty"`
	User User   `json:"user,omitempty"`
}


func (api *Api) RequestCashReceiptByBootpay(cashReceipt CashReceiptData) (APIResponse, error) {

	postBody, _ := json.Marshal(cashReceipt)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/request/receipt/cash/publish", body)
	if err != nil {
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode
	return result, nil
}

func (api *Api) RequestCashReceiptCancelByBootpay(cancel CancelData) (APIResponse, error) {

	postBody, _ := json.Marshal(cancel)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodDelete, "/request/receipt/cash/cancel/" + cancel.ReceiptId, body)
	if err != nil {
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)

	if result == nil { result =  map[string]interface{}{} }

	result["http_status"] = res.StatusCode
	return result, nil
}

func (api *Api) RequestCashReceipt(cashReceipt CashReceiptData) (APIResponse, error) {
	postBody, _ := json.Marshal(cashReceipt)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/request/cash/receipt", body)
	if err != nil {
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode
	return result, nil
}

func (api *Api) RequestCashReceiptCancel(cancel CancelData) (APIResponse, error) {

	postBody, _ := json.Marshal(cancel)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodDelete, "/request/cash/receipt/" + cancel.ReceiptId, body)
	if err != nil {
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode
	return result, nil
}
