package storage

import "github.com/TauAdam/timer-bot/internal/timer"

type Storage interface {
	AddTimer(id string, timer timer.Timer) error
	GetTimer(id string) (timer.Timer, error)
	ResetTimer(id string) error
}
