package storage

import (
	"fmt"
	"time"
)

type Event struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	Description  string    `json:"description"`
	OwnerID      string    `json:"ownerId"`
	NotifyBefore int32     `json:"notifyBefore"`
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
