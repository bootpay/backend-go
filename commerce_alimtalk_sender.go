package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
)

// AlimtalkSenderModule handles the alimtalk sender profile (kakao channel) lifecycle
// GET /v1/alimtalk/categories · /senders 계열
//
// 카테고리 조회 → OTP 발송 → 발신프로필 등록 → 목록/상세 → 연동 해지 순으로 쓴다.
// 등록이 끝나면 서버가 그룹키 등록까지 자동으로 하므로, 공식 템플릿은 별도 채택 없이 바로 발송된다.
//
// ⚠️ 실제 부작용: Otp 는 채널 관리자 휴대폰으로 **문자를 실제 발송**하고,
//
//	Create 는 카카오에 발신프로필을 **실제 등록**한다. 샌드박스가 없다.
type AlimtalkSenderModule struct {
	api *CommerceApi
}

// Categories retrieves the kakao category list
// GET /v1/alimtalk/categories
// 발신프로필 등록 시 필요한 category_code 후보다. 벤더 응답을 그대로 프록시한다.
func (m *AlimtalkSenderModule) Categories() (map[string]interface{}, error) {
	return m.api.getWithHeaders("alimtalk/categories", alimtalkHeaders())
}

// Otp sends an OTP to the channel admin phone
// POST /v1/alimtalk/senders/otp
// ⚠️ 실제로 문자가 나간다. 여기서 받은 인증번호를 Create 의 Otp 로 넘긴다.
func (m *AlimtalkSenderModule) Otp(params AlimtalkSenderOtpParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("alimtalk/senders/otp", params, alimtalkHeaders())
}

// Create registers a sender profile
// POST /v1/alimtalk/senders
// ⚠️ 카카오에 발신프로필이 실제 등록된다. 같은 yellow_id 를 다시 등록하면 기존 프로필을 재사용한다(dedup).
// 등록 성공 시 그룹키 등록까지 서버가 수행하므로 공식 카탈로그 전체를 바로 발송할 수 있다.
func (m *AlimtalkSenderModule) Create(params AlimtalkSenderCreateParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("alimtalk/senders", params, alimtalkHeaders())
}

// List retrieves the linked channel list
// GET /v1/alimtalk/senders
// 자체 DB 만 조회하며 벤더를 호출하지 않는다. 응답은 { list: [...], count: N }.
func (m *AlimtalkSenderModule) List() (map[string]interface{}, error) {
	return m.api.getWithHeaders("alimtalk/senders", alimtalkHeaders())
}

// Detail retrieves channel details
// GET /v1/alimtalk/senders/{ksp_id}
// sync 가 true 면 벤더에서 채널 상태를 다시 읽어 반영한다(느리다). nil 이면 자체 DB 만 본다.
// ⚠️ 미연동/미존재 채널은 404, 다른 프로젝트의 채널은 403 으로 오며 둘 다 error_code 는 3024 다.
func (m *AlimtalkSenderModule) Detail(kspId string, sync *bool) (map[string]interface{}, error) {
	query := url.Values{}
	if sync != nil {
		query.Set("sync", strconv.FormatBool(*sync))
	}
	return m.api.getWithHeaders(withQuery(fmt.Sprintf("alimtalk/senders/%s", kspId), query), alimtalkHeaders())
}

// Release unlinks the channel from this project
// DELETE /v1/alimtalk/senders/{ksp_id}
// 이 프로젝트와의 연동만 끊는다 — 채널 모델과 템플릿은 보존된다. 성공 시 본문은 null 이다.
func (m *AlimtalkSenderModule) Release(kspId string) (map[string]interface{}, error) {
	return m.api.deleteWithHeaders(fmt.Sprintf("alimtalk/senders/%s", kspId), nil, alimtalkHeaders())
}

// VariableExamples updates the channel variable example dictionary
// PUT /v1/alimtalk/senders/{ksp_id}/variable_examples
// 템플릿 미리보기에서 #{user_name} 대신 '홍길동' 처럼 읽히게 하는 **표시용** 값이다.
// ⚠️ 발송값이 아니다 — 벤더로 전송되지 않으므로 검수 상태와 무관하다. 보낸 키만 덮어쓴다(부분 갱신).
// examples: { "user_name": "홍길동", "company_name": "부트페이몰" } — 키에 '.' 이나 선행 '$' 는 쓸 수 없다.
func (m *AlimtalkSenderModule) VariableExamples(kspId string, examples map[string]interface{}) (map[string]interface{}, error) {
	body := map[string]interface{}{}
	if examples != nil {
		body["examples"] = examples
	}
	return m.api.putWithHeaders(fmt.Sprintf("alimtalk/senders/%s/variable_examples", kspId), body, alimtalkHeaders())
}
