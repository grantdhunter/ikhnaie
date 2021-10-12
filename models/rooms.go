package models

import (
	"log"
	"sync"
)

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
