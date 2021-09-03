package pgx4

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
)

func Pgx4() {
	conn, err := pgx.Connect(context.Background(), "postgres://myuser:secret@localhost:5432/mydb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot connect to Postgres: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var res string
	err = conn.QueryRow(context.Background(), "select 'hello'").Scan(&res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot query rows: %v\n", err)
	}
	fmt.Println(res)
}
