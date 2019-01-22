package intent

import (
	"munzini/protocol"
	"munzini/question"
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

var qData question.QData

var repIdx int = 0
var repMax int
var detIdx int = 0
var detMax int

func GetSQPAnswer(intentName string) (protocol.CEKResponsePayload, int) {
	var statusDelta int = 0
	var responseValue string
	var shoudEndSession bool
	
	switch intentName {
	case "YesIntent":
		qData = question.PrepareRep(qData)	// prepare representative questions
		responseValue = qData.RawData.QCWP[qData.QRepIdx[repIdx++]][question.QUESTION]	// next question
		shouldEndSession = false
		statusDelta = 1
	case "NoIntent":
		responseValue = "다음에 언제든지 불러주세요"
		shouldEndSession = true
	case default:
		responseValue = "예 또는 아니오로 대답해주세요."
		shouldEndSession = false
	}
	
	// make an answer
	responsePayload := protocol.CEKResponsePayload{
		OutputSpeech: protocol.MakeOutputSpeechList(
			protocol.Value{
				Lang:  "ko",
				Value: responseValue
				Type:  "PlainText",
			},
		),
		ShouldEndSession: shoudEndSession,
	}
	
	return responsePayload, statusDelta
}

func GetSQSAnswer(intentName string, slots protocol.CEKRequest.Request.Intent.Slots) (protocol.CEKResponsePayload, int) {

}

func GetDQPAnswer(intentName string) (protocol.CEKResponsePayload, int) {

}

func GetDQSAnswer(intentName string, slots protocol.CEKRequest.Request.Intent.Slots) (protocol.CEKResponsePayload, int) {

}

func GetRAnswer(intentName string) (protocol.CEKResponsePayload, int) {

}
