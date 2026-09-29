package bootpay

import (
	"net/http"
	"testing"
)

// 게시판 API 27종(FAQ 5 · 공지사항 5 · 1:1 문의 6 · 상품문의 6 · 상품평 6)이
// 각각 약속된 method · uri · BOOTPAY-ROLE 로 나가는지 확인한다.
//
// role 이 특히 중요하다 — 같은 경로를 고객 모드와 운영자 모드가 나눠 쓰기 때문에,
// 헤더를 붙이지 않으면 운영자 전용 조회가 인스턴스 기본값(user)으로 조용히 나가고
// 서버는 비공개 글을 빼고 응답한다(에러가 아니라서 눈에 띄지 않는다).
func TestCommerceBoardEndpointContract(t *testing.T) {
	cases := []struct {
		name   string
		call   func(api *CommerceApi) (map[string]interface{}, error)
		method string
		url    string
		role   string
	}{
		// FAQ
		{"FaqList", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Faq.List(nil)
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/faqs?limit=20&page=1", "user"},
		{"FaqListSupervisor", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Faq.List(&FaqListParams{Page: 2, Limit: 50, Keyword: "배송", View: "all", Supervisor: true})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/faqs?keyword=%EB%B0%B0%EC%86%A1&limit=50&page=2&view=all", "supervisor"},
		{"FaqDetail", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Faq.Detail("f1", false)
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/faqs/f1", "user"},
		{"FaqDetailSupervisor", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Faq.Detail("f1", true)
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/faqs/f1", "supervisor"},
		{"FaqCreate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Faq.Create(FaqCreateParams{Title: "배송 문의", Content: "본문"})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/faqs", "supervisor"},
		{"FaqUpdate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Faq.Update(FaqUpdateParams{FaqId: "f1", Title: "변경"})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/faqs/f1", "supervisor"},
		{"FaqDelete", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Faq.Delete("f1")
		}, http.MethodDelete, COMMERCE_DEVELOPMENT + "/faqs/f1", "supervisor"},

		// 공지사항
		{"NoticeList", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Notice.List(nil)
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/notices?limit=20&page=1", "user"},
		{"NoticeListSupervisor", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Notice.List(&NoticeListParams{View: "all", Supervisor: true})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/notices?limit=20&page=1&view=all", "supervisor"},
		{"NoticeDetail", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Notice.Detail("n1", false)
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/notices/n1", "user"},
		{"NoticeCreate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Notice.Create(NoticeCreateParams{Title: "점검 안내", Content: "본문"})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/notices", "supervisor"},
		{"NoticeUpdate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Notice.Update(NoticeUpdateParams{NoticeId: "n1", Content: "변경"})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/notices/n1", "supervisor"},
		{"NoticeDelete", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Notice.Delete("n1")
		}, http.MethodDelete, COMMERCE_DEVELOPMENT + "/notices/n1", "supervisor"},

		// 1:1 문의
		{"InquiryList", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Inquiry.List(&InquiryListParams{UserId: "u1"})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/inquiries?limit=20&page=1&user_id=u1", "user"},
		{"InquiryListSupervisor", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Inquiry.List(&InquiryListParams{Answered: BoolPtr(false), Supervisor: true})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/inquiries?answered=false&limit=20&page=1", "supervisor"},
		{"InquiryDetail", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Inquiry.Detail(InquiryDetailParams{InquiryId: "i1", LoginId: "hong"})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/inquiries/i1?login_id=hong", "user"},
		{"InquiryCreate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Inquiry.Create(InquiryCreateParams{Content: "문의 본문", UserId: "u1"})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/inquiries", "user"},
		{"InquiryUpdate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Inquiry.Update(InquiryUpdateParams{InquiryId: "i1", Content: "수정"})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/inquiries/i1", "user"},
		{"InquiryDelete", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Inquiry.Delete(InquiryDeleteParams{InquiryId: "i1", UserId: "u1"})
		}, http.MethodDelete, COMMERCE_DEVELOPMENT + "/inquiries/i1?user_id=u1", "user"},
		{"InquiryDeleteSupervisor", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Inquiry.Delete(InquiryDeleteParams{InquiryId: "i1", Supervisor: true})
		}, http.MethodDelete, COMMERCE_DEVELOPMENT + "/inquiries/i1", "supervisor"},
		{"InquiryAnswer", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.Inquiry.Answer("i1", "답변 본문")
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/inquiries/i1/answer", "supervisor"},

		// 상품문의
		{"ProductQnaList", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductQna.List(&ProductQnaListParams{ProductId: "p1"})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/product-qnas?limit=20&page=1&product_id=p1", "user"},
		{"ProductQnaListSupervisor", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductQna.List(&ProductQnaListParams{View: "all", Supervisor: true})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/product-qnas?limit=20&page=1&view=all", "supervisor"},
		{"ProductQnaDetail", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductQna.Detail(ProductQnaDetailParams{ProductQnaId: "q1", UserId: "u1"})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/product-qnas/q1?user_id=u1", "user"},
		{"ProductQnaCreate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductQna.Create(ProductQnaCreateParams{ProductId: "p1", Content: "문의"})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/product-qnas", "user"},
		{"ProductQnaUpdate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductQna.Update(ProductQnaUpdateParams{ProductQnaId: "q1", Content: "수정"})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/product-qnas/q1", "user"},
		{"ProductQnaDelete", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductQna.Delete(ProductQnaDeleteParams{ProductQnaId: "q1", GuestPassword: "1234"})
		}, http.MethodDelete, COMMERCE_DEVELOPMENT + "/product-qnas/q1", "user"},
		{"ProductQnaAnswer", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductQna.Answer("q1", "답변 본문")
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/product-qnas/q1/answer", "supervisor"},

		// 상품평
		{"ProductReviewList", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductReview.List(&ProductReviewListParams{ProductId: "p1"})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/reviews?limit=20&page=1&product_id=p1", "user"},
		{"ProductReviewListSupervisor", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductReview.List(&ProductReviewListParams{View: "all", Supervisor: true})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/reviews?limit=20&page=1&view=all", "supervisor"},
		{"ProductReviewDetail", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductReview.Detail(ProductReviewDetailParams{ProductReviewId: "r1"})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/reviews/r1", "user"},
		{"ProductReviewCreate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductReview.Create(ProductReviewCreateParams{
				OrderId: "o1", ProductId: "p1", Rating: 5, Content: "좋아요",
			})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/reviews", "user"},
		{"ProductReviewUpdate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductReview.Update(ProductReviewUpdateParams{ProductReviewId: "r1", Rating: 4})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/reviews/r1", "user"},
		{"ProductReviewDelete", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductReview.Delete(ProductReviewDeleteParams{ProductReviewId: "r1", Reason: "욕설", Supervisor: true})
		}, http.MethodDelete, COMMERCE_DEVELOPMENT + "/reviews/r1?reason=%EC%9A%95%EC%84%A4", "supervisor"},
		{"ProductReviewReply", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.ProductReview.Reply("r1", "소중한 후기 감사합니다")
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/reviews/r1/reply", "supervisor"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var captured []capturedCommerceRequest
			api := newMockCommerceApi(&captured)
			// 인스턴스 기본 role 이 무엇이든 요청별 role 이 이겨야 한다
			api.SetRole("manager")

			if _, err := tc.call(api); err != nil {
				t.Fatalf("%s: %v", tc.name, err)
			}
			req := lastRequest(t, captured)
			if req.Method != tc.method || req.URL != tc.url {
				t.Fatalf("%s: got %s %s, want %s %s", tc.name, req.Method, req.URL, tc.method, tc.url)
			}
			if got := req.Header.Get("BOOTPAY-ROLE"); got != tc.role {
				t.Fatalf("%s: BOOTPAY-ROLE = %q, want %q", tc.name, got, tc.role)
			}
			if req.Header.Get("Idempotency-Key") == "" {
				t.Fatalf("%s: Idempotency-Key is not attached", tc.name)
			}
		})
	}
}

// 게시판 API 는 회원 JWT 가 있을 때만 Bootpay-User-JWT 를 붙인다.
// 빈 문자열을 헤더로 내보내면 서버가 JWT 파싱에 실패해 비회원 조회가 401 이 된다.
func TestCommerceBoardUserJwtHeader(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.Inquiry.List(&InquiryListParams{UserJwt: "jwt-token"}); err != nil {
		t.Fatal(err)
	}
	if got := lastRequest(t, captured).Header.Get("Bootpay-User-JWT"); got != "jwt-token" {
		t.Fatalf("Bootpay-User-JWT = %q, want %q", got, "jwt-token")
	}

	if _, err := api.Inquiry.List(&InquiryListParams{}); err != nil {
		t.Fatal(err)
	}
	if _, exists := lastRequest(t, captured).Header["Bootpay-User-Jwt"]; exists {
		t.Fatal("Bootpay-User-JWT must not be attached when no member JWT is given")
	}

	// supervisor 전용 경로(답변·답글)는 JWT 없이 supervisor 로만 나간다
	if _, err := api.ProductReview.Reply("r1", "답글 본문입니다"); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if _, exists := req.Header["Bootpay-User-Jwt"]; exists {
		t.Fatal("supervisor reply must not carry Bootpay-User-JWT")
	}
	if body := decodeBody(t, req); body["content"] != "답글 본문입니다" {
		t.Fatalf("reply content missing: %+v", body)
	}
}

// images 는 3-state 다 — nil(그대로 둠) · 빈 배열(모두 삭제) · 값(통째 교체).
// omitempty 를 그대로 썼다면 빈 배열이 사라져 "사진 전부 삭제" 를 표현할 수 없다.
func TestCommerceBoardImagesTriState(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.Faq.Update(FaqUpdateParams{FaqId: "f1", Title: "제목만 변경"}); err != nil {
		t.Fatal(err)
	}
	if body := decodeBody(t, lastRequest(t, captured)); body["images"] != nil {
		t.Fatalf("nil images must not be sent: %+v", body)
	}

	if _, err := api.Faq.Update(FaqUpdateParams{FaqId: "f1", Images: []interface{}{}}); err != nil {
		t.Fatal(err)
	}
	body := decodeBody(t, lastRequest(t, captured))
	images, ok := body["images"].([]interface{})
	if !ok || len(images) != 0 {
		t.Fatalf("empty images must be sent as [] (모두 삭제): %+v", body)
	}

	if _, err := api.ProductReview.Create(ProductReviewCreateParams{
		OrderId: "o1", ProductId: "p1", Rating: 5, Content: "좋아요",
		Images: []interface{}{"https://img/1.png", map[string]interface{}{"url": "https://img/2.png"}},
	}); err != nil {
		t.Fatal(err)
	}
	body = decodeBody(t, lastRequest(t, captured))
	images, ok = body["images"].([]interface{})
	if !ok || len(images) != 2 || images[0] != "https://img/1.png" {
		t.Fatalf("images must be sent as given: %+v", body)
	}
}

// is_display · is_notice · is_secret 는 explicit false 가 "끔" 이라서 pointer type 이다.
// 값 타입이었다면 omitempty 가 false 를 지워 "끔" 요청이 통째로 사라진다.
func TestCommerceBoardExplicitFalseFlags(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.Notice.Create(NoticeCreateParams{
		Title: "점검 안내", Content: "본문",
		IsNotice: BoolPtr(false), IsDisplay: BoolPtr(false),
	}); err != nil {
		t.Fatal(err)
	}
	body := decodeBody(t, lastRequest(t, captured))
	if body["is_notice"] != false || body["is_display"] != false {
		t.Fatalf("explicit false must be sent: %+v", body)
	}

	if _, err := api.Notice.Create(NoticeCreateParams{Title: "점검 안내", Content: "본문"}); err != nil {
		t.Fatal(err)
	}
	body = decodeBody(t, lastRequest(t, captured))
	if _, exists := body["is_notice"]; exists {
		t.Fatalf("unset flags must not be sent: %+v", body)
	}
	if body["title"] != "점검 안내" || body["content"] != "본문" {
		t.Fatalf("required fields missing: %+v", body)
	}

	if _, err := api.ProductQna.Create(ProductQnaCreateParams{
		ProductId: "p1", Content: "문의", IsSecret: BoolPtr(false),
	}); err != nil {
		t.Fatal(err)
	}
	if body = decodeBody(t, lastRequest(t, captured)); body["is_secret"] != false {
		t.Fatalf("explicit is_secret=false must be sent: %+v", body)
	}
}

// 1:1 문의 수정의 title 은 빈 문자열이 "제목 삭제" 다 — 미전송(nil)과 구분되어야 한다.
func TestCommerceBoardInquiryUpdateTitleClear(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.Inquiry.Update(InquiryUpdateParams{InquiryId: "i1", Title: StringPtr(""), Content: "본문"}); err != nil {
		t.Fatal(err)
	}
	body := decodeBody(t, lastRequest(t, captured))
	title, exists := body["title"]
	if !exists || title != "" {
		t.Fatalf(`title must be sent as "" to clear it: %+v`, body)
	}

	if _, err := api.Inquiry.Update(InquiryUpdateParams{InquiryId: "i1", Content: "본문"}); err != nil {
		t.Fatal(err)
	}
	if _, exists := decodeBody(t, lastRequest(t, captured))["title"]; exists {
		t.Fatalf("nil title must not be sent: %s", string(lastRequest(t, captured).Body))
	}
}

// 삭제 요청은 엔드포인트마다 실리는 자리가 다르다.
//   - 상품문의: guest_password 를 **본문**으로 (쿼리는 접근로그에 비밀번호가 남는다)
//   - 1:1 문의 · 상품평: 회원 식별값·사유를 **쿼리**로
func TestCommerceBoardDeletePayloadPlacement(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.ProductQna.Delete(ProductQnaDeleteParams{
		ProductQnaId: "q1", GuestPassword: "1234", LoginId: "hong",
	}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/product-qnas/q1" {
		t.Fatalf("guest_password must not leak into the query: %s", req.URL)
	}
	body := decodeBody(t, req)
	if body["guest_password"] != "1234" || body["login_id"] != "hong" {
		t.Fatalf("delete payload must be sent in the body: %+v", body)
	}

	if _, err := api.ProductReview.Delete(ProductReviewDeleteParams{
		ProductReviewId: "r1", LoginId: "hong", Reason: "abuse", Supervisor: true,
	}); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/reviews/r1?login_id=hong&reason=abuse" {
		t.Fatalf("review delete params must go into the query: %s", req.URL)
	}
	if len(req.Body) != 0 {
		t.Fatalf("review delete must not carry a body: %s", string(req.Body))
	}
}

// 명시한 Idempotency-Key 는 그대로 전달되고 바디에는 실리지 않는다.
func TestCommerceBoardExplicitIdempotencyKey(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.Faq.Create(FaqCreateParams{Title: "제목", Content: "본문", IdempotencyKey: "fixed-key"}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if got := req.Header.Get("Idempotency-Key"); got != "fixed-key" {
		t.Fatalf("explicit Idempotency-Key must win: %q", got)
	}
	if _, exists := decodeBody(t, req)["idempotency_key"]; exists {
		t.Fatalf("Idempotency-Key must not be serialized into the body: %s", string(req.Body))
	}

	if _, err := api.Faq.Detail("f1", true, "detail-key"); err != nil {
		t.Fatal(err)
	}
	if got := lastRequest(t, captured).Header.Get("Idempotency-Key"); got != "detail-key" {
		t.Fatalf("explicit Idempotency-Key must win: %q", got)
	}
}

// 식별자가 비면 네트워크로 나가기 전에 막는다 — 빈 id 는 상위 컬렉션 경로를 때린다
// (DELETE /v1/faqs/ 가 DELETE /v1/faqs 로 해석되는 사고를 막기 위함).
func TestCommerceBoardRequiredIdValidation(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	calls := map[string]func() (map[string]interface{}, error){
		"FaqDetail":           func() (map[string]interface{}, error) { return api.Faq.Detail("", false) },
		"FaqUpdate":           func() (map[string]interface{}, error) { return api.Faq.Update(FaqUpdateParams{}) },
		"FaqDelete":           func() (map[string]interface{}, error) { return api.Faq.Delete("") },
		"NoticeDetail":        func() (map[string]interface{}, error) { return api.Notice.Detail("", false) },
		"NoticeUpdate":        func() (map[string]interface{}, error) { return api.Notice.Update(NoticeUpdateParams{}) },
		"NoticeDelete":        func() (map[string]interface{}, error) { return api.Notice.Delete("") },
		"InquiryDetail":       func() (map[string]interface{}, error) { return api.Inquiry.Detail(InquiryDetailParams{}) },
		"InquiryUpdate":       func() (map[string]interface{}, error) { return api.Inquiry.Update(InquiryUpdateParams{}) },
		"InquiryDelete":       func() (map[string]interface{}, error) { return api.Inquiry.Delete(InquiryDeleteParams{}) },
		"InquiryAnswer":       func() (map[string]interface{}, error) { return api.Inquiry.Answer("", "답변") },
		"ProductQnaDetail":    func() (map[string]interface{}, error) { return api.ProductQna.Detail(ProductQnaDetailParams{}) },
		"ProductQnaUpdate":    func() (map[string]interface{}, error) { return api.ProductQna.Update(ProductQnaUpdateParams{}) },
		"ProductQnaDelete":    func() (map[string]interface{}, error) { return api.ProductQna.Delete(ProductQnaDeleteParams{}) },
		"ProductQnaAnswer":    func() (map[string]interface{}, error) { return api.ProductQna.Answer("", "답변") },
		"ProductReviewDetail": func() (map[string]interface{}, error) { return api.ProductReview.Detail(ProductReviewDetailParams{}) },
		"ProductReviewUpdate": func() (map[string]interface{}, error) { return api.ProductReview.Update(ProductReviewUpdateParams{}) },
		"ProductReviewDelete": func() (map[string]interface{}, error) { return api.ProductReview.Delete(ProductReviewDeleteParams{}) },
		"ProductReviewReply":  func() (map[string]interface{}, error) { return api.ProductReview.Reply("", "답글") },
	}

	for name, call := range calls {
		if _, err := call(); err == nil {
			t.Fatalf("%s: empty id must return an error", name)
		}
	}
	if len(captured) != 0 {
		t.Fatalf("no request must be sent for an empty id, got %d", len(captured))
	}
}
