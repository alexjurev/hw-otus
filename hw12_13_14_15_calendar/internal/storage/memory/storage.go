package memorystorage

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/util"
)

type Storage struct {
	mu           sync.RWMutex
	data         map[string]storage.Event
	idSeq        int
	firstWeekDay time.Weekday
}

func New() *Storage {
	return &Storage{data: make(map[string]storage.Event), firstWeekDay: time.Monday}
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func (s *Storage) AddEvent(_ context.Context, e *storage.Event) error {
	if !e.EndTime.After(e.StartTime) {
		return fmt.Errorf("start time of the event must be in the future: %w", storage.ErrIncorrectEventTime)
	}

	if e.StartTime.Before(time.Now()) {
		return storage.ErrIncorrectEventTime
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[e.ID]; ok {
		return fmt.Errorf("duplicate ID %q: %w", e.ID, storage.ErrDuplicateEventID)
	}
	if e.ID == "" {
		e.ID = s.nextID()
	}
	s.data[e.ID] = *e
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, id string, e storage.Event) error {
	if e.StartTime.Before(time.Now()) {
		return fmt.Errorf("start time of the event must be in the future: %w", storage.ErrIncorrectEventTime)
	}
	if !e.EndTime.After(e.StartTime) {
		return fmt.Errorf("event end time should be after of start time: %w", storage.ErrIncorrectEventTime)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[id]; !ok {
		return fmt.Errorf("failed to update event with id %q: %w", id, storage.ErrNotFoundEvent)
	}
	e.ID = id
	s.data[e.ID] = e
	return nil
}

func (s *Storage) RemoveEvent(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[id]; !ok {
		return fmt.Errorf("failed to remove event with id %q: %w", id, storage.ErrNotFoundEvent)
	}
	delete(s.data, id)
	return nil
}

func (s *Storage) GetEventsForDay(_ context.Context, date time.Time) ([]storage.Event, error) {
	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endTime := startTime.Add(24 * time.Hour)
	return s.selectByRange(startTime, endTime)
}

func (s *Storage) GetEventsForWeek(_ context.Context, startDate time.Time) ([]storage.Event, error) {
	startTime := util.TruncateToDay(startDate)
	if startTime.Weekday() != s.firstWeekDay {
		return nil, storage.ErrIncorrectStartDate
	}
	endTime := startTime.AddDate(0, 0, 7)
	return s.selectByRange(startTime, endTime)
}

func (s *Storage) GetEventsForMonth(_ context.Context, startDate time.Time) ([]storage.Event, error) {
	startTime := util.TruncateToDay(startDate)
	if startTime.Day() != 1 {
		return nil, storage.ErrIncorrectStartDate
	}
	endTime := startTime.AddDate(0, 1, 0)
	return s.selectByRange(startTime, endTime)
}

func (s *Storage) GetEventsByNotifier(
	ctx context.Context,
	limit int,
	endTime time.Time,
) ([]storage.Event, error) {
	events := make([]storage.Event, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, event := range s.data {
		notifyTime := event.StartTime.Add(time.Hour * time.Duration(event.NotifyBefore))
		if event.NotifyBefore > 0 && !event.IsSent && notifyTime.Before(endTime) {
			events = append(events, event)
			if len(events) == limit {
				return events, nil
			}
		}
	}

	return events, nil
}

func (s *Storage) MarkSentEvents(
	ctx context.Context,
	events []storage.Event,
) error {
	if len(events) == 0 {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := make([]string, len(events))
	for i, e := range events {
		ids[i] = e.ID
	}
	for _, event := range s.data {
		event.IsSent = contains(ids, event.ID)
	}

	return nil
}

func (s *Storage) RemoveAfter(ctx context.Context, time time.Time) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, event := range s.data {
		if event.StartTime.After(time) {
			delete(s.data, k)
		}
	}
	return nil
}

// Select in range [startTime:endTime).
func (s *Storage) selectByRange(startTime time.Time, endTime time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, event := range s.data {
		if (event.StartTime.Equal(startTime) || event.StartTime.After(startTime)) && event.StartTime.Before(endTime) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *Storage) nextID() string {
	s.idSeq++
	return strconv.Itoa(s.idSeq)
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func (s *Storage) AddSenderLog(_ context.Context, _ *rabbit.Message) error {
	return nil
}
