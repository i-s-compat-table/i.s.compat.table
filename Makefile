SCRAPER_DBS= ./data/mariadb/docs/db.sqlite
SCRAPER_DBS+=./data/mssql/docs/db.sqlite
SCRAPER_DBS+=./data/postgres/docs/db.sqlite

TSVS= ./data/postgres.tsv
TSVS+=./data/mariadb.tsv
TSVS+=./data/mssql.tsv


# build the scraper binaries
./bin/scrape_mariadb_docs: ./pkg/dbs/mariadb/docs/scraper.go pkg/common/schema/db.go pkg/common/schema/db.sql
	go build -o ./bin/scrape_mariadb_docs ./cmd/scrape_mariadb_docs/main.go
./bin/scrape_mssql_docs: ./pkg/dbs/mssql/docs/scraper.go pkg/common/schema/db.go pkg/common/schema/db.sql
	go build -o ./bin/scrape_mssql_docs ./cmd/scrape_mssql_docs/main.go
./bin/scrape_postgres_docs: ./pkg/dbs/postgres/docs/scraper.go pkg/common/schema/db.go pkg/common/schema/db.sql
	go build -o ./bin/scrape_postgres_docs ./cmd/scrape_postgres_docs/main.go

# run the scaper binaries
./data/mariadb/docs.sqlite: ./bin/scrape_mariadb_docs
	./bin/scrape_mariadb_docs
./data/mssql/docs.sqlite: ./bin/scrape_mssql_docs
	./bin/scrape_mssql_docs
./data/postgres/docs.sqlite: ./bin/scrape_postgres_docs
	./bin/scrape_postgres_docs

.PHONY: doc_dbs
doc_dbs=./data/mariadb/docs.sqlite ./data/mssql/docs.sqlite ./data/postgres/docs.sqlite
doc-dbs: $(doc_dbs)

./data/columns.sqlite: ./db.sql $(TSVS)
	./scripts/assemble.sh


# TODO: maybe split the table into a separate markdown file?
README: # should depend on the ouput db

.PHONY: README clean-scraped-docs


clean-scraped-docs:
	rm -f $(doc_dbs)

clean:
	rm -f $(SCRAPER_DBS)
