package postgresql

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

// Transact runs fn inside a sqlx transaction on db.
// It calls Beginx, executes fn, and commits on success.
// On fn failure it rolls back; if both fn and Rollback fail, both errors are combined via errors.Join.
func Transact(db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	if rbErr := tx.Rollback(); rbErr != nil {
		return errors.Join(err, rbErr)
	}

	return err
}
