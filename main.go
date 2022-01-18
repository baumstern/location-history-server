package main

import (
	"location-history-server/server"
	"log"
	"net/http"
	"os"
)

func main() {
	s := server.NewServer()
	http.Handle("/location/", server.AppHandler(s.HandleLocation))

	var port string
	port, ok := os.LookupEnv("HISTORY_SERVER_LISTEN_ADDR")
	if !ok {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
