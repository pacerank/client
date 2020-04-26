package main

import (
	"github.com/pacerank/client/internal/operation"
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

	start := time.Now()

	sys := system.New()

	processes, err := sys.Processes()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get list of processes")
	}

	log.Debug().Str("duration", time.Now().Sub(start).String()).Msg("duration for processing processes")

	for _, process := range processes {
		log.Debug().
			Int64("pid", process.Pid).
			Int64("ppid", process.Ppid).
			Ints64("modules", process.ModulePid).
			Str("executable", process.FileName).
			Str("checksum", process.Checksum).
			Msg("")
	}
}
