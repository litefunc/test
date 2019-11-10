package http

import (
	"cloud/lib/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

type hd struct {
	N int `json:"N"`
	n int `json:"n"`
}

func (h *hd) PointerGet(w http.ResponseWriter, r *http.Request) {

	logger.Debug(h)
	by, err := json.Marshal(h)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(string(by))
	w.Write(by)
}

func (h *hd) PointerAdd(w http.ResponseWriter, r *http.Request) {
	h.N++
	h.n++
	logger.Debug(h.n)
}

func (h hd) ValueAdd(w http.ResponseWriter, r *http.Request) {
	h.N++
	h.n++
	logger.Debug(h.n)
}

func TestPointerReceiver(t *testing.T) {

	h := new(hd)

	http.HandleFunc("/get/pointer", h.PointerGet)
	http.HandleFunc("/add/pointer", h.PointerAdd)
	http.HandleFunc("/add/value", h.ValueAdd)

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Second)

	_, err := http.Get("http://localhost:8080/add/pointer")
	if err != nil {
		t.Error(err)
	}

	if want := 1; h.n != want {
		t.Errorf(`want:%d, got:%d`, want, h.n)
	}

	_, err = http.Get("http://localhost:8080/add/value")
	if err != nil {
		t.Error(err)
	}

	if want := 1; h.n != want {
		t.Errorf(`want:%d, got:%d`, want, h.n)
	}

	resp, err := http.Get("http://localhost:8080/get/pointer")
	if err != nil {
		t.Error(err)
	}
	if err := getn(resp, hd{1, 0}); err != nil {
		t.Error(err)
	}

}

func getn(resp *http.Response, want hd) error {
	var h hd
	if err := json.NewDecoder(resp.Body).Decode(&h); err != nil {
		return err
	}
	if h != want {
		return fmt.Errorf(`want:%+v, got:%+v`, want, h)
	}
	return nil
}
