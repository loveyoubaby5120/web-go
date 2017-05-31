package pgdb

import (
	"fmt"

	"bytes"

	"time"

	"jdy/pkg/base/sentry"
	"jdy/pkg/base/trace"
	"jdy/pkg/util/errs"

	"github.com/jackc/pgx"
)

// Tx reprsents a transcation.
type Tx struct {
	runner   txRunner
	scanning bool
}

type txRunner interface {
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
}

// NoTransactionUseWithCaution returns a transcation object that execute statement directly.
func (env *Env) NoTransactionUseWithCaution() *Tx {
	return &Tx{runner: env.conn}
}

func (env *Env) tryTransaction(f func(tx *Tx) error) (serializeError *pgx.PgError, err error) {
	pgtx, err := env.conn.BeginIso(pgx.Serializable)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer func() {
		err := pgtx.Rollback()
		if err != nil && err != pgx.ErrTxClosed {
			panic(err)
		}
	}()

	sErr, err := func() (serializeError *pgx.PgError, e error) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					if pgErr, ok := errs.Unwrap(err).(pgx.PgError); ok {
						if pgErr.Code == "40001" {
							serializeError = &pgErr
							return
						}
					}
				}
				panic(r)
			}
		}()
		err := f(&Tx{
			runner: pgtx,
		})
		if err != nil {
			if pgErr, ok := errs.Unwrap(err).(pgx.PgError); ok {
				if pgErr.Code == "40001" {
					return &pgErr, nil
				}
			}
			return nil, err
		}
		return nil, nil
	}()
	if err != nil {
		return nil, err
	}
	if sErr != nil {
		return sErr, nil
	}

	err = pgtx.Commit()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return nil, nil
}

// Transaction starts a new transaction.
func (env *Env) Transaction(f func(tx *Tx) error) error {
	var sErr *pgx.PgError
	var err error
	for i := 0; i < 100; i++ {
		sErr, err = env.tryTransaction(f)
		if err != nil {
			return err
		}
		if sErr == nil {
			return nil
		}
	}
	return fmt.Errorf("cannot serialize transcation after 100 attempts, %v", sErr)
}

// ReadCommittedTransaction is a weaker transaction which does not gurantee
// read data not be changed during the transaction by a different transaction.
// It is typically used for read-only data import.
func (env *Env) ReadCommittedTransaction(f func(tx *Tx) error) (err error) {
	pgtx, err := env.conn.BeginIso(pgx.ReadCommitted)
	if err != nil {
		return errs.Wrap(err)
	}
	defer func() {
		rbErr := pgtx.Rollback()
		if rbErr == nil || rbErr == pgx.ErrTxClosed {
			return
		}
		panic(rbErr)
	}()
	err = f(&Tx{
		runner: pgtx,
	})
	if err != nil {
		return err
	}
	err = pgtx.Commit()
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

// Snapshot returns a read-only transcation object that execute statement in snapshot transaction.
func (env *Env) Snapshot() (*Tx, func()) {
	pgtx, err := env.conn.BeginIso(pgx.ReadCommitted + " READ ONLY")
	if err != nil {
		panic(err)
	}
	return &Tx{runner: pgtx}, func() {
		err := pgtx.Commit()
		if err == nil || err == pgx.ErrTxClosed {
			return
		}
		sentry.Error(err)
	}
}

// Exec delegates to pgx.Exec.
func (tx *Tx) Exec(tr trace.T, sql string, args ...interface{}) (pgx.CommandTag, error) {
	tr.Printf("SQL: %s", sql)
	defer tr.Printf("done.")
	return tx.runner.Exec(sql, args...)
}

// Query delegates to pgx.Query.
func (tx *Tx) Query(tr trace.T, sql string, args ...interface{}) *Rows {
	tr.Printf("SQL: %s", sql)
	defer tr.Printf("done.")

	rows, err := tx.runner.Query(sql, args...)
	errs.Ignore(err) // Handled by Row.Scan.
	return &Rows{rows: rows, tx: tx}
}

// QueryRow delegates to pgx.QueryRow.
func (tx *Tx) QueryRow(tr trace.T, sql string, args ...interface{}) *Row {
	tr.Printf("SQL: %s", sql)
	defer tr.Printf("done.")

	rows, err := tx.runner.Query(sql, args...)
	errs.Ignore(err) // Handled by Row.Scan.
	return &Row{rows: rows, standalone: true}
}

// SetStatementTimeout sets the statement timeout for current transaction.
func (tx *Tx) SetStatementTimeout(tr trace.T, dur time.Duration) error {
	_, err := tx.Exec(tr, fmt.Sprintf(`SET LOCAL statement_timeout = %d`, dur/time.Millisecond))
	return err
}

// Row is almost the same as pgx.Row but more robust with panic.
type Row struct {
	rows       *pgx.Rows
	standalone bool
}

// Scan works the same as (*Rows Scan) with the following exceptions. If no
// rows were found it returns ErrNoRows. If multiple rows are returned it
// ignores all but the first.
func (r *Row) Scan(dest ...interface{}) (err error) {
	rows := r.rows
	if r.standalone {
		defer rows.Close()
	}
	if rows.Err() != nil {
		return rows.Err()
	}
	if r.standalone && !rows.Next() {
		if rows.Err() == nil {
			return pgx.ErrNoRows
		}
		return rows.Err()
	}
	return rows.Scan(dest...)
}

// Values returns an array of the row values.
func (r *Row) Values() ([]interface{}, error) {
	rows := r.rows
	if r.standalone {
		defer rows.Close()
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	if r.standalone && !rows.Next() {
		if rows.Err() == nil {
			return nil, pgx.ErrNoRows
		}
		return nil, rows.Err()
	}
	return rows.Values()
}

// Rows provides easier access to pgx.Rows which will be simpler in handling rows.Close.
type Rows struct {
	rows *pgx.Rows
	tx   *Tx
}

// fieldDescriptions returns the result field descriptions.
func (r *Rows) fieldDescriptions() []pgx.FieldDescription {
	return r.rows.FieldDescriptions()
}

// Close closes the rows.
func (r *Rows) Close() {
	r.rows.Close()
}

// Each iterate over the rows and close it.
func (r *Rows) Each(action func(r *Row) error) error {
	rows := r.rows
	defer func() {
		r.tx.scanning = false
	}()
	defer rows.Close()
	if rows.Err() != nil {
		return rows.Err()
	}
	for rows.Next() {
		err := action(&Row{rows: rows, standalone: false})
		if err != nil {
			return err
		}
	}
	return rows.Err()
}

// DebugString consumes the query and covert the result as debug string.
func (r *Rows) DebugString() string {
	var buf bytes.Buffer
	if r.rows.Err() != nil {
		fmt.Fprintf(&buf, "%v", r.rows.Err())
		return buf.String()
	}

	for _, f := range r.fieldDescriptions() {
		fmt.Fprintf(&buf, "%s,", f.Name)
	}
	buf.WriteString("\n")
	err := r.Each(func(r *Row) error {
		v, err := r.Values()
		if err != nil {
			fmt.Fprintf(&buf, "error: %v", err)
			return err
		}
		fmt.Fprintf(&buf, "%#v\n", v)
		return nil
	})
	if err != nil {
		fmt.Fprintf(&buf, "error: %v", err)
	}
	return buf.String()
}

// One returns a row that will close after scan.
func (r *Rows) One() *Row {
	return &Row{rows: r.rows, standalone: true}
}
