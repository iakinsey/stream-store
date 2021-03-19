package controller

import (
	"crypto/sha1"
	"io"
	"net"

	"github.com/iakinsey/stream-store/config"
	"github.com/iakinsey/stream-store/util"
)

// Download ...
func Download(conn net.Conn) {
	checksum, err := util.ReadChecksum(conn)

	if err != nil {
		return
	}

	exists, f, err := util.GetStoreFile(checksum)

	if !exists {
		util.Respond(conn, config.ResponseNotFound)
		return
	} else if err != nil {
		util.RespondInternalError(conn, err)
		return
	}

	h := sha1.New()
	complete := false

	if err := util.Respond(conn, config.ResponseContinue); err != nil {
		return
	}

	for {
		buf := make([]byte, config.ChunkSize)
		n, err := io.ReadFull(f, buf)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			complete = true
		} else if err != nil {
			util.RespondWriteError(conn, err)
			return
		}

		chunk := buf[0:n]

		// Update hash
		if _, err := h.Write(chunk); err != nil {
			util.RespondWriteError(conn, err)
			return
		}

		// Write hash to client
		if _, err = conn.Write(h.Sum(nil)); err != nil {
			util.RespondWriteError(conn, err)
			return
		}

		// Write bytes to client
		if _, err = conn.Write(chunk); err != nil {
			util.RespondWriteError(conn, err)
			return
		}

		if complete {
			return
		}
	}
}
