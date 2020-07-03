//go:generate goversioninfo -icon=resources/original_icon_large.ico -manifest=resources/pacerank.exe.manifest -64

package main

import (
	tool "github.com/GeertJohan/go.rice"
	"github.com/getlantern/systray"
	"github.com/pacerank/client/internal/gui"
	"github.com/pacerank/client/internal/inspect"
	"github.com/pacerank/client/internal/operation"
	"github.com/pacerank/client/internal/store"
	"github.com/pacerank/client/internal/watcher"
	"github.com/pacerank/client/pkg/api"
	"github.com/pacerank/client/pkg/system"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	// Operational setup
	err := operation.Setup()
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup application")
	}

	// Gracefully shutdown application on termination
	defer func() {
		err = operation.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("could not terminate application gracefully")
		}
	}()

	systray.Run(onReady, onExit)
}

func onReady() {
	apiClient := api.New("https://digest.pacerank.io")

	storage, err := store.New()
	if err != nil {
		log.Fatal().Err(err).Msgf("could start store")
	}

	defer storage.Close()

	if storage.AuthorizationToken() != "" {
		apiClient.AddAuthorizationToken(storage.AuthorizationToken())
	}

	directories, err := storage.Directories()
	for _, directory := range directories {
		go watcher.Code(directory.Directory, func(event watcher.CodeEvent) {
			if event.Err != nil {
				log.Error().Err(event.Err).Msg("could not watch code")
				return
			}

			err = storage.AddHeap(store.InHeap{
				Id:       event.Id,
				Language: event.Language,
				Branch:   event.Branch,
				FileName: event.FilePath,
				Project:  event.Project,
				Git:      event.Git,
			})
			if err != nil {
				log.Error().Err(err).Msg("could not save code activity to store")
			}

			log.Info().
				Str("language", event.Language).
				Str("filepath", event.FilePath).
				Str("filename", event.FileName).
				Str("project", event.Project).
				Str("branch", event.Branch).
				Str("git", event.Git).
				Str("id", event.Id).
				Msg("found code change")
		})
	}

	// Poll store to see if anything should be queued for dispatch to digest service
	go watcher.Sessions(storage)

	// Poll store to see if any messages are in queue, and send them
	go watcher.Queue(storage, apiClient, func(structure *api.DefaultReplyStructure, err error) {
		if err != nil {
			log.Error().Err(err).Msgf("could not send message to digest service: %s", structure.CorrelationId)
			return
		}

		log.Info().Msgf("digest has acknowledged the message: %s", structure.CorrelationId)
	})

	// Setup a new session and meta
	err = storage.NewSession()

	// Start recording typing
	sys := system.New()
	go watcher.Keyboard(func(key watcher.KeyEvent) {
		process, err := sys.ActiveProcess()
		if err != nil {
			log.Error().Err(err).Msgf("could not get active process")
			return
		}

		editor, ok := inspect.Editor(process.Executable)
		if !ok {
			return
		}

		err = storage.MetaTypingActivity(editor)
		if err != nil {
			log.Error().Err(err).Msg("could not record typing activity to store")
			return
		}

		log.Debug().Msgf("recorded typing activity in %s", process.FileName)
	})

	box, err := tool.FindBox("resources")
	if err != nil {
		log.Error().Err(err).Msg("could not find resources")
		return
	}

	systray.SetIcon(box.MustBytes("original_icon_large.ico"))
	systray.SetTitle("PaceRank")
	systray.SetTooltip("Currently collecting your programming measurements")
	start := systray.AddMenuItem("Show", "Show the app")
	quit := systray.AddMenuItem("Quit", "Quit the app")

	for {
		select {
		case <-start.ClickedCh:
			win := gui.Start(storage, apiClient)
			if win == nil {
				return
			}

			win.Show()
			win.Run()
		case <-quit.ClickedCh:
			systray.Quit()
			os.Exit(0)
		}
	}
}

func onExit() {
	os.Exit(0)
}
