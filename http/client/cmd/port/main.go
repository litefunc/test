package main

import (
	"cloud/server/ota/service/ui/asset/src/backend/logger"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type connDialer struct {
	c net.Conn
}

func (cd connDialer) Dial(network, addr string) (net.Conn, error) {
	return cd.c, nil
}

type dialer struct {
	src *net.TCPAddr
}

func (cd dialer) Dial(network, addr string) (net.Conn, error) {
	dst, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	conn, err := net.DialTCP(network, cd.src, dst)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if err := conn.SetLinger(0); err != nil {
		logger.Error(err)
	}
	return conn, nil
}

func dial() {
	localPort := 9001
	dialer := dialer{src: &net.TCPAddr{Port: localPort}}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}
	defer client.CloseIdleConnections()
	for i := 0; i < 10; i++ {
		get(client, "http://localhost:8090")
	}
	time.Sleep(time.Second * 10)
}

func dial1() {

	localPort := 9001
	netAddr := &net.TCPAddr{Port: localPort}

	RemoteEP := net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8090}
	conn, err := net.DialTCP("tcp", netAddr, &RemoteEP)
	if err != nil {
		logger.Error(err)
		return
	}
	conn.SetLinger(0)
	client := &http.Client{
		Transport: &http.Transport{
			Dial: connDialer{conn}.Dial,
		},
	}
	defer client.CloseIdleConnections()

	for i := 0; i < 10; i++ {
		get(client, "http://localhost:8090")
	}
	time.Sleep(time.Second * 10)
}

func main() {
	dial()
}

func get(cli *http.Client, url string) {
	res, err := cli.Get(url)
	if err != nil {
		logger.Error(err)
		return
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	fmt.Println(string(content))
}
