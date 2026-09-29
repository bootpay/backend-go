package bootpay

import (
	"fmt"
	"net/url"
)

// ProductQnaModule handles 상품문의 operations
//
// 고객 모드는 ProductId 가 필수다. 다른 사람의 비밀글은 is_viewable: false 로 본문 없이 온다 —
// 회원(UserId·LoginId·UserJwt)을 보내면 그 회원이 쓴 비밀글 본문이 보인다.
// Supervisor 가 true 이고 View 가 "all" 이면 몰 전체(숨김 포함)를 보며, 이때 ProductId 는 선택 필터다.
//
// 작성은 회원 또는 비회원(GuestName + GuestPassword, 몰이 비회원 작성을 허용할 때)이 할 수 있다.
type ProductQnaModule struct {
	api *CommerceApi
}

// List retrieves the product Q&A list
// GET /v1/product-qnas
// page/limit 은 미지정 시 1/20 으로 항상 실린다.
func (m *ProductQnaModule) List(params *ProductQnaListParams) (map[string]interface{}, error) {
	if params == nil {
		params = &ProductQnaListParams{}
	}
	query := url.Values{}
	if params.ProductId != "" {
		query.Set("product_id", params.ProductId)
	}
	if params.View != "" {
		query.Set("view", params.View)
	}
	boardPageParams(query, params.Page, params.Limit)
	if params.UserId != "" {
		query.Set("user_id", params.UserId)
	}
	if params.LoginId != "" {
		query.Set("login_id", params.LoginId)
	}
	return m.api.getWithHeaders(withQuery("product-qnas", query),
		commerceBoardHeaders(params.Supervisor, params.UserJwt, params.IdempotencyKey))
}

// Detail retrieves a single product Q&A
// GET /v1/product-qnas/{product_qna_id}
// ⚠️ 다른 사람의 비밀글은 403 PRODUCT_QNA_SECRET_FORBIDDEN.
// Supervisor 면 비밀글·숨김 글도 본다.
func (m *ProductQnaModule) Detail(params ProductQnaDetailParams) (map[string]interface{}, error) {
	if params.ProductQnaId == "" {
		return nil, fmt.Errorf("product_qna_id is required")
	}
	query := url.Values{}
	if params.UserId != "" {
		query.Set("user_id", params.UserId)
	}
	if params.LoginId != "" {
		query.Set("login_id", params.LoginId)
	}
	return m.api.getWithHeaders(withQuery(fmt.Sprintf("product-qnas/%s", params.ProductQnaId), query),
		commerceBoardHeaders(params.Supervisor, params.UserJwt, params.IdempotencyKey))
}

// Create writes a product Q&A
// POST /v1/product-qnas
// 회원(UserId·LoginId·UserJwt) 또는 비회원(GuestName + GuestPassword)으로 쓴다.
// ⚠️ 회원 문의면 Guest* 는 무시된다.
func (m *ProductQnaModule) Create(params ProductQnaCreateParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("product-qnas", params,
		commerceBoardHeaders(false, params.UserJwt, params.IdempotencyKey))
}

// Update updates a product Q&A (작성 회원 본인 또는 비회원 비밀번호, 답변 전만)
// PUT /v1/product-qnas/{product_qna_id}
// ⚠️ 답변이 달린 문의는 409 PRODUCT_QNA_ALREADY_ANSWERED,
// 비회원 비밀번호가 틀리면 403 PRODUCT_QNA_GUEST_PASSWORD_INVALID.
func (m *ProductQnaModule) Update(params ProductQnaUpdateParams) (map[string]interface{}, error) {
	if params.ProductQnaId == "" {
		return nil, fmt.Errorf("product_qna_id is required")
	}
	return m.api.putWithHeaders(fmt.Sprintf("product-qnas/%s", params.ProductQnaId), params,
		commerceBoardHeaders(false, params.UserJwt, params.IdempotencyKey))
}

// Delete deletes a product Q&A (작성 회원 본인·비회원 비밀번호 또는 supervisor)
// DELETE /v1/product-qnas/{product_qna_id}
// ⚠️ GuestPassword 는 URL 쿼리가 아닌 본문으로 보낸다(서버는 둘 다 받지만, 쿼리는 접근로그에 남는다).
func (m *ProductQnaModule) Delete(params ProductQnaDeleteParams) (map[string]interface{}, error) {
	if params.ProductQnaId == "" {
		return nil, fmt.Errorf("product_qna_id is required")
	}
	return m.api.deleteWithHeaders(fmt.Sprintf("product-qnas/%s", params.ProductQnaId), params,
		commerceBoardHeaders(params.Supervisor, params.UserJwt, params.IdempotencyKey))
}

// Answer registers or updates the answer of a product Q&A (supervisor scope)
// PUT /v1/product-qnas/{product_qna_id}/answer
// 답변이 없으면 만들고, 있으면 내용을 바꾼다. 응답은 답변이 반영된 상품문의 객체다.
func (m *ProductQnaModule) Answer(productQnaId string, content string, idempotencyKey ...string) (map[string]interface{}, error) {
	if productQnaId == "" {
		return nil, fmt.Errorf("product_qna_id is required")
	}
	body := map[string]interface{}{"content": content}
	return m.api.putWithHeaders(fmt.Sprintf("product-qnas/%s/answer", productQnaId), body,
		commerceRoleHeaders("supervisor", firstOrEmpty(idempotencyKey)))
}
