package main

import (
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/cheggaaa/pb/v3"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/schema"
	"github.com/i-s-compat-table/i.s.compat.table/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var postgres = &commonSchema.Database{Name: "postgres"}
var license = &commonSchema.License{
	License:     "PostgreSQL",
	Attribution: "© 1996-2022 The PostgreSQL Global Development Group",
	Url:         &commonSchema.Url{Url: "https://www.postgresql.org/about/licence/"},
}
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

var dedentRe *regexp.Regexp = regexp.MustCompile(`[ \r\n\t]+`)

func scrape12Minus(html *colly.HTMLElement, tableName string, version string) []commonSchema.ColVersion {
	tableEl := html.DOM.Find("table.table, table.CALSTABLE").First()
	ths := tableEl.Find("thead").First().Find("th")
	headers := make([]string, ths.Length())
	ths.Each(func(i int, th *goquery.Selection) {
		headers[i] = utils.NormalizeString(th.Text())
	})
	trs := tableEl.Find("tbody tr")
	if trs.Length() == 0 {
		log.Warnf("no rows: %s\n", html.Request.URL)
		return nil
	}
	isCurrent := false
	order := commonSchema.AsOrder(version)
	dbVersion := &commonSchema.Version{
		Db: postgres, Version: version, IsCurrent: &isCurrent, Order: &order,
	}
	url := &commonSchema.Url{Url: html.Request.URL.String()}
	table := &commonSchema.Table{Name: tableName}
	cols := make([]commonSchema.ColVersion, trs.Length())
	trs.Each(func(i int, tr *goquery.Selection) {
		cols[i].DbVersion = dbVersion
		cols[i].Url = url
		cols[i].Nullable = commonSchema.Unknown
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			text := utils.NormalizeString(td.Text())
			switch headers[j] {
			case "Name":
				cols[i].Column = &commonSchema.Column{
					Table: table,
					Name:  strings.ToLower(utils.NormalizeString(text)),
				}
			case "Data Type":
				cols[i].Type = &commonSchema.Type{
					Name: strings.ToUpper(utils.NormalizeString(text)),
				}
			case "Description":
				cols[i].Notes = &commonSchema.Note{
					Note:    sectionRef.ReplaceAllString(dedentRe.ReplaceAllString(text, " "), ""),
					License: license,
				}
			default:
				panic("unknown header: " + headers[j])
			}
		})
	})
	return cols
}

func scrape13Plus(page *colly.HTMLElement, tableName string, version string) []commonSchema.ColVersion {
	isCurrent := version == allPgVersions[0]
	order := commonSchema.AsOrder(version)
	dbVersion := &commonSchema.Version{Db: postgres, Version: version, IsCurrent: &isCurrent, Order: &order}
	table := &commonSchema.Table{Name: tableName}
	rows := page.DOM.Find("table.table, table.CALSTABLE").First().Find("table tbody tr")
	url := &commonSchema.Url{Url: page.Request.URL.String()}
	if isCurrent {
		url.Url = strings.Replace(url.Url, version, "current", 1)
	}
	if rows.Length() == 0 {
		log.Warnf("no rows: %s\n", url)
	}
	cols := make([]commonSchema.ColVersion, rows.Length())
	rows.Each(func(i int, tr *goquery.Selection) {
		row := tr.Find(".column_definition").First()
		colNames := row.Find(".structfield")
		colName := strings.ToLower(utils.NormalizeString(colNames.First().Text()))
		columnType := strings.ToUpper(
			utils.NormalizeString(tr.Find(".type").First().Text()))
		notes := utils.NormalizeString(row.NextAll().Filter("p").Text())
		notes = dedentRe.ReplaceAllString(notes, " ")
		notes = sectionRef.ReplaceAllString(notes, "")
		cols[i] = commonSchema.ColVersion{
			DbVersion: dbVersion,
			Number:    i,
			Column:    &commonSchema.Column{Table: table, Name: colName},
			Type:      &commonSchema.Type{Name: columnType},
			Nullable:  commonSchema.Unknown,
			Notes:     &commonSchema.Note{Note: notes, License: license},
			Url:       url,
		}
	})
	return cols
}
func deriveVersion(el *colly.HTMLElement) (version string) {
	version = strings.Split(el.Request.URL.Path, "/")[2]
	//  /docs/13/etc
	// ^0 ^1  ^2
	if version == "current" {
		version = allPgVersions[0]
	}
	return version
}

var sectionRef = regexp.MustCompile(`\s*See Section [0-9.]+ for more information\s*[.]?`)

func Scrape(cacheDir string, dbPath string, dbg bool) {

	collector := colly.NewCollector(
		colly.Async(true),
		colly.CacheDir(cacheDir),
		colly.AllowedDomains("www.postgresql.org"),
		colly.URLFilters(
			regexp.MustCompile(`^https://www.postgresql.org/docs/[0-9.]+/infoschema-.+$`),
			regexp.MustCompile(`^https://www.postgresql.org/docs/[0-9.]+/information-schema.html$`),
		),
		colly.DisallowedURLFilters(
			regexp.MustCompile(`infoschema-datatypes.html$`),
			regexp.MustCompile(`infoschema-schema.html$`),
		),
	)
	if dbg {
		collector.SetDebugger(&debug.LogDebugger{})
	}
	nRoutines := runtime.NumCPU()*2 - 1
	collector.Limit(&colly.LimitRule{Parallelism: nRoutines - runtime.NumCPU()})
	log.Infof("scraping versions %s - %s", allPgVersions[0], allPgVersions[len(allPgVersions)-1])
	wg := sync.WaitGroup{}

	urlChan := make(chan string)
	urls := []string{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if url, ok := <-urlChan; ok {
				urls = append(urls, url)
			} else {
				break
			}
		}
	}()
	collector.OnHTML("div.toc a[href], div.TOC a[href]", func(a *colly.HTMLElement) {
		href := a.Request.AbsoluteURL(a.Attr("href"))
		urlChan <- href
	})
	for _, version := range allPgVersions {
		err := collector.Visit("https://www.postgresql.org/docs/" + version + "/information-schema.html")
		if err != nil {
			panic(err)
		}
	}
	collector.Wait()
	close(urlChan)
	wg.Wait()

	sort.Strings(urls)
	collector = collector.Clone() // clear callbacks

	rowChan := make(chan []commonSchema.ColVersion)
	wg.Add(1)
	go commonSchema.BulkInsert(dbPath, rowChan, &wg)
	nRows := 0

	re := regexp.MustCompile(`[a-zA-Z_]+$`)
	collector.OnHTML("html", func(html *colly.HTMLElement) {
		version := deriveVersion(html)
		title := utils.NormalizeString(html.DOM.Find("title").Text()) // 13+ sometimes have 0-width or nonbreaking spaces thrown in
		tableName := strings.Trim(strings.ToLower(re.FindString(title)), " ")
		if v, err := strconv.ParseFloat(version, 32); err == nil {
			if v >= 13 {
				rows := scrape13Plus(html, tableName, version)
				nRows += len(rows)
				rowChan <- rows
			} else {
				rows := scrape12Minus(html, tableName, version)
				nRows += len(rows)
				rowChan <- rows
			}
		} else {
			panic(err)
		}
	})
	for _, url := range urls {
		collector.Visit(url)
	}
	collector.Wait()
	log.Infof("scraped %d column-versions", nRows)
	close(rowChan)
	wg.Wait()
}

func main() {
	Scrape("./.cache", "./data/postgres/docs.sqlite", false)
}
