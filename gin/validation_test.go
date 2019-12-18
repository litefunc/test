package gin

import (
	"MediaImage/test/request"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

type va struct {
	A string `form:"a" json:"a" binding:"required"`
	B string `form:"b" json:"b" binding:"required"`
	C string `form:"c" json:"c" `
}

func (md va) m(exclude ...string) map[string]interface{} {
	m := map[string]interface{}{
		"a": md.A,
		"b": md.B,
		"c": md.C,
		"d": "d",
	}
	for _, k := range exclude {
		delete(m, k)
	}
	return m
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/query", func(c *gin.Context) {
		var q va
		if err := c.ShouldBind(&q); err != nil {
			c.JSON(http.StatusOK, err.Error())
			return
		}

		c.JSON(http.StatusOK, q)
	})

	router.POST("/json", func(c *gin.Context) {
		var json va
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusOK, err.Error())
			return
		}

		c.JSON(http.StatusOK, json)
	})

	router.POST("/upload", func(c *gin.Context) {

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusOK, err.Error())
			return
		}

		c.String(http.StatusOK, file.Filename)
	})

	return router
}

func TestQueryString(t *testing.T) {

	router := setupRouter()

	for _, v := range []struct {
		in   map[string]interface{}
		want interface{}
	}{
		{va{"a", "b", "c"}.m(), va{"a", "b", "c"}},
		{va{"a", "b", ""}.m(), va{"a", "b", ""}},
		{va{"", "", ""}.m(), "Key: 'va.A' Error:Field validation for 'A' failed on the 'required' tag\nKey: 'va.B' Error:Field validation for 'B' failed on the 'required' tag"},
	} {
		w := httptest.NewRecorder()
		url := fmt.Sprintf(`/query?a=%s&b=%s&c=%s`, v.in["a"], v.in["b"], v.in["c"])
		req := request.Get(t, url, nil)
		router.ServeHTTP(w, req)
		if err := equal(w, v.want); err != nil {
			t.Error(err)
		}
	}
}

func TestPostValidation(t *testing.T) {

	router := setupRouter()

	for _, v := range []struct {
		in   map[string]interface{}
		want interface{}
	}{
		{va{"a", "b", "c"}.m(), va{"a", "b", "c"}},
		{va{"a", "b", ""}.m(), va{"a", "b", ""}},
		{va{"", "", ""}.m(), "Key: 'va.A' Error:Field validation for 'A' failed on the 'required' tag\nKey: 'va.B' Error:Field validation for 'B' failed on the 'required' tag"},
	} {
		w := httptest.NewRecorder()
		req := request.Post(t, "/json", nil, v.in)
		router.ServeHTTP(w, req)
		if err := equal(w, v.want); err != nil {
			t.Error(err)
		}
	}
}

func TestFile(t *testing.T) {

	router := setupRouter()

	for _, v := range []struct {
		in   map[string]string
		file string
		want interface{}
	}{
		{nil, "", http.ErrMissingFile.Error()},
	} {
		w := httptest.NewRecorder()
		var ff request.FormFile
		if v.file != "" {
			ff = request.FormFile{"file", v.file}
		}

		req := request.PostFormFile(t, "/upload", nil, v.in, ff)
		router.ServeHTTP(w, req)
		if err := equal(w, v.want); err != nil {
			t.Error(err)
		}
	}
}

func equal(w *httptest.ResponseRecorder, want interface{}) error {

	by1, _ := jsoniter.Marshal(want)
	var o interface{}
	json.Unmarshal(by1, &o)

	by2, err := ioutil.ReadAll(w.Body)
	if err != nil {
		return err
	}

	var got interface{}
	if err := json.Unmarshal(by2, &got); err != nil {
		return err
	}

	if !reflect.DeepEqual(o, got) {
		return fmt.Errorf(`want:%s, got:%s`, string(by1), string(by2))
	}
	return nil

}
