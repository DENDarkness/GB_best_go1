package main

import (
	"context"
	"time"
)

type delayRequesterWrap struct {
	requester Requester
	delay time.Duration
}

func NewDelayRequesterWrap(delay time.Duration, requester Requester) Requester {
	return delayRequesterWrap{
		requester: requester,
		delay: delay,
	}
}


func (d delayRequesterWrap) Get(ctx context.Context, url string) (Page, error) {
	time.Sleep(5 * time.Second)
	return d.requester.Get(ctx, url)
}