// // package main

// // import (
// // 	"fmt"
// // 	"strings"
// // 	//"time"
// // 	//	"bufio"
// // 	//"encoding/csv"
// // 	//"munzini/DB"
// // 	//"os"
// // 	//"gopkg.in/mgo.v2"
// // 	//"gopkg.in/mgo.v2/bson"
// // )

// // // func main() {

// // // 	recordID := bson.NewObjectId()

// // // 	temp := DB.MedicalRecord{

// // // 		RecordID:     recordID,
// // // 		UserID:       "123",
// // // 		TimeStamp:    time.Now(),
// // // 		QuestionType: 1,
// // // 		Pattern:      []string{"담읍", "심혈"},
// // // 		TherapyID:    "123",
// // // 	}

// // // 	DB.InsertMedicalRecord(temp)
// // // }

// // // type Game struct {
// // // 	Winner       string    `bson:"winner"`
// // // 	OfficialGame bool      `bson:"official_game"`
// // // 	Location     string    `bson:"location"`
// // // 	StartTime    time.Time `bson:"start"`
// // // 	EndTime      time.Time `bson:"end"`
// // // 	Players      []Player  `bson:"players"`
// // // }

// // // type Player struct {
// // // 	Name   string    `bson:"name"`
// // // 	Decks  [2]string `bson:"decks"`
// // // 	Points uint8     `bson:"points"`
// // // 	Place  uint8     `bson:"place"`
// // // }

// // // //생성
// // // func NewPlayer(name, firstDeck, secondDeck string, points, place uint8) Player {
// // // 	return Player{
// // // 		Name:   name,
// // // 		Decks:  [2]string{firstDeck, secondDeck},
// // // 		Points: points,
// // // 		Place:  place,
// // // 	}
// // // }

// // // var isDropMe = true

// // func main() {
// // 	// Host := []string{
// // 	// 	"127.0.0.1:27017",
// // 	// 	// replica set addrs...
// // 	// }

// // 	// //상수값 설정
// // 	// const (
// // 	// 	Username   = "YOUR_USERNAME"
// // 	// 	Password   = "YOUR_PASS"
// // 	// 	Database   = "YOUR_DB"
// // 	// 	Collection = "YOUR_COLLECTION"
// // 	// )
// // 	// session, err := mgo.DialWithInfo(&mgo.DialInfo{
// // 	// 	Addrs: Host,
// // 	// 	// Username: Username,
// // 	// 	// Password: Password,
// // 	// 	// Database: Database,
// // 	// 	// DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
// // 	// 	// 	return tls.Dial("tcp", addr.String(), &tls.Config{})
// // 	// 	// },
// // 	// })
// // 	// if err != nil {
// // 	// 	panic(err)
// // 	// }
// // 	// defer session.Close()

// // 	// rc_file, _ := os.Open("resources/data/CDI_AISpeaker_ResultAndCuration0317.csv") //result&curation file
// // 	// rc_reader := csv.NewReader(bufio.NewReader(rc_file))
// // 	// rows, _ := rc_reader.ReadAll()

// // 	// for i, row := range rows {
// // 	// 	for j := range row {
// // 	// 		fmt.Printf("%s", rows[i][j])
// // 	// 	}
// // 	// 	fmt.Println()
// // 	// }

// // 	a := []string{"123", "  "}
// // 	b := strings.Join(a, " ")
// // 	c := strings.Trim(b, " ")

// // 	fmt.Println(c)

// // 	/////
// // 	// qcwp, _ := qcwp_reader.ReadAll() // read QCWP.csv
// // 	// ptoc, _ := ptoc_reader.ReadAll() // read cutoff.csv

// // 	// // test code to show how qcwp data looks like
// // 	// /*
// // 	// 	printArray(qcwp)
// // 	// 	printArray(ptoc)
// // 	// */
// // 	// // make map from slice ptoc
// // 	// var ptocMap map[string]int
// // 	// ptocMap = make(map[string]int)
// // 	// for i := FIRST_IDX; i < len(ptoc); i++ {
// // 	// 	ptocMap[ptoc[i][PTOC_PATTERN]], _ = strconv.Atoi(ptoc[i][PTOC_CUTOFF])
// // 	// }

// // 	// // test code to show how ptocMap looks like
// // 	// // fmt.Print(ptocMap)

// // 	// qdatacon := qDataConst{
// // 	// 	QCWP: qcwp,
// // 	// 	PtoC: ptocMap,
// // 	// }

// // 	// return qdatacon

// // 	////

// // 	// fmt.Printf("Connected to %v!\n", session.LiveServers())

// // 	// game := Game{
// // 	// 	Winner:       "Dave",
// // 	// 	OfficialGame: true,
// // 	// 	Location:     "Austin",
// // 	// 	StartTime:    time.Date(2015, time.February, 12, 04, 11, 0, 0, time.UTC),
// // 	// 	EndTime:      time.Now(),
// // 	// 	Players: []Player{
// // 	// 		NewPlayer("Dave", "Wizards", "Steampunk", 21, 1),
// // 	// 		NewPlayer("Javier", "Zombies", "Ghosts", 18, 2),
// // 	// 		NewPlayer("George", "Aliens", "Dinosaurs", 17, 3),
// // 	// 		NewPlayer("Seth", "Spies", "Leprechauns", 10, 4),
// // 	// 	},
// // 	// }

// // 	// //DB 삭제
// // 	// if isDropMe {
// // 	// 	err = session.DB("TEST").DropDatabase()
// // 	// 	if err != nil {
// // 	// 		panic(err)
// // 	// 	}
// // 	// }

// // 	// // Collection
// // 	// c := session.DB(Database).C(Collection)

// // 	// // Insert
// // 	// if err := c.Insert(game); err != nil {
// // 	// 	panic(err)
// // 	// }

// // 	// // Find and Count
// // 	// player := "Dave"
// // 	// gamesWon, err := c.Find(bson.M{"winner": player}).Count()
// // 	// if err != nil {
// // 	// 	panic(err)
// // 	// }
// // 	// fmt.Printf("%s has won %d games.\n", player, gamesWon)

// // 	// // Find One (with Projection)
// // 	// var result Game
// // 	// err = c.Find(bson.M{"winner": player, "location": "Austin"}).Select(bson.M{"official_game": 1}).One(&result)
// // 	// if err != nil {
// // 	// 	panic(err)
// // 	// }
// // 	// fmt.Println("Is game in Austin Official?", result.OfficialGame)

// // 	// // Find All
// // 	// var games []Game
// // 	// err = c.Find(nil).Sort("-start").All(&games)
// // 	// if err != nil {
// // 	// 	panic(err)
// // 	// }
// // 	// fmt.Println("Number of Games", len(games))

// // 	// // Update
// // 	// newPlayer := "John"
// // 	// selector := bson.M{"winner": player}
// // 	// updator := bson.M{"$set": bson.M{"winner": newPlayer}}
// // 	// if err := c.Update(selector, updator); err != nil {
// // 	// 	panic(err)
// // 	// }

// // 	// // Update All
// // 	// info, err := c.UpdateAll(selector, updator)
// // 	// if err != nil {
// // 	// 	panic(err)
// // 	// }
// // 	// fmt.Println("Updated", info.Updated)

// // 	// // Remove
// // 	// info, err = c.RemoveAll(bson.M{"winner": newPlayer})
// // 	// if err != nil {
// // 	// 	panic(err)
// // 	// }
// // 	// fmt.Println("Removed", info.Removed)
// // }

// // package main

// // import (
// // 	"bufio"
// // 	"encoding/csv"

// // 	"fmt"
// // 	"os"
// // )

// // func main() {
// // 	// open QCWP file	- Use CWP ( Category-Weight-Pattern )
// // 	qcwp_file, _ := os.Open("resources/data/QCWP.csv")

// // 	// create csv Reader
// // 	rdr := csv.NewReader(bufio.NewReader(qcwp_file))

// // 	// read csv file
// // 	rows, _ := rdr.ReadAll()
// // 	printArray(rows)
// // }

// // func printArray(arr [][]string) {
// // 	for i := 1; i < len(arr); i++ {
// // 		for j := 0; j < len(arr[i]); j++ {
// // 			fmt.Printf("%s ", arr[i][j])
// // 		}
// // 		fmt.Println()
// // 	}
// // }

// // package main

// // import "fmt"

// // type AandB struct {
// // 	A string
// // 	B string
// // }

// // func main() {
// // 	var X AandB
// // 	var Y AandB
// // 	X.A = "칠정"
// // 	X.B = "Neuro"
// // 	Y.A = "노권"
// // 	Y.B = "Neuro"
// // 	var a map[AandB]string
// // 	a = make(map[AandB]string)
// // 	a[X] = X.A
// // 	a[Y] = Y.A
// // 	fmt.Println(a[X])
// // 	fmt.Println(a)
// // }

// package main

// import (
// 	"bufio"
// 	"encoding/csv"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"test"
// )

// type QueryData struct { // Query Data : 총 23개
// 	Pattern              string // 변증 이름
// 	Category             string // 카테고리 이름
// 	Half_Of_Category_Num int    // 카테고리별 질문 수의 절반 = 가중치 / 2
// 	ShouldBeQueried      bool   // 추천 DB에 쿼리를 날려야하는가?	- 1. 정밀 진단 결과 해당하는 변증인가? ( Key = Pattern ), 2. 진단 결과 HOCN 의 값이 양인가?
// }

// type PatternCat struct { // Queries의 Key 구조체
// 	Pattern  string
// 	Category string
// }

// type Queries struct {
// 	QueryCore    map[PatternCat]QueryData // Pattern & Category ( = Key )로 QueryData ( = Value ) 접근
// 	QueryStrings []string                 // Query문들
// }

// func loadData() Queries {
// 	// open QCWP file	- Use CWP ( Category-Weight-Pattern )
// 	qcwp_file, _ := os.Open("../resources/data/QCWP.csv")

// 	// create csv Reader
// 	qcwp_reader := csv.NewReader(bufio.NewReader(file))

// 	// read csv file
// 	qcwp, _ := qcwp_reader.ReadAll()

// 	/*
// 		TODO
// 		1. QueryCore를 초기화한다.
// 			2. QueryCore를 초기화하기 위하여 PatternCat 리스트를 만든다.
// 			3. QueryCore의 Key값에 PatternCat을 넣고, 이에 따른 QueryData를 작성하는 로직을 만든다.
// 	*/

// 	// 1. PatternCat 초기화 ( QueryCore의 Key 값 )
// 	var patcat []PatternCat

// 	// TODO: PatternCat리스트 초기화 ( 23개 - 카테고리 개수)
// 	var row int = FIRST_IDX
// 	var temp_weight []int	// 추후에 QueryCore의 Value값 중 Half_Of_Category_Num에 값을 담아놓기 위해 가중치 값들을 미리 저장해놓는 슬라이스
// 	var patcat_idx int = 0
// 	for row < len(qcwp) {
// 		patcat[patcat_idx].Pattern = qcwp[row][PATTERN_IDX]
// 		patcat[patcat_idx].Category = qcwp[row][CATEGORY_IDX]
// 		temp_weight[patcat_idx] = qcwp[row][WEIGHT_IDX]
// 		row = row + qcwp[row][WEIGHT_IDX] // 가중치만큼 Forwarding
// 		pat_idx++
// 	}

// 	// 2. QueryCore 초기화 ( PatternCat - QueryData : Pattern / Category / Half_Of_Category_Num / ShouldBeQueried )
// 	var queryCore map[PatternCat]QueryData

// 	// TODO: PatternCat의 값을 QueryCore의 Key값에 넣고, 그에 해당하는 QueryData를 작성한다.
// 	for qd_idx := 0; qd_idx < len(patcat); qd_idx++ {
// 		// TODO:구조체 초기화
// 		queryCore[patcat[qd_idx]] = QueryData{
// 			Pattern:  patcat[qd_idx].Pattern,
// 			Category: patcat[qd_idx].Category,
// 			Half_Of_Category_Num : temp_weight[qd_idx] / 2,
// 			ShouldBeQueried : true
// 		}
// 	}

// 	// 3. Queries 작성
// 	var queries Queries = Queries{
// 		QueryCore: queryCore,
// 		QueryStrings : nil
// 	}

// 	return queries
// }

// func main() {
// 	// qcwp_file, _ := os.Open("../resources/data/QCWP.csv")
// 	// qcwp_reader := csv.NewReader(bufio.NewReader(qcwp_file))
// 	// qcwp, _ := qcwp_reader.ReadAll()

// 	// fmt.Printf("%T\n", qcwp)

// 	// var temp int = test.GetX()
// 	// fmt.Println("From test : " + strconv.Itoa(temp))

// 	// x := 0
// 	// for x < 100 {
// 	// 	// println(x)
// 	// 	x = x + 4
// 	// }

// 	// fmt.Println(x)

// 	// var A PatternCat
// 	// A = PatternCat{Pattern: "하", Category: "하~"}
// 	// fmt.Println(A.Category)

// 	// fmt.Println(5 / 2)

// 	// var AA Queries = Queries{QueryStrings: nil}
// 	// fmt.Println(AA.QueryCore)
// 	// fmt.Println(AA.QueryStrings)

// 	var queries Queries = loadData()
// 	fmt.Println(queries.QueryCore)
// 	fmt.Println(queries.QueryStrings)

// 	// var patcat []PatternCat

// 	// for row := 1; row < len(qcwp); row++ {
// 	// 	var tempPatcat PatternCat
// 	// 	tempPatcat.Category = qcwp[row][0]
// 	// 	tempPatcat.Pattern = qcwp[row][3]
// 	// 	patcat = append(patcat, tempPatcat)
// 	// 	fmt.Println(patcat)
// 	// }

// 	// a := qcwp[1]
// 	// fmt.Println(a)
// }

package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	// "strings"
	// "go.mongodb.org/mongo-driver/bson"
)

const FIRST_IDX = 1    // QCWP.csv를 담아올 때 접근해야하는 첫번째 인덱스
const CATEGORY_IDX = 1 // QCWP.csv에서 Category에 접근하기 위한 인덱스
const PATTERN_IDX = 3  // QCWP.csv에서 Pattern에 접근하기 위한 인덱스
const WEIGHT_IDX = 2   // QCWP.csv에서 Weight에 접근하기 위한 인덱스

var queries Queries = loadData()

type QueryData struct { // Query Data : 총 23개
	Pattern              string // 변증 이름
	Category             string // 카테고리 이름
	Half_Of_Category_Num int    // 카테고리별 질문 수의 절반 = 가중치 / 2
	ShouldBeQueried      bool   // 추천 DB에 쿼리를 날려야하는가?	- 1. 정밀 진단 결과 해당하는 변증인가? ( Key = Pattern ), 2. 진단 결과 HOCN 의 값이 양인가?
}

type PatternCat struct { // Queries의 Key 구조체
	Pattern  string
	Category string
}

func (pc PatternCat) toString() string {
	return pc.Pattern + " " + pc.Category
}

type SimpleDoc struct {
	Pattern  string `bson:"pattern"`
	Category string `bson:"category"`
	FoodNm   string `bson:"foodNm"`
}

// CEKSessionAttributes를 통하여 주고받아야할 구조체
type Queries struct {
	QueryCore map[string]QueryData // Pattern & Category ( = Key )로 QueryData ( = Value ) 접근
	// 확장을 위하여 남겨두었음.
	// QueryStrings []string
	// QueryOutput [][]SimpleDoc
	// QueryStrings map[PatternCat]string                 // Query문들
}

func loadData() Queries {
	// open QCWP file	- Use CWP ( Category-Weight-Pattern )
	qcwp_file, _ := os.Open("../resources/data/QCWP.csv")

	// create csv Reader
	qcwp_reader := csv.NewReader(bufio.NewReader(qcwp_file))

	// read csv file
	qcwp, _ := qcwp_reader.ReadAll()

	var a []int
	a = append(a, 1)
	fmt.Println("test")
	fmt.Println(len(a))
	fmt.Println("test2")

	fmt.Printf("%T\n", qcwp)

	/*
		TODO
		1. QueryCore를 초기화한다.
			2. QueryCore를 초기화하기 위하여 PatternCat 리스트를 만든다.
			3. QueryCore의 Key값에 PatternCat을 넣고, 이에 따른 QueryData를 작성하는 로직을 만든다.
	*/
	fmt.Println("Step 0 Completed.")
	fmt.Println(qcwp[30][1])
	fmt.Println(len(qcwp))
	fmt.Println(FIRST_IDX)

	// 1. PatternCat 초기화 ( QueryCore의 Key 값 )
	var patcat []PatternCat

	// TODO: PatternCat리스트 초기화 ( 23개 - 카테고리 개수)
	var row int = FIRST_IDX
	var weight []int // 추후에 QueryCore의 Value값 중 Half_Of_Category_Num에 값을 담아놓기 위해 가중치 값들을 미리 저장해놓는 슬라이스

	fmt.Println("Step 1 Started.")
	for row < len(qcwp) {
		// fmt.Println(patcat_idx)
		// patcat[patcat_idx].Pattern = qcwp[row][PATTERN_IDX]temp_weight
		temp_patcat := PatternCat{
			Pattern:  qcwp[row][PATTERN_IDX],
			Category: qcwp[row][CATEGORY_IDX],
		}
		patcat = append(patcat, temp_patcat)
		// patcat[patcat_idx].Category = qcwp[row][CATEGORY_IDX]
		temp_weight, _ := strconv.Atoi(qcwp[row][WEIGHT_IDX])
		// fmt.Println(temp_weight)
		weight = append(weight, temp_weight)
		row_gap, _ := strconv.Atoi(qcwp[row][WEIGHT_IDX])
		row = row + row_gap // 가중치만큼 Forwarding
	}
	fmt.Println("Step 1 Completed.")
	// 2. QueryCore 초기화 ( PatternCat - QueryData : Pattern / Category / Half_Of_Category_Num / ShouldBeQueried )
	var queryCore map[string]QueryData = make(map[string]QueryData)
	fmt.Println(len(queryCore))

	fmt.Println(len(patcat))
	fmt.Println(patcat[0])
	fmt.Println(weight[0])

	// TODO: PatternCat의 값을 QueryCore의 Key값에 넣고, 그에 해당하는 QueryData를 작성한다.
	for qd_idx := 0; qd_idx < len(patcat); qd_idx++ {
		// TODO:구조체 초기화
		queryCore[patcat[qd_idx].toString()] = QueryData{
			Pattern:              patcat[qd_idx].Pattern,
			Category:             patcat[qd_idx].Category,
			Half_Of_Category_Num: weight[qd_idx] / 2,
			ShouldBeQueried:      true,
		}
	}

	// 3. Queries 작성
	var queries Queries = Queries{
		QueryCore: queryCore,
	}
	return queries
}

func main() {

	// fmt.Println("Main Testing")
	// fmt.Println(queries.QueryCore)
	// fmt.Println(len(queries.QueryCore))

	// fmt.Println("07-31 Testing")
	// fmt.Println(3 >= 3)
	// var a []int
	// a = append(a, 3)
	// fmt.Println(a)
	// v := "가 나 다 라 마 바"
	// v_list := strings.Split(v, " ")
	// fmt.Println(v_list[0])

	// var dd map[int]int = make(map[int]int)
	// dd[3] = 5
	// dd[2] = 2
	// dd[1] = 6

	// for k, v := range dd {
	// 	fmt.Println(k, v)
	// }

	// for i := 0; i < len(a); i++ {
	// 	fmt.Println(a[i])
	// }

	// var asd []bson.M

	// fmt.Println("08-02 Testing")
	// fmt.Print(asd)
	// fmt.Printf("%T\n", asd)

	fmt.Println("Main Testing is started.")

	loadData()

	var q Queries
	q.QueryCore = make(map[string]QueryData)
	q.QueryCore[PatternCat{Pattern: "칠정", Category: "카테고리1"}.toString()] = QueryData{
		Pattern:              "칠정",
		Category:             "카테고리1",
		Half_Of_Category_Num: 3,
		ShouldBeQueried:      true,
	}
	for k, v := range q.QueryCore {
		fmt.Println(k)
		fmt.Println(v)
	}

	fmt.Println("JSON Testing is started.")

	doc, err := json.Marshal(q.QueryCore)

	if err != nil {
		fmt.Println("Error occured!!")
		fmt.Println("Error Type :")
		fmt.Printf("%T\n", err)
	} else {
		fmt.Println("JSON Encoding Succeded!")
	}

	fmt.Println(string(doc))

}
