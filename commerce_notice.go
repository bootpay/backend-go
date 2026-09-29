package bootpay

import (
	"fmt"
	"net/url"
)

// NoticeModule handles 공지사항 게시판 operations
//
// FAQ 와 같은 구조다 — 등록·수정·삭제는 supervisor 전용이고, 목록·단건만 고객 모드로도 열린다.
// ⚠️ 고객 모드는 몰 공지사항 사용여부가 꺼져 있으면 BOARD_FEATURE_DISABLED 로 거절된다.
type NoticeModule struct {
	api *CommerceApi
}

// List retrieves the notice list
// GET /v1/notices
// page/limit 은 미지정 시 1/20 으로 항상 실린다.
// Supervisor 가 true 면 운영자 모드 — View: "all" 로 비공개 공지까지 볼 수 있다.
func (m *NoticeModule) List(params *NoticeListParams) (map[string]interface{}, error) {
	if params == nil {
		params = &NoticeListParams{}
	}
	query := url.Values{}
	boardPageParams(query, params.Page, params.Limit)
	if params.Keyword != "" {
		query.Set("keyword", params.Keyword)
	}
	if params.View != "" {
		query.Set("view", params.View)
	}
	return m.api.getWithHeaders(withQuery("notices", query), commerceBoardHeaders(params.Supervisor, "", params.IdempotencyKey))
}

// Detail retrieves a single notice
// GET /v1/notices/{notice_id}
// ⚠️ 비공개 공지는 supervisor 가 true 일 때만 보인다(아니면 404 POST_NOT_FOUND).
func (m *NoticeModule) Detail(noticeId string, supervisor bool, idempotencyKey ...string) (map[string]interface{}, error) {
	if noticeId == "" {
		return nil, fmt.Errorf("notice_id is required")
	}
	return m.api.getWithHeaders(fmt.Sprintf("notices/%s", noticeId),
		commerceBoardHeaders(supervisor, "", firstOrEmpty(idempotencyKey)))
}

// Create registers a notice (supervisor scope)
// POST /v1/notices
func (m *NoticeModule) Create(params NoticeCreateParams) (map[string]interface{}, error) {
	body, err := boardPayload(params, params.Images)
	if err != nil {
		return nil, err
	}
	return m.api.postWithHeaders("notices", body, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// Update updates a notice (supervisor scope)
// PUT /v1/notices/{notice_id}
// 보낸 필드만 바뀐다. Images 는 보내면 목록 전체를 교체한다(빈 슬라이스면 모두 삭제).
func (m *NoticeModule) Update(params NoticeUpdateParams) (map[string]interface{}, error) {
	if params.NoticeId == "" {
		return nil, fmt.Errorf("notice_id is required")
	}
	body, err := boardPayload(params, params.Images)
	if err != nil {
		return nil, err
	}
	return m.api.putWithHeaders(fmt.Sprintf("notices/%s", params.NoticeId), body, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// Delete deletes a notice (supervisor scope)
// DELETE /v1/notices/{notice_id}
func (m *NoticeModule) Delete(noticeId string, idempotencyKey ...string) (map[string]interface{}, error) {
	if noticeId == "" {
		return nil, fmt.Errorf("notice_id is required")
	}
	return m.api.deleteWithHeaders(fmt.Sprintf("notices/%s", noticeId), nil,
		commerceRoleHeaders("supervisor", firstOrEmpty(idempotencyKey)))
}
