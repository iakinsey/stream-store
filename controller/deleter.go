package controller

import (
	"net/http"

	"github.com/iakinsey/stream-store/util"
)

// Deleter ...
func Deleter(w http.ResponseWriter, r *http.Request) {
	checksum := util.GetChecksumFromURL(r.URL.Path)

	if checksum == "" {
		util.Respond(w, http.StatusBadRequest, "No checksum specified in URL")
		return
	}

	if exists, err := util.DeleteStoreFile(checksum); err != nil {
		util.RespondInternalError(w, err)
	} else if !exists {
		util.RespondStandard(w, http.StatusNotFound)
	} else {
		util.RespondStandard(w, http.StatusOK)
	}
}
