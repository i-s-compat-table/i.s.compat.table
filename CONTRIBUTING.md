# Contributing

## Reporting incorrect information

<!-- where? issue tracker -->

As with all open-source, the only way to ensure a fix gets written is to write it yourself.

## Adding a new database

Gathering and maintaining correct, license-compliant information on a new database would take a significant amount of manual work.
Thus, adding a new database requires writing at least one of a documentation-scraper or a database-observer.

### Repo structure

Directories, files in the order you should visit them:

```tree
.
├── README.md
├── CONTRIBUTING.md
├── pkg
│   └── common
│       ├── observer/*.{sql,go}
│       ├── schema/*.{sql,go}
│       └── utils/*.go
├── Makefile
├── cmd
│   └── ${DB}
│       ├── observer/*.go
│       └── scrape_docs/*.go
├── bin/* # built binaries
├── scripts
│   ├── scrape/${DB}_docs.sh
│   ├── observe/${DB}.sh
│   └── *.sh
├── data
│   ├── README.md
│   ├── ${DB}
│   │   ├── docs.sqlite
│   │   ├── patch.sql
│   │   ├── observed.sqlite
│   │   ├── merged.sqlite
│   │   └── columns.tsv
│   ├── columns.sqlite
│   └── columns.tsv
└── dist/* # future web artifacts
```

### required development tools

- `go >= 1.17`
- `sqlite3`
- `make`
- `bash`
- `docker`
- `docker-compose` (version 2 with the bridge should be ok)
- `python3.9`
- `poetry`
- several GiB of storage for database docker images.

### Workflow

```mermaid
graph TD
  remote_docs[remote docs]  -->|scrape | docs.sqlite
  patch.sql                 -->|correct| docs.sqlite

  docker_db[live database] --> |observe| observed.sqlite
  observed.sqlite          --> |merge  | merged.sqlite
  docs.sqlite              --> |merge  | merged.sqlite
  merged.sqlite            --> |render | columns.tsv
  columns.tsv              --> |inform | patch.sql
```

Each sqlite database containing data on `information_schema` implementations has the following schema, located in [`./pkg/schema/db.sql`](./pkg/schema/db.sql):

<!-- [[[cog
  from scripts.get_mermaid_erd import get_mermaid_erd
  print("```mermaid")
  print(get_mermaid_erd())
  print("```")
]]] -->

```mermaid
erDiagram
  column_versions {
    INTEGER id PK
    INTEGER column_id FK
    INTEGER version_id FK
    INTEGER type_id FK
    BOOLEAN nullable
    INTEGER url_id FK
    INTEGER note_id FK
    INTEGER note_license_id FK
  }
  column_versions }o--|| licenses : note_license_id
  column_versions }o--|| notes : note_id
  column_versions }o--|| urls : url_id
  column_versions }o--|| types : type_id
  column_versions }|--|| versions : version_id
  column_versions }|--|| columns : column_id
  columns {
    INTEGER id PK
    INTEGER table_id FK
    TEXT name
  }
  columns }|--|| tables : table_id
  dbs {
    INTEGER id PK
    TEXT name
  }
  licenses {
    INTEGER id PK
    TEXT license
    TEXT attribution
    INTEGER link_id FK
  }
  licenses }o--|| urls : link_id
  notes {
    INTEGER id PK
    TEXT note
  }
  tables {
    INTEGER id PK
    TEXT name
  }
  types {
    INTEGER id PK
    TEXT name
  }
  urls {
    INTEGER id PK
    TEXT url
  }
  versions {
    INTEGER id PK
    INTEGER db_id FK
    TEXT version
    TEXT release_date
    BOOLEAN is_current
  }
  versions }|--|| dbs : db_id
```

<!-- [[[end]]] -->

### When to add a documentation scraper

If the database's documentation is offered under a creative commons license such as `cc0`, `CC-BY`, or `CC-BY-SA`, consider writing a new scraper for it.

### How to create a new documentation scraper

<!-- TODO: script templating new documentation scrapers! -->
<!-- TODO: use `cog` to read in the list from the PR template -->

- [ ] Check that the documentation is licensed under a Creative Commons-compatible license such as `CC BY 4.0` or `CY BY-SA 4.0`
- [ ] Create a scraper script at `cmd/scrape_${DB}_docs/main.go`
- [ ] Create a target in `Makefile` to compile `cmd/${DB}/scrape_docs/main.go` into `bin/scrape_${DB}_docs`
- [ ] Create a target in `Makefile` to run your scraper to create `data/${DB}/docs.sqlite`
- [ ] Update the other targets in `Makefile` to include the new dataset (i.e. `./data/${DB}/merged.sqlite` if it exists, `./data/merged.docs.sqlite`)
- [ ] Update [`./.github/workflows/scrape.yaml`](./.github/workflows/scrape.yaml) with a new doc-scraping job
- [ ] Update the README to list support for the new database!
<!-- - [ ] Commit the changes to `./data/${DB}/columns.tsv` and `./data/columns.tsv` -->

```diff
  .
  ├── cmd
+ │   └── scrape_${DB}_docs/main.go
+ ├── Makefile # add bin, docs.sqlite targets
  ├── bin
  │   ├── README.md
+ │   └── scrape_${DB}_docs # generated
  ├── data
+ │   └── ${DB}
+ │       └── docs
+ │           ├── docs.sqlite # generated
+ │           └── patch.sql
+ ├── README.md
```

Finally, create a pull request. You'll be reminded of all the steps above.

### When to create a new database-observer

When anyone has the right to spin up a local instance of the database server.
In the future, we may create observers for remote commercial databases that offer a free tier.

### How to create a new database-observer

- [ ] Check that the database docker image is licensed and available for unrestricted use
- [ ] add services to the dockerfile
- [ ] Create a target in `Makefile` to compile `cmd/${DB}/observe/main.go` into `bin/observe_${DB}`
- [ ] Create a target in `Makefile` to run your observer to create `data/${DB}/observed.sqlite`
- [ ] Update the other targets in `Makefile` to include the new dataset (i.e. `./data/${DB}/merged.sqlite` if it exists, `./data/merged.docs.sqlite`)
- [ ] Update [`./.github/workflows/scrape.yaml`](./.github/workflows/scrape.yaml) with new binary-compilation steps and a new db-observation job
<!-- - [ ] Commit the changes to `./data/${DB}/columns.tsv` and `./data/columns.tsv` -->
