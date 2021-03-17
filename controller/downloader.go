package controller

import (
	"crypto/sha1"
	"io"
	"net/http"

	"github.com/iakinsey/stream-store/config"
	"github.com/iakinsey/stream-store/util"
)

// Downloader ...
func Downloader(w http.ResponseWriter, r *http.Request) {
	checksum := util.GetChecksumFromURL(r.URL.Path)

	if checksum == "" {
		util.Respond(w, http.StatusBadRequest, "No checksum specified in URL")
		return
	}

	exists, f, err := util.GetStoreFile(checksum)

	if !exists {
		util.RespondStandard(w, http.StatusNotFound)
		return
	} else if err != nil {
		util.RespondInternalError(w, err)
		return
	}

	h := sha1.New()
	complete := false

	w.WriteHeader(http.StatusAccepted)

	for {
		buf := make([]byte, config.ChunkSize)
		n, err := io.ReadFull(f, buf)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			complete = true
		} else if err != nil {
			util.IssueWriteError(w, err)
			return
		}

		chunk := buf[0:n]

		// Update hash
		if _, err := h.Write(chunk); err != nil {
			util.IssueWriteError(w, err)
			return
		}

		// Write hash to client
		if _, err = w.Write(h.Sum(nil)); err != nil {
			util.IssueWriteError(w, err)
			return
		}

		// Write bytes to client
		if _, err = w.Write(chunk); err != nil {
			util.IssueWriteError(w, err)
			return
		}

		if complete {
			return
		}
	}
}
