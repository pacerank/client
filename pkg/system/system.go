package system

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

type System interface {
	Processes() ([]Process, error)
}

type Process struct {
	Pid        int64
	FileName   string
	Checksum   string
	Executable string
}

var (
	ErrOperatingSystemNotSupported = errors.New("operating system is not supported")
)

func New(os string) (System, error) {
	switch os {
	case "windows":
		return &windows{}, nil
	case "linux":
		return &linux{}, nil
	default:
		return nil, ErrOperatingSystemNotSupported
	}
}

// Create a checksum of given file
func checksum(path string) (string, error) {
	hasher := sha256.New()

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
