package main

import (
	"log"
	"munzini/handler"
	"net/http"
	//"os"
)

func main() {
	fileServer := http.FileServer(http.Dir("resources"))
	http.Handle("/resources/", http.StripPrefix("/resources/", fileServer))
	http.HandleFunc("/munzini", handler.Dispatch)
	http.HandleFunc("/health_check", handler.HealthCheck)
	http.HandleFunc("/monitor/l7check", handler.HealthCheck)
	//port := os.Getenv("PORT") // for server
	port := "13780"
	log.Fatalln(http.ListenAndServe(":"+port, nil)) // updated: ":13780" -> ":" + port
}
