package schema

import (
	"database/sql"
	_ "embed"
	"log"
	"os"
	"strings"

	_ "github.com/cespare/xxhash/v2"
)

//go:embed db.sql
var Schema string

type Nullability int8

const (
	Unknown     Nullability = -1
	Nullable    Nullability = 0
	NotNullable Nullability = 1
)

type Column interface {
	// xxhash3_64 of some at least Db, Table, Column, maybe type, nullable, and/or notes
	Id() int64
	Db() string
	Table() string
	Column() string
	Type() *string // may be nil
	Nullable() Nullability
	Notes() string

	Url() string
	Version() string
}

func Connect(path string) (db *sql.DB, err error) {
	if f, e := os.Stat(path); e != nil {
		// initialize if the path doesn't exist
		db, err = sql.Open("sqlite3", path)
		if err != nil {
			panic(err)
		}
		for _, cmd := range strings.Split(Schema, ";") {
			_, err = db.Exec(cmd)
			if err != nil {
				panic(err)
			}
		}
		return db, err
	} else if f.IsDir() {
		log.Panicf("%s is a directory", path)
	} else {
		db, err = sql.Open("sqlite3", path)
	}
	return db, err
}

func MustConnect(path string) *sql.DB {
	db, err := Connect(path)
	if err != nil {
		panic(err)
	}
	return db
}

// func BulkInsert(txn *sql.Tx, columns *[]Column, onConflict string) (error) {
// 	s := strings.Builder{}
// 	s.WriteString("INSERT INTO columns ")
// 	s.WriteString("VALUES")
// 	for _,
// 	s.WriteString("")
// 	s.Wr
// // id
// // db_name
// // table_name
// // column_name
// // column_index
// // column_type
// // column_nullable
// // notes
// // created_at
// // last_modified
// }

// func BulkInsert(db *sql.DB, progress *pb.ProgressBar) {

// 	writer := func(columns <-chan *[]Column) {
// 		txn, err := db.Begin()
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		defer txn.Commit()

// 		s := strings.Builder{}
// 			s.WriteString("INSERT INTO predictions")
// 			s.WriteString("(statement_id, oracle_id, language_id, message, error, valid)")
// 			s.WriteString(" VALUES ")

// 	}
// 	// TODO: some <- chan magic to do worker pooling with a single batched writer
// }
