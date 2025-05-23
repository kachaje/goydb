package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/kachaje/goydb/pkg/model"
	"github.com/kachaje/goydb/pkg/port"
)

func (d *Database) Changes(ctx context.Context, options *model.ChangesOptions) ([]*model.Document, int, error) {
	var pending int
	var docs []*model.Document
	wait := false

start:
	if options.SinceNow() || wait { // wait for new database changes
		wait := make(chan struct{})
		defer close(wait)
		t := time.AfterFunc(options.Timeout, func() { wait <- struct{}{} })
		err := d.AddListener(ctx, port.ChangeListenerFunc(func(ctx context.Context, doc *model.Document) error {
			wait <- struct{}{}
			options.Since = strconv.FormatInt(int64(doc.LocalSeq-1), 10)
			return context.Canceled // only wait for the next document
		}))
		if err != nil {
			return nil, 0, err
		}
		<-wait
		t.Stop()
	}

	err := d.Transaction(ctx, func(tx *Transaction) error {
		index := d.ChangesIndex()
		opts, err := index.IteratorOptions(ctx)
		if err != nil {
			return err
		}

		i := NewIterator(tx, WithOptions(opts))
		i.SetLimit(options.Limit)

		if !options.SinceNow() {
			i.SetStartKey([]byte(options.Since))
		}

		for doc := i.First(); i.Continue(); doc = i.Next() {
			docs = append(docs, doc)
		}

		// get number of remaining changes
		pending = i.Remaining()

		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	if len(docs) == 0 && options.Limit != 0 && !wait {
		wait = true
		goto start
	}

	return docs, pending, nil
}
