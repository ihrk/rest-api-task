package database

import (
	"context"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exec"
)

type TxRunner interface {
	WithTx(fn func(*goqu.TxDatabase) error) error
}

type SessionRunner interface {
	Insert(table interface{}) *goqu.InsertDataset
	Select(cols ...interface{}) *goqu.SelectDataset
	Update(table interface{}) *goqu.UpdateDataset
	Delete(table interface{}) *goqu.DeleteDataset
}

var ErrNotFound = errors.New("not found")

func SelectStructQuery(
	ctx context.Context,
	q exec.QueryExecutor,
	dst interface{},
) error {
	found, err := q.ScanStructContext(ctx, dst)
	if err != nil {
		return err
	}

	if !found {
		return ErrNotFound
	}

	return nil
}
