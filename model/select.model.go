package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/rifqiakrm/go-microservice-lib/cache"
	"github.com/rifqiakrm/go-microservice/resources"
	"github.com/rifqiakrm/go-microservice/utils/cache_tag"
)

func GetSample(ctx context.Context, db *sql.DB, id int64) (*resources.Sample, error) {
	var sample resources.Sample

	var query = `
		SELECT 
			id,
			message
		FROM 
			sample 
		WHERE 
			id = $1`

	bytes, _ := cache.Get(fmt.Sprintf(cache_tag.GET_SAMPLE, id))

	if bytes != nil {
		if err := json.Unmarshal(bytes, &sample); err != nil {
			return nil, err
		}

		return &sample, nil
	}

	err := db.QueryRowContext(ctx, query, id).Scan(
		&sample.ID,
		&sample.Message,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, fmt.Errorf("there was an error on model/select.model.go : %v", err)
		}
	}

	if err := cache.Set(fmt.Sprintf(cache_tag.GET_SAMPLE, id), &sample, 3600); err != nil {
		return nil, fmt.Errorf("there was an error while caching query on model/select.model.go (GET_SAMPLE) : %v", err.Error())
	}

	return &sample, nil
}

func GetSamples(ctx context.Context, db *sql.DB) ([]*resources.Sample, error) {
	var samples []*resources.Sample

	var query = `
		SELECT 
			id,
			message
		FROM 
			sample`

	bytes, _ := cache.Get(fmt.Sprintf(cache_tag.GET_SAMPLES))

	if bytes != nil {
		if err := json.Unmarshal(bytes, &samples); err != nil {
			return nil, err
		}

		return samples, nil
	}

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, fmt.Errorf("there was an error on model/select.model.go : %v", err)
		}
	}

	defer rows.Close()

	for rows.Next() {
		var sample resources.Sample

		errScan := rows.Scan(
			&sample.ID,
			&sample.Message,
		)

		if errScan != nil {
			return nil, fmt.Errorf("there was an error on model/select.model.go : %v", err)
		}

		samples = append(samples, &sample)
	}

	if err := cache.Set(fmt.Sprintf(cache_tag.GET_SAMPLES), &samples, 3600); err != nil {
		return nil, fmt.Errorf("there was an error while caching query on model/select.model.go (GET_SAMPLES) : %v", err.Error())
	}

	return samples, nil
}
