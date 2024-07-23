package inmemdb

import (
	"github.com/TauAdam/timer-bot/internal/timer"
)

type InMemoryDB struct {
	timers map[string]timer.Timer
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		timers: make(map[string]timer.Timer),
	}
}

func (db *InMemoryDB) AddTimer(id string, t timer.Timer) error {
	db.timers[id] = t
	return nil
}
