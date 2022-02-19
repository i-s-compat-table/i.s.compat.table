# build the scraper binaries
_common_backend=pkg/schema/db.go pkg/schema/db.sql
./bin/scrape_mariadb_docs: ./cmd/mariadb/scrape_docs/main.go $(_common_backend)
	go build -o ./bin/scrape_mariadb_docs ./cmd/mariadb/scrape_docs/main.go
./bin/scrape_mssql_docs: ./cmd/mssql/scrape_docs/main.go $(_common_backend)
	go build -o ./bin/scrape_mssql_docs ./cmd/mssql/scrape_docs/main.go
./bin/scrape_postgres_docs: ./cmd/postgres/scrape_docs/main.go $(_common_backend)
	go build -o ./bin/scrape_postgres_docs ./cmd/postgres/scrape_docs/main.go

# build the observer binaries
_observer_common=pkg/observer/observer.go pkg/observer/columns.sql
./bin/observe_mariadb: ./cmd/mariadb/observe/main.go $(_common_backend) $(_observer_common)
	go build -o ./bin/observe_mariadb ./cmd/mariadb/observe/main.go
./bin/observe_mssql: ./cmd/mssql/observe/main.go $(_common_backend) $(_observer_common)
	go build -o ./bin/observe_mssql ./cmd/mssql/observe/main.go
./bin/observe_mysql: ./cmd/mysql/observe/main.go $(_common_backend) $(_observer_common)
	go build -o ./bin/observe_mysql ./cmd/mysql/observe/main.go
./bin/observe_postgres: ./cmd/postgres/observe/main.go $(_common_backend) $(_observer_common)
	go build -o ./bin/observe_postgres ./cmd/postgres/observe/main.go

# run the scaper binaries
mariadb-docs: ./data/mariadb/docs.sqlite
./data/mariadb/docs.sqlite: ./bin/scrape_mariadb_docs
	./bin/scrape_mariadb_docs
	touch -m ./data/mariadb/docs.sqlite
mssql-docs: ./data/mssql/docs.sqlite
./data/mssql/docs.sqlite: ./bin/scrape_mssql_docs
	./bin/scrape_mssql_docs
	touch -m ./data/mssql/docs.sqlite
pg-docs: ./data/postgres/docs.sqlite
./data/postgres/docs.sqlite: ./bin/scrape_postgres_docs
	./bin/scrape_postgres_docs
	touch -m ./data/postgres/docs.sqlite


# run the observer binaries
mariadb_services=mariadb-10.2
mariadb_services+=mariadb-10.3
mariadb_services+=mariadb-10.4
mariadb_services+=mariadb-10.5
mariadb_services+=mariadb-10.6
mariadb_services+=mariadb-10.7
mariadb-observations:./data/mariadb/observed.sqlite
./data/mariadb/observed.sqlite:./bin/observe_mariadb
	docker-compose up -d $(mariadb_services)
	./bin/observe_mariadb
	docker-compose down

mysql-observations: ./bin/observe_mysql
	docker-compose up -d mysql-5.7 mysql-8.0
	./bin/observe_mysql
	docker-compose down

pg-observations: ./data/postgres/observed.sqlite
pg_services=postgres-10 postgres-11 postgres-12 postgres-13 postgres-14
./data/postgres/observed.sqlite:./bin/observe_postgres
	docker-compose up -d $(pg_services)
	./bin/observe_postgres
	docker-compose down

.PHONY: doc_dbs
doc_dbs=./data/mariadb/docs.sqlite ./data/mssql/docs.sqlite ./data/postgres/docs.sqlite
doc-dbs: $(doc_dbs)
merge_scripts=./scripts/merge/dbs.sh ./scripts/merge/merge.sql
./data/merged.docs.sqlite: $(merge_scripts) $(doc_dbs)
	./scripts/merge/dbs.sh ./data/merged.docs.sqlite $(doc_dbs)
	touch -m ./data/merged.docs.sqlite

observation_dbs= ./data/mariadb/observed.sqlite
observation_dbs+=./data/mysql/observed.sqlite
observation_dbs+=./data/postgres/observed.sqlite

./data/merged.observations.sqlite: $(merge_scripts) $(observation_dbs)
	./scripts/merge/dbs.sh ./data/merged.observations.sqlite $(observation_dbs)
	touch -m ./data/merged.observations.sqlite

# TODO: maybe split the table into a separate markdown file?
README: # should depend on the ouput db

.PHONY: README clean-scraped-docs


clean-scraped-docs:
	rm -f $(doc_dbs)

clean:
	rm -f $(SCRAPER_DBS)
