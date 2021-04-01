package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rifqiakrm/go-microservice-lib/cache"
	"github.com/rifqiakrm/go-microservice-lib/formatter"
	"github.com/rifqiakrm/go-microservice/resources"
	"github.com/rifqiakrm/go-microservice/utils/cache_tag"
	"time"
)

func InsertSample(ctx context.Context, db *sql.DB, req *resources.Sample) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	query := fmt.Sprintf(`
		INSERT INTO sample (
			message
			created_at,
			updated_at
		) VALUES(
			$1, $2, $3
		)
		RETURNING
			id`)

	errInsert := tx.QueryRowContext(ctx, query,
		req.Message,
		formatter.DateWithTimestamp(time.Now().Local()),
		formatter.DateWithTimestamp(time.Now().Local()),
	).Scan(
		&req.ID,
	)

	if errInsert != nil {
		return 0, fmt.Errorf("there was an error on model/product (InsertProduct) : %v", errInsert)
	}

	if err := cache.Remove(fmt.Sprintf(cache_tag.GET_SAMPLE, req.ID)); err != nil {
		return 0, fmt.Errorf("there was an error while removing cache_tag (GET_SAMPLE) : %v", err)
	}

	defer tx.Commit()
	return req.ID, nil
}
