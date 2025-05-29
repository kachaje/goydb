package tests

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kachaje/goydb/pkg/model"
)

type MockBase struct{}

func (m *MockBase) Authenticate(w http.ResponseWriter, r *http.Request) (*model.Session, bool) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		return nil, false
	}

	if user != "admin" || pass != "secret" {
		return nil, false
	}

	return nil, true
}

func (b *MockBase) PutDoc(w http.ResponseWriter, r *http.Request) (map[string]any, int, error) {
	var doc map[string]any
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	doc["_rev"] = uuid.NewString()

	return doc, http.StatusOK, nil
}
