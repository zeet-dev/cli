package api

import "time"

type LogEntry struct {
	Text      string
	Timestamp time.Time
}
