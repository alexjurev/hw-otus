package storage

import (
	"fmt"
	"time"
)

type Event struct {
	ID           string
	Title        string
	StartTime    time.Time
	EndTime      time.Time
	Description  string
	OwnerID      string
	NotifyBefore time.Duration
}

func (e *Event) Validate() error {
	if !e.EndTime.After(e.StartTime) {
		return fmt.Errorf("start time of the event must be in the future: %w", ErrIncorrectEventTime)
	}

	if e.StartTime.Before(time.Now()) {
		return ErrIncorrectEventTime
	}

	return nil
}
