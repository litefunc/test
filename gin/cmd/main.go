package main

import (
	"cloud/lib/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

type handlers []http.Handler

func (hds handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, hd := range hds {
		hd.ServeHTTP(w, r)
	}
}

func routerA() *gin.Engine {
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

	return r
}

func routerB() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/b", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "bpong",
		})
	})

	return r
}

func main() {

	hds := handlers{routerA()}
	// r.Run() // listen and serve on 0.0.0.0:8080

	addr := fmt.Sprintf(`:%d`, 8080)
	srv := &http.Server{
		Addr:    addr,
		Handler: hds,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	} else {
		logger.Warn("Server Shutdown succeed")
		cancel()
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		logger.Warn("timeout of 5 seconds.")
	}
	logger.Warn("Server exiting")

}
