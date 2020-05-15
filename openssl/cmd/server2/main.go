package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"mstore/logger"
	"net"
	"net/http"
	"os"
	"path"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func main() {

	p := flag.Int("p", 80, "port")

	gopath := os.Getenv("GOPATH")
	cert := flag.String("cert", path.Join(gopath, "/src/test/openssl/static/server.crt"), "cert")
	key := flag.String("key", path.Join(gopath, "/src/test/openssl/static/server.key"), "key")

	flag.Parse()
	port := fmt.Sprintf(`:%d`, *p)

	logger.Debug(port)
	logger.Debug(*cert)
	logger.Debug(*key)

	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(port, *cert, *key, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// func main() {
// 	log.SetFlags(log.Lshortfile)

// 	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	config := &tls.Config{Certificates: []tls.Certificate{cer}}
// 	ln, err := tls.Listen("tcp", ":443", config)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer ln.Close()

// 	for {
// 		conn, err := ln.Accept()
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		go handleConnection(conn)
// 	}
// }

func handleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		println(msg)

		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}
