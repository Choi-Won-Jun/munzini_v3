package main

import (
	//"fmt"
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
	//CDI_AISpeaker_ResultAndCuration0317.csv 파일 업데이트시 아래 코드의 주석을 해제하여 DB에 업데이트
	//DB.SaveResult_and_CurationDataAtDB()
	fileServer := http.FileServer(http.Dir("resources"))
	http.Handle("/resources/", http.StripPrefix("/resources/", fileServer))
	http.HandleFunc("/munzini", handler.Dispatch)
	http.HandleFunc("/health_check", handler.HealthCheck)
	http.HandleFunc("/monitor/l7check", handler.HealthCheck)
	port := os.Getenv("PORT") // for server
	// port := "443"                                   // for local test
	log.Fatalln(http.ListenAndServe(":"+port, nil)) // updated: ":13780" -> ":" + port
}
