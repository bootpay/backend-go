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
func (m *UserGroupModule) UserCreate(userGroupId string, userId string) (map[string]interface{}, error) {
	data := map[string]string{
		"user_id": userId,
	}
	return m.api.Post(fmt.Sprintf("user-groups/%s/add_user", userGroupId), data)
}

// UserDelete removes a user from a group
func (m *UserGroupModule) UserDelete(userGroupId string, userId string) (map[string]interface{}, error) {
	return m.api.Delete(fmt.Sprintf("user-groups/%s/remove_user?user_id=%s", userGroupId, userId))
}

// Limit sets group limit settings
func (m *UserGroupModule) Limit(params UserGroupLimitParams) (map[string]interface{}, error) {
	if params.UserGroupId == "" {
		return nil, fmt.Errorf("user_group_id is required")
	}
	return m.api.Put(fmt.Sprintf("user-groups/%s/limit", params.UserGroupId), params)
}

// AggregateTransaction sets group aggregate transaction settings
func (m *UserGroupModule) AggregateTransaction(params UserGroupAggregateTransactionParams) (map[string]interface{}, error) {
	if params.UserGroupId == "" {
		return nil, fmt.Errorf("user_group_id is required")
	}
	return m.api.Put(fmt.Sprintf("user-groups/%s/aggregate-transaction", params.UserGroupId), params)
}
