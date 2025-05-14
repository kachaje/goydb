package port

import "github.com/kachaje/goydb/pkg/model"

type SessionBuilder interface {
	Session() *model.Session
}
