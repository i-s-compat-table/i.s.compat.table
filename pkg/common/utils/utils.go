package utils

import (
	"database/sql"
	"log"
	"strings"
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

func MustExec(stmt *sql.Stmt, args ...interface{}) {
	if _, err := stmt.Exec(args...); err != nil {
		log.Panic(append([]interface{}{err}, args...)...)
	}
}
