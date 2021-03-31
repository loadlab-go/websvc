package main

import (
	"github.com/loadlab-go/websvc/idl/proto/authpb"
	"github.com/loadlab-go/websvc/idl/proto/userpb"
	"go.uber.org/zap"
)

var (
	authClient authpb.AuthClient
	userClient userpb.UserClient
)

func mustDiscoverServices() error {
	authcc, err := grpcDial("auth-svc")
	if err != nil {
		logger.Panic("grpc dial auth-svc failed", zap.Error(err))
	}
	authClient = authpb.NewAuthClient(authcc)

	usercc, err := grpcDial("user-svc")
	if err != nil {
		logger.Panic("grpc dial user-svc failed", zap.Error(err))
	}
	userClient = userpb.NewUserClient(usercc)

	return nil
}
