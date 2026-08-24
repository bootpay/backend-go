package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
)

// UserGroupModule handles user group-related operations
type UserGroupModule struct {
	api *CommerceApi
}

// Create creates a new user group
func (m *UserGroupModule) Create(userGroup CommerceUserGroup) (map[string]interface{}, error) {
	return m.api.Post("user-groups", userGroup)
}

// List retrieves user group list
func (m *UserGroupModule) List(params *UserGroupListParams) (map[string]interface{}, error) {
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
		if params.CorporateType > 0 {
			queryParams.Set("corporate_type", strconv.Itoa(params.CorporateType))
		}
		if len(queryParams) > 0 {
			query = "?" + queryParams.Encode()
		}
	}
	return m.api.Get("user-groups" + query)
}

// Detail retrieves user group details
func (m *UserGroupModule) Detail(userGroupId string) (map[string]interface{}, error) {
	return m.api.Get(fmt.Sprintf("user-groups/%s", userGroupId))
}

// Update updates user group information
func (m *UserGroupModule) Update(userGroup CommerceUserGroup) (map[string]interface{}, error) {
	if userGroup.UserGroupId == "" {
		return nil, fmt.Errorf("user_group_id is required")
	}
	return m.api.Put(fmt.Sprintf("user-groups/%s", userGroup.UserGroupId), userGroup)
}

// UserCreate adds a user to a group
// ⚠️ The server requires manager scope (scope_invalid!) — BOOTPAY-ROLE must be sent.
// idempotencyKey is optional; when omitted a key is generated per call.
func (m *UserGroupModule) UserCreate(userGroupId string, userId string, idempotencyKey ...string) (map[string]interface{}, error) {
	data := map[string]string{
		"user_id": userId,
	}
	return m.api.postWithHeaders(fmt.Sprintf("user-groups/%s/user", userGroupId), data,
		commerceRoleHeaders("manager", firstOrEmpty(idempotencyKey)))
}

// UserDelete removes a user from a group
// ⚠️ The server requires manager scope (scope_invalid!) — BOOTPAY-ROLE must be sent.
// idempotencyKey is optional; when omitted a key is generated per call.
func (m *UserGroupModule) UserDelete(userGroupId string, userId string, idempotencyKey ...string) (map[string]interface{}, error) {
	return m.api.deleteWithHeaders(fmt.Sprintf("user-groups/%s/user/%s", userGroupId, userId), nil,
		commerceRoleHeaders("manager", firstOrEmpty(idempotencyKey)))
}

// Limit sets group purchase limit settings (manager scope)
// PUT /v1/user-groups/{user_group_id}/limit
// ⚠️ Limits are never applied through Update — the server's user_groups_controller#update
// explicitly strips use_limit / limit_message / limit_month_purchase / limit_week_purchase.
// This dedicated route is the only way to change them. Server scope: manager:limit
func (m *UserGroupModule) Limit(params UserGroupLimitParams) (map[string]interface{}, error) {
	if params.UserGroupId == "" {
		return nil, fmt.Errorf("user_group_id is required")
	}
	return m.api.putWithHeaders(fmt.Sprintf("user-groups/%s/limit", params.UserGroupId), params, commerceRoleHeaders("manager", params.IdempotencyKey))
}

// AggregateTransaction sets group aggregate transaction settings (manager scope)
// PUT /v1/user-groups/{user_group_id}/aggregate-transaction
// Update has same-named arguments but the server only processes this dedicated route.
func (m *UserGroupModule) AggregateTransaction(params UserGroupAggregateTransactionParams) (map[string]interface{}, error) {
	if params.UserGroupId == "" {
		return nil, fmt.Errorf("user_group_id is required")
	}
	return m.api.putWithHeaders(fmt.Sprintf("user-groups/%s/aggregate-transaction", params.UserGroupId), params, commerceRoleHeaders("manager", params.IdempotencyKey))
}
