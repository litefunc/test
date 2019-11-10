package gin

import (
	"cloud/lib/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type hd struct {
	N int `json:"N"`
	n int `json:"n"`
}

func (h *hd) PointerGet(c *gin.Context) {

	logger.Debug(h)
	c.JSON(200, h)
}

func (h *hd) PointerAdd(c *gin.Context) {
	h.N++
	h.n++
	logger.Debug(h.n)
}

func (h hd) ValueAdd(c *gin.Context) {
	h.N++
	h.n++
	logger.Debug(h.n)
}

func TestPointerReceiver(t *testing.T) {

	h := new(hd)

	r := gin.New()
	r.GET("/get/pointer", h.PointerGet)
	r.GET("/add/pointer", h.PointerAdd)
	r.GET("/add/value", h.ValueAdd)

	go func() {
		r.Run(":8080")
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
