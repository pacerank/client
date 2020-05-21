package store

import (
	"github.com/boltdb/bolt"
	"github.com/pacerank/client/pkg/system"
	"os"
	"path"
)

type Store struct {
	db *bolt.DB
}

func New() (*Store, error) {
	_, err := os.Stat(path.Join(system.HomePath(), ".pacerank"))
	if os.IsNotExist(err) {
		err = os.Mkdir(path.Join(system.HomePath(), ".pacerank"), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(path.Join(system.HomePath(), ".pacerank", "cache"), 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() {
	s.db.Close()
}
