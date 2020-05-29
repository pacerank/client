package inspect

import (
	"errors"
	"github.com/go-enry/go-enry/v2"
	"io/ioutil"
	"strings"
)

func AnalyzeFile(path string, filename string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	if enry.IsImage(path) {
		return "", errors.New("file is image, skip")
	}

	if enry.IsBinary(b) {
		return "", errors.New("file is binary, skip")
	}

	if enry.IsConfiguration(path) {
		return "", errors.New("file is configuration, skip")
	}

	if enry.IsDocumentation(path) {
		return "", errors.New("file is documentation, skip")
	}

	if enry.IsDotFile(path) {
		return "", errors.New("file is dotfile, skip")
	}

	if enry.IsVendor(path) {
		return "", errors.New("file is vendor, skip")
	}

	if IsIgnored(path, filename) {
		return "", errors.New("matches file ignore pattern, skip")
	}

	lang := strings.ToLower(enry.GetLanguage(filename, b))

	if lang == "Text" {
		return lang, errors.New("file is text, skip")
	}

	if lang == "" {
		return lang, errors.New("could not determine language, skip")
	}

	return lang, nil
}

func IsIgnored(path, filename string) bool {
	if ignoreExtension(filename) {
		return true
	}

	if filename[len(filename):] == "~" {
		return true
	}

	return false
}

func ignoreExtension(filename string) bool {
	for _, ignore := range ignoreExtensions {
		if strings.Contains(filename, ignore) {
			return true
		}
	}

	return false
}

var ignoreExtensions = []string{
	".log",
	"package-lock.json",
	".gitignore",
}
