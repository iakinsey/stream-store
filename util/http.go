package util

import (
	"log"
	"net/http"

	"github.com/iakinsey/stream-store/config"
)

// Respond issues an http response to the client
func Respond(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	w.Write([]byte(body))
}

// RespondStandard responds with standard http status text
func RespondStandard(w http.ResponseWriter, code int) {
	Respond(w, code, http.StatusText(code))
}

// RespondInternalError issues a 500 internal server error response
func RespondInternalError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	message := http.StatusText(code)

	log.Println(err.Error())
	Respond(w, code, message)
}

// IssueWriteError informs client of internal error during write
func IssueWriteError(w http.ResponseWriter, err error) {
	log.Println(err.Error())
	w.Write(config.WriteErrorContent)
}
