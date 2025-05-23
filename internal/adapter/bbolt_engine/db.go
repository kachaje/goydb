package bbolt_engine

import (
	"bytes"
	"os"

	"github.com/kachaje/goydb/internal/adapter/index"
	"github.com/kachaje/goydb/pkg/model"
	"github.com/kachaje/goydb/pkg/port"
	"go.etcd.io/bbolt"
)

var _ port.DatabaseEngine = (*DB)(nil)

type DB struct {
	db *bbolt.DB
}

func Open(path string) (*DB, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &DB{
		db: db,
	}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) ReadTransaction(fn func(tx port.EngineReadTransaction) error) error {
	return db.db.View(func(btx *bbolt.Tx) error {
		return fn(NewReadTransaction(btx))
	})
}

// WriteTransaction executes the given function in a read transaction
// that collects all database updates into an operation log that will
// be executed at the end of the transaction execution as one transaction.
// This method is designed to allow more concurrent write transactions due
// to less time spend waiting between the different write operations.
// If no writes are made, the update transaction is omitted.
func (db *DB) WriteTransaction(fn func(tx port.EngineWriteTransaction) error) error {
	var wtx *WriteTransaction
	err := db.db.View(func(btx *bbolt.Tx) error {
		wtx = NewWriteTransaction(btx)
		return fn(wtx)
	})
	if err != nil {
		return err
	}

	// only attempt the update transaction if there is something to do
	if len(wtx.opLog) > 0 {
		return db.db.Update(func(btx *bbolt.Tx) error {
			return wtx.Commit(btx)
		})
	}

	return nil
}

func (db *DB) Stats() (stats model.DatabaseStats, err error) {
	fi, err := os.Stat(db.db.Path())
	if err != nil {
		return stats, err
	}
	stats.FileSize = uint64(fi.Size())
	err = db.db.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			s := b.Stats()
			// only take the doc count from the docs bucket
			if bytes.Equal(name, model.DocsBucket) {
				stats.DocCount += uint64(s.KeyN)
			}
			// if deleted index
			if bytes.Equal(name, []byte(index.DeletedIndexName)) {
				stats.DocCount -= uint64(s.KeyN)
				stats.DocDelCount = uint64(s.KeyN)
			}

			// accumulate all numbers to have accurate database statistics
			stats.Alloc += uint64(s.BranchAlloc + s.LeafAlloc)
			stats.InUse += uint64(s.BranchInuse + s.LeafInuse)
			return nil
		})
	})
	return
}
