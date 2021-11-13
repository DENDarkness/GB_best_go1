package main

import (
	"context"
	"time"
)

type delayCrawlerWrap struct {
	crawler Crawler
	delay time.Duration
}

func NewDelayCrawlerWrap(delay time.Duration, crawler Crawler) Crawler {
	return delayCrawlerWrap{
		crawler: crawler,
		delay: delay,
	}
}

func (d delayCrawlerWrap)  Scan(ctx context.Context, url string, inc int64) {
	time.Sleep(10 * time.Second)
	d.crawler.Scan(ctx, url, inc)
}

func (d delayCrawlerWrap) ChanResult() <-chan CrawlResult {
	return d.crawler.ChanResult()
}

func (d delayCrawlerWrap) AddDepth() {
	d.crawler.AddDepth()
}
