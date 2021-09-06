package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type EasyUserTokenPayload struct {
	RestConfig
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Gender int    `json:"gender"`
	Birth  string `json:"birth"`
	Phone  string `json:"phone"`
}

type EasyUserToken struct {
	APIResponse
	Data    struct {
		UserToken        string `json:"user_token"`
		ExpiredUnixtime  int64  `json:"expired_unixtime"`
		ExpiredLocaltime string `json:"expired_localtime"`
	} `json:"data"`
}

func (c *Client) GetUserToken(userToken EasyUserTokenPayload) (EasyUserToken, error) {
	if userToken.ApplicationId == "" {
		userToken.ApplicationId = c.applicationId
	}
	if userToken.PrivateKey == "" {
		userToken.PrivateKey = c.privateKey
	}
	postBody, _ := json.Marshal(userToken)
	body := bytes.NewBuffer(postBody)
	req, err := c.NewRequest(http.MethodPost, "/request/user/token", body)
	if err != nil {
		errors.New("bootpay: ReserveCancelSubscribe error: " + err.Error())
		return EasyUserToken{}, err
	}
	req.Header.Set("Authorization", c.token)
	res, err := c.httpClient.Do(req)

	defer res.Body.Close()

	result := EasyUserToken{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}
