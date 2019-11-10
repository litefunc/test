package main

import (
	"cloud/lib/logger"

	"github.com/gin-gonic/gin"
)

type handler struct {
	n int
}

func (h handler) GetN(c *gin.Context) {
	logger.Debug(h.n)
	c.JSON(200, h.n)
}

func (h *handler) GetN1(c *gin.Context) {
	logger.Debug(h.n)
	c.JSON(200, h.n)
}

func (h *handler) AddN(c *gin.Context) {
	h.n++
	logger.Debug(h.n)
	c.JSON(200, h.n)
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	var h handler
	r.GET("/get", h.GetN)
	r.GET("/get1", h.GetN)
	r.GET("/add", h.AddN)

	h1 := &handler{}
	r.GET("/get/p", h1.GetN)
	r.GET("/get1/p", h1.GetN)
	r.GET("/add/p", h1.AddN)

	r.Run() // listen and serve on 0.0.0.0:8080
}
