package main

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/i-s-compat-table/i.s.compat.table/pkg/observer"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/schema"
)

const outputPath = "./data/clickhouse/observed.sqlite"
const host = "127.0.0.1:9000"
const _version = "22.8.5"

var dbRecord = &commonSchema.Database{Name: "clickhouse"}
var isCurrent = true
var order = commonSchema.AsOrder(_version)
var dbVersion = &commonSchema.Version{
	Db: dbRecord, IsCurrent: &isCurrent, Version: _version, Order: &order,
}

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
	return db
}
func waitFor(db *sql.DB, retries int) {
	ticker := time.NewTicker(time.Second)
	var err error
	for i := 0; i <= retries; i++ {
		if err = db.Ping(); err == nil {
			return
		}
		<-ticker.C
	}
	log.Panicf("failed to connect to %s after %d seconds", host, retries)
}

// there should already be one or more Mysql's running
func main() {
	colChan := make(chan []commonSchema.ColVersion)
	waitForWrites := sync.WaitGroup{}
	waitForWrites.Add(1)
	go commonSchema.BulkInsert(outputPath, colChan, &waitForWrites)
	db := openConn()
	waitFor(db, 100)
	colChan <- observer.Observe(db, dbVersion, nil)
	close(colChan)
	waitForWrites.Wait()
}
