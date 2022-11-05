package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/schema"
	"github.com/i-s-compat-table/i.s.compat.table/pkg/utils"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var tidb = &commonSchema.Database{Name: "tidb"}
var license = &commonSchema.License{
	License:     "Apache-2.0",
	Attribution: "© 2022 PingCAP",
	Url:         &commonSchema.Url{Url: "https://github.com/pingcap/tidb/blob/master/LICENSE"},
}
var versions = [...]string{
	"6.1",
	"6.0",
	"5.4",
	"5.3",
	"5.2",
	"5.1",
	"5.0",
	"4.0",
	// "3.1",
	// "3.0",
	// "2.1",
}

const (
	baseRawUrl  = "https://raw.githubusercontent.com/pingcap/docs/"
	baseNiceUrl = "https://docs.pingcap.com/tidb/"
)

func normalize(s string) string {
	return strings.ToLower(strings.Trim(s, " \t\n\r"))
}

func isCurrent(version string) bool {
	return version == versions[0]
}
func getNiceUrl(version string, table string) string {
	result := strings.Builder{}
	result.WriteString(baseNiceUrl)
	if isCurrent(version) {
		result.WriteString("stable")
	} else {
		result.WriteRune('v')
		result.WriteString(version)
	}
	result.WriteString("/information-schema-")
	result.WriteString(strings.ReplaceAll(table, "_", "-"))
	return result.String()
}
func getRawUrl(branch string, file string) string {
	return baseRawUrl + branch + "/information-schema/" + file
}

func getVersionBranch(version string) string {
	if isCurrent(version) {
		return "master"
	}
	return "release-" + version
}

func deriveVersion(url string) (version string) {
	m := strings.Split(url, "/")[5]
	if len(m) > 8 && m[:8] == "release-" {
		version = m[8:]
	} else if m == "master" {
		version = versions[0]
	} else {
		log.Panic(url, m)
	}
	return
}

type plusTableDesc struct {
	Field, Type, Null, Key, Default, Extra string
}

func mdToHtml(md goldmark.Markdown, body []byte) *goquery.Document {
	buf := bytes.Buffer{}
	if err := md.Convert(body, &buf); err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		panic(err)
	}
	return doc
}

func collectPagesToScrape(collector *colly.Collector, md goldmark.Markdown, wg *sync.WaitGroup) (pages []string) {
	pageChan := make(chan string)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if page, ok := <-pageChan; ok {
				pages = append(pages, page)
			} else {
				break
			}
		}
	}()
	collector.OnResponse(func(r *colly.Response) {
		doc := mdToHtml(md, r.Body)
		doc.Find("a[href]").Each(func(i int, a *goquery.Selection) {
			url := r.Request.AbsoluteURL(
				strings.Replace(r.Request.URL.Path, "/information-schema/information-schema.md", a.AttrOr("href", ""), 1),
			)
			pageChan <- url
		})
	})
	for _, version := range versions {
		rawUrl := getRawUrl(getVersionBranch(version), "information-schema.md")
		collector.Visit(rawUrl)
	}
	collector.Wait()
	close(pageChan)
	wg.Wait()
	sort.Strings(pages)
	return
}

// field descriptions are stored in list. This function parses them into a mapping
// of field name to field description
func parseFieldList(list *goquery.Selection) map[string]string {
	lis := list.Children().Filter("li")
	definitions := make(map[string]string, lis.Length())
	if lis.Length() == 0 {
		return definitions
	}
	lis.Each(func(i int, li *goquery.Selection) {
		text := li.Text()
		split := strings.SplitN(text, ":", 2)
		var key, definition string
		if len(split) == 2 {
			key = normalize(split[0])
			definition = strings.Trim(split[1], " \t\n\r")
		} else {
			code := li.Contents().First().Filter("code").Text()
			key = strings.Trim(normalize(code), "\": ")
			if len(key) == 0 {
				log.Warnf("no key: %s", text)
			}
			definition = split[0][len(key):]
			definition = strings.Trim(definition, "\": ")
		}
		if prev, ok := definitions[key]; ok {
			log.Warnf("previous definition for %s: %s", key, prev)
		}
		definitions[key] = definition
	})
	return definitions
}

func isDefinitionList(list *goquery.Selection) (ok bool) {
	lis := list.Children().Filter("li")
	ok = lis.Length() > 0
	if !ok {
		log.Debug("no lis")
	}
	lis.EachWithBreak(func(i int, li *goquery.Selection) bool {
		text := li.Text()
		ok = li.Contents().First().Filter("code").Length() != 0
		if !ok {
			log.Debug("no code", text)
		}
		ok = ok || len(strings.Split(text, ":")) >= 2
		if !ok {
			log.Debug("no split", text)
		}
		return ok
	})
	return ok
}

type void struct{}

func findColDescriptions(doc *goquery.Document) []map[string]string {
	expected := map[*goquery.Selection]void{}
	handleList := func(p *goquery.Selection) {
		nextList := p.NextAllFiltered("ul, ol").First()
		if isDefinitionList(nextList) {
			if _, ok := expected[nextList]; ok {
				log.Warnf("duplicate: %+v", nextList)
			} else {
				expected[nextList] = void{}
			}
		}
	}
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		normalized := normalize(s.Text())
		if strings.Contains(normalized, "field description") {
			handleList(s)
			return
		} else {
			for _, x := range []string{"field", "column", "description", "describe", "columns"} {
				if strings.Contains(normalized, x) {
					handleList(s)
					return
				} else {
					log.Debug("rejected:", normalized)
				}
			}
		}
	})

	results := make([]map[string]string, 0, len(expected))
	for s := range expected {
		results = append(results, parseFieldList(s))
	}

	return results
}

// parse a plus table into an array of rows, each of which is an array of cells
func parsePlusTable(text string) [][]string {
	lines := strings.Split(text, "\n")

	result := make([][]string, 0, len(lines)) // slightly over-allocate the backing array

	for _, line := range lines {
		row := strings.Split(line, "|")
		if len(row) < 3 {
			// it's a delimiter like `+----+----+`
			continue
		}
		for i := 0; i < len(row); i++ {
			row[i] = strings.Trim(row[i], "\t ")
		}
		result = append(result, row[1:len(row)-1])
	}
	return result
}

// parse the result of `DESC tbl` calls.
func parseDescriptions(doc *goquery.Document, tableName string, url string) map[string][]*plusTableDesc {
	matches := doc.Find("pre>code.language-sql").FilterFunction(func(i int, s *goquery.Selection) bool {
		normalized := strings.ToLower(s.Text())
		return strings.Contains(normalized, "use information_schema;") && strings.Contains(normalized, "desc ")
	})
	allCols := make(map[string][]*plusTableDesc, matches.Length())
	matches.Each(func(i int, s *goquery.Selection) {
		describedTableName := strings.Split(strings.Split(strings.ToLower(s.Text()), ";")[1], "desc ")[1]
		if tableName != describedTableName && fmt.Sprintf("cluster_%s", tableName) != describedTableName {
			log.Panicf("title %s != table %s", tableName, describedTableName)
		}
		text := s.Parent().Next().Text()
		plusTable := parsePlusTable(text)
		if len(plusTable) == 0 {
			log.Panicf("unable to parse +table @ %s", url)
		}
		if header := strings.Join(plusTable[0], ","); header != "Field,Type,Null,Key,Default,Extra" {
			log.Panicf("unexpected header @ %s : %s", url, header)
		}
		cols := make([]*plusTableDesc, len(plusTable[1:]))
		for i, row := range plusTable[1:] {
			cols[i] = &plusTableDesc{
				Field:   strings.ToLower(utils.NormalizeString(row[0])),
				Type:    strings.ToUpper(utils.NormalizeString(row[1])),
				Null:    strings.ToUpper(utils.NormalizeString(row[2])),
				Key:     row[3],
				Default: row[4],
				Extra:   row[5],
			}
		}
		allCols[describedTableName] = cols
	})
	return allCols
}

func scrapePage(doc *goquery.Document, tableName string, url string, rowChan chan []commonSchema.ColVersion) {
	versionNumber := deriveVersion(url)
	order := commonSchema.AsOrder(versionNumber)
	isCurrent := versionNumber == versions[0]
	version := commonSchema.Version{
		Db: tidb, IsCurrent: &isCurrent, Version: versionNumber, Order: &order,
	}
	DESCs := parseDescriptions(doc, tableName, url)
	lists := findColDescriptions(doc)
	if len(lists) == 0 || len(DESCs) == 0 || len(DESCs) > len(lists) {
		log.Warnf("DESCs\t%03d\tdefns\t%03d\t%s\t%s\t", len(DESCs), len(lists), tableName, url)
	} else {
		log.Debugf("DESCs\t%03d\tdefns\t%03d\t%s\t%s\t", len(DESCs), len(lists), tableName, url)
	}
	for tableName, desc := range DESCs {
		commonUrl := commonSchema.Url{Url: getNiceUrl(versionNumber, tableName)}
		table := commonSchema.Table{Name: tableName}
		resultSet := make([]commonSchema.ColVersion, len(desc))
		for i, colDesc := range desc {
			col := commonSchema.Column{
				Table: &table,
				Name:  colDesc.Field,
			}
			nullable := commonSchema.Unknown
			if colDesc.Null != "" {
				normalized := strings.ToUpper(utils.NormalizeString(colDesc.Null))
				if normalized == "YES" {
					nullable = commonSchema.Nullable
				} else if normalized == "NO" {
					nullable = commonSchema.NotNullable
				} else {
					log.Panic(colDesc.Null)
				}
			}
			colVersion := commonSchema.ColVersion{
				DbVersion: &version,
				Column:    &col,
				Type:      &commonSchema.Type{Name: strings.ToUpper(utils.NormalizeString(colDesc.Type))},
				Url:       &commonUrl,
				Nullable:  nullable,
				Number:    i,
			}
			for _, definitions := range lists {
				if notes, ok := definitions[colDesc.Field]; ok {
					colVersion.Notes = &commonSchema.Note{Note: notes, License: license}
					break
				}
			}
			resultSet[i] = colVersion
		}
		rowChan <- resultSet
	}
}

func scrape(cacheDir string, dbPath string, dbg bool) {
	collector := colly.NewCollector(
		colly.CacheDir(cacheDir),
		colly.AllowedDomains("raw.githubusercontent.com"),
		colly.Async(true),
	)
	if dbg {
		collector.SetDebugger(&debug.LogDebugger{})
	}
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	wg := sync.WaitGroup{}
	pages := collectPagesToScrape(collector.Clone(), md, &wg)
	bar := pb.New(len(pages))
	rowChan := make(chan []commonSchema.ColVersion)

	wg.Add(1)
	go commonSchema.BulkInsert(dbPath, rowChan, &wg)

	collector.OnResponse(func(r *colly.Response) {
		defer bar.Increment()
		url := r.Request.URL
		parts := strings.Split(url.Path, "/")
		if len(parts) <= 5 {
			return
		}
		prefix, rest := parts[5][:19], parts[5][19:]
		if prefix != "information-schema-" {
			return
		}
		tableName := strings.ReplaceAll(strings.SplitN(rest, ".", 2)[0], "-", "_")
		doc := mdToHtml(md, r.Body)
		scrapePage(doc, tableName, url.String(), rowChan)
	})
	for _, page := range pages {
		collector.Visit(page)
	}
	collector.Wait()
	bar.Finish()
	close(rowChan)
	wg.Wait()

}

func main() {
	scrape("./.cache", "./data/tidb/docs.sqlite", false)
}
