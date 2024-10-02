package x_clone_auth_srv

import (
	"context"
	"time"

	"github.com/go-kit/log"
)

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type Middleware func(Service) Service

func (mw loggingMiddleware) Authenticate(ctx context.Context, username, password string) (jwtToken string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Authenticate", "username", username, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Authenticate(ctx, username, password)
}

func (mw loggingMiddleware) SignUp(ctx context.Context, username, password string) (jwtToken string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "SignUp", "username", username, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.SignUp(ctx, username, password)
}
