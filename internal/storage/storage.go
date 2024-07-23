package storage

import "github.com/TauAdam/timer-bot/internal/timer"

type Storage interface {
	AddTimer(id string, timer timer.Timer) error
}
