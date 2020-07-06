// +build windows

package inspect

import "strings"

func Editor(process string) (string, bool) {
	var (
		result string
	)

	switch true {
	// START JETBRAINS
	case strings.Contains(process, "idea64.exe"):
		result = "IntelliJ IDEA"
		break
	case strings.Contains(process, "idea.exe"):
		result = "IntelliJ IDEA"
		break
	case strings.Contains(process, "goland64.exe"):
		result = "GoLand"
		break
	case strings.Contains(process, "goland.exe"):
		result = "GoLand"
		break
	case strings.Contains(process, "datagrip64.exe"):
		result = "DataGrip"
		break
	case strings.Contains(process, "datagrip.exe"):
		result = "DataGrip"
		break
	case strings.Contains(process, "phpstorm64.exe"):
		result = "PhpStorm"
		break
	case strings.Contains(process, "phpstorm.exe"):
		result = "PhpStorm"
		break
	case strings.Contains(process, "pycharm64.exe"):
		result = "PyCharm"
		break
	case strings.Contains(process, "rubymine64.exe"):
		result = "RubyMine"
		break
	case strings.Contains(process, "webstorm.exe"):
		result = "WebStorm"
		break
	case strings.Contains(process, "webstorm64.exe"):
		result = "WebStorm"
		break
	case strings.Contains(process, "clion64.exe"):
		result = "CLion"
		break
	case strings.Contains(process, "rider64.exe"):
		result = "Jetbrains Rider"
		break
		// END JETBRAINS
		// START MICROSOFT
	case strings.Contains(process, "Code.exe"):
		result = "Visual Studio Code"
		break
	case strings.Contains(process, "devenv.exe"):
		result = "Visual Studio"
		break
		// END MICROSOFT
		// START VARIOUS
	case strings.Contains(process, "atom.exe"):
		result = "Atom"
		break
	case strings.Contains(process, "sublime_text.exe"):
		result = "Sublime Text"
		break
	case strings.Contains(process, "notepad++.exe"):
		result = "Notepad++"
		break
	case strings.Contains(process, "vim.exe"):
		result = "VIM"
		break
	case strings.Contains(process, "bluefish.exe"):
		result = "Bluefish"
		break
	case strings.Contains(process, "Unity.exe"):
		result = "Unity"
		break
	case strings.Contains(process, "Brackets.exe"):
		result = "Brackets"
		break
	}

	return result, result != ""
}
