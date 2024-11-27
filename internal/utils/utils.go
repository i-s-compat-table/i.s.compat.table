package utils

import (
	"database/sql"
	"strings"
)

func NormalizeString(input string) (normalized string) {
	normalized = strings.Trim(input, " \t\r\n")
	normalized = strings.ReplaceAll(normalized, `“`, `"`)
	normalized = strings.ReplaceAll(normalized, `”`, `"`)
	normalized = strings.ReplaceAll(normalized, `‘`, `'`)
	normalized = strings.ReplaceAll(normalized, `’`, `'`)
	normalized = strings.ReplaceAll(normalized, "\u00A0", ` `) // non-breaking space
	normalized = strings.ReplaceAll(normalized, "\u200B", "")  // 0-width space
	return normalized
}

func MustPrepare(txn *sql.Tx, query string) *sql.Stmt {
	stmt, err := txn.Prepare(query)
	if err != nil {
		panic(err)
	}
	return stmt
}

func Last[T any](arr []T) *T {
	if len(arr) == 0 {
		return nil
	}
	return &arr[len(arr)-1]
}

func MustExec(stmt *sql.Stmt, args ...interface{}) int64 {
	if result, err := stmt.Exec(args...); err != nil {
		panic(append([]interface{}{err.Error()}, args...))
	} else {
		if n, err := result.RowsAffected(); err != nil {
			return 0
		} else {
			return n
		}
	}
}

func Series(min, max int) <-chan int {
	ch := make(chan int, 1)
	go func() {
		for i := min; i < max; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

var ephemeralPorts <-chan int

func init() {
	ephemeralPorts = Series(49152, 65535)
}

func GetEphemeralPort() int {
	return <-ephemeralPorts
}
