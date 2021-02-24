package internal

import (
	"cloud/lib/logger"
	"io"
	"net/http"
	"os"
)

func writeError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}

func ok(w http.ResponseWriter, by []byte) {
	w.WriteHeader(http.StatusOK)
	if len(by) != 0 {
		w.Write(by)
	}
}

func FormFile(w http.ResponseWriter, r *http.Request) {

	f, fh, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	logger.Debug(fh)
	defer f.Close()

	f1, err := os.OpenFile("test", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Error(err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer f1.Close()

	if _, err := io.Copy(f1, f); err != nil {
		logger.Error(err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	ok(w, nil)
}
