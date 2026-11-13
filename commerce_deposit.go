package bootpay

import (
	"fmt"
	"net/url"
)

// DepositModule handles 뱅크알람 입금대기 operations.
// 금액·입금자명이 맞는 입금 알림을 주문에 연결하고 입금 웹훅으로 알린다.
// 전부 supervisor 전용이며 인스턴스 role 은 바꾸지 않는다.
type DepositModule struct {
	api *CommerceApi
}

// Accounts retrieves the seller's deposit accounts.
// GET /v1/sdk/deposits/accounts
// 첫 계좌가 대표 계좌(primary: true)다.
func (m *DepositModule) Accounts(idempotencyKey ...string) (map[string]interface{}, error) {
	return m.api.getWithHeaders("sdk/deposits/accounts", commerceRoleHeaders("supervisor", firstOrEmpty(idempotencyKey)))
}

// WebhookSecret reissues the webhook signing secret.
// POST /v1/sdk/deposits/webhook_secret
// 원문은 이 응답에서만 준다. 다시 부르면 바뀌며 이미 만든 웹훅은 이전 비밀값으로 서명된다.
func (m *DepositModule) WebhookSecret(idempotencyKey ...string) (map[string]interface{}, error) {
	return m.api.postWithHeaders("sdk/deposits/webhook_secret", map[string]interface{}{},
		commerceRoleHeaders("supervisor", firstOrEmpty(idempotencyKey)))
}

// Create registers a pending deposit.
// POST /v1/sdk/deposits
// 같은 merchant_order_id 로 금액·입금자명·계좌가 같은 재호출은 기존 주문을 돌려준다(하나라도 다르면 18110).
// WebhookUrl 은 비밀값을 먼저 발급받아야 한다(18131). Notify 는 유료다(준비가 안 됐으면 18133).
func (m *DepositModule) Create(params DepositCreateParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("sdk/deposits", params, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// Detail retrieves a pending deposit by merchant order ID or Bootpay order number.
// GET /v1/sdk/deposits/{id}
// 기한이 지났으면 조회 시점에 expired 로 닫힌다.
func (m *DepositModule) Detail(depositId string, idempotencyKey ...string) (map[string]interface{}, error) {
	if depositId == "" {
		return nil, fmt.Errorf("deposit_id is required")
	}
	return m.api.getWithHeaders("sdk/deposits/"+url.PathEscape(depositId),
		commerceRoleHeaders("supervisor", firstOrEmpty(idempotencyKey)))
}

// Cancel cancels a pending deposit.
// DELETE /v1/sdk/deposits/{id}
// 이미 닫혔으면 그대로 돌려주며 기한이 지났으면 만료로 닫는다. 입금 확인된 주문은 18108.
func (m *DepositModule) Cancel(depositId string, idempotencyKey ...string) (map[string]interface{}, error) {
	if depositId == "" {
		return nil, fmt.Errorf("deposit_id is required")
	}
	return m.api.deleteWithHeaders("sdk/deposits/"+url.PathEscape(depositId), nil,
		commerceRoleHeaders("supervisor", firstOrEmpty(idempotencyKey)))
}
