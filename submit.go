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

func (c *Client) ServerSubmit(receiptId string) (APIResponse, error) {
	sub := Submit{}
	if sub.ApplicationId == "" {
		sub.ApplicationId = c.applicationId
	}
	if sub.PrivateKey == "" {
		sub.PrivateKey = c.privateKey
	}
	sub.ReceiptId = receiptId
	postBody, _ := json.Marshal(sub)
	body := bytes.NewBuffer(postBody)
	req, err := c.NewRequest(http.MethodPost, "/submit", body)
	if err != nil {
		errors.New("bootpay: Submit error: " + err.Error())
		return APIResponse{}, err
	}
	req.Header.Set("Authorization", c.token)
	res, err := c.httpClient.Do(req)

	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}
