package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)
type TokenData struct {
	Token      string `json:"token"`
	ServerTime int64  `json:"server_time"`
	ExpiredAt  int64  `json:"expired_at"`
}
type Token struct {
	APIResponse
	Data TokenData `json:"data"`
}
func (bootpay *Bootpay) GetToken() (Token, error) {
	config := RestConfig{bootpay.applicationId, bootpay.privateKey}
	postBody, _ := json.Marshal(config)
	body := bytes.NewBuffer(postBody)
	req, err := bootpay.NewRequest(http.MethodPost, "/request/token", body)
	if err != nil {
		errors.New("bootpay: getToken error: " + err.Error())
		return Token{}, err
	}
	res, err := bootpay.client.Do(req)
	defer res.Body.Close()

	result := Token{}
	json.NewDecoder(res.Body).Decode(&result)
	fmt.Println(result)

	if result.Status == 200 {
		if result.Data.Token != "" { bootpay.token = result.Data.Token }
	}

	return result, nil
}
