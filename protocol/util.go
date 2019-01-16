package protocol

// MakeCEKResponse creates CEKResponse instance with given params
func MakeCEKResponse(responsePayload CEKResponsePayload) CEKResponse {
	response := CEKResponse{
		Version:  "0.1.0",
		Response: responsePayload,
	}

	return response
}

// 멀티턴을 위한 리스폰스 세팅
func SetMultiturn(response CEKResponse, sessionAtt SessionAttributes) CEKResponse {
  response.Response.ShouldEndSession = false
	response.SessionAttributes = sessionAtt

	return response
}

// MakeOutputSpeech creates OutputSpeech instance with given params
func MakeSimpleOutputSpeech(msg string) OutputSpeech {
	return OutputSpeech{
		Type: "SimpleSpeech",
		Values: Value{
			Lang:  "ko",
			Value: msg,
			Type:  "PlainText",
		},
	}
}

// MakeOutputSpeech creates OutputSpeech instance with given params
func MakeOutputSpeechList(value ...Value) OutputSpeech {
	return OutputSpeech{
		Type:   "SpeechList",
		Values: value,
	}
}
