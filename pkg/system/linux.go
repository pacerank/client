// +build linux

package system

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type target struct{}

func (t *target) Processes() ([]*Process, error) {
	var result = make([]*Process, 0)

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

		process, err := getProcess(pid)
		if err != nil {
			continue
		}

		result = append(result, process)
	}

	return result, nil
}

func (t *target) ActiveProcess() (*Process, error) {
	if !activeWindowSupport() {
		return nil, errors.New("your window manager does not support _NET_ACTIVE_WINDOW")
	}

	return &Process{}, nil
}

// Function looks if the current xorg supports _NET_ACTIVE_WINDOW
func activeWindowSupport() bool {
	cmd := exec.Command("xprop", "-root", "_NET_SUPPORTED")
	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		return false
	}

	str := strings.ReplaceAll(strings.Trim(output.String(), "\n"), " ", "")
	str = str[strings.LastIndex(str, "=")+1:]
	strArr := strings.Split(str, ",")
	for _, supported := range strArr {
		if supported == "_NET_ACTIVE_WINDOW" {
			return true
		}
	}

	return false
}

// Create process struct for given processID
func getProcess(processID int64) (*Process, error) {
	path, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", processID))
	if os.IsPermission(err) {
		return nil, errors.New("no permission to read pid")
	}

	if err != nil {
		return nil, err
	}

	cs, err := checksum(path)
	if err != nil {
		return nil, err
	}

	return &Process{
		ProcessID:  processID,
		Parent:     0,
		Children:   nil,
		FileName:   getExecutableName(path),
		Checksum:   cs,
		Executable: path,
	}, nil
}

// Get name of executable from path
func getExecutableName(path string) string {
	return path[strings.LastIndex(path, "/")+1:]
}
