package main

import (
	"fmt"
	"os"

	"encoding/csv"

	"github.com/i-s-compat-table/i.s.compat.table/internal/observer"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const driver = "postgres"
const dsn = "host=db port=5432 user=postgres password=password sslmode=disable"

var query = `
	SELECT
		col.table_name
		, col.column_name
		, col.ordinal_position
		, col.is_nullable
		, coalesce(col.domain_name, col.data_type) -- not all dbs have domains
	FROM information_schema.columns AS col
	WHERE lower(table_schema) = 'information_schema'
	ORDER BY
		col.table_name
		, col.column_name
		, col.ordinal_position
		, col.is_nullable
		, col.domain_name
		, col.data_type
	;
`

func main() {
	version := os.Args[1]
	db, err := observer.WaitFor(driver, dsn, 100)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	log.Infof("connected to postgres %s", version)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := tx.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	for rows.Next() {
		var table, column, nullable, type_ string
		var number int

		err := rows.Scan(&table, &column, &number, &nullable, &type_)
		if err != nil {
			log.Fatal(err)
		}
		writer.Write([]string{table, column, fmt.Sprintf("%d", number), nullable, type_})
	}
}
