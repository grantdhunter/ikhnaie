package models

type Location struct {
	Lat     float32 `json:"lat"`
	Lng     float32 `json:"lng"`
	Acc     float32 `json:"acc"`
	Heading float32 `json:"heading"`
	Speed   float32 `json:"speed"`
	TS      uint64  `json:"ts"`
}
