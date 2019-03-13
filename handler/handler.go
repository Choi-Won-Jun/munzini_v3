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

		cekIntent := req.Request.Intent // CEKIntent

		// 사용자의 발화에 대한 응답을 현재 상태에 따라 세팅한다. 필요한 경우 응답을 세팅하는 과정에서 슬롯에 대한 처리를 포함한다.
		switch status {
		case SQP_S:
			result, statusDelta, qdata = intent.GetSQPAnswer(cekIntent, qdata)
		case SQS_S:
			result, statusDelta, qdata = intent.GetSQSAnswer(cekIntent, qdata)
		case DQP_S:
			result, statusDelta, qdata = intent.GetDQPAnswer(cekIntent, qdata)
		case DQS_S:
			result, statusDelta, qdata = intent.GetDQSAnswer(cekIntent, qdata) // 개발노트) qData.SQSProb에 따라 다르게 처리 하도록 구현해야 함.
		case R_S:
			result, statusDelta, qdata = intent.GetRAnswer(cekIntent, qdata)
		}
		response = protocol.MakeCEKResponse(result) // 응답 구조체 작성
		status += statusDelta                       // 상태 변화 적용

		var sessionAttributesRes protocol.CEKSessionAttributes
		sessionAttributesRes.Status = status
		sessionAttributesRes.QData = qdata
		response = protocol.SetSessionAttributes(response, sessionAttributesRes) // json:status 값 추가
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(&response)
	w.Write(b)
}

func handleLaunchRequest() protocol.CEKResponsePayload {

	// uri := os.Getenv("MONGODB_URI")
	// if uri == "" {
	// 	fmt.Println("no connection string provided")
	// 	os.Exit(1)
	// }
	// session, err := mgo.Dial(uri)
	// if err != nil {
	// 	fmt.Printf("Can't connect to mongo, go error %v\n", err)
	// 	os.Exit(1)
	// }
	// defer session.Close()

	// // // Insert
	// c := session.DB(DB.Database).C(DB.MRCollection)
	// recordID := bson.NewObjectId()

	// temp := DB.MedicalRecord{

	// 	RecordID:     recordID,
	// 	UserID:       "125",
	// 	TimeStamp:    time.Now(),
	// 	QuestionType: 1,
	// 	Pattern:      []string{"담읍", "심혈"},
	// 	TherapyID:    "125",
	// }

	// // Insert
	// if err := c.Insert(temp); err != nil {
	// 	panic(err)
	// }

	return protocol.CEKResponsePayload{
		OutputSpeech:     protocol.MakeSimpleOutputSpeech("안녕하세요, 문지니입니다.? 오늘의 문진을 시작해볼까요?"),
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
