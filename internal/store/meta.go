package store

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
	"time"
)

var ErrNoActivity = errors.New("no activity exist currently")

func (s *Store) NewSession() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		var err error

		_ = tx.DeleteBucket([]byte("heaps"))
		_ = tx.DeleteBucket([]byte("meta"))

		b := tx.Bucket([]byte("meta"))
		if b == nil {
			b, err = tx.CreateBucketIfNotExists([]byte("meta"))
		}

		if err != nil {
			return err
		}

		err = b.Put([]byte("session_id"), []byte(uuid.NewV4().String()))
		if err != nil {
			return err
		}

		return b.Put([]byte("keypress_count"), itob(0))
	})
}

func (s *Store) SessionId() (string, error) {
	var (
		result string
		err    error
	)

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("meta"))
		if b == nil {
			return errors.New("meta is not initialized")
		}

		v := b.Get([]byte("session_id"))
		if v == nil {
			return errors.New("meta is not initialized")
		}

		result = string(v)
		return nil
	})

	return result, err
}

func (s *Store) KeyPressCount() (uint64, error) {
	var (
		result uint64
		err    error
	)

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("meta"))
		if b == nil {
			return errors.New("meta is not initialized")
		}

		v := b.Get([]byte("keypress_count"))
		if v == nil {
			return errors.New("meta is not initialized")
		}

		result = binary.BigEndian.Uint64(v)
		return nil
	})

	return result, err
}

func (s *Store) SetKeyPressCount(count uint64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("meta"))
		if b == nil {
			return errors.New("meta is not initialized")
		}

		return b.Put([]byte("keypress_count"), itob(count))
	})
}

func (s *Store) UpdateTypingActivity(editor string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("meta"))
		if b == nil {
			return errors.New("meta is not initialized")
		}

		if b.Get([]byte("first_editor_activity")) == nil {
			err := b.Put([]byte("first_editor_activity"), []byte(time.Now().Format(time.RFC3339)))
			if err != nil {
				return err
			}
		}

		v := b.Get([]byte("keypress_count"))
		count := binary.BigEndian.Uint64(v)

		// Increment keypress
		err := b.Put([]byte("keypress_count"), itob(count+1))
		if err != nil {
			return err
		}

		err = b.Put([]byte("editors"), appendToBytes(b.Get([]byte("editors")), editor))
		if err != nil {
			return err
		}

		return b.Put([]byte("last_editor_activity"), []byte(time.Now().Format(time.RFC3339)))
	})
}

type Meta struct {
	SessionId           string
	KeypressCount       uint64
	Editors             []string
	FirstEditorActivity time.Time
	LastEditorActivity  time.Time
}

func (s *Store) CurrentMeta() (Meta, error) {
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
		result.KeypressCount = binary.BigEndian.Uint64(b.Get([]byte("keypress_count")))
		v := b.Get([]byte("editors"))
		if v != nil {
			err = json.NewDecoder(bytes.NewBuffer(v)).Decode(&result.Editors)
			if err != nil {
				return err
			}
		}

		v = b.Get([]byte("first_editor_activity"))
		if v == nil {
			return ErrNoActivity
		}

		result.FirstEditorActivity, err = time.Parse(time.RFC3339, string(v))
		if err != nil {
			return err
		}

		v = b.Get([]byte("last_editor_activity"))
		if v == nil {
			return ErrNoActivity
		}

		result.LastEditorActivity, err = time.Parse(time.RFC3339, string(v))
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
