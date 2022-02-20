# build the scraper binaries
.PHONY: all
all: ./data/columns.tsv ./data/mariadb/columns.tsv ./data/mssql/columns.tsv ./data/mysql/columns.tsv ./data/postgres/columns.tsv
_common_backend=pkg/schema/db.go pkg/schema/db.sql
./bin/scrape_mariadb_docs: ./cmd/mariadb/scrape_docs/main.go $(_common_backend)
	go build -o ./bin/scrape_mariadb_docs ./cmd/mariadb/scrape_docs/main.go
./bin/scrape_mssql_docs: ./cmd/mssql/scrape_docs/main.go $(_common_backend)
	go build -o ./bin/scrape_mssql_docs ./cmd/mssql/scrape_docs/main.go
./bin/scrape_postgres_docs: ./cmd/postgres/scrape_docs/main.go $(_common_backend)
	go build -o ./bin/scrape_postgres_docs ./cmd/postgres/scrape_docs/main.go

scraper_binaries=./bin/scrape_mariadb_docs ./bin/scrape_mssql_docs ./bin/scrape_postgres_docs
# build the observer binaries
_observer_common=pkg/observer/observer.go pkg/observer/columns.sql
./bin/observe_mariadb: ./cmd/mariadb/observe/main.go $(_common_backend) $(_observer_common)
	go build -o ./bin/observe_mariadb ./cmd/mariadb/observe/main.go
./bin/observe_mysql: ./cmd/mysql/observe/main.go $(_common_backend) $(_observer_common)
	go build -o ./bin/observe_mysql ./cmd/mysql/observe/main.go
./bin/observe_postgres: ./cmd/postgres/observe/main.go $(_common_backend) $(_observer_common)
	go build -o ./bin/observe_postgres ./cmd/postgres/observe/main.go

observer_binaries=./bin/observe_mariadb ./bin/observe_mssql ./bin/observe_mysql ./bin/observe_postgres

# run the scaper binaries
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

pg-docs: ./data/postgres/docs.sqlite
./data/postgres/docs.sqlite: ./bin/scrape_postgres_docs
	mkdir -p ./data/postgres
	rm -f ./data/postgres/docs.sqlite
	./bin/scrape_postgres_docs


# run the observer binaries
mariadb_services=mariadb-10.2.41
mariadb_services+=mariadb-10.3.32
mariadb_services+=mariadb-10.4.22
mariadb_services+=mariadb-10.5.13
mariadb_services+=mariadb-10.6.5
mariadb_services+=mariadb-10.7.1
mariadb-observations:./data/mariadb/observed.sqlite
./data/mariadb/observed.sqlite:./bin/observe_mariadb
	mkdir -p ./data/mariadb
	rm -f ./data/mariadb/observed.sqlite
	docker-compose up -d $(mariadb_services)
	./bin/observe_mariadb
	docker-compose down

mysql-observations: ./data/mysql/observed.sqlite
./data/mysql/observed.sqlite: ./bin/observe_mysql
	mkdir -p ./data/mysql
	rm -f ./data/mysql/observed.sqlite
	docker-compose up -d mysql-5.7 mysql-8.0
	./bin/observe_mysql
	docker-compose down

pg-observations: ./data/postgres/observed.sqlite
pg_services=postgres-10 postgres-11 postgres-12 postgres-13 postgres-14
./data/postgres/observed.sqlite:./bin/observe_postgres
	mkdir -p ./data/postgres
	rm -f ./data/postgres/observed.sqlite
	docker-compose up -d $(pg_services)
	./bin/observe_postgres
	docker-compose down

.PHONY: doc_dbs
doc_dbs=./data/mariadb/docs.sqlite ./data/mssql/docs.sqlite ./data/postgres/docs.sqlite
doc-dbs: $(doc_dbs)
merge_scripts=./scripts/merge/dbs.sh ./scripts/merge/merge.sql
./data/merged.docs.sqlite: $(merge_scripts) $(doc_dbs)
	rm -f ./data/merged.docs.sqlite
	./scripts/merge/dbs.sh ./data/merged.docs.sqlite $(doc_dbs)

observation_dbs= ./data/mariadb/observed.sqlite
observation_dbs+=./data/mysql/observed.sqlite
observation_dbs+=./data/postgres/observed.sqlite

./data/merged.observations.sqlite: $(merge_scripts) $(observation_dbs)
	./scripts/merge/dbs.sh ./data/merged.observations.sqlite $(observation_dbs)
	touch -m ./data/merged.observations.sqlite

./data/mariadb/merged.sqlite: $(merge_scripts) ./data/mariadb/docs.sqlite ./data/mariadb/observed.sqlite
	rm -f ./data/mariadb/merged.sqlite
	./scripts/merge/dbs.sh ./data/mariadb/merged.sqlite ./data/mariadb/observed.sqlite ./data/mariadb/docs.sqlite
# ./data/mssql/merged.sqlite: ./data/mssql/docs.sqlite
# 	ln -sf 
# ./data/mysql/merged.sqlite: ./data/mysql/observed.sqlite
# 	ln -sf ./data/mysql/merged.sqlite ./data/mysql/observed.sqlite
./data/postgres/merged.sqlite: $(merge_scripts) ./data/postgres/docs.sqlite ./data/postgres/observed.sqlite
	rm -f ./data/postgres/merged.sqlite
	./scripts/merge/dbs.sh ./data/postgres/merged.sqlite ./data/postgres/observed.sqlite ./data/postgres/docs.sqlite

tsv_dump_scripts= ./scripts/dump_tsv.sh ./pkg/schema/views.sql
./data/mssql/columns.tsv: $(tsv_dump_scripts) ./data/mssql/docs.sqlite
	./scripts/dump_tsv.sh --output=./data/mssql/columns.tsv ./data/mssql/docs.sqlite
./data/postgres/columns.tsv: $(tsv_dump_scripts) ./data/postgres/merged.sqlite
	./scripts/dump_tsv.sh --output=./data/postgres/columns.tsv ./data/postgres/merged.sqlite
./data/mysql/columns.tsv: $(tsv_dump_scripts) ./data/mysql/observed.sqlite
	./scripts/dump_tsv.sh --output=./data/mysql/columns.tsv ./data/mysql/observed.sqlite
./data/mariadb/columns.tsv: $(tsv_dump_scripts) ./data/mariadb/merged.sqlite
	./scripts/dump_tsv.sh --output=./data/mariadb/columns.tsv ./data/mariadb/merged.sqlite
./data/columns.sqlite: $(merge_scripts) ./data/merged.observations.sqlite ./data/merged.docs.sqlite
	./scripts/merge/dbs.sh ./data/columns.sqlite ./data/merged.observations.sqlite ./data/merged.docs.sqlite
	touch -m ./data/columns.sqlite

./data/columns.tsv: $(tsv_dump_scripts) ./data/columns.sqlite
	./scripts/dump_tsv.sh --output=./data/columns.tsv ./data/columns.sqlite

# TODO: use cog to write a mermaid+markdown ERD to CONTRIBUTING.md
# TODO: create html/markdown tables out of columns.tsv or columns.sqlite


all_shell_scripts=$(shell find . -type f -name '*.sh')
.PHONY: shellcheck
shellcheck:
	shellcheck -x $(all_shell_scripts)

.PHONY: clean-scraped-docs clean-observations clean-scraper-binaries clean-observer-binaries
clean-scraped-docs:
	rm -f $(doc_dbs)
clean-observations:
	rm -f $(observation_dbs)
clean-scraper-binaries:
	rm -f $(scraper_binaries)
clean-observer-binaries:
	rm -f $(observer_binaries)

clean: clean-scraped-docs clean-observations clean-scraper-binaries clean-observer-binaries
