package controller

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
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
	// TODO move and delete temp file
	f, err := ioutil.TempFile(os.TempDir(), "streamstore")

	if err != nil {
		status := http.StatusInternalServerError
		util.Respond(w, status, http.StatusText(status))
		log.Fatalf(err.Error())
	}

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
			err := finalizeFile(f, *checksum)
			message := http.StatusText(http.StatusInternalServerError)

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

func readBlock(body io.ReadCloser, h hash.Hash, f *os.File) (finalChecksum *string, err error) {
	checksum, err := getChecksum(body)

	if err != nil {
		return nil, err
	}

	buf := make([]byte, chunkSize)
	n, err := body.Read(buf)

	if err != nil {
		return nil, err
	}

	chunk := buf[0:n]

	if err = checksumBlock(chunk, checksum, h); err != nil {
		return nil, err
	} else if n != chunkSize {
		checksum := hex.EncodeToString(checksum)
		return &checksum, err
	}

	return nil, nil
}

func getChecksum(body io.ReadCloser) ([]byte, error) {
	checksum := make([]byte, checksumSize)
	n, err := body.Read(checksum)

	if n != checksumSize {
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
	storeDir := getOrCreateStoreDir()
	finalPath := path.Join(storeDir, checksum)

	return os.Rename(f.Name(), finalPath)
}

func getOrCreateStoreDir() string {
	execPath, err := os.Executable()

	if err != nil {
		log.Fatalf(err.Error())
	}

	storeDir := path.Join(path.Dir(execPath), "store")

	if _, err = os.Stat(storeDir); os.IsNotExist(err) {
		err = os.Mkdir(storeDir, 0755)

		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	return storeDir
}
