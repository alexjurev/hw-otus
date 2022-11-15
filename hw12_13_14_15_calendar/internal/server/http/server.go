package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type Config struct {
	Host string
	Port int
}

type Server struct {
	srv  *http.Server
	addr string
	app  *app.App
}

func NewServer(config Config, app *app.App) *Server {
	return &Server{
		addr: net.JoinHostPort(config.Host, strconv.Itoa(config.Port)),
		srv:  &http.Server{Addr: net.JoinHostPort(config.Host, strconv.Itoa(config.Port))}, //nolint
		app:  app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HELLO !!!"))
	})
	mux.HandleFunc("/add", s.AddEvent)
	mux.HandleFunc("/update", s.UpdateEvent)
	mux.HandleFunc("/remove", s.RemoveEvent)
	mux.HandleFunc("/events/day", s.GetEventsForDay)
	mux.HandleFunc("/events/week", s.GetEventsForWeek)
	mux.HandleFunc("/events/month", s.GetEventsForMonth)

	s.srv.Handler = loggingMiddleware(mux)

	err := s.srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server failed: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func getIP(req *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	if parsed := net.ParseIP(ip); parsed == nil {
		return "", fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return ip, nil
}

func (s *Server) AddEvent(w http.ResponseWriter, r *http.Request) {
	event := storage.Event{}
	res, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(res, &event)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := s.app.CreateEvent(context.Background(), event)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

type UpdateReq struct {
	ID    string        `json:"id"`
	Event storage.Event `json:"event"`
}

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	updateEvent := UpdateReq{}
	res, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(res, &updateEvent)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.app.UpdateEvent(context.Background(), updateEvent.ID, updateEvent.Event)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) RemoveEvent(w http.ResponseWriter, r *http.Request) {
	event := storage.Event{}
	res, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(res, &event)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.app.RemoveEvent(context.Background(), event.ID)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type ReqDate struct {
	Date time.Time `json:"date"`
}

const (
	day   = "day"
	week  = "week"
	month = "month"
)

func (s *Server) GetEvents(w http.ResponseWriter, r *http.Request, period string) {
	date := ReqDate{}
	res, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(res, &date)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var event []storage.Event
	switch period {
	case day:
		event, err = s.app.GetEventsForDay(context.Background(), date.Date)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case week:
		event, err = s.app.GetEventsForWeek(context.Background(), date.Date)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case month:
		event, err = s.app.GetEventsForMonth(context.Background(), date.Date)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event)
}

func (s *Server) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	s.GetEvents(w, r, day)
}

func (s *Server) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	s.GetEvents(w, r, week)
}

func (s *Server) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	s.GetEvents(w, r, month)
}
