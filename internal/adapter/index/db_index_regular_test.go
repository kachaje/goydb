package index_test

import (
	"context"
	"testing"

	"github.com/kachaje/goydb/internal/adapter/index"
	"github.com/kachaje/goydb/internal/adapter/storage"
	"github.com/kachaje/goydb/pkg/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestRegularIndex(t *testing.T) {
	WithTestDatabase(t, func(ctx context.Context, db *storage.Database) {
		ri := index.NewRegularIndex(
			&model.DesignDocFn{
				Type:        model.ViewFn,
				DesignDocID: "doc",
				FnName:      "fn",
			},
			func(ctx context.Context, doc *model.Document) ([][]byte, [][]byte) {
				var keys, values [][]byte

				// index all values
				for k, v := range doc.Data {
					keys = append(keys, []byte(k))
					out, err := bson.Marshal(model.Document{
						ID:    doc.ID,
						Value: v,
					})
					assert.NoError(t, err)
					if err == nil {
						values = append(values, out)
					}
				}

				return keys, values
			},
		)
		err := db.Transaction(ctx, func(tx *storage.Transaction) error {
			err := ri.Ensure(ctx, tx)
			assert.NoError(t, err)

			t.Run("delete on unknown document", func(t *testing.T) {
				err := ri.DocumentDeleted(ctx, tx, &model.Document{ID: "unknown"})
				assert.NoError(t, err)
			})

			t.Run("stats with no documents", func(t *testing.T) {
				stats, err := ri.Stats(ctx, tx)
				assert.NoError(t, err)
				assert.Equal(t, uint64(0), stats.Documents)
				assert.Equal(t, uint64(0), stats.Keys)
			})

			t.Run("iterator with no documents", func(t *testing.T) {
				iter, err := db.IndexIterator(ctx, tx, ri)
				assert.NoError(t, err)
				assert.Equal(t, 0, iter.Remaining())
			})

			t.Run("with documents", func(t *testing.T) {
				err := ri.DocumentStored(ctx, tx, &model.Document{
					ID: "test",
					Data: map[string]any{
						"name": "Foo",
						"test": 123,
					},
				})
				assert.NoError(t, err)
				// same record twice
				err = ri.DocumentStored(ctx, tx, &model.Document{
					ID: "test",
					Data: map[string]any{
						"name": "Foo",
						"test": 123,
					},
				})
				assert.NoError(t, err)
				err = ri.DocumentStored(ctx, tx, &model.Document{
					ID: "test1",
					Data: map[string]any{
						"name": "Foo",
						"test": 234,
					},
				})
				assert.NoError(t, err)
			})

			return nil
		})
		assert.NoError(t, err)

		err = db.Transaction(ctx, func(tx *storage.Transaction) error {
			t.Run("iterator", func(t *testing.T) {
				// FIXME: fix test
				t.SkipNow()
				iter, err := db.IndexIterator(ctx, tx, ri)
				assert.NoError(t, err)
				var docs []*model.Document
				for doc := iter.First(); iter.Continue(); doc = iter.Next() {
					docs = append(docs, doc)
				}
				assert.EqualValues(t, []*model.Document{
					&model.Document{
						ID:          "test",
						Key:         "name",
						Value:       "Foo",
						Data:        map[string]any{},
						Attachments: map[string]*model.Attachment{},
					},
					&model.Document{
						ID:          "test1",
						Key:         "name",
						Value:       "Foo",
						Data:        map[string]any{},
						Attachments: map[string]*model.Attachment{},
					},
					&model.Document{
						ID:          "test",
						Key:         "test",
						Value:       int(123),
						Data:        map[string]any{},
						Attachments: map[string]*model.Attachment{},
					},
					&model.Document{
						ID:          "test1",
						Key:         "test",
						Value:       int(234),
						Data:        map[string]any{},
						Attachments: map[string]*model.Attachment{},
					},
				}, docs)
			})

			t.Run("stats", func(t *testing.T) {
				// FIXME: fix test
				t.SkipNow()
				stats, err := ri.Stats(ctx, tx)
				assert.NoError(t, err)
				assert.Equal(t, uint64(2), stats.Documents)
				assert.Equal(t, uint64(4), stats.Keys)
			})

			t.Run("document removed", func(t *testing.T) {
				// FIXME: fix test
				t.SkipNow()
				err := ri.DocumentDeleted(ctx, tx, &model.Document{ID: "test"})
				assert.NoError(t, err)

				t.Run("iterator", func(t *testing.T) {
					iter, err := db.IndexIterator(ctx, tx, ri)
					assert.NoError(t, err)
					var docs []*model.Document
					for doc := iter.First(); iter.Continue(); doc = iter.Next() {
						docs = append(docs, doc)
					}
					assert.EqualValues(t, []*model.Document{
						&model.Document{
							ID:          "test1",
							Key:         "name",
							Value:       "Foo",
							Data:        map[string]any{},
							Attachments: map[string]*model.Attachment{},
						},
						&model.Document{
							ID:          "test1",
							Key:         "test",
							Value:       int(234),
							Data:        map[string]any{},
							Attachments: map[string]*model.Attachment{},
						},
					}, docs)
				})
			})

			return ri.Remove(ctx, tx)
		})
		assert.NoError(t, err)
	})
}
