package watcher

import (
	"github.com/pacerank/client/internal/store"
	"github.com/rs/zerolog/log"
	"time"
)

func Heaps(storage *store.Store) {
	// Infinite loop that checks if any time thresholds has been met
	for {
		time.Sleep(time.Second * 5)

		meta, err := storage.CurrentMeta()
		if err != nil {
			log.Error().Err(err).Msg("could not get current meta")
			continue
		}

		if err == store.ErrNoActivity {
			continue
		}

		log.Debug().Time("last_editor_activity", meta.LastEditorActivity).Time("first_editor_activity", meta.FirstEditorActivity).Uint64("keypress_count", meta.KeypressCount).Str("session_id", meta.SessionId).Msg("did a check on meta and heaps")

		// If last editor activity was more than 3 minutes ago, queue all heaps
		// and clean the current heap stack
		if time.Now().Add(-time.Minute * 3).After(meta.LastEditorActivity) {
			heaps, err := storage.Heaps()
			if err != nil {
				log.Error().Err(err).Msg("couldn't get any heaps")
				continue
			}

			// If no heaps, continue
			if len(heaps) == 0 {
				continue
			}

			log.Debug().Strs("heaps", heaps).Msg("activity has passed - TODO: Send activity to digest")
		}

		// If last editor activity was more than 30 minutes ago, start a new session
		if time.Now().Add(-time.Minute * 30).After(meta.LastEditorActivity) {
			err = storage.NewSession()
			if err != nil {
				log.Error().Err(err).Msg("could not start a new session")
			}

			log.Debug().Msg("30 minutes has passed, clear session")
		}
	}

}
