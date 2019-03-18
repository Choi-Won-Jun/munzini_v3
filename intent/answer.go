package intent

import (
	"fmt"       // 디버그 관련
	"math/rand" // 임의 추출 관련

	"munzini/DB"
	"munzini/nlp"      // 맞장구 관련
	"munzini/protocol" // CEK 관련 구조체
	"munzini/question" // 문진 데이터 관련
	"strconv"          // 문자열 함수 관련
	"strings"
	"time" // 임의 추출 관련
)

// 구 대답 리스트
/*
var answers = []string{
	"저의 말에는 관심이 없으시네요.", "우울증을 의심해 보세요.", "슬퍼만 하기엔 인생은 너무나 짧죠.",
	 "결정 장애를 의심해 보세요.", "기분이 좋다고 해서 다른 사람도 기분이 좋을 거라는 생각이 실수를 만들죠.",
	 "조증을 의심해 보세요.",
}
*/

const SIMPLE_QUESTION_TYPE = 0
const DETAIL_QUESTION_TYPE = 1

//TODO Therapy ID Update
const therapyID = "will be updated later"

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
		responseValue = "그럼, 이제부터 문진을 시작할게요. 질문을 듣고 긍정 혹은 부정의 뜻으로 말씀해주시면 됩니다. 첫 질문입니다. " + question.RAW_DATA.QCWP[qData.QRepIdx[qData.RepIdx]][question.QUESTION] // current question
		statusDelta = 1                                                                                                                                           // next status
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
func GetSQSAnswer(intent protocol.CEKIntent, qData question.QData, userID string) (protocol.CEKResponsePayload, int, question.QData) {
	//var score int = 0
	var statusDelta int = 0
	var responseValue string
	var shouldEndSession bool = false
	var intentName = intent.Name
	// var slots = intent.Slots
	var qNum int = 0
	var playUptoMessage string

	switch intentName {

	case "Clova.YesIntent": // 질문에 대해 문제가 있다고 이야기 할 때,
		qData.Answer[qData.QRepIdx[qData.RepIdx]] = question.YES_SCORE // 점수 부여
		fmt.Println(qData.QRepIdx[qData.RepIdx])
		playUptoMessage = nlp.GetPlayUptoMessage(question.YES_SCORE, qData.QRepIdx[qData.RepIdx])
		qData.RepIdx++ //질문 index
		// 대표 질문이 끝났을 때
		if qData.RepIdx == qData.RepMax {
			qData = question.PrepareDet(qData) // 대표 질문들에 대한 컷오프 계산 후 문제가 있는 변증 관련 데이터 준비
			for i := 0; i < len(qData.SQSProbPatternIdx); i++ {
				qNum += len(qData.QDetailIdx[qData.SQSProbPatternIdx[i]]) // 질문의 개수 이야기 해주기 위함. 모든 정밀 진단 질문 개수.
			}
			qData.QDetailNum = qNum // 정밀 진단 질문 개수 기록

			if qData.SQSProb == true {
				var SQSResult string = makeSQSResult(qData, userID) // 간단 문지 결과.
				SQSResult += " 총 " + strconv.Itoa(qNum) + "개의 질문에 대답해 주셔야 해요."
				responseValue = SQSResult
				statusDelta = 1
			} else { // NOSQSProbPatternIdx가 채워졌을 때

				responseValue = "기쁜 소식이예요! 현재 건강 발랜스가 매우 좋습니다. 지금처럼만 유지하신다면 매일매일 건강한 하루를 보내실 수 있습니다. 하지만 자만은 금물이예요! 그래도 혹시 모르니깐 더 자세한 문진을 시작해 볼까요? 총" + strconv.Itoa(question.Q_NUM-question.SQ_NUM) + "개의 질문에 대답해 주셔야 해요."

				saveUserMedicalResult(userID, SIMPLE_QUESTION_TYPE, strings.Split(DB.PATTERN_NON, " "), therapyID)

				statusDelta = 1 // next status
			}
		} else { // 간단진단 질문을 진행할 때, 특정 지점에서 남은 질문의 개수를 알려준다.
			if qData.RepIdx == question.REP_HALF {
				responseValue = "이제 절반 남았어요! 다음 질문입니다. " + question.RAW_DATA.QCWP[qData.QRepIdx[qData.RepIdx]][question.QUESTION] // next question
			} else if qData.RepIdx == question.REP_FINAL {
				responseValue = "이제 5개의 질문이 남았어요! 다음 질문입니다. " + question.RAW_DATA.QCWP[qData.QRepIdx[qData.RepIdx]][question.QUESTION] // next question
			} else {
				rand_seed := rand.NewSource(time.Now().UnixNano())
				r := rand.New(rand_seed) // 정해진 확률로 맞장구 추가하기 위함.

				responseValue = question.RAW_DATA.QCWP[qData.QRepIdx[qData.RepIdx]][question.QUESTION] // next question
				if randomPick := r.Intn(question.PROB_PLAYUPTO); randomPick == 0 {                     // 1/PROB_PLAYUPTO 확률로 점수에 해당하는 맞장구를 추가한다.
					responseValue = playUptoMessage + responseValue // nlp.PlayUpto 이제 설계 해야한다.
				}
			}
		}
	case "Clova.NoIntent": // 질문에 대해 문제가 있다고 이야기 할 때,
		qData.Answer[qData.QRepIdx[qData.RepIdx]] = question.NO_SCORE                            // 점수 부여
		playUptoMessage = nlp.GetPlayUptoMessage(question.NO_SCORE, qData.QRepIdx[qData.RepIdx]) // 맞장구 메시지 가져오
		qData.RepIdx++
		// 대표 질문이 끝났을 때
		if qData.RepIdx == qData.RepMax {
			qData = question.PrepareDet(qData) // 대표 질문들에 대한 컷오프 계산 후 문제가 있는 변증 관련 데이터 준비
			for i := 0; i < len(qData.SQSProbPatternIdx); i++ {
				qNum += len(qData.QDetailIdx[qData.SQSProbPatternIdx[i]]) // 질문의 개수 이야기 해주기 위함. 모든 SQS 정밀 진단 질문 개수.
			}
			qData.QDetailNum = qNum // 정밀 진단 질문 개수 기록

			if qData.SQSProb == true {
				var SQSResult string = makeSQSResult(qData, userID)
				SQSResult += " 총 " + strconv.Itoa(qNum) + "개의 질문에 대답해 주셔야 해요."
				responseValue = SQSResult
				statusDelta = 1 // next status
			} else {
				responseValue = "기쁜 소식이예요! 현재 건강 발랜스가 매우 좋습니다. 지금처럼만 유지하신다면 매일매일 건강한 하루를 보내실 수 있습니다. 하지만 자만은 금물이예요! 그래도 혹시 모르니깐 더 자세한 문진을 시작해 볼까요? 총" + strconv.Itoa(question.Q_NUM-question.SQ_NUM) + "개의 질문에 대답해 주셔야 해요." // 68

				saveUserMedicalResult(userID, SIMPLE_QUESTION_TYPE, strings.Split(DB.PATTERN_NON, " "), therapyID)
			}
		} else { // 간단진단 질문을 진행할 때, 특정 지점에서 남은 질문의 개수를 알려준다.
			if qData.RepIdx == question.REP_HALF {
				responseValue = "이제 절반 남았어요! 다음 질문입니다. " + question.RAW_DATA.QCWP[qData.QRepIdx[qData.RepIdx]][question.QUESTION] // next question
			} else if qData.RepIdx == question.REP_FINAL {
				responseValue = "이제 5개의 질문이 남았어요! 다음 질문입니다. " + question.RAW_DATA.QCWP[qData.QRepIdx[qData.RepIdx]][question.QUESTION] // next question
			} else {
				rand_seed := rand.NewSource(time.Now().UnixNano())
				r := rand.New(rand_seed)                                                               // 정해진 확률로 맞장구 추가하기 위함.
				responseValue = question.RAW_DATA.QCWP[qData.QRepIdx[qData.RepIdx]][question.QUESTION] // next question
				if randomPick := r.Intn(question.PROB_PLAYUPTO); randomPick == 0 {                     // 1/PROB_PLAYUPTO 확률로 점수에 해당하는 맞장구를 추가한다.
					responseValue = playUptoMessage + responseValue
				}
			}
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

	switch intentName {
	case "Clova.YesIntent":
		if qData.SQSProb == true { // 간단문진 결과 문제가 있는데, 정밀진단을 진행한다고 한 경우
			responseValue = "그럼, 이제부터 정밀 문진을 시작할게요. 이번에는 더 정확한 문진을 위해 질문들에 해당하는 정도를 1번에서 5번까지 말해주시면 되요. 자 그럼 시작할게요! " + question.RAW_DATA.QCWP[qData.QDetailIdx[qData.SQSProbPatternIdx[0]][qData.DetIdx]][question.QUESTION] // Detail Question 중 첫번째 질문을 이어서 내보낸다.
			statusDelta = 1
		} else { //  간단문진 결과 문제가 없음에도, 정밀진단을 진행한다고 한 경우
			responseValue = "그럼, 이제부터 정밀 문진을 시작할게요. 이번에는 더 정확한 문진을 위해 질문들에 해당하는 정도를 1번에서 5번까지 말해주시면 되요. 자 그럼 시작할게요! " + question.RAW_DATA.QCWP[qData.QDetailIdx[qData.NoSQSProbPatternIdx[0]][qData.DetIdx]][question.QUESTION] // 간단문진 이후 계속 하겠다고 했을 때, 시작// + question.RAW_DATA.QCWP[qData.QDetailIdx[qData.SQSProbPatternIdx[0]][qData.DetIdx]][question.QUESTION] // Detail Question 중 첫번째 질문을 이어서 내보낸다.
			statusDelta = 1
		}
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
func GetDQSAnswer(intent protocol.CEKIntent, qData question.QData, userID string) (protocol.CEKResponsePayload, int, question.QData) {

	var responsePayload protocol.CEKResponsePayload
	var statusDelta int = 0

	if qData.SQSProb == true { // 간단문진 결과 문제 패턴(SQSProbPatternIdx)이 있는 경우
		responsePayload, statusDelta, qData = GetDQSAnswer_S(intent, qData, userID)
	} else { // 간단문진 결과 문제 패턴이 없는데, 정밀검사를 진행하는 경우.
		responsePayload, statusDelta, qData = GetDQSAnswer_NS(intent, qData, userID)
	}
	return responsePayload, statusDelta, qData
}

func GetDQSAnswer_S(intent protocol.CEKIntent, qData question.QData, userID string) (protocol.CEKResponsePayload, int, question.QData) { // SQSProbPatternIdx가 존재할 때의 질문
	var score int = 0
	var statusDelta int = 0
	var responseValue string
	var shouldEndSession bool = false
	var intentName = intent.Name
	var slots = intent.Slots
	var playUptoMessage string

	qData.DetMax = len(qData.QDetailIdx[qData.SQSProbPatternIdx[qData.DetPat]])

	/*
		for i := 0; i < len(qData.SQSProbPatternIdx); i++ {
			fmt.Println(qData.SQSProbPatternIdx[i])
		}
	*/
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
			qData.Answer[qData.QDetailIdx[qData.SQSProbPatternIdx[qData.DetPat]][qData.DetIdx]] = score // score 값 저장
			playUptoMessage = nlp.GetPlayUptoMessage(score, qData.QDetailIdx[qData.SQSProbPatternIdx[qData.DetPat]][qData.DetIdx])

			// 개발 노트)
			// 이부분에 qData.QDetailIdx[qData.SQSProbPatternIdx[qData.DetPat]][qData.DetIdx] ( 방금 던진 질문 인덱스 )
			// 를 이용해서 nlp.PlayUptoMessage 를 초기화해주어야 한다. GetSQSAnswer / GetDQSAnswer_S / GetDQSAnswer_NS
			// 에도 구현해 주어야 한다.

			qData.DetIdx++       // next question
			qData.QDetailCount++ // 전체 정밀 진단 질문 수 카운트

			if qData.DetIdx == qData.DetMax {
				qData.DetPat++ // 다음 패턴
				qData.DetIdx = 0
				if qData.DetPat == len(qData.SQSProbPatternIdx) {
					qData = question.PrepareFin(qData)                // PrepareFin
					qData = makeFinalScoreNotification(qData, userID) // 최종 결과에 대한 대답을 지정해준다.
					responseValue = qData.FinalScoreNotification + " 문진 결과를 다시 알려드릴까요?"
					statusDelta = 1
				} else { // 다음 패턴 첫질문을 준비한다.
					if qData.QDetailCount%question.DETAIL_GAP == 0 {
						responseValue = "앞으로 " + strconv.Itoa(qData.QDetailNum-qData.QDetailCount) + "개의 질문이 남았어요! 다음 질문입니다. " + question.RAW_DATA.QCWP[qData.QDetailIdx[qData.SQSProbPatternIdx[qData.DetPat]][qData.DetIdx]][question.QUESTION] // next question
					} else {
						responseValue = question.RAW_DATA.QCWP[qData.QDetailIdx[qData.SQSProbPatternIdx[qData.DetPat]][qData.DetIdx]][question.QUESTION] // next question

						rand_seed := rand.NewSource(time.Now().UnixNano())
						r := rand.New(rand_seed) // 정해진 확률로 맞장구 추가하기 위함.

						if randomPick := r.Intn(question.PROB_PLAYUPTO); randomPick == 0 { // 1/PROB_PLAYUPTO 확률로 점수에 해당하는 맞장구를 추가한다.
							responseValue = /*nlp.PlayUpto(score)*/ playUptoMessage + responseValue // next question
						}
					}
				}
			} else { // 같은 패턴 내 다음 질문을 준비한다.
				if qData.QDetailCount%question.DETAIL_GAP == 0 {
					responseValue = "앞으로 " + strconv.Itoa(qData.QDetailNum-qData.QDetailCount) + "개의 질문이 남았어요! 다음 질문입니다. " + question.RAW_DATA.QCWP[qData.QDetailIdx[qData.SQSProbPatternIdx[qData.DetPat]][qData.DetIdx]][question.QUESTION] // next question
				} else {
					responseValue = question.RAW_DATA.QCWP[qData.QDetailIdx[qData.SQSProbPatternIdx[qData.DetPat]][qData.DetIdx]][question.QUESTION] // next question

					rand_seed := rand.NewSource(time.Now().UnixNano())
					r := rand.New(rand_seed) // 정해진 확률로 맞장구 추가하기 위함.

					if randomPick := r.Intn(question.PROB_PLAYUPTO); randomPick == 0 { // 1/PROB_PLAYUPTO 확률로 점수에 해당하는 맞장구를 추가한다.
						responseValue = /*nlp.PlayUpto(score)*/ playUptoMessage + responseValue // next question
					}
				}
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

func GetDQSAnswer_NS(intent protocol.CEKIntent, qData question.QData, userID string) (protocol.CEKResponsePayload, int, question.QData) { // SQSProbPattern이 존재하지 않을 때 질문
	var score int = 0
	var statusDelta int = 0
	var responseValue string
	var shouldEndSession bool = false
	var intentName = intent.Name
	var slots = intent.Slots
	var playUptoMessage string

	qData.DetMax = len(qData.QDetailIdx[qData.NoSQSProbPatternIdx[qData.DetPat]]) // 해당하는 패턴에 대한 질문 개수

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
			qData.Answer[qData.QDetailIdx[qData.NoSQSProbPatternIdx[qData.DetPat]][qData.DetIdx]] = score                            // score 값 저장
			playUptoMessage = nlp.GetPlayUptoMessage(score, qData.QDetailIdx[qData.NoSQSProbPatternIdx[qData.DetPat]][qData.DetIdx]) // 맞장구 가져오
			qData.DetIdx++                                                                                                           // next question
			qData.QDetailCount++                                                                                                     // 전체 정밀 진단 질문 수 카운트

			if qData.DetIdx == qData.DetMax {
				qData.DetPat++ // 다음 패턴
				qData.DetIdx = 0
				if qData.DetPat == len(qData.NoSQSProbPatternIdx) {
					qData = question.PrepareFin(qData)                // PrepareFin
					qData = makeFinalScoreNotification(qData, userID) // 최종 결과에 대한 대답을 지정해준다.
					responseValue = qData.FinalScoreNotification + " 문진 결과를 다시 알려드릴까요?"
					statusDelta = 1
				} else { // 다음 패턴 첫질문을 준비한다.
					if qData.QDetailCount%question.DETAIL_GAP == 0 {
						responseValue = "앞으로 " + strconv.Itoa((question.Q_NUM-question.SQ_NUM)-qData.QDetailCount) + "개의 질문이 남았어요! 다음 질문입니다. " + question.RAW_DATA.QCWP[qData.QDetailIdx[qData.NoSQSProbPatternIdx[qData.DetPat]][qData.DetIdx]][question.QUESTION] // next question
					} else { // Q_NUM - SQ_NUM : NoSQSProbPatternIdx의 질문개수를 의미한다.
						responseValue = question.RAW_DATA.QCWP[qData.QDetailIdx[qData.NoSQSProbPatternIdx[qData.DetPat]][qData.DetIdx]][question.QUESTION] // next question

						rand_seed := rand.NewSource(time.Now().UnixNano())
						r := rand.New(rand_seed) // 정해진 확률로 맞장구 추가하기 위함.

						if randomPick := r.Intn(question.PROB_PLAYUPTO); randomPick == 0 { // 1/PROB_PLAYUPTO 확률로 점수에 해당하는 맞장구를 추가한다.
							responseValue = /*nlp.PlayUpto(score)*/ playUptoMessage + responseValue // next question
						}
					}
				}
			} else { // 같은 패턴 내 다음 질문을 준비한다.
				if qData.QDetailCount%question.DETAIL_GAP == 0 {
					responseValue = "앞으로 " + strconv.Itoa((question.Q_NUM-question.SQ_NUM)-qData.QDetailCount) + "개의 질문이 남았어요! 다음 질문입니다. " + question.RAW_DATA.QCWP[qData.QDetailIdx[qData.NoSQSProbPatternIdx[qData.DetPat]][qData.DetIdx]][question.QUESTION] // next question
				} else {
					responseValue = question.RAW_DATA.QCWP[qData.QDetailIdx[qData.NoSQSProbPatternIdx[qData.DetPat]][qData.DetIdx]][question.QUESTION] // next question

					rand_seed := rand.NewSource(time.Now().UnixNano())
					r := rand.New(rand_seed) // 정해진 확률로 맞장구 추가하기 위함.

					if randomPick := r.Intn(question.PROB_PLAYUPTO); randomPick == 0 { // 1/PROB_PLAYUPTO 확률로 점수에 해당하는 맞장구를 추가한다.
						responseValue = /*nlp.PlayUpto(score)*/ playUptoMessage + responseValue // next question
					}
				}
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
		// makeFinalScoreNotification(qData)
		responseValue = qData.FinalScoreNotification + "검진 결과를 다시 들으시겠어요?" // 최종 검사 결과
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

func makeSQSResult(qData question.QData, userID string) string { // SQSProbPattern이 NULL이 아닌 경우, 간단문진 결과 출

	var sqsResult string  // 간단 문진 결과
	var identifier string // 문제 패턴 조사
	var sortedSQS []int   // identifier 초기화에 이용

	var recentCKU_result string //Recent Check up result, 최근 문진 결과를 분석하여 sqsResult에 반영
	var isDataENOUGH bool       //Recent Check up 분석을 할 최근 문진기록이 충분한지 여부를 저장하는 변수

	if len(qData.SQSProbPatternIdx) >= question.SERIOUS_SQS { // 간단문진 결과 발생한 문제가 SERIOUS_SQS개 이상일 시

		recentCKU_result, isDataENOUGH := makeRecentCheckUPResult(userID, strings.Split(DB.COMPLECATION, " "))
		saveUserMedicalResult(userID, SIMPLE_QUESTION_TYPE, strings.Split(DB.COMPLECATION, " "), therapyID)
		if isDataENOUGH == false {
			sqsResult = "문진 결과를 알려드릴께요. 현재 건강상태는 여러 가지 원인들이 합쳐서 복잡한 문제들이 나타나고 있는 상황이예요. 몸과 마음이 많이 지쳐있고, 이로 인해 삶의 질이 많이 저하된 상태예요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
		} else {
			sqsResult = "문진 결과를 알려드릴께요. 현재 건강상태는 여러 가지 원인들이 합쳐서 복잡한 문제들이 나타나고 있는 상황이예요. 몸과 마음이 많이 지쳐있고, 이로 인해 삶의 질이 많이 저하된 상태예요." + recentCKU_result + "그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
		}

		return sqsResult
	}

	sortedSQS = append(sortedSQS, qData.SQSProbPatternIdx...) // Copy SQSProbPatternIdx

	for i := 0; i < len(sortedSQS); i++ { // 오름차순으로 정리
		var minIdx int = i
		for j := i; j < len(sortedSQS); j++ {
			if sortedSQS[minIdx] > sortedSQS[j] {
				minIdx = j
			}
		}
		temp := sortedSQS[i]
		sortedSQS[i] = sortedSQS[minIdx]
		sortedSQS[minIdx] = temp
	}

	for i := 0; i < len(sortedSQS); i++ {
		identifier += question.PATTERN_NAME[sortedSQS[i]]
		if i < len(sortedSQS)-1 { // 후에 질병들을 " "를 기준으로 Split하기 위해 추가
			identifier += " "
		}
	}
	// Identifier를 이용해 Medical Record저장 수행

	//질환이 발견되지 않은 경우

	patterns := strings.Split(identifier, " ")
	recentCKU_result, flag := makeRecentCheckUPResult(userID, patterns)
	saveUserMedicalResult(userID, SIMPLE_QUESTION_TYPE, patterns, therapyID)

	switch identifier {
	case "칠정":
		sqsResult = "문진 결과를 알려드릴께요. 현재 정신적인 스트레스로 건강상태가 좋지 않아요. 스트레스가가 지속되면 식욕이 줄고 수면의 질이 나빠질 수 있어요. 그리고, 이유 없이 불안하거나 가슴이 내려앉는 느낌이 종종 나타날 수도 있고요. 그러니까 하루 빨리 스트레스에서 벗어나야 해요! 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "노권":
		sqsResult = "문진 결과를 알려드릴께요. 계속 무리를 하셔셔 몸이 항상 피곤한 상태시군요. 말하는 것조차 싫으시죠? 이런 경우 휴식을 취해도 피로가 쉽게 회복되지 않거나 오히려 심해지는 경우가 많아요. 이럴땐 휴식을 취하는 게 가장 좋은 방법이예요! 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요? "
	case "담음":
		sqsResult = "문진 결과를 알려드릴께요. 전반적으로 대사기능이 저하되어 있으시군요. 이런 상태가 계속되면 자주 두통이 생기거나 어지러울 수 있고, 소화기 기능이 약해질 수 있어요. 특히 몸의 이곳저곳이 아프거나, 기침이 자주 난다거나 가래가 끓는 호흡기의 문제로도 나타날 수 있어요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "식적":
		sqsResult = "문진 결과를 알려드릴께요. 식사가 불규칙하시거나 식습관이 좋지 않은 것 같아요. 요즘 들어 입맛이 없거나, 소화가 잘 안되거나 복통 등이 있었나 생각해 보세요! 이런 증상이 계속되면 배가 아프고 설사를 할 수도 있어요. 타고나게 소화기가 약하신 분들이라면 증상이 더 안 좋게 나타나실 수 있어요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "어혈":
		sqsResult = "문진 결과를 알려드릴께요. 현재 몸에 혈액이 제대로 돌지 못해서 한 곳에 정체되어 있을 수 있어요. 만약 몸의 특정 부위에 통증이 있거나 특히 밤에 심하게 아플 수 있어요. 특히 여성의 경우 생리통이 심하거나 생리 주기가 불규칙한 경우가 있을 수 있어요. 이런 증상이 계속되면 통증이 심해져서 하루하루가 괴로울 수 있어요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "칠정 노권":
		sqsResult = "문진 결과를 알려드릴께요. 현재 몸과 마음이 모두 지친상태시군요. 몸과 마음이 서로를 자극해 건강이 안좋은 상태예요. 이 상태를 극복하기 위해서는 충분히 휴식을 취하고, 당신만의 시간을 가져보는 것이 반드시 필요해요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "칠정 담음":
		sqsResult = "문진 결과를 알려드릴께요. 정신적 스트레스로 매우 다양한 증상에 시달리고 계시군요. 만병의 근원이 스트레스라고 하잖아요? 스트레스 때문에 몸이 아프고 몸이 아프니 스트레스 받고... 삶의 활력이 떨어진 상태네요. 우선 소화기 계통과 수면 관련 증상들을 해결해야 해요! 그럼 조금씩 건강 상태가 좋아지실 꺼예요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "칠정 식적":
		sqsResult = "문진 결과를 알려드릴께요. 정신적인 스트레스가 심해서 소화 장애가 생긴 상태예요. 이런 상태일수록 폭식과 과식, 불규칙하게 음식을 먹으면 절대 앙데여. 한동안 소화가 잘되는 음식을 위주로 식사를 하시는 게 좋을 거 같아요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "칠정 어혈":
		sqsResult = "문진 결과를 알려드릴께요. 스트레스로 인해 몸에 혈액이 제대로 돌지 못해 컨디션이 별로인 상태시군요. 그러니 기분을 풀어줄 수 있는 간단한 운동을 권해 드려요. 운동을 하고 나면 기분이 좋아져서 스트레스도 한방에 없어질 꺼예요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "노권 담음":
		sqsResult = "문진 결과를 알려드릴께요. 지금 지친 육체가 무너지려는 신호를 보내고 있어요. 육체적인 스트레스로 소화기, 호흡기, 신경계 등에 다양한 증상들이 건강을 위협하고 있는 거예요. 온몸에 활력이 떨어져 몸에 에너지를 복구하는 속도도 많이 느려져 있어요. 이런 상태를 방치하면 기운이 다 떨어져서 다시 체력을 회복하기에 힘이 많이 들꺼 같아요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "노권 식적":
		sqsResult = "문진 결과를 알려드릴께요. 체력 소모가 심해서 소화 기능이 많이 줄어든 상태예요. 식사 후에 유난히 졸리고 피곤하다면 약을 먹는 것 보다 속을 빠르게 비워주고 적절한 양의 규칙적인 식사 습관을 유지하는 것이 도움이 될꺼예요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "노권 어혈":
		sqsResult = "문진 결과를 알려드릴께요. 지금의 건강상태는 일반적으론 드물게 나타나는 증상이예요. 피로가 심해서 평소보다 다친 곳이 잘 낫지 않을 수도 있어요. 낮에는 너무 피곤하고, 저녁에는 통증에 시달릴 수도 있어요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "담음 식적":
		sqsResult = "문진 결과를 알려드릴께요. 현재 건상상태는 특히 소화기 문제가 안좋은 상태예요. 식욕이 떨어지거나 잦은 소화불량에 두통, 속쓰림, 토가 나올 것 같은 느낌이 나타날 수 있어요. 이런 증상이 지속되면 소중한 피부도 탄력을 잃고 거칠어 질 꺼예요. 약해진 소화 기능을 회복하는 것이 가장 시급해요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "담음 어혈":
		sqsResult = "문진 결과를 알려드릴께요. 현재 건강상태는 매우 다양한 증상들이 뒤섞여 나타나고 있어요.하루하루를 더욱 지치고 힘들게 할꺼예요. 그래도 지금부터 건강관리를 한다면 빠르게 회복할 수 있어요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	case "식적 어혈":
		sqsResult = "문진 결과를 알려드릴께요. 소화불량이나 아랫배를 눌렀을 때 찌르는 통증이 있을 수 있어요. 이런 증상들은 한번에 관리하기 보다는 하나씩 줄여나가는 것이 최선이예요. 그럼, 더 자세한 건강상태 확인을 위해 추가 문진을 시작해 볼까요?"
	default:
		sqsResult = ""
	}
	return sqsResult + recentCKU_result
}

/**
* Author: Jun
* 문진 결과, 판별된 Pattern들을 DB에 저장, in form of Medical Record
* questionTYPE(0: 간단 문진, 1: 정밀 문진)
 */
func saveUserMedicalResult(userID string, questionTYPE int, patterns []string, therapyID string) {

	DB.InsertMedicalRecord(userID, questionTYPE, patterns, therapyID)
}

/**
* Author: Jun
* 최근 세 번의 건강 검진 결과를 바탕으로 건강 상태에 대한 정보를 추가 제공
*
 */
func makeRecentCheckUPResult(userID string, patterns []string) (string, bool) {
	mrTABLE, mrRecords, flag := DB.GetMedicalRecordTable(userID)
	var notification string
	if flag == false { // DB에 세번 이상의 문진기록이 저장되어있지 않는 경우
		return notification, flag //종합적인 문진 결과를 notify할 수 없음
	} else {
		if patterns[0] == DB.COMPLECATION { //현재 진행중인 문진을 통한 진단결과가 미병의심(3 가지 이상 패턴의 조합)인 경우
			year_of_Record := string(mrRecords[DB.NUM_MR_to_CHECK-1].TimeStamp.Year())
			month_of_Record := string(mrRecords[DB.NUM_MR_to_CHECK-1].TimeStamp.Month())
			day_of_Record := string(mrRecords[DB.NUM_MR_to_CHECK-1].TimeStamp.Day())

			//NUM_MR_to_CHECK는 DB에서 최신순으로 불러올 Medical Record들의 수,  mrTABLE[DB.COMPLECATION_INDEX][DB.NUM_MR_to_CHECK]의 자리에는 현재 진행된 문진의 결과가 저장되어있으므로 그 이전 기록을 조회하기 위해 -1
			if mrTABLE[DB.COMPLECATION_INDEX][DB.NUM_MR_to_CHECK-1] == 1 { // case : mrTABLE[DB.COMPLECATION_INDEX][DB.NUM_MR_to_CHECK -1] == 1 => 이전 문진에서도 미병의심 진단을 받음

				notification := year_of_Record + "년 " + month_of_Record + "월 " + day_of_Record + "일 " + "부터 지금까지 종합적인 건강수치가 좋지 못한 상태에요."
			} else {
				notification := "이전 " + year_of_Record + "년 " + month_of_Record + "월 " + day_of_Record + "일 " + "문진 결과보다 종합적인 건강 상태가 나빠졌어요. "
			}
			return notification, flag

		} else if patterns[0] == DB.PATTERN_NON { // 현재 진행중인 문진을 통한 진단결과가 건강(의심되는 패턴이 없음)인 경우
			year_of_Record := string(mrRecords[DB.NUM_MR_to_CHECK-1].TimeStamp.Year())
			month_of_Record := string(mrRecords[DB.NUM_MR_to_CHECK-1].TimeStamp.Month())
			day_of_Record := string(mrRecords[DB.NUM_MR_to_CHECK-1].TimeStamp.Day())

			if mrTABLE[DB.PATTERN_NON_INDEX][DB.NUM_MR_to_CHECK-1] == 1 {
				notification := "최근 건강 상태가 아주 훌륭하시네요!"

			} else {
				notification := "이전 " + year_of_Record + "년 " + month_of_Record + "월 " + day_of_Record + "일 문진결과와 비교했을 때, " + strings.Join(mrRecords[DB.NUM_MR_to_CHECK-1]Pattern, " ") + "이 치료되었어요!"
			}
			return notification, flag
		} else { // 현재 진행중인 문진을 통한 결과가 하나 혹은 두 가지 패턴의 조합 인 경우 ex)'칠정 노권', '칠정 담음'
			return notification, flag
		}
	}
}

// 최종 문진 결과 생성
func makeFinalScoreNotification(qData question.QData, userID string) question.QData {

	var identifier string // 문제 패턴 조사
	var probNum int = 0   // 문제 패턴 개수

	if qData.SQSProb == true { // SQSProbPatternIdx 를 기반으로 정밀검사를 했을 때
		//finalScoreNotification
		//qData.FinalScore[] 사용
		sqslength := len(qData.SQSProbPatternIdx)
		for i := 0; i < sqslength; i++ { // identifier 초기화
			if qData.FinalScore[qData.SQSProbPatternIdx[i]] > question.PROB_EXIST_SCORE {
				identifier += question.PATTERN_NAME[qData.SQSProbPatternIdx[i]]
				if i < sqslength-1 { // 후에 질병들을 " "를 기준으로 Split하기 위해 추가
					identifier += " "
				}
				probNum++
			}
		}
		if probNum >= question.SERIOUS_DQS { // 3가지 이상의 문제 패턴이 있을 시,
			qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 현재 건강상태는 여러 가지 원인들이 합쳐서 복잡한 문제들이 나타나고 있는 상황이예요. 몸과 마음이 많이 지쳐있고, 이로 인해 삶의 질이 많이 저하된 상태예요. 건강에 대해 여러가지 불편이 발생하고 있어서 혼자 해결하려고 하기 보다는 가급적 의사상담을 권해 드리고 싶어요. 무엇보다 지금은 스스로의 건강에 많은 관심을 가지고, 적극적으로 관리를 꼭 하셔야해요. 주변에 가장 실력 좋은 의사선생님을 추천해 드릴까요?"

			//Insert Generated MedicalReuslt to DB
			//TODO Therapy ID Update
			therapyID := "will be updated later"
			saveUserMedicalResult(userID, DETAIL_QUESTION_TYPE, []string{DB.COMPLECATION}, therapyID)
			return qData
		}
	} else { // NoSQSProbPatternIdx 를 기반으로 정밀검사를 했을 때
		nosqslength := len(qData.NoSQSProbPatternIdx)
		for i := 0; i < nosqslength; i++ { // identifier 초기화
			if qData.FinalScore[qData.NoSQSProbPatternIdx[i]] > question.PROB_EXIST_SCORE {
				identifier += question.PATTERN_NAME[qData.NoSQSProbPatternIdx[i]]
				if i < nosqslength-1 { // 후에 질병들을 " "를 기준으로 Split하기 위해 추가
					identifier += " "
				}
				probNum++
			}
		}
		if probNum >= question.SERIOUS_DQS { // 3가지 이상의 문제 패턴이 있을 시,
			qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 현재 건강상태는 여러 가지 원인들이 합쳐서 복잡한 문제들이 나타나고 있는 상황이예요. 몸과 마음이 많이 지쳐있고, 이로 인해 삶의 질이 많이 저하된 상태예요. 건강에 대해 여러가지 불편이 발생하고 있어서 혼자 해결하려고 하기 보다는 가급적 의사상담을 권해 드리고 싶어요. 무엇보다 지금은 스스로의 건강에 많은 관심을 가지고, 적극적으로 관리를 꼭 하셔야해요. 주변에 가장 실력 좋은 의사선생님을 추천해 드릴까요?"

			//Insert Generated MedicalReuslt to DB
			//TODO Therapy ID Update
			therapyID := "will be updated later"
			saveUserMedicalResult(userID, DETAIL_QUESTION_TYPE, []string{DB.COMPLECATION}, therapyID)
			return qData
		}
	}

	fmt.Println(strconv.Itoa(probNum) + "문제 있음.")
	fmt.Println(identifier)

	if probNum == 0 { // 정밀문진 결과 문제되는 패턴이 없을 때
		qData.FinalScoreNotification = "기쁜 소식이예요! 현재 건강 발랜스가 매우 좋습니다. 지금처럼만 유지하신다면 매일매일 건강한 하루를 보내실 수 있습니다. 하지만 자만은 금물이예요! 오늘도 화이팅 하세요!"
		return qData
	}

	switch identifier {
	case "칠정":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 현재 정신적인 스트레스로 건강상태가 좋지 않아요. 스트레스가 지속되면 식욕이 줄고 수면의 질이 나빠질 수 있어요. 그리고, 이유 없이 불안하거나 가슴이 내려앉는 느낌이 종종 나타날 수도 있고요. 그러니까 하루 빨리 스트레스에서 벗어나야 해요! 그럼, 5분만 산책을 해보면 어떨까요? "
	case "노권":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 계속 무리를 하셔셔 몸이 항상 피곤한 상태시군요. 말하는 것조차 싫으시죠? 이런 경우 휴식을 취해도 피로가 쉽게 회복되지 않거나 오히려 심해지는 경우가 많아요. 이럴땐 휴식을 취하는 게 가장 좋은 방법이예요! 이번 주말엔 집에서 푹 쉬시는건 어떨까요? 그럼 빠르게 회복하실 수 있을 꺼예요. "
	case "담음":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 전반적으로 대사기능이 저하되어 있으시군요. 이런 상태가 계속되면 자주 두통이 생기거나 어지러울 수 있고, 소화기 기능이 약해질 수 있어요. 특히 몸의 이곳저곳이 아프거나, 기침이 자주 난다거나 가래가 끓는 호흡기의 문제로도 나타날 수 있어요. 특별히 증상을 명확히 설명하기도 힘든 상태라면 많은 병의 원인이 될 수 있으니 빠른 시일 내에 병원을 한번 방문해 보는 것이 좋은 방법 이예요. 주변에 적합한 병원 예약을 해드릴까요?"
	case "식적":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 식사가 불규칙하시거나 식습관이 좋지 않은 것 같아요. 요즘 들어 입맛이 없거나, 소화가 잘 안되거나 복통 등이 있었나 생각해 보세요! 이런 증상이 계속되면 배가 아프고 설사를 할 수도 있어요. 타고나게 소화기가 약하신 분들이라면 증상이 더 안 좋게 나타나실 수 있어요. 참지 말고 가까운 병원을 방문해 보세요! 금방 좋아지실 꺼 예요. "
	case "어혈":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 현재 몸에 혈액이 제대로 돌지 못해서 한 곳에 정체되어 있을 수 있어요. 만약 몸의 특정 부위에 통증이 있거나 특히 밤에 심하게 아플 수 있어요. 특히 여성의 경우 생리통이 심하거나 생리 주기가 불규칙한 경우가 있을 수 있어요. 이런 증상이 계속되면 통증이 심해져서 하루하루가 괴로울 수 있어요. 가까운 한의원에 가서 의사선생님과 상담을 해보시는건 어떨까요?"
	case "칠정 노권":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 현재 몸과 마음이 모두 지친상태시군요. 몸과 마음이 서로를 자극해 건강이 안좋은 상태예요. 이 상태를 극복하기 위해서는 충분히 휴식을 취하고, 당신만의 시간을 가져보는 것이 반드시 필요해요. 시간적 여유가 없다면 하루에 10분씩 명상을 해보는건 어떨까요?"
	case "칠정 담음":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 정신적 스트레스로 매우 다양한 증상에 시달리고 계시군요. 만병의 근원이 스트레스라고 하잖아요? 스트레스 때문에 몸이 아프고 몸이 아프니 스트레스 받고... 삶의 활력이 떨어진 상태네요. 우선 소화기 계통과 수면 관련 증상들을 해결해야 해요! 그럼 조금씩 건강 상태가 좋아지실 꺼예요. 그 비법이 궁금하시죠? 제가 알려드릴 수 있어요."
	case "칠정 식적":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 정신적인 스트레스가 심해서 소화 장애가 생긴 상태예요. 이런 상태일수록 폭식과 과식, 불규칙하게 음식을 먹으면 절대 앙데여. 한동안 소화가 잘되는 음식을 위주로 식사를 하시는 게 좋을 거 같아요. 제가 다양한 식이요법을 알고 있어요. 알려드릴까요?"
	case "칠정 어혈":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 스트레스로 인해 몸에 혈액이 제대로 돌지 못해 컨디션이 별로인 상태시군요. 그러니 기분을 풀어줄 수 있는 간단한 운동을 권해 드려요. 운동을 하고 나면 기분이 좋아져서 스트레스도 한방에 없어질 꺼예요. 도움이 될만한 좋은 운동 영상을 추천해 드릴까요?"
	case "노권 담음":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 지금 지친 육체가 무너지려는 신호를 보내고 있어요. 육체적인 스트레스로 소화기, 호흡기, 신경계 등에 다양한 증상들이 건강을 위협하고 있는 거예요. 온몸에 활력이 떨어져 몸에 에너지를 복구하는 속도도 많이 느려져 있어요. 이런 상태를 방치하면 기운이 다 떨어져서 다시 체력을 회복하기에 힘이 많이 들꺼 같아요. 빨리 건강 밧데리를 충전해야 해요! 비법이 궁금하시면 알려드릴 수 있어요!"
	case "노권 식적":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 체력 소모가 심해서 소화 기능이 많이 줄어든 상태예요. 식사 후에 유난히 졸리고 피곤하다면 약을 먹는 것 보다 속을 빠르게 비워주고 적절한 양의 규칙적인 식사 습관을 유지하는 것이 도움이 될꺼예요. 그럼 건강한 식사 습관 방법을 추천해 볼까요?"
	case "노권 어혈":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 지금의 건강상태는 일반적으론 드물게 나타나는 증상이예요. 피로가 심해서 평소보다 다친 곳이 잘 낫지 않을 수도 있어요. 낮에는 너무 피곤하고, 저녁에는 통증에 시달릴 수도 있어요. 아무래도 병원을 한 번 방문해서 상담을 받아보는걸 권하고 싶어요. 가까운 한의원을 예약해 드릴까요?"
	case "담음 식적":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 현재 건상상태는 특히 소화기 문제가 안좋은 상태예요. 식욕이 떨어지거나 잦은 소화불량에 두통, 속쓰림, 토가 나올 것 같은 느낌이 나타날 수 있어요. 이런 증상이 지속되면 소중한 피부도 탄력을 잃고 거칠어 질 꺼예요. 약해진 소화 기능을 회복하는 것이 가장 시급해요. 그러니까 음식의 양과 종류에 특히 신경을 써야 해요! 이럴 때 딱 맞는 음식들이 있어요. 레시피를 추천해 드려볼까요?"
	case "담음 어혈":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 현재 건강상태는 매우 다양한 증상들이 뒤섞여 나타나고 있어요.하루하루를 더욱 지치고 힘들게 할꺼예요. 그래도 지금부터 건강관리를 한다면 빠르게 회복할 수 있어요. 몸에 좋은 음식과 간단한 운동을 지속적으로 한다면 간단한 증상들은 모두 사라질꺼예요. 제가 좋은 비법을 알려 드릴까요?"
	case "식적 어혈":
		qData.FinalScoreNotification = "문진 결과를 알려드릴께요. 소화불량이나 아랫배를 눌렀을 때 찌르는 통증이 있을 수 있어요. 이런 증상들은 한번에 관리하기 보다는 하나씩 줄여나가는 것이 최선이예요. 건강개선을 위한 스케쥴을 추천해 드릴까요?"
	default:
		qData.FinalScoreNotification = ""
	}

	return qData
}
