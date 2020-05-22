package store

import (
	"errors"
	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
)

// Get the authorization token from store
func (s *Store) AuthorizationToken() string {
	var result string
	_ = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("settings"))
		if b == nil {
			return nil
		}

		v := b.Get([]byte("authorization_token"))
		if v == nil {
			return nil
		}

		result = string(v)
		return nil
	})

	return result
}

// Set authorization token in store
func (s *Store) SetAuthorizationToken(apiKey string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("settings"))
		if err != nil {
			return err
		}
		return b.Put([]byte("authorization_token"), []byte(apiKey))
	})
}

// Save the authorization flow
func (s *Store) SaveAuthorizationFlow(id, link string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("authorization_flow"))
		if err != nil {
			return err
		}

		err = b.Put([]byte("authorization_id"), []byte(id))
		if err != nil {
			return err
		}

		err = b.Put([]byte("authorization_link"), []byte(link))
		if err != nil {
			return err
		}

		return nil
	})
}

// Get values for the authorization flow
func (s *Store) AuthorizationFlow() (string, string, error) {
	var (
		authorizationId   []byte
		authorizationLink []byte
		err               error
	)

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("authorization_flow"))
		authorizationId = b.Get([]byte("authorization_id"))
		authorizationLink = b.Get([]byte("authorization_link"))
		return nil
	})

	if authorizationId == nil {
		return "", "", errors.New("authorization_id does not exist in store")
	}

	if authorizationLink == nil {
		return "", "", errors.New("authorization_link does not exist in store")
	}

	return string(authorizationId), string(authorizationLink), err
}

// Get the user signature name from store
func (s *Store) UserSignatureName() string {
	var result string
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("settings"))
		result = string(b.Get([]byte("user_signature_name")))
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("could not get user_signature_name")
	}

	return result
}

// Set the signature username in store
func (s *Store) SetUserSignatureName(userSignatureName string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("settings"))
		if err != nil {
			return err
		}
		return b.Put([]byte("user_signature_name"), []byte(userSignatureName))
	})
}
