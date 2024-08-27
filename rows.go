package iter2

import (
	"database/sql"
	"iter"
)

// Row is a row of [sql.Rows].
type Row struct {
	rows *sql.Rows
}

// Scan copies the columns in the current row into the values pointed at by dest.
// See Scan method of [sql.Rows].
func (row Row) Scan(dest ...any) error {
	return row.rows.Scan(dest...)
}

// All returns an iterator over rows in the *sql.Rows.
func AllRows(rows *sql.Rows) iter.Seq[Row] {
	return func(yield func(Row) bool) {
		for rows.Next() {
			if !yield(Row{rows}) {
				return
			}
		}
	}
}

// MustAllRows is similar to the [AllRows], but with two differences:
// 1. It has an additional err parameter. If err is not nil, MustAllRows will call panic(err).
// Otherwise, it calls AllRows(rows) and returns its result.
// 2. Once the returned Seq  is iterated, it will call rows.Close() when done.
//
// MustAllRows is convenient when calling with Query of sql. For example:
//
//	seq := iter2.MustAllRows(db.Query(q))
func MustAllRows(rows *sql.Rows, err error) iter.Seq[Row] {
	if err != nil {
		panic(err)
	}
	return func(yield func(Row) bool) {
		defer rows.Close()
		for rows.Next() {
			if !yield(Row{rows}) {
				return
			}
		}
	}
}
