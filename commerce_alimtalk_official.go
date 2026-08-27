package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
)

// AlimtalkOfficialModule handles the Bootpay official alimtalk template catalog
// GET/POST /v1/alimtalk/official 계열
//
// 부트페이가 미리 카카오 승인을 받아 둔 템플릿이라, 그룹키가 등록된 채널이면 **검수 없이 즉시 발송**된다.
// AlimtalkSender.Create 로 채널을 등록하면 그룹 등록이 함께 끝나므로 따로 채택할 것이 없다.
// (26-08-27 채택 엔드포인트는 서버에서 비활성화되어 SDK 에도 두지 않는다)
//
// 전부 조회 계열이라 부작용이 없다(자체 DB 만 본다).
type AlimtalkOfficialModule struct {
	api *CommerceApi
}

// List searches the official template catalog
// GET /v1/alimtalk/official
// keyword 는 본문·이름·분류를 부분일치(대소문자 무시)로 훑는다.
// 응답: { list: [...], count:, page:, per:, categories: [...] }
func (m *AlimtalkOfficialModule) List(params *AlimtalkOfficialListParams) (map[string]interface{}, error) {
	query := url.Values{}
	if params != nil {
		// 서버는 q 를 먼저 보고 없으면 keyword 를 본다 — 정본 키인 q 로 보낸다
		if params.Keyword != "" {
			query.Set("q", params.Keyword)
		}
		if params.Category != "" {
			query.Set("category", params.Category)
		}
		if params.MsgType != "" {
			query.Set("msg_type", params.MsgType)
		}
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
		if params.Per > 0 {
			query.Set("per", strconv.Itoa(params.Per))
		}
		if params.KspId != "" {
			query.Set("ksp_id", params.KspId)
		}
	}
	return m.api.getWithHeaders(withQuery("alimtalk/official", query), alimtalkHeaders())
}

// Recommend recommends official templates for the text you want to send
// POST /v1/alimtalk/official/recommend
// 유사도 score(0~1) 내림차순으로 돌려준다.
func (m *AlimtalkOfficialModule) Recommend(params AlimtalkOfficialRecommendParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("alimtalk/official/recommend", params, alimtalkHeaders())
}

// Detail retrieves an official template
// GET /v1/alimtalk/official/{code}
// code 는 서버 채번 코드(슬래시를 포함하지 않는다). 없거나 미노출이면 404(3015).
// kspId 를 주면 그 채널의 변수 예문 사전으로 variable_examples 를 채워 준다(표시용).
func (m *AlimtalkOfficialModule) Detail(code string, kspId string) (map[string]interface{}, error) {
	query := url.Values{}
	if kspId != "" {
		query.Set("ksp_id", kspId)
	}
	return m.api.getWithHeaders(withQuery(fmt.Sprintf("alimtalk/official/%s", code), query), alimtalkHeaders())
}
