package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
)

// AlimtalkMessageModule handles alimtalk send-history and aggregate operations
// GET /v1/alimtalk/messages 계열
//
// **유료** 알림톡만 조회된다(무료 커머스 알림톡은 포함되지 않는다).
// 상태는 벤더 결과 동기화로 확정되므로 접수 직후에는 requested 로 보인다.
type AlimtalkMessageModule struct {
	api *CommerceApi
}

// List retrieves the alimtalk send history
// GET /v1/alimtalk/messages
//
// ⚠️ 기간 기본값은 최근 30일이고 최대 조회 폭은 92일이다 — 초과분은 거부하지 않고 시작일을
// 당겨 잘라낸다. 실제 적용된 구간은 응답의 period 로 확인한다.
// 응답: { list: [...], count:, page:, per:, period: { from:, to: } }
func (m *AlimtalkMessageModule) List(params *AlimtalkMessageListParams) (map[string]interface{}, error) {
	query := url.Values{}
	if params != nil {
		if params.TemplateCode != "" {
			query.Set("template_code", params.TemplateCode)
		}
		if params.Status != "" {
			query.Set("status", params.Status)
		}
		if params.RefId != "" {
			query.Set("ref_id", params.RefId)
		}
		if params.To != "" {
			query.Set("to", params.To)
		}
		if params.SAt != "" {
			query.Set("s_at", params.SAt)
		}
		if params.EAt != "" {
			query.Set("e_at", params.EAt)
		}
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
	}
	return m.api.getWithHeaders(withQuery("alimtalk/messages", query), alimtalkHeaders())
}

// Stats retrieves the period aggregate
// GET /v1/alimtalk/messages/stats
// 일자별 집계 원장에서 읽으므로 응답이 빠르다.
// 응답: { period:, totals: { sent, success, failed, fallback, opted_out_hit, rejected, canceled,
//
//	success_rate }, daily: [...], billing: { billable_count, unit_price, fallback_count, ..., amount } }
//
// ⚠️ billing.unit_price_source 가 'default' 면 **잠정 단가**다(확정 청구액이 아니다).
// ⚠️ billable_count 는 성공 − 폴백이다 — 폴백분은 LMS 단가로 따로 계산된다.
func (m *AlimtalkMessageModule) Stats(sAt string, eAt string) (map[string]interface{}, error) {
	query := url.Values{}
	if sAt != "" {
		query.Set("s_at", sAt)
	}
	if eAt != "" {
		query.Set("e_at", eAt)
	}
	return m.api.getWithHeaders(withQuery("alimtalk/messages/stats", query), alimtalkHeaders())
}

// Detail retrieves a single send result
// GET /v1/alimtalk/messages/{receipt_id}
// 실패 사유는 error_code·error_message 에 담긴다.
// fallback_type 은 폴백이 꺼진 건이면 null, 켜진 건이면 LMS 다.
// 다른 프로젝트의 건이거나 없으면 404(3025).
func (m *AlimtalkMessageModule) Detail(receiptId string) (map[string]interface{}, error) {
	return m.api.getWithHeaders(fmt.Sprintf("alimtalk/messages/%s", receiptId), alimtalkHeaders())
}
