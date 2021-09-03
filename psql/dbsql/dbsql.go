package dbsql

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
)

func Dbsql() {
	db, err := sql.Open("pgx", "postgres://myuser:secret@localhost:5432/mydb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var res string
	err = db.QueryRow("select 'hello from db sql'").Scan(&res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
