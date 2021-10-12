package models

import "sync"

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
