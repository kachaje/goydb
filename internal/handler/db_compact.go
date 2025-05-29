package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kachaje/goydb/pkg/port"
)

type DBDocsCompact struct {
	Base
}

func (s *DBDocsCompact) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	db := Database{Base: s.Base}.Do(w, r)
	if db == nil {
		return
	}

	if _, ok := (Authenticator{Base: s.Base}.DB(w, r, db)); !ok {
		return
	}

	var q port.AllDocsQuery
	q.IncludeDocs = true

	docs, _, err := db.AllDocs(r.Context(), q)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	dbname := db.Name()

	tmpDb := fmt.Sprintf("%s_tmp", dbname)

	newDb, err := s.Storage.CreateDatabase(r.Context(), tmpDb)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer func() {
		err := s.Storage.DeleteDatabase(r.Context(), tmpDb)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}()

	for _, doc := range docs {
		newDb.PutDocument(r.Context(), doc)
	}

	err = s.Storage.DeleteDatabase(r.Context(), dbname)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedDb, err := s.Storage.CreateDatabase(r.Context(), dbname)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, doc := range docs {
		updatedDb.PutDocument(r.Context(), doc)
	}

	response := CompactResponse{
		OK: true,
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

type CompactResponse struct {
	OK bool `json:"ok"`
}
