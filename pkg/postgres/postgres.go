package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Pantry struct {
	QR   string
	UUID string
}

type PantryStore struct {
	pool *pgxpool.Pool
}

func NewPantryStore(pool *pgxpool.Pool) *PantryStore {
	return &PantryStore{
		pool: pool,
	}
}

func (p *PantryStore) ListPantries(ctx context.Context) ([]Pantry, error) {
	rows, err := p.pool.Query(ctx, "SELECT qr, uuid FROM pantry")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pantries := []Pantry{}

	for rows.Next() {
		pantry := Pantry{}
		if err := rows.Scan(&pantry.QR, &pantry.UUID); err != nil {
			return nil, err
		}

		pantries = append(pantries, pantry)
	}

	return pantries, nil
}
