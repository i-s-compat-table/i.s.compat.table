package main

import (
	"fmt"
	"sync"

	"github.com/i-s-compat-table/i.s.compat.table/pkg/observer"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/schema"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const outputPath = "./data/postgres/observed.sqlite"
const driver = "postgres"
const dsnTemplate = "user=postgres password=password host=localhost sslmode=disable port=%d"

var dbRecord = &commonSchema.Database{Name: "postgres"}
var versionPorts = map[string]int{
	// needs to keep in sync with docker-compose.yaml
	"10": 5432,
	"11": 5433,
	"12": 5434,
	"13": 5435,
	"14": 5436,
}
var query = `SELECT
col.table_name
, col.column_name
, col.ordinal_position
, col.is_nullable
, coalesce(col.domain_name, col.data_type) -- not all dbs have domains
FROM information_schema.columns AS col
WHERE lower(table_schema) = 'information_schema';`

// there should already be one or more Mysql's running
func main() {
	colChan := make(chan []commonSchema.ColVersion)
	waitForWrites := sync.WaitGroup{}
	waitForWrites.Add(1)
	go commonSchema.BulkInsert(outputPath, colChan, &waitForWrites)
	waitForObservations := sync.WaitGroup{}
	for version, port := range versionPorts {
		waitForObservations.Add(1)
		go func(version string, portNumber int) {
			defer waitForObservations.Done()
			dbVersion := &commonSchema.Version{Db: dbRecord, Version: version}
			dsn := fmt.Sprintf(dsnTemplate, portNumber)
			db := observer.WaitFor(driver, dsn)
			log.Infof("connected to postgres %s on port %d", version, portNumber)
			colChan <- observer.Observe(db, dbVersion, &query)
		}(version, port)
	}
	waitForObservations.Wait()
	close(colChan)
	waitForWrites.Wait()
}
