package observer

import (
	"database/sql"
	_ "embed"
	"fmt"
	"strings"
	"time"

	log "log/slog"

	commonSchema "github.com/i-s-compat-table/i.s.compat.table/internal/schema"
	"github.com/i-s-compat-table/i.s.compat.table/internal/utils"
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
	log.Info("found columns", "db", dbVersion.Db.Name, "version", dbVersion.Version, "nColumns", len(cols))
	return cols
}

func WaitFor(driverName, dsn string, retries int) (db *sql.DB, finalErr error) {

	ticker := time.NewTicker(time.Second)
	for i := 0; i <= retries; i++ {
		<-ticker.C // wait for a tick

		if db, err := sql.Open(driverName, dsn); err == nil {
			if err := db.Ping(); err == nil {
				log.Info("connected", "dsn", dsn)
				return db, nil
			} else {
				finalErr = err
				fmt.Printf(".")
			}
		} else {
			finalErr = err
			log.Debug(err.Error())
		}
	}
	return nil, finalErr
}
