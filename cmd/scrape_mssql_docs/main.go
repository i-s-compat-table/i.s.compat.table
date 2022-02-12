package main

import "github.com/i-s-compat-table/i.s.compat.table/pkg/dbs/mssql/docs"

func main() {
	docs.Scrape("./pkg/dbs/mssql/.cache", "./data/mssql/docs/db.sqlite", false)
}
