// +build linux

package system

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type target struct{}

func (t *target) Processes() ([]Process, error) {
	var result = make([]Process, 0)

	proc, err := os.Open("/proc")
	if err != nil {
		return result, err
	}

	defer proc.Close()

	names, err := proc.Readdirnames(0)
	if err == io.EOF {
		return result, nil
	}

	if err != nil {
		return result, err
	}

	for _, name := range names {
		pid, err := strconv.ParseInt(name, 10, 32)
		if err != nil {
			continue
		}

		path, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
		if os.IsPermission(err) {
			continue
		}

		if err != nil {
			return result, err
		}

		cs, err := checksum(path)
		if err != nil {
			return result, err
		}

		result = append(result, Process{
			ProcessID:  pid,
			Parent:     0,
			Children:   nil,
			Modules:    nil,
			FileName:   getExecutableName(path),
			Checksum:   cs,
			Executable: path,
		})
	}

	return result, nil
}

func (t *target) ActiveProcess() (*Process, error) {
	return &Process{}, nil
}

// Get name of executable from path
func getExecutableName(path string) string {
	return path[strings.LastIndex(path, "/")+1:]
}
