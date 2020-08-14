package main

import (
	"cloud/server/ota/service/ui/asset/src/backend/logger"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

// func main() {
// 	url := "http://localhost:8090/"
// 	socksAddress := "127.0.0.1:8081"

// 	socks, err := proxy.SOCKS5("tcp", socksAddress, nil, &net.Dialer{
// 		Timeout:   30 * time.Second,
// 		KeepAlive: 30 * time.Second,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	client := &http.Client{
// 		Transport: &http.Transport{
// 			Dial:                socks.Dial,
// 			TLSHandshakeTimeout: 10 * time.Second,
// 		},
// 	}

// 	res, err := client.Get(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	content, err := ioutil.ReadAll(res.Body)
// 	res.Body.Close()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%s", string(content))
// }

// func main() {

// 	proxyUrl := "http://127.0.0.1:8081"

// 	u, _ := url.Parse(proxyUrl)

// 	client := &http.Client{
// 		Transport: &http.Transport{
// 			Proxy:               http.ProxyURL(u),
// 			TLSHandshakeTimeout: 10 * time.Second,
// 		},
// 	}

// 	res, err := client.Get("http://localhost:8090")
// 	if err != nil {
// 		panic(err)
// 	}
// 	content, err := ioutil.ReadAll(res.Body)
// 	res.Body.Close()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%s", string(content))
// }

// func DialCustom(network, address string, timeout time.Duration, localIP []byte, localPort int) (net.Conn, error) {
// 	netAddr := &net.TCPAddr{Port: localPort}

// 	if len(localIP) != 0 {
// 		netAddr.IP = localIP
// 	}

// 	fmt.Println("netAddr:", netAddr)

// 	d := net.Dialer{Timeout: timeout, LocalAddr: netAddr}

// 	return d.Dial(network, address)
// }

func DialCustom(network, address string, localIP []byte, localPort int) *net.Dialer {
	netAddr := &net.TCPAddr{Port: localPort}

	if len(localIP) != 0 {
		netAddr.IP = localIP
	}

	fmt.Println("netAddr:", netAddr)

	d := net.Dialer{LocalAddr: netAddr}

	return &d
}

func main() {

	serverAddr := "127.0.0.1:8090"

	// 172.28.0.180
	//localIP := []byte{0xAC, 0x1C, 0, 0xB4}  // 指定IP
	localIP := []byte{} //  any IP，不指定IP
	localPort := 9001   // 指定端口
	// conn, err := DialCustom("tcp", serverAddr, time.Second*10, localIP, localPort)
	// if err != nil {
	// 	fmt.Println("dial failed:", err)
	// 	os.Exit(1)
	// }
	// defer conn.Close()

	conn := DialCustom("tcp", serverAddr, localIP, localPort)

	client := &http.Client{
		Transport: &http.Transport{
			Dial: conn.Dial,
		},
	}
	defer client.CloseIdleConnections()

	res, err := client.Get("http://localhost:8090")
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
	fmt.Printf("%s", string(content))

	// buffer := make([]byte, 512)
	// reader := bufio.NewReader(conn)

	// n, err2 := reader.Read(buffer)
	// if err2 != nil {
	// 	fmt.Println("Read failed:", err2)
	// 	return
	// }

	// fmt.Println("count:", n, "msg:", string(buffer))

	// select {}
}
