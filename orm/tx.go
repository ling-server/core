package orm

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/beego/beego/orm"
	"github.com/google/uuid"
)

type CommittedKey struct{}

// HasCommittedKey checks whether exist committed key in context.
func HasCommittedKey(ctx context.Context) bool {
	if value := ctx.Value(CommittedKey{}); value != nil {
		return true
	}

	return false
}

// ormerTx transaction which support savepoint
type ormerTx struct {
	orm.Ormer
	savepoint string
}

func (o *ormerTx) savepointMode() bool {
	return o.savepoint != ""
}

func (o *ormerTx) createSavepoint() error {
	val := uuid.New()
	o.savepoint = fmt.Sprintf("p%s", hex.EncodeToString(val[:]))

	_, err := o.Raw(fmt.Sprintf("SAVEPOINT %s", o.savepoint)).Exec()
	return err
}

func (o *ormerTx) releaseSavepoint() error {
	_, err := o.Raw(fmt.Sprintf("RELEASE SAVEPOINT %s", o.savepoint)).Exec()
	return err
}

func (o *ormerTx) rollbackToSavepoint() error {
	_, err := o.Raw(fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", o.savepoint)).Exec()
	return err
}

func (o *ormerTx) Begin() error {
	err := o.Ormer.Begin()
	if err == orm.ErrTxHasBegan {
		// transaction has began for the ormer, so begin nested transaction by savepoint
		return o.createSavepoint()
	}

	return err
}

func (o *ormerTx) Commit() error {
	if o.savepointMode() {
		return o.releaseSavepoint()
	}

	return o.Ormer.Commit()
}

func (o *ormerTx) Rollback() error {
	if o.savepointMode() {
		return o.rollbackToSavepoint()
	}

	return o.Ormer.Rollback()
}
