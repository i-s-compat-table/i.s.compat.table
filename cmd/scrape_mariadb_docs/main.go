package main

import "github.com/i-s-compat-table/i.s.compat.table/pkg/dbs/mariadb/docs"

func main() {
	docs.Scrape("./pkg/dbs/mariadb/.cache", "./data/mariadb/docs/db.sqlite", false)
}
