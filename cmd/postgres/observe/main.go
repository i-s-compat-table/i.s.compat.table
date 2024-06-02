package main

import (
	_ "embed"
	"fmt"
	"sync"

	"github.com/i-s-compat-table/i.s.compat.table/internal/observer"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/internal/schema"
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
	"15": 5437,
}

//go:embed query.sql
var query string

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
			order := commonSchema.AsOrder(version)
			dbVersion := &commonSchema.Version{Db: dbRecord, Version: version, Order: &order}
			dsn := fmt.Sprintf(dsnTemplate, portNumber)
			db, err := observer.WaitFor(driver, dsn, 100)
			if err != nil {
				panic(err)
			}
			log.Infof("connected to postgres %s on port %d", version, portNumber)
			colChan <- observer.Observe(db, dbVersion, &query)
		}(version, port)
	}
	waitForObservations.Wait()
	close(colChan)
	waitForWrites.Wait()
}
