package main

import (
	"context"
	"github.com/GlidingTracks/gt-crawler/auth"
	"github.com/GlidingTracks/gt-crawler/chrome"
	"github.com/GlidingTracks/gt-crawler/sites"
	"github.com/MarkusAJacobsen/jConfig-go"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	conf, err := getConfig()
	if err != nil {
		logrus.Fatal("Could not get config, storing result locally")
		// CHECK STORAGE AND UPLOAD RESIDUALS
	}

	links := crawl(ctx)
	upload(ctx, conf, links)
}

func crawl(ctx context.Context) (links []string) {
	c := &chrome.Chrome{}
	cSites := []sites.ChromeSite{&sites.XContestChrome{}}
	crawlRes := make(chan []string)

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
