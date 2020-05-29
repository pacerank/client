package store

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
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

type InHeap struct {
	Id       string
	Language string
	Branch   string
	FileName string
	Project  string
	Git      string
}

func (s *Store) AddHeap(heap InHeap) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("heaps"))
		if err != nil {
			return err
		}

		pb := b.Bucket([]byte(heap.Id))
		if pb == nil {
			pb, err = b.CreateBucketIfNotExists([]byte(heap.Id))
			if err != nil {
				return err
			}
		}

		err = pb.Put([]byte("languages"), appendToBytes(pb.Get([]byte("languages")), heap.Language))
		if err != nil {
			return err
		}

		err = pb.Put([]byte("files"), appendToBytes(pb.Get([]byte("files")), heap.FileName))
		if err != nil {
			return err
		}

		err = pb.Put([]byte("project"), []byte(heap.Project))
		if err != nil {
			return err
		}

		err = pb.Put([]byte("git"), []byte(heap.Git))
		if err != nil {
			return err
		}

		err = pb.Put([]byte("branch"), []byte(heap.Branch))
		if err != nil {
			return err
		}

		return nil
	})
}

type OutHeap struct {
	Id        string
	Languages []string
	Files     []string
	Project   string
	Git       string
	Branch    string
}

func (s *Store) HeapById(id string) (OutHeap, error) {
	var (
		heap OutHeap
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
			return errors.New(fmt.Sprintf("heap with id %s does not exist", id))
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

func (s *Store) DeleteHeap(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("heaps"))
		if b == nil {
			return errors.New("heaps has not been initiated")
		}

		pb := b.Bucket([]byte(id))
		if pb == nil {
			return errors.New(fmt.Sprintf("heap with id %s does not exist", id))
		}

		return b.DeleteBucket([]byte(id))
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
