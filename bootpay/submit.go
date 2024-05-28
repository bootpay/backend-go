package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Submit struct {
	RestConfig
	ReceiptId string `json:"receipt_id"`
}

func (api *Api) ServerSubmit(receiptId string) (APIResponse, error) {
	sub := Submit{}
	if sub.ApplicationId == "" {
		sub.ApplicationId = api.applicationId
	}
	if sub.PrivateKey == "" {
		sub.PrivateKey = api.privateKey
	}
	sub.ReceiptId = receiptId
	postBody, _ := json.Marshal(sub)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/submit", body)
	if err != nil {
		errors.New("bootpay: Submit error: " + err.Error())
		return APIResponse{}, err
	}
	req.Header.Set("Authorization", api.token)
	res, err := api.client.Do(req)

	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode
	return result, nil
}
