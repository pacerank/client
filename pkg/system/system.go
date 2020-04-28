package system

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

type System interface {
	Processes() ([]Process, error)
	ActiveProcess() (*Process, error)
}

type Process struct {
	ProcessID  int64
	Parent     int64
	Children   []int64
	FileName   string
	Checksum   string
	Executable string
}

type Module struct {
	ProcessID  int64
	ParentID   int64
	Checksum   string
	Executable string
	FileName   string
}

func New() System {
	return &target{}
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
