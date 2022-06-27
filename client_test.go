package webclient_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fernandoocampo/webclient"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultClient(t *testing.T) {
	expectedHttpClient := http.Client{}
	got := webclient.New("http://localhost:8999")
	assert.NotNil(t, got)
	assert.Equal(t, &expectedHttpClient, got.HTTPClient())
}

func TestGetRequest(t *testing.T) {
	cases := map[string]struct {
		path      string
		timeoutms int
		isError   error
		handler   *handlerMock
		want      *webclient.Response
	}{
		"success_with_data": {
			path:    "/people",
			handler: newHandlerMock(200, `[{name: "vane", age:  43}]`),
			want: &webclient.Response{
				StatusCode: 200,
				Data:       []byte(`[{name: "vane", age:  43}]`),
			},
		},
		"success_without_data": {
			path:    "/people",
			handler: newHandlerMock(200, ""),
			want: &webclient.Response{
				StatusCode: 200,
				Data:       []byte(""),
			},
		},
		"request_with_timeout": {
			path:      "/people",
			timeoutms: 200,
			handler:   newHandlerMockWithTimeout(100, 200, ""),
			isError:   errors.New("context deadline exceeded"),
		},
	}
	for name, data := range cases {
		t.Run(name, func(st *testing.T) {
			mockServer := httptest.NewServer(data.handler)
			defer mockServer.Close()
			ctx := context.TODO()
			var cancel context.CancelFunc
			if data.timeoutms != 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(data.timeoutms)*time.Millisecond)
				defer cancel()
			}

			client := webclient.New(mockServer.URL)
			got, err := client.Get(ctx, data.path)
			if data.isError != nil {
				assert.Error(st, err)
				assert.Contains(st, err.Error(), data.isError.Error())
			} else {
				assert.NoError(st, err)
			}
			assert.Equal(st, data.want, got)
		})
	}
}

type handlerMock struct {
	resp    []byte
	code    int
	sleepms int
}

func newHandlerMock(code int, resp string) *handlerMock {
	var respdata []byte
	if resp != "" {
		respdata = []byte(resp)
	}
	return &handlerMock{
		code: code,
		resp: respdata,
	}
}

func newHandlerMockWithTimeout(code, sleepms int, resp string) *handlerMock {
	var respdata []byte
	if resp != "" {
		respdata = []byte(resp)
	}
	return &handlerMock{
		code:    code,
		resp:    respdata,
		sleepms: sleepms,
	}
}

func (w *handlerMock) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if w.sleepms != 0 {
		time.Sleep(time.Duration(w.sleepms) * time.Millisecond)
	}
	rw.WriteHeader(w.code)
	rw.Write(w.resp)
}
