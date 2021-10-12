package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/grantdhunter/ikhnaie/models"
)

var upgrader = websocket.Upgrader{}
var rooms = models.NewRooms()
var broadcast = make(chan models.Event)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	clientId := query["clientId"][0]
	log.Print("CONNECTED: ", clientId)
	for {
		var evt models.Event
		err := conn.ReadJSON(&evt)
		if err != nil {
			log.Print("Error during message reading:", err)
			break
		}
		log.Print("Received:", evt)
		room := rooms.Get(evt.RoomId)
		client := models.NewClient(evt.ClientId, evt.Name, conn)
		room.Join(evt.ClientId, &client)

		log.Print("Recieved:", evt)
		broadcast <- evt
	}
}

func EventHandler() {
	for {
		evt := <-broadcast
		log.Print("Broadcast:", evt)

		for room := range rooms.Iter() {
			for client := range room.Iter() {

				if client.Id == evt.ClientId {
					continue
				}
				log.Print("Sent:", evt)

				err := client.Conn.WriteJSON(evt)
				if err != nil {
					log.Printf("error: %v", err)
					client.Conn.Close()
					room.Leave(client.Id)
				}
			}
		}
	}
}
