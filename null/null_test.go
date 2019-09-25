package null

import (
	"cloud/lib/logger"
	"encoding/json"
	"testing"
)

type null struct {
	Int64 Int64 `json:"int64"`
	Int   int    `json:"int,omitempty"`
}

func (n null) Json() json.RawMessage {
	by, _ := json.Marshal(n)
	return by
}

func (n null) String() string {
	return string(n.Json())
}

func TestNull(t *testing.T) {

	var n null
	logger.Debug(n.String())
	var n1 null
	if err := json.Unmarshal(n.Json(), &n1);err !=nil {
		t.Error(err)
	}
	logger.Debug(n1.String())
}
