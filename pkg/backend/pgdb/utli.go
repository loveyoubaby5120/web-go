package pgdb

import "github.com/jackc/pgx"

// MustExec panic if a exec result has error.
func MustExec(c pgx.CommandTag, err error) pgx.CommandTag {
	if err != nil {
		panic(err)
	}
	return c
}
