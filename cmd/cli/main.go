package main

import (
	"github.com/pacerank/client/internal/operation"
	"github.com/pacerank/client/pkg/system"
	"github.com/rs/zerolog/log"
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

	_, err = system.New(system.Linux)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create new system")
	}

	log.Info().Msg("application started")
}
