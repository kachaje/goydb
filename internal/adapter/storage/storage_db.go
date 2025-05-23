package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/kachaje/goydb/internal/adapter/bbolt_engine"
	"github.com/kachaje/goydb/internal/adapter/index"
	"github.com/kachaje/goydb/pkg/model"
	"github.com/kachaje/goydb/pkg/port"
)

type Database struct {
	name        string
	databaseDir string
	db          port.DatabaseEngine

	listener sync.Map

	indices        map[string]port.DocumentIndex
	viewEngines    map[string]port.ViewServerBuilder
	reducerEngines map[string]port.ReducerServerBuilder
}

func (d *Database) ChangesIndex() port.DocumentIndex {
	return d.indices[index.ChangesIndexName]
}

func (d *Database) Indices() map[string]port.DocumentIndex {
	return d.indices
}

func (d *Database) Name() string {
	return d.name
}

func (d *Database) String() string {
	stats, err := d.db.Stats()
	if err == nil {
		return fmt.Sprintf("<Database name=%q stats=%+v>", d.name, stats)
	}

	return fmt.Sprintf("<Database name=%q stats=%v>", d.name, err)
}

func (d *Database) Stats(ctx context.Context) (stats model.DatabaseStats, err error) {
	return d.db.Stats()
}

func (d *Database) Sequence(ctx context.Context) (string, error) {
	var seq uint64
	err := d.Transaction(ctx, func(tx *Transaction) error {
		seq = tx.Sequence(model.DocsBucket)
		return nil
	})
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(seq, 10), nil
}

func (d *Database) ViewEngine(name string) port.ViewServerBuilder {
	return d.viewEngines[name]
}

func (d *Database) ReducerEngine(name string) port.ReducerServerBuilder {
	return d.reducerEngines[name]
}

func (s *Storage) CreateDatabase(ctx context.Context, name string) (*Database, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	databaseDir := path.Join(s.path, name+".d")

	log.Println("Open..")
	db, err := bbolt_engine.Open(path.Join(s.path, name))
	if err != nil {
		return nil, err
	}
	log.Println("Done..")

	database := &Database{
		name:        name,
		databaseDir: databaseDir,
		db:          db,
		indices: map[string]port.DocumentIndex{
			index.ChangesIndexName: index.NewChangesIndex(),
			index.DeletedIndexName: index.NewDeletedIndex(),
		},
		viewEngines:    s.viewEngines,
		reducerEngines: s.reducerEngines,
	}
	s.dbs[name] = database

	// create all required database Indices
	log.Println("BuildIndices")
	err = database.Transaction(ctx, func(tx *Transaction) error {
		tx.EnsureBucket(model.DocsBucket)

		err := database.BuildIndices(ctx, tx, false)
		if err != nil {
			return err
		}

		for _, index := range database.Indices() {
			log.Printf("ENSURE INDEX %s", index)
			err := index.Ensure(ctx, tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return database, nil
}

func (s *Storage) DeleteDatabase(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	db, ok := s.dbs[name]
	if !ok {
		return fmt.Errorf("%w: %q", ErrUnknownDatabase, name)
	}

	err := db.db.Close()
	if err != nil {
		return err
	}

	err = os.Remove(path.Join(s.path, name))
	if err != nil {
		return err
	}

	delete(s.dbs, name)

	return nil
}

func (s *Storage) Databases(ctx context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, len(s.dbs))
	var i int
	for name := range s.dbs {
		names[i] = name
		i++
	}

	return names, nil
}

func (s *Storage) Database(ctx context.Context, name string) (*Database, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	db, ok := s.dbs[name]
	if !ok {
		return nil, fmt.Errorf("database %q not found", name)
	}

	return db, nil
}
