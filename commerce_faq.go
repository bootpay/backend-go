package bootpay

import (
	"fmt"
	"net/url"
)

// FaqModule handles FAQ (자주 묻는 질문) 게시판 operations
//
// 고객 모드(BOOTPAY-ROLE: user)와 운영자 모드(supervisor)가 같은 경로를 공유한다 —
// 등록·수정·삭제는 supervisor 전용이고, 목록·단건만 고객 모드로도 열린다.
// ⚠️ 고객 모드는 몰 FAQ 사용여부가 꺼져 있으면 BOARD_FEATURE_DISABLED 로 거절된다.
type FaqModule struct {
	api *CommerceApi
}

// List retrieves the FAQ list
// GET /v1/faqs
// page/limit 은 미지정 시 1/20 으로 항상 실린다.
// Supervisor 가 true 면 운영자 모드 — View: "all" 로 비공개 FAQ 까지 볼 수 있다.
func (m *FaqModule) List(params *FaqListParams) (map[string]interface{}, error) {
	if params == nil {
		params = &FaqListParams{}
	}
	query := url.Values{}
	boardPageParams(query, params.Page, params.Limit)
	if params.Keyword != "" {
		query.Set("keyword", params.Keyword)
	}
	if params.View != "" {
		query.Set("view", params.View)
	}
	return m.api.getWithHeaders(withQuery("faqs", query), commerceBoardHeaders(params.Supervisor, "", params.IdempotencyKey))
}

// Detail retrieves a single FAQ
// GET /v1/faqs/{faq_id}
// ⚠️ 비공개 FAQ 는 supervisor 가 true 일 때만 보인다(아니면 404 POST_NOT_FOUND).
func (m *FaqModule) Detail(faqId string, supervisor bool, idempotencyKey ...string) (map[string]interface{}, error) {
	if faqId == "" {
		return nil, fmt.Errorf("faq_id is required")
	}
	return m.api.getWithHeaders(fmt.Sprintf("faqs/%s", faqId),
		commerceBoardHeaders(supervisor, "", firstOrEmpty(idempotencyKey)))
}

// Create registers a FAQ (supervisor scope)
// POST /v1/faqs
func (m *FaqModule) Create(params FaqCreateParams) (map[string]interface{}, error) {
	body, err := boardPayload(params, params.Images)
	if err != nil {
		return nil, err
	}
	return m.api.postWithHeaders("faqs", body, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// Update updates a FAQ (supervisor scope)
// PUT /v1/faqs/{faq_id}
// 보낸 필드만 바뀐다. Images 는 보내면 목록 전체를 교체한다(빈 슬라이스면 모두 삭제).
func (m *FaqModule) Update(params FaqUpdateParams) (map[string]interface{}, error) {
	if params.FaqId == "" {
		return nil, fmt.Errorf("faq_id is required")
	}
	body, err := boardPayload(params, params.Images)
	if err != nil {
		return nil, err
	}
	return m.api.putWithHeaders(fmt.Sprintf("faqs/%s", params.FaqId), body, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// Delete deletes a FAQ (supervisor scope)
// DELETE /v1/faqs/{faq_id}
func (m *FaqModule) Delete(faqId string, idempotencyKey ...string) (map[string]interface{}, error) {
	if faqId == "" {
		return nil, fmt.Errorf("faq_id is required")
	}
	return m.api.deleteWithHeaders(fmt.Sprintf("faqs/%s", faqId), nil,
		commerceRoleHeaders("supervisor", firstOrEmpty(idempotencyKey)))
}
