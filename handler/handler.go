package handler

import (
	"encoding/json"
	"log"
	"munzini/intent"
	"munzini/protocol"
	"net/http"
)

const SQP_S = 0 // Simple Question Proceed Status
const SQS_S = 1 // Simple Question Score Status
const DQP_S = 2 // Detail Question Proceed Status
const DQS_S = 3 // Detail Question Score Status
const R_S = 4   // Repeat Status

var status int = SQP_S // initial status is SQP_S

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

	// 요청 타입에 따른 기능 수행
	switch reqType {
	case "LaunchRequest": // 앱 실행 요청 시
		response = protocol.MakeCEKResponse(handleLaunchRequest())
	case "SessionEndedRequest": // 앱 종료 요청 시
		response = protocol.MakeCEKResponse(handleEndRequest())
	case "IntentRequest": // 의도가 담긴 요청 시
		intentName := req.Request.Intent.Name // 의도의 이름
		slots := req.Request.Intent.Slots     // 슬롯

		// 사용자의 발화에 대한 응답을 현재 상태에 따라 세팅한다. 필요한 경우 응답을 세팅하는 과정에서 슬롯에 대한 처리를 포함한다.
		switch status {
		case SQP_S:
			result, statusDelta := intent.GetSQPAnswer(intentName)
		case SQS_S:
			result, statusDelta := intent.GetSQSAnswer(intentName, slots)
		case DQP_S:
			result, statusDelta := intent.GetDQPAnswer(intentName)
		case DQS_S:
			result, statusDelta := intent.GetDQSAnswer(intentName, slots)
		case R_S:
			result, statusDelta := intent.GetRAnswer(intentName)
		}
		response = protocol.MakeCEKResponse(result) // 응답 구조체 작성
		status += statusDelta                       // 상태 변화 적용

		// 슬롯 파싱 코드 (참조용)
		/*
			switch intentName {
			case "FeelingIntent":
				var feelingScore int = 0
				var err error
				if (slots != nil) {	// slots가 nil인 경우 if문 건너뜀
					if (len(slots) != 0) {	// slots 요소 개수가 0이 아니어야 함
						feelingScoreSlotValue := slots["feelingScore"].Value
						feelingScore, err = strconv.Atoi(feelingScoreSlotValue)
						if (err != nil) {	// feelingScore를 int형으로 변환한 값이 올바른 값이 아닐 때
							feelingScore = 0
						}
					}
				}
				opt := feelingScore
				if result, err := intent.GetAnswer(opt); err == nil {
					response = protocol.MakeCEKResponse(result)
				}
				break
			case "Clova.GuideIntent":
			default:
				response = protocol.MakeCEKResponse(handleWrongAnnounce())
			}
		*/
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(&response)
	w.Write(b)
}

func handleLaunchRequest() protocol.CEKResponsePayload {
	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeSimpleOutputSpeech("부르셨나요? 오늘의 문진을 시작해볼까요?"),
		ShouldEndSession: false,
	}
}

func handleEndRequest() protocol.CEKResponsePayload {
	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeSimpleOutputSpeech("편하실 때 다시 불러주세요."),
		ShouldEndSession: true,
	}
}

func handleWrongAnnounce() protocol.CEKResponsePayload {
	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeSimpleOutputSpeech("다시 말씀해주세요."),
		ShouldEndSession: false,
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
