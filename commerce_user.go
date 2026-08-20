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

// UserLogin performs a V1 Mall API member login
// POST /v1/users/login (v1/users/login#create)
// ⚠️ v1 has no singular user/* routes, and POST /v1/users/session has no create action — do not send there.
// ⚠️ The server (LoginService) reads only login_id/password; corporate_type is ignored even when sent.
func (m *UserModule) UserLogin(params MallUserLoginParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("users/login", params, commerceMallHeaders("", params.IdempotencyKey))
}

// UserSession retrieves the member session (V1 Mall API)
// GET /v1/users/session
// userJwt is the member JWT issued at login (attached as Bootpay-User-JWT only when present).
// idempotencyKey is auto-generated when empty.
func (m *UserModule) UserSession(userJwt string, idempotencyKey string) (map[string]interface{}, error) {
	return m.api.getWithHeaders("users/session", commerceMallHeaders(userJwt, idempotencyKey))
}

// UserLogout logs the member out (V1 Mall API)
// DELETE /v1/users/session
func (m *UserModule) UserLogout(userJwt string, idempotencyKey string) (map[string]interface{}, error) {
	return m.api.deleteWithHeaders("users/session", nil, commerceMallHeaders(userJwt, idempotencyKey))
}

// UserJoin joins a member (V1 Mall API — general signup)
// POST /v1/users/join
// ⚠️ Calls the same endpoint as Join(user) but with a different purpose —
// this one is the general signup using password/corporate_type/group, while Join is the
// external-uid signup using uid/login_email/login_pw. The server branches on the parameter set.
func (m *UserModule) UserJoin(params MallUserJoinParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("users/join", params, commerceMallHeaders("", params.IdempotencyKey))
}

// UserJoinCheck checks signup duplication (V1 Mall API — generic form)
// GET /v1/users/join/{type}?pk={pk}
// checkType: email-exist, id-exist, phone-exist, uid-exist, group-business-number-exist
// (see MALL_USER_JOIN_CHECK_* constants; new server keys also work without an SDK update)
func (m *UserModule) UserJoinCheck(checkType string, pk string, idempotencyKey string) (map[string]interface{}, error) {
	return m.api.getWithHeaders(fmt.Sprintf("users/join/%s?pk=%s", checkType, url.QueryEscape(pk)), commerceMallHeaders("", idempotencyKey))
}

// UidExist checks external uid (ex_uid) duplication
// GET /v1/users/join/uid-exist?pk={uid}
// Dedicated form alongside email-exist / id-exist / phone-exist / group-business-number-exist.
func (m *UserModule) UidExist(uid string, idempotencyKey string) (map[string]interface{}, error) {
	headers := commerceMallHeaders("", idempotencyKey)
	headers["BOOTPAY-ROLE"] = "user"
	return m.api.getWithHeaders("users/join/uid-exist?pk="+url.QueryEscape(uid), headers)
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
