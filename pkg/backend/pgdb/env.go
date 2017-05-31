package pgdb

import (
	"fmt"
	"os"

	"jdy/pkg/util/must"

	"github.com/jackc/pgx"
)

// Env holds connections for pgdb.
type Env struct {
	conn *pgx.ConnPool
}

// NewEnv creates a new env.
// Sample DSN: user=username password=password host=1.2.3.4 port=5432 dbname=mydb sslmode=disable
func NewEnv(dsn string) *Env {
	config, err := pgx.ParseDSN(dsn)
	must.Must(err)
	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 32,
		AfterConnect: func(c *pgx.Conn) error {
			_, err := c.Exec(fmt.Sprintf(`SET statement_timeout = %d`, 5*60*1000))
			return err
		},
	})
	must.Must(err)
	return &Env{conn: conn}
}

// NewTestEnv creates a new env for unit test.
// pgweb: postgres://postgres:postgres@192.168.59.103:5432/a2_go_test?sslmode=disable
func NewTestEnv() *Env {
	source := os.Getenv("A2_TESTING_POSTGRES_SOURCE")
	if source == "" {
		host := "127.0.0.1"
		source = fmt.Sprint("user=postgres password=postgres dbname=a2_go_test host=", host, " port=5432 sslmode=disable")
	}
	env := NewEnv(source)
	MustExec(env.conn.Exec("drop schema if exists public cascade; create schema public;"))
	return env
}

// ConnPool returns the underlying connection pool.
func (env *Env) ConnPool() *pgx.ConnPool {
	return env.conn
}
