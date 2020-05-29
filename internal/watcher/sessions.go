package watcher

import (
	"github.com/pacerank/client/internal/store"
	"github.com/rs/zerolog/log"
	"time"
)

// Function is used to see if a session has passed threshold times
// if it has, it will add the session to the send queue and update
// the state with new data.
func Sessions(storage *store.Store) {
	// Infinite loop that checks if any time thresholds has been met
	for {
		time.Sleep(time.Second * 5)

		meta, err := storage.Meta()
		if err != nil {
			log.Error().Err(err).Msg("could not get current meta")
			continue
		}

		if err == store.ErrNoActivity {
			continue
		}

		log.Debug().Interface("meta", meta).Msg("did a check on meta and heaps")

		// Send heap to queue if the last activity surpass x minutes and hasn't been sent already.
		if time.Now().Add(-time.Minute*1).After(meta.LastActivity) && !meta.HeapAddedToQueue {
			heaps, err := storage.Heaps()
			if err != nil {
				log.Error().Err(err).Msg("couldn't get any heaps")
				continue
			}

			// Updated so that heap is added to queue (Skip queueing heap data that has already been queued)
			err = storage.SentToQueue()
			if err != nil {
				log.Error().Err(err).Msg("could not update heap_added_to_queue state")
				continue
			}

			// If no heaps, continue
			if len(heaps) == 0 {
				continue
			}

			log.Debug().Strs("heaps", heaps).Msg("activity has passed - TODO: Send activity to digest")
		}

		// If last editor activity was more than 30 minutes ago, start a new session
		if time.Now().Add(-time.Minute * 30).After(meta.LastActivity) {
			err = storage.NewSession()
			if err != nil {
				log.Error().Err(err).Msg("could not start a new session")
			}

			log.Debug().Msg("30 minutes has passed, clear session")
		}
	}

}
