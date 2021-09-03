// +build integration

package storage_test

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"log"
	"psql/storage"
	"testing"
)

func TestPG_Search(t *testing.T) {
	pool := connect(context.Background())
	defer pool.Close()

	tests := []struct {
		name    string
		store   *storage.PG
		ctx     context.Context
		prefix  string
		limit   int
		prepare func(*pgxpool.Pool)
		check   func(*testing.T, []storage.DBInfoRes, error)
	}{
		{
			"ok",
			storage.NewPG(pool),
			context.Background(),
			"Yevg",
			2,
			func(pool *pgxpool.Pool) {
				pool.Exec(context.Background(), `insert into users (first_name, last_name, email, phone, updated) values 
('Yevgen', 'Va', 'my@mail.pod', '1234567', to_timestamp('25-08-2021 15:36:38', 'dd-mm-yyyy hh24:mi:ss'))`)
			},
			func(t *testing.T, res []storage.DBInfoRes, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(pool)
			res, err := tt.store.Search(tt.ctx, tt.prefix, tt.limit)
			tt.check(t, res, err)
		})
	}
}

func connect(ctx context.Context) *pgxpool.Pool {
	pool, err := pgxpool.Connect(ctx, "postgres://myuser:secret@localhost:5432/mydb")
	if err != nil {
		log.Fatal(err)
	}
	return pool

}
