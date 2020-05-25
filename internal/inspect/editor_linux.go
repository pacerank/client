// +build linux

package inspect

import "strings"

func Editor(process string) (string, bool) {
	var (
		result string
	)

	switch true {
	case strings.Contains(process, "idea64"):
		result = "IntelliJ IDEA"
		break
	case strings.Contains(process, "atom"):
		result = "Atom"
		break
	}

	return result, result != ""
}
