package DB

import (
	"fmt"

	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//TODO 동일한 ID값을 가진 유저의 계정에 Medical Record의 key값을 추가하고, Medicalrecord collection에 해당 mr 추가
func InsertMedicalRecord(mr MedicalRecord) {

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer session.Close()
	fmt.Printf("Connected to %v!\n", session.LiveServers())

	// Insert medical-record to the DB
	c := session.DB(Database).C(MRCollection)
	if err := c.Insert(mr); err != nil {
		panic(err)
	}

	//Push medical-record ID to the repective user's record
	URc := session.DB(Database).C(URCollection)
	query := bson.M{"userID": mr.UserID}
	change := bson.M{"$push": bson.M{"recordID": mr.RecordID}}
	updateErr := URc.Update(query, change)

	if updateErr != nil {
		panic(updateErr)
	}

	// // Find Example

	// c := session.DB(Database).C(URCollection)

	// result := c.Find(bson.M{"userID": "123"})

	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Results All: ", result)

}

func InsertUserRecord(ur UserRecord) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer session.Close()
	fmt.Printf("Connected to %v!\n", session.LiveServers())

	c := session.DB(Database).C(URCollection)

	// Insert
	if err := c.Insert(ur); err != nil {
		panic(err)
	}
}

func sample_main() {

	recordID := bson.NewObjectId()

	temp := MedicalRecord{

		RecordID:     recordID,
		UserID:       "123",
		TimeStamp:    time.Now(),
		QuestionType: 1,
		Pattern:      []string{"담읍", "심혈"},
		TherapyID:    "123",
	}
	InsertMedicalRecord(temp)

	//TODO UserRecord Insert Sample
	// temp_user := UserRecord{
	// 	UserID:           "125",
	// 	UserName:         "125",
	// 	RecordID:         []string{"obj23412", "129dhflb"},
	// 	RegistrationDate: time.Now(),
	// }
	// insertUserRecord(temp_user)
}
