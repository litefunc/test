package main

import (
	"VodoPlay/logger"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// An artificial input source.
	// const input = "Now is the winter of our discontent,\n\nMade glorious summer by this sun of York . abc\n\n"
	// scanner := bufio.NewScanner(strings.NewReader(input))

	f, _ := os.Open("a.txt")
	scanner := bufio.NewScanner(f)
	// Set the split function for the scanning operation.
	// scanner.Split(bufio.ScanLines)
	// Count the words.
	count := 0
	for scanner.Scan() {
		var m map[int]int
		if err := json.Unmarshal(scanner.Bytes(), &m); err != nil {
			logger.Error(err)
		}
		fmt.Printf("%d %s\n", count, scanner.Text())
		logger.Debug(count, m)
		count++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	fmt.Printf("%d\n", count)
}
