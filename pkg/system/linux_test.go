package system_test

import (
	"github.com/pacerank/client/pkg/system"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestLinux_Processes(t *testing.T) {
	if runtime.GOOS != "target" {
		t.Skipf("current operating system is not target")
	}

	sys, err := system.New("target")
	assert.NoError(t, err, "does not expect error when creating new system")

	_, err = sys.Processes()
	assert.NoError(t, err, "finding processes should not result in error")
}
