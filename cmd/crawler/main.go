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

//Config - структура для конфигурации
type Config struct {
	MaxDepth   int64
	MaxResults int
	MaxErrors  int
	Url        string
	Timeout    int //in seconds
}

func main() {

	debug := flag.Bool("debug", false, "Debug mode" )
	t := flag.Duration("t", time.Minute, "Maximum program execution time")

	flag.Parse()

	cfg := Config{
		MaxDepth:   5,
		MaxResults: 15,
		MaxErrors:  5,
		Url:        "https://telegram.com",
		Timeout:    10,
	}
	var cr crawler.Crawler
	var r requester.Requester
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	logger := l.Sugar()
	defer logger.Sync()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(*t))

	if *debug {
		r = requester.NewRequester(time.Duration(cfg.Timeout) * time.Second, logger, *debug)
		rl := requester.NewloggerRequesterWrap(r, logger)
		cr = crawler.NewCrawler(rl, cfg.MaxDepth)
		crl := crawler.NewloggerCrawlerWrap(cr, logger)
		go crl.Scan(ctx, cfg.Url, 1) //Запускаем краулер в отдельной рутине
		go processResult(ctx, cancel, crl, cfg, logger) //Обрабатываем результаты в отдельной рутине
	}

	r = requester.NewRequester(time.Duration(cfg.Timeout) * time.Second, logger, *debug)
	cr = crawler.NewCrawler(r, cfg.MaxDepth)
	go cr.Scan(ctx, cfg.Url, 1) //Запускаем краулер в отдельной рутине
	go processResult(ctx, cancel, cr, cfg, logger) //Обрабатываем результаты в отдельной рутине

	sigCh := make(chan os.Signal)        //Создаем канал для приема сигналов
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGUSR1) //Подписываемся на сигнал SIGINT
	for {
		select {
		case <-ctx.Done(): //Если всё завершили - выходим
			return
		case s := <-sigCh:
			if s == syscall.SIGINT {
				logger.Info("Stop the program")
				cancel() //Если пришёл сигнал SigInt - завершаем контекст
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
