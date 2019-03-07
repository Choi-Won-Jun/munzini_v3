package main

import (
	"log"
	"munzinis_project/handler"
	"net/http"
	"os"
)

// Imported Packages only for Local Test
// import (
// 	"fmt"
// 	"log"
// 	"munzinis_project/handler"
// 	"net/http"
// )

func main() {
	fileServer := http.FileServer(http.Dir("resources"))

	http.Handle("/resources/", http.StripPrefix("/resources/", fileServer))
	http.HandleFunc("/munzini", handler.Dispatch)
	http.HandleFunc("/health_check", handler.HealthCheck)
	http.HandleFunc("/monitor/l7check", handler.HealthCheck)
	port := os.Getenv("PORT") // for server
	//port := "13780"                                 // for local test
	log.Fatalln(http.ListenAndServe(":"+port, nil)) // updated: ":13780" -> ":" + port
}
