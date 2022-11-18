package internalhttp

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func mockApp() *app.App {
	testStorage := memorystorage.New()
	appForTest := app.New(testStorage)

	return appForTest
}

func TestServer_AddEvent(t *testing.T) {
	body := []byte(`{
    "id":"123",
    "title":"dadadasfasfaf",
    "startTime": "2099-01-02T15:04:05Z",
    "endTime": "2099-01-03T15:04:05Z",
    "description":"md509201@mail.ru",
    "ownerID":"Москва",
    "notifyBefore": 11
}`)
	req := http.Request{Body: io.NopCloser(bytes.NewReader(body))}
	type fields struct {
		app *app.App
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		expectedMessage string
	}{
		{
			name: "test success",
			fields: fields{
				app: mockApp(),
			},
			args: args{
				r: &req,
			},
			expectedMessage: "123",
		},
		{
			name: "test failed",
			fields: fields{
				app: mockApp(),
			},
			args: args{
				r: &http.Request{Body: io.NopCloser(bytes.NewReader([]byte(`{"startTime": "2002-01-02T15:04:05Z",
"endTime": "2001-01-02T15:04:05Z"}`)))},
			},
			expectedMessage: "start time of the event must be in the future: incorrect event time",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				app: tt.fields.app,
			}
			resp := httptest.NewRecorder()
			s.AddEvent(resp, tt.args.r)
			require.Equal(t, tt.expectedMessage, resp.Body.String())
		})
	}
}

func TestServer_UpdateEvent(t *testing.T) {
	body := []byte(`{
    "id":"123",
    "title":"dadadasfasfaf",
    "startTime": "2099-01-02T15:04:05Z",
    "endTime": "2099-01-03T15:04:05Z",
    "description":"md509201@main.ru",
    "ownerID":"Москва",
    "notifyBefore": 11
}`)
	successBody := []byte(`{
    "id":"123",
	"event": {
    	"id":"123",
    	"title":"dadadasfasfaf",
    	"startTime": "2099-01-02T15:04:05Z",
    	"endTime": "2099-01-03T15:04:05Z",
    	"description":"md509201@main.ru",
    	"ownerID":"Москва",
    	"notifyBefore": 11
	}
}`)
	successReq := http.Request{Body: io.NopCloser(bytes.NewReader(successBody))}
	req := http.Request{Body: io.NopCloser(bytes.NewReader(body))}
	type fields struct {
		srv  *http.Server
		addr string
		app  *app.App
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		expectedMessage string
	}{
		{
			name: "test event not found",
			fields: fields{
				app: mockApp(),
			},
			args: args{
				r: &successReq,
			},
			expectedMessage: "failed to update event with id \"123\": event not found",
		},
		{
			name: "test failed",
			fields: fields{
				app: mockApp(),
			},
			args: args{
				r: &req,
			},
			expectedMessage: "start time of the event must be in the future: incorrect event time",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				srv:  tt.fields.srv,
				addr: tt.fields.addr,
				app:  tt.fields.app,
			}
			resp := httptest.NewRecorder()
			s.UpdateEvent(resp, tt.args.r)
			require.Equal(t, tt.expectedMessage, resp.Body.String())
		})
	}
}

func TestServer_RemoveEvent(t *testing.T) {
	type fields struct {
		srv  *http.Server
		addr string
		app  *app.App
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		expectedMessage string
	}{
		{
			name: "test failed",
			fields: fields{
				app: mockApp(),
			},
			args: args{
				r: &http.Request{Body: io.NopCloser(bytes.NewReader([]byte(`{"id":"123"}`)))},
			},
			expectedMessage: "failed to remove event with id \"123\": event not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				srv:  tt.fields.srv,
				addr: tt.fields.addr,
				app:  tt.fields.app,
			}
			resp := httptest.NewRecorder()
			s.RemoveEvent(resp, tt.args.r)
			require.Equal(t, tt.expectedMessage, resp.Body.String())
		})
	}
}

func TestServer_GetEvents(t *testing.T) {
	type fields struct {
		app *app.App
	}
	type args struct {
		r      *http.Request
		period string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		expectedMessage string
	}{
		{
			name: "test success",
			fields: fields{
				app: mockApp(),
			},
			args: args{
				r: &http.Request{Body: io.NopCloser(bytes.NewReader([]byte(`{"date":"2019-01-03T15:04:05Z"}`)))},
			},
			expectedMessage: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				app: tt.fields.app,
			}
			resp := httptest.NewRecorder()
			s.GetEvents(resp, tt.args.r, tt.args.period)
			require.Equal(t, tt.expectedMessage, resp.Body.String())
		})
	}
}
