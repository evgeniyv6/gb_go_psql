package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"net"
	"os"
	"time"
)

type (
	psqlConn struct { pool *pgxpool.Pool }
	DBAction interface {
		Close()
		SearchJID(ctx context.Context, table string, jid string) ([]int, error)
		Insert(ctx context.Context, table string, regData RegistryJobInfo) (RegistryInsertId, error)
	}

	RegistryInsertId int

	RegistryJobInfo struct {
		S_Protocol string
		S_Host string
		S_Domain string
		S_Jenkins_path string
		J_Build_type string
		J_Token string
		J_Notes string
		J_Mdata string
		J_Published string
		J_Version string
		J_Version_Tag string
		S_User_Name string
		S_User_Token string
		RegIDs []RegistryJobId
	}

	RegistryJobId struct {
		J_Id string
		J_Path string
		J_Job_name string
	}
)

func NewPool(ctx context.Context, addr, port, db, credPrefix string) (DBAction, error) {
	_, ok := os.LookupEnv(credPrefix + "_USR")
	if !ok {
		zap.S().Info("Couldnot get postgres user env variable.")
	}

	_, ok = os.LookupEnv(credPrefix + "_PSW")
	if !ok {
		zap.S().Info("Couldnot get postgres pwd env variable.")
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s","myuser", "secret", addr, port, db)
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		zap.S().Errorw("Couldnot connect to postgres", "err", err)
	}
	// Pool соединений обязательно ограничивать сверху,
	// иначе есть потенциальная опасность превысить лимит соединений с базой
	// TODO пенести параметры в конфиг файл
	cfg.MaxConns = 8
	cfg.MinConns = 4
	// HealthCheckPeriod - частота проверки работоспособности
	// соединения с Postgres
	cfg.HealthCheckPeriod = 1 * time.Minute
	// MaxConnLifetime - сколько времени будет жить соединение.
	// Так как большого смысла удалять живые соединения нет,
	// можно устанавливать большие значения
	cfg.MaxConnLifetime = 24 * time.Hour
	// MaxConnIdleTime - время жизни неиспользуемого соединения,
	// если запросов не поступало, то соединение закроется.
	cfg.MaxConnIdleTime = 30 * time.Minute
	// ConnectTimeout устанавливает ограничение по времени
	// на весь процесс установки соединения и аутентификации.
	cfg.ConnConfig.ConnectTimeout = 1 * time.Second
	// Лимиты в net.Dialer позволяют достичь предсказуемого
	// поведения в случае обрыва сети.
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		// Timeout на установку соединения гарантирует,
		// что не будет зависаний при попытке установить соединение.
		Timeout: cfg.ConnConfig.ConnectTimeout,
	}).DialContext
	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		zap.S().Errorw("Couldnot apply postgres config parameters.")
		return nil, err
	}
	return &psqlConn{dbpool}, nil
}

func (p *psqlConn) Close() {
	p.pool.Close()
}

func (p *psqlConn) SearchJID(ctx context.Context, table, jid string) ([]int, error) {
	var (
		sql = `select id from ` + table + ` where j_id = $1;`
		id int
		idList []int
	)


	rows, err := p.pool.Query(ctx,sql,jid)
	defer rows.Close()
	if err != nil {
		zap.S().Errorw("Coudnot get info from registry.", "err", err)
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			zap.S().Errorw("Couldnot get id from registry", "err", err)
		} else {
			idList = append(idList, id)
		}
	}

	return idList, nil
}

func (p *psqlConn) Insert(ctx context.Context, table string,regData RegistryJobInfo) (RegistryInsertId, error) {
	var (
		id RegistryInsertId
		sql = `INSERT INTO ` + table + ` (s_protocol, s_host,s_domain, s_jenkins_path, j_path, j_job_name,
	j_build_type, j_token, j_id, j_notes, j_mdata, j_published, j_version,
	j_version_tag, s_user_name, s_user_token) values ($1, $2, $3, $4, $5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING id;`
	)

	err := p.pool.QueryRow(ctx, sql,
		regData.S_Protocol,
		regData.S_Host,
		regData.S_Domain,
		regData.S_Jenkins_path,
		regData.J_Path,
		regData.J_Job_name,
		regData.J_Build_type,
		regData.J_Token,
		regData.J_Id,
		regData.J_Notes,
		regData.J_Mdata,
		regData.J_Published,
		regData.J_Version,
		regData.J_Version_Tag,
		regData.S_User_Name,
		regData.S_User_Token,
		).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			zap.S().Infow("No rows returned, but this is not error. Check DB record.", "err", err)
		} else {
			zap.S().Infow("Couldnot insert into registry", "err", err)
			return 0, err
		}

	}
	return id, nil
}

//func Search(ctx context.Context, dbpool *pgxpool.Pool, table string,jid string) {
//	sql := `select id from ` + table + ` where j_id = $1;`
//
//	rows, err := dbpool.Query(ctx, sql, jid)
//	if err != nil {
//		//return nil, fmt.Errorf("failed to query data: %w", err)
//		fmt.Println("++",err)
//	}
//
//	defer rows.Close()
//
//	var id int
//	var idList []int
//	for rows.Next() {
//		err = rows.Scan(&id)
//		if err != nil {
//			fmt.Println("---->>", err)
//		}
//		fmt.Println("id = ", id)
//		idList = append(idList, id)
//	}
//
//	if len(idList) > 0 {
//		fmt.Println(idList)
//	}
//
//
//
//}
