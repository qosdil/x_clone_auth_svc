package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	app "x_clone_auth_svc"
	config "x_clone_auth_svc/config"
	"x_clone_auth_svc/internal/pkg/user"
	userGrpcSvc "x_clone_auth_svc/internal/pkg/user/grpc/service"

	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to User Svc on gRPC
	userGrpcClientConn, err := grpc.Dial(config.GetEnv("USER_GRPC_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer userGrpcClientConn.Close()

	// gRPC client of User Service
	userGrpcClient := userGrpcSvc.NewServiceClient(userGrpcClientConn)

	userRepo := user.NewRepository(userGrpcClient)
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
