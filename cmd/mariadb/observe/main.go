package main

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/go-sql-driver/mysql"
	"github.com/i-s-compat-table/i.s.compat.table/pkg/observer"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/schema"
)

const outputPath = "./data/mariadb/observed.sqlite"
const driver = "mysql"
const dsnTemplate = "root:password@tcp(127.0.0.1:%d)/"

var dbRecord = &commonSchema.Database{Name: "mariadb"}
var versionPorts = map[string]int{
	// needs to keep in sync with docker-compose.yaml
	"10.2.41": 3308,
	"10.3.32": 3309,
	"10.4.22": 3310,
	"10.5.13": 3311,
	"10.6.5":  3312,
	"10.7.1":  3313,
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
			dbVersion := &commonSchema.Version{Db: dbRecord, Version: version}
			dsn := fmt.Sprintf(dsnTemplate, portNumber)
			db, err := observer.WaitFor(driver, dsn, 30)
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
