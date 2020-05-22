// +build linux

package system

import (
	"os/exec"
)

type target struct{}

func OpenBrowser(url string) error {
	return exec.Command("xdg-open", url).Start()
}
