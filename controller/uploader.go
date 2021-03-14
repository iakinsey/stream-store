package controller

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"

	"github.com/iakinsey/stream-store/errs"
	"github.com/iakinsey/stream-store/util"
)

// BadRequest ...
type BadRequest error

// ConflictError ...
type ConflictError error

// Size of checksum block
const checksumSize = sha1.Size

// Size of content chunk block
const chunkSize = 128 << 10

// Uploader ...
func Uploader(w http.ResponseWriter, r *http.Request) {
	// TODO get tempfile path and move on success
	f, err := os.Open("")

	if err != nil {
		panic(err)
	}

	h := sha1.New()
	body := r.Body

	for {
		n, checksum, err := readBlock(body, h, f)

		if err != nil {
			code := 500
			message := "Internal server error."

			if e, ok := err.(*errs.HTTPError); ok {
				code = e.Code
				message = e.Err.Error()
			}

			util.Respond(w, code, message)
			return
		}

		if n != chunkSize {
			util.Respond(w, 200, hex.EncodeToString(checksum))
			return
		}
	}

}

func readBlock(body io.ReadCloser, h hash.Hash, f *os.File) (n int, checksum []byte, err error) {
	buf := make([]byte, chunkSize)
	n, err = body.Read(buf)

	if err != nil {
		return n, checksum, err
	}

	chunk := buf[0:n]
	checksum, err = checksumBlock(chunk, body, h)

	if err == nil {
		f.Write(chunk)
	}

	return n, checksum, err
}

func checksumBlock(chunk []byte, body io.ReadCloser, h hash.Hash) (result []byte, err error) {
	expected := make([]byte, checksumSize)
	n, err := body.Read(expected)

	if err != nil {
		return result, err
	} else if n != checksumSize {
		return result, &errs.HTTPError{
			Code: 400,
			Err:  fmt.Errorf("Invalid checksum size: %d", n),
		}
	}

	h.Write(chunk)
	actual := h.Sum(nil)

	if actual := h.Sum(nil); bytes.Compare(expected, actual) != 0 {
		return result, &errs.HTTPError{
			Code: 409,
			Err:  fmt.Errorf("Integrity error, wanted %x, found %x", expected, actual),
		}
	}

	return actual, nil
}
