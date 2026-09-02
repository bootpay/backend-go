package bootpay

import "fmt"

// AlimtalkSendModule handles alimtalk sending
// POST /v1/alimtalk/send · /send/bulk · DELETE /send/{receipt_id}
//
// ⚠️ **실제로 카카오톡이 발송되고 과금된다. 샌드박스가 없다.**
//
// 처리 순서: 멱등 확인 → 템플릿·채널 해석 → 발송권한 → 지갑 자격 → 발송제어 → 폴백 확정(발신번호 확보)
//
//		→ 수신거부 대조 → 변수 치환·규격검증 → 접수(READY) → 워커 전송
//
//	  - **멱등**: 같은 (프로젝트, ref_id) 로 재요청하면 기존 receipt 를 그대로 돌려준다. 실패한 건만 재발송된다.
//	  - **필수 변수**: 템플릿 응답의 required_variables 를 모두 채워야 한다. 하나라도 비면 3017 로 거부된다.
//	    ⚠️ 다만 실제로 치환되어 나가는 건 본문·강조 타이틀·버튼 링크뿐이다 — 보조문구와 아이템리스트형
//	    요소는 발송 페이로드에 자리가 없어 카카오가 등록된 템플릿 문구 그대로 렌더한다.
type AlimtalkSendModule struct {
	api *CommerceApi
}

// Send sends a single alimtalk message
// POST /v1/alimtalk/send
// 응답: { receipt_id:, ref_id:, to:, status: } — 접수 직후 status 는 requested
func (m *AlimtalkSendModule) Send(params AlimtalkSendParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("alimtalk/send", params, alimtalkHeaders())
}

// Bulk sends one request to N recipients
// POST /v1/alimtalk/send/bulk
// ⚠️ 수신자 수만큼 실제 발송되고 과금된다.
//   - 쿼터를 넘으면 요청 시점에 **전체 거부**된다(3022) — 일부만 나가지 않는다.
//   - 개별 수신자의 실패는 건별 rejected 로 표시되고 나머지는 정상 발송된다.
//   - 수신거부 번호는 skipped 이며 **과금되지 않고 발송 기록도 만들지 않는다**.
//
// 응답: { count:, requested:, skipped:, rejected:, receipts: [...] }
func (m *AlimtalkSendModule) Bulk(params AlimtalkSendBulkParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("alimtalk/send/bulk", params, alimtalkHeaders())
}

// Cancel cancels a reserved send
// DELETE /v1/alimtalk/send/{receipt_id}
// 접수(READY) 상태의 예약 건만 취소할 수 있다 — 이미 전송에 들어갔으면 3023 이다.
func (m *AlimtalkSendModule) Cancel(receiptId string) (map[string]interface{}, error) {
	return m.api.deleteWithHeaders(fmt.Sprintf("alimtalk/send/%s", receiptId), nil, alimtalkHeaders())
}
