package main

import (
	"fmt"
	"log"

	// "munzini/DB"
	"munzini/handler"
	"net/http"
	"os"

	// "time"

	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

// import (
// 	"log"
// 	"munzini/DB"
// 	"munzini/handler"
// 	"net/http"
// 	"os"
// )

func main() {
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

	// // // Insert
	// c := session.DB(DB.Database).C(DB.MRCollection)
	// recordID := bson.NewObjectId()

	// temp := DB.MedicalRecord{

	// 	RecordID:     recordID,
	// 	UserID:       "123",
	// 	TimeStamp:    time.Now(),
	// 	QuestionType: 1,
	// 	Pattern:      []string{"담읍", "심혈"},
	// 	TherapyID:    "123",
	// }

	// // Insert
	// if err := c.Insert(temp); err != nil {
	// 	panic(err)
	// }

	fileServer := http.FileServer(http.Dir("resources"))
	http.Handle("/resources/", http.StripPrefix("/resources/", fileServer))
	http.HandleFunc("/munzini", handler.Dispatch)
	http.HandleFunc("/health_check", handler.HealthCheck)
	http.HandleFunc("/monitor/l7check", handler.HealthCheck)
	port := os.Getenv("PORT") // for server
	//port := "443"                                   // for local test
	log.Fatalln(http.ListenAndServe(":"+port, nil)) // updated: ":13780" -> ":" + port
}
