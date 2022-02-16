package main

import (
	"github.com/i-s-compat-table/i.s.compat.table/pkg/dbs/postgres/docs"
)

func main() {
	docs.Scrape("./pkg/dbs/postgres/.cache", "./data/postgres/docs.sqlite", false)
}
