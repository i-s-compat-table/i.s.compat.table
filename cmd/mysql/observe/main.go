package main

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/go-sql-driver/mysql"
	"github.com/i-s-compat-table/i.s.compat.table/internal/observer"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/internal/schema"
)

const outputPath = "./data/mysql/observed.sqlite"
const driver = "mysql"
const dsnTemplate = "root:password@tcp(127.0.0.1:%d)/"

var dbRecord = &commonSchema.Database{Name: "mysql"}
var versionPorts = map[string]int{
	// needs to keep in sync with docker-compose.yaml
	"5.7": 3306,
	"8.0": 3307,
}

type myLogger struct {
	logs []interface{}
}

func (m *myLogger) Print(logs ...interface{}) {
	m.logs = append(m.logs, logs...)
}

// there should already be one or more Mysql's running
func main() {
	colChan := make(chan []commonSchema.ColVersion)
	waitForWrites := sync.WaitGroup{}
	waitForWrites.Add(1)
	go commonSchema.BulkInsert(outputPath, colChan, &waitForWrites)
	waitForObservations := sync.WaitGroup{}
	logger := myLogger{}
	mysql.SetLogger(&logger)
	for version, port := range versionPorts {
		waitForObservations.Add(1)
		go func(version string, portNumber int) {
			defer waitForObservations.Done()
			order := commonSchema.AsOrder(version)
			dbVersion := &commonSchema.Version{Db: dbRecord, Version: version, Order: &order}
			dsn := fmt.Sprintf(dsnTemplate, portNumber)
			db, err := observer.WaitFor(driver, dsn, 100)
			if err != nil {
				for _, message := range logger.logs {
					log.Debug(message)
				}
				log.Panicf("failed to connect to %s: %v", dsn, err)
			}
			colChan <- observer.Observe(db, dbVersion, nil)
		}(version, port)
	}
	waitForObservations.Wait()
	close(colChan)
	waitForWrites.Wait()
}
