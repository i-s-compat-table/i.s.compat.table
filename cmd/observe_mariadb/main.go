package main

import (
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/i-s-compat-table/i.s.compat.table/pkg/common/observer"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/common/schema"
)

const outputPath = "./data/mariadb/observed.sqlite"
const driver = "mysql"
const dsnTemplate = "root:password@tcp(127.0.0.1:%d)/"

var dbRecord = &commonSchema.Database{Name: ""}
var versionPorts = map[string]int{
	// needs to keep in sync with docker-compose.yaml
	// "10.2": 3038,
	// "10.3": 3039,
	// "10.4": 3040,
	// "10.5": 3041,
	// "10.6": 3042,
	"10.7": 3043,
}

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
			colChan <- observer.Observe(db, dbVersion)
		}(version, port)
	}
	waitForObservations.Wait()
	close(colChan)
	waitForWrites.Wait()
}
