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

	w.WriteHeader(http.StatusContinue)

	for {
		buf := make([]byte, config.ChunkSize)
		_, err := io.ReadFull(f, buf)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		} else if err != nil {
			util.IssueWriteError(w, err)
			return
		}

		if _, err := h.Write(buf); err != nil {
			util.IssueWriteError(w, err)
			return
		}

		if _, err = w.Write(buf); err != nil {
			util.IssueWriteError(w, err)
			return
		}
	}
}
