package main

import (
	"flag"
	"net"
	"net/http"

	"go.uber.org/zap"
)

var (
	flagListenAddr = flag.String("l", ":8080", "listen address")
)

func main() {
	flag.Parse()

	srv := &http.Server{
		Handler: buildRouter(),
	}

	l, err := net.Listen("tcp", *flagListenAddr)
	if err != nil {
		logger.Fatal("listen failed", zap.Error(err))
	}
	logger.Info("web server startup", zap.String("listen", l.Addr().String()))

	err = srv.Serve(l)
	if err != nil {
		logger.Info("web server serve failed", zap.Error(err))
	}
}
