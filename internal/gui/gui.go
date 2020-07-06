package gui

import (
	tool "github.com/GeertJohan/go.rice"
	"github.com/pacerank/client/internal/operation"
	"github.com/pacerank/client/internal/store"
	"github.com/pacerank/client/internal/watcher"
	"github.com/pacerank/client/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/rice"
	"github.com/sciter-sdk/go-sciter/window"
	"net/url"
	"runtime"
	"strings"
)

func init() {
	runtime.LockOSThread()
}

func Start(storage *store.Store, apiClient *api.Api) *window.Window {
	sciter.SetOption(sciter.SCITER_SET_SCRIPT_RUNTIME_FEATURES, sciter.ALLOW_SOCKET_IO|sciter.SCITER_SET_SCRIPT_RUNTIME_FEATURES|sciter.ALLOW_FILE_IO|sciter.ALLOW_EVAL|sciter.ALLOW_SYSINFO)
	sciter.SetOption(sciter.SCITER_SET_DEBUG_MODE, 1)

	win, err := window.New(sciter.SW_TITLEBAR|sciter.SW_CONTROLS|sciter.SW_MAIN|sciter.SW_ENABLE_DEBUG, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("create window error")
	}

	// Add handler for rice
	rice.HandleDataLoad(win.Sciter)
	_, err = tool.FindBox("resources")
	if err != nil {
		log.Error().Err(err).Msg("could not load resources box")
		return nil
	}

	win.SetTitle("PaceRank")

	// Get user signature name
	win.DefineFunction("user_signature_name", func(args ...*sciter.Value) *sciter.Value {
		return sciter.NewValue(storage.UserSignatureName())
	})

	// Login user
	win.DefineFunction("login", func(args ...*sciter.Value) *sciter.Value {
		err = win.LoadFile("rice://resources/loggingIn.html")
		if err != nil {
			log.Error().Err(err).Msgf("could not load loggingIn.html")
			return nil
		}

		go func() {
			err = operation.AuthorizationFlow(apiClient, storage)
			if err != nil {
				err = win.LoadFile("rice://resources/index.html")
				if err != nil {
					log.Error().Err(err).Msgf("could not load index.html")
				}

				log.Error().Err(err).Msg("could not complete authorization flow")
				return
			}

			apiClient.AddAuthorizationToken(storage.AuthorizationToken())
			err = win.LoadFile("rice://resources/main.html")
			if err != nil {
				log.Error().Err(err).Msgf("could not load main.html")
			}
		}()

		return nil
	})

	// Add a directory
	win.DefineFunction("add_directory", func(args ...*sciter.Value) *sciter.Value {
		folder := strings.Replace(args[0].String(), "file://", "", -1)
		if folder == "" {
			log.Error().Err(err).Msg("no folder was selected")
			return nil
		}

		folder, err := url.QueryUnescape(folder)
		if err != nil {
			log.Error().Err(err).Msg("could not decode folder path")
			return nil
		}

		err = storage.AddDirectory(folder)
		if err != nil {
			log.Error().Err(err).Msgf("could not add folder to storage")
			return nil
		}

		go watcher.Code(folder, func(event watcher.CodeEvent) {
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
		})

		err = win.LoadFile("rice://resources/main.html")
		if err != nil {
			log.Error().Err(err).Msgf("could not load loggingIn.html")
		}
		return nil
	})

	// Delete a directory
	win.DefineFunction("delete_directory", func(args ...*sciter.Value) *sciter.Value {
		err = storage.DeleteDirectory(args[0].String())
		if err != nil {
			log.Error().Err(err).Msgf("could not delete folder from storage")
		}

		err = win.LoadFile("rice://resources/main.html")
		if err != nil {
			log.Error().Err(err).Msgf("could not load loggingIn.html")
		}
		return nil
	})

	// Get all directories
	win.DefineFunction("directories", func(args ...*sciter.Value) *sciter.Value {
		directories, err := storage.Directories()
		if err != nil {
			log.Error().Err(err).Msgf("could not get folders")
			return nil
		}

		a := sciter.NewValue()
		for _, v := range directories {
			err = a.Append(v.Directory)
			if err != nil {
				log.Error().Err(err).Msg("could not append directory")
			}
		}

		return a
	})

	win.DefineFunction("log", func(args ...*sciter.Value) *sciter.Value {
		for _, arg := range args {
			log.Info().Interface("v", arg.String()).Msg("")
		}

		return nil
	})

	// Check if user is authenticated
	if storage.AuthorizationToken() != "" {
		// Load html
		err = win.LoadFile("rice://resources/main.html")
		if err != nil {
			log.Fatal().Err(err).Msg("could not load file")
		}
	} else {
		// Load html
		err = win.LoadFile("rice://resources/index.html")
		if err != nil {
			log.Fatal().Err(err).Msg("could not load file")
		}
	}

	return win
}
