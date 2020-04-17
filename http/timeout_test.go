package http

import (
	"VodoPlay/logger"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientTimeout(t *testing.T) {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := map[string]interface{}{
			"id":    "12",
			"scope": "test-scope",
		}

		time.Sleep(100 * time.Millisecond) //<- Any value > 20ms
		b, err := json.Marshal(d)
		if err != nil {
			t.Error(err)
		}
		io.WriteString(w, string(b))
		w.WriteHeader(http.StatusOK)
	})

	backend := httptest.NewServer(http.TimeoutHandler(handlerFunc, 20*time.Millisecond, "server timeout"))

	url := backend.URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Error("Request error", err)
		return
	}

	netClient := &http.Client{
		Timeout: 10 * time.Millisecond,
	}

	resp, err := netClient.Do(req)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			logger.Debug(true)
		}
		t.Error("Response error", err)
		return
	}
	defer resp.Body.Close()
	logger.Debug(resp.StatusCode)
}
