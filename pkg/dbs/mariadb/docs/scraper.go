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

	// _ "github.com/cheggaaa/pb/v3"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	log "github.com/sirupsen/logrus"

	commonSchema "github.com/i-s-compat-table/i.s.compat.table/pkg/common/schema"
	"github.com/i-s-compat-table/i.s.compat.table/pkg/common/utils"
	_ "github.com/mattn/go-sqlite3"
)

var mariaDb = commonSchema.Database{Name: "mariadb"}
var license commonSchema.License = commonSchema.License{
	License:     "CC BY-SA 3.0 AND GFDL-1.3",
	Attribution: "© 2022 MariaDB",
	Url: &commonSchema.Url{
		Url: "https://mariadb.com/kb/en/information-schema-tables/+license/",
	},
}

const ( // used only for scraping
	basePath              = "https://mariadb.com/kb/en"
	infoSchemaTablesIndex = basePath + "/information-schema-tables/"
	tableArticleLink      = "li.article > a[href]"
	tableHeader           = "tr > th"
)

var kbLink = regexp.MustCompile("/kb/en/information-schema-([a-z_]+)-table/")
var oddLink = regexp.MustCompile("/kb/en/([a-z-_]*)-table-information-schema/")
var semverPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
var dedentRe = regexp.MustCompile(`[ \r\n\t]+`)

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

func scrapePage(
	page *colly.HTMLElement,
	tableName string,
	versions []*semver.Version,
) []commonSchema.ColVersion {
	doc := page.DOM
	pageUrl := page.Request.URL.String()
	if dl := doc.Find("#sidebar-first div.node_info dl"); dl.Length() == 1 {
		txt := dl.Find("dt + dd").Text()
		if !strings.Contains(txt, "CC BY-SA") {
			log.Warnf("%s is missing cc-by-sa -license", pageUrl)
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
	data := doc.Find("table").First()
	headers := data.Find(tableHeader).Map(func(i int, tr *goquery.Selection) string {
		return strings.Trim(strings.ToLower(tr.Text()), " \t\n\r")
	})

	rows := data.Find("tr").FilterFunction(func(i int, tr *goquery.Selection) bool {
		return tr.Find("td").Length() > 0 && tr.Find("th").Length() == 0
	})
	if rows.Length() == 0 {
		log.Warnf("%s : no rows", pageUrl)
	}
	resultRows := []commonSchema.ColVersion{}

	table := &commonSchema.Table{Name: tableName}
	baseResult := commonSchema.ColVersion{
		Url: &commonSchema.Url{Url: pageUrl},
	}
	rows.Each(func(i int, tr *goquery.Selection) {
		thisCol := baseResult.Clone()
		columns := tr.Find("td").Map(func(i int, td *goquery.Selection) string {
			return strings.Trim(td.Text(), " \t\n\r")
		})
		for j, col := range columns {
			thisCol.Number = j
			col = utils.NormalizeString(col)
			switch headers[j] {
			case "column", "field", "column name":
				thisCol.Column = &commonSchema.Column{
					Table: table, Name: strings.ToLower(col),
				}
			case "description", "notes":
				thisCol.Notes = &commonSchema.Note{Note: col, License: &license}
			case "type":
				thisCol.Type = &commonSchema.Type{Name: strings.ToUpper(col)}
			case "null":
				switch strings.Trim(strings.ToLower(col), " \t\r\n") {
				case "yes":
					thisCol.Nullable = commonSchema.Nullable
				case "no":
					thisCol.Nullable = commonSchema.NotNullable
				default:
					thisCol.Nullable = commonSchema.Unknown
				}
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

		for i, ver := range versions {
			if ver.GreaterThan(firstVersion) || ver.Equal(firstVersion) {
				c := thisCol.Clone()
				isCurrent := i == len(versions)-1
				c.DbVersion = &commonSchema.Version{
					Db:        &mariaDb,
					Version:   ver.String(),
					IsCurrent: &isCurrent,
				}
				resultRows = append(resultRows, c)
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
	)
	if dbg {
		collector.SetDebugger(&debug.LogDebugger{})
	}

	links := scrapeIndex(collector)
	versions := getVersions(collector)
	log.Infof("%d links", len(links))
	log.Infof("%d versions", len(versions))
	colChan := make(chan []commonSchema.ColVersion)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go commonSchema.BulkInsert(dbPath, colChan, &wg)
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
