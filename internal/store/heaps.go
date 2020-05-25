package store

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
)

func (s *Store) Heaps() ([]string, error) {
	var (
		result []string
		err    error
	)

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("heaps"))
		if b == nil {
			return errors.New("no heaps are active")
		}

		return b.ForEach(func(k, v []byte) error {
			pb := b.Bucket(k)
			if pb != nil {
				result = append(result, string(k))
			}
			return nil
		})
	})

	return result, err
}

type Beat struct {
	Id       string
	Language string
	Branch   string
	FileName string
	Project  string
	Git      string
}

func (s *Store) AddHeap(beat Beat) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("heaps"))
		if err != nil {
			return err
		}

		pb := b.Bucket([]byte(beat.Id))
		if pb == nil {
			pb, err = b.CreateBucketIfNotExists([]byte(beat.Id))
			if err != nil {
				return err
			}
		}

		err = pb.Put([]byte("languages"), appendToBytes(pb.Get([]byte("languages")), beat.Language))
		if err != nil {
			return err
		}

		err = pb.Put([]byte("files"), appendToBytes(pb.Get([]byte("files")), beat.FileName))
		if err != nil {
			return err
		}

		err = pb.Put([]byte("project"), []byte(beat.Project))
		if err != nil {
			return err
		}

		err = pb.Put([]byte("git"), []byte(beat.Git))
		if err != nil {
			return err
		}

		err = pb.Put([]byte("branch"), []byte(beat.Branch))
		if err != nil {
			return err
		}

		return nil
	})
}

type Heap struct {
	Id        string
	Languages []string
	Files     []string
	Project   string
	Git       string
	Branch    string
}

func (s *Store) HeapByProjectId(id string) (Heap, error) {
	var (
		heap Heap
		err  error
	)

	heap.Id = id

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("heaps"))
		if b == nil {
			return errors.New("heaps has not been initiated")
		}

		pb := b.Bucket([]byte(id))
		if pb == nil {
			return errors.New("heap with project id does not exist")
		}

		heap.Project = string(pb.Get([]byte("project")))
		heap.Git = string(pb.Get([]byte("git")))
		heap.Branch = string(pb.Get([]byte("branch")))
		err = json.NewDecoder(bytes.NewBuffer(pb.Get([]byte("files")))).Decode(&heap.Files)
		if err != nil {
			return err
		}

		return json.NewDecoder(bytes.NewBuffer(pb.Get([]byte("languages")))).Decode(&heap.Languages)
	})

	return heap, err
}

func (s *Store) ClearHeaps() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte("heaps"))
	})
}

// returns a new byte array with appended value if it doesn't exist
func appendToBytes(b []byte, value string) []byte {
	var strs []string
	var exists bool
	if b != nil {
		_ = json.NewDecoder(bytes.NewBuffer(b)).Decode(&strs)
		for _, v := range strs {
			if v == value {
				exists = true
			}
		}
	}

	if !exists {
		strs = append(strs, value)
	}

	result, _ := json.Marshal(strs)
	return result
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
