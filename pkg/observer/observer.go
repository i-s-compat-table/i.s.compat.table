package observer

import (
	"database/sql"
	_ "embed"
	"fmt"
	"strings"
	"time"

	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/schema"
	"github.com/i-s-compat-table/i.s.compat.table/pkg/utils"
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
		table.Name = strings.ToLower(utils.NormalizeString(table.Name))
		column.Name = strings.ToLower(utils.NormalizeString(column.Name))
		type_.Name = strings.ToUpper(utils.NormalizeString(type_.Name))
		if err != nil {
			panic(err)
		}
		colVersion.Nullable = commonSchema.FromString(nullable)
		cols = append(cols, colVersion)
	}
	log.Infof("%s %s: found %d columns", dbVersion.Db.Name, dbVersion.Version, len(cols))
	return cols
}

func WaitFor(driverName, dsn string, retries int) (db *sql.DB, finalErr error) {
	log.SetLevel(log.DebugLevel)
	ticker := time.NewTicker(time.Second)
	for i := 0; i <= retries; i++ {
		<-ticker.C // wait for a tick

		if db, err := sql.Open(driverName, dsn); err == nil {
			if err := db.Ping(); err == nil {
				log.Infof("connected to %s", dsn)
				return db, nil
			} else {
				finalErr = err
				fmt.Printf(".")
			}
		} else {
			finalErr = err
			log.Debugf("%+v", err)
			fmt.Printf("_")
		}
	}
	return nil, finalErr
}
