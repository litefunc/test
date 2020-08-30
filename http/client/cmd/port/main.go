package main

import (
	"cloud/server/ota/service/ui/asset/src/backend/logger"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func DialCustom(network, address string, localIP []byte, localPort int) *net.Dialer {
	netAddr := &net.TCPAddr{Port: localPort}

	if len(localIP) != 0 {
		netAddr.IP = localIP
	}

	fmt.Println("netAddr:", netAddr)

	d := net.Dialer{LocalAddr: netAddr}

	return &d
}

type connDialer struct {
	c net.Conn
}

func (cd connDialer) Dial(network, addr string) (net.Conn, error) {
	return cd.c, nil
}

func main() {

	localIP := []byte{} //  any IP，不指定IP
	localPort := 9001   // 指定端口

	netAddr := &net.TCPAddr{Port: localPort}
	if len(localIP) != 0 {
		netAddr.IP = localIP
	}

	fmt.Println("netAddr:", netAddr)

	// dialer := net.Dialer{LocalAddr: netAddr, Timeout: time.Second}

	RemoteEP := net.TCPAddr{Port: 8090}
	// RemoteEP := net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8090}
	conn, err := net.DialTCP("tcp", netAddr, &RemoteEP)
	// conn, err := dialer.Dial("tcp", "localhost:8090")
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
	// if err := conn.Close(); err != nil {
	// 	logger.Error(err)
	// }
	time.Sleep(time.Second * 10)
	// for i := 0; i < 10; i++ {
	// 	get(client, "http://localhost:8090")
	// }

	// client.CloseIdleConnections()

	// time.Sleep(time.Second * 60)
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
