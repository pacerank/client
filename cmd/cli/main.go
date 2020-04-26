package main

import (
	"github.com/pacerank/client/internal/operation"
	"github.com/pacerank/client/pkg/system"
	"github.com/rs/zerolog/log"
	"runtime"
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

	sys, err := system.New(runtime.GOOS)
	if err != nil {
		log.Fatal().Str("os", runtime.GOOS).Err(err).Msg("could not create new system")
	}

	_, err = sys.Processes()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get list of processes")
	}

	log.Info().Msg("application started")
}
