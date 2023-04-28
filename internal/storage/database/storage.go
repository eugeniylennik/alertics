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
	InsertMetrics(ctx context.Context, m metrics.Metrics) error
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

func NewStorage(client database.Client) Repository {
	return &Storage{
		client,
	}
}
