package models

import (
	"encoding/json"
	"fmt"
)

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
