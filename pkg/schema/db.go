package schema

import (
	"database/sql"
	_ "embed"
	"encoding/binary"
	"os"
	"strings"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/i-s-compat-table/i.s.compat.table/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

//go:embed db.sql
var Schema string

type Nullability int8

const (
	Unknown     Nullability = -1
	Nullable    Nullability = 0
	NotNullable Nullability = 1
)

func (n Nullability) ToBool() *bool {
	if n == Unknown {
		return nil
	}
	nullable := (n == Nullable)
	return &nullable
}
func FromBool(nullable *bool) Nullability {
	if nullable == nil {
		return Unknown
	} else if *nullable {
		return Nullable
	} else {
		return NotNullable
	}
}
func FromString(nullable string) Nullability {
	s := strings.ToLower(utils.NormalizeString(nullable))
	switch s {
	case "yes":
		return Nullable
	case "no":
		return NotNullable
	case "":
		return Unknown
	}
	log.Panicf("unknown value: '%s'", s)
	return Unknown
}

func xxhash3_64(text ...string) int64 {
	d := xxhash.Digest{}
	for _, str := range text {
		d.WriteString(str)
	}
	return int64(d.Sum64())
}

func digestInt64(digest *xxhash.Digest, i int64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	digest.Write(b)
}

func DeriveDbId(dbName string) int64   { return xxhash3_64(dbName) }
func DeriveTableId(table string) int64 { return xxhash3_64(table) }
func DeriveUrlId(url string) int64     { return xxhash3_64(url) }
func DeriveNoteId(note string) int64   { return xxhash3_64(note) }
func DeriveTypeId(name string) int64   { return xxhash3_64(name) }
func DeriveVersionId(dbId int64, version string) int64 {
	d := xxhash.Digest{}
	digestInt64(&d, dbId)
	d.WriteString(version)
	return int64(d.Sum64())
}
func DeriveColumnId(tableId int64, colName string) int64 {
	d := xxhash.Digest{}
	digestInt64(&d, tableId)
	d.WriteString(colName)
	return int64(d.Sum64())
}
func DeriveLicenseId(license string, attribution string) int64 {
	return xxhash3_64(license, attribution)
}
func DeriveColumnVersionId(colId int64, versionId int64) int64 {
	d := xxhash.Digest{}
	digestInt64(&d, colId)
	digestInt64(&d, versionId)
	return int64(d.Sum64())
}

const (
	InsertDbQ    = "INSERT INTO dbs(id, name) VALUES (?, ?) ON CONFLICT DO NOTHING;"
	InsertTableQ = "INSERT INTO tables(id, name) VALUES (?, ?) ON CONFLICT DO NOTHING;"
	InsertTypeQ  = "INSERT INTO types(id, name) VALUES (?, ?) ON CONFLICT DO NOTHING;"
	InsertUrlQ   = "INSERT INTO urls(id, url) VALUES (?, ?) ON CONFLICT DO NOTHING;"
	InsertNoteQ  = "INSERT INTO notes(id, note) VALUES (?, ?) ON CONFLICT DO NOTHING;"

	InsertColumnQ = "INSERT INTO columns(id, table_id, name) VALUES (?, ?, ?) " +
		"ON CONFLICT DO NOTHING;"

	InsertVersionQ = "INSERT INTO versions(id, db_id, version, is_current) " +
		"VALUES (?, ?, ?, ?) ON CONFLICT DO UPDATE " +
		"SET is_current = coalesce(excluded.is_current, is_current);"
	InsertLicenseQ = "INSERT INTO licenses(id, license, attribution, link_id) " +
		"VALUES (?, ?, ?, ?) ON CONFLICT DO NOTHING;" // ?
	InsertColumnVersionQ = "INSERT INTO column_versions(" +
		"id, column_id, version_id, type_id, nullable, url_id, note_id, note_license_id" +
		") VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET" +
		" type_id = coalesce(excluded.type_id, type_id)" +
		" , url_id = coalesce(excluded.url_id, url_id)" +
		" , note_id = coalesce(excluded.note_id, note_id)" +
		" , note_license_id = coalesce(excluded.note_license_id, note_license_id);"
)

type Database struct {
	Name string
}

func (db Database) Id() int64 {
	return DeriveDbId(db.Name)
}

type Version struct {
	Db        *Database
	IsCurrent *bool
	Version   string
}

func (v Version) Id() int64 {
	return DeriveVersionId(v.Db.Id(), v.Version)
}

type Table struct{ Name string }

func (t Table) Id() int64 {
	return DeriveTableId(t.Name)
}

type Column struct {
	Table *Table
	Name  string
}

func (c Column) Id() int64 {
	return DeriveColumnId(c.Table.Id(), c.Name)
}

type Type struct {
	Name string
}

func (t Type) Id() int64 {
	return DeriveTypeId(t.Name)
}

type License struct {
	License     string
	Attribution string
	Url         *Url
}

func (l License) Id() int64 {
	return DeriveLicenseId(l.License, l.Attribution)
}

type Note struct {
	Note    string
	License *License
}

func (n Note) Id() int64 {
	return DeriveNoteId(n.Note)
}

type Url struct {
	Url string
}

func (u Url) Id() int64 {
	return DeriveUrlId(u.Url)
}

type ColVersion struct {
	DbVersion *Version
	Column    *Column
	Type      *Type
	Notes     *Note
	Url       *Url
	Nullable  Nullability
	Number    int
}

func (c ColVersion) Id() int64 {
	return DeriveColumnVersionId(c.Column.Id(), c.DbVersion.Id())
}

func (c ColVersion) Clone() ColVersion {
	return ColVersion{
		DbVersion: c.DbVersion,
		Column:    c.Column,
		Type:      c.Type,
		Notes:     c.Notes,
		Url:       c.Url,
		Nullable:  c.Nullable,
		Number:    c.Number,
	}
}

func BulkInsert(outputPath string, cols <-chan []ColVersion, wg *sync.WaitGroup) {
	db := MustConnect(outputPath)
	txn, err := db.Begin()
	if err != nil {
		panic(err)
	}
	count := 0

	defer func() {
		if err := txn.Commit(); err != nil {
			panic(err)
		}
		if err := db.Close(); err != nil {
			panic(err)
		}
		wg.Done()
		log.Infof("done; inserted %d column-versions into %s", count, outputPath)
	}()

	insertDb := utils.MustPrepare(txn, InsertDbQ)
	insertVersion := utils.MustPrepare(txn, InsertVersionQ)
	insertTable := utils.MustPrepare(txn, InsertTableQ)
	insertColumn := utils.MustPrepare(txn, InsertColumnQ)
	insertType := utils.MustPrepare(txn, InsertTypeQ)
	insertNote := utils.MustPrepare(txn, InsertNoteQ)
	insertLicense := utils.MustPrepare(txn, InsertLicenseQ)
	insertUrl := utils.MustPrepare(txn, InsertUrlQ)
	insertColVersion := utils.MustPrepare(txn, InsertColumnVersionQ)
	for {
		if batch, ok := <-cols; ok {
			for _, col := range batch {
				count++
				dbId := col.DbVersion.Db.Id()

				utils.MustExec(insertDb, dbId, col.DbVersion.Db.Name)

				versionId := DeriveVersionId(dbId, col.DbVersion.Version)
				utils.MustExec(insertVersion,
					versionId, dbId, col.DbVersion.Version, col.DbVersion.IsCurrent)

				tableId := col.Column.Table.Id()
				utils.MustExec(insertTable, tableId, col.Column.Table.Name)

				colId := DeriveColumnId(tableId, col.Column.Name)
				utils.MustExec(insertColumn, colId, tableId, col.Column.Name)

				var typeId *int64
				if col.Type != nil {
					id := col.Type.Id()
					typeId = &id
					utils.MustExec(insertType, typeId, col.Type.Name)
				}

				var urlId *int64 = nil
				if col.Url != nil {
					id := col.Url.Id()
					urlId = &id
					utils.MustExec(insertUrl, id, col.Url.Url)
				}

				var licenseId *int64
				var noteId *int64
				if notes := col.Notes; notes != nil {
					if notes.License != nil {
						id := notes.License.Id()
						licenseId = &id
						var licenseUrlId *int64
						if notes.License.Url != nil {
							id := notes.License.Url.Id()
							licenseUrlId = &id
							utils.MustExec(insertUrl, id, notes.License.Url.Url)
						}
						if notes.Note != "" {
							utils.MustExec(insertLicense, licenseId, notes.License.License, notes.License.Attribution, licenseUrlId)
							id := notes.Id()
							noteId = &id
							utils.MustExec(insertNote, id, notes.Note)
						}
					}
				}
				colVersionId := DeriveColumnVersionId(colId, versionId)

				utils.MustExec(insertColVersion,
					colVersionId,
					colId,
					versionId,
					typeId,
					col.Nullable.ToBool(),
					urlId,
					noteId,
					licenseId,
				)
			}
		} else {
			break
		}
	}
}

func Connect(path string) (db *sql.DB, err error) {
	if f, e := os.Stat(path); e != nil {
		// initialize if the path doesn't exist
		db, err = sql.Open("sqlite3", path)
		if err != nil {
			panic(err)
		}
		for _, cmd := range strings.Split(Schema, ";") {
			_, err = db.Exec(cmd)
			if err != nil {
				panic(err)
			}
		}
		return db, err
	} else if f.IsDir() {
		log.Panicf("%s is a directory", path)
	} else {
		db, err = sql.Open("sqlite3", path)
	}
	return db, err
}

func MustConnect(path string) *sql.DB {
	db, err := Connect(path)
	if err != nil {
		panic(err)
	}
	return db
}
