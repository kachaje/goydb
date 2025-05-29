package tests

import (
	"net/http"

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
