package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/pacerank/client/internal/inspect"
	"github.com/pacerank/client/internal/operation"
	"github.com/pacerank/client/internal/store"
	"github.com/pacerank/client/internal/watcher"
	"github.com/pacerank/client/pkg/api"
	"github.com/pacerank/client/pkg/system"
	"github.com/rs/zerolog/log"
	"time"
)

// Possible flags
type Options struct {
	Verbose bool     `short:"v" long:"verbose" description:"Show verbose debug information"`
	Folders []string `short:"f" long:"folders" description:"Folders to watch for file changes"`
}

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

	var opts Options
	_, err = flags.Parse(&opts)
	if err != nil {
		return
	}

	if opts.Verbose {
		operation.EnableDebug()
	}

	if len(opts.Folders) == 0 {
		log.Fatal().Msg("must give at least one folder to watch")
	}

	apiClient := api.New("https://digest.pacerank.io")

	storage, err := store.New()
	if err != nil {
		log.Fatal().Err(err).Msgf("could start store")
	}

	defer storage.Close()

	token := storage.AuthorizationToken()
	for token == "" {
		log.Info().Msg("api key does not exist, initialize authorization flow")
		err = operation.AuthorizationFlow(apiClient, storage)
		if err != nil {
			log.Fatal().Err(err).Msg("could not complete authorization flow")
		}

		token = storage.AuthorizationToken()
	}

	log.Info().Msgf("Hello %s", storage.UserSignatureName())

	// Setup a new session and meta
	err = storage.NewSession()

	apiClient.AddAuthorizationToken(token)

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

	for _, folder := range opts.Folders {
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

			log.Debug().
				Str("language", event.Language).
				Str("filepath", event.FilePath).
				Str("filename", event.FileName).
				Str("project", event.Project).
				Str("branch", event.Branch).
				Str("git", event.Git).
				Str("id", event.Id).
				Msg("code change found")
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

	for {
		time.Sleep(time.Second)
	}
}
