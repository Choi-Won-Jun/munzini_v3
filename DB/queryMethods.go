package DB

import (
	//"fmt"

	//"log"
	//"munzini/question"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//TODO 동일한 ID값을 가진 유저의 계정에 Medical Record의 key값을 추가하고, Medicalrecord collection에 해당 mr 추가
/*
* Author: Jun
* 동일한 ID값을 가진 유저의 계정에 Medical Record의 key값을 추가하고, Medicalrecord collection에 해당 mr 추가
 */
func InsertMedicalRecord(userID string, questionTYPE int, patterns []string, therapyID string) {

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		//	fmt.Println("no connection string provided")
		os.Exit(1)
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		//	fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer session.Close()
	//fmt.Printf("Connected to %v!\n", session.LiveServers())

	recordID := bson.NewObjectId().Hex()
	medicalRecord := MedicalRecord{
		RecordID:     recordID,
		UserID:       userID,
		TimeStamp:    time.Now(),
		QuestionType: questionTYPE,
		Pattern:      patterns,
		TherapyID:    therapyID,
	}

	// Insert medical-record to the DB
	insertC := session.DB(Database).C(MRCollection)
	if insertErr := insertC.Insert(medicalRecord); insertErr != nil {
		panic(insertErr)
	}

	// Find First, If user is not exist in database, add his data
	findC := session.DB(Database).C(URCollection)

	var result []UserRecord
	if findErr := findC.Find(bson.M{"userID": userID}).All(&result); findErr != nil {
		panic(findErr)
	}

	if len(result) == 0 {
		temp_user := UserRecord{
			UserID:           userID,
			UserName:         "nil",
			RecordID:         []string{},
			RegistrationDate: time.Now(),
		}
		InsertUserRecord(temp_user)
	}

	//Push medical-record ID to the repective user's record
	updateC := session.DB(Database).C(URCollection)
	query := bson.M{"userID": medicalRecord.UserID}
	change := bson.M{"$push": bson.M{"recordID": medicalRecord.RecordID}}
	updateErr := updateC.Update(query, change)

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

/*
* Author: Jun
* 사용자 정보를 DB안의 UserRecordCollection에 추가
 */
func InsertUserRecord(ur UserRecord) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		//fmt.Println("no connection string provided")
		os.Exit(1)
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		//fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer session.Close()
	//fmt.Printf("Connected to %v!\n", session.LiveServers())

	c := session.DB(Database).C(URCollection)

	// Insert
	if err := c.Insert(ur); err != nil {
		panic(err)
	}
}

/*
* Author: Jun
* Look up the recent medical records by userID
 */
func RetreiveRecentMedicalRecordByUserID(userID string) []MedicalRecord {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		//fmt.Println("no connection string provided")
		os.Exit(1)
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		//fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer session.Close()
	//fmt.Printf("Connected to %v!\n", session.LiveServers())

	// Find First, If user is not exist in database, add his data
	findC := session.DB(Database).C(URCollection)

	var urRecord UserRecord
	iter := findC.Find(bson.M{"userID": userID}).Limit(NUM_MR_to_CHECK).Iter()
	findErr := iter.All(&urRecord)
	if findErr != nil {
		panic(err)
	}
	// if findErr := findC.Find(bson.M{"userID": userID}).One(&urRecord); findErr != nil {
	// 	panic(findErr)
	// }

	findMR := session.DB(Database).C(MRCollection)

	// List of IDs of Medical Records
	mrIDs := urRecord.RecordID

	medicalRecords := []MedicalRecord{}

	for _, mrID := range mrIDs {
		var tempMR MedicalRecord

		if FindMRError := findMR.Find(bson.M{"_id": mrID}).One(&tempMR); FindMRError != nil {
			panic(FindMRError)
		}
		medicalRecords = append(medicalRecords, tempMR)

	}

	return medicalRecords

}

// func getMedicalRecordTable(userID string) ([][]int, bool) {
// 	//var PATTERN_NAME = []string{"칠정", "노권", "담음", "식적", "어혈"}                       // 변증 이름
// 	//var PATTERN_INDEX = map[string]int{"칠정": 0, "노권": 1, "담음": 2, "식적": 3, "어혈": 4} // 변증 인덱스 : 이름
// 	//const PATTERN_NUM = 5

// 	// 확인할 최근 문진 기록들의 수
// 	var num_Records = 3

// 	if len(medicalRecords) < 3{ // 저장된 문진 기록들이 충분하지 않은 경우

// 		flag := 0
// 		return nil, flag
// 	} else {
// 		medicalRecords := RetreiveRecentMedicalRecordByUserID(userID)
// 	patternRecords :=
// 	mrTable = [question.PATTERN_NUM][num_Records]int

// 		flag := 1
// 	}

// }

/*
* Author: Jun
* DB와의 Connection을 생성 뒤 반환
 */
// func CreateSession() {

// }

// func sample_main() {

// 	recordID := bson.NewObjectId()

// 	temp := MedicalRecord{

// 		RecordID:     recordID,
// 		UserID:       "123",
// 		TimeStamp:    time.Now(),
// 		QuestionType: 1,
// 		Pattern:      []string{"담읍", "심혈"},
// 		TherapyID:    "123",
// 	}
// 	InsertMedicalRecord(temp)

// 	//TODO UserRecord Insert Sample
// 	// temp_user := UserRecord{
// 	// 	UserID:           "125",
// 	// 	UserName:         "125",
// 	// 	RecordID:         []string{"obj23412", "129dhflb"},
// 	// 	RegistrationDate: time.Now(),
// 	// }
// 	// insertUserRecord(temp_user)
// }
