package watcher

import (
	"encoding/json"
	"github.com/pacerank/client/internal/store"
	"github.com/pacerank/client/pkg/model"
	"github.com/rs/zerolog/log"
	"strconv"
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
			continue
		}

		// Send heap to queue if the last activity surpass x minutes and hasn't been sent already.
		if time.Now().Add(-time.Minute*3).After(meta.LastActivity) && !meta.HeapAddedToQueue {
			heaps, err := storage.Heaps()
			if err != nil {
				log.Error().Err(err).Msg("couldn't get any heaps")
				continue
			}

			// If no heaps, continue
			if len(heaps) == 0 {
				continue
			}

			for _, heapId := range heaps {
				heap, err := storage.HeapById(heapId)
				if err != nil {
					log.Error().Err(err).Msgf("couldn't get heap by id %s", heapId)
					break
				}

				record := model.Record{
					Start:     meta.FirstActivity,
					Stop:      meta.LastActivity,
					Activity:  model.ActivityCoding,
					SessionId: meta.SessionId,
					Labels: []model.Label{
						{
							Category: model.CategoryProject,
							Value:    heap.Project,
						},
						{
							Category: model.CategoryBranch,
							Value:    heap.Branch,
						},
						{
							Category: model.CategoryGit,
							Value:    heap.Git,
						},
						{
							Category: model.CategoryKeyCount,
							Value:    strconv.FormatUint(meta.KeypressCount, 10),
						},
					},
				}

				for _, file := range heap.Files {
					record.Labels = append(record.Labels, model.Label{
						Category: model.CategoryFilename,
						Value:    file,
					})
				}

				for _, language := range heap.Languages {
					record.Labels = append(record.Labels, model.Label{
						Category: model.CategoryLanguage,
						Value:    language,
					})
				}

				for _, editor := range meta.Editors {
					record.Labels = append(record.Labels, model.Label{
						Category: model.CategoryEditor,
						Value:    editor,
					})
				}

				b, err := json.Marshal(record)
				if err != nil {
					log.Error().Err(err).Msg("could not marshal record into byte array")
					break
				}

				err = storage.AddToQueue(b)
				if err != nil {
					log.Error().Err(err).Msg("could not add record to queue")
					break
				}

				err = storage.DeleteHeap(heapId)
				if err != nil {
					log.Error().Err(err).Msg("could not delete heap by id")
					break
				}
			}

			// Updated so that heap is added to queue (Skip queueing heap data that has already been queued)
			err = storage.SentToQueue()
			if err != nil {
				log.Error().Err(err).Msg("could not update heap_added_to_queue state")
				continue
			}
		}

		// If last editor activity was more than 30 minutes ago, start a new session
		if time.Now().Add(-time.Minute * 30).After(meta.LastActivity) {
			err = storage.NewSession()
			if err != nil {
				log.Error().Err(err).Msg("could not start a new session")
			}

			log.Info().Msg("30 minutes has passed, clear session")
		}
	}
}
