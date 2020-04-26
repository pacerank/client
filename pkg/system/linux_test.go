package system_test

import (
	"github.com/pacerank/client/pkg/system"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestLinux_Processes(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skipf("current operating system is not linux")
	}

	sys, err := system.New("linux")
	assert.NoError(t, err, "does not expect error when creating new system")

	_, err = sys.Processes()
	assert.NoError(t, err, "finding processes should not result in error")
}
