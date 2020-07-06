package store

import (
	"github.com/boltdb/bolt"
	"os"
	"os/user"
	"path"
)

type Store struct {
	db                      *bolt.DB
	queueRecord             []byte
	id                      []byte
	NotifyListenToDirectory chan string
}

func New() (*Store, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(path.Join(usr.HomeDir, ".pacerank"))
	if os.IsNotExist(err) {
		err = os.Mkdir(path.Join(usr.HomeDir, ".pacerank"), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(path.Join(usr.HomeDir, ".pacerank", "cache"), 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Store{
		db:                      db,
		NotifyListenToDirectory: make(chan string),
	}, nil
}

func (s *Store) Close() {
	_ = s.db.Close()
}
