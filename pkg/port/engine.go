package port

import (
	"errors"

	"github.com/kachaje/goydb/pkg/model"
)

var ErrNotFound = errors.New("resource not found")
var ErrConflict = errors.New("rev doesn't match for update")

type DatabaseEngine interface {
	Stats() (stats model.DatabaseStats, err error)
	ReadTransaction(fn func(tx EngineReadTransaction) error) error
	WriteTransaction(fn func(tx EngineWriteTransaction) error) error
	Close() error
}

// KeyWithSeq should return a new key based on the given
// key and a sequence. The function may return a new key or new
// data. If the returned data is nil, the original data is used.
type KeyWithSeq func(key, value []byte, seq uint64) (newKey []byte, newValue []byte)

type EngineWriteTransaction interface {
	EnsureBucket(bucket []byte)
	DeleteBucket(bucket []byte)
	Put(bucket, k, v []byte)
	// PutWithSequence will get the next sequence for the bucket
	// and then call the fn func using the passed key and seq to
	// generate the final key
	PutWithSequence(bucket, k, v []byte, fn KeyWithSeq)
	// PutWithReusedSequence will perform the same operation as
	// PutWithSequence but reused the generated sequence for
	// PutWithSequence
	PutWithReusedSequence(bucket, k, v []byte, fn KeyWithSeq)
	Delete(bucket, k []byte)
	EngineReadTransaction
}

type EngineReadTransaction interface {
	BucketStats(bucket []byte) *model.IndexStats
	Cursor(bucket []byte) EngineCursor
	Get(bucket, key []byte) ([]byte, error)
	Sequence(bucket []byte) uint64
}

type EngineCursor interface {
	First() (key []byte, value []byte)
	Last() (key []byte, value []byte)
	Next() (key []byte, value []byte)
	Prev() (key []byte, value []byte)
	Seek(seek []byte) (key []byte, value []byte)
}
