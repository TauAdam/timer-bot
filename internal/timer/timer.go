package timer

import "time"

type Timer struct {
	Duration  time.Duration
	StartTime time.Time
	Label     string
}
