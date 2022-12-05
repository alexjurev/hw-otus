//go:build integration_test
// +build integration_test

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/rabbit"

	internalhttp "github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/server/http"

	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	httpServerHost = "127.0.0.1"
	httpServerPort = 10080
	grpcServerHost = "127.0.0.1"
	grpcServerPort = 10081
	pgHost         = "127.0.0.1"
	pgPort         = 5432
	pgDatabase     = "postgres"
	pgUsername     = "postgres"
	pgPassword     = "postgres"
	storageType    = "sql"
	httpServerURL  = ""
)

func TestMain(m *testing.M) {
	logger.PrepareLogger(logger.Config{Level: "ERROR"})

	host := os.Getenv("TEST_HTTP_SERVER_HOST")
	if host != "" {
		httpServerHost = host
	}
	host = os.Getenv("TEST_GRPC_SERVER_HOST")
	if host != "" {
		grpcServerHost = host
	}

	port := os.Getenv("TEST_HTTP_SERVER_PORT")
	if port != "" {
		httpServerPort, _ = strconv.Atoi(port)
	}
	port = os.Getenv("TEST_GRPC_SERVER_PORT")
	if port != "" {
		grpcServerPort, _ = strconv.Atoi(port)
	}

	host = os.Getenv("TEST_POSTGRES_HOST")
	if host != "" {
		pgHost = host
	}
	port = os.Getenv("TEST_POSTGRES_PORT")
	if port != "" {
		var err error
		pgPort, err = strconv.Atoi(port)
		if err != nil {
			log.Printf("failed to parse port '%s': %v", port, err)
			os.Exit(-1)
		}
	}

	opt := os.Getenv("TEST_POSTGRES_DB")
	if opt != "" {
		pgDatabase = opt
	}
	opt = os.Getenv("TEST_POSTGRES_USERNAME")
	if host != "" {
		pgUsername = opt
	}
	opt = os.Getenv("TEST_POSTGRES_PASSWORD")
	if host != "" {
		pgPassword = opt
	}

	storage := os.Getenv("TEST_STORAGE_TYPE")
	if storage != "" {
		storageType = storage
	}

	httpServerURL = fmt.Sprintf("http://%s/", net.JoinHostPort(httpServerHost, strconv.Itoa(httpServerPort)))

	cleanupDB()
	code := m.Run()
	os.Exit(code)
}

func TestStorage(t *testing.T) {
	t.Run("add event", func(t *testing.T) {
		require.NoError(t, cleanupDB())
		event := createEvent()
		jsonStr, err := json.Marshal(event)
		require.NoError(t, err)

		resp := sendRequest(t, "POST", httpServerURL, "add", jsonStr)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "failed to read body")
		require.NotEmpty(t, string(body))
		require.Equal(t, event.ID, string(body))
	})

	t.Run("update get event", func(t *testing.T) {
		require.NoError(t, cleanupDB())
		event := createEvent()
		jsonStr, err := json.Marshal(event)

		require.NoError(t, err)

		resp := sendRequest(t, "POST", httpServerURL, "add", jsonStr)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "failed to read body")
		require.Equal(t, event.ID, string(body))

		expected := createEvent()
		expected.Title = "new title"
		jsonStr, err = json.Marshal(internalhttp.UpdateReq{
			ID:    event.ID,
			Event: expected,
		})
		require.NoError(t, err)
		updResp := sendRequest(t, "POST", httpServerURL, "update", jsonStr)
		defer updResp.Body.Close()

		require.Equal(t, 200, updResp.StatusCode)
		body, err = ioutil.ReadAll(updResp.Body)
		require.NoError(t, err, "failed to read body")
		require.Equal(t, "", string(body))
		getResp := sendRequest(
			t,
			"POST",
			httpServerURL,
			"events/day",
			[]byte(`{"date": "`+expected.StartTime.Local().Format(time.RFC3339)+`"}`),
		)
		defer getResp.Body.Close()
		require.Equal(t, 200, getResp.StatusCode)
		body, err = ioutil.ReadAll(getResp.Body)
		require.NoError(t, err, "failed to read body")
		var actual []storage.Event
		require.NoError(t, json.Unmarshal(body, &actual), "failed to parse response")
		require.Equal(t, 1, len(actual))
		compareEvents(t, expected, actual[0])
	})

	t.Run("remove event", func(t *testing.T) {
		require.NoError(t, cleanupDB())
		event := createEvent()
		jsonStr, err := json.Marshal(event)
		require.NoError(t, err)

		resp := sendRequest(t, "POST", httpServerURL, "add", jsonStr)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "failed to read body")
		require.Equal(t, event.ID, string(body))

		updResp := sendRequest(t, "POST", httpServerURL, "remove", []byte(`{"id": "`+event.ID+`"}`))
		defer updResp.Body.Close()

		require.Equal(t, 200, updResp.StatusCode)
		body, err = ioutil.ReadAll(updResp.Body)
		require.NoError(t, err, "failed to read body")
		require.Equal(t, "", string(body))

		getResp := sendRequest(
			t,
			"POST",
			httpServerURL,
			"events/day",
			[]byte(`{"date": "`+event.StartTime.Local().Format(time.RFC3339)+`"}`),
		)
		defer getResp.Body.Close()
		require.Equal(t, 200, getResp.StatusCode)
		body, err = ioutil.ReadAll(getResp.Body)
		require.NoError(t, err, "failed to read body")
		var actual []storage.Event
		require.NoError(t, json.Unmarshal(body, &actual), "failed to parse response")
		require.Equal(t, 0, len(actual))
	})
}

func TestGetEvents(t *testing.T) {
	require.NoError(t, cleanupDB())
	initDate := time.Date(2300, 0o1, 0o1, 0, 0, 0, 0, time.UTC)
	event := createEvent()
	event.StartTime = initDate
	event.EndTime = initDate.Add(2 * time.Hour)
	events := make([]storage.Event, 0, 60)

	for i := 0; i < 60; i++ {
		jsonStr, err := json.Marshal(event)
		require.NoError(t, err)
		resp := sendRequest(t, "POST", httpServerURL, "add", jsonStr)
		require.Equal(t, 200, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		require.NoError(t, err, "failed to read body")
		require.Equal(t, event.ID, string(body))
		events = append(events, event)
		if i < 9 {
			event.ID = fmt.Sprint("c96506d7-22c4-4b11-bfdb-76faabde3b1" + strconv.Itoa(i+1))
		} else {
			event.ID = fmt.Sprint("c96506d7-22c4-4b11-bfdb-76faabde3c" + strconv.Itoa(i+1))
		}
		event.Title += strconv.Itoa(i)
		event.StartTime = event.StartTime.AddDate(0, 0, 1)
		event.EndTime = event.EndTime.AddDate(0, 0, 1)
	}

	t.Run("get day", func(t *testing.T) {
		resp := sendRequest(
			t,
			"POST",
			httpServerURL,
			"events/day",
			[]byte(`{"date": "`+initDate.Local().Format(time.RFC3339)+`"}`),
		)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "failed to read body")

		var actual []storage.Event
		require.NoError(t, json.Unmarshal(body, &actual), "failed to parse response")
		require.Equal(t, 1, len(actual))
		compareEvents(t, events[0], actual[0])
	})

	t.Run("get week", func(t *testing.T) {
		resp := sendRequest(
			t,
			"POST",
			httpServerURL,
			"events/week",
			[]byte(`{"date": "`+initDate.Local().Format(time.RFC3339)+`"}`),
		)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "failed to read body")

		var actual []storage.Event
		require.NoError(t, json.Unmarshal(body, &actual), "failed to parse response")
		require.Equal(t, 7, len(actual))
		sort.Slice(actual, func(i, j int) bool {
			return actual[i].StartTime.Before(actual[j].StartTime)
		})

		for i := 0; i < 7; i++ {
			compareEvents(t, events[i], actual[i])
		}
	})

	t.Run("get month", func(t *testing.T) {
		resp := sendRequest(
			t,
			"POST",
			httpServerURL,
			"events/month",
			[]byte(`{"date": "`+initDate.Local().Format(time.RFC3339)+`"}`),
		)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "failed to read body")

		var actual []storage.Event
		require.NoError(t, json.Unmarshal(body, &actual), "failed to parse response")
		require.Equal(t, 31, len(actual))
		sort.Slice(actual, func(i, j int) bool {
			return actual[i].StartTime.Before(actual[j].StartTime)
		})

		for i := 0; i < 31; i++ {
			compareEvents(t, events[i], actual[i])
		}
	})

	t.Run("get month 28 days", func(t *testing.T) {
		resp := sendRequest(
			t,
			"POST",
			httpServerURL,
			"events/month",
			[]byte(`{"date": "`+initDate.AddDate(0, 1, 0).Local().Format(time.RFC3339)+`"}`),
		)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "failed to read body")

		var actual []storage.Event
		require.NoError(t, json.Unmarshal(body, &actual), "failed to parse response")
		require.Equal(t, 28, len(actual))
		sort.Slice(actual, func(i, j int) bool {
			return actual[i].StartTime.Before(actual[j].StartTime)
		})

		for i := 0; i < 28; i++ {
			compareEvents(t, events[i+31], actual[i])
		}
	})
}

func TestErrors(t *testing.T) {
	t.Run("add no event", func(t *testing.T) {
		resp := sendRequest(t, "POST", httpServerURL, "add", []byte(`{}`))
		defer resp.Body.Close()
		require.Equal(t, 500, resp.StatusCode)
	})

	t.Run("remove non exists event", func(t *testing.T) {
		resp := sendRequest(t, "POST", httpServerURL, "remove", []byte(`{"id": "_non_exists_"}`))
		defer resp.Body.Close()
		require.Equal(t, 500, resp.StatusCode)
	})

	t.Run("update non exists event", func(t *testing.T) {
		event := createEvent()
		jsonStr, err := json.Marshal(event)
		require.NoError(t, err)

		resp := sendRequest(t, "POST", httpServerURL, "update", jsonStr)
		defer resp.Body.Close()
		require.Equal(t, 500, resp.StatusCode)
	})
}

func sendRequest(t *testing.T, method string, url string, path string, requestBody []byte) *http.Response {
	t.Helper()
	req, err := http.NewRequestWithContext(
		context.Background(),
		method,
		url+path,
		bytes.NewBuffer(requestBody),
	)
	require.NoError(t, err, "failed to send request")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		require.NoError(t, err, "failed send request")
	}
	return resp
}

func createEvent() storage.Event {
	return storage.Event{
		ID:          "c96506d7-22c4-4b11-bfdb-76faabde3a24",
		Title:       "Test",
		StartTime:   time.Now().Truncate(time.Second).Add(5 * time.Minute),
		EndTime:     time.Now().Truncate(time.Second).Add(20 * time.Minute),
		Description: "TestDescription",
		OwnerID:     "OwnId",
	}
}

func compareEvents(t *testing.T, expected storage.Event, actual storage.Event) {
	t.Helper()
	require.True(
		t,
		expected.StartTime.Equal(actual.StartTime),
		"start time is not equals %q != %q", expected.StartTime, actual.StartTime)
	require.True(
		t,
		expected.StartTime.Equal(actual.StartTime),
		"start time is not equals %q != %q", expected.StartTime, actual.StartTime)
	expected.StartTime = actual.StartTime
	expected.EndTime = actual.EndTime
	require.Equal(t, expected, actual)
}

func cleanupDB() error {
	if storageType != "sql" {
		return nil
	}
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"sslmode=disable host=%s port=%d dbname=%s user=%s password=%s",
			pgHost,
			pgPort,
			pgDatabase,
			pgUsername,
			pgPassword,
		),
	)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("TRUNCATE TABLE Events")
	if err != nil {
		return err
	}
	_, err = db.Exec("TRUNCATE TABLE Sender_logs")

	return err
}

func TestSender(t *testing.T) {
	t.Run("add event", func(t *testing.T) {
		require.NoError(t, cleanupDB())
		event := storage.Event{
			ID:           "c96506d7-22c4-4b11-bfdb-76faabde3a99",
			Title:        "sender_check",
			StartTime:    time.Now().Add(24 * time.Hour),
			EndTime:      time.Now().Add(24*time.Hour + 2*time.Minute),
			Description:  "1234sender",
			OwnerID:      "sender",
			NotifyBefore: 1,
		}
		jsonStr, err := json.Marshal(event)
		require.NoError(t, err)

		resp := sendRequest(t, "POST", httpServerURL, "add", jsonStr)
		defer resp.Body.Close()
		require.Equal(t, 200, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "failed to read body")
		require.NotEmpty(t, string(body))
		require.Equal(t, event.ID, string(body))

		message, err := waitForEventSenderLogDB(event.ID)
		require.NoError(t, err)
		require.Equal(t, event.ID, message.ID)
		require.Equal(t, event.OwnerID, message.OwnerID)
		cleanupDB()
	})
}

func waitForEventSenderLogDB(id string) (rabbit.Message, error) {
	if storageType != "sql" {
		return rabbit.Message{}, nil
	}
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"sslmode=disable host=%s port=%d dbname=%s user=%s password=%s",
			pgHost,
			pgPort,
			pgDatabase,
			pgUsername,
			pgPassword,
		),
	)
	if err != nil {
		return rabbit.Message{}, nil
	}
	defer db.Close()
	logMessage := make([]rabbit.Message, 0)

	ticker := time.NewTicker(time.Second)
	c1 := make(chan rabbit.Message)
	go func(logMessage []rabbit.Message) {
		for range ticker.C {
			err = db.SelectContext(
				context.Background(),
				&logMessage,
				"SELECT id, name, time, owner_id AS ownerId "+
					"FROM Sender_logs WHERE id = $1",
				id,
			)
			switch len(logMessage) {
			case 0:
				continue
			default:
				ticker.Stop()
				c1 <- logMessage[0]
				break
			}
		}
	}(logMessage)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	for {
		select {
		case msg := <-c1:
			return msg, nil
		case <-ctx.Done():
			cancel()
			return rabbit.Message{}, errors.New("sender didn't send a message")
		}
	}

	return rabbit.Message{}, err
}
