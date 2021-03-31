package main

import (
	authpb "github.com/loadlab-go/pkg/proto/auth"
	userpb "github.com/loadlab-go/pkg/proto/user"
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
