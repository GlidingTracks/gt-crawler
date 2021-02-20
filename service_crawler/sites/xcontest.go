package sites

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
	"gt-crawler/crawler"
	"strings"
	"time"
)

// xContestChrome represents instance of Crawlable
type xContestSite struct {
	baseUrl string
	baseQuery string
	filterDateQuery string
	pagination string
}

/*

	In xcontest.org the format for the URL looks like

	https://www.xcontest.org/2018/world/en/flights/daily-score-pg/#filter[date]=2018-08-01@flights[start]=100

 */

func NewXContestSite() CrawlableSite {
	return xContestSite {
		baseUrl: 			"https://www.xcontest.org/",
		baseQuery: 			"world/en/flights/daily-score-pg/",
		filterDateQuery:  	"#filter[date]=",
		pagination: 		"@flights[start]=",
	}
}


// Crawl crawls the daily score records from xcontest.org
func (xcc xContestSite) Crawl(ctx context.Context, start time.Time, howLong time.Duration) {
	var nodes []*cdp.Node
	url := xcc.baseUrl + "2018/" + xcc.baseQuery + xcc.filterDateQuery + crawler.Timestamp2DateString(start.Unix(),0)

	logrus.Infof("Crawl() in xContestChrome, for URL: %s", url)

	task := xcc.getLinksToDetailPages(url, &nodes)

	err := chromedp.Run(ctx, task)
	if err != nil {
		logrus.Error("intial Run() call", err)
		return
	}

	visitQueue := xcc.filterRealFlightDetailsPages(nodes)
	logrus.Info("visitQueue length: ", len(visitQueue))

	for _, pageNum := range visitQueue {
		err = xcc.visitDetailsPageAndExtract(pageNum, ctx)
		if err != nil {
			logrus.Errorf("Problems processing page: %s", pageNum)
		}
	}
}


func (xcc xContestSite) getLinksToDetailPages(url string, nodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks {
		chromedp.Navigate(url),
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(`#flights`),
		chromedp.WaitReady(`//a`, chromedp.BySearch),
		chromedp.Nodes(`//a`, nodes, chromedp.BySearch),
	}
}

func (xcc xContestSite) filterRealFlightDetailsPages(found []*cdp.Node) (visitQueue []string) {
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

func (xcc xContestSite) visitDetailsPageAndExtract(url string, ctx context.Context) error {
	logrus.Info("visitDetailsPageAndExtract")

	var nodes []*cdp.Node
	sl := make([]string,0)
	task := xcc.getIGCFileLink(url, &nodes)

	err := chromedp.Run(ctx, task)
	if err != nil {
		logrus.Error(err)
	}

	for _, node := range nodes {
		if strings.Contains(node.AttributeValue("href"), ".igc") {
			sl = append(sl, node.AttributeValue("href"))
		}
	}

	logrus.Infof("Found %v .igc files to download", len(sl))
	for i, igcUrl := range sl {

		// TODO debugging only, just doing it for 2 files
		if i > 2 { return nil }
		buf, err := Download(ctx, igcUrl)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof("File downloaded, file size %v", len(buf))

	}
	return nil
}


func (xcc xContestSite) getIGCFileLink(url string, nodes *[]*cdp.Node) chromedp.Tasks {
	logrus.Info("getSourceLinks")
	return chromedp.Tasks {
		chromedp.Navigate(url),
		chromedp.WaitReady(`#flight`, chromedp.BySearch),
		chromedp.Nodes(`th[class="igc"] > a`, nodes, chromedp.BySearch),
	}
}




/**

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

	for i, url := range urls {
		if i > 4 { return }
		var nodes []*cdp.Node
		chromedp.New
		task := getIGCFileLink(url, &nodes)

		var bytes []byte
		err = chromedp.Run(ctx, chromedp.ActionFunc(func(cxt context.Context) error {
			bytes, err = network.GetResponseBody(requestId).Do(cxt)
			return err
		}))

		err := chromedp.Run(ctx, task)
		if err != nil {
			logrus.Error(err)
		}

		for _, node := range nodes {
			if strings.Contains(node.AttributeValue("href"), ".igc") {
				sl = append(sl, node.AttributeValue("href"))
			}
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

**/