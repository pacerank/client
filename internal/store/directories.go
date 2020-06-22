package store

import (
	"bytes"
	"github.com/boltdb/bolt"
)

type Directory struct {
	Id        string
	Directory string
}

func (s *Store) AddDirectory(directory string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("directories"))
		if err != nil {
			return err
		}

		id, _ := b.NextSequence()

		return b.Put(itob(id), []byte(directory))
	})
}

func (s *Store) Directories() ([]Directory, error) {
	var (
		result []Directory
		err    error
	)
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("directories"))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			result = append(result, Directory{
				Id:        string(k),
				Directory: string(v),
			})
		}

		return nil
	})

	return result, err
}

func (s *Store) DeleteDirectory(directory string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("directories"))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.Compare(v, []byte(directory)) == 0 {
				return b.Delete(k)
			}
		}

		return nil
	})
}
