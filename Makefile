# variables ###################################################################
# try to keep each list in alphabetical order
_common_backend=internal/schema/db.go internal/schema/db.sql
COMPOSE ?= ${COMPOSE}
observer_binaries=\
	./bin/observe_mariadb \
	./bin/observe_mssql \
	./bin/observe_mysql \
	./bin/observe_postgres

scraper_binaries=\
	./bin/scrape_cockroachdb_docs \
	./bin/scrape_mariadb_docs \
	./bin/scrape_mssql_docs \
	./bin/scrape_postgres_docs

mariadb_services=\
	mariadb-10.2.41 \
	mariadb-10.3.32 \
	mariadb-10.4.22 \
	mariadb-10.5.13 \
	mariadb-10.6.5 \
	mariadb-10.7.1 \

doc_dbs=\
	./data/cockroachdb/docs.sqlite \
	./data/mariadb/docs.sqlite \
	./data/mssql/docs.sqlite \
	./data/postgres/docs.sqlite \
	./data/tidb/docs.sqlite \

observation_dbs=./data/mariadb/observed.sqlite \
	./data/mysql/observed.sqlite \
	./data/postgres/observed.sqlite \
	./data/trino/observed.sqlite \
	./data/clickhouse/observed.sqlite \

tsv_dump_scripts=./scripts/dump_tsv.sh ./internal/schema/views.sql

# phony targets --------------------------------------------------------------
# try to keep these in alphabetical order
.PHONY: all \
	clean \
	clean-merged-dbs \
	clean-observer-binaries  \
	clean-observations \
	clean-scraped-docs \
	clean-scraper-binaries \
	doc-dbs \
	shellcheck \
	update-docs

# targets ====================================================================
.PHONY: all
all: ./data/columns.tsv ./data/mariadb/columns.tsv ./data/mssql/columns.tsv ./data/mysql/columns.tsv ./data/postgres/columns.tsv

# build the scraper binaries -------------------------------------------------
./bin/scrape_cockroachdb_docs: ./cmd/cockroachdb/scrape_docs/main.go $(_common_backend)
	go build -o ./bin/scrape_cockroachdb_docs ./cmd/cockroachdb/scrape_docs/main.go
./bin/scrape_mariadb_docs: ./cmd/mariadb/scrape_docs/main.go $(_common_backend)
	go build -o ./bin/scrape_mariadb_docs ./cmd/mariadb/scrape_docs/main.go
./bin/scrape_mssql_docs: ./cmd/mssql/scrape_docs/main.go $(_common_backend)
	go build -o ./bin/scrape_mssql_docs ./cmd/mssql/scrape_docs/main.go
./bin/scrape_postgres_docs: ./cmd/postgres/scrape_docs/main.go $(_common_backend)
	go build -o ./bin/scrape_postgres_docs ./cmd/postgres/scrape_docs/main.go
./bin/scrape_tidb_docs: ./cmd/tidb/scrape_docs/main.go $(_common_backend)
	go build -o ./bin/scrape_tidb_docs ./cmd/tidb/scrape_docs/main.go
# build the observer binaries ------------------------------------------------
_observer_common=internal/observer/observer.go internal/observer/columns.sql
./bin/observe_mariadb: ./cmd/mariadb/observe/main.go $(_common_backend) $(_observer_common)
	go build -o ./bin/observe_mariadb ./cmd/mariadb/observe/main.go
./bin/observe_mysql: ./cmd/mysql/observe/main.go $(_common_backend) $(_observer_common)
	go build -o ./bin/observe_mysql ./cmd/mysql/observe/main.go
./bin/observe_postgres: ./cmd/postgres/observe/main.go $(_common_backend) $(_observer_common)
	go build -o ./bin/observe_postgres ./cmd/postgres/observe/main.go
./bin/observe_trino: ./cmd/trino/observe/main.go $(_common_backend) $(_observer_common)
	go build -o ./bin/observe_trino ./cmd/trino/observe/main.go
./bin/observe_clickhouse: ./cmd/clickhouse/observe/main.go $(_common_backend) $(_observer_common)
	go build -o ./bin/observe_clickhouse ./cmd/clickhouse/observe/main.go

# run the scaper binaries ----------------------------------------------------
cockroachdb-docs: ./data/cockroachdb/docs.sqlite
./data/cockroachdb/docs.sqlite: ./bin/scrape_cockroachdb_docs
	mkdir -p ./data/cockroachdb
	rm -f ./data/cockroachdb/docs.sqlite
	./bin/scrape_cockroachdb_docs

mariadb-docs: ./data/mariadb/docs.sqlite
./data/mariadb/docs.sqlite: ./bin/scrape_mariadb_docs
	mkdir -p ./data/mariadb
	rm -f ./data/mariadb/docs.sqlite
	./bin/scrape_mariadb_docs

mssql-docs: ./data/mssql/docs.sqlite
./data/mssql/docs.sqlite: ./bin/scrape_mssql_docs
	mkdir -p ./data/mssql
	rm -f ./data/mssql/docs.sqlite
	./bin/scrape_mssql_docs
	touch -m ./data/mssql/docs.sqlite

pg-docs: ./data/postgres/docs.sqlite
./data/postgres/docs.sqlite: ./bin/scrape_postgres_docs ./data/postgres/patch.sql
	mkdir -p ./data/postgres
	rm -f ./data/postgres/docs.sqlite
	./bin/scrape_postgres_docs
	sqlite3 ./data/postgres/docs.sqlite <./data/postgres/patch.sql
tidb-docs: ./data/tidb/docs.sqlite 
./data/tidb/docs.sqlite: ./bin/scrape_tidb_docs
	mkdir -p ./data/tidb
	rm -f ./data/tidb/docs.sqlite
	./bin/scrape_tidb_docs

# run the observer binaries --------------------------------------------------
mariadb-observations:./data/mariadb/observed.sqlite
./data/mariadb/observed.sqlite: ./bin/observe_mariadb
	mkdir -p ./data/mariadb
	rm -f ./data/mariadb/observed.sqlite
	${COMPOSE} up -d $(mariadb_services)
	./bin/observe_mariadb
	${COMPOSE} down
	touch -m ./data/mariadb/observed.sqlite

mysql-observations: ./data/mysql/observed.sqlite
./data/mysql/observed.sqlite: ./bin/observe_mysql
	mkdir -p ./data/mysql
	rm -f ./data/mysql/observed.sqlite
	${COMPOSE} up -d mysql-5.7 mysql-8.0
	./bin/observe_mysql
	touch -m ./data/mysql/observed.sqlite
	${COMPOSE} down

pg-observations: ./data/postgres/observed.sqlite
pg_services=\
	postgres-10\
	postgres-11\
	postgres-12\
	postgres-13\
	postgres-14\
	postgres-15\
	postgres-16\

./data/postgres/observed.sqlite:./bin/observe_postgres
	mkdir -p ./data/postgres
	rm -f ./data/postgres/observed.sqlite
	${COMPOSE} up -d $(pg_services)
	./bin/observe_postgres
	touch -m ./data/postgres/observed.sqlite
	${COMPOSE} down

trino-observations: ./data/trino/observed.sqlite
./data/trino/observed.sqlite: ./bin/observe_trino
	mkdir -p ./data/trino
	rm -rf ./data/trino/observed.sqlite
	${COMPOSE} up -d trino
	./bin/observe_trino
	touch -m ./data/trino/observed.sqlite
	${COMPOSE} down
clickhouse-observations: ./data/clickhouse/observed.sqlite
./data/clickhouse/observed.sqlite: ./bin/observe_clickhouse
	mkdir -p ./data/clickhouse
	rm -rf ./data/clickhouse/observed.sqlite
	${COMPOSE} up -d clickhouse
	./bin/observe_clickhouse
	${COMPOSE} down

# merge dataset as sqlite ----------------------------------------------------
doc-dbs: $(doc_dbs)
merge_scripts=./scripts/merge/dbs.sh ./scripts/merge/merge.sql
./data/merged.docs.sqlite: $(merge_scripts) $(doc_dbs)
	rm -f ./data/merged.docs.sqlite
	./scripts/merge/dbs.sh ./data/merged.docs.sqlite $(doc_dbs)
	touch -m ./data/merged.docs.sqlite

./data/merged.observations.sqlite: $(merge_scripts) $(observation_dbs)
	./scripts/merge/dbs.sh ./data/merged.observations.sqlite $(observation_dbs)
	touch -m ./data/merged.observations.sqlite

./data/mariadb/merged.sqlite: $(merge_scripts) ./data/mariadb/docs.sqlite ./data/mariadb/observed.sqlite
	rm -f ./data/mariadb/merged.sqlite
	./scripts/merge/dbs.sh ./data/mariadb/merged.sqlite ./data/mariadb/docs.sqlite ./data/mariadb/observed.sqlite

./data/postgres/merged.sqlite: $(merge_scripts) ./data/postgres/docs.sqlite ./data/postgres/observed.sqlite
	rm -f ./data/postgres/merged.sqlite
	./scripts/merge/dbs.sh ./data/postgres/merged.sqlite ./data/postgres/observed.sqlite ./data/postgres/docs.sqlite

# dump tsvs ------------------------------------------------------------------
./data/cockroachdb/columns.tsv: $(tsv_dump_scripts) ./data/cockroachdb/docs.sqlite
	./scripts/dump_tsv.sh --output ./data/cockroachdb/columns.tsv ./data/cockroachdb/docs.sqlite
./data/mariadb/columns.tsv: $(tsv_dump_scripts) ./data/mariadb/merged.sqlite
	./scripts/dump_tsv.sh --output ./data/mariadb/columns.tsv ./data/mariadb/merged.sqlite
./data/mssql/columns.tsv: $(tsv_dump_scripts) ./data/mssql/docs.sqlite
	./scripts/dump_tsv.sh --output ./data/mssql/columns.tsv ./data/mssql/docs.sqlite
./data/mysql/columns.tsv: $(tsv_dump_scripts) ./data/mysql/observed.sqlite
	./scripts/dump_tsv.sh --output ./data/mysql/columns.tsv ./data/mysql/observed.sqlite
./data/postgres/columns.tsv: $(tsv_dump_scripts) ./data/postgres/merged.sqlite
	./scripts/dump_tsv.sh --output ./data/postgres/columns.tsv ./data/postgres/merged.sqlite
./data/tidb/columns.tsv: $(tsv_dump_scripts) ./data/tidb/docs.sqlite
	./scripts/dump_tsv.sh --output ./data/tidb/columns.tsv ./data/tidb/docs.sqlite
./data/trino/columns.tsv: $(tsv_dump_scripts) ./data/trino/observed.sqlite
	./scripts/dump_tsv.sh --output ./data/trino/columns.tsv ./data/trino/observed.sqlite
./data/clickhouse/columns.tsv: $(tsv_dump_scripts) ./data/clickhouse/observed.sqlite
	./scripts/dump_tsv.sh --output ./data/clickhouse/columns.tsv ./data/clickhouse/observed.sqlite

./data/columns.sqlite: $(merge_scripts) ./data/merged.observations.sqlite ./data/merged.docs.sqlite
	./scripts/merge/dbs.sh ./data/columns.sqlite ./data/merged.observations.sqlite ./data/merged.docs.sqlite
	touch -m ./data/columns.sqlite
./data/columns.tsv: $(tsv_dump_scripts) ./data/columns.sqlite
	./scripts/dump_tsv.sh --output ./data/columns.tsv ./data/columns.sqlite

# TODO: create html/markdown tables out of columns.tsv or columns.sqlite

update-docs:
	poetry run cog -PUre ./CONTRIBUTING.md

shellcheck:
	shellcheck --external-sources --source-path=SCRIPTDIR $(shell find . -type f -name '*.sh')

clean-scraped-docs:
	rm -f $(doc_dbs)
clean-observations:
	rm -f $(observation_dbs)
clean-scraper-binaries:
	rm -f $(scraper_binaries)
clean-observer-binaries:
	rm -f $(observer_binaries)
clean-merged-dbs:
	find . -name 'merge*.sqlite' -type f | xargs rm -f ./data/columns.sqlite
clean: clean-scraped-docs clean-observations clean-merged-dbs clean-scraper-binaries clean-observer-binaries

./bin/debug_id: go.mod go.sum ./scripts/debug_id/*.go internal/schema/*
	go build -o bin/debug_id ./scripts/debug_id
