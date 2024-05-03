package main

import (
	"encoding/json"
	"regexp"
	"sort"
	"sync"

	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"

	// "github.com/cespare/xxhash/v2"
	// // _ "github.com/cheggaaa/pb/v3"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	commonSchema "github.com/i-s-compat-table/i.s.compat.table/internal/schema"
	"github.com/i-s-compat-table/i.s.compat.table/internal/utils"
	_ "github.com/mattn/go-sqlite3"
)

const basePath = "https://docs.microsoft.com/en-us/sql/"

// "relational-databases/system-information-schema-views"
const familyTreesByMonikerUrl = "https://docs.microsoft.com/_api/familyTrees/bymoniker/sql-server-ver15"
const tocJsonUrl = "https://docs.microsoft.com/en-us/sql/toc.json"

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

func scrapeFamilyTree(base *colly.Collector) map[string][]string {
	collector := base.Clone()
	tree := FamilyTree{}
	collector.OnResponse(func(r *colly.Response) {
		if err := json.Unmarshal(r.Body, &tree); err != nil {
			panic(err)
		}
	})
	collector.Visit(familyTreesByMonikerUrl)
	collector.Wait()
	result := map[string][]string{}

	set := map[string]bool{}
	for _, product := range tree.Products {
		packages := []string{}
		for _, pkg := range product.Packages {
			result[pkg.MonikerDisplayName] = []string{pkg.MonikerName}
			packages = append(packages, pkg.MonikerName)
		}
		result[product.ProductName] = packages
		if strings.Split(product.ProductName, " ")[0] == "SQL" {
			for _, pkg := range packages {
				set[pkg] = true
			}
		}
	}
	allSql := make([]string, 0, len(set))
	for pkg := range set {
		allSql = append(allSql, pkg)
	}
	sort.Strings(allSql)
	result["SQL Server (all supported versions)"] = allSql
	return result
}

var dedentRe *regexp.Regexp = regexp.MustCompile(`[ \r\n\t]+`)
var mssql = &commonSchema.Database{Name: "mssql"}
var license = &commonSchema.License{
	License:     "CC BY 4.0",
	Attribution: "© Microsoft 2022",
	Url: &commonSchema.Url{
		Url: "https://github.com/MicrosoftDocs/sql-docs/blob/live/LICENSE",
	},
}

// scrape a page for the mssql db versions the page pertains to
func getVersions(page *colly.HTMLElement, fam map[string][]string) (result []string) {
	imgs := page.DOM.Find("p>img").First()
	if imgs.Length() == 0 {
		panic("unable to find versions on " + page.Request.URL.String())
	}
	for n := imgs.Parent().Nodes[0].FirstChild; n != nil; n = n.NextSibling {
		if n.Type == html.TextNode {
			name := utils.NormalizeString(n.Data)
			if name == "" {
				continue
			}
			if monikers, ok := fam[name]; ok {
				result = append(result, monikers...)
			} else {
				log.Panicf("unknown name %s", name)
			}
		}
	}
	return result
}

func scrapePage(
	doc *colly.HTMLElement,
	tableName string,
	versions []string,
	fam map[string][]string,
) []commonSchema.ColVersion {
	url := &commonSchema.Url{Url: doc.Request.URL.String()}
	if len(versions) == 0 {
		versions = getVersions(doc, fam)
	}
	table := &commonSchema.Table{Name: utils.NormalizeString(tableName)}
	data := doc.DOM.Find("main table").First()
	headers := data.Find("thead > tr > th").
		Map(func(i int, tr *goquery.Selection) string {
			return strings.ToLower(utils.NormalizeString(tr.Text()))
		})
	rows := data.Find("tbody > tr").
		FilterFunction(func(i int, tr *goquery.Selection) bool {
			return tr.Find("td").Length() > 0 && tr.Find("th").Length() == 0
		})
	resultRows := []commonSchema.ColVersion{} // TODO: pre-allocate
	rows.Each(func(i int, tr *goquery.Selection) {

		columns := tr.Find("td").Map(func(i int, td *goquery.Selection) string {
			return utils.NormalizeString(td.Text())
		})
		if len(columns) == 0 {
			log.Panicf("no columns in " + url.Url)
		}
		result := commonSchema.ColVersion{Url: url}
		for j, col := range columns {
			col = utils.NormalizeString(col)
			switch headers[j] {
			case "column", "field", "column name":
				result.Column = &commonSchema.Column{
					Table: table, Name: strings.ToLower(col),
				}
			case "description", "notes":
				result.Notes = &commonSchema.Note{
					Note:    dedentRe.ReplaceAllString(col, " "),
					License: license,
				}
			case "type", "data type":
				result.Type = &commonSchema.Type{Name: strings.ToUpper(col)}
			default:
				log.Warnf("unknown header '%s' with value '%s'", headers[j], col)
			}
			if result.Column == nil {
				log.Panicf("%+v <- %s", result, url.Url)
			}
			allVersions := fam["SQL Server (all supported versions)"]

			for _, version := range versions {
				var versionNumber int64
				for i := 0; i < len(allVersions); i++ {
					if version == allVersions[i] {
						break
					}
				}
				versionNumber = int64(i)

				colVersion := result.Clone()
				colVersion.DbVersion = &commonSchema.Version{Db: mssql, Version: version, Order: &versionNumber}
				resultRows = append(resultRows, colVersion)
			}
		}
	})
	return resultRows
}

func Scrape(cacheDir string, dbPath string, dbg bool) {
	collector := colly.NewCollector(
		colly.CacheDir(cacheDir),
		colly.Async(true),
		colly.AllowedDomains("docs.microsoft.com"),
	)
	if dbg {
		collector.SetDebugger(&debug.LogDebugger{})
	}
	items := scrapeToc(collector)
	fam := scrapeFamilyTree(collector)
	log.Infof("%d items", len(items))
	wg := sync.WaitGroup{}
	wg.Add(1)
	colChan := make(chan []commonSchema.ColVersion)
	go commonSchema.BulkInsert(dbPath, colChan, &wg)
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
		colChan <- scrapePage(html, strings.ToLower(tableName), versions, fam)
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

func main() {
	Scrape("./.cache", "./data/mssql/docs.sqlite", false)
}
