package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

var (
	flagListenAddr      = flag.String("listen", ":8080", "listen address")
	flagEtcdEndpoints   = flag.String("etcd-endpoints", "localhost:2379", "etcd endpoints")
	flagAdvertiseClient = flag.String("advertise-client", "localhost:8080", "advertise client url")
)

func main() {
	flag.Parse()

	tracerCloser, err := initTracer()
	if err != nil {
		logger.Panic("init tracer failed", zap.Error(err))
	}
	defer tracerCloser.Close()

	mustInitEtcdCli(*flagEtcdEndpoints)

	go func() {
		err := registerEndpointWithRetry(*flagAdvertiseClient)
		if err != nil {
			logger.Panic("register endpoint faield", zap.Error(err))
		}
	}()

	mustDiscoverServices()
	logger.Info("discover services succeed")

	tracer := opentracing.GlobalTracer()
	mw := nethttp.Middleware(tracer, buildRouter(), nethttp.OperationNameFunc(func(r *http.Request) string {
		return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path)
	}))
	srv := &http.Server{
		Handler: mw,
	}

	l, err := net.Listen("tcp", *flagListenAddr)
	if err != nil {
		logger.Fatal("listen failed", zap.Error(err))
	}
	logger.Info("web server startup", zap.String("listen", l.Addr().String()))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go signalSet(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_ = srv.Shutdown(ctx)
	})

	err = srv.Serve(l)
	if err == http.ErrServerClosed {
		logger.Info("http server closed")
		return
	}
	if err != nil {
		logger.Panic("web server serve failed", zap.Error(err))
	}
}

func signalSet(cb func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	s := <-sigCh
	logger.Warn("exit signal", zap.String("signal", s.String()))

	cb()
}

func initTracer() (io.Closer, error) {
	cfg := config.Configuration{ServiceName: "websvc"}
	_, err := cfg.FromEnv()
	if err != nil {
		return nil, err
	}
	tracer, tracerCloser, err := cfg.NewTracer()
	if err != nil {
		return nil, err
	}
	opentracing.InitGlobalTracer(tracer)
	return tracerCloser, nil
}
