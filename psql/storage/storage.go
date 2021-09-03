package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PG struct {
	pool *pgxpool.Pool
}

func NewPG(pool *pgxpool.Pool) *PG {
	return &PG{pool}
}

type (
	FirstName string
	LastName  string
)

type DBInfoRes struct {
	FirstName FirstName
	LastName  LastName
}

func (s *PG) Search(ctx context.Context, prefix string, limit int) ([]DBInfoRes, error) {
	const sqlQuery = `
select first_name, last_name from users where first_name like $1 limit $2;
`
	pattern := prefix + "%"
	rows, err := s.pool.Query(ctx, sqlQuery, pattern, limit)

	if err != nil {
		return nil, fmt.Errorf("cannot query data: %w", err)
	}

	var res []DBInfoRes

	for rows.Next() {
		var r DBInfoRes

		err = rows.Scan(&r.FirstName, &r.LastName)

		if err != nil {
			return nil, fmt.Errorf("cannot query data: %w", err)
		}

		res = append(res, r)

	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("cannot read response: %w", rows.Err())
	}
	return res, nil
}
