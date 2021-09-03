package psqlinsert

import (
	"context"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type (
	Id    int
	Phone string
	Email string
)

type User struct {
	FirstName string
	LastName  string
	Email     Email
	Phone     Phone
	Updated   pgtype.Timestamp
}

func Psqlinsert() {
	ctx := context.Background()
	cfg, err := pgxpool.ParseConfig("postgres://myuser:secret@localhost:5432/mydb")
	if err != nil {
		log.Fatal(err)
	}
	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	myUser := User{
		"Yevg",
		"Va",
		"hi@mail.org",
		"9887654333",
		pgtype.Timestamp{
			time.Now().UTC(),
			pgtype.Present,
			pgtype.InfinityModifier(0),
		},
	}
	id, err := insert(ctx, pool, myUser)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(id)
}

func insert(ctx context.Context, pool *pgxpool.Pool, usr User) (Id, error) {
	const sqlExec = `insert into users (first_name, last_name, email, phone, updated) values 
($1, $2, $3, $4, $5::timestamp)
returning user_id;`
	var id Id
	err := pool.QueryRow(ctx, sqlExec,
		usr.FirstName,
		usr.LastName,
		usr.Email,
		usr.Phone,
		usr.Updated).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot insert: %w\n", err)
	}

	return id, nil

}
