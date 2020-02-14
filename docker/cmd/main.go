package main

import (
	"cloud/lib/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func cmd(args ...string) (string, error) {
	out, err := exec.Command("docker", args...).CombinedOutput()
	logger.Warn("docker", strings.Join(args, " "))
	if err != nil && len(out) != 0 {
		logger.Error(err.Error() + ": " + string(out))
		return "", errors.New(err.Error() + ": " + string(out))
	}

	return strings.TrimSuffix(string(out), "\n"), err
}

func prune() error {

	args := []string{"system", "prune", "-a", "-f"}

	// find containers which based on img
	out, err := cmd(args...)
	if err != nil {
		return err
	}
	logger.Debug(out)
	return nil
}

func main() {
	// logger.Debug(prune())

	logger.Debug(PullImage("nginx:latest"))
}

// Struct representing events returned from image pulling
type pullEvent struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Error          string `json:"error,omitempty"`
	Progress       string `json:"progress,omitempty"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

// Actual image pulling function
func PullImage(dockerImageName string) bool {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)
	resp, err := cli.ImagePull(ctx, dockerImageName, types.ImagePullOptions{})

	if err != nil {
		panic(err)
	}

	cursor := Cursor{}
	layers := make([]string, 0)
	oldIndex := len(layers)

	var event *pullEvent
	decoder := json.NewDecoder(resp)

	fmt.Printf("\n")
	cursor.hide()

	for {
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}
		logger.Debug(event)
		imageID := event.ID

		// Check if the line is one of the final two ones
		if strings.HasPrefix(event.Status, "Digest:") || strings.HasPrefix(event.Status, "Status:") {
			fmt.Printf("%s\n", event.Status)
			continue
		}

		// Check if ID has already passed once
		index := 0
		for i, v := range layers {
			if v == imageID {
				index = i + 1
				break
			}
		}

		// Move the cursor
		if index > 0 {
			diff := index - oldIndex

			if diff > 1 {
				down := diff - 1
				cursor.moveDown(down)
			} else if diff < 1 {
				up := diff*(-1) + 1
				cursor.moveUp(up)
			}

			oldIndex = index
		} else {
			layers = append(layers, event.ID)
			diff := len(layers) - oldIndex

			if diff > 1 {
				cursor.moveDown(diff) // Return to the last row
			}

			oldIndex = len(layers)
		}

		cursor.clearLine()

		if event.Status == "Pull complete" {
			fmt.Printf("%s: %s\n", event.ID, event.Status)
		} else {
			fmt.Printf("%s: %s %s\n", event.ID, event.Status, event.Progress)
		}

	}

	cursor.show()

	if strings.Contains(event.Status, fmt.Sprintf("Downloaded newer image for %s", dockerImageName)) {
		return true
	}

	return false
}

// Cursor structure that implements some methods
// for manipulating command line's cursor
type Cursor struct{}

func (cursor *Cursor) hide() {
	fmt.Printf("\033[?25l")
}

func (cursor *Cursor) show() {
	fmt.Printf("\033[?25h")
}

func (cursor *Cursor) moveUp(rows int) {
	fmt.Printf("\033[%dF", rows)
}

func (cursor *Cursor) moveDown(rows int) {
	fmt.Printf("\033[%dE", rows)
}

func (cursor *Cursor) clearLine() {
	fmt.Printf("\033[2K")
}
