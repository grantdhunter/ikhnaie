package main

import (
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port           string
	Mapbox_api_key string
}

var config Config

func main() {

	if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
		log.Fatal(err)
		return
	}

	fileServer := http.FileServer(http.Dir("./static"))

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/ws/", wsHandler)
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Print("Starting server on: ", ":"+config.Port)

	go EventHandler()
	err := http.ListenAndServeTLS(":"+config.Port, "server.crt", "server.key", nil)

	if err != nil {
		log.Fatal("Error Starting the HTTP Server :", err)
		return
	}
}
