package main

import (
	"github.com/pacerank/client/internal/operation"
	"github.com/pacerank/client/pkg/keyboard"
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

	sys := system.New()

	go keyboard.Listen(func(key keyboard.KeyEvent) {
		process, err := sys.ActiveProcess()
		if err != nil {
			log.Error().Err(err).Msgf("could not get active process")
			return
		}

		log.Info().Msgf("process active: %s", process.Executable)
	})

	for {
		time.Sleep(time.Second)
	}
}
