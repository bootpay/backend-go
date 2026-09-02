package bootpay

import (
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 알림톡 v1 API 35종이 각각 약속된 method · uri 로 나가는지 확인한다.
// 두 가지가 함께 검증된다.
//   - BOOTPAY-ROLE 은 인스턴스 기본값과 무관하게 항상 user (알림톡 스코프 키가 전부 user:alimtalk_*)
//   - Idempotency-Key 는 붙지 않는다 — 알림톡 API 는 이 헤더를 읽지 않는다.
//     붙이면 서버가 주지 않는 멱등 보장을 주는 것처럼 보인다(멱등은 발송의 ref_id 로만 성립).
func TestCommerceAlimtalkEndpointContract(t *testing.T) {
	dir := t.TempDir()
	image := filepath.Join(dir, "template.png")
	if err := os.WriteFile(image, []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name   string
		call   func(api *CommerceApi) (map[string]interface{}, error)
		method string
		url    string
	}{
		// messages
		{"MessageList", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkMessage.List(nil)
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/messages"},
		{"MessageStats", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkMessage.Stats("2026-08-01", "2026-08-27")
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/messages/stats?e_at=2026-08-27&s_at=2026-08-01"},
		{"MessageDetail", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkMessage.Detail("r1")
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/messages/r1"},

		// official catalog
		{"OfficialList", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkOfficial.List(&AlimtalkOfficialListParams{Keyword: "주문", Per: 50})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/official?per=50&q=%EC%A3%BC%EB%AC%B8"},
		{"OfficialRecommend", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkOfficial.Recommend(AlimtalkOfficialRecommendParams{Text: "주문이 접수되었습니다"})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/official/recommend"},
		{"OfficialDetail", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkOfficial.Detail("BP0001", "k1")
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/official/BP0001?ksp_id=k1"},

		// optouts
		{"OptoutList", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkOptout.List(&AlimtalkOptoutListParams{Phone: "0101234", Page: 2})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/optouts?page=2&phone=0101234"},
		{"OptoutCreate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkOptout.Create(AlimtalkOptoutCreateParams{Phone: "01012345678", Reason: "고객요청"})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/optouts"},
		{"OptoutCheck", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkOptout.Check(AlimtalkOptoutCheckParams{Phones: []string{"01012345678"}})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/optouts/check"},
		{"OptoutRelease", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkOptout.Release("01012345678")
		}, http.MethodDelete, COMMERCE_DEVELOPMENT + "/alimtalk/optouts/01012345678"},

		// send
		{"Send", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkSend.Send(AlimtalkSendParams{TemplateCode: "T1", To: "01012345678"})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/send"},
		{"SendBulk", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkSend.Bulk(AlimtalkSendBulkParams{
				TemplateCode: "T1",
				Recipients:   []AlimtalkSendRecipient{{To: "01012345678", RefId: "bulk-0001"}},
			})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/send/bulk"},
		{"SendCancel", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkSend.Cancel("r1")
		}, http.MethodDelete, COMMERCE_DEVELOPMENT + "/alimtalk/send/r1"},

		// senders
		{"SenderCategories", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkSender.Categories()
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/categories"},
		{"SenderOtp", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkSender.Otp(AlimtalkSenderOtpParams{YellowId: "@bootpay", Phone: "01012345678"})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/senders/otp"},
		{"SenderCreate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkSender.Create(AlimtalkSenderCreateParams{
				Otp: "123456", YellowId: "@bootpay", Phone: "01012345678", CategoryCode: "001001",
			})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/senders"},
		{"SenderList", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkSender.List()
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/senders"},
		{"SenderDetail", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkSender.Detail("k1", BoolPtr(true))
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/senders/k1?sync=true"},
		{"SenderRelease", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkSender.Release("k1")
		}, http.MethodDelete, COMMERCE_DEVELOPMENT + "/alimtalk/senders/k1"},
		{"SenderVariableExamples", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkSender.VariableExamples("k1", map[string]interface{}{"user_name": "홍길동"})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/alimtalk/senders/k1/variable_examples"},

		// templates
		{"TemplateList", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkTemplate.List(&AlimtalkTemplateListParams{Ins: "3", Sort: "latest"})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/templates?ins=3&sort=latest"},
		{"TemplateCreate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkTemplate.Create(AlimtalkTemplateCreateParams{KspId: "k1", Name: "주문완료"})
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/templates"},
		{"TemplateDetail", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkTemplate.Detail("t1", BoolPtr(false))
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/templates/t1?sync=false"},
		{"TemplateUpdate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkTemplate.Update("t1", AlimtalkTemplateUpdateParams{Name: "주문완료v2"})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/alimtalk/templates/t1"},
		{"TemplateDelete", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkTemplate.Delete("t1")
		}, http.MethodDelete, COMMERCE_DEVELOPMENT + "/alimtalk/templates/t1"},
		{"TemplateRegister", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkTemplate.Register("t1")
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/templates/t1/register"},
		{"TemplateInspect", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkTemplate.Inspect("t1")
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/templates/t1/inspect"},
		{"TemplateExport", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkTemplate.Export(nil)
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/templates/export?format=json"},
		{"TemplateImage", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkTemplate.Image(image, "")
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/templates/image"},
		{"TemplateHighlightImage", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkTemplate.HighlightImage(image, "")
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/templates/highlight_image"},

		// webhook
		{"WebhookDetail", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkWebhook.Detail()
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/webhook"},
		{"WebhookUpdate", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkWebhook.Update(AlimtalkWebhookUpdateParams{Url: "https://example.com/hook"})
		}, http.MethodPut, COMMERCE_DEVELOPMENT + "/alimtalk/webhook"},
		{"WebhookTest", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkWebhook.Test()
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/webhook/test"},
		{"WebhookRotateSecret", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkWebhook.RotateSecret()
		}, http.MethodPost, COMMERCE_DEVELOPMENT + "/alimtalk/webhook/secret"},
		{"WebhookDeliveries", func(api *CommerceApi) (map[string]interface{}, error) {
			return api.AlimtalkWebhook.Deliveries(&AlimtalkWebhookDeliveriesParams{Page: 1, Limit: 100})
		}, http.MethodGet, COMMERCE_DEVELOPMENT + "/alimtalk/webhook/deliveries?limit=100&page=1"},
	}

	if len(cases) != 35 {
		t.Fatalf("알림톡 SDK 메서드는 35종이다, got %d", len(cases))
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var captured []capturedCommerceRequest
			api := newMockCommerceApi(&captured)
			// 인스턴스 role 이 supervisor 여도 알림톡은 user 로 나가야 한다
			api.SetRole("supervisor")

			if _, err := tc.call(api); err != nil {
				t.Fatalf("%s: %v", tc.name, err)
			}
			req := lastRequest(t, captured)
			if req.Method != tc.method || req.URL != tc.url {
				t.Fatalf("%s: got %s %s, want %s %s", tc.name, req.Method, req.URL, tc.method, tc.url)
			}
			if got := req.Header.Get("BOOTPAY-ROLE"); got != "user" {
				t.Fatalf("%s: BOOTPAY-ROLE = %q, want user", tc.name, got)
			}
			if got := req.Header.Get("Idempotency-Key"); got != "" {
				t.Fatalf("%s: 알림톡 API 는 Idempotency-Key 를 읽지 않는다, got %q", tc.name, got)
			}
		})
	}
}

// fallback 은 미지정(nil)과 false 가 다르다 — nil 이면 프로젝트 기본값을 따르고 false 는 명시적으로 끈다.
// 값 타입이었다면 false 가 omitempty 로 사라져 "껐는데 문자가 나가는" 상태가 된다.
func TestCommerceAlimtalkSendFallbackTriState(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.AlimtalkSend.Send(AlimtalkSendParams{
		TemplateCode: "T1",
		To:           "01012345678",
		Variables:    map[string]interface{}{"user_name": "홍길동"},
		RefId:        "order-1",
	}); err != nil {
		t.Fatal(err)
	}
	body := decodeBody(t, lastRequest(t, captured))
	if _, exists := body["fallback"]; exists {
		t.Fatalf("미지정 fallback 은 전송되면 안 된다: %+v", body)
	}
	if body["template_code"] != "T1" || body["to"] != "01012345678" || body["ref_id"] != "order-1" {
		t.Fatalf("send body mismatch: %+v", body)
	}
	variables, ok := body["variables"].(map[string]interface{})
	if !ok || variables["user_name"] != "홍길동" {
		t.Fatalf("variables mismatch: %+v", body)
	}

	if _, err := api.AlimtalkSend.Send(AlimtalkSendParams{
		TemplateCode: "T1", To: "01012345678", Fallback: BoolPtr(false),
	}); err != nil {
		t.Fatal(err)
	}
	body = decodeBody(t, lastRequest(t, captured))
	if body["fallback"] != false {
		t.Fatalf("명시한 false 는 그대로 실려야 한다: %+v", body)
	}

	if _, err := api.AlimtalkSend.Bulk(AlimtalkSendBulkParams{
		TemplateCode: "T1",
		Recipients: []AlimtalkSendRecipient{
			{To: "01011112222", RefId: "b-1", Variables: map[string]interface{}{"user_name": "김"}},
			{To: "01033334444", RefId: "b-2"},
		},
		Fallback: BoolPtr(true),
	}); err != nil {
		t.Fatal(err)
	}
	body = decodeBody(t, lastRequest(t, captured))
	if body["fallback"] != true {
		t.Fatalf("bulk fallback mismatch: %+v", body)
	}
	recipients, ok := body["recipients"].([]interface{})
	if !ok || len(recipients) != 2 {
		t.Fatalf("recipients mismatch: %+v", body)
	}
	first, _ := recipients[0].(map[string]interface{})
	if first["to"] != "01011112222" || first["ref_id"] != "b-1" {
		t.Fatalf("recipient[0] mismatch: %+v", first)
	}
	second, _ := recipients[1].(map[string]interface{})
	if _, exists := second["variables"]; exists {
		t.Fatalf("빈 variables 는 전송되면 안 된다: %+v", second)
	}
}

// register 를 명시적으로 false 로 주지 않으면 생성 즉시 대행사·카카오에 실제 등록된다.
// 그래서 false 가 바디에서 사라지면 안 된다(pointer type).
func TestCommerceAlimtalkTemplateCreateBody(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.AlimtalkTemplate.Create(AlimtalkTemplateCreateParams{
		KspId:    "k1",
		Name:     "주문완료",
		Content:  "#{user_name}님 주문이 완료되었습니다.",
		Register: BoolPtr(false),
		MsgType:  "BA",
		Buttons: []AlimtalkTemplateButton{
			{Name: "주문조회", LinkType: "WL", LinkMo: "https://m.example.com", Ordering: 1},
		},
		Examples: map[string]interface{}{"user_name": "홍길동"},
		Extra:    map[string]interface{}{"comment": "내부 메모", "dropped": nil},
	}); err != nil {
		t.Fatal(err)
	}
	body := decodeBody(t, lastRequest(t, captured))
	if body["register"] != false {
		t.Fatalf("명시한 register=false 는 그대로 실려야 한다: %+v", body)
	}
	if body["ksp_id"] != "k1" || body["msg_type"] != "BA" {
		t.Fatalf("create body mismatch: %+v", body)
	}
	if body["comment"] != "내부 메모" {
		t.Fatalf("Extra 키가 병합되지 않았다: %+v", body)
	}
	if _, exists := body["dropped"]; exists {
		t.Fatalf("nil Extra 값은 전송되면 안 된다: %+v", body)
	}
	if _, exists := body["security_flag"]; exists {
		t.Fatalf("미지정 값은 전송되면 안 된다: %+v", body)
	}
	buttons, ok := body["buttons"].([]interface{})
	if !ok || len(buttons) != 1 {
		t.Fatalf("buttons mismatch: %+v", body)
	}
	button, _ := buttons[0].(map[string]interface{})
	// ⚠️ 등록 API 는 linkType/linkMo 를 읽는다. 발송 포맷(type/url_mobile)으로 보내면
	//    서버가 linkType 을 못 읽어 "지원하지 않는 버튼 타입입니다" 로 거부한다.
	if button["name"] != "주문조회" || button["linkType"] != "WL" ||
		button["linkMo"] != "https://m.example.com" || button["ordering"] != float64(1) {
		t.Fatalf("button 등록 포맷 mismatch: %+v", button)
	}
	if _, exists := button["type"]; exists {
		t.Fatalf("발송 포맷 키(type)가 등록 요청에 실리면 안 된다: %+v", button)
	}
	if _, exists := button["url_mobile"]; exists {
		t.Fatalf("발송 포맷 키(url_mobile)가 등록 요청에 실리면 안 된다: %+v", button)
	}

	// 수정은 부분 수정이 아니다. 이미지 삭제는 빈 문자열을 명시해야 하는데 omitempty 가 지우므로
	// Extra 로 실어 보낸다 — 그 경로가 실제로 동작하는지 확인한다.
	if _, err := api.AlimtalkTemplate.Update("t1", AlimtalkTemplateUpdateParams{
		Name:  "주문완료v2",
		Extra: map[string]interface{}{"storage_image_url": ""},
	}); err != nil {
		t.Fatal(err)
	}
	body = decodeBody(t, lastRequest(t, captured))
	if body["name"] != "주문완료v2" {
		t.Fatalf("update body mismatch: %+v", body)
	}
	value, exists := body["storage_image_url"]
	if !exists || value != "" {
		t.Fatalf("이미지 삭제 의도(빈 문자열)가 전송되지 않았다: %+v", body)
	}
}

// 공식 카탈로그 검색어는 서버 정본 키인 q 로 나가야 한다.
// 서버는 q 를 먼저 보고 없으면 keyword 를 보는데, keyword 로 보내면 정본 경로를 타지 못한다.
func TestCommerceAlimtalkOfficialListUsesQueryKeyQ(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.AlimtalkOfficial.List(&AlimtalkOfficialListParams{Keyword: "order", MsgType: "BA", Page: 2, KspId: "k1"}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if !strings.Contains(req.URL, "q=order") {
		t.Fatalf("검색어는 q 로 나가야 한다: %s", req.URL)
	}
	if strings.Contains(req.URL, "keyword=") {
		t.Fatalf("keyword 는 보내지 않는다: %s", req.URL)
	}
	if !strings.Contains(req.URL, "msg_type=BA") || !strings.Contains(req.URL, "page=2") || !strings.Contains(req.URL, "ksp_id=k1") {
		t.Fatalf("official list query mismatch: %s", req.URL)
	}
}

// format=csv 는 JSON 이 아니라 공용 파서를 통과하지 못한다.
// 파싱을 시도하면 성공한 요청이 "통신 실패" 로 보고되므로 원문 경로로 빠져야 한다.
func TestCommerceAlimtalkTemplateExportCsvKeepsRawBody(t *testing.T) {
	var captured []capturedCommerceRequest
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			var body []byte
			if req.Body != nil {
				body, _ = io.ReadAll(req.Body)
			}
			captured = append(captured, capturedCommerceRequest{
				Method: req.Method,
				URL:    req.URL.String(),
				Header: req.Header.Clone(),
				Body:   body,
			})
			header := make(http.Header)
			header.Set("Content-Type", "text/csv; charset=utf-8")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("code,name\nT1,주문완료\n")),
				Header:     header,
			}, nil
		}),
	}
	api := NewCommerceAPI("ck", "sk", client, "development")

	result, err := api.AlimtalkTemplate.Export(&AlimtalkTemplateExportParams{
		Format: "csv", Scope: "private", KspId: "k1", IncludeContent: BoolPtr(true),
	})
	if err != nil {
		t.Fatalf("csv export 는 파싱 실패로 떨어지면 안 된다: %v", err)
	}
	if result["body"] != "code,name\nT1,주문완료\n" {
		t.Fatalf("원문이 그대로 담겨야 한다: %+v", result)
	}
	if !strings.HasPrefix(result["content_type"].(string), "text/csv") {
		t.Fatalf("content_type mismatch: %+v", result)
	}

	req := lastRequest(t, captured)
	for _, want := range []string{"format=csv", "scope=private", "ksp_id=k1", "include_content=true"} {
		if !strings.Contains(req.URL, want) {
			t.Fatalf("export query 에 %s 가 없다: %s", want, req.URL)
		}
	}
	if req.Header.Get("BOOTPAY-ROLE") != "user" {
		t.Fatalf("export role mismatch: %q", req.Header.Get("BOOTPAY-ROLE"))
	}
}

// 템플릿 이미지 업로드는 multipart/form-data 이고 Content-Type 의 boundary 가 유지돼야 한다.
// boundary 를 덮어쓰면 서버가 본문을 null 로 파싱한다.
func TestCommerceAlimtalkTemplateImageMultipart(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	dir := t.TempDir()
	image := filepath.Join(dir, "banner.png")
	if err := os.WriteFile(image, []byte("png-bytes"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := api.AlimtalkTemplate.Image(image, "https://cdn.example.com/old.png"); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	if req.Method != http.MethodPost || req.URL != COMMERCE_DEVELOPMENT+"/alimtalk/templates/image" {
		t.Fatalf("image upload mismatch: %s %s", req.Method, req.URL)
	}

	mediaType, mediaParams, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		t.Fatal(err)
	}
	if mediaType != "multipart/form-data" || mediaParams["boundary"] == "" {
		t.Fatalf("boundary 가 유지돼야 한다: %q", req.Header.Get("Content-Type"))
	}

	reader := multipart.NewReader(strings.NewReader(string(req.Body)), mediaParams["boundary"])
	fields := map[string]string{}
	files := map[string]string{}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		content, err := io.ReadAll(part)
		if err != nil {
			t.Fatal(err)
		}
		if part.FileName() != "" {
			files[part.FormName()] = part.FileName()
		} else {
			fields[part.FormName()] = string(content)
		}
	}
	if files["image"] != "banner.png" {
		t.Fatalf("image 파트가 없다: %+v", files)
	}
	if fields["replace_url"] != "https://cdn.example.com/old.png" {
		t.Fatalf("replace_url mismatch: %+v", fields)
	}

	// replace_url 이 비면 파트를 만들지 않는다 (Ruby 의 present? 가드와 같다)
	if _, err := api.AlimtalkTemplate.HighlightImage(image, ""); err != nil {
		t.Fatal(err)
	}
	req = lastRequest(t, captured)
	if req.URL != COMMERCE_DEVELOPMENT+"/alimtalk/templates/highlight_image" {
		t.Fatalf("highlight image uri mismatch: %s", req.URL)
	}
	if strings.Contains(string(req.Body), "replace_url") {
		t.Fatalf("빈 replace_url 은 파트를 만들지 않는다: %s", string(req.Body))
	}
}

// 지정하지 않은 조회 필터는 쿼리에 실리지 않는다 (status=&ref_id= 노이즈 방지).
func TestCommerceAlimtalkMessageListOmitsUnsetFilters(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.AlimtalkMessage.List(&AlimtalkMessageListParams{
		TemplateCode: "T1", Status: "success", To: "01012345678", Limit: 100,
	}); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	want := COMMERCE_DEVELOPMENT + "/alimtalk/messages?limit=100&status=success&template_code=T1&to=01012345678"
	if req.URL != want {
		t.Fatalf("message list query mismatch:\n got %s\nwant %s", req.URL, want)
	}
}

// 웹훅 설정은 enabled=false 를 명시적으로 보낼 수 있어야 한다 (끄기 의도).
func TestCommerceAlimtalkWebhookUpdateBody(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.AlimtalkWebhook.Update(AlimtalkWebhookUpdateParams{
		Url:     "https://example.com/alimtalk-hook",
		Events:  []int{301, 302, 310},
		Enabled: BoolPtr(false),
	}); err != nil {
		t.Fatal(err)
	}
	body := decodeBody(t, lastRequest(t, captured))
	if body["enabled"] != false {
		t.Fatalf("명시한 enabled=false 는 그대로 실려야 한다: %+v", body)
	}
	if body["url"] != "https://example.com/alimtalk-hook" {
		t.Fatalf("url mismatch: %+v", body)
	}
	events, ok := body["events"].([]interface{})
	if !ok || len(events) != 3 || events[0] != float64(301) {
		t.Fatalf("events mismatch: %+v", body)
	}

	// 파라미터 없는 POST 도 빈 JSON 오브젝트를 보낸다 (Ruby payload 기본값 {} 과 동일)
	if _, err := api.AlimtalkWebhook.RotateSecret(); err != nil {
		t.Fatal(err)
	}
	req := lastRequest(t, captured)
	sent := map[string]interface{}{}
	if err := json.Unmarshal(req.Body, &sent); err != nil {
		t.Fatalf("본문이 JSON 이 아니다: %s", string(req.Body))
	}
	if len(sent) != 0 {
		t.Fatalf("빈 오브젝트여야 한다: %+v", sent)
	}
}

// 변수 예문 사전은 examples 키로 감싸 보낸다 (부분 갱신 — 보낸 키만 덮어쓴다).
func TestCommerceAlimtalkSenderVariableExamplesBody(t *testing.T) {
	var captured []capturedCommerceRequest
	api := newMockCommerceApi(&captured)

	if _, err := api.AlimtalkSender.VariableExamples("k1", map[string]interface{}{
		"user_name": "홍길동", "company_name": "부트페이몰",
	}); err != nil {
		t.Fatal(err)
	}
	body := decodeBody(t, lastRequest(t, captured))
	examples, ok := body["examples"].(map[string]interface{})
	if !ok || examples["user_name"] != "홍길동" || examples["company_name"] != "부트페이몰" {
		t.Fatalf("examples mismatch: %+v", body)
	}
}
