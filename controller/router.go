package controller

import (
	"net/http"

	"github.com/iakinsey/stream-store/util"
)

// Router ...
func Router(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method == http.MethodPut {
		Uploader(w, r)
	} else if method == http.MethodGet {
		Downloader(w, r)
	} else if method == http.MethodDelete {
		Deleter(w, r)
	} else {
		text := http.StatusText(http.StatusMethodNotAllowed)
		util.Respond(w, http.StatusMethodNotAllowed, text)
	}
}
