package watcher

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pacerank/client/pkg/system"
	notify "github.com/radovskyb/watcher"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CodeEvent struct {
	FilePath string
	FileName string
	Language string
	GitRepo  string
	Err      error
}

type CodeCallback func(event CodeEvent)

var ignoreDirectories = []string{
	"node_modules",
	".git",
	".idea",
}

var ignoreExtensions = []string{
	".exe",
	".log",
	".sum",
	"package-lock.json",
	".gitignore",
}

func Code(directory string, c CodeCallback) {
	// setup fsnotify
	fs, err := fsnotify.NewWatcher()
	if err != nil {
		c(CodeEvent{Err: err})
		return
	}

	defer func() {
		err := fs.Close()
		if err != nil {
			c(CodeEvent{Err: err})
		}
	}()

	// setup file polling
	w := notify.New()
	defer w.Close()

	w.FilterOps(notify.Write, notify.Create, notify.Remove)

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-fs.Events:
				if !ok {
					break
				}

				if ignoreExtension(event.Name) {
					break
				}

				info, err := os.Stat(event.Name)
				if err != nil {
					c(CodeEvent{Err: err})
					break
				}

				if event.Op == fsnotify.Write {
					if info.IsDir() {
						break
					}

					c(CodeEvent{
						FilePath: event.Name,
						FileName: "",
						Language: "",
						GitRepo:  "",
						Err:      nil,
					})
					break
				}

				if info.IsDir() {
					// TODO: Add fsnotify if create folder
					if event.Op == fsnotify.Create {

					}

					// TODO: Add remove if remove folder
					if event.Op == fsnotify.Remove {

					}
				}
			case err, ok := <-fs.Errors:
				if !ok {
					break
				}

				c(CodeEvent{Err: err})
			}
		}
	}()

	go func() {
		for {
			select {
			case event := <-w.Event:
				if ignoreExtension(event.Name()) {
					break
				}

				if event.Op == notify.Write {
					if event.IsDir() {
						break
					}

					c(CodeEvent{
						FilePath: event.Path,
						FileName: event.Name(),
						Language: "",
						GitRepo:  "",
						Err:      nil,
					})
				}

				if event.IsDir() {
					// TODO: Add folder to watch
					if event.Op == notify.Create {
						err = w.Add(event.Path)
						if err != nil {
							c(CodeEvent{Err: err})
						}
					}

					// Todo: Add remove if remove folder
					if event.Op == notify.Remove {
						err = w.Remove(event.Path)
						if err != nil {
							c(CodeEvent{Err: err})
						}
					}
				}
			case err := <-w.Error:
				c(CodeEvent{Err: err})
			case <-w.Closed:
				return
			}
		}
	}()

	for _, d := range watchList(directory) {
		// Try to add to fsnotify first, for performance
		err := fs.Add(d)
		if err == nil {
			log.Debug().Str("path", d).Msg("folder is watched by fsnotify")
			continue
		}

		// If it fails, fall back to file polling
		err = w.Add(d)
		if err != nil {
			c(CodeEvent{Err: err})
			return
		}

		log.Debug().Str("path", d).Msg("folder is watched by file system polling")
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Second * 15); err != nil {
		c(CodeEvent{Err: err})
		return
	}

	<-done
}

func watchList(folder string) []string {
	var result []string

	abs, err := filepath.Abs(folder)
	if err != nil {
		log.Error().Err(err).Msg("could not get absolute file path")
		return result
	}

	result = append(result, abs)

	list, err := ioutil.ReadDir(abs)
	if err != nil {
		log.Error().Err(err).Msg("could not read dir")
		return result
	}

	for _, file := range list {
		if file.IsDir() {
			if ignoreDirectory(fmt.Sprintf("%s/%s", abs, file.Name())) {
				continue
			}

			result = append(result, watchList(fmt.Sprintf("%s/%s", abs, file.Name()))...)
		}
	}

	return result
}

func ignoreDirectory(path string) bool {
	for _, ignore := range ignoreDirectories {
		if strings.Contains(filepath.ToSlash(filepath.FromSlash(path)), filepath.ToSlash(filepath.FromSlash(ignore))) {
			return true
		}
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

type Key uint16

type KeyEvent struct {
	Key  Key
	Rune rune
	Err  error
}

type KeyboardCallback func(event KeyEvent)

// Listen for keyboard inputs
func Keyboard(c KeyboardCallback) {
	channel := make(chan byte)
	sys := system.New()
	go sys.ListenKeyboard(channel)

	for {
		key := <-channel
		c(KeyEvent{
			Key:  Key(key),
			Rune: rune(key),
			Err:  nil,
		})
	}
}
