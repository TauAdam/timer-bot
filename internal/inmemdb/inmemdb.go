package inmemdb

import (
	"fmt"
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

func (db *InMemoryDB) GetTimer(id string) (timer.Timer, error) {
	t, exists := db.timers[id]
	if !exists {
		return timer.Timer{}, fmt.Errorf("timer not found")
	}
	return t, nil
}

func (db *InMemoryDB) DeleteTimer(id string) error {
	_, exists := db.timers[id]
	if !exists {
		return fmt.Errorf("timer not found")
	}
	delete(db.timers, id)
	return nil
}
