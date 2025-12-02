package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
)

// UserModule handles user-related operations
type UserModule struct {
	api *CommerceApi
}

// Token issues a user token
func (m *UserModule) Token(userId string) (map[string]interface{}, error) {
	data := map[string]string{
		"user_id": userId,
	}
	return m.api.Post("users/login/token", data)
}

// Join creates a new user
func (m *UserModule) Join(user CommerceUser) (map[string]interface{}, error) {
	return m.api.Post("users/join", user)
}

// CheckExist checks if a user exists by key and value
func (m *UserModule) CheckExist(key string, value string) (map[string]interface{}, error) {
	encodedValue := url.QueryEscape(value)
	return m.api.Get(fmt.Sprintf("users/join/%s?pk=%s", key, encodedValue))
}

// AuthenticationData retrieves authentication data by standId
func (m *UserModule) AuthenticationData(standId string) (map[string]interface{}, error) {
	return m.api.Get(fmt.Sprintf("users/authenticate/%s", standId))
}

// Login performs user login
func (m *UserModule) Login(loginId string, loginPw string) (map[string]interface{}, error) {
	data := map[string]string{
		"login_id": loginId,
		"login_pw": loginPw,
	}
	return m.api.Post("users/login", data)
}

// List retrieves user list
func (m *UserModule) List(params *UserListParams) (map[string]interface{}, error) {
	query := ""
	if params != nil {
		queryParams := url.Values{}
		if params.Page > 0 {
			queryParams.Set("page", strconv.Itoa(params.Page))
		}
		if params.Limit > 0 {
			queryParams.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Keyword != "" {
			queryParams.Set("keyword", params.Keyword)
		}
		if params.MemberType > 0 {
			queryParams.Set("member_type", strconv.Itoa(params.MemberType))
		}
		if params.Type != "" {
			queryParams.Set("type", params.Type)
		}
		if len(queryParams) > 0 {
			query = "?" + queryParams.Encode()
		}
	}
	return m.api.Get("users" + query)
}

// Detail retrieves user details
func (m *UserModule) Detail(userId string) (map[string]interface{}, error) {
	return m.api.Get(fmt.Sprintf("users/%s", userId))
}

// Update updates user information
func (m *UserModule) Update(user CommerceUser) (map[string]interface{}, error) {
	if user.UserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	return m.api.Put(fmt.Sprintf("users/%s", user.UserId), user)
}

// Delete deletes a user
func (m *UserModule) Delete(userId string) (map[string]interface{}, error) {
	return m.api.Delete(fmt.Sprintf("users/%s", userId))
}
