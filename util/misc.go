package util

import (
	"io/ioutil"
	"log"
	"os"
	"path"
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
	tempDir := GetOrCreateAppRelativeDir("download")
	f, err := ioutil.TempFile(tempDir, "streamstore")

	if err != nil {
		log.Fatalf(err.Error())
	}

	return f
}
