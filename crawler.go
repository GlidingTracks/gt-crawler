package main

import (
	"context"
	"github.com/MarkusAJacobsen/jConfig-go"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
	"gt-crawler/auth"
	"gt-crawler/chrome"
	"gt-crawler/sites"
	"sync"
)

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	conf, err := getConfig()
	if err != nil {
		logrus.Fatal("Could not get config, storing result locally")
		// CHECK STORAGE AND UPLOAD RESIDUALS
	}

	var wg sync.WaitGroup

	links := crawl(ctx, &wg)
	upload(ctx, conf, links)

	wg.Wait()
}

func crawl(ctx context.Context, wg *sync.WaitGroup) (links []string) {
	defer wg.Done()

	c := &chrome.Chrome{}
	cSites := []sites.ChromeSite{sites.NewXContest()}
	crawlRes := make(chan []string)

	wg.Add(1)
	go c.Crawl(ctx, cSites, crawlRes)
	links = <-crawlRes
	close(crawlRes)

	return
}

func upload(ctx context.Context, conf State, links []string) {
	if len(links) == 0 {
		logrus.Info("No links uploaded, empty input")
		return
	}

	up := &Upload{
		Auth: auth.FAuth{},
	}

	uploaded := up.UploadLinks(ctx, links, conf)
	if uploaded {
		logrus.Info("Links uploaded")
	}
}

func getConfig() (state State, err error) {
	conf := jConfigGo.Config{}

	if err = conf.CreateConfig("state"); err != nil {
		return
	}

	state = State{}
	err = conf.Get(&state)

	return
}
