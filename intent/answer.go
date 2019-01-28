package intent

import (
	"fmt"
	"munzini/nlp"
	"munzini/protocol"
	"munzini/question"
	"strconv"
	//"math/rand"
	//"time"
)

// 구 대답 리스트
/*
var answers = []string{
	"저의 말에는 관심이 없으시네요.", "우울증을 의심해 보세요.", "슬퍼만 하기엔 인생은 너무나 짧죠.",
	 "결정 장애를 의심해 보세요.", "기분이 좋다고 해서 다른 사람도 기분이 좋을 거라는 생각이 실수를 만들죠.",
	 "조증을 의심해 보세요.",
}
*/

// 1. Get Simple Question Proceed Answer: 간단 문진 시작 여부 및 첫 질문 출력
func GetSQPAnswer(intent protocol.CEKIntent, qData question.QData) (protocol.CEKResponsePayload, int, question.QData) {
	var statusDelta int = 0
	var responseValue string
	var shouldEndSession bool = false
	var intentName = intent.Name

	switch intentName {
	case "Clova.YesIntent":
		qData = question.PrepareRep(qData) // prepare representative questions

		qData.RepMax = len(qData.QRepIdx)
		responseValue = "그럼, 이제부터 문진을 시작할게요. " + qData.RawData.QCWP[qData.QRepIdx[qData.RepIdx]][question.QUESTION] // current question
		statusDelta = 1                                                                                             // next status
	case "Clova.NoIntent":
		responseValue = "다음에 언제든지 불러주세요."
		shouldEndSession = true
	default:
		responseValue = "예 또는 아니오로 대답해주세요."
	}
	// make an answer
	responsePayload := protocol.CEKResponsePayload{
		OutputSpeech: protocol.MakeOutputSpeechList(
			protocol.Value{
				Lang:  "ko",
				Value: responseValue,
				Type:  "PlainText",
			},
		),
		ShouldEndSession: shouldEndSession,
	}

	return responsePayload, statusDelta, qData
}

// 2. Get Simple Question Score Answer: 간단 문진 질문에 대한 응답 입력 및 전체 간단 문진 질문에 대한 점수 계산
func GetSQSAnswer(intent protocol.CEKIntent, qData question.QData) (protocol.CEKResponsePayload, int, question.QData) {
	var score int = 0
	var statusDelta int = 0
	var responseValue string
	var shouldEndSession bool = false
	var intentName = intent.Name
	var slots = intent.Slots

	switch intentName {
	case "ScoreIntent":
		// slot에 있는 score 값 파싱
		if slots != nil { // slots가 nil이 아니어야
			if len(slots) != 0 { // slots 요소 개수가 0이 아니어야 함
				slotScore := nlp.ConvertInquiryScore(slots["inquiryScore"].Value)
				sc, err := strconv.Atoi(slotScore) // score 값 부여
				score = sc
				if err != nil { // feelingScore를 int형으로 변환한 값이 올바른 값이 아닐 때
					score = 0 // score 값에 문제가 있으므로 0으로 재부여
				}
			}
		}

		// score 값이 0이면 오류, 답을 재요구
		if score == 0 {
			responseValue = "다시 말씀해주세요."
		} else if score > 0 && score <= question.SCORE_MAX { // score 값이 정상적으로 부여된 경우
			qData.Answer[qData.QRepIdx[qData.RepIdx]] = score // score 값 저장
			qData.RepIdx++                                    // next question

			// 대표 질문이 끝났을 때
			if qData.RepIdx == qData.RepMax {
				qData = question.PrepareDet(qData) // 대표 질문들에 대한 컷오프 계산 후 문제가 있는 변증 관련 데이터 준비
				if len(qData.SQSProbPatternIdx) == 0 {
					responseValue = "간단 문진 결과 의심되는 문제가 없습니다. 앞으로도 쭈욱 건강하시고, 제가 그리우시면 언제든지 다시 불러주세요!"
					shouldEndSession = true
				} else {
					responseValue = "간단 문진 결과 " + string(len(qData.SQSProbPatternIdx)) + "개의 문제가 의심됩니다. 정밀 진단을 진행할까요?"
					statusDelta = 1 // next status
				}
			} else {
				responseValue = qData.RawData.QCWP[qData.QRepIdx[qData.RepIdx]][question.QUESTION] // next question
			}
		} else { // score 값이 1~5 가 아닌 경우
			responseValue = "1번에서 5번까지 다시 말씀해주세요."
		}
	default:
		responseValue = "다시 말씀해주세요."
	}

	// make an answer
	responsePayload := protocol.CEKResponsePayload{
		OutputSpeech: protocol.MakeOutputSpeechList(
			protocol.Value{
				Lang:  "ko",
				Value: responseValue,
				Type:  "PlainText",
			},
		),
		ShouldEndSession: shouldEndSession,
	}

	return responsePayload, statusDelta, qData
}

// 3. Get Detail Question Proceed Answer: 정밀 문진에 대한 진행 여부 및 첫 정밀 문진 질문 출력
func GetDQPAnswer(intent protocol.CEKIntent, qData question.QData) (protocol.CEKResponsePayload, int, question.QData) {
	var statusDelta int = 0
	var responseValue string
	var shouldEndSession bool = false
	var intentName = intent.Name
	// qData

	switch intentName {
	case "Clova.YesIntent":
		responseValue = "그럼, 이제부터 정밀 문진을 시작할게요. " + qData.RawData.QCWP[qData.QDetailIdx[qData.SQSProbPatternIdx[0]][qData.DetIdx]][question.QUESTION] // Detail Question 중 첫번째 질문을 이어서 내보낸다.
		statusDelta = 1
	case "Clova.NoIntent":
		responseValue = "검사하시느라 수고하셨어요. 다음에 또 불러주세요!"
		shouldEndSession = true
	default:
		responseValue = "예 또는 아니오로 대답해주세요."
	}
	// make an answer
	responsePayload := protocol.CEKResponsePayload{
		OutputSpeech: protocol.MakeOutputSpeechList(
			protocol.Value{
				Lang:  "ko",
				Value: responseValue,
				Type:  "PlainText",
			},
		),
		ShouldEndSession: shouldEndSession,
	}

	return responsePayload, statusDelta, qData
}

// 4. Get Detail Question Score Answer: 정밀 진단 질문에 대한 응답 점수 입력 및 최종 점수 계산 및 문진 결과 출력
func GetDQSAnswer(intent protocol.CEKIntent, qData question.QData) (protocol.CEKResponsePayload, int, question.QData) {
	var score int = 0
	var statusDelta int = 0
	var responseValue string
	var shouldEndSession bool = false
	var intentName = intent.Name
	var slots = intent.Slots

	qData.DetMax = len(qData.QDetailIdx[qData.SQSProbPatternIdx[qData.DetPat]])

	switch intentName {
	case "ScoreIntent":
		// slot에 있는 score 값 파싱
		if slots != nil { // slots가 nil이 아니어야
			if len(slots) != 0 { // slots 요소 개수가 0이 아니어야 함
				slotScore := nlp.ConvertInquiryScore(slots["inquiryScore"].Value) // map[string]CEKSlot , CEKSlot - Name, Value
				sc, err := strconv.Atoi(slotScore)                                // score 값 부여
				score = sc
				if err != nil { // feelingScore를 int형으로 변환한 값이 올바른 값이 아닐 때
					score = 0 // score 값에 문제가 있으므로 0으로 재부여
				}
			}
		}

		// score 값이 0이면 오류, 답을 재요구
		if score == 0 {
			responseValue = "다시 말씀해주세요."
		} else if score > 0 && score <= question.SCORE_MAX { // score 값이 정상적으로 부여된 경우
			qData.Answer[qData.QDetailIdx[qData.DetPat][qData.DetIdx]] = score // score 값 저장
			qData.DetIdx++                                                     // next question

			if qData.DetIdx == qData.DetMax {
				qData.DetPat++ // 다음 패턴
				qData.DetIdx = 0
				if qData.DetPat == len(qData.SQSProbPatternIdx) {
					qData = question.PrepareFin(qData)        // PrepareFin
					qData = makeFinalScoreNotification(qData) // 최종 결과에 대한 대답을 지정해준다.
					responseValue = qData.FinalScoreNotification + " 문진 결과를 다시 알드려릴까요?"
					statusDelta = 1
				} else {
					responseValue = qData.RawData.QCWP[qData.QDetailIdx[qData.DetPat][qData.DetIdx]][question.QUESTION] // next question                                                                                      // next question
				}
			} else {
				responseValue = qData.RawData.QCWP[qData.QDetailIdx[qData.DetPat][qData.DetIdx]][question.QUESTION] // next question
			}
		} else {
			responseValue = "1번부터 5번까지 다시 말씀해주세요."
		}
	default:
		responseValue = "다시 말씀해주시면 좋겠어요."
	}
	//make Answer
	responsePayload := protocol.CEKResponsePayload{
		OutputSpeech: protocol.MakeOutputSpeechList(
			protocol.Value{
				Lang:  "ko",
				Value: responseValue,
				Type:  "PlainText",
			},
		),
		ShouldEndSession: shouldEndSession,
	}
	return responsePayload, statusDelta, qData
}

// 5. Get Repeat Answer: 최종 문진 결과에 대한 다시 듣기 여부 처리
func GetRAnswer(intent protocol.CEKIntent, qData question.QData) (protocol.CEKResponsePayload, int, question.QData) {
	var statusDelta int = 0
	var responseValue string
	var shouldEndSession bool = false
	var intentName = intent.Name

	switch intentName {
	case "Clova.YesIntent":
		makeFinalScoreNotification(qData)
		responseValue = qData.FinalScoreNotification // 최종 검사 결과
	case "Clova.NoIntent":
		responseValue = "검사하느라 수고 많으셨어요. 다음에도 또 불러주세요."
		shouldEndSession = true
	default:
		responseValue = "예 또는 아니오로 대답해주세요."
	}

	// make an answer
	responsePayload := protocol.CEKResponsePayload{
		OutputSpeech: protocol.MakeOutputSpeechList(
			protocol.Value{
				Lang:  "ko",
				Value: responseValue,
				Type:  "PlainText",
			},
		),
		ShouldEndSession: shouldEndSession,
	}
	return responsePayload, statusDelta, qData
}

// 최종 문진 결과 생성
func makeFinalScoreNotification(qData question.QData) question.QData {

	//finalScoreNotification
	//qData.FinalScore[] 사용
	sqslength := len(qData.SQSProbPatternIdx)
	qData.FinalScoreNotification = "검진결과가 나왔어요. "
	for i := 0; i < sqslength; i++ {
		finalScoreString := fmt.Sprintf("%.2f", qData.FinalScore[qData.SQSProbPatternIdx[i]])
		qData.FinalScoreNotification += question.PATTERN_NAME[qData.SQSProbPatternIdx[i]] +
			"부분에 있어서의 점수는 " + finalScoreString + "점, "
	}
	qData.FinalScoreNotification += "입니다."
	// 나쁜 피가 뭉쳐있는 것(피멍, 혈액순환)
	// 담음이랑 어혈이 같이 옴.
	// 담음 : 몸속의 노폐물이 많음

	return qData
}
