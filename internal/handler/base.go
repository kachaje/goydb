package handler

import (
	"net/http"

	"github.com/kachaje/goydb/internal/adapter/storage"
	"github.com/kachaje/goydb/pkg/model"

	"github.com/gorilla/sessions"
)

type IBase interface {
	Authenticate(http.ResponseWriter, *http.Request) (*model.Session, bool)
}

type Base struct {
	Storage      *storage.Storage
	SessionStore sessions.Store
	Admins       model.AdminUsers
}

func (b *Base) Authenticate(w http.ResponseWriter, r *http.Request) (*model.Session, bool) {
	return Authenticator{Base: *b, RequiresAdmin: true}.Do(w, r)
}
