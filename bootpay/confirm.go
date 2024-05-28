package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
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
		errors.New("bootpay: ServerConfirm error: " + err.Error())
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
