package handler

import (
	"github.com/kachaje/goydb/internal/adapter/storage"
	"github.com/kachaje/goydb/pkg/model"

	"github.com/gorilla/sessions"
)

type Base struct {
	Storage      *storage.Storage
	SessionStore sessions.Store
	Admins       model.AdminUsers
}
