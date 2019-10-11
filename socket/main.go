package main

import (
	"cloud/lib/logger"
	"fmt"
	"net"
	"os"
	"test/socket/client"
	"time"
)

const address = "localhost:1024"

func main() {

	for i := 1; i <= 10; i++ {
		go func(i int) {
			time.Sleep(time.Second)
			cli := client.NewClient(address)
			conn := cli.Conn()

			cli.Write(conn, fmt.Sprintf(`"hello %d %d"`, i, 1))
			cli.Write(conn, fmt.Sprintf(`"hello %d %d"`, i, 2))
		}(i)
	}

	//建立socket，監聽埠
	netListen, err := net.Listen("tcp", address)
	CheckError(err)
	defer netListen.Close()

	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			logger.Error(err)
			continue
		}

		Log(conn.RemoteAddr().String(), " tcp connect success")
		handleConnection(conn)
	}

}

//處理連線
func handleConnection(conn net.Conn) {

	buffer := make([]byte, 2048)

	for {

		n, err := conn.Read(buffer)

		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		Log(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))

	}

}
func Log(v ...interface{}) {
	logger.Info(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
