package internal

import (
	"encoding/json"
	"test/logger"
)

type S1 struct {
	A json.RawMessage `json:"a"`
	B json.RawMessage `json:"b"`
	C json.RawMessage `json:"c"`
	D json.RawMessage `json:"d,omitempty"`
}

const s = `{"a":{"first":"Janet","last":"Prichard"},"b":47}`
const a = `{
	"auths": {
		"aa": {
			"auth": "123"
		},
		"ab": {
			"auth": "456"
		},
		"ac": {
			"auth": "789"
		}
	},
	"HttpHeaders": {
		"User-Agent": "Docker-Client/19.03.13 (linux)"
	},
	"b": {
		"b": "b"
	}
}`

func F() {

	var s1 S1
	um([]byte(s), &s1)
	j(s1)

	var a1 Auth
	um([]byte(a), &a1)
	j(a1)
}

func j(o interface{}) {
	logger.Debug(o)
	by, err := json.Marshal(o)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(string(by))
}

func um(data []byte, o interface{}) {
	logger.Debug(string(data))
	logger.Debug(o)

	if err := json.Unmarshal(data, o); err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(o)
}

type AutoGenerated struct {
	Auths struct {
		Aa struct {
			Auth string `json:"auth"`
		} `json:"aa"`
		Ab struct {
			Auth string `json:"auth"`
		} `json:"ab"`
		Ac struct {
			Auth string `json:"auth"`
		} `json:"ac"`
	} `json:"auths"`
	HTTPHeaders struct {
		UserAgent string `json:"User-Agent"`
	} `json:"HttpHeaders"`
}

type Auth struct {
	Auths map[string]struct {
		Auth string `json:"auth"`
	} `json:"auths"`
	json.RawMessage `json:",omitempty"`
	// M
}

type M map[string]json.RawMessage

func S() {
	var r json.RawMessage
	if err := json.Unmarshal([]byte(`{"a": 1,   "b":{  }}`), &r); err != nil {
		logger.Fatal(err)
	}
	logger.Debug(string(r))
}
