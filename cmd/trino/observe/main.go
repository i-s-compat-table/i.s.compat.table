package main

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/i-s-compat-table/i.s.compat.table/pkg/observer"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/schema"
	_ "github.com/trinodb/trino-go-client/trino"
)

const outputPath = "./data/trino/observed.sqlite"
const driver = "trino"
const dsn = "http://user@127.0.0.1:8080?catalog=tpch"

var dbRecord = &commonSchema.Database{Name: "trino"}
var isCurrent = true
var dbVersion = &commonSchema.Version{Db: dbRecord, IsCurrent: &isCurrent, Version: "396"}

var query = `SELECT
    lower(col.table_name)
  , lower(col.column_name)
  , col.ordinal_position
  , col.is_nullable
  , upper(col.data_type)
FROM information_schema.columns AS col
WHERE lower(table_schema) = 'information_schema'`

// there should already be one or more Mysql's running
func main() {
	colChan := make(chan []commonSchema.ColVersion)
	waitForWrites := sync.WaitGroup{}
	waitForWrites.Add(1)
	go commonSchema.BulkInsert(outputPath, colChan, &waitForWrites)
	db, err := observer.WaitFor(driver, dsn, 100)
	if err != nil {
		log.Panicf("failed to connect to %s: %v", dsn, err)
	}
	colChan <- observer.Observe(db, dbVersion, &query)
	close(colChan)
	waitForWrites.Wait()
}
