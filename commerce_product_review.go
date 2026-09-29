package bootpay

import (
	"fmt"
	"net/url"
)

// ProductReviewModule handles 상품평 operations
//
// 회원 모드는 UserId·LoginId·UserJwt 로 정한 회원이 쓴 상품평만 본다.
// Supervisor 가 true 면 몰 전체 상품평 — 이때 UserId 는 작성 회원 필터(선택)다.
//
// ⚠️ 상품 상세의 공개 상품평 목록은 여기가 아니라 기존 GET /v1/products/{id}/reviews 를 쓴다.
type ProductReviewModule struct {
	api *CommerceApi
}

// List retrieves the product review list
// GET /v1/reviews
// page/limit 은 미지정 시 1/20 으로 항상 실린다.
func (m *ProductReviewModule) List(params *ProductReviewListParams) (map[string]interface{}, error) {
	if params == nil {
		params = &ProductReviewListParams{}
	}
	query := url.Values{}
	boardPageParams(query, params.Page, params.Limit)
	if params.ProductId != "" {
		query.Set("product_id", params.ProductId)
	}
	if params.View != "" {
		query.Set("view", params.View)
	}
	if params.UserId != "" {
		query.Set("user_id", params.UserId)
	}
	if params.LoginId != "" {
		query.Set("login_id", params.LoginId)
	}
	return m.api.getWithHeaders(withQuery("reviews", query),
		commerceBoardHeaders(params.Supervisor, params.UserJwt, params.IdempotencyKey))
}

// Detail retrieves a single product review
// GET /v1/reviews/{product_review_id}
// 공개 상품평은 누구나, 숨김 상품평은 작성 회원 본인 또는 supervisor 만 본다.
func (m *ProductReviewModule) Detail(params ProductReviewDetailParams) (map[string]interface{}, error) {
	if params.ProductReviewId == "" {
		return nil, fmt.Errorf("product_review_id is required")
	}
	query := url.Values{}
	if params.UserId != "" {
		query.Set("user_id", params.UserId)
	}
	if params.LoginId != "" {
		query.Set("login_id", params.LoginId)
	}
	return m.api.getWithHeaders(withQuery(fmt.Sprintf("reviews/%s", params.ProductReviewId), query),
		commerceBoardHeaders(params.Supervisor, params.UserJwt, params.IdempotencyKey))
}

// Create writes a product review (회원 전용)
// POST /v1/reviews
// OrderId 는 회원의 구매확정 주문, ProductId(·ProductOptionId)는 그 주문에 담긴 상품이어야 한다.
// ⚠️ Images 는 사진 URL 최대 5개다(문자열 또는 {"url": ...} 맵). 파일 업로드는 지원하지 않는다.
func (m *ProductReviewModule) Create(params ProductReviewCreateParams) (map[string]interface{}, error) {
	body, err := boardPayload(params, params.Images)
	if err != nil {
		return nil, err
	}
	return m.api.postWithHeaders("reviews", body,
		commerceBoardHeaders(false, params.UserJwt, params.IdempotencyKey))
}

// Update updates a product review (작성 회원 본인, 작성 후 7일 이내)
// PUT /v1/reviews/{product_review_id}
// Images 는 보내면 통째로 교체한다(빈 슬라이스면 모두 삭제).
// ⚠️ 7일이 지나면 REVIEW_EDIT_TIME_EXPIRED.
func (m *ProductReviewModule) Update(params ProductReviewUpdateParams) (map[string]interface{}, error) {
	if params.ProductReviewId == "" {
		return nil, fmt.Errorf("product_review_id is required")
	}
	body, err := boardPayload(params, params.Images)
	if err != nil {
		return nil, err
	}
	return m.api.putWithHeaders(fmt.Sprintf("reviews/%s", params.ProductReviewId), body,
		commerceBoardHeaders(false, params.UserJwt, params.IdempotencyKey))
}

// Delete deletes a product review (작성 회원 본인 또는 supervisor)
// DELETE /v1/reviews/{product_review_id}
// 지급된 적립금 회수·상품 통계 차감이 함께 일어난다. Reason 은 운영자 모드의 삭제 사유다.
// ⚠️ 회원 식별값·사유는 body 가 아니라 query 로 간다.
func (m *ProductReviewModule) Delete(params ProductReviewDeleteParams) (map[string]interface{}, error) {
	if params.ProductReviewId == "" {
		return nil, fmt.Errorf("product_review_id is required")
	}
	query := url.Values{}
	if params.UserId != "" {
		query.Set("user_id", params.UserId)
	}
	if params.LoginId != "" {
		query.Set("login_id", params.LoginId)
	}
	if params.Reason != "" {
		query.Set("reason", params.Reason)
	}
	return m.api.deleteWithHeaders(withQuery(fmt.Sprintf("reviews/%s", params.ProductReviewId), query), nil,
		commerceBoardHeaders(params.Supervisor, params.UserJwt, params.IdempotencyKey))
}

// Reply registers or updates the seller reply of a product review (supervisor scope)
// PUT /v1/reviews/{product_review_id}/reply
// 답글이 없으면 만들고, 있으면 내용을 바꾼다(10자 이상). 응답은 답글이 반영된 상품평 객체다.
func (m *ProductReviewModule) Reply(productReviewId string, content string, idempotencyKey ...string) (map[string]interface{}, error) {
	if productReviewId == "" {
		return nil, fmt.Errorf("product_review_id is required")
	}
	body := map[string]interface{}{"content": content}
	return m.api.putWithHeaders(fmt.Sprintf("reviews/%s/reply", productReviewId), body,
		commerceRoleHeaders("supervisor", firstOrEmpty(idempotencyKey)))
}
