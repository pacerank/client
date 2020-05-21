// +build linux

package system

import (
	"os/user"
)

type target struct{}

func HomePath() string {
	usr, _ := user.Current()
	return usr.HomeDir
}
