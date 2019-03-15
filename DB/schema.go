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
}

//TODO ObjectID가 type이 깨져서 Get이 안될 시에는 ObjectID.Hex()-> String 형태로 디비에 저장할
type MedicalRecord struct {
	RecordID     string    `json:"id" bson:"_id,omitempty"`
	UserID       string    `bson:"userID"`
	TimeStamp    time.Time `bson:"timeStamp"`
	QuestionType int       `bson:"questionType"` //0 = 간단문진, 1 = 전체문진, 2 = Interrupt
	Pattern      []string  `bson:"pattern"`
	TherapyID    string    `bson:"therapyID"`
}

const (
	Username     = "partnersnco"
	Password     = "sc06250625"
	Database     = "heroku_7v6nqjgb"
	MRCollection = "MEDICALRECORD_COLLECTION"
	URCollection = "USERRECOLD_COLLECTION"
	DB_URL       = "mongodb://partnersnco:sc06250625@ds111791.mlab.com:11791/heroku_7v6nqjgb"

	COMPLECATION = "미병의심"
	PATTERN_NON  = "건강"
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
