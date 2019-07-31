package handler

import (
	"encoding/json"

	"log"
	//"munzini/DB"
	"munzini/intent"
	"munzini/protocol"

	"net/http"
)

const SQP_S = 0 // Simple Question Proceed Status
const SQS_S = 1 // Simple Question Score Status
const DQP_S = 2 // Detail Question Proceed Status
const DQS_S = 3 // Detail Question Score Status
const R_S = 4   // Repeat Status

// ServeHTTP handles CEK requests
func Dispatch(w http.ResponseWriter, r *http.Request) {

	var req protocol.CEKRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("JSON decoding failed")
		respondError(w, "서버와의 연결이 원활하지 않네요")
		return
	}

	reqType := req.Request.Type

	var response protocol.CEKResponse
	var result protocol.CEKResponsePayload
	var statusDelta int

	// 요청 타입에 따른 기능 수행
	switch reqType {
	case "LaunchRequest": // 앱 실행 요청 시
		response = protocol.MakeCEKResponse(handleLaunchRequest())

		var sessionAttributesRes protocol.CEKSessionAttributes
		sessionAttributesRes.Status = 0
		response = protocol.SetSessionAttributes(response, sessionAttributesRes)

	case "SessionEndedRequest": // 앱 종료 요청 시
		response = protocol.MakeCEKResponse(handleEndRequest())

	case "IntentRequest": // 의도가 담긴 요청 시

		sesstionAttributesReq := req.Session.SessionAttributes
		status := sesstionAttributesReq.Status
		qdata := sesstionAttributesReq.QData
		// Author : Wonjun
		fqcore := sesstionAttributesReq.FQCore

		//userID := sesstionAttributesReq.

		cekIntent := req.Request.Intent // CEKIntent

		// 사용자의 발화에 대한 응답을 현재 상태에 따라 세팅한다. 필요한 경우 응답을 세팅하는 과정에서 슬롯에 대한 처리를 포함한다.
		switch status {
		case SQP_S: // status가 0인 경우
			result, statusDelta, qdata, fqcore = intent.GetSQPAnswer(cekIntent, qdata, fqcore, req.Session.User.UserId)
		case SQS_S:
			result, statusDelta, qdata, fqcore = intent.GetSQSAnswer(cekIntent, qdata, fqcore, req.Session.User.UserId)
		case DQP_S:
			result, statusDelta, qdata, fqcore = intent.GetDQPAnswer(cekIntent, qdata, fqcore)
		case DQS_S:
			result, statusDelta, qdata, fqcore = intent.GetDQSAnswer(cekIntent, qdata, fqcore, req.Session.User.UserId) // 개발노트) qData.SQSProb에 따라 다르게 처리 하도록 구현해야 함.
		case R_S:
			result, statusDelta, qdata, fqcore = intent.GetRAnswer(cekIntent, qdata, fqcore)
		}
		response = protocol.MakeCEKResponse(result) // 응답 구조체 작성
		status += statusDelta                       // 상태 변화 적용

		var sessionAttributesRes protocol.CEKSessionAttributes
		sessionAttributesRes.Status = status
		sessionAttributesRes.QData = qdata
		sessionAttributesRes.FQcore = fqcore
		response = protocol.SetSessionAttributes(response, sessionAttributesRes) // json:status 값 추가
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(&response)
	w.Write(b)
}

func handleLaunchRequest() protocol.CEKResponsePayload {
	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeSimpleOutputSpeech("안녕하세요, 히포입니다.? 오늘의 문진을 시작해볼까요?"),
		ShouldEndSession: false,
	}
}

func handleEndRequest() protocol.CEKResponsePayload {
	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeSimpleOutputSpeech("편하실 때 다시 불러주세요."),
		ShouldEndSession: true,
	}
}

func respondError(w http.ResponseWriter, msg string) {
	response := protocol.MakeCEKResponse(
		protocol.CEKResponsePayload{
			OutputSpeech: protocol.MakeSimpleOutputSpeech(msg),
		})

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(&response)
	w.Write(b)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {}
