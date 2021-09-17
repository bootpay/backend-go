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

func (bootpay *Bootpay) ServerSubmit(receiptId string) (APIResponse, error) {
	sub := Submit{}
	if sub.ApplicationId == "" {
		sub.ApplicationId = bootpay.applicationId
	}
	if sub.PrivateKey == "" {
		sub.PrivateKey = bootpay.privateKey
	}
	sub.ReceiptId = receiptId
	postBody, _ := json.Marshal(sub)
	body := bytes.NewBuffer(postBody)
	req, err := bootpay.NewRequest(http.MethodPost, "/submit", body)
	if err != nil {
		errors.New("bootpay: Submit error: " + err.Error())
		return APIResponse{}, err
	}
	req.Header.Set("Authorization", bootpay.token)
	res, err := bootpay.client.Do(req)

	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}
