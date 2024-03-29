package app

import (
	"context"
	"time"

	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Storage storage.Storage
}

func New(storage storage.Storage) *App {
	return &App{Storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, e storage.Event) (string, error) {
	if err := a.Storage.AddEvent(ctx, &e); err != nil {
		return "", err
	}
	return e.ID, nil
}

func (a *App) UpdateEvent(ctx context.Context, id string, e storage.Event) error {
	return a.Storage.UpdateEvent(ctx, id, e)
}

func (a *App) RemoveEvent(ctx context.Context, id string) error {
	return a.Storage.RemoveEvent(ctx, id)
}

func (a *App) GetEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.Storage.GetEventsForDay(ctx, date)
}

func (a *App) GetEventsForWeek(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	return a.Storage.GetEventsForWeek(ctx, startDate)
}

func (a *App) GetEventsForMonth(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	return a.Storage.GetEventsForMonth(ctx, startDate)
}
