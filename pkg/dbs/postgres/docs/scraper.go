package docs

import (
	"database/sql"
	"fmt"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/cespare/xxhash/v2"
	_ "github.com/cheggaaa/pb/v3"
	"github.com/gocolly/colly/v2"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/common/schema"
	"github.com/i-s-compat-table/i.s.compat.table/pkg/common/utils"
	_ "github.com/mattn/go-sqlite3"
)

var allPgVersions = [...]string{
	"14",
	"13",
	"12",
	"11",
	"10",
	"9.6",
	"9.5",
	"9.4",
	"9.3",
	"9.2",
	"9.1",
	"9.0",
	"8.4",
	"8.3",
	"8.2",
	"8.1",
	"8.0",
	"7.4",
}

type col struct {
	table      string
	column     string
	notes      string
	columnType string
	url        string
	version    string

	index    uint8
	nullable commonSchema.Nullability
}

func (c *col) Id() int64 {
	if c.table == "" || c.column == "" {
		panic("insufficient details to derive id")
	}
	digest := xxhash.Digest{}
	digest.WriteString("postgres")
	digest.WriteString(c.table)
	digest.WriteString(c.column)
	digest.WriteString(c.columnType)
	// digest.Write([]byte{byte(c.nullable)})
	digest.WriteString(c.notes) // nitpicky: "can" -> "may" creates a new id

	return int64(digest.Sum64())
}

func (c *col) Table() string {
	return c.table
}

func (c *col) Column() string {
	return c.column
}

func (c *col) Type() *string {
	return &c.columnType
}

func (c *col) Nullable() *bool {
	if c.nullable == commonSchema.Unknown {
		return nil
	} else {
		nullable := (c.nullable == commonSchema.Nullable)
		return &nullable
	}
}

func (c *col) Notes() string {
	return c.notes
}

// TODO
func (c *col) Versions() []string {
	return []string{}
}
func (c *col) Url() string {
	return c.url
}

var dedentRe *regexp.Regexp = regexp.MustCompile(`[ \r\n\t]+`)

func scrape12Minus(html *colly.HTMLElement, tableName string, version string) []col {
	tableEl := html.DOM.Find("table.table, table.CALSTABLE").First()
	ths := tableEl.Find("thead").First().Find("th")
	headers := make([]string, ths.Length())
	ths.Each(func(i int, th *goquery.Selection) {
		headers[i] = utils.NormalizeString(th.Text())
	})
	trs := tableEl.Find("tbody tr")
	if trs.Length() == 0 {
		fmt.Printf("no rows: %s\n", html.Request.URL)
		return nil
	}
	cols := make([]col, trs.Length())
	trs.Each(func(i int, tr *goquery.Selection) {
		cols[i].table = tableName
		cols[i].url = html.Request.URL.String()
		cols[i].version = version
		cols[i].nullable = commonSchema.Unknown
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			text := utils.NormalizeString(td.Text())
			switch headers[j] {
			case "Name":
				cols[i].column = strings.ToLower(text)
			case "Data Type":
				cols[i].columnType = strings.ToLower(text)
			case "Description":
				cols[i].notes = dedentRe.ReplaceAllString(text, " ")
			default:
				panic("unknown header: " + headers[j])
			}
		})
	})
	return cols
}

func scrape13Plus(page *colly.HTMLElement, tableName string, version string) []col {
	rows := page.DOM.Find("table.table, table.CALSTABLE").First().Find("table tbody tr")
	if rows.Length() == 0 {
		fmt.Printf("no rows: %s\n", page.Request.URL)
	}
	cols := make([]col, rows.Length())
	rows.Each(func(i int, tr *goquery.Selection) {
		row := tr.Find(".column_definition").First()
		colNames := row.Find(".structfield")
		colName := strings.ToLower(utils.NormalizeString(colNames.First().Text()))
		columnType := tr.Find(".type").First().Text()
		notes := utils.NormalizeString(row.NextAll().Filter("p").Text())
		notes = dedentRe.ReplaceAllString(notes, " ")
		cols[i] = col{
			table:      tableName,
			column:     colName,
			index:      uint8(i),
			columnType: columnType,
			nullable:   commonSchema.Unknown,
			notes:      notes,
			version:    version,
			url:        page.Request.URL.String(),
		}
	})
	return cols
}
func deriveVersion(el *colly.HTMLElement) string {
	return strings.Split(el.Request.URL.Path, "/")[2]
	//  /docs/13/etc
	// ^0 ^1  ^2
}

func BulkInsert(db *sql.DB, inputs <-chan []col) {
	txn, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer txn.Commit()

	insertCol, err := txn.Prepare(
		"INSERT INTO columns(id, db_name, table_name, column_name, column_type, column_nullable, notes) " +
			"VALUES (?, 'postgres', ?, ?, ?, ?, ?) ON CONFLICT DO NOTHING;",
	)
	if err != nil {
		panic(err)
	}
	defer insertCol.Close()

	insertVersion, err := txn.Prepare(
		"INSERT INTO versions(id, db_name, version) " +
			"VALUES (?, 'postgres', ?) ON CONFLICT DO NOTHING;",
	)
	if err != nil {
		panic(err)
	}
	defer insertVersion.Close()

	insertColVersion, err := txn.Prepare(
		"INSERT INTO version_columns(version_id, column_id, column_number) " +
			"VALUES (?, ?, ?) ON CONFLICT DO NOTHING;",
	)
	if err != nil {
		panic(err)
	}
	defer insertCol.Close()

	versions := make(map[string]int16, len(allPgVersions))
	for {
		if cols, ok := <-inputs; ok {
			for _, col := range cols {
				id := col.Id()
				// nullable := col.nullable == commonSchema.Nullable
				if col.nullable == commonSchema.Unknown {

					_, err = insertCol.Exec(
						id, col.Table(), col.Column(), col.Type(), nil, col.Notes())
					if err != nil {
						fmt.Println(versions)
						fmt.Println(id, col.Table(), col.Column(), col.Type(), nil, col.Notes())
						panic(err)
					}
				} else {
					nullable := col.nullable == commonSchema.Nullable
					_, err = insertCol.Exec(
						id, col.Table(), col.Column(), col.Type(), &nullable, col.Notes())
					if err != nil {
						panic(err)
					}
				}
				// xxh3_64 is roughly as fast as memcpy, so hashing every row's version id
				// is affordable
				versionId := int64(xxhash.Sum64([]byte("postgres" + col.version)))
				versions[col.version]++
				if versions[col.version] <= 1 {
					_, err = insertVersion.Exec(versionId, col.version)
					if err != nil {
						panic(err)
					}
				}
				_, err := insertColVersion.Exec(versionId, id, col.index)
				if err != nil {
					panic(err)
				}
			}
		} else {
			fmt.Println(versions)
			break
		}
	}
}

// TODO: also pass sqlite3 db location
func Scrape(cacheDir string, dbPath string) {
	db, err := commonSchema.Connect(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// TODO: if the tsv exists, write it out?

	collector := colly.NewCollector(
		colly.Async(true),
		colly.CacheDir(cacheDir),
		colly.AllowedDomains("www.postgresql.org"),
		colly.URLFilters(
			regexp.MustCompile(`^https://www.postgresql.org/docs/[0-9.]+/infoschema-.+$`),
			regexp.MustCompile(`^https://www.postgresql.org/docs/[0-9.]+/information-schema.html$`),
		),
	)
	nRoutines := runtime.NumCPU()*2 - 1
	collector.Limit(&colly.LimitRule{Parallelism: nRoutines - runtime.NumCPU()})
	for _, version := range allPgVersions {
		err := collector.Visit("https://www.postgresql.org/docs/" + version + "/information-schema.html")
		if err != nil {
			panic(err)
		}
	}
	urlChan := make(chan string)
	urls := []string{}

	go func() {
		for {
			if url, ok := <-urlChan; ok {
				urls = append(urls, url)
			} else {
				break
			}
		}
	}()

	collector.OnHTML("div.toc a[href], div.TOC a[href]", func(a *colly.HTMLElement) {
		href := a.Attr("href")
		if strings.HasSuffix(href, ".html") {
			urlChan <- a.Request.AbsoluteURL(href)
		}
	})
	collector.Wait()
	close(urlChan)

	sort.Strings(urls)
	collector = collector.Clone() // clear callbacks

	rowChan := make(chan []col)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		BulkInsert(db, rowChan)
		wg.Done()
	}()

	re := regexp.MustCompile(`[a-zA-Z_]+$`)
	collector.OnHTML("html", func(html *colly.HTMLElement) {
		version := deriveVersion(html)
		title := re.FindString(html.DOM.Find("title").Text())
		tableName := strings.Trim(strings.ToLower(title), " ")
		if v, err := strconv.ParseFloat(version, 32); err == nil {
			if v >= 13 {
				rowChan <- scrape13Plus(html, tableName, version)
			} else {
				rowChan <- scrape12Minus(html, tableName, version)
			}
		}
	})
	for _, page := range urls {
		collector.Visit(page)
	}
	collector.Wait()
	close(rowChan)
	wg.Wait()
}
