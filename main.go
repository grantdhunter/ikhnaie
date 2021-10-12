package main

import (
	"log"
	"net/http"

	"github.com/grantdhunter/ikhnaie/handlers"
	"github.com/grantdhunter/ikhnaie/middleware"
	"github.com/grantdhunter/ikhnaie/models"
)

func main() {

	fileServer := http.FileServer(http.Dir("./static"))

	intMux := http.NewServeMux()
	intMux.HandleFunc("/", handlers.TrackHandler)
	intMux.HandleFunc("/ws/", handlers.WsHandler)
	authMux := middleware.NewAuth(intMux)
	//	authMux := intMux

	pubMux := http.NewServeMux()
	pubMux.HandleFunc("/auth/", handlers.AuthHandler)
	pubMux.HandleFunc("/oauth/redirect/", handlers.OauthRedirectHandler)
	pubMux.HandleFunc("/", handlers.IndexHandler)

	pubMux.Handle("/static/", http.StripPrefix("/static", fileServer))

	pubMux.Handle("/app/", authMux)
	go handlers.EventHandler()

	log.Print("Starting server on: ", ":"+models.Config.Port)
	err := http.ListenAndServeTLS(":"+models.Config.Port, "server.crt", "server.key", pubMux)

	if err != nil {
		log.Fatal("Error Starting the HTTP Server :", err)
		return
	}
}
