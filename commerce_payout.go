package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
)

// PayoutModule handles 입점사 지급 요청 operations.
// 실제 송금은 별도 배치가 처리한다. 전부 supervisor 전용이며 인스턴스 role 은 바꾸지 않는다.
type PayoutModule struct {
	api *CommerceApi
}

// Create registers a payout request.
// POST /v1/sdk/payouts
// 같은 merchant_payout_id 재호출은 처음 만든 항목을 돌려준다.
func (m *PayoutModule) Create(params PayoutCreateParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("sdk/payouts", params, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// Detail retrieves a payout by merchant payout ID or item ID.
// GET /v1/sdk/payouts/{id}
func (m *PayoutModule) Detail(payoutId string, idempotencyKey ...string) (map[string]interface{}, error) {
	if payoutId == "" {
		return nil, fmt.Errorf("payout_id is required")
	}
	return m.api.getWithHeaders("sdk/payouts/"+url.PathEscape(payoutId),
		commerceRoleHeaders("supervisor", firstOrEmpty(idempotencyKey)))
}

// List retrieves payout items in recent creation order.
// GET /v1/sdk/payouts
// 미지정 필터와 page/limit 은 보내지 않으며 서버 응답의 bare array 는 data 배열로 감싸 반환한다.
func (m *PayoutModule) List(params *PayoutListParams) (map[string]interface{}, error) {
	if params == nil {
		params = &PayoutListParams{}
	}
	query := url.Values{}
	if params.Status != "" {
		query.Set("status", string(params.Status))
	}
	if params.SubSellerId != "" {
		query.Set("sub_seller_id", params.SubSellerId)
	}
	if params.From != "" {
		query.Set("from", params.From)
	}
	if params.To != "" {
		query.Set("to", params.To)
	}
	if params.Page != 0 {
		query.Set("page", strconv.Itoa(params.Page))
	}
	if params.Limit != 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	path := "sdk/payouts"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	return m.api.getWithHeaders(path, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}
