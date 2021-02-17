package sites

import (
	"context"
	"github.com/MarkusAJacobsen/jConfig-go"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
	"gt-crawler/crawlTime"
	"sort"
	"strings"
	"time"
)


// xContestChrome represents instance of Crawlable
type xContestChrome struct {
	baseUrl string
	filterDateQuery string
	pagination string
	crawledDates []string
}

func NewXContest() ChromeSite {
	return xContestChrome {
		baseUrl: 			"https://www.xcontest.org/world/en/flights/daily-score-pg/",
		filterDateQuery:  	"#filter[date]=",
		pagination: 		"@flights[start]=",
	}
}

// CrawledDates returns the already crawled dates for this ChromeSites
func (xcc xContestChrome) CrawledDates() []string {
	return xcc.crawledDates
}

// Crawl crawls the daily score records from xcontest.org
func (xcc xContestChrome) Crawl(ctx context.Context) (sl []string, err error) {
	var nodes []*cdp.Node
	url, date, _ := xcc.getURL(true)
	task := getLinksFromUrl(url, &nodes)

	logrus.Info("Crawl() in xContestChrome")

	err = chromedp.Run(ctx, task)
	if err != nil {
		logrus.Error("intial Run() call", err)
		return
	}

	visitQueue := filterRealFlightLinks(nodes)
	logrus.Info("visitQueue length: ", len(visitQueue))

	sl, err = visitDetailsPagesAndExtract(visitQueue, ctx)
	if err == nil || len(sl) == 0 {
		return
	}

	// Ensure crawled date has been written to config before launching a new run
	dateUpdated := xcc.writeCrawledDateToConfig(date)

	if !dateUpdated {
		logrus.Info("Could not update last crawled date, aborting")
		return
	}

	// If still finding links, recursively continue to next date
	csl, err := xcc.Crawl(ctx)
	if err != nil {
		logrus.Error("Error during crawling, err: ", err)
		return nil, err
	}
	logrus.Infof("XContest date: %v crawled", date)
	sl = append(sl, csl...)

	return
}

func (xcc xContestChrome) writeCrawledDateToConfig(date string) (updated bool) {
	conf := &jConfigGo.Config{}
	err := conf.CreateConfig("xcontest")
	if err != nil {
		logrus.Error(err)
		return false
	}

	if err := conf.Get(&xcc); err != nil {
		logrus.Errorf("Could not write date: %v to config", date)
		return false
	}

	xcc.crawledDates = append(xcc.crawledDates, date)
	if err := conf.Write(xcc); err != nil {
		logrus.Errorf("Could not write date: %v to config", date)
		return false
	}

	return true
}

func visitDetailsPagesAndExtract(urls []string, ctx context.Context) (sl []string, err error) {
	logrus.Info("visitDetailsPagesAndExtract")

	for i := range urls {
		var nodes []*cdp.Node

		task := getSourceLink(urls[i], &nodes)
		err := chromedp.Run(ctx, task)
		if err != nil {
			logrus.Error(err)
		}

		for i := range nodes {
			if strings.Contains(nodes[i].AttributeValue("href"), ".igc") {
				sl = append(sl, nodes[i].AttributeValue("href"))
			}
		}
	}

	return
}

func getLinksFromUrl(url string, nodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(`#flights`),
		chromedp.WaitReady(`//a`, chromedp.BySearch),
		chromedp.Nodes(`//a`, nodes, chromedp.BySearch),
	}
}

func getSourceLink(url string, nodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady(`#flight`, chromedp.BySearch),
		chromedp.Nodes(`th[class="igc"] > a`, nodes, chromedp.BySearch),
	}
}

func filterRealFlightLinks(found []*cdp.Node) (visitQueue []string) {
	pseudoHits := make([]string, len(found))

	for i := range found {
		pseudoHits[i] = found[i].AttributeValue("href")
	}

	for k := range pseudoHits {
		if strings.Contains(pseudoHits[k], "flights/detail") {
			visitQueue = append(visitQueue, pseudoHits[k])
		}
	}

	return
}

func (xcc xContestChrome) getURL(pagination bool) (url string, date string, err error) {
	conf := &jConfigGo.Config{}
	err = conf.CreateConfig("xcontest")
	if err != nil {
		logrus.Error(err)
		return
	}

	if err = conf.Get(&xcc); err != nil {
		return
	}

	if len(xcc.crawledDates) == 0 {
		date = crawlTime.GetDateString(0)
		url = xcc.baseUrl + xcc.filterDateQuery + date
		return
	}

	sort.Strings(xcc.crawledDates)

	// Earliest date crawled
	date = crawlTime.FindPreviousDate(xcc.crawledDates[0])
	url = xcc.baseUrl + xcc.filterDateQuery + date

	return
}

