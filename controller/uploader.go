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
		n, checksum, err := readBlock(body, h, f)

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

		if n != chunkSize {
			code := 200
			message := hex.EncodeToString(checksum)
			err := finalizeFile(f, checksum)

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

func clearTempFile(f *os.File) {
	log.Println(os.Remove(f.Name()))
}

func finalizeFile(f *os.File, checksum []byte) error {
	checksumStr := hex.EncodeToString(checksum)
	storeDir := getOrCreateStoreDir()
	finalPath := path.Join(storeDir, checksumStr)

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
