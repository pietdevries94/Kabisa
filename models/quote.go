package models

type Quote struct {
	ID     int64  `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}
