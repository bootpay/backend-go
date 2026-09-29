package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
)

// InquiryModule handles 1:1 문의 operations
//
// 회원 모드는 UserId(Bootpay 회원 _id 또는 외부 회원 ID)·LoginId·UserJwt 중 하나로 회원을 정하고
// 본인 문의만 다룬다. Supervisor 가 true 면 몰 전체 문의를 본다 — 이때 UserId 는 작성 회원 필터(선택)다.
//
// ⚠️ 작성·수정은 회원 전용(BOOTPAY-ROLE: user)이고, 답변은 supervisor 전용이다.
type InquiryModule struct {
	api *CommerceApi
}

// List retrieves the 1:1 inquiry list
// GET /v1/inquiries
// page/limit 은 미지정 시 1/20 으로 항상 실린다.
// Answered 는 tri-state 다 — true 답변 완료만 / false 미답변만 / nil 전체.
func (m *InquiryModule) List(params *InquiryListParams) (map[string]interface{}, error) {
	if params == nil {
		params = &InquiryListParams{}
	}
	query := url.Values{}
	if params.UserId != "" {
		query.Set("user_id", params.UserId)
	}
	if params.LoginId != "" {
		query.Set("login_id", params.LoginId)
	}
	if params.Answered != nil {
		query.Set("answered", strconv.FormatBool(*params.Answered))
	}
	boardPageParams(query, params.Page, params.Limit)
	return m.api.getWithHeaders(withQuery("inquiries", query),
		commerceBoardHeaders(params.Supervisor, params.UserJwt, params.IdempotencyKey))
}

// Detail retrieves a single 1:1 inquiry
// GET /v1/inquiries/{inquiry_id}
// 회원 모드는 작성 회원 본인 문의만 본다(다른 회원 문의는 404). Supervisor 면 몰의 모든 문의.
func (m *InquiryModule) Detail(params InquiryDetailParams) (map[string]interface{}, error) {
	if params.InquiryId == "" {
		return nil, fmt.Errorf("inquiry_id is required")
	}
	query := url.Values{}
	if params.UserId != "" {
		query.Set("user_id", params.UserId)
	}
	if params.LoginId != "" {
		query.Set("login_id", params.LoginId)
	}
	return m.api.getWithHeaders(withQuery(fmt.Sprintf("inquiries/%s", params.InquiryId), query),
		commerceBoardHeaders(params.Supervisor, params.UserJwt, params.IdempotencyKey))
}

// Create writes a 1:1 inquiry (회원 전용)
// POST /v1/inquiries
// ProductId 는 이 몰의 상품, OrderId 는 작성 회원의 주문이어야 한다.
// Option 은 상품 옵션 문구다(ProductId 가 함께 필요하다).
func (m *InquiryModule) Create(params InquiryCreateParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("inquiries", params,
		commerceBoardHeaders(false, params.UserJwt, params.IdempotencyKey))
}

// Update updates a 1:1 inquiry (작성 회원 본인, 답변 전만)
// PUT /v1/inquiries/{inquiry_id}
// ⚠️ 답변이 달린 문의는 409 INQUIRY_ALREADY_ANSWERED.
// Title 에 빈 문자열을 보내면 제목을 지운다 — 그래서 pointer type 이다(StringPtr("")).
func (m *InquiryModule) Update(params InquiryUpdateParams) (map[string]interface{}, error) {
	if params.InquiryId == "" {
		return nil, fmt.Errorf("inquiry_id is required")
	}
	return m.api.putWithHeaders(fmt.Sprintf("inquiries/%s", params.InquiryId), params,
		commerceBoardHeaders(false, params.UserJwt, params.IdempotencyKey))
}

// Delete deletes a 1:1 inquiry (작성 회원 본인 또는 supervisor)
// DELETE /v1/inquiries/{inquiry_id}
// 달린 답변도 함께 삭제된다. ⚠️ 회원 식별값은 body 가 아니라 query 로 간다.
func (m *InquiryModule) Delete(params InquiryDeleteParams) (map[string]interface{}, error) {
	if params.InquiryId == "" {
		return nil, fmt.Errorf("inquiry_id is required")
	}
	query := url.Values{}
	if params.UserId != "" {
		query.Set("user_id", params.UserId)
	}
	if params.LoginId != "" {
		query.Set("login_id", params.LoginId)
	}
	return m.api.deleteWithHeaders(withQuery(fmt.Sprintf("inquiries/%s", params.InquiryId), query), nil,
		commerceBoardHeaders(params.Supervisor, params.UserJwt, params.IdempotencyKey))
}

// Answer registers or updates the answer of a 1:1 inquiry (supervisor scope)
// PUT /v1/inquiries/{inquiry_id}/answer
// 답변이 없으면 만들고, 있으면 내용을 바꾼다. 응답은 답변이 반영된 문의 객체다.
func (m *InquiryModule) Answer(inquiryId string, content string, idempotencyKey ...string) (map[string]interface{}, error) {
	if inquiryId == "" {
		return nil, fmt.Errorf("inquiry_id is required")
	}
	body := map[string]interface{}{"content": content}
	return m.api.putWithHeaders(fmt.Sprintf("inquiries/%s/answer", inquiryId), body,
		commerceRoleHeaders("supervisor", firstOrEmpty(idempotencyKey)))
}
