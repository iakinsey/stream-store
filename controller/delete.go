package controller

import (
	"net"

	"github.com/iakinsey/stream-store/config"
	"github.com/iakinsey/stream-store/util"
)

// Delete ...
func Delete(conn net.Conn) {
	checksum, err := util.ReadChecksum(conn)

	if err != nil {
		return
	}

	if exists, err := util.DeleteStoreFile(checksum); err != nil {
		util.RespondInternalError(conn, err)
	} else if !exists {
		util.Respond(conn, config.ResponseNotFound)
	} else {
		util.Respond(conn, config.ResponseSuccess)
	}
}
