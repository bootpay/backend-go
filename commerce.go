package bootpay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (api *Api) commerceGet(path string) (APIResponse, error) {
	req, err := api.NewRequest(http.MethodGet, path, nil)
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
	result["http_status_code"] = res.StatusCode
	return result, nil
}

func (api *Api) commerceWrite(method string, path string, payload interface{}) (APIResponse, error) {
	postBody, _ := json.Marshal(payload)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(method, path, body)
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
	result["http_status_code"] = res.StatusCode
	return result, nil
}

// Store
func (api *Api) GetStore() (APIResponse, error) {
	return api.commerceGet("/store")
}

func (api *Api) StoreInfo() (APIResponse, error) {
	return api.GetStore()
}

func (api *Api) GetStoreDetail() (APIResponse, error) {
	return api.commerceGet("/store/detail")
}

func (api *Api) StoreDetail() (APIResponse, error) {
	return api.GetStoreDetail()
}

// User
func (api *Api) UserLogin(loginId string, loginPw string) (APIResponse, error) {
	return api.commerceWrite(http.MethodPost, "/users/login", APIResponse{
		"login_id": loginId,
		"login_pw": loginPw,
	})
}

func (api *Api) UserJoin(user APIResponse) (APIResponse, error) {
	return api.commerceWrite(http.MethodPost, "/users/join", user)
}

func (api *Api) UserJoinCheck(checkType string, pk string) (APIResponse, error) {
	encoded := url.QueryEscape(pk)
	return api.commerceGet(fmt.Sprintf("/users/join/%s?pk=%s", checkType, encoded))
}

// Mall aliases
func (api *Api) MallUserLogin(loginId string, loginPw string) (APIResponse, error) {
	return api.UserLogin(loginId, loginPw)
}

func (api *Api) MallUserJoin(user APIResponse) (APIResponse, error) {
	return api.UserJoin(user)
}

func (api *Api) MallUserJoinCheck(checkType string, pk string) (APIResponse, error) {
	return api.UserJoinCheck(checkType, pk)
}

// Product
func (api *Api) ProductList(params APIResponse) (APIResponse, error) {
	q := url.Values{}
	if params != nil {
		if v, ok := params["page"]; ok {
			q.Set("page", fmt.Sprintf("%v", v))
		}
		if v, ok := params["limit"]; ok {
			q.Set("limit", fmt.Sprintf("%v", v))
		}
		if v, ok := params["keyword"]; ok {
			q.Set("keyword", fmt.Sprintf("%v", v))
		}
		if v, ok := params["type"]; ok {
			switch t := v.(type) {
			case int:
				q.Set("type", strconv.Itoa(t))
			default:
				q.Set("type", fmt.Sprintf("%v", v))
			}
		}
	}
	path := "/products"
	if encoded := q.Encode(); encoded != "" {
		path = path + "?" + encoded
	}
	return api.commerceGet(path)
}

func (api *Api) ProductDetail(productId string) (APIResponse, error) {
	return api.commerceGet("/products/" + productId)
}

// Mall aliases
func (api *Api) Products(params APIResponse) (APIResponse, error) {
	return api.ProductList(params)
}

func (api *Api) MallProductDetail(productId string) (APIResponse, error) {
	return api.ProductDetail(productId)
}
