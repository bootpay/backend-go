package bootpay

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type TokenData struct {
	Token      string `json:"access_token"`
	ServerTime int64  `json:"server_time"`
	ExpiredAt  int64  `json:"expired_at"`
}

//	type Token struct {
//		Data interface{}
//		//ExpireIn int64 `json:"expire_in"`
//		//AccessToken string `json:"access_token"`
//	}

// GetAccessToken issues a PG access token — an alias of GetToken.
//
// 26-08-24: added because GetToken has opposite meanings across the two API surfaces.
// Api.GetToken() *issues* a token while CommerceApi.GetToken() merely *reads* the stored
// one. This name matches CommerceApi.GetAccessToken() and the other language SDKs.
// GetToken is kept as-is for backward compatibility.
func (api *Api) GetAccessToken() (APIResponse, error) {
	return api.GetToken()
}

func (api *Api) GetToken() (APIResponse, error) {
	if err := api.validateCredentials(); err != nil {
		return APIResponse{}, err
	}
	// client_key/secret_key 인증은 매 요청에 Basic Auth 헤더가 자동 부착되므로
	// access token 발급이 불필요하다. 호환을 위해 합성 응답을 반환한다.
	if api.clientKey != "" && api.secretKey != "" {
		api.token = ""
		return APIResponse{
			"access_token": "",
			"expire_in":    0,
			"http_status":  http.StatusOK,
		}, nil
	}

	config := RestConfig{}
	config.ApplicationId = api.applicationId
	config.PrivateKey = api.privateKey
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
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil && err != io.EOF {
		return APIResponse{}, err
	}
	if result == nil {
		result = map[string]interface{}{}
	}
	result["http_status"] = res.StatusCode

	if accessToken, ok := result["access_token"].(string); ok {
		api.token = accessToken
	} else {
		api.token = ""
	}

	return result, nil
}
