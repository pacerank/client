package store

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"
)

var ErrNoActivity = errors.New("no activity exist currently")

// Add a new session, this happens when the meta and heaps should all
// be reset into a new session.
func (s *Store) NewSession() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		var err error

		_ = tx.DeleteBucket([]byte("heaps"))
		_ = tx.DeleteBucket([]byte("meta"))

		b := tx.Bucket([]byte("meta"))
		if b == nil {
			b, err = tx.CreateBucketIfNotExists([]byte("meta"))
			if err != nil {
				return err
			}
		}

		return b.Put([]byte("session_id"), []byte(uuid.NewV4().String()))
	})
}

// This function will react on keypress and update the meta state
func (s *Store) MetaTypingActivity(editor string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("meta"))
		if b == nil {
			return errors.New("meta is not initialized")
		}

		if b.Get([]byte("first_activity")) == nil {
			err := b.Put([]byte("first_activity"), []byte(time.Now().Format(time.RFC3339)))
			if err != nil {
				return err
			}
		}

		var count uint64
		v := b.Get([]byte("keypress_count"))
		if v == nil {
			count = 0
		} else {
			count = binary.BigEndian.Uint64(v)
		}

		// Increment keypress
		err := b.Put([]byte("keypress_count"), itob(count+1))
		if err != nil {
			return err
		}

		err = b.Put([]byte("editors"), appendToBytes(b.Get([]byte("editors")), editor))
		if err != nil {
			return err
		}

		err = b.Put([]byte("heap_added_to_queue"), []byte("false"))
		if err != nil {
			return err
		}

		return b.Put([]byte("last_activity"), []byte(time.Now().Format(time.RFC3339)))
	})
}

// Update heap_added_to_queue status, this gets called when heap has been added to queue
// and it shouldn't happen again until next activity. It also reset first_activity
func (s *Store) SentToQueue() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("meta"))
		if b == nil {
			return errors.New("meta is not initialized")
		}

		err := b.Delete([]byte("first_activity"))
		if err != nil {
			return err
		}

		return b.Put([]byte("heap_added_to_queue"), []byte("true"))
	})
}

type Meta struct {
	SessionId        string
	KeypressCount    uint64
	Editors          []string
	HeapAddedToQueue bool
	FirstActivity    time.Time
	LastActivity     time.Time
}

func (s *Store) Meta() (Meta, error) {
	var (
		result Meta
		err    error
	)

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("meta"))
		if b == nil {
			return errors.New("meta is not initialized")
		}

		result.SessionId = string(b.Get([]byte("session_id")))
		kpc := b.Get([]byte("keypress_count"))
		if kpc != nil {
			result.KeypressCount = binary.BigEndian.Uint64(kpc)
		}
		result.HeapAddedToQueue, _ = strconv.ParseBool(string(b.Get([]byte("heap_added_to_queue"))))

		v := b.Get([]byte("editors"))
		if v != nil {
			err = json.NewDecoder(bytes.NewBuffer(v)).Decode(&result.Editors)
			if err != nil {
				return err
			}
		}

		v = b.Get([]byte("first_activity"))
		if v != nil {
			result.FirstActivity, err = time.Parse(time.RFC3339, string(v))
			if err != nil {
				return err
			}
		}

		v = b.Get([]byte("last_activity"))
		if v == nil {
			return ErrNoActivity
		}

		result.LastActivity, err = time.Parse(time.RFC3339, string(v))
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
