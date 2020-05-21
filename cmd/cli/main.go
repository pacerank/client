package main

import (
	"github.com/pacerank/client/internal/operation"
	"github.com/pacerank/client/internal/store"
	"github.com/pacerank/client/internal/watcher"
	"github.com/pacerank/client/pkg/system"
	"github.com/rs/zerolog/log"
	"time"
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

	storage, err := store.New()
	if err != nil {
		log.Error().Err(err).Msgf("could start store")
		return
	}

	defer storage.Close()

	sys := system.New()
	go watcher.Keyboard(func(key watcher.KeyEvent) {
		process, err := sys.ActiveProcess()
		if err != nil {
			log.Error().Err(err).Msgf("could not get active process")
			return
		}

		log.Info().Msgf("process active: %s", process.Executable)
	})

	go watcher.Code("/home/kansuler/workspace", func(event watcher.CodeEvent) {
		if event.Err != nil {
			log.Error().Err(event.Err).Msg("could not watch code")
			return
		}

		log.Info().Str("language", event.Language).Str("filepath", event.FilePath).Str("filename", event.FileName).Msg("code watched")
	})

	for {
		time.Sleep(time.Second)
	}
}
