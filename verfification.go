package bootpay

import (
	"encoding/json"
	"errors"
	"net/http"
)

//type VerifyData struct {
//	ReceiptId        string                 `json:"receipt_id"`
//	OrderId          string                 `json:"order_id"`
//	Name             string                 `json:"name"`
//	Price            float64                `json:"price"`
//	TaxFree          float64                `json:"tax_free"`
//	RemainPrice      float64                `json:"remain_price"`
//	RemainTaxFree    float64                `json:"remain_tax_free"`
//	CancelledPrice   float64                `json:"cancelled_price"`
//	CancelledTaxFree float64                `json:"cancelled_tax_free"`
//	ReceiptUrl       string                 `json:"receipt_url"`
//	Unit             string                 `json:"unit"`
//	Pg               string                 `json:"pg"`
//	Method           string                 `json:"method"`
//	PgName           string                 `json:"pg_name"`
//	MethodName       string                 `json:"method_name"`
//	Params           map[string]interface{} `json:"params"`
//	//PaymentData      map[string]interface{} `json:"payment_data"`
//	PaymentData struct {
//		Username  string `json:"username"`
//		Phone     string `json:"phone"`
//		Birth     string `json:"birth"`
//		Gender    string `json:"gender"`
//		Unique    string `json:"unique"`
//		Di        string `json:"di"`
//		ReceiptId string `json:"receipt_id"`
//		N         string `json:"n"`
//		P         int    `json:"p"`
//		Tid       string `json:"tid"`
//		Pg        string `json:"pg"`
//		Pm        string `json:"pm"`
//		PgA       string `json:"pg_a"`
//		PmA       string `json:"pm_a"`
//		OId       string `json:"o_id"`
//		PAt       string `json:"p_at"`
//		S         int    `json:"s"`
//		G         int    `json:"g"`
//	} `json:"payment_data"`
//
//	RequestedAt      string                 `json:"requested_at"`
//	PurchasedAt      string                 `json:"purchased_at"`
//	EscrowStatus     int                    `json:"escrow_status"`
//	EscrowStatusEn   string         		`json:"escrow_status_en"`
//	EscrowStatusKo   string         		`json:"escrow_status_ko"`
//	Status           int                    `json:"status"`
//	StatusEn         string                 `json:"status_en"`
//	StatusKo         string                 `json:"status_ko"`
//}
//type Verify struct {
//	APIResponse
//	Data map[string]interface{} `json:"data"`
//	//Data VerifyData `json:"data"`
//}

func (api *Api) GetReceipt(receiptId string) (APIResponse, error) {
	req, err := api.NewRequest(http.MethodGet, "/receipt/" + receiptId, nil)
	if err != nil {
		errors.New("bootpay: Verify error: " + err.Error())
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

//type Certificate struct {
//	APIResponse
//	Data struct {
//		ReceiptId   string                 `json:"receipt_id"`
//		OrderId     string                 `json:"order_id"`
//		Pg          string                 `json:"pg"`
//		Method      string                 `json:"method"`
//		PgName      string                 `json:"pg_name"`
//		MethodName  string                 `json:"method_name"`
//		//Certificate map[string]interface{} `json:"certificate"`
//		Certificate struct {
//			Username  string `json:"username"`
//			Phone     string `json:"phone"`
//			Birth     string `json:"birth"`
//			Gender    string `json:"gender"`
//			Unique    string `json:"unique"`
//			Di        string `json:"di"`
//			Foreigner bool   `json:"foreigner"`
//			Carrier string   `json:"carrier"`
//		}
//		PaymentData struct {
//			Username  string `json:"username"`
//			Phone     string `json:"phone"`
//			Birth     string `json:"birth"`
//			Gender    string `json:"gender"`
//			Unique    string `json:"unique"`
//			Di        string `json:"di"`
//			ReceiptId string `json:"receipt_id"`
//			N         string `json:"n"`
//			P         int    `json:"p"`
//			Tid       string `json:"tid"`
//			Pg        string `json:"pg"`
//			Pm        string `json:"pm"`
//			PgA       string `json:"pg_a"`
//			PmA       string `json:"pm_a"`
//			OId       string `json:"o_id"`
//			PAt       string `json:"p_at"`
//			S         int    `json:"s"`
//			G         int    `json:"g"`
//		} `json:"payment_data"`
//	} `json:"data"`
//}

func (api *Api) Certificate(receiptId string) (APIResponse, error) {
	req, err := api.NewRequest(http.MethodGet, "/certificate/" + receiptId, nil)
	if err != nil {
		errors.New("bootpay: Certificate error: " + err.Error())
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
