package protocol

import (
	"munzinis_project/question"
)

// 요청 json.Request.Intent
type CEKIntent struct {
	IntentType string             `json:"intent"` //RULE 'json"intent' would be deprecated(keep this for compatibility for a while). instead of, use 'json:"name"'.
	Name       string             `json:"name"`   //RULE should be set same as IntentType value.
	Slots      map[string]CEKSlot `json:"slots"`
}

// 요청 json.Reqeust.Intent.Slots 리스트 요소
type CEKSlot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// 요청 json.Request의 한 요소
type RequestCommon struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	Locale    string `json:"locale"`
}

// 요청 json.Session
type CEKSession struct {
	New       bool   `json:"new"`
	SessionId string `json:"sessionId`
	User      struct {
		AccessToken string `json:"accessToken"`
		UserId      string `json:"userId`
	}
	SessionAttributes CEKSessionAttributes `json:"sessionAttributes"`
}

// 요청 json.Request
type CEKRequestPayload struct {
	RequestCommon
	Intent CEKIntent   `json:"intent,omitempty"`
	Event  interface{} `json:"event,omitempty"`
}

// 요청 json
type CEKRequest struct {
	Version  string                 `json:"version"`
	Session  CEKSession             `json:"session"`
	Contexts map[string]interface{} `json:"context"`
	Request  CEKRequestPayload      `json:"request"`
}

// 응답 json.Response.OutputSpeech.Values 리스트 요소
type Value struct {
	Lang  string `json:"lang"`
	Type  string `json:"type"` // Will be deprecated
	Value string `json:"value"`
}

// 응답 json.Response.OutputSpeech
type OutputSpeech struct {
	Type   string      `json:"type"`
	Values interface{} `json:"values"`
}

// 응답 json.Response
type CEKResponsePayload struct {
	OutputSpeech     OutputSpeech `json:"outputSpeech"`
	Card             interface{}  `json:"card,omitempty"`
	Directives       interface{}  `json:"directives"`
	ShouldEndSession bool         `json:"shouldEndSession"`
}

// 응답 json
type CEKResponse struct {
	Version           string             `json:"version"`
	SessionAttributes interface{}        `json:"sessionAttributes"`
	Response          CEKResponsePayload `json:"response"`
}

type Card struct {
	ActionList    []CardValue `json:"actionList"`
	BgUrl         CardValue   `json:"bgUrl`
	HighlightText CardValue   `json:"highlightText"`
	MainText      CardValue   `json:"mainText"`
	ParagraphText CardValue   `json:"paragraphText"`
	ReferenceText CardValue   `json:"referenceText"`
	ReferenceUrl  CardValue   `json:"referenceUrl"`
	SentenceText  CardValue   `json:"sentenceText"`
	SubText       CardValue   `json:"subText"`
	TableList     []CardValue `json:"tableList"`
	emotionCode   CardValue   `json:"emotionCode"`
	motionCode    CardValue   `json:"motionCode"`
	Type          string      `json:"type"`
}

type CardValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// SesstionAttributes 값
type CEKSessionAttributes struct {
	Status int            `json:"status"`
	QData  question.QData `json:"qdata"`
}
