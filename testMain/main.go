package main

import (
	"fmt"
	"strings"
	//"time"
	//	"bufio"
	//"encoding/csv"
	//"munzini/DB"
	//"os"
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

// func main() {

// 	recordID := bson.NewObjectId()

// 	temp := DB.MedicalRecord{

// 		RecordID:     recordID,
// 		UserID:       "123",
// 		TimeStamp:    time.Now(),
// 		QuestionType: 1,
// 		Pattern:      []string{"담읍", "심혈"},
// 		TherapyID:    "123",
// 	}

// 	DB.InsertMedicalRecord(temp)
// }

// type Game struct {
// 	Winner       string    `bson:"winner"`
// 	OfficialGame bool      `bson:"official_game"`
// 	Location     string    `bson:"location"`
// 	StartTime    time.Time `bson:"start"`
// 	EndTime      time.Time `bson:"end"`
// 	Players      []Player  `bson:"players"`
// }

// type Player struct {
// 	Name   string    `bson:"name"`
// 	Decks  [2]string `bson:"decks"`
// 	Points uint8     `bson:"points"`
// 	Place  uint8     `bson:"place"`
// }

// //생성
// func NewPlayer(name, firstDeck, secondDeck string, points, place uint8) Player {
// 	return Player{
// 		Name:   name,
// 		Decks:  [2]string{firstDeck, secondDeck},
// 		Points: points,
// 		Place:  place,
// 	}
// }

// var isDropMe = true

func main() {
	// Host := []string{
	// 	"127.0.0.1:27017",
	// 	// replica set addrs...
	// }

	// //상수값 설정
	// const (
	// 	Username   = "YOUR_USERNAME"
	// 	Password   = "YOUR_PASS"
	// 	Database   = "YOUR_DB"
	// 	Collection = "YOUR_COLLECTION"
	// )
	// session, err := mgo.DialWithInfo(&mgo.DialInfo{
	// 	Addrs: Host,
	// 	// Username: Username,
	// 	// Password: Password,
	// 	// Database: Database,
	// 	// DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
	// 	// 	return tls.Dial("tcp", addr.String(), &tls.Config{})
	// 	// },
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()

	// rc_file, _ := os.Open("resources/data/CDI_AISpeaker_ResultAndCuration0317.csv") //result&curation file
	// rc_reader := csv.NewReader(bufio.NewReader(rc_file))
	// rows, _ := rc_reader.ReadAll()

	// for i, row := range rows {
	// 	for j := range row {
	// 		fmt.Printf("%s", rows[i][j])
	// 	}
	// 	fmt.Println()
	// }

	a := []string{"123", "  "}
	b := strings.Join(a, " ")
	c := strings.Trim(b, " ")

	fmt.Println(c)

	/////
	// qcwp, _ := qcwp_reader.ReadAll() // read QCWP.csv
	// ptoc, _ := ptoc_reader.ReadAll() // read cutoff.csv

	// // test code to show how qcwp data looks like
	// /*
	// 	printArray(qcwp)
	// 	printArray(ptoc)
	// */
	// // make map from slice ptoc
	// var ptocMap map[string]int
	// ptocMap = make(map[string]int)
	// for i := FIRST_IDX; i < len(ptoc); i++ {
	// 	ptocMap[ptoc[i][PTOC_PATTERN]], _ = strconv.Atoi(ptoc[i][PTOC_CUTOFF])
	// }

	// // test code to show how ptocMap looks like
	// // fmt.Print(ptocMap)

	// qdatacon := qDataConst{
	// 	QCWP: qcwp,
	// 	PtoC: ptocMap,
	// }

	// return qdatacon

	////

	// fmt.Printf("Connected to %v!\n", session.LiveServers())

	// game := Game{
	// 	Winner:       "Dave",
	// 	OfficialGame: true,
	// 	Location:     "Austin",
	// 	StartTime:    time.Date(2015, time.February, 12, 04, 11, 0, 0, time.UTC),
	// 	EndTime:      time.Now(),
	// 	Players: []Player{
	// 		NewPlayer("Dave", "Wizards", "Steampunk", 21, 1),
	// 		NewPlayer("Javier", "Zombies", "Ghosts", 18, 2),
	// 		NewPlayer("George", "Aliens", "Dinosaurs", 17, 3),
	// 		NewPlayer("Seth", "Spies", "Leprechauns", 10, 4),
	// 	},
	// }

	// //DB 삭제
	// if isDropMe {
	// 	err = session.DB("TEST").DropDatabase()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// // Collection
	// c := session.DB(Database).C(Collection)

	// // Insert
	// if err := c.Insert(game); err != nil {
	// 	panic(err)
	// }

	// // Find and Count
	// player := "Dave"
	// gamesWon, err := c.Find(bson.M{"winner": player}).Count()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s has won %d games.\n", player, gamesWon)

	// // Find One (with Projection)
	// var result Game
	// err = c.Find(bson.M{"winner": player, "location": "Austin"}).Select(bson.M{"official_game": 1}).One(&result)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Is game in Austin Official?", result.OfficialGame)

	// // Find All
	// var games []Game
	// err = c.Find(nil).Sort("-start").All(&games)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Number of Games", len(games))

	// // Update
	// newPlayer := "John"
	// selector := bson.M{"winner": player}
	// updator := bson.M{"$set": bson.M{"winner": newPlayer}}
	// if err := c.Update(selector, updator); err != nil {
	// 	panic(err)
	// }

	// // Update All
	// info, err := c.UpdateAll(selector, updator)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Updated", info.Updated)

	// // Remove
	// info, err = c.RemoveAll(bson.M{"winner": newPlayer})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Removed", info.Removed)
}
