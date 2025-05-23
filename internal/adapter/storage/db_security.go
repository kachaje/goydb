package storage

import (
	"context"

	"github.com/kachaje/goydb/pkg/model"
	"github.com/kachaje/goydb/pkg/port"
)

var (
	internalDocsBucket = []byte("_internal")
	securityDoc        = []byte("_security")
)

func (d *Database) PutSecurity(ctx context.Context, sec *model.Security) error {
	return d.Transaction(ctx, func(tx *Transaction) error {
		tx.SetBucketName(internalDocsBucket)
		return tx.PutRaw(ctx, securityDoc, sec)
	})
}

func (d *Database) GetSecurity(ctx context.Context) (*model.Security, error) {
	var sec model.Security
	err := d.Transaction(ctx, func(tx *Transaction) error {
		tx.SetBucketName(internalDocsBucket)
		err := tx.GetRaw(ctx, securityDoc, &sec)
		return err
	})
	if err == port.ErrNotFound {
		return model.DefaultSecurity(), nil
	}
	if err != nil {
		return nil, err
	}
	return &sec, err
}
