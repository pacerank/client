package watcher

import (
	"errors"
	"fmt"
	"github.com/pacerank/client/internal/store"
	"github.com/pacerank/client/pkg/api"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func Queue(storage *store.Store, apiClient *api.Api) {
	for {
		time.Sleep(time.Second * 5)
		err := storage.NextInQueue(func(b []byte) error {
			result, _, err := apiClient.SendRecord(b)
			if err != nil {
				return err
			}

			if result.Status != http.StatusOK {
				return errors.New(fmt.Sprintf("did not acknowledge record in service with status %d and error: %s", result.Status, result.Error))
			}

			log.Info().Str("correlation-id", result.CorrelationId).Str("service", result.Service).Int("status", result.Status).Msg("pacerank service has acknowledged record")
			return nil
		})

		if err != nil {
			log.Error().Err(err).Msg("could not handle queue")
		}
	}
}
