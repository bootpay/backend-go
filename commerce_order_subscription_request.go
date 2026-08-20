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
// GET /v1/order-subscription-requests
// With ProjectId the request runs in supervisor mode (project-wide search), otherwise
// only the caller's own requests are returned. page/limit default to 1/20 when unset.
func (m *OrderSubscriptionRequestModule) List(params *OrderSubscriptionRequestListParams) (map[string]interface{}, error) {
	if params == nil {
		params = &OrderSubscriptionRequestListParams{}
	}
	queryParams := url.Values{}
	if params.ProjectId != "" {
		queryParams.Set("project_id", params.ProjectId)
	}
	if params.OrderSubscriptionId != "" {
		queryParams.Set("order_subscription_id", params.OrderSubscriptionId)
	}
	page := params.Page
	if page <= 0 {
		page = 1
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}
	queryParams.Set("page", strconv.Itoa(page))
	queryParams.Set("limit", strconv.Itoa(limit))
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
	if params.UserId != "" {
		queryParams.Set("user_id", params.UserId)
	}
	if params.UserGroupId != "" {
		queryParams.Set("user_group_id", params.UserGroupId)
	}
	return m.api.getWithHeaders("order-subscription-requests?"+queryParams.Encode(), requestScopeHeaders(params.ProjectId, params.IdempotencyKey))
}

// Detail retrieves a single request (user / supervisor common)
// GET /v1/order-subscription-requests/{id}
func (m *OrderSubscriptionRequestModule) Detail(orderSubscriptionRequestHistoryId string, projectId string) (map[string]interface{}, error) {
	query := ""
	if projectId != "" {
		queryParams := url.Values{}
		queryParams.Set("project_id", projectId)
		query = "?" + queryParams.Encode()
	}
	return m.api.getWithHeaders(fmt.Sprintf("order-subscription-requests/%s%s", orderSubscriptionRequestHistoryId, query), requestScopeHeaders(projectId, ""))
}

// Update approves/rejects a request (supervisor only)
// PUT /v1/order-subscription-requests/{id}
// ⚠️ Approve and reject are not separate actions — the route set is index/show/update only,
// branched by approval: "approve" | "reject" (the key is named approval because the server
// uses params[:action] as a Rails reserved word).
func (m *OrderSubscriptionRequestModule) Update(params OrderSubscriptionRequestUpdateParams) (map[string]interface{}, error) {
	if params.OrderSubscriptionRequestHistoryId == "" {
		return nil, fmt.Errorf("order_subscription_request_history_id is required")
	}
	historyId := params.OrderSubscriptionRequestHistoryId
	body := OrderSubscriptionRequestUpdateBody{
		Approval:            params.Approval,
		Reason:              params.Reason,
		Price:               params.Price,
		TaxFreePrice:        params.TaxFreePrice,
		TerminationFee:      params.TerminationFee,
		LastBillRefundPrice: params.LastBillRefundPrice,
		FinalFee:            params.FinalFee,
		ServiceEndAt:        params.ServiceEndAt,
		Extra:               params.Extra,
	}
	return m.api.putWithHeaders(fmt.Sprintf("order-subscription-requests/%s", historyId), body, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// requestScopeHeaders returns lookup headers — supervisor when project_id is present, user otherwise.
func requestScopeHeaders(projectId string, idempotencyKey string) map[string]string {
	role := "user"
	if projectId != "" {
		role = "supervisor"
	}
	return commerceRoleHeaders(role, idempotencyKey)
}
