package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(url string, dst io.Writer) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(dst, resp.Body)
	return err
}

type writer struct{}

func (rec writer) Write(p []byte) (n int, err error) {
	n = len(p)
	p = p[:0]
	return
}

func main() {
	reader := strings.NewReader("ab cde fghij k l  m")
	s := make([]byte, 5)
	var writer bytes.Buffer

	for {
		n, err := reader.Read(s)
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		str := string(s[:n])
		fmt.Println(str)
		wr, err := writer.Write([]byte(str))
		fmt.Println(wr) // 寫了多少 byte
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}
	fmt.Println(writer.String())
}
