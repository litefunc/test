package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {

	for i := 0; i < 100; i++ {
		go func(i int) {
			cmd := exec.Command("./script.sh")
			cmd.Stdout = os.Stdout
			err := cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			if err := cmd.Wait(); err != nil {
				log.Fatal(err)
			}

			log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
			log.Println(i)
		}(i)
	}

	time.Sleep(time.Second * 600)
}
