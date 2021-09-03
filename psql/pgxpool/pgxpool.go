package pgxpool

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"time"
)

func Pgxpool() {
	ctx := context.Background()
	cfg, err := pgxpool.ParseConfig("postgres://myuser:secret@localhost:5432/mydb")
	if err != nil {
		log.Fatal(err)
	}
	cfg.MaxConns = 8
	cfg.MinConns = 4

	cfg.HealthCheckPeriod = 1 * time.Minute
	cfg.MaxConnLifetime = 24 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	cfg.ConnConfig.ConnectTimeout = 1 * time.Second

	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		Timeout:   cfg.ConnConfig.ConnectTimeout,
	}).DialContext

	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	var res string
	err = dbpool.QueryRow(ctx, "select 'hello from pool'").Scan(&res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
