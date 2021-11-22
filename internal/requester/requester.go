package requester

import (
	"context"
	"lesson1/internal/page"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Requester interface {
	Get(ctx context.Context, url string) (page.Page, error)
}

type requester struct {
	timeout time.Duration
	logger *zap.SugaredLogger
	debugType bool
}

func NewRequester(timeout time.Duration, logger *zap.SugaredLogger, dt bool) requester {
	return requester{
		timeout: timeout,
		logger: logger,
		debugType: dt,
	}
}

func (r requester) Get(ctx context.Context, url string) (page.Page, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		cl := &http.Client{
			Timeout: r.timeout,
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			r.logger.Errorf("Failed to create NewRequest: %v", err)
			return nil, err
		}
		body, err := cl.Do(req)
		if err != nil {
			r.logger.Errorf("Failed to complete the request: %v", err)
			return nil, err
		}
		defer body.Body.Close()
		pg, err := page.NewPage(body.Body, url)
		if err != nil {
			r.logger.Errorf("Failed to create the NewPage: %v", err)
			return nil, err
		}
		if r.debugType {

			p := page.NewloggerPageWrap(pg, r.logger)
			return p, nil
		}

		return pg, nil
	}
}