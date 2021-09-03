package psqlselect

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type (
	Phone string
	Email string
)

type EmailSearchHint struct {
	Phone Phone
	Email Email
}

func Psqlselect() {
	cfg, err := pgxpool.ParseConfig("postgres://myuser:secret@localhost:5432/mydb")
	if err != nil {
		log.Fatal(err)
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	res, err := Search(context.Background(), pool, "hi", 2)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range res {
		fmt.Println(r.Email, r.Phone)
	}
}

func Search(ctx context.Context, pool *pgxpool.Pool, prefix string, limit int) ([]EmailSearchHint, error) {
	const sqlQuery = `select email, phone from users where email like $1 order by email limit $2;`

	pattern := prefix + "%"

	rows, err := pool.Query(ctx, sqlQuery, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("query err: %w\n", err)
	}
	defer rows.Close()

	var res []EmailSearchHint
	for rows.Next() {
		var r EmailSearchHint
		err = rows.Scan(&r.Email, &r.Phone)
		if err != nil {
			return nil, fmt.Errorf("cannot scan rows: %w\n", err)
		}

		res = append(res, r)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to read respponse: %w", err)
	}
	return res, nil

}
