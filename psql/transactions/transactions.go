package transactions

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type TransactionFunc func(context.Context, pgx.Tx) error

type (
	User_id int
	FName   string
	LName   string
)

func Transactions() {
	cfg, err := pgxpool.ParseConfig("postgres://myuser:secret@localhost:5432/mydb")
	if err != nil {
		log.Fatal(err)
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	defer pool.Close()

	var id1, id2 User_id = 1000111, 1000112
	var lname LName = "Vakh"
	var name FName = "Yevg"
	err = update(context.Background(), pool, id1, id2, name, lname)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("data updated")
}

func update(ctx context.Context, pool *pgxpool.Pool, id1, id2 User_id, name FName, surname LName) error {
	err := inTx(ctx, pool, func(ctx context.Context, tx pgx.Tx) error {
		const sqlQuery = `update users set last_name = $1 where user_id = $2;`

		_, err := tx.Exec(ctx, sqlQuery, surname, id1)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, sqlQuery, name, id2)
		if err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		return err
	}

	return nil
}

func inTx(ctx context.Context, pool *pgxpool.Pool, f TransactionFunc) error {
	transact, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	err = f(ctx, transact)

	if err != nil {
		rbErr := transact.Rollback(ctx)
		if rbErr != nil {
			fmt.Println(rbErr)
		}
		return err
	}

	err = transact.Commit(ctx)
	if err != nil {
		rbErr := transact.Rollback(ctx)
		if rbErr != nil {
			fmt.Println(rbErr)
		}
		return err
	}

	return nil
}
