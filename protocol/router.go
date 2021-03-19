package protocol

import (
	"fmt"
	"io"
	"net"

	"github.com/iakinsey/stream-store/config"
	"github.com/iakinsey/stream-store/controller"
	"github.com/iakinsey/stream-store/util"
)

func Listen(listener net.Listener) error {
	for {
		conn, err := listener.Accept()

		if err != nil {
			return err
		}

		// TODO make into a channel?
		go onConnect(conn)
	}
}

func onConnect(conn net.Conn) {
	requestBuf := make([]byte, 1)

	n, err := io.ReadFull(conn, requestBuf)

	if err != nil {
		util.RespondError(conn, config.ResponseInternalError, err)
		return
	} else if n != 1 {
		err = fmt.Errorf("Expected to read 1 byte, got %d", n)
		util.RespondError(conn, config.ResponseInternalError, err)
		return
	}

	switch code := requestBuf[0]; code {
	case config.RequestUpload:
		controller.Upload(conn)
	case config.RequestDownload:
		controller.Download(conn)
	case config.RequestDelete:
		controller.Delete(conn)
	default:
		err = fmt.Errorf("Invalid Operation requested: %d", code)
		util.RespondError(conn, config.ResponseInvalidOperation, err)
		return
	}
}
