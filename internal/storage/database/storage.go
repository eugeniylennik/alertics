package database

import (
	"context"
	"github.com/eugeniylennik/alertics/internal/database"
	"github.com/eugeniylennik/alertics/internal/metrics"
)

type Storage struct {
	database.Client
}

type Repository interface {
	SelectMetricById(ctx context.Context, m metrics.Metrics) (metrics.Metrics, error)
	InsertMetrics(ctx context.Context, m metrics.Metrics) error
	InsertMetricsStatement(ctx context.Context, m []metrics.Metrics) error
}

func (s *Storage) SelectMetricById(ctx context.Context, m metrics.Metrics) (metrics.Metrics, error) {
	q := `
        SELECT id, type, delta, value, hash 
        FROM "public".metrics
        WHERE id=$1 AND type=$2
        `
	var r metrics.Metrics
	err := s.QueryRow(ctx, q, m.ID, m.MType).Scan(&r.ID, &r.MType, &r.Delta, &r.Value, &r.Hash)
	if err != nil {
		return metrics.Metrics{}, err
	}
	return r, nil
}

func (s *Storage) InsertMetrics(ctx context.Context, m metrics.Metrics) error {
	q := `
        INSERT INTO public."metrics" (id, type, delta, value, hash) 
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE
        SET type = excluded.type,
            delta = excluded.delta,
            value = excluded.value,
            hash = excluded.hash`
	if _, err := s.Exec(ctx, q, m.ID, m.MType, m.Delta, m.Value, m.Hash); err != nil {
		return err
	}
	return nil
}

func (s *Storage) InsertMetricsStatement(ctx context.Context, m []metrics.Metrics) error {
	tx, err := s.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := `
		INSERT INTO public."metrics" (id, type, delta, value, hash) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET type = excluded.type,
			delta = excluded.delta,
			value = excluded.value,
			hash = excluded.hash`

	_, err = tx.Prepare(ctx, "insert-metrics", q)
	if err != nil {
		return err
	}

	for _, metric := range m {
		_, err = tx.Exec(context.Background(), "insert-metrics",
			metric.ID, metric.MType, metric.Delta, metric.Value, metric.Hash)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func NewStorage(client database.Client) Repository {
	return &Storage{
		client,
	}
}
