package client

import (
	"cloud/lib/logger"
	"fmt"
	"net"
	"os"
)

type Client struct {
	target  string
	tcpAddr *net.TCPAddr
}

func NewClient(target string) *Client {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	return &Client{target, tcpAddr}
}

func (cli Client) Conn() *net.TCPConn {
	conn, err := net.DialTCP("tcp", nil, cli.tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	return conn
}

func (cli Client) Write(conn *net.TCPConn, s string) int {
	n, err := conn.Write([]byte(s))
	if err != nil {
		logger.Error(err)
	}
	logger.Debug("client write:", s)
	return n
}
