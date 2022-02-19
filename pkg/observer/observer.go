package observer

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/schema"
	log "github.com/sirupsen/logrus"
)

var note = &commonSchema.Note{License: &commonSchema.License{License: "CC0-1.0"}}

//go:embed columns.sql
var infoSchemaColumnsQuery string

// TODO: nullable query parameter
func Observe(db *sql.DB, dbVersion *commonSchema.Version, query *string) []commonSchema.ColVersion {
	var getInfoSchemaColumns string
	if query == nil {
		getInfoSchemaColumns = infoSchemaColumnsQuery
	} else {
		getInfoSchemaColumns = *query
	}
	rows, err := db.Query(getInfoSchemaColumns)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	cols := []commonSchema.ColVersion{}
	for rows.Next() {
		table := commonSchema.Table{}
		column := commonSchema.Column{Table: &table}
		type_ := commonSchema.Type{}
		colVersion := commonSchema.ColVersion{
			DbVersion: dbVersion,
			Notes:     note,
			Column:    &column,
			Type:      &type_,
		}
		var nullable string = ""
		err := rows.Scan(
			&table.Name, &column.Name, &colVersion.Number, &nullable, &type_.Name,
		)
		if err != nil {
			panic(err)
		}
		colVersion.Nullable = commonSchema.FromString(nullable)
		cols = append(cols, colVersion)
	}
	log.Infof("%s %s: found %d columns", dbVersion.Db.Name, dbVersion.Version, len(cols))
	return cols
}

func WaitFor(driverName, dsn string) (db *sql.DB) {
	log.SetLevel(log.DebugLevel)
	ticker := time.NewTicker(time.Second)
	var finalErr error
	for i := 0; i <= 15; i++ {
		<-ticker.C // wait for a tick

		if db, err := sql.Open(driverName, dsn); err == nil {
			// db.SetMaxIdleConns(0)
			// db.SetConnMaxLifetime(time.Second * 1000)
			// db.SetConnMaxIdleTime(time.Second * 1000)
			if err := db.Ping(); err == nil {
				log.Infof("connected to %s", dsn)
				return db
			} else {
				finalErr = err
				log.Debugf("%s >> %+v", dsn, err)
				// fmt.Printf(".")
			}
		} else {
			finalErr = err
			log.Debugf("%+v", err)
			fmt.Printf("_")
		}
	}
	log.Panicf("unable to connect to %s : %v", dsn, finalErr)
	return nil
}
