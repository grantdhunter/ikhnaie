package models

import "github.com/gorilla/websocket"

type Client struct {
	Id   string
	Name string
	Conn *websocket.Conn
}

func NewClient(id string, name string, conn *websocket.Conn) Client {
	return Client{Id: id, Name: name, Conn: conn}
}
