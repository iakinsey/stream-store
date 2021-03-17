package util

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/iakinsey/stream-store/config"
)

// Environ ...
func Environ(name string, fallback string) string {
	result := os.Getenv(name)

	if result == "" {
		result = fallback
	}

	return result
}

// GetOrCreateAppRelativeDir ...
func GetOrCreateAppRelativeDir(name string) string {
	execPath, err := os.Executable()

	if err != nil {
		log.Fatalf(err.Error())
	}

	storeDir := path.Join(path.Dir(execPath), name)

	if _, err = os.Stat(storeDir); os.IsNotExist(err) {
		err = os.Mkdir(storeDir, 0755)

		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	return storeDir
}

// NewTempFile ...
func NewTempFile() (f *os.File) {
	tempDir := GetOrCreateAppRelativeDir(config.DownloadFolderName)
	f, err := ioutil.TempFile(tempDir, config.TempDownloadPrefix)

	if err != nil {
		log.Fatalf(err.Error())
	}

	return f
}

// GetStoreFile ...
func GetStoreFile(name string) (bool, *os.File, error) {
	tempDir := GetOrCreateAppRelativeDir(config.StoreFolderName)
	filePath := path.Join(tempDir, name)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, nil, nil
	}

	f, err := os.Open(filePath)

	return true, f, err
}

// DeleteStoreFile ...
func DeleteStoreFile(name string) (bool, error) {
	tempDir := GetOrCreateAppRelativeDir(config.StoreFolderName)
	filePath := path.Join(tempDir, name)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, nil
	}

	if err := os.Remove(filePath); err != nil {
		return true, err
	}

	return true, nil
}
