package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"mstore/logger"
	"net/http"
	"os"
	"path"
)

func main() {
	log.SetFlags(log.Lshortfile)

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}

	println(string(buf[:n]))

	get("https://www.google.com/")
	get("https://127.0.0.1:443/hello")
	getTlsInsecure("https://127.0.0.1:443/hello")
	getTls("https://127.0.0.1:443/hello", path.Join(os.Getenv("GOPATH"), "/src/test/openssl/static/server.crt"))

	// getTls("https://127.0.0.1:40000/token/new", path.Join(os.Getenv("GOPATH"), "/src/test/openssl/static/server.crt"))
	// getTls("http://127.0.0.1:40000/token/new", path.Join(os.Getenv("GOPATH"), "/src/MediaImage/assets/ssl/server.crt"))
	// getTls("https://127.0.0.1:40000/token/new", path.Join(os.Getenv("GOPATH"), "/src/MediaImage/assets/ssl/server.crt"))
	// getTls("https://192.168.2.2:40001/token/new", path.Join(os.Getenv("GOPATH"), "/src/MediaImage/assets/ssl/server.crt"))
	// getTls("https://192.168.2.2:40001/vod/list", path.Join(os.Getenv("GOPATH"), "/src/MediaImage/assets/ssl/server.crt"))
}

func get(url string) {

	cli := http.Client{}
	resp, err := cli.Get(url)
	if err != nil {
		logger.Error(err)
		return
	}
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(string(by))
}

func getTlsInsecure(url string) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: conf}
	cli := http.Client{Transport: tr}
	resp, err := cli.Get(url)
	if err != nil {
		logger.Error(err)
		return
	}
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(string(by))
}

func getTls(url, crt string) {

	conf := &tls.Config{
		RootCAs: loadCA(crt),
	}
	tr := &http.Transport{TLSClientConfig: conf}
	cli := http.Client{Transport: tr}
	resp, err := cli.Get(url)
	if err != nil {
		logger.Error(err)
		return
	}
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(string(by))
}

func loadCA(caFile string) *x509.CertPool {
	pool := x509.NewCertPool()

	if ca, e := ioutil.ReadFile(caFile); e != nil {
		log.Fatal("ReadFile: ", e)
	} else {
		pool.AppendCertsFromPEM(ca)
	}
	return pool
}
