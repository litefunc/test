package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/auth"
	"cloud/server/portal/web/response"
	"context"
	"net/http"
	"net/url"
)

// UserContextKey is the key to retrieve a web.User from a http.Request.
const UserContextKey = contextKey("user")

type contextKey string

var mds = []Middleware{allowOrigin, AuthRequired, logRequest}

func allowOrigin(h http.HandlerFunc) http.HandlerFunc {
	hd := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, HEAD, DELETE")
		h(w, r)
	}
	return hd
}

type HTTPLog struct {
	Method     string              `json:"method"`
	Host       string              `json:"host"`
	URL        *url.URL            `json:"url"`
	Proto      string              `json:"proto"`
	RemoteAddr string              `json:"remote_addr"`
	Header     map[string][]string `json:"header"`
}

func logRequest(h http.HandlerFunc) http.HandlerFunc {
	hd := func(w http.ResponseWriter, r *http.Request) {

		httpLog := HTTPLog{Method: r.Method, Host: r.Host, URL: r.URL, Proto: r.Proto, RemoteAddr: r.RemoteAddr}
		httpLog.Header = make(map[string][]string)

		for name, headers := range r.Header {
			if name != "Authorization" {
				httpLog.Header[name] = headers
			}
		}

		logger.HTTP(httpLog)
		h(w, r)
	}
	return hd
}

func AuthRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := r.Header.Get("Authorization")
		if len(s) > 7 && s[:7] == "Bearer " {
			u, err := auth.AuthByToken(s[7:])
			if err != nil {
				response.ErrTokenInvalid().Send(w)
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey, u)
			h(w, r.WithContext(ctx))
			return
		}
		response.ErrUnauthorized().Send(w)
	}
}
