package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/gorilla/websocket"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, config)
}

type Location struct {
	Lat     float32 `json:"lat"`
	Lng     float32 `json:"lng"`
	Acc     float32 `json:"acc"`
	Heading float32 `json:"heading"`
	Speed   float32 `json:"speed"`
	TS      uint64  `json:"ts"`
}

type Event struct {
	RoomId   string   `json:"roomId"`
	ClientId string   `json:"clientId"`
	Name     string   `json:"name"`
	Loc      Location `json:"loc"`
}

func (e Event) String() string {
	str, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s", str)

}

func NewClient(id string, name string, conn *websocket.Conn) Client {
	return Client{Id: id, Name: name, Conn: conn}
}

type Client struct {
	Id   string
	Name string
	Conn *websocket.Conn
}

func NewRoom(id string) Room {
	return Room{Id: id, Clients: make(map[string]*Client)}
}

type Room struct {
	sync.RWMutex
	Id      string
	Clients map[string]*Client
}

func (r *Room) Join(id string, c *Client) {
	r.Lock()
	defer r.Unlock()
	r.Clients[id] = c
}

func (r *Room) Leave(id string) {
	r.Lock()
	defer r.Unlock()
	delete(r.Clients, id)
}

func (r *Room) Iter() <-chan *Client {
	c := make(chan *Client)
	f := func() {
		r.Lock()
		defer r.Unlock()

		for _, v := range r.Clients {
			c <- v
		}
		close(c)
	}
	go f()
	return c

}

func NewRooms() Rooms {
	return Rooms{Rooms: make(map[string]*Room)}
}

type Rooms struct {
	sync.RWMutex
	Id    string
	Rooms map[string]*Room
}

func (rs *Rooms) Add(id string, r *Room) {
	rs.Lock()
	defer rs.Unlock()
	rs.Rooms[id] = r
}

func (rs *Rooms) Remove(id string) {
	rs.Lock()
	defer rs.Unlock()
	delete(rs.Rooms, id)
}

func (rs *Rooms) Get(id string) *Room {
	rs.Lock()
	defer rs.Unlock()

	if rs.Rooms[id] == nil {
		log.Println("No Room found:", id)
		r := NewRoom(id)
		rs.Rooms[id] = &r
	}
	return rs.Rooms[id]
}

func (rs *Rooms) Iter() <-chan *Room {
	c := make(chan *Room)
	f := func() {
		rs.Lock()
		defer rs.Unlock()

		for _, r := range rs.Rooms {
			c <- r
		}
		close(c)
	}
	go f()
	return c

}

var upgrader = websocket.Upgrader{}
var rooms = NewRooms()
var broadcast = make(chan Event)

func wsHandler(w http.ResponseWriter, r *http.Request) {
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
		var evt Event
		err := conn.ReadJSON(&evt)
		if err != nil {
			log.Print("Error during message reading:", err)
			break
		}
		log.Print("Received:", evt)
		room := rooms.Get(evt.RoomId)
		client := NewClient(evt.ClientId, evt.Name, conn)
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
