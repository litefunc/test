package main

import (
	"log"
	"net"
	"time"
)

const (
	srcAddr = "224.100.100.3:5001"
	srvAddr = "224.100.100.3:5000"
)

//send multicast data
func ping(srcAddr, a string) {

	src, err := net.ResolveUDPAddr("udp", srcAddr)
	if err != nil {
		log.Fatal(err)
	}

	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", src, addr)
	for {
		c.Write([]byte("hello, world\n"))
		time.Sleep(1 * time.Second)
	}
}

func main() {
	ping(srcAddr, srvAddr)
}
