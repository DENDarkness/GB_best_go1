package main

import (
	"context"
	"log"
)

type loggerCrawlerWrap struct {
	crawler Crawler
}

func NewLoggerCrawlerWrap(crawler Crawler) *loggerCrawlerWrap {
	return &loggerCrawlerWrap{
		crawler: crawler,
	}
}

func (lw *loggerCrawlerWrap)  Scan(ctx context.Context, url string, inc int64) {
	log.Println("interface: Crawler; method: Scan; input params: ctx, url, inc; output: not")
	lw.crawler.Scan(ctx, url, inc)
}

func (lw *loggerCrawlerWrap) ChanResult() <-chan CrawlResult {
	log.Println("interface: Crawler; method: ChanResult; input params: not; output: chan type CrawlResult")
	return lw.crawler.ChanResult()
}

func (lw *loggerCrawlerWrap) AddDepth() {
	log.Println("interface: Crawler; method: AddDepth; input params: not; output: not")
	lw.crawler.AddDepth()
}

type loggerRequesterWrap struct {
	requester Requester
}

func NewLoggerRequesterWrap(r Requester) *loggerRequesterWrap {
	return &loggerRequesterWrap{
		requester: r,
	}
}

func (lw *loggerRequesterWrap) Get(ctx context.Context, url string) (Page, error) {
	log.Println("interface: Requester; method: Get; input params: ctx, url; output: Page interface, error")
	return lw.requester.Get(ctx, url)
}

type loggerPageWrap struct {
	page Page
}

func NewLoggerPageWrap(page Page) *loggerPageWrap {
	return &loggerPageWrap{
		page: page,
	}
}


func (lw *loggerPageWrap) GetTitle() string {
	log.Println("interface: Page; method: GetTitle; input params: not; output: string")
	return lw.page.GetTitle()
}

func (lw *loggerPageWrap) GetLinks() []string {
	log.Println("interface: Page; method: GetLinks; input params: not; output: slice string")
	return lw.page.GetLinks()
}
