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
	lw.logger.Infof("Call Crawler -> Scan(%v, %d)", ctx, inc)
	lw.crawler.Scan(ctx, url, inc)
}

func (lw *loggerCrawlerWrap) ChanResult() <-chan CrawlResult {
	lw.logger.Info("Call Crawler -> ChanResult()")
	res := lw.crawler.ChanResult()
	lw.logger.Infof("Result ChanResult: %v", res)
	return res
}

func (lw *loggerCrawlerWrap) AddDepth() {
	lw.logger.Info("Call Crawler -> AddDepth()")
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
	lw.logger.Infof("Call Requester -> Get(%v, %s)", ctx, url)
	p, err := lw.requester.Get(ctx, url)
	lw.logger.Infof("Result Get: %v, %v", p, err)
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
	lw.logger.Info("Call Page -> GetTitle()")
	s := lw.page.GetTitle()
	lw.logger.Infof("Result GetTitle: %s", s)
	return s
}

func (lw *loggerPageWrap) GetLinks() []string {
	lw.logger.Info("Call Page -> GetLinks()")
	s := lw.page.GetLinks()
	lw.logger.Infof("Result GetLinks: %v", s)
	return s
}
