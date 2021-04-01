package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rifqiakrm/go-microservice/resources"
)

func UpdateSample(ctx context.Context, db *sql.DB, req *resources.Sample) error {
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

	query := fmt.Sprintf(`
		UPDATE 
			sample 
		SET
			name = $1
		WHERE
			id = $2
		RETURNING
			id
`)

	errUpdate := tx.QueryRowContext(ctx, query,
		req.Message,
		req.ID,
	).Scan(
		&req.ID,
	)

	if errUpdate != nil {
		return fmt.Errorf("there was an error on model/sample (UpdateSample) : %v", errUpdate)
	}

	defer tx.Commit()
	return nil
}
