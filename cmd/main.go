package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	app "x_clone_auth_srv"
	config "x_clone_auth_srv/config"
	"x_clone_auth_srv/user"

	"github.com/go-kit/log"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	userRepo := user.NewRepository()
	var (
		httpAddr = flag.String("http.addr", ":"+config.GetEnv("PORT"), "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s app.Service
	{
		s = app.NewService(userRepo, config.GetEnv("JWT_SECRET"))
		s = app.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = app.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
