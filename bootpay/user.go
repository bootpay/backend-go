package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// UserTokenRequest represents the request parameters for user token
type UserTokenRequest struct {
	UserId   string `json:"user_id"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Gender   int    `json:"gender,omitempty"`
	Birth    string `json:"birth,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

// RequestUserToken requests a user token for the given user
func (api *Api) RequestUserToken(request UserTokenRequest) (APIResponse, error) {
	postBody, _ := json.Marshal(request)
	body := bytes.NewBuffer(postBody)

	req, err := api.NewRequest(http.MethodPost, "/request/user/token", body)
	if err != nil {
		errors.New("bootpay: RequestUserToken error: " + err.Error())
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil {
		result = map[string]interface{}{}
	}
	result["http_status"] = res.StatusCode
	return result, nil
}
