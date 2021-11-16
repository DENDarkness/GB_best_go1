package main

// Code generated by gowrap. DO NOT EDIT.
// template: template/zap
// gowrap: http://github.com/hexdigest/gowrap

//go:generate gowrap gen -p lesson1 -i Requester -t template/zap -o loggerRequesterZapWrap.go -l ""

import (
	"context"

	"go.uber.org/zap"
)

type loggerRequesterWrap struct {
	logger *zap.SugaredLogger
	base   Requester
}

func NewloggerRequesterWrap(base Requester, log *zap.SugaredLogger) loggerRequesterWrap {
	return loggerRequesterWrap{
		base:   base,
		logger: log,
	}
}

func (d loggerRequesterWrap) Get(ctx context.Context, url string) (p1 Page, err error) {
	d.logger.Debug("Call Requester -> Get(ctx context.Context, url string) (p1 Page, err error)")
	res, err := d.base.Get(ctx, url)
	d.logger.Debugf("Result Get: %v, %v", res, err)
	return res, err

}
