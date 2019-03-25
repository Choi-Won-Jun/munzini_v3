//package DB
package DB

/*
* Author: Jun
* DB에 저장될 데이터들의 Schema(Document)의 Structure 정의
 */
import (
	// "fmt"
	"time"
	// "gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

type UserRecord struct {
	UserName         string                 `bson:"userName"`             //사용자 이름, 현재에는 클로바 API에서 사용자의 이름을 제공해주지 않기에 NIL or NULL로 초기화하여 사용중입니다.
	UserID           string                 `bson:"userID"`               // JSON내의 Session Attribute에서 USER ID를 파싱하여 초기화합니다.
	RecordID         []string               `bson:"recordID"`             // 해당 사용자의 문진 기록들의 고유 ID값들을 저장하는 변수
	RegistrationDate time.Time              `bson:"registrationDate"`     //회원가입한 시간, 날짜
	SimpleMRs        []Simple_MedicalRecord `bson:"simpleMedicalRecords"` // 해당 사용자의 문진 기록들중 처방, 문진타입, 추천 등의 필드들만 간략하게 저장하는 변수
}

type MedicalRecord struct {
	RecordID     string    `json:"id" bson:"_id,omitempty"` //고유 식별 ID
	UserID       string    `bson:"userID"`                  //해당 문진을 받은 사용자의 ID
	TimeStamp    time.Time `bson:"timeStamp"`               //문진을 받은 시간, 날짜
	QuestionType int       `bson:"questionType"`            //0 = 간단문진, 1 = 전체문진, 2 = Interrupt
	Pattern      []string  `bson:"pattern"`                 // 사용자에게 해당되는 병증
	CurationType int       `bson:curationType`              // 0= NONE(간단문진의 경우 추천요법이 없음) 1 = 식이요법, 2 = 운동요법 ... //TODO 추후 양생요법, 만성질환 요법 등 업데이트 필
	Curation     string    `bson:"curation"`                //추천
}

type ResultAndCuration struct {
	Pattern          string   `bson:"pattern"` // 병증 이름
	Description      string   `bson:"description"`
	Explanation      []string `bson:"explanation"`       //병증에 대한 클로바가 말할 설명이 저장된 변
	DietCuration     []string `bson:"diet_curation"`     //(식이요법)
	ExerciseCuration []string `bson:"exercise_curation"` // (운동요법)
	YangSangCuration []string `bson:"yangsang_curation"` // (양생)
	CDM_Curation     []string `bson:"CDM_curation"`      //chronic disease menagement curation (만성질환 관리)

}

// Light version of MedicalRecord which is stored in UserRecord Schema
type Simple_MedicalRecord struct {
	RecordID     string    `bson:"recordID"`     //고유 식별 ID
	TimeStamp    time.Time `bson:"timeStamp"`    //문진을 받은 시간, 날짜
	QuestionType int       `bson:"questionType"` //0 = 간단문진, 1 = 전체문진, 2 = Interrupt
	Pattern      []string  `bson:"pattern"`
	CurationType int       `bson:curationType` // 0= NONE(간단문진의 경우 추천요법이 없음) 1 = 식이요법, 2 = 운동요법 ... //TODO 추후 양생요법, 만성질환 요법 등 업데이트 필
	Curation     string    `bson:"curation"`
}

const (

	// DB 관련 변수 및 데이터가 저장될 TABLE의 이름 선언
	Username      = "partnersnco"
	Password      = "sc06250625"
	Database      = "heroku_7v6nqjgb"
	MRCollection  = "MEDICALRECORD_COLLECTION"                                                 //문진기록 Table Name
	URCollection  = "USERRECOLD_COLLECTION"                                                    //사용자 Table Name
	RnCCollection = "RESULT_AND_CURATION_COLLECTION"                                           // 해설 및 처방 Table Name
	DB_URL        = "mongodb://partnersnco:sc06250625@ds111791.mlab.com:11791/heroku_7v6nqjgb" // Heroku 서버에 업로드하여 구동하는 경우 DB_URL, Username, Password 불필요(‘uri := os.Getenv("MONGODB_URI")')와 같이 간단히 접근 가

	COMPLECATION       = "복합적미병" //3 가지 이상의 병증이 혼합된 상태
	COMPLECATION_INDEX = 5       // 0~4까지는 기존의 담음 식적 ... 등의 5 가지 병증
	PATTERN_NON        = "건강"    // 0~5까지의 병증중 해당되는 것이 없는 상
	PATTERN_NON_INDEX  = 6

	NUM_MR_to_CHECK = 3 //최근 기록 조회하여 사용자의 건강 추세를 분석할 때 기준이 되는 분석할 최근 문진 기록의 수

	RAC_SQS_EXPLANATION_INDEX = 2 //ResultAndCuration 내의 Explanation Field에 간단 문진 진단 결과가 저장된 인덱스
	RAC_DQS_EXPLANATION_INDEX = 3 //상세 문진 진단 결과가 저장된 인덱스

	SIMPLE_QUESTION_TYPE = 0 //간단 문진
	DETAIL_QUESTION_TYPE = 1 //상세문진

	//DB에 저장되는 추천 건강 요법들 관련 Index
	CURATION_NON_INDEX      = 0
	DIET_CURATION_INDEX     = 1
	Exercise_Curation_INDEX = 2
	YangSang_Curation_INDEX = 3
	CDM_Curation_INDEX      = 4
	//CURATION_MAP = map[int]
	//var PATTERN_INDEX = map[string]int{"칠정": 0, "노권": 1, "담음": 2, "식적": 3, "어혈": 4} // 변증 인덱스 : 이름
)
