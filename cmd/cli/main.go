package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/pacerank/client/internal/operation"
	"github.com/pacerank/client/internal/store"
	"github.com/pacerank/client/internal/watcher"
	"github.com/pacerank/client/pkg/api"
	"github.com/pacerank/client/pkg/system"
	"github.com/rs/zerolog/log"
	"time"
)

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

	apiClient := api.New("https://digest.development.pacerank.io")

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

	apiClient.AddAuthorizationToken(token)

	sys := system.New()
	go watcher.Keyboard(func(key watcher.KeyEvent) {
		process, err := sys.ActiveProcess()
		if err != nil {
			log.Error().Err(err).Msgf("could not get active process")
			return
		}

		log.Debug().Msgf("process active: %s", process.Executable)
	})

	for _, folder := range opts.Folders {
		go watcher.Code(folder, func(event watcher.CodeEvent) {
			if event.Err != nil {
				log.Error().Err(event.Err).Msg("could not watch code")
				return
			}

			log.Debug().Str("language", event.Language).Str("filepath", event.FilePath).Str("filename", event.FileName).Msg("code watched")
		})
	}

	for {
		time.Sleep(time.Second)
	}
}
