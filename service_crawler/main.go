package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
	"gt-crawler/crawler"
	"gt-crawler/sites"
	"log"
	"sync"
	"time"
)

func main() {
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	var wg sync.WaitGroup


	cSites := []sites.CrawlableSite{
		/*sites.NewXContestSite()*/
		sites.NewXCPortalSite(),
	}
	for _, site := range cSites {
		go crawl(ctx, &wg, site)
		// upload(ctx, conf, links)
		// downloadLinks(ctx, links)
	}
	time.Sleep(2*time.Second)
	wg.Wait()
}

func crawl(ctx context.Context, wg *sync.WaitGroup, site sites.CrawlableSite) {
	defer wg.Done()

	wg.Add(1)
	start, err := crawler.ParseDate("2018-08-01")
	if err != nil {
		logrus.Error(err)
	}
	site.Crawl(ctx, start, time.Duration(1) * 24 * time.Hour)
}