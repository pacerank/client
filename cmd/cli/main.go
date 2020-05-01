package main

import (
	"github.com/pacerank/client/internal/operation"
	"github.com/pacerank/client/pkg/keyboard"
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

	go keyboard.Listen(func(key keyboard.KeyEvent) {
		log.Info().Int16("key", int16(key.Key)).Err(err).Msgf("rune %q", key.Rune)
	})

	for {
		time.Sleep(time.Second)
	}
}
