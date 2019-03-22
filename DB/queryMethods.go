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
func InsertMedicalRecord(userID string, questionTYPE int, patterns []string, curationType int, curation string) {

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
	timeStamp := time.Now()

	medicalRecord := MedicalRecord{
		RecordID:     recordID,
		UserID:       userID,
		TimeStamp:    timeStamp,
		QuestionType: questionTYPE,
		Pattern:      patterns,
		CurationType: curationType,
		Curation:     curation,
	}

	//UserRecord에 저장될 간략한 문진 내용 Struct 생성
	simpleMR := Simple_MedicalRecord{
		RecordID:     recordID,
		TimeStamp:    timeStamp,
		QuestionType: questionTYPE,
		Pattern:      patterns,
		CurationType: curationType,
		Curation:     curation,
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
			SimpleMRs:        []Simple_MedicalRecord{},
		}
		InsertUserRecord(temp_user)
	}

	//Push medical-record ID to the repective user's record
	updateC := session.DB(Database).C(URCollection)
	query := bson.M{"userID": medicalRecord.UserID}
	change := bson.M{"$push": bson.M{"recordID": medicalRecord.RecordID, "simpleMedicalRecords": simpleMR}}
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

func SaveUserRecord(userID string) {
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
			SimpleMRs:        []Simple_MedicalRecord{},
		}
		InsertUserRecord(temp_user)
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

	count := 0                                      // 최근 건강 추세 조회를 위한 문진 기록의 수를 count 하기 위한 변수
	flag := false                                   // 최근 건강 추세 조회를 위한 문진 기록의 수가 NUM_MR_to_CHECK이상인지 확인하기 위한 변수
	mrRecords_to_Return := []MedicalRecord{}        //반환할 최근 문진기록 리스트
	for i := len(medicalRecords) - 1; i >= 0; i-- { //최근기록일수록 리스트의 뒤에 위치하기에 역순으로 탐색

		if medicalRecords[i].QuestionType == DETAIL_QUESTION_TYPE { //i번째의 문진 기록이 정밀 검진인 경우
			mrRecords_to_Return = append(mrRecords_to_Return, medicalRecords[i])
			count++
			if medicalRecords[i-1].QuestionType == SIMPLE_QUESTION_TYPE { // 정밀검진은 항상 간단 검진 결과 뒤에 저장된다.
				//최근 기록을 조회하여 건강 추세를 분석할 때 간단 검진과 정밀 검진 결과가 모두 있다면 정밀 검진 결과만 사용하므로, 그 앞의 간단 문진 결과를 포함하지 않기 위함
				i -= 1
			}
		} else { //medicalRecords[i].QuestionType ==SIMPLE_QUESTION_TYPE
			mrRecords_to_Return = append(mrRecords_to_Return, medicalRecords[i])
			count++

		}

		if count >= NUM_MR_to_CHECK {
			flag = true
			break
		}
	}

	//위 For문을 통해 Append를 한 경우에는 가장 최근의 문진 기록일 수록 Array의 앞에 위치하게 된다. 최근 건강 추세를 분석하기 위해서는 이전 문진부터 최근 문진 순으로 Array에 저장되는 것이 바람직 하므로 이를 Reverse 해준다.
	for i, j := 0, len(mrRecords_to_Return)-1; i < j; i, j = i+1, j-1 {
		mrRecords_to_Return[i], mrRecords_to_Return[j] = mrRecords_to_Return[j], mrRecords_to_Return[i]
	}

	log.Println(mrRecords_to_Return, flag)
	return mrRecords_to_Return, flag

	// // 충분한 수의 문진기록이 저장되어있는지 확인
	// if len(medicalRecords) < NUM_MR_to_CHECK {
	// 	flag := false
	// 	return nil, flag
	// }

	// flag := true
	// // 가장 최근의 문진 기록만을 반환
	// return medicalRecords[len(medicalRecords)-NUM_MR_to_CHECK : len(medicalRecords)], flag

}

/*
* Author: Jun
* UserID에 해당하는 사용자의 문진 기록을 불러와 반환
 */
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
				if pattern == COMPLECATION { //복합적 미병인 경우, 5 가지의 모든 패턴은 1로 저장(3 가지 이상 패턴의 조합인데, 어떠한 조합인지 알 수 없으므로 모두 1로 설정)
					mrTable[COMPLECATION_INDEX][index] = 1
					for _, pattern_Temp := range question.PATTERN_NAME { //칠정, 노권, 담음, 식적, 어혈 필드를 모두 1로 설정
						mrTable[question.PATTERN_INDEX[pattern_Temp]][index] = 1
					}
				} else if pattern == PATTERN_NON {
					mrTable[PATTERN_NON_INDEX][index] = 1
				} else {
					mrTable[question.PATTERN_INDEX[pattern]][index] = 1
				}

			}
		}
		//log.Println(mrTable)
		return mrTable, medicalRecords, flag
	}

}

/*
* Author: Jun
* 질환들의 설명, 해설, 및 처방들을 Resource 폴더의 CDI_AISpeaker_Result_And_Curation.csv로부터 DB에 업데이트 하는 함수
 */
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

	//Before update all data, clear the Collection(DB)
	c.RemoveAll(bson.M{})

	rc_file, _ := os.Open("resources/data/CDI_AISpeaker_ResultAndCuration0317.csv") //result&curation file
	rc_reader := csv.NewReader(bufio.NewReader(rc_file))
	rows, _ := rc_reader.ReadAll()

	for i := question.FIRST_IDX; i < len(rows); i++ {

		pattern := rows[i][0]

		//복합 질환인 경우 pattern 변수하나에 두 질환을 합쳐 저
		if rows[i][1] != "" {
			pattern += (" " + rows[i][1])
		}

		description := rows[i][2]
		explanation := []string{rows[i][3], rows[i][4], rows[i][5], rows[i][6]}
		var dietCuration []string

		for j := 7; j < len(rows[i]); j++ {

			// 해당  필드가 비어있는 경우의 예외처리
			if rows[i][j] != "" {
				dietCuration = append(dietCuration, rows[i][j])
			}
		}
		temp := ResultAndCuration{
			Pattern:      pattern,      // Pattern     []string `bson:"pattern"`
			Description:  description,  // Description string   `bson:"description"`
			Explanation:  explanation,  // Explanation []string `bson:"explanation"`
			DietCuration: dietCuration, // Curation    []string `bson:"curation"`

			//TODO 추후 CDI_AISpeaker_Result and Curation 파일에 아래 3 가지 요법에 대한 추천항목이 업데이트 될시 아래 코드에 업데이트 필요
			ExerciseCuration: []string{},
			YangSangCuration: []string{},
			CDM_Curation:     []string{},
		}
		if err := c.Insert(temp); err != nil {
			panic(err)
		}
	}
}

/*
* Author: Jun
* Result and Curation Collection(DB)로 부터 Pattern에 해당하는 문진 설명을 불러온다.
 */
func GetResult_and_Explanation(pattern string) string {

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

	// Find First, If user is not exist in database, add his data
	findC := session.DB(Database).C(RnCCollection)

	var rncInfo ResultAndCuration //result&curation Info that matches to the given userID
	if findErr := findC.Find(bson.M{"pattern": pattern}).One(&rncInfo); findErr != nil {
		panic(findErr)
	}

	//Index 2: Where the explnation is for Simple Question Score
	return rncInfo.Explanation[RAC_SQS_EXPLANATION_INDEX]
}

func GetResult_and_Curation(pattern string) ResultAndCuration {

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

	// Find First, If user is not exist in database, add his data
	findC := session.DB(Database).C(RnCCollection)

	var rncInfo ResultAndCuration //result&curation Info that matches to the given userID
	if findErr := findC.Find(bson.M{"pattern": pattern}).One(&rncInfo); findErr != nil {
		panic(findErr)
	}

	return rncInfo
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
