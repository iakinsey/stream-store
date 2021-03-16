package controller

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"os"
	"path"

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
	f := util.NewTempFile()
	h := sha1.New()
	body := r.Body

	for {
		checksum, err := readBlock(body, h, f)

		if err != nil {
			code := 500
			message := http.StatusText(http.StatusInternalServerError)

			if e, ok := err.(*errs.HTTPError); ok {
				code = e.Code
				message = e.Err.Error()
			} else {
				log.Println(message)
			}

			util.Respond(w, code, message)
			clearTempFile(f)
			return
		}

		if checksum != nil {
			code := 200
			message := *checksum
			err := finalizeFile(f, message)

			if err != nil {
				code = 500
				message = http.StatusText(http.StatusInternalServerError)
				log.Println(message)

				clearTempFile(f)
			}

			util.Respond(w, code, message)
			return
		}
	}
}

func readBlock(body io.Reader, h hash.Hash, f *os.File) (finalChecksum *string, err error) {
	checksum, err := getChecksum(body)

	if err == io.EOF {
		finalChecksumVal := hex.EncodeToString(h.Sum(nil))
		return &finalChecksumVal, nil
	} else if err != nil {
		return nil, err
	}

	buf := make([]byte, chunkSize)
	n, err := io.ReadFull(body, buf)

	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	chunk := buf[0:n]

	if err = checksumBlock(chunk, checksum, h); err != nil {
		return nil, err
	}

	f.Write(chunk)

	if n != chunkSize {
		checksum := hex.EncodeToString(checksum)
		return &checksum, err
	}

	return nil, nil
}

func getChecksum(body io.Reader) ([]byte, error) {
	checksum := make([]byte, checksumSize)
	n, err := io.ReadFull(body, checksum)

	if err == io.EOF {
		return checksum, err
	} else if n != checksumSize {
		return checksum, &errs.HTTPError{
			Code: 400,
			Err:  fmt.Errorf("Invalid checksum size %d", n),
		}
	}

	return checksum, err
}

func checksumBlock(chunk, expectedChecksum []byte, h hash.Hash) error {
	h.Write(chunk)

	if actualChecksum := h.Sum(nil); bytes.Compare(expectedChecksum, actualChecksum) != 0 {
		return &errs.HTTPError{
			Code: 409,
			Err:  fmt.Errorf("Integrity error, wanted %x, found %x", expectedChecksum, actualChecksum),
		}
	}

	return nil
}

func clearTempFile(f *os.File) {
	log.Println(os.Remove(f.Name()))
}

func finalizeFile(f *os.File, checksum string) error {
	storeDir := util.GetOrCreateAppRelativeDir("store")
	finalPath := path.Join(storeDir, checksum)

	err := os.Rename(f.Name(), finalPath)

	if e, ok := err.(*os.LinkError); ok {
		return e
	}

	return err
}
