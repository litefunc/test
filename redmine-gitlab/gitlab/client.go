package gitlab

import (
	"bytes"
	"fmt"
  	"cloud/lib/logger"
	"os"
	"os/exec"
	"strings"
)

type Client struct {
	token string
}

func NewClient(tk string) *Client {
	return &Client{tk}
}

func (cli Client) CloneWithHTTP(dir, repo string) error {

	ss := strings.Split(repo, "://")

	url := fmt.Sprintf(`%s://oauth2:%s@%s`, ss[0], cli.token, ss[1])
	return Clone(dir, url)
}

func Clone(dir, url string) error {

	if err := os.Chdir(dir); err != nil {
		return err
	}

	ss := strings.Split(url, "/")
	if n := len(ss); n != 0 {

		_, err := os.Stat(ss[n-1])
		if err != nil {
			if os.IsNotExist(err) {

				_, err := output("git", "clone", "--mirror", url)
				if err != nil {
					return err
				}
				return nil

			}

			logger.Error(err)
			return err
		}

		return nil
	}

	return fmt.Errorf(`invalid url:%s`, url)

}

func Fetch(dir string) error {

	if err := os.Chdir(dir); err != nil {
		return err
	}

	_, err := output("git", "fetch", "--all")
	if err != nil {
		return err
	}

	return nil

}

func output(name string, arg ...string) (string, error) {
	logger.Debug(name, strings.Join(arg, " "))

	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		e := fmt.Errorf(`%s: %s`, fmt.Sprint(err), stderr.String())
		logger.Debug(e)
		return out.String(), e
	}

	logger.Debug(out.String())
	return out.String(), nil
}

func run(name string, arg ...string) error {
	logger.Debug(name, strings.Join(arg, " "))

	cmd := exec.Command(name, arg...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		e := fmt.Errorf(`%s: %s`, fmt.Sprint(err), stderr.String())
		logger.Debug(e)
		return e
	}

	return nil
}
