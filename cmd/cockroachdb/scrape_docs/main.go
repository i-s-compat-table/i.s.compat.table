package main

import (
	"bytes"
	"regexp"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/internal/schema"
	"github.com/i-s-compat-table/i.s.compat.table/internal/utils"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/net/html"
)

var cockroachdb = &commonSchema.Database{Name: "cockroachdb"}
var license = &commonSchema.License{
	License:     "CC BY 4.0",
	Attribution: "© 2022 CockroachDB",
	Url: &commonSchema.Url{
		Url: "https://github.com/cockroachdb/docs/blob/master/LICENSE",
	},
}

var versions = [...]string{ // see https://github.com/cockroachdb/docs
	"21.2",
	"21.1",
	"20.2",
	"20.1",
	"19.2",
	"19.1",
	"2.1",
	"2.0",
	"1.1",
	"1.0",
}

const (
	baseUrl     = "https://raw.githubusercontent.com/cockroachdb/docs/master/"
	gitPath     = "/information-schema.md"
	cosmeticUrl = "https://www.cockroachlabs.com/docs"
)

func getRawUrl(version string) string {
	return baseUrl + "v" + version + gitPath
}
func getNiceUrl(version string) string {
	return cosmeticUrl + "/v" + version + "/information-schema"
}
func getVersion(url string) string {
	offset := len(baseUrl) + 1
	return strings.Split(url[offset:], "/")[0]
}

var viewRe = regexp.MustCompile(`^[a-zA-Z_]+$`)

func scrape(cacheDir string, dbPath string, dbg bool) {
	collector := colly.NewCollector(
		colly.CacheDir("./.cache"),
		colly.AllowedDomains("raw.githubusercontent.com"),
		colly.Async(true),
	)
	if dbg {
		collector.SetDebugger(&debug.LogDebugger{})
	}
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		// goldmark.WithParserOptions(
		// 	parser.WithAutoHeadingID(),
		// ),
	)
	colChan := make(chan []commonSchema.ColVersion)
	collector.OnResponse(func(r *colly.Response) {
		version := getVersion(r.Request.URL.String())
		niceUrl := getNiceUrl(version)
		isCurrent := version == versions[0]
		if isCurrent {
			niceUrl = strings.Replace(niceUrl, "v"+version, "stable", 1)
		}
		buf := bytes.Buffer{}
		if err := md.Convert(r.Body, &buf); err != nil {
			panic(err)
		}
		doc, err := goquery.NewDocumentFromReader(&buf)
		if err != nil {
			panic(err)
		}
		versionOrder := commonSchema.AsOrder(version)
		baseColumnVersion := commonSchema.ColVersion{
			DbVersion: &commonSchema.Version{
				Db: cockroachdb, Version: version, IsCurrent: &isCurrent, Order: &versionOrder,
			},
			Url: &commonSchema.Url{Url: niceUrl},
		}

		doc.Find("h3").Each(func(i int, h3 *goquery.Selection) {
			var title string
			n := h3.Nodes[0].FirstChild
			for n != nil {
				if n.Type == html.TextNode && len(strings.Trim(n.Data, " \t\r\n")) > 0 {
					title = utils.NormalizeString(n.Data)
					break
				}
				n = n.NextSibling
			}
			if !viewRe.Match([]byte(title)) {
				return
			}
			tableBaseColumnVersion := baseColumnVersion.Clone()
			table := &commonSchema.Table{Name: title}
			tableBaseColumnVersion.Url = &commonSchema.Url{
				Url: niceUrl + "#" + title,
			}
			tableEl := h3.NextAllFiltered("table").First()
			headers := tableEl.Find("th").Map(func(i int, s *goquery.Selection) string {
				return utils.NormalizeString(strings.ToLower(s.Text()))
			})
			trs := tableEl.Find("tbody > tr")
			rows := make([]commonSchema.ColVersion, 0, trs.Length())

			trs.Each(func(i int, tr *goquery.Selection) {
				row := tableBaseColumnVersion.Clone()
				row.Nullable = commonSchema.Unknown
				tr.Find("td").Each(func(j int, td *goquery.Selection) {
					text := utils.NormalizeString(td.Text())
					switch headers[j] {
					case "column":
						row.Column = &commonSchema.Column{Name: strings.ToLower(text), Table: table}
					case "description":
						row.Notes = &commonSchema.Note{Note: text, License: license}
						if text == "Always NULL." {
							row.Nullable = commonSchema.Nullable
						}
					default:
						log.Warnf("unknown column %s @ %s", text, r.Request.URL)
					}
				})
				if row.Column == nil {
					log.Panicf("no col %+v @ %s", row, row.Url.Url)
				}
				rows = append(rows, row)
			})
			colChan <- rows
		})
	})

	for _, version := range versions {
		collector.Visit(getRawUrl(version))
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go commonSchema.BulkInsert(dbPath, colChan, &wg)
	collector.Wait()
	close(colChan)
	wg.Wait()
}

func main() {
	scrape("./.cache", "./data/cockroachdb/docs.sqlite", false)
}
