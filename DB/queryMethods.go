package DB

import (
	//"fmt"

	"bufio"
	"encoding/csv"
	"log"
	"munzini/question"
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
func RetreiveRecentMedicalRecordByUserID(userID string) ([]MedicalRecord, bool) {
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
	// iter := findC.Find(bson.M{"userID": userID}).Limit(NUM_MR_to_CHECK).Iter()
	// findErr := iter.All(&urRecord)
	// if findErr != nil {
	// 	panic(err)
	// }
	if findErr := findC.Find(bson.M{"userID": userID}).One(&urRecord); findErr != nil {
		panic(findErr)
	}

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

	// 충분한 수의 문진기록이 저장되어있는지 확인
	if len(medicalRecords) < NUM_MR_to_CHECK {
		flag := false
		return nil, flag
	}

	flag := true
	return medicalRecords[len(medicalRecords)-NUM_MR_to_CHECK : len(medicalRecords)], flag

}

func GetMedicalRecordTable(userID string) ([question.PATTERN_NUM + 2][NUM_MR_to_CHECK]int, []MedicalRecord, bool) { //기본 5가지의 패턴과 미병의심, 건강의 2 가지 패턴을 추가하여 (총 7가지의 패턴) 테이블을 구성
	//var PATTERN_NAME = []string{"칠정", "노권", "담음", "식적", "어혈"}                       // 변증 이름
	//var PATTERN_INDEX = map[string]int{"칠정": 0, "노권": 1, "담음": 2, "식적": 3, "어혈": 4} // 변증 인덱스 : 이름
	//const PATTERN_NUM = 5

	medicalRecords, flag := RetreiveRecentMedicalRecordByUserID(userID)

	if flag == false { // 충분한 수의 문진 기록이 없는 경우 : False Flag를 반환하여, 최근 문진 기록에 대한 설명기능을 비활성
		var nilTable [question.PATTERN_NUM + 2][NUM_MR_to_CHECK]int
		return nilTable, medicalRecords, flag
	} else {
		var mrTable [question.PATTERN_NUM + 2][NUM_MR_to_CHECK]int //기본 5가지의 패턴과 미병의심, 건강의 2 가지 패턴을 추가하여 (총 7가지) 테이블을 구성

		for index, mrRecord := range medicalRecords {
			for _, pattern := range mrRecord.Pattern {
				// 사용자가 해당 질환(패턴)을 가진 경우 Table 내의 값은 1로 저장
				if pattern == COMPLECATION {
					mrTable[COMPLECATION_INDEX][index] = 1
				} else if pattern == PATTERN_NON {
					mrTable[PATTERN_NON_INDEX][index] = 1
				} else {
					mrTable[question.PATTERN_INDEX[pattern]][index] = 1
				}

			}
		}

		//log.Println(medicalRecords[0].TimeStamp.Year())
		////////
		// var mrTable_ChgReport [question.PATTERN_NUM][NUM_MR_to_CHECK - 1]int // MRTABLE내 사용자의 질환기록 중 변화(쾌유 혹은 발병 등)를 저장하는 테이블 (Change Report)
		// for i := 0; i < question.PATTERN_NUM; i++ {
		// 	for j := 0; j < NUM_MR_to_CHECK-1; j++ {
		// 		mrTable_ChgReport[i][j] = (mrTable[i][j+1] - mrTable[i][j])
		// 	}
		// }
		////////
		return mrTable, medicalRecords, flag
	}

}

func SaveResult_and_CurationDataAtDB() {

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
	c := session.DB(Database).C(RnCCollection)

	rc_file, _ := os.Open("resources/data/CDI_AISpeaker_ResultAndCuration0317.csv") //result&curation file
	rc_reader := csv.NewReader(bufio.NewReader(rc_file))
	rows, _ := rc_reader.ReadAll()

	for i, row := range rows {
		for j := range row {
			log.Printf("%s", rows[i][j])
		}
		log.Println()
		break
	}

	for i := question.FIRST_IDX; i < len(rows); i++ {

		pattern := rows[i][0]

		//복합 질환인 경우 pattern 변수하나에 두 질환을 합쳐 저
		if rows[i][1] != "" {
			pattern += rows[i][1]
		}

		description := rows[i][2]
		explanation := []string{rows[i][3], rows[i][4], rows[i][5], rows[i][6]}
		var curation []string
		for j := 7; j < len(rows[i]); j++ {
			curation = append(curation, rows[i][j])
		}
		temp := ResultAndCuration{
			Pattern:     pattern,     // Pattern     []string `bson:"pattern"`
			Description: description, // Description string   `bson:"description"`
			Explanation: explanation, // Explanation []string `bson:"explanation"`
			Curation:    curation,    // Curation    []string `bson:"curation"`
		}
		if err := c.Insert(temp); err != nil {
			panic(err)
		}
	}
}

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
