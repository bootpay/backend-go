package bootpay

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// AlimtalkTemplateModule handles the merchant's own alimtalk templates
// /v1/alimtalk/templates 계열 (CRUD · 등록 · 검수)
//
// 흐름: (초안 생성 → 확인 → 대행사 등록) → 검수 요청 → 승인(APR) → 발송 가능
//
//	Create 에 Register=false 로 초안만 만들고, 내용을 확인한 뒤 Register 로 올리는 것을 권장한다.
//
// ⚠️ Register 를 명시적으로 false 로 주지 않으면 **생성 즉시 대행사·카카오에 실제 등록**된다.
// ⚠️ 본문 변수는 `#{변수명}` 형식이고 템플릿 전체에서 최대 40개다.
type AlimtalkTemplateModule struct {
	api *CommerceApi
}

// List retrieves the own-template list
// GET /v1/alimtalk/templates
// ⚠️ 페이지네이션이 없다 — 필터에 걸린 템플릿을 한 번에 모두 돌려준다.
func (m *AlimtalkTemplateModule) List(params *AlimtalkTemplateListParams) (map[string]interface{}, error) {
	query := url.Values{}
	if params != nil {
		if params.Ins != "" {
			query.Set("ins", params.Ins)
		}
		if params.Sort != "" {
			query.Set("sort", params.Sort)
		}
		if params.Keyword != "" {
			query.Set("keyword", params.Keyword)
		}
	}
	return m.api.getWithHeaders(withQuery("alimtalk/templates", query), alimtalkHeaders())
}

// Create creates an own template
// POST /v1/alimtalk/templates
// ⚠️ params.Register 를 false 로 주지 않으면 대행사·카카오에 **실제 등록**된다(되돌리려면 삭제해야 한다).
func (m *AlimtalkTemplateModule) Create(params AlimtalkTemplateCreateParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("alimtalk/templates", params, alimtalkHeaders())
}

// Detail retrieves an own template
// GET /v1/alimtalk/templates/{template_id}
// templateId 는 문서 id 이고, ObjectId 형식이 아니면 **템플릿 코드**로 해석한다.
// ⚠️ sync 는 서버 기본값이 **true** 라 조회만 해도 벤더 상태 동기화가 일어난다.
//
//	초안(등록 전)을 조회할 때는 sync 를 false 로 넘기는 것을 권장한다.
func (m *AlimtalkTemplateModule) Detail(templateId string, sync *bool) (map[string]interface{}, error) {
	query := url.Values{}
	if sync != nil {
		query.Set("sync", strconv.FormatBool(*sync))
	}
	return m.api.getWithHeaders(withQuery(fmt.Sprintf("alimtalk/templates/%s", templateId), query), alimtalkHeaders())
}

// Update updates an own template
// PUT /v1/alimtalk/templates/{template_id}
// ⚠️ **부분 수정이 아니다.** 보내지 않은 필드는 nil 로 덮어써지므로 항상 전체 필드를 보낸다.
// ⚠️ 수정 가능 상태는 초안 / REG(등록) / REJ(승인반려) / KRR(등록거절) 뿐이다 — APR·REQ 는 거부된다.
func (m *AlimtalkTemplateModule) Update(templateId string, params AlimtalkTemplateUpdateParams) (map[string]interface{}, error) {
	return m.api.putWithHeaders(fmt.Sprintf("alimtalk/templates/%s", templateId), params, alimtalkHeaders())
}

// Delete deletes an own template
// DELETE /v1/alimtalk/templates/{template_id}
// 초안(등록 전)은 대행사 거부와 무관하게 로컬에서 삭제된다.
// ⚠️ 등록분은 **대행사 삭제가 성공해야** 삭제된다 — 승인(APR) 템플릿은 카카오가 거부하므로
//
//	500(3013)이 오고 템플릿은 남는다. 같은 코드가 대행사에 선점된 채 로컬만 사라지는 것을 막기 위함이다.
func (m *AlimtalkTemplateModule) Delete(templateId string) (map[string]interface{}, error) {
	return m.api.deleteWithHeaders(fmt.Sprintf("alimtalk/templates/%s", templateId), nil, alimtalkHeaders())
}

// Register registers a draft with the agency
// POST /v1/alimtalk/templates/{template_id}/register
// ⚠️ 대행사·카카오에 실제 등록된다. 등록 전(초안) 상태에서만 호출할 수 있다.
func (m *AlimtalkTemplateModule) Register(templateId string) (map[string]interface{}, error) {
	return m.api.postWithHeaders(fmt.Sprintf("alimtalk/templates/%s/register", templateId), map[string]interface{}{}, alimtalkHeaders())
}

// Inspect requests inspection
// POST /v1/alimtalk/templates/{template_id}/inspect
// ⚠️ **카카오에 검수를 요청하며 취소할 수 없다.**
// 대행사 등록이 끝난 대기(R) + REG(등록) 상태에서만 호출할 수 있다 — 초안은 먼저 Register 를 부른다.
// 반려(REJ/KRR)된 건은 재요청이 아니라 **수정 후 재요청**이다. 반려 사유는 응답의 comments 에 담긴다.
func (m *AlimtalkTemplateModule) Inspect(templateId string) (map[string]interface{}, error) {
	return m.api.postWithHeaders(fmt.Sprintf("alimtalk/templates/%s/inspect", templateId), map[string]interface{}{}, alimtalkHeaders())
}

// Export exports the template list
// GET /v1/alimtalk/templates/export
// ⚠️ 기본 format 을 **json 으로 둔다** — 서버 기본은 csv 지만, csv 본문은 JSON 이 아니라서
// 공용 파서를 통과하지 못한다. csv 를 주면 파싱 없이 { "body", "content_type" } 으로 원문을 돌려준다.
// 1회 5,000건을 넘으면 3031 로 거부되므로 채널·상태 필터로 좁힌다.
func (m *AlimtalkTemplateModule) Export(params *AlimtalkTemplateExportParams) (map[string]interface{}, error) {
	format := "json"
	query := url.Values{}
	if params != nil {
		if params.Format != "" {
			format = params.Format
		}
		if params.Scope != "" {
			query.Set("scope", params.Scope)
		}
		if params.KspId != "" {
			query.Set("ksp_id", params.KspId)
		}
		if params.Status != "" {
			query.Set("status", params.Status)
		}
		if params.IncludeContent != nil {
			query.Set("include_content", strconv.FormatBool(*params.IncludeContent))
		}
	}
	query.Set("format", format)

	uri := withQuery("alimtalk/templates/export", query)
	if strings.EqualFold(format, "csv") {
		return m.api.getRaw(uri, alimtalkHeaders())
	}
	return m.api.getWithHeaders(uri, alimtalkHeaders())
}

// Image uploads the source image of an image-type template
// POST /v1/alimtalk/templates/image
// 돌려받은 image_url 을 템플릿 생성/수정의 StorageImageUrl 로 넘긴다.
// 규격을 업로드 **전에** 서버가 검사한다 — jpg/png · 500KB 이하 · 가로 500px 이상 · 2:1.
// replaceUrl 을 주면 업로드 성공 후에 기존 파일을 지운다.
func (m *AlimtalkTemplateModule) Image(imagePath string, replaceUrl string) (map[string]interface{}, error) {
	return m.uploadImage("alimtalk/templates/image", imagePath, replaceUrl)
}

// HighlightImage uploads the highlight thumbnail of an item-list template
// POST /v1/alimtalk/templates/highlight_image
// ⚠️ 본문 이미지와 **규격이 다르다** — jpg/png · 500KB 이하 · 가로 **108px** 이상 · **1:1**.
//
//	본문 이미지 엔드포인트로 올리면 거부된다.
//
// 돌려받은 image_url 은 item_highlight.storage_image_url 로 넘긴다.
// ⚠️ 썸네일을 붙이면 하이라이트 글자 한도가 줄어든다(타이틀 30→21, 설명 19→13).
func (m *AlimtalkTemplateModule) HighlightImage(imagePath string, replaceUrl string) (map[string]interface{}, error) {
	return m.uploadImage("alimtalk/templates/highlight_image", imagePath, replaceUrl)
}

// uploadImage posts an image as multipart/form-data.
// ⚠️ Content-Type must keep the boundary generated by the writer — overwriting it makes
// the server parse the body as null.
func (m *AlimtalkTemplateModule) uploadImage(uri string, imagePath string, replaceUrl string) (map[string]interface{}, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	if replaceUrl != "" {
		if err := writer.WriteField("replace_url", replaceUrl); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	return m.api.postMultipart(uri, body, writer.FormDataContentType(), alimtalkHeaders())
}
