package main

import (
	"log"
	"magicball/handler"
	"net/http"
	"os"	// added
)

func main() {
	fileServer := http.FileServer(http.Dir("resources"))
	http.Handle("/resources/", http.StripPrefix("/resources/", fileServer))
	http.HandleFunc("/magicball", handler.Dispatch)
	http.HandleFunc("/health_check", handler.HealthCheck)
	http.HandleFunc("/monitor/l7check", handler.HealthCheck)

	port := os.Getenv("PORT")	// added
	log.Fatalln(http.ListenAndServe(":" + port, nil))	// updated: ":13780" -> ":" + port
}
