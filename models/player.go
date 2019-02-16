package models

import "sync"

type Player struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Age   string   `json:"age"`
	Teams []string `json:"teams"`
}

type PlayerCollection struct {
	Mutex         *sync.Mutex
	PlayersLookup map[string]Player
	TeamCount     int
}
