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

	"github.com/iakinsey/stream-store/config"
	"github.com/iakinsey/stream-store/errs"
	"github.com/iakinsey/stream-store/util"
)

// Uploader ...
func Uploader(w http.ResponseWriter, r *http.Request) {
	f := util.NewTempFile()
	h := sha1.New()
	body := r.Body

	for {
		checksum, err := readBlock(body, h, f)

		if err != nil {
			if e, ok := err.(*errs.HTTPError); ok {
				util.Respond(w, e.Code, e.Err.Error())
			} else {
				util.RespondInternalError(w, err)
			}

			clearTempFile(f)
			return
		}

		if checksum != nil {
			err := finalizeFile(f, *checksum)

			if err != nil {
				util.RespondInternalError(w, err)
				clearTempFile(f)
			} else {
				util.Respond(w, http.StatusCreated, *checksum)
			}

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

	buf := make([]byte, config.ChunkSize)
	n, err := io.ReadFull(body, buf)

	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	chunk := buf[0:n]

	if err = checksumBlock(chunk, checksum, h); err != nil {
		return nil, err
	}

	f.Write(chunk)

	if n != config.ChunkSize {
		checksum := hex.EncodeToString(checksum)
		return &checksum, err
	}

	return nil, nil
}

func getChecksum(body io.Reader) ([]byte, error) {
	checksum := make([]byte, config.ChecksumSize)
	n, err := io.ReadFull(body, checksum)

	if err == io.EOF {
		return checksum, err
	} else if n != config.ChecksumSize {
		return checksum, &errs.HTTPError{
			Code: 400,
			Err:  fmt.Errorf("Invalid checksum size %d", n),
		}
	}

	return checksum, err
}

func checksumBlock(chunk, expectedChecksum []byte, h hash.Hash) error {
	if _, err := h.Write(chunk); err != nil {
		return err
	}

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
	storeDir := util.GetOrCreateAppRelativeDir(config.StoreFolderName)
	finalPath := path.Join(storeDir, checksum)

	err := os.Rename(f.Name(), finalPath)

	if e, ok := err.(*os.LinkError); ok {
		return e
	}

	return err
}
