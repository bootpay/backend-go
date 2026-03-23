package bootpay

import (
	"bytes"
	"encoding/json"
	"net/http"
)
type TokenData struct {
	Token      string `json:"access_token"`
	ServerTime int64  `json:"server_time"`
	ExpiredAt  int64  `json:"expired_at"`
}
//type Token struct {
//	Data interface{}
//	//ExpireIn int64 `json:"expire_in"`
//	//AccessToken string `json:"access_token"`
//}
func (api *Api) GetToken() (APIResponse, error) {
	config := RestConfig{api.applicationId, api.privateKey}
	postBody, _ := json.Marshal(config)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/request/token", body)

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

	if result["access_token"] != nil {
		api.token = result["access_token"].(string)
	}

	return result, nil
}
