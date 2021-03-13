package util

import (
	"net/http"
)

// Respond issues an http response to the client
func Respond(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	w.Write([]byte(body))
}
