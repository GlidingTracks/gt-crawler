package sites

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"os"
	"time"
)

// CrawlableSite represents a single crawlable site.
// The site persistent state is backed up by DB.
// The site knows how to crawl a single resource URL,
// usually setup by date, through Crawl method and
// how to keep the persistent state in the DB.
type CrawlableSite interface {

	// Crawl crawls this particular crawlable site, starting
	// with a given date for the given duration.
	Crawl(ctx context.Context, from time.Time, howLong time.Duration)

	getLinksToDetailPages(url string, nodes *[]*cdp.Node) chromedp.Tasks

	filterRealFlightDetailsPages(found []*cdp.Node) (visitQueue []string)

	visitDetailsPageAndExtract(url string, ctx context.Context) error

	getIGCFileLink(url string, nodes *[]*cdp.Node) chromedp.Tasks
}

// Download utility to download a specific resource from a given URL.
func Download(ctx context.Context, url string) ([]byte, error) {
	done := make(chan bool)
	tab, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var requestId network.RequestID
	chromedp.ListenTarget(tab, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
			req := ev.Request
			if req.URL == url {
				requestId = ev.RequestID
			}

		case *network.EventLoadingFinished:
			if ev.RequestID == requestId {
				close(done)
			}
		}
	})

	err := chromedp.Run(tab, chromedp.Tasks {
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorAllow).WithDownloadPath(os.TempDir()),
		chromedp.Navigate(url),
	})
	if err != nil {
		return nil, err
	}

	<- done

	var bytes []byte
	err = chromedp.Run(tab, chromedp.ActionFunc(func(cxt context.Context) error {
		bytes, err = network.GetResponseBody(requestId).Do(cxt)
		return err
	}))

	return bytes, nil
}