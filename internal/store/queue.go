package store

import "github.com/boltdb/bolt"

// Add a json in byte format to queue
func (s *Store) AddToQueue(record []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("queue"))
		if err != nil {
			return err
		}

		id, _ := b.NextSequence()
		return b.Put(itob(id), record)
	})
}

// Send in a callback that receives payload, if no error occur in callback
// delete payload from queue
func (s *Store) NextInQueue(Callback func([]byte) error) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("queue"))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := Callback(v)
			if err != nil {
				return err
			}

			err = b.Delete(k)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
