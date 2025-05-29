package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kachaje/goydb/internal/adapter/storage"
	"github.com/kachaje/goydb/pkg/model"
	"github.com/mitchellh/mapstructure"

	"github.com/gorilla/sessions"
)

type IBase interface {
	Authenticate(http.ResponseWriter, *http.Request) (*model.Session, bool)
	PutDoc(http.ResponseWriter, *http.Request) (map[string]any, int, error)
}

type Base struct {
	Storage      *storage.Storage
	SessionStore sessions.Store
	Admins       model.AdminUsers
}

func (b *Base) Authenticate(w http.ResponseWriter, r *http.Request) (*model.Session, bool) {
	return Authenticator{Base: *b, RequiresAdmin: true}.Do(w, r)
}

func (b *Base) PutDoc(w http.ResponseWriter, r *http.Request) (map[string]any, int, error) {
	db := Database{Base: *b}.Do(w, r)
	if db == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to connect to database")
	}

	var doc map[string]any
	err := json.NewDecoder(r.Body).Decode(&doc)

	if doc["_id"] == nil {
		doc["_id"] = uuid.NewString()
	}

	docID, docIDok := doc["_id"].(string)
	if err != nil || !docIDok {
		return nil, http.StatusBadRequest, err
	}

	var attachments map[string]*model.Attachment
	err = mapstructure.Decode(doc["_attachments"], &attachments)
	if err != nil || !docIDok {
		return nil, http.StatusBadRequest, err
	}

	rev, err := db.PutDocument(r.Context(), &model.Document{
		ID:          docID,
		Data:        doc,
		Deleted:     doc["_deleted"] == "true",
		Attachments: attachments,
	})
	if errors.Is(err, storage.ErrConflict) {
		WriteError(w, http.StatusConflict, err.Error())
		return nil, http.StatusConflict, err
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil, http.StatusInternalServerError, err
	}
	doc["_rev"] = rev

	return doc, http.StatusOK, nil
}
