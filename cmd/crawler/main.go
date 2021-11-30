package main

import (
	"context"
	"flag"
	"lesson1/internal/crawler"
	"lesson1/internal/requester"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// Config - структура для конфигурации
type Config struct {
	MaxDepth   int64
	MaxResults int
	MaxErrors  int
	Url        string
	Timeout    int // in seconds
}

func main() {

	debug := flag.Bool("debug", false, "debug mode" )
	t := flag.Duration("t", time.Minute, "maximum program execution time")
	url := flag.String("url", "https://telegram.com", "url to parser")
	d := flag.Int64("d", 5, "maximum entry depth")
	e := flag.Int("e", 5, "maximum number of errors")
	r := flag.Int("r", 15, "maximum number of results")
	flag.Parse()

	// TODO: проверять правильность урла...

	cfg := Config{
		MaxDepth:   *d,
		MaxResults: *r,
		MaxErrors:  *e,
		Url:        *url,
		Timeout:    10,
	}

	var cr crawler.Crawler
	var req requester.Requester
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	logger := l.Sugar()
	defer logger.Sync()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(*t))

	if *debug {
		req = requester.NewRequester(time.Duration(cfg.Timeout) * time.Second, logger, *debug)
		rl := requester.NewloggerRequesterWrap(req, logger)
		cr = crawler.NewCrawler(rl, cfg.MaxDepth)
		crl := crawler.NewloggerCrawlerWrap(cr, logger)
		go crl.Scan(ctx, cfg.Url, 1) // запускаем краулер в отдельной рутине
		go processResult(ctx, cancel, crl, cfg, logger) // обрабатываем результаты в отдельной рутине
	} else {
		req = requester.NewRequester(time.Duration(cfg.Timeout) * time.Second, logger, *debug)
		cr = crawler.NewCrawler(req, cfg.MaxDepth)
		go cr.Scan(ctx, cfg.Url, 1) // запускаем краулер в отдельной рутине
		go processResult(ctx, cancel, cr, cfg, logger) // обрабатываем результаты в отдельной рутине
	}

	sigCh := make(chan os.Signal)        // создаем канал для приема сигналов
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGUSR1) // подписываемся на сигнал SIGINT
	for {
		select {
		case <-ctx.Done(): // если всё завершили - выходим
			return
		case s := <-sigCh:
			if s == syscall.SIGINT {
				logger.Info("Stop the program")
				cancel() // если пришёл сигнал SigInt - завершаем контекст
			}
			if s == syscall.SIGUSR1 {
				logger.Info("Increase the depth by two")
				cr.AddDepth()
			}

		}
	}
}

func processResult(ctx context.Context, cancel func(), cr crawler.Crawler, cfg Config, logger *zap.SugaredLogger) {
	var maxResult, maxErrors = cfg.MaxResults, cfg.MaxErrors
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-cr.ChanResult():
			if msg.Err != nil {
				maxErrors--
				logger.Warnf("crawler result return err: %s\n", msg.Err.Error())
				if maxErrors <= 0 {
					cancel()
					return
				}
			} else {
				maxResult--
				logger.Infof("crawler result: [url: %s] Title: %s\n", msg.Url, msg.Title)
				if maxResult <= 0 {
					cancel()
					return
				}
			}
		}
	}
}
