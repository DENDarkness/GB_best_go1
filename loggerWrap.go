package main

import (
	"context"

	"go.uber.org/zap"
)

type loggerCrawlerWrap struct {
	crawler Crawler
	logger *zap.SugaredLogger
}

func NewLoggerCrawlerWrap(crawler Crawler, logger  *zap.SugaredLogger) *loggerCrawlerWrap {
	return &loggerCrawlerWrap{
		crawler: crawler,
		logger: logger,
	}
}

func (lw *loggerCrawlerWrap)  Scan(ctx context.Context, url string, inc int64) {
	lw.logger.Debugf("Call Crawler -> Scan(%v, %d)", ctx, inc)
	lw.crawler.Scan(ctx, url, inc)
}

func (lw *loggerCrawlerWrap) ChanResult() <-chan CrawlResult {
	lw.logger.Debugf("Call Crawler -> ChanResult()")
	res := lw.crawler.ChanResult()
	lw.logger.Debugf("Result ChanResult: %v", res)
	return res
}

func (lw *loggerCrawlerWrap) AddDepth() {
	lw.logger.Debugf("Call Crawler -> AddDepth()")
	lw.crawler.AddDepth()
}

type loggerRequesterWrap struct {
	requester Requester
	logger *zap.SugaredLogger
}

func NewLoggerRequesterWrap(r Requester, logger *zap.SugaredLogger) *loggerRequesterWrap {
	return &loggerRequesterWrap{
		requester: r,
		logger: logger,
	}
}

func (lw *loggerRequesterWrap) Get(ctx context.Context, url string) (Page, error) {
	lw.logger.Debugf("Call Requester -> Get(%v, %s)", ctx, url)
	p, err := lw.requester.Get(ctx, url)
	lw.logger.Debugf("Result Get: %v, %v", p, err)
	return p, err
}

type loggerPageWrap struct {
	page Page
	logger *zap.SugaredLogger
}

func NewLoggerPageWrap(page Page, logger *zap.SugaredLogger) *loggerPageWrap {
	return &loggerPageWrap{
		page: page,
		logger: logger,
	}
}


func (lw *loggerPageWrap) GetTitle() string {
	lw.logger.Debugf("Call Page -> GetTitle()")
	s := lw.page.GetTitle()
	lw.logger.Debugf("Result GetTitle: %s", s)
	return s
}

func (lw *loggerPageWrap) GetLinks() []string {
	lw.logger.Debugf("Call Page -> GetLinks()")
	s := lw.page.GetLinks()
	lw.logger.Debugf("Result GetLinks: %v", s)
	return s
}
