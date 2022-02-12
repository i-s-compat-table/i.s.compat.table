package docs

import (
	"sort"
	"sync"

	// "log"

	// "net/url"
	// "os"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/PuerkitoBio/goquery"
	"github.com/cespare/xxhash/v2"

	// _ "github.com/cheggaaa/pb/v3"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	log "github.com/sirupsen/logrus"

	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/common/schema"
	. "github.com/i-s-compat-table/i.s.compat.table/pkg/common/utils"
	_ "github.com/mattn/go-sqlite3"
)

const basePath = "https://mariadb.com/kb/en"
const infoSchemaTablesIndex = basePath + "/information-schema-tables/"
const tableArticleLink = "li.article > a[href]"
const tableHeader = "tr > th"

var kbLink = regexp.MustCompile("/kb/en/information-schema-([a-z_]+)-table/")
var oddLink = regexp.MustCompile("/kb/en/([a-z-_]*)-table-information-schema/")
var semverPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
var dedentRe = regexp.MustCompile(`[ \r\n\t]+`)

type col struct {
	tableName      string
	columnName     string
	columnType     string
	columnNullable string
	index          int8
	url            string
	versions       []*semver.Version
	notes          string
}

// TODO: impl commonSchema.Column for col

func (c *col) Id() int64 {
	if c.tableName == "" || c.columnName == "" {
		panic("insufficient details to derive id")
	}
	digest := xxhash.Digest{}
	digest.WriteString("postgres")
	digest.WriteString(c.tableName)
	digest.WriteString(c.columnName)
	// digest.WriteString(c.columnType)
	// digest.Write([]byte{byte(c.nullable)})
	digest.WriteString(c.notes)
	return int64(digest.Sum64())
}

func getVersions(base *colly.Collector) semver.Collection {
	rawVersions := make([]string, 0)

	collector := base.Clone() // ensure no callbacks attached
	collector.OnHTML("html", func(html *colly.HTMLElement) {
		rawVersions = html.DOM.Find("tr > td:first-child > a[href]").
			FilterFunction(func(i int, s *goquery.Selection) bool {
				return semverPattern.Match([]byte(s.Text()))
			}).
			Map(func(i int, s *goquery.Selection) string {
				return s.Text()
			})
	})
	collector.Visit("https://mariadb.org/mariadb/all-releases/")
	collector.Wait()
	if len(rawVersions) == 0 {
		panic("no versions")
	}
	versions := make([]*semver.Version, len(rawVersions))
	for i, r := range rawVersions {
		version, err := semver.NewVersion(r)
		if err != nil {
			log.Errorf("Error parsing version: %s", err)
		}
		versions[i] = version
	}
	collection := semver.Collection(versions)
	sort.Sort(collection)
	return collection
}

// get a list of urls, one per information_schema table
func scrapeIndex(collector *colly.Collector) []string {
	collector = collector.Clone()
	links := make([]string, 0)
	collector.OnHTML("html", func(html *colly.HTMLElement) {
		links = html.DOM.Find(tableArticleLink).FilterFunction(
			func(i int, selection *goquery.Selection) bool {
				return !strings.Contains(
					strings.ToLower(selection.Text()),
					"columnstore",
				) && len(selection.First().AttrOr("href", "")) > 0
			},
		).Map(func(i int, selection *goquery.Selection) string {
			href, exists := selection.First().Attr("href")
			if !exists {
				panic("filter didn't catch non-existent href")
			}
			return html.Request.AbsoluteURL(href)
		})
	})
	collector.Visit(infoSchemaTablesIndex)
	collector.Wait()
	return links
}

func scrapePage(page *colly.HTMLElement, tableName string, versions []*semver.Version) []col {
	doc := page.DOM
	pageUrl := page.Request.URL.String()
	if dl := doc.Find("#sidebar-first div.node_info dl"); dl.Length() == 1 {
		txt := dl.Find("dt + dd").Text()
		if !strings.Contains(txt, "CC BY-SA") {
			log.Warn(pageUrl, "is missing cc-by-sa -license")
		}
	} else {
		log.Warnf("dl missing: %s", pageUrl)
	}
	tableName = strings.ReplaceAll(tableName, "-", "_")
	firstVersion := versions[0]
	caveats := doc.Find("div.answer > div.product:first-child")
	caveats.Each(
		func(i int, caveat *goquery.Selection) {
			ver := semver.MustParse(caveat.Find(".product_title > a[href]").First().Text())
			t := strings.ToLower(caveat.First().Text())
			if strings.Contains(t, "starting with") ||
				strings.Contains(t, "introduced") ||
				strings.Contains(t, " added ") {
				log.Debugf("%d\t:\t%s\t:\t%s\t:\t%s", i, ver, t, pageUrl)
				firstVersion = ver
			} else {
				log.Warnf("%d\t:\t%s\t:\t%s\t:\t%s", i, ver, t, pageUrl)
			}
		},
	)
	caveats = doc.Find("div.product")
	if caveats.Length() > 1 {
		caveats.Each(func(i int, caveat *goquery.Selection) {
			if i > 1 {
				log.Warnf("extra caveat %d\t:\t%s\n%s", i, caveat.Text(), pageUrl)
			}
		})
	}
	table := doc.Find("table").First()
	headers := table.Find(tableHeader).Map(func(i int, tr *goquery.Selection) string {
		return strings.Trim(strings.ToLower(tr.Text()), " \t\n\r")
	})

	rows := table.Find("tr").FilterFunction(func(i int, tr *goquery.Selection) bool {
		return tr.Find("td").Length() > 0 && tr.Find("th").Length() == 0
	})
	resultRows := make([]col, rows.Length())
	rows.Each(func(i int, tr *goquery.Selection) {
		resultRows[i].tableName = tableName
		resultRows[i].url = pageUrl
		columns := tr.Find("td").Map(func(i int, td *goquery.Selection) string {
			return strings.Trim(td.Text(), " \t\n\r")
		})
		for j, col := range columns {
			resultRows[i].index = int8(j)
			col = strings.Trim(col, " \t\n\r")
			switch headers[j] {
			case "column", "field", "column name":
				resultRows[i].columnName = strings.ToLower(col)
			case "description", "notes":
				resultRows[i].notes = NormalizeString(col)
			case "type":
				resultRows[i].columnType = col
			case "null":
				resultRows[i].columnNullable = col
			case "added", "introduced":
				if col == "" {
					break
				}
				matches := semverPattern.FindStringSubmatch(col)
				if len(matches) == 0 {
					break
				} else if len(matches) != 1 {
					log.Errorf("%v", matches)
				} else {
					col = matches[0]
				}
				ver, err := semver.NewVersion(col)
				if err == nil {
					log.Infof("version added: %s", ver)
					firstVersion = ver
				}
			default:
				log.Warnf("unknown header '%s' with value '%s'", headers[j], col)
			}
		}
		for _, ver := range versions {
			if ver.GreaterThan(firstVersion) || ver.Equal(firstVersion) {
				resultRows[i].versions = append(resultRows[i].versions, ver)
			}
		}
	})
	return resultRows
}

func Scrape(cacheDir string, dbPath string, dbg bool) {
	collector := colly.NewCollector(
		colly.CacheDir(cacheDir),
		colly.Async(true),
		colly.AllowedDomains("mariadb.com", "mariadb.org"),
		// colly.URLFilters(
		// 	regexp.MustCompile("^https://mariadb.com/kb/en/")
		// ),
	)
	if dbg {
		collector.SetDebugger(&debug.LogDebugger{})
	}
	db, err := commonSchema.Connect(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	links := scrapeIndex(collector)
	versions := getVersions(collector)
	log.Infof("%d links", len(links))
	log.Infof("%d versions", len(versions))
	colChan := make(chan []col)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		txn, err := db.Begin()
		if err != nil {
			panic(err)
		}
		defer func() {
			txn.Commit()
		}()

		insertCol := MustPrepare(txn,
			"INSERT INTO columns(id, db_name, table_name, column_name, column_type, notes) "+
				"VALUES (?, 'mariadb', ?, ?, ?, ?) ON CONFLICT DO NOTHING;",
		)
		defer insertCol.Close()

		insertVersion := MustPrepare(txn,
			"INSERT INTO versions(id, db_name, version) "+
				"VALUES (?, 'mariadb', ?) ON CONFLICT DO NOTHING;",
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

		versions := map[string]int16{}
		for {
			if cols, ok := <-colChan; ok {
				for _, col := range cols {
					colId := col.Id()
					urlId := int64(xxhash.Sum64([]byte(col.url)))
					var columnType *string = nil
					if col.columnType != "" {
						columnType = &col.columnType
					}

					MustExec(insertCol, colId, col.tableName, col.columnName, columnType, col.notes)
					MustExec(insertUrl, urlId, col.url)
					MustExec(insertUrlRef, colId, urlId)
					for _, v := range col.versions {
						version := v.String()
						versionId := int64(xxhash.Sum64([]byte("postgres" + version)))
						if _, ok := versions[version]; ok {
							MustExec(insertVersion, versionId, version)
						}
						versions[version] += 1
						MustExec(insertColVersion, versionId, colId, nil)
					}
				}
			} else {
				break
			}
		}
	}()
	collector.OnHTML("html", func(h *colly.HTMLElement) {
		link := h.Request.URL.String()
		tableMatch := kbLink.FindStringSubmatch(link)
		if len(tableMatch) == 0 {
			tableMatch = oddLink.FindStringSubmatch(link)
		}
		if len(tableMatch) > 1 {
			colChan <- scrapePage(h, tableMatch[1], versions)
		} else {
			log.Warnf("unable to parse table name from %s : %s\n", link, tableMatch)
		}
	})
	for _, link := range links {
		collector.Visit(link)
	}
	collector.Wait()
	close(colChan)
	wg.Wait()
}
