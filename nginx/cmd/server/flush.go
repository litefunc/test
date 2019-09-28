package main

import (
	"fmt"
	"net/http"
	"time"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Step 1\n")) // 浏览器不会立即输出Step1
	time.Sleep(5 * time.Second) // 转圈5秒
	w.Write([]byte("Step 2"))   // Step1和Step2同时输出
}

func testFlushHandler(w http.ResponseWriter, r *http.Request) {
	b := "Step 1. 立即响应client, handler继续处理"
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b))) // 要设置Content-Length否则客户端可能会认为返回的数据不完整
	w.Write([]byte(b))
	flushBody(w) // 立即flush并断开连接
	time.Sleep(5 * time.Second)
	fmt.Println("Do something...") // 之后处理自己的业务
}

func testFlushHandler1(w http.ResponseWriter, r *http.Request) {
	b := "Step 1. 立即响应client, handler继续处理"
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b))) // 要设置Content-Length否则客户端可能会认为返回的数据不完整
	w.Write([]byte(b))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Do something...") // 之后处理自己的业务
	b = "Step 2. 立即响应client, handler继续处理"
	w.Write([]byte(b))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func flushBody(w http.ResponseWriter) bool {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
		if hj, ok := w.(http.Hijacker); ok { // 从ResponseWriter获取链接控制权
			if conn, _, err := hj.Hijack(); err == nil {
				if err := conn.Close(); err == nil {
					return true
				}
			}
		}
	}
	return false
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", testHandler)         // 测试不立即响应
	mux.HandleFunc("/flush", testFlushHandler)   // 测试立即响应后同步处理业务
	mux.HandleFunc("/flush1", testFlushHandler1) // 测试立即响应后同步处理业务
	http.ListenAndServe(":8080", mux)
}
