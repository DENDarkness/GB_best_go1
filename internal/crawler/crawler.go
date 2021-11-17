package crawler

import (
	"context"
	"fmt"
	"lesson1/internal/requester"
	"sync"
	"sync/atomic"
)

//Crawler - интерфейс (контракт) краулера
type Crawler interface {
	Scan(ctx context.Context, url string, inc int64)
	ChanResult() <-chan CrawlResult
	AddDepth()
}

type CrawlResult struct {
	Err   error
	Title string
	Url   string
}

type crawler struct {
	r        requester.Requester
	maxDepth int64
	res      chan CrawlResult
	visited  map[string]struct{}
	mu       sync.RWMutex
}

func NewCrawler(r requester.Requester, depth int64) *crawler {
	return &crawler{
		r:       r,
		maxDepth: depth,
		res:     make(chan CrawlResult),
		visited: make(map[string]struct{}),
		mu:      sync.RWMutex{},
	}
}

func (c *crawler) Scan(ctx context.Context, url string, inc int64) {

	if c.maxDepth <= inc {
		return
	}

	if inc > 4 {
		panic(fmt.Sprintf("search depth exceeded: %d", inc))
	}

	c.mu.RLock()
	_, ok := c.visited[url] //Проверяем, что мы ещё не смотрели эту страницу
	c.mu.RUnlock()
	if ok {
		return
	}
	select {
	case <-ctx.Done(): //Если контекст завершен - прекращаем выполнение
		return
	default:
		page, err := c.r.Get(ctx, url) //Запрашиваем страницу через Requester
		if err != nil {
			c.res <- CrawlResult{Err: err} //Записываем ошибку в канал
			return
		}
		c.mu.Lock()
		c.visited[url] = struct{}{} //Помечаем страницу просмотренной
		c.mu.Unlock()
		c.res <- CrawlResult{ //Отправляем результаты в канал
			Title: page.GetTitle(),
			Url:   url,
		}

		atomic.AddInt64(&inc, 1)

		for _, link := range page.GetLinks() {
			go c.Scan(ctx, link, inc) //На все полученные ссылки запускаем новую рутину сборки
		}
	}
}

func (c *crawler) ChanResult() <-chan CrawlResult {
	return c.res
}

func (c *crawler) AddDepth() {
	c.maxDepth += 2
}
