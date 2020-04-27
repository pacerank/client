package system_test

import (
	"github.com/pacerank/client/pkg/system"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestLinux_Processes(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skipf("current operating system is not target")
	}

	sys := system.New()

	_, err := sys.Processes()
	assert.NoError(t, err, "finding processes should not result in error")
}
