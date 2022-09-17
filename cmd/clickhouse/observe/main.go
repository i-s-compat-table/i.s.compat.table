package main

import (
	"database/sql"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/i-s-compat-table/i.s.compat.table/pkg/observer"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/schema"
)

const outputPath = "./data/clickhouse/observed.sqlite"
const host = "127.0.0.1:9000"

var dbRecord = &commonSchema.Database{Name: "clickhouse"}
var isCurrent = true
var dbVersion = &commonSchema.Version{Db: dbRecord, IsCurrent: &isCurrent, Version: "22.8.5"}

func openConn() *sql.DB {
	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{host},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		TLS: nil,
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
			Level:  1,
		},
		Debug: true,
	})
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}

// there should already be one or more Mysql's running
func main() {
	colChan := make(chan []commonSchema.ColVersion)
	waitForWrites := sync.WaitGroup{}
	waitForWrites.Add(1)
	go commonSchema.BulkInsert(outputPath, colChan, &waitForWrites)
	db := openConn()
	if db == nil {
		log.Panicf("failed to connect to %s", host)
	}
	colChan <- observer.Observe(db, dbVersion, nil)
	close(colChan)
	waitForWrites.Wait()
}
