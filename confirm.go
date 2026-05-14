package bootpay

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Confirm struct {
	//RestConfig
	ReceiptId string `json:"receipt_id"`
}

func (api *Api) ServerConfirm(receiptId string) (APIResponse, error) {
	sub := Confirm{}
	//if sub.ApplicationId == "" {
	//	sub.ApplicationId = api.applicationId
	//}
	//if sub.PrivateKey == "" {
	//	sub.PrivateKey = api.privateKey
	//}
	sub.ReceiptId = receiptId
	postBody, _ := json.Marshal(sub)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/confirm", body)
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
