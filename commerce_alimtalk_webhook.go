package bootpay

import (
	"net/url"
	"strconv"
)

// AlimtalkWebhookModule handles alimtalk send-result / inspection-result webhook settings
// /v1/alimtalk/webhook 계열
//
// ⚠️ **주문·구독 통합 웹훅과 완전히 별개다.** 알림톡 이벤트를 기존 주문 웹훅 URL 로 태우면
//
//	그 수신 서버가 모르는 payload 를 받아 기존 연동이 깨진다. 그래서 수신 URL 을 따로 둔다.
//	(Webhook.SendTest 는 주문 웹훅용이다 — 이 모듈의 Test 와 혼동하지 말 것)
//
// # 서명 검증
//
// 요청에 다음 헤더가 붙는다.
//
//	X-Bootpay-Signature: sha256=HMAC_SHA256(secret, "{X-Bootpay-Timestamp}.{raw_body}")
//
// 타임스탬프가 5분 이상 지난 요청은 거부한다(replay 방지).
type AlimtalkWebhookModule struct {
	api *CommerceApi
}

// Detail retrieves the webhook setting
// GET /v1/alimtalk/webhook
// 시크릿은 앞 12자만 노출된다. 미설정이면 { configured: false } 로 온다.
func (m *AlimtalkWebhookModule) Detail() (map[string]interface{}, error) {
	return m.api.getWithHeaders("alimtalk/webhook", alimtalkHeaders())
}

// Update saves the webhook setting
// PUT /v1/alimtalk/webhook
// Url 은 **https 만** 허용한다(아니면 3028). 최초 저장 시 서명 시크릿이 자동 발급된다.
func (m *AlimtalkWebhookModule) Update(params AlimtalkWebhookUpdateParams) (map[string]interface{}, error) {
	return m.api.putWithHeaders("alimtalk/webhook", params, alimtalkHeaders())
}

// Test sends one test event
// POST /v1/alimtalk/webhook/test
// ⚠️ **설정된 URL 로 실제 HTTP 요청이 나간다.** 구독 여부와 무관하게 보낸다.
// 웹훅이 설정돼 있지 않으면 3029. 응답: { delivery_id:, url:, queued: }
func (m *AlimtalkWebhookModule) Test() (map[string]interface{}, error) {
	return m.api.postWithHeaders("alimtalk/webhook/test", map[string]interface{}{}, alimtalkHeaders())
}

// RotateSecret reissues the signing secret
// POST /v1/alimtalk/webhook/secret
// ⚠️ **이 응답에서만 secret 원문을 돌려준다**(이후 조회는 마스킹된다).
// ⚠️ 이미 큐에 있는 전송 건은 발송 당시 시크릿으로 서명된다.
func (m *AlimtalkWebhookModule) RotateSecret() (map[string]interface{}, error) {
	return m.api.postWithHeaders("alimtalk/webhook/secret", map[string]interface{}{}, alimtalkHeaders())
}

// Deliveries retrieves the webhook delivery history
// GET /v1/alimtalk/webhook/deliveries
// 성공·실패를 모두 남긴다. 응답: { list: [{ delivery_id:, event:, event_code:, url:, status:,
//
//	retry_count:, max_retry:, tags:, created_at: }], count:, page:, per: }
func (m *AlimtalkWebhookModule) Deliveries(params *AlimtalkWebhookDeliveriesParams) (map[string]interface{}, error) {
	query := url.Values{}
	if params != nil {
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
	}
	return m.api.getWithHeaders(withQuery("alimtalk/webhook/deliveries", query), alimtalkHeaders())
}
