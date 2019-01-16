package intent

import (
	"magicball/protocol"
	//"math/rand"
	//"time"
)

var answers = []string{
	"저의 말에는 관심이 없으시네요.", "우울증을 의심해 보세요.", "슬퍼만 하기엔 인생은 너무나 짧죠.",
	 "결정 장애를 의심해 보세요.", "기분이 좋다고 해서 다른 사람도 기분이 좋을 거라는 생각이 실수를 만들죠.",
	 "조증을 의심해 보세요.",
}

func GetAnswer(opt int) (protocol.CEKResponsePayload, error) {
	//rand.Seed(time.Now().UTC().UnixNano())
	//randomIndex := rand.Intn(len(answers))

	if opt < 0 || opt > 5 {	// 값이 0 미만 또는 5 초과일 때 예외처리
		responsePayload := protocol.CEKResponsePayload{
			OutputSpeech: protocol.MakeOutputSpeechList(
				protocol.Value{
					Lang:  "ko",
					Value: "제가 말한 범위를 벗어났어요.",
					Type:  "PlainText",
				},
			),
			ShouldEndSession: false,
		}

		return responsePayload, nil
	}

	responsePayload := protocol.CEKResponsePayload{
		OutputSpeech: protocol.MakeOutputSpeechList(
			protocol.Value{
				Lang:  "",
				//Value: "https://ssl.pstatic.net/static/clova/service/native_extensions/magicball/magic_ball_sound.mp3",
				Value: "resources/magicball.magic_ball_sound.mp3",
				Type:  "URL",
			},
			protocol.Value{
				Lang:  "ko",
				Value: answers[opt],
				Type:  "PlainText",
			},
		),
		ShouldEndSession: false,
	}

	return responsePayload, nil
}
