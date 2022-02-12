package docs

import (
	"encoding/json"
	"regexp"
	"sort"
	"sync"

	// "log"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cespare/xxhash/v2"
	log "github.com/sirupsen/logrus"

	// "github.com/cespare/xxhash/v2"
	// // _ "github.com/cheggaaa/pb/v3"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/common/schema"
	. "github.com/i-s-compat-table/i.s.compat.table/pkg/common/utils"
	_ "github.com/mattn/go-sqlite3"
)

const basePath = "https://docs.microsoft.com/en-us/sql/"

// "relational-databases/system-information-schema-views"
const familyTreesByMonikerUrl = "https://docs.microsoft.com/_api/familyTrees/bymoniker/sql-server-ver15"
const tocJsonUrl = "https://docs.microsoft.com/en-us/sql/toc.json"

func normalize(s string) string {
	return strings.Trim(strings.ToLower(s), " \t\n\r")
}

type col struct {
	tableName  string
	columnName string
	columnType string
	url        string
	versions   []string
	notes      string
}

func getVersionId(version string) int64 {
	return int64(xxhash.Sum64([]byte("mssql" + version)))
}
func (c *col) Id() int64 {
	sum := xxhash.Digest{}
	sum.WriteString("mssql")
	sum.WriteString(c.Table())
	sum.WriteString(c.Column())
	sum.WriteString(c.Notes())
	return int64(sum.Sum64())
}
func (c *col) Table() string {
	return c.tableName
}
func (c *col) Column() string {
	return c.columnName
}
func (c *col) Type() *string {
	if c.columnType != "" {
		return &c.columnType
	} else {
		log.Panicf("missing column type for %s.%s in %s", c.Table(), c.Column(), c.Url())
		return nil
	}
}

func (c *col) Notes() string {
	return c.notes
}
func (c *col) Url() string {
	return c.url
}

type FamilyMember struct {
	MonikerName        string
	MonikerDisplayName string
	VersionDisplayName string
}
type Family struct {
	ProductName string
	Packages    []FamilyMember
}
type FamilyTree struct {
	FamilyName string
	Products   []Family `json:"products"`
}

type TocTree struct {
	Items []Content
}
type Content struct {
	Href     string
	Monikers []string `json:"monikers"`
	TocTitle string   `json:"toc_title"`
	Children []Content
}

func walk(tree Content, callback func(Content)) {
	callback(tree)
	for _, child := range tree.Children {
		walk(child, callback)
	}
}

func scrapeToc(base *colly.Collector) []Content {
	collector := base.Clone()
	toc := TocTree{}
	result := []Content{}
	collector.OnResponse(func(r *colly.Response) {
		log.Info("...")
		if err := json.Unmarshal(r.Body, &toc); err != nil {
			panic(err)
		}
		for _, item := range toc.Items {
			walk(item, func(item Content) {
				if strings.Contains(item.Href, "relational-databases/system-information-schema-views/") {
					sort.Strings(item.Monikers)
					item.Href = r.Request.AbsoluteURL(item.Href)
					result = append(result, item)
				}
			})
		}
	})
	collector.Visit(tocJsonUrl)
	collector.Wait()
	return result
}

var dedentRe *regexp.Regexp = regexp.MustCompile(`[ \r\n\t]+`)

func scrapePage(
	doc *colly.HTMLElement,
	tableName string,
	versions []string,
) []col {
	url := doc.Request.URL.String()
	if len(versions) == 0 {
		log.Warnf("no versions for %s", url)
	}
	tableName = normalize(tableName)
	// TODO: check passed versions match passed versions in page
	table := doc.DOM.Find("main table").First()
	headers := table.Find("thead > tr > th").
		Map(func(i int, tr *goquery.Selection) string {
			return normalize(tr.Text())
		})
	rows := table.Find("tbody > tr").
		FilterFunction(func(i int, tr *goquery.Selection) bool {
			return tr.Find("td").Length() > 0 && tr.Find("th").Length() == 0
		})
	resultRows := make([]col, rows.Length())
	rows.Each(func(i int, tr *goquery.Selection) {
		resultRow := col{tableName: tableName, url: url, versions: versions}
		columns := tr.Find("td").Map(func(i int, td *goquery.Selection) string {
			return strings.Trim(td.Text(), " \t\n\r")
		})
		for j, col := range columns {
			col = NormalizeString(col)
			switch headers[j] {
			case "column", "field", "column name":
				resultRow.columnName = strings.ToLower(col)
			case "description", "notes":
				resultRow.notes = dedentRe.ReplaceAllString(col, " ")
			case "type", "data type":
				resultRow.columnType = col
			default:
				log.Warnf("unknown header '%s' with value '%s'", headers[j], col)
			}
		}
		resultRows[i] = resultRow
	})
	return resultRows
}

func Scrape(cacheDir string, dbPath string, dbg bool) {
	db := commonSchema.MustConnect(dbPath)
	collector := colly.NewCollector(
		colly.CacheDir(cacheDir),
		colly.Async(true),
		colly.AllowedDomains("docs.microsoft.com"),
	)
	if dbg {
		collector.SetDebugger(&debug.LogDebugger{})
	}
	items := scrapeToc(collector)
	log.Infof("%d items", len(items))
	wg := sync.WaitGroup{}
	wg.Add(1)
	colChan := make(chan []col)
	go func() {
		defer wg.Done()
		txn, err := db.Begin()
		if err != nil {
			panic(err)
		}
		defer txn.Commit()

		insertCol := MustPrepare(txn,
			"INSERT INTO columns(id, db_name, table_name, column_name, column_type, notes) "+
				"VALUES (?, 'mssql', ?, ?, ?, ?) ON CONFLICT DO NOTHING;",
		)
		defer insertCol.Close()

		insertVersion := MustPrepare(txn,
			"INSERT INTO versions(id, db_name, version) "+
				"VALUES (?, 'mssql', ?) ON CONFLICT DO NOTHING;",
		)
		defer insertVersion.Close()

		insertColVersion := MustPrepare(txn,
			"INSERT INTO version_columns(version_id, column_id, column_number) "+
				"VALUES (?, ?, ?) ON CONFLICT DO NOTHING;",
		)
		defer insertColVersion.Close()

		insertUrl := MustPrepare(txn,
			"INSERT INTO urls(id, url) VALUES (?, ?) ON CONFLICT DO NOTHING;")
		defer insertUrl.Close()

		insertUrlRef := MustPrepare(txn,
			"INSERT INTO column_reference_urls(column_id, url_id) VALUES (?, ?) "+
				"ON CONFLICT DO NOTHING;")
		defer insertUrlRef.Close()

		versions := map[int64]bool{}
		urls := map[int64]bool{}

		for {
			if cols, ok := <-colChan; ok {

				for _, c := range cols {
					id := c.Id()
					MustExec(insertCol, id, c.Table(), c.Column(), c.Type(), c.Notes())

					urlId := int64(xxhash.Sum64([]byte(c.Url())))
					if _, ok := urls[urlId]; !ok {
						urls[urlId] = true
						MustExec(insertUrl, urlId, c.url)
					}
					MustExec(insertUrlRef, id, urlId)

					for _, v := range c.versions {
						versionId := getVersionId(v)
						if _, ok := versions[versionId]; !ok {
							versions[versionId] = true
							MustExec(insertVersion, versionId, v)
						}
						MustExec(insertColVersion, versionId, id, nil)
					}
				}
			} else {
				break
			}
		}
	}()
	lookup := func(url string) (viewName string, versions []string) {
		for i := 0; i < len(items); i++ {
			if items[i].Href == url {
				return items[i].TocTitle, items[i].Monikers
			}
		}
		// might be faster than constructing a map and doing a map lookup?
		log.Panic("unable to find "+url, items)
		return "", []string{}
	}
	collector.OnHTML("html", func(html *colly.HTMLElement) {
		url := html.Request.URL
		tableName, versions := lookup(url.Scheme + "://" + url.Host + url.Path)
		colChan <- scrapePage(html, tableName, versions)
	})
	for _, item := range items {
		sort.Strings(item.Monikers)
		viewName := item.TocTitle
		if viewName != "System information schema views" {
			collector.Visit(item.Href)
		}
	}
	collector.Wait()
	close(colChan)
	wg.Wait()
}
