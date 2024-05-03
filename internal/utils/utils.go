package utils

import (
	"database/sql"
	"strings"

	log "github.com/sirupsen/logrus"
)

func NormalizeString(input string) (normalized string) {
	normalized = strings.Trim(input, " \t\r\n")
	normalized = strings.ReplaceAll(normalized, `“`, `"`)
	normalized = strings.ReplaceAll(normalized, `”`, `"`)
	normalized = strings.ReplaceAll(normalized, `‘`, `'`)
	normalized = strings.ReplaceAll(normalized, `’`, `'`)
	normalized = strings.ReplaceAll(normalized, "\u00A0", ` `) // nonbreaking space
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

func MustExec(stmt *sql.Stmt, args ...interface{}) int64 {
	if result, err := stmt.Exec(args...); err != nil {
		log.Panic(append([]interface{}{err}, args...)...)
		return -1
	} else {
		if n, err := result.RowsAffected(); err != nil {
			return 0
		} else {
			return n
		}
	}
}
