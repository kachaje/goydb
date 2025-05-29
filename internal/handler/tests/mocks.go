package tests

import (
	"net/http"

	"github.com/kachaje/goydb/pkg/model"
)

type MockBase struct{}

func (m *MockBase) Authenticate(http.ResponseWriter, *http.Request) (*model.Session, bool) {
	return nil, true
}
