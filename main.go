package main

import (
	"fmt"
	"log"
	"munzini/handler"
	"net/http"
	"os"
	//"gopkg.in/mgo.v2"
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
	// sess, err := mgo.Dial(uri)
	// if err != nil {
	// 	//fmt.Printf("Can't connect to mongo, go error %v\n", err)
	// 	os.Exit(1)
	// }
	// defer sess.Close()

	fileServer := http.FileServer(http.Dir("resources"))
	http.Handle("/resources/", http.StripPrefix("/resources/", fileServer))
	http.HandleFunc("/munzini", handler.Dispatch)
	http.HandleFunc("/health_check", handler.HealthCheck)
	http.HandleFunc("/monitor/l7check", handler.HealthCheck)
	port := os.Getenv("PORT") // for server
	//port := "443"                                   // for local test
	log.Fatalln(http.ListenAndServe(":"+port, nil)) // updated: ":13780" -> ":" + port
}
