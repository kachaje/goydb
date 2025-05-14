package goydb

import (
	"net/http"

	"github.com/kachaje/goydb/internal/adapter/storage"
)

type Goydb struct {
	*storage.Storage
	Handler http.Handler
}
