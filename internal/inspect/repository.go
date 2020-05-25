package inspect

import (
	"encoding/base64"
	"github.com/go-git/go-git/v5"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ProjectInfo struct {
	Id       string
	Project  string
	Git      string
	Branch   string
	FileName string
}

func Project(filePath, watcherPath string) (ProjectInfo, error) {
	var (
		pi  ProjectInfo
		err error
	)

	filePath = filepath.ToSlash(filePath)

	var gitFound bool

	for {
		if !gitFound {
			pi.Git, pi.Branch, gitFound = gitInfo(filePath)
		}

		if gitFound {
			pi.Project = filePath[strings.LastIndex(filePath, string(os.PathSeparator))+1:]
			pi.Id = base64.StdEncoding.EncodeToString([]byte(filePath))
			break
		}

		if path.Clean(filePath) == path.Clean(watcherPath) {
			break
		}

		filePath = filepath.Dir(filePath)
	}

	return pi, err
}

func gitInfo(filePath string) (urls string, branch string, ok bool) {
	repository, err := git.PlainOpen(filePath)
	if err != nil {
		return
	}

	ok = true

	remote, err := repository.Remote("origin")
	if err != nil {
		return
	}

	urls = strings.Join(remote.Config().URLs, ",")

	head, err := repository.Head()
	if err != nil {
		return
	}

	branch = head.Name().String()
	return
}
