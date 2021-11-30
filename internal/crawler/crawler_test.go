package crawler

import (
	"context"
	"fmt"
	"testing"
)

type mockCrawler struct {
	maxDepth int64
	res      chan CrawlResult
}

func (mc *mockCrawler) Scan(ctx context.Context, url string, inc int64) {
	if mc.maxDepth <= inc {
		return
	}

	if inc > 4 {
		panic(fmt.Sprintf("search depth exceeded: %d", inc))
	}
}

func (mc *mockCrawler) ChanResult() <-chan CrawlResult {
	return mc.res
}
func (mc *mockCrawler) AddDepth() {}

func Test_Scan(t *testing.T) {


}