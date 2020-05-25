// +build windows

package inspect

import "strings"

func Editor(process string) (string, bool) {
	var (
		result string
	)

	switch true {
	case strings.Contains(process, "idea64.exe"):
		result = "IntelliJ IDEA"
		break
	case strings.Contains(process, "atom.exe"):
		result = "Atom"
		break
	}

	return result, result != ""
}
