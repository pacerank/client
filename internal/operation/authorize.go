package operation

import (
	"fmt"
	"github.com/pacerank/client/internal/store"
	"github.com/pacerank/client/pkg/api"
	"github.com/pacerank/client/pkg/system"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type authorizeError struct {
	Message string
	Err     error
}

func (ae authorizeError) Unwrap() error {
	return ae.Err
}

func (ae authorizeError) Error() string {
	return fmt.Sprintf("%s: %s", ae.Message, ae.Err.Error())
}

// This endpoint will start an authorization flow, and check against api
// until the user is authorized. This is a blocking function.
func AuthorizationFlow(apiClient *api.Api, storage *store.Store) error {
	defaultReply, endpointReply, err := apiClient.InitializeAuthorizationFlow()
	if err != nil {
		return authorizeError{
			Message: "could not initialize authorization flow",
			Err:     err,
		}
	}

	if defaultReply.Status != http.StatusOK {
		return authorizeError{
			Message: fmt.Sprintf("unexpected server error occurred, returned status %d", defaultReply.Status),
			Err:     err,
		}
	}

	err = storage.SaveAuthorizationFlow(endpointReply.AuthorizationId, endpointReply.AuthorizationUrl)
	if err != nil {
		return authorizeError{
			Message: "could not save authorization flow values",
			Err:     err,
		}
	}

	err = system.OpenBrowser(endpointReply.AuthorizationUrl)
	if err != nil {
		log.Info().Msgf("could not open browser automatically, please authorize this client here: %s", endpointReply.AuthorizationUrl)
	}

	for {
		time.Sleep(time.Second * 5)

		id, _, err := storage.AuthorizationFlow()
		if err != nil {
			return authorizeError{
				Message: "could not retrieve authorization flow from store",
				Err:     err,
			}
		}

		defaultReply, endpointReply, err := apiClient.ConfirmAuthorizationFlow(id)
		if defaultReply.Status != http.StatusOK {
			return authorizeError{
				Message: fmt.Sprintf("unexpected status code from api: %d", defaultReply.Status),
				Err:     err,
			}
		}

		if endpointReply.Authorized {
			log.Info().Msgf("user has authorized client: %s", endpointReply.ApiKey)
			err = storage.SetAuthorizationToken(endpointReply.ApiKey)
			if err != nil {
				return authorizeError{
					Message: "could not set authorization token in store",
					Err:     err,
				}
			}

			err = storage.SetUserSignatureName(endpointReply.UserSignatureName)
			if err != nil {
				return authorizeError{
					Message: "could not set user signature name in store",
					Err:     err,
				}
			}

			break
		}
	}

	return nil
}
