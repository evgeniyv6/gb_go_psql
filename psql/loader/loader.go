package loader

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"psql/psqlselect"
	"sync"
	"sync/atomic"
	"time"
)

type LoadRes struct {
	Duration         time.Duration
	Threads          int
	QueriesPerformed uint64
}

func Loader() {
	cfg, err := pgxpool.ParseConfig("postgres://myuser:secret@localhost:5432/mydb")
	if err != nil {
		log.Fatal(err)
	}
	cfg.MinConns = 8
	cfg.MaxConns = 8

	pool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	dur := time.Duration(10 * time.Second)
	threads := 1000
	fmt.Println("start load")
	res := load(context.Background(), dur, threads, pool)
	fmt.Println("duration: ", res.Duration)
	fmt.Println("threads: ", res.Threads)
	fmt.Println("queries: ", res.QueriesPerformed)
	qps := res.QueriesPerformed / uint64(res.Duration.Seconds())
	fmt.Println("QPS: ", qps)

}

func load(ctx context.Context, dur time.Duration, threads int, pool *pgxpool.Pool) LoadRes {
	var queries uint64
	loader := func(stopAt time.Time) {
		for {
			_, err := psqlselect.Search(ctx, pool, "hi", 5)
			if err != nil {
				log.Fatal(err)
			}
			atomic.AddUint64(&queries, 1)

			if time.Now().After(stopAt) {
				return
			}
		}
	}
	var wg sync.WaitGroup
	wg.Add(threads)
	startAt := time.Now()
	stopAt := startAt.Add(dur)
	for i := 0; i < threads; i++ {
		go func() {
			loader(stopAt)
			wg.Done()
		}()
	}

	wg.Wait()

	return LoadRes{
		time.Now().Sub(startAt),
		threads,
		queries,
	}
}
