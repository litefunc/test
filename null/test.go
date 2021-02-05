package null

import (
	"encoding/json"
	"fmt"
	"testing"
)

func js(t *testing.T, v interface{}) string {
	by, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(by)
}

func jsequal(t *testing.T, v interface{}, want string) error {
	if got := js(t, v); got != want {
		return fmt.Errorf(`want:%s, got:%s`, want, got)
	}
	return nil
}

func um(t *testing.T, data []byte, v interface{}) {
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatal(err)
	}
}
