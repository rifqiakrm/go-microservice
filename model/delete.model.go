package model

import (
	"context"
	"database/sql"
	"fmt"
)

func DeleteSample(ctx context.Context, db *sql.DB, id int64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	query := `
		DELETE FROM 
			samples 
		WHERE 
			id = $1`
	_, errDelete := tx.ExecContext(ctx, query, id)

	if errDelete != nil {
		return fmt.Errorf("there was an error on model/product (DeleteSample) : %v", errDelete)
	}

	defer tx.Commit()
	return nil
}

