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
	UserName         string    `bson:"userName"`
	UserID           string    `bson:"userID"`
	RecordID         []string  `bson:"recordID"`
	RegistrationDate time.Time `bson:"registrationDate"`
	SimpleMRs []Simple_MedicalRecord 'bson:"simple_Medical_Records"'
}

type MedicalRecord struct {
	RecordID     string    `json:"id" bson:"_id,omitempty"`
	UserID       string    `bson:"userID"`
	TimeStamp    time.Time `bson:"timeStamp"`
	QuestionType int       `bson:"questionType"` //0 = 간단문진, 1 = 전체문진, 2 = Interrupt
	Pattern      []string  `bson:"pattern"`
	CurationType int       `bson:curationType` // 0= NONE(간단문진의 경우 추천요법이 없음) 1 = 식이요법, 2 = 운동요법 ...
	Curation     string    `bson:"curation"`
}

type ResultAndCuration struct {
	Pattern      string   `bson:"pattern"`
	Description  string   `bson:"description"`
	Explanation  []string `bson:"explanation"`
	DietCuration []string `bson:"diet_curation"`
	// 운동요법 등등 다른 처방 필드 아래에 추가
}

// Light version of MedicalRecord which is stored in UserRecord Schema
type Simple_MedicalRecord struct {
	TimeStamp    time.Time `bson:"timeStamp"`
	QuestionType int       `bson:"questionType"` //0 = 간단문진, 1 = 전체문진, 2 = Interrupt
	Pattern      []string  `bson:"pattern"`
	CurationType int       `bson:curationType` // 0= NONE(간단문진의 경우 추천요법이 없음) 1 = 식이요법, 2 = 운동요법 ...
	Curation     string    `bson:"curation"`
}

const (
	Username      = "partnersnco"
	Password      = "sc06250625"
	Database      = "heroku_7v6nqjgb"
	MRCollection  = "MEDICALRECORD_COLLECTION"
	URCollection  = "USERRECOLD_COLLECTION"
	RnCCollection = "RESULT_AND_CURATION_COLLECTION"
	DB_URL        = "mongodb://partnersnco:sc06250625@ds111791.mlab.com:11791/heroku_7v6nqjgb"

	COMPLECATION       = "복합적미병"
	COMPLECATION_INDEX = 5
	PATTERN_NON        = "건강"
	PATTERN_NON_INDEX  = 6

	NUM_MR_to_CHECK = 3

	RAC_SQS_ExPLANATION_INDEX = 2

	SIMPLE_QUESTION_TYPE = 0

	//DB에 저장되는 추천 건강 요법들 관련 Index
	CURATION_NON_INDEX  = 0
	DIET_CURATION_INDEX = 1
)

// func main() {
// 	fmt.Print("1")
// 	Host := []string{
// 		"127.0.0.1:27017",
// 		// replica set addrs...
// 	}
// 	session, err := mgo.DialWithInfo(&mgo.DialInfo{
// 		Addrs: Host,
// 		// Username: Username,
// 		// Password: Password,
// 		// Database: Database,
// 		// DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
// 		// 	return tls.Dial("tcp", addr.String(), &tls.Config{})
// 		// },
// 	})

// 	fmt.Print("1")

// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()
// 	fmt.Printf("Connected to %v!\n", session.LiveServers())

// 	temp := MedicalRecord{

// 		RecordID:     "123",
// 		TimeStamp:    time.Now().String(),
// 		QuestionType: 1,
// 		Pattern:      []string{"담읍", "심혈"},
// 		TherapyID:    "123",
// 	}

// 	// Collection
// 	c := session.DB(Database).C(Collection)

// 	// Insert
// 	if err := c.Insert(temp); err != nil {
// 		panic(err)
// 	}

// }
