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

	for {
		start := time.Now()

		sys := system.New()

		process, err := sys.ActiveProcess()
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

		log.Debug().Str("duration", time.Now().Sub(start).String()).Int64("pid", process.ProcessID).Ints64("children", process.Children).Str("checksum", process.Checksum).Str("name", process.FileName).Msg("")

		time.Sleep(time.Second * 2)
	}
}
