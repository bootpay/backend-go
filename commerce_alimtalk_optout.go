package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
)

// AlimtalkOptoutModule handles alimtalk opt-out (수신거부) operations
// /v1/alimtalk/optouts 계열 (가맹점 CRM 수신거부 동기화용)
//
// 발송 판정과 **같은 기준**으로 다룬다 — 부트페이 전역(global) + 내 프로젝트.
// ⚠️ 전역 건은 **조회는 되지만 해제할 수 없다**(releasable: false).
//
//	이걸 보지 않으면 "화면엔 수신거부가 아닌데 발송은 3021 로 막히는" 상태가 된다.
type AlimtalkOptoutModule struct {
	api *CommerceApi
}

// List retrieves the opt-out list
// GET /v1/alimtalk/optouts
// phone 은 숫자만 남겨 **부분일치**로 찾는다(정확 매칭이 아니다). 50건 단위로 페이징된다.
// 응답: { list: [{ id:, phone:, scope:, global:, releasable:, source:, reason:, opted_out_at:,
//
//	created_at: }], count:, page: }
func (m *AlimtalkOptoutModule) List(params *AlimtalkOptoutListParams) (map[string]interface{}, error) {
	query := url.Values{}
	if params != nil {
		if params.Phone != "" {
			query.Set("phone", params.Phone)
		}
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
	}
	return m.api.getWithHeaders(withQuery("alimtalk/optouts", query), alimtalkHeaders())
}

// Create registers an opt-out
// POST /v1/alimtalk/optouts
// 내 프로젝트 스코프로 등록된다(source: api). 같은 번호를 다시 등록해도 멱등이다.
func (m *AlimtalkOptoutModule) Create(params AlimtalkOptoutCreateParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("alimtalk/optouts", params, alimtalkHeaders())
}

// Check pre-checks opt-outs before sending
// POST /v1/alimtalk/optouts/check
// 발송 판정과 **같은 축**으로 대조하므로, 벌크에서 skipped 로 낭비될 건을 미리 뺄 수 있다.
// 응답: { list: [{ phone:, opted_out:, global:, releasable:, opted_out_at: }], count:, opted_out_count: }
func (m *AlimtalkOptoutModule) Check(params AlimtalkOptoutCheckParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("alimtalk/optouts/check", params, alimtalkHeaders())
}

// Release releases an opt-out
// DELETE /v1/alimtalk/optouts/{phone}
// 내 프로젝트 스코프 건만 해제되며 멱등이다(없어도 성공).
// ⚠️ 전역 차단은 해제되지 않고 global_blocked: true 로 알려 준다 —
//
//	"지웠는데 여전히 막히는" 상태를 응답으로 드러내기 위함이다.
//
// 응답: { phone:, released:, global_blocked: }
func (m *AlimtalkOptoutModule) Release(phone string) (map[string]interface{}, error) {
	return m.api.deleteWithHeaders(fmt.Sprintf("alimtalk/optouts/%s", phone), nil, alimtalkHeaders())
}
