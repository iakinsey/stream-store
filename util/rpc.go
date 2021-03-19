package util

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"

	"github.com/iakinsey/stream-store/config"
	"github.com/iakinsey/stream-store/errs"
)

func RespondError(conn net.Conn, code byte, theErr error) {
	log.Printf("RespondError (Code %d): %v", int(code), theErr)

	if _, err := conn.Write([]byte{code}); err != nil {
		log.Printf("Respond error failure: %v", err)
	}

	if err := conn.Close(); err != nil {
		log.Printf("Close connection failure: %v", err)
	}
}

func Respond(conn net.Conn, code byte) error {
	if _, err := conn.Write([]byte{code}); err != nil {
		RespondError(conn, config.ResponseInternalError, fmt.Errorf("Error during respond: %v", err))
		return err
	}

	return nil
}

func HandleResponseError(conn net.Conn, err errs.ResponseError) {
	RespondError(conn, err.Code, err.Err)
}

func RespondInternalError(conn net.Conn, err error) {
	RespondError(conn, config.ResponseInternalError, err)
}

func RespondWriteError(conn net.Conn, err error) {
	log.Printf("Sending write error: %v", err)
	conn.Write(config.ResponseWriteError)

	if err := conn.Close(); err != nil {
		log.Printf("Close connection failure: %v", err)
	}
}

func ReadChecksum(conn net.Conn) (string, error) {
	buf := make([]byte, config.ChecksumSize)

	if n, err := conn.Read(buf); err != nil {
		RespondInternalError(conn, err)
		return "", err
	} else if n != config.ChecksumSize {
		err := fmt.Errorf("Invalid checksum size, expected %d, got %d", n, config.ChecksumSize)
		RespondError(conn, config.ResponseMalformedRequest, err)
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
