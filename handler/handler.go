package handler

import (
	"encoding/json"
	"log"
	"magicball/intent"
	"magicball/protocol"
	"net/http"
	"strconv"
)

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

	switch reqType {
	case "LaunchRequest":	// 앱 실행 직후
		response = protocol.SetMultiturn(protocol.MakeCEKResponse(handleLaunchRequest()), protocol.SessionAttributes{
			Intent:	"FeelingIntent",
		})
	case "SessionEndedRequest":	// 앱 종료 시
		response = protocol.MakeCEKResponse(handleEndRequest())

	case "IntentRequest":	// 의도에 따른 기능 수행
		intentName := req.Request.Intent.Name
		slots := req.Request.Intent.Slots

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
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(&response)
	w.Write(b)
}

func handleLaunchRequest() protocol.CEKResponsePayload {
	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeSimpleOutputSpeech("여기는 비밀의 구슬입니다. 기분이 좋은 정도를 1번에서 5번까지라고 한다면, 지금 기분이 어떠세요?"),
		ShouldEndSession: false,
	}
}

func handleEndRequest() protocol.CEKResponsePayload {
	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeSimpleOutputSpeech("비밀의 구슬은 이만 퇴갤하겠습니다."),
		ShouldEndSession: true,
	}
}

func handleWrongAnnounce() protocol.CEKResponsePayload {
	return protocol.CEKResponsePayload{
		OutputSpeech:			protocol.MakeSimpleOutputSpeech("1번에서 5번까지 중에 골라서 말씀해 주세요."),
		ShouldEndSession:	false,
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
