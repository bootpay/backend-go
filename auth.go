package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)


type Authentication struct {
	Pg             		      string         			 `json:"pg"`
	Method         			  string         			 `json:"method"`
	Username      			  string         			 `json:"username"`
	IdentityNo     			  string         			 `json:"identity_no"`

	Carrier        			  string         			 `json:"carrier"`
	Phone          			  string         			 `json:"phone"`

	SiteUrl        			  string         		     `json:"site_url"`
	OrderName      			  string         			 `json:"order_name"`
	AuthenticationId          string         			 `json:"authentication_id"`
}


type AuthenticationParams struct {
	ReceiptId      			  string         			 `json:"receipt_id"`
	Otp          			  string         			 `json:"otp,omitempty"`
}



func (api *Api) RequestAuthentication(authentication Authentication) (APIResponse, error) {
	postBody, _ := json.Marshal(authentication)
	body := bytes.NewBuffer(postBody)

	req, err := api.NewRequest(http.MethodPost, "/request/authentication", body)
	if err != nil {
		errors.New("bootpay: RequestAuthentication error: " + err.Error())
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


func (api *Api) ConfirmAuthentication(params AuthenticationParams) (APIResponse, error) {
	postBody, _ := json.Marshal(params)
	body := bytes.NewBuffer(postBody)

	req, err := api.NewRequest(http.MethodPost, "/authenticate/confirm", body)
	if err != nil {
		errors.New("bootpay: ConfirmAuthentication error: " + err.Error())
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



func (api *Api) RealarmAuthentication(params AuthenticationParams) (APIResponse, error) {
	postBody, _ := json.Marshal(params)
	body := bytes.NewBuffer(postBody)

	req, err := api.NewRequest(http.MethodPost, "/authenticate/realarm", body)
	if err != nil {
		errors.New("bootpay: RealarmAuthentication error: " + err.Error())
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
