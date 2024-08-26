package iter2

import (
	"database/sql"
	"iter"
)

// All returns an iterator over rows in the *sql.Rows.
func AllRows[T any](rows *sql.Rows) iter.Seq[*T] {
	return func(yield func(*T) bool) {
		for rows.Next() {
			var v T
			if err := rows.Scan(&v); err != nil {
				return
			}
			if !yield(&v) {
				return
			}
		}
	}
}
