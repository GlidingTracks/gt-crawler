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

/*
	In xcportal.pl the format for the URL looks like this:

	https://xcportal.pl/flights-table/2018-08-01

	there seem to be no pagination
 */

func NewXCPortalSite() CrawlableSite {
	return xCPortalSite{
		baseUrl: 			"https://xcportal.pl",
		baseQuery: 			"/flights-table/",
		portalFilePrefix:	"xcportal_",
	}
}

// xContestChrome represents instance of Crawlable.
// Note, the type itself is not exported.
type xCPortalSite struct {
	baseUrl string
	baseQuery string
	portalFilePrefix string
}

// Crawl crawls the daily score records
func (xcc xCPortalSite) Crawl(ctx context.Context, start time.Time, howLong time.Duration) {
	var nodes []*cdp.Node
	url := xcc.baseUrl + xcc.baseQuery + crawler.Timestamp2DateString(start.Unix(),0)

	logrus.Infof("Crawl() in xContestChrome, for URL: %s", url)

	task := xcc.getLinksToDetailPages(url, &nodes)

	err := chromedp.Run(ctx, task)
	if err != nil {
		logrus.Error("initial Run() call error: ", err)
		return
	}

	visitQueue := xcc.filterRealFlightDetailsPages(nodes)
	logrus.Info("visitQueue length: ", len(visitQueue))

	for _, pageNum := range visitQueue {
		logrus.Infof("Visiting details page %s", pageNum)
		err = xcc.visitDetailsPageAndExtract(pageNum, ctx)
		if err != nil {
			logrus.Errorf("Problems processing page: %s", pageNum)
		}
	}
}


func (xcc xCPortalSite) getLinksToDetailPages(url string, nodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks {
		chromedp.Navigate(url),
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(`#main-content`),
		chromedp.WaitReady(`//a`, chromedp.BySearch),
		chromedp.Nodes(`//a`, nodes, chromedp.BySearch),
	}
}

func (xcc xCPortalSite) filterRealFlightDetailsPages(found []*cdp.Node) (visitQueue []string) {
	pseudoHits := make([]string, len(found))

	for i := range found {
		pseudoHits[i] = found[i].AttributeValue("href")
	}

	for k := range pseudoHits {
		if strings.Contains(pseudoHits[k], "node") && !strings.Contains(pseudoHits[k], "add") {
			visitQueue = append(visitQueue, pseudoHits[k])
		}
	}

	return
}

func (xcc xCPortalSite) visitDetailsPageAndExtract(rurl string, ctx context.Context) error {
	logrus.Info("visitDetailsPageAndExtract")

	var nodes []*cdp.Node
	sl := make([]string,0)
	url := xcc.baseUrl + rurl
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
	for _, igcUrl := range sl {
		buf, err := Download(ctx, igcUrl)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof("File downloaded, file size %v", len(buf))
	}
	return nil
}


func (xcc xCPortalSite) getIGCFileLink(url string, nodes *[]*cdp.Node) chromedp.Tasks {
	logrus.Infof("getSourceLinks for %s", url)
	return chromedp.Tasks {
		chromedp.Navigate(url),
		chromedp.WaitReady(`#main-content`, chromedp.BySearch),
		chromedp.Nodes(`div[class="file"] > a`, nodes, chromedp.BySearch),
	}
}
