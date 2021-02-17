package chrome

import (
	"context"
	"github.com/sirupsen/logrus"
	"gt-crawler/sites"
)

type Chrome struct{}

func (ch Chrome) Crawl(ctx context.Context, vs []sites.ChromeSite, pipe chan []string) {
	for _, v := range vs {

		links, err := v.Crawl(ctx)
		if err != nil {
			logrus.Error("Error during crawling, err: ", err.Error())
			return
		}

		for _, l := range links {
			logrus.Info(l)
		}

		pipe <- links
	}

	return
}
