package main

import (
	// "fmt"
	"log"

	//"munzini/DB"
	"munzini/handler"
	"net/http"
	"os"
	//"time"
	// "gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

func main() {

	// //TODO UserRecord Insert Sample
	// temp_user := DB.UserRecord{
	// 	UserID:           "124",
	// 	UserName:         "124",
	// 	RecordID:         []string{"obj23412", "129dhflb"},
	// 	RegistrationDate: time.Now(),
	// }
	// DB.InsertUserRecord(temp_user)

	// recordID := bson.NewObjectId()

	// temp := DB.MedicalRecord{

	// 	RecordID:     recordID,
	// 	UserID:       "125",
	// 	TimeStamp:    time.Now(),
	// 	QuestionType: 1,
	// 	Pattern:      []string{"담읍", "심혈"},
	// 	TherapyID:    "124",
	// }
	// DB.InsertMedicalRecord(temp)

	fileServer := http.FileServer(http.Dir("resources"))
	http.Handle("/resources/", http.StripPrefix("/resources/", fileServer))
	http.HandleFunc("/munzini", handler.Dispatch)
	http.HandleFunc("/health_check", handler.HealthCheck)
	http.HandleFunc("/monitor/l7check", handler.HealthCheck)
	port := os.Getenv("PORT") // for server
	//port := "443"                                   // for local test
	log.Fatalln(http.ListenAndServe(":"+port, nil)) // updated: ":13780" -> ":" + port
}
