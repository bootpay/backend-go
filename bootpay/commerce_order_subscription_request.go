package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
)

// OrderSubscriptionRequestModule handles V1 OrderSubscription Request lookup/approval operations
//
// User mode (role=user): call without project_id -> own request list/detail
// Supervisor mode (role=supervisor): include project_id -> project-wide list + update (approve/reject)
//
// For buyer-side request creation (pause/resume/termination), use
// commerce.OrderSubscription.RequestIng.* module instead.
type OrderSubscriptionRequestModule struct {
	api *CommerceApi
}

// List retrieves request list (user / supervisor common)
func (m *OrderSubscriptionRequestModule) List(params *OrderSubscriptionRequestListParams) (map[string]interface{}, error) {
	query := ""
	if params != nil {
		queryParams := url.Values{}
		if params.ProjectId != "" {
			queryParams.Set("project_id", params.ProjectId)
		}
		if params.Page > 0 {
			queryParams.Set("page", strconv.Itoa(params.Page))
		}
		if params.Limit > 0 {
			queryParams.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.RequestType != 0 {
			queryParams.Set("request_type", strconv.Itoa(params.RequestType))
		}
		if params.Status != 0 {
			queryParams.Set("status", strconv.Itoa(params.Status))
		}
		if params.SAt != "" {
			queryParams.Set("s_at", params.SAt)
		}
		if params.EAt != "" {
			queryParams.Set("e_at", params.EAt)
		}
		if params.Keyword != "" {
			queryParams.Set("keyword", params.Keyword)
		}
		if len(queryParams) > 0 {
			query = "?" + queryParams.Encode()
		}
	}
	return m.api.Get("order-subscription-requests" + query)
}

// Detail retrieves a single request (user / supervisor common)
func (m *OrderSubscriptionRequestModule) Detail(orderSubscriptionRequestHistoryId string, projectId string) (map[string]interface{}, error) {
	query := ""
	if projectId != "" {
		queryParams := url.Values{}
		queryParams.Set("project_id", projectId)
		query = "?" + queryParams.Encode()
	}
	return m.api.Get(fmt.Sprintf("order-subscription-requests/%s%s", orderSubscriptionRequestHistoryId, query))
}

// Update approves/rejects a request (supervisor only)
func (m *OrderSubscriptionRequestModule) Update(params OrderSubscriptionRequestUpdateParams) (map[string]interface{}, error) {
	if params.OrderSubscriptionRequestHistoryId == "" {
		return nil, fmt.Errorf("order_subscription_request_history_id is required")
	}
	historyId := params.OrderSubscriptionRequestHistoryId
	body := OrderSubscriptionRequestUpdateBody{
		Approval: params.Approval,
		Reason:   params.Reason,
		Extra:    params.Extra,
	}
	return m.api.Put(fmt.Sprintf("order-subscription-requests/%s", historyId), body)
}
