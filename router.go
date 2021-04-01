package main

import (
	"context"

	"github.com/gin-gonic/gin"
	authpb "github.com/loadlab-go/pkg/proto/auth"
	"github.com/loadlab-go/websvc/middleware"
)

type jwtValidator struct {
}

func (v *jwtValidator) Validate(ctx context.Context, jwt string) (middleware.Claims, error) {
	validateResp, err := jwtClient.ValidateJWT(ctx, &authpb.ValidateJWTRequest{Token: jwt})
	if err != nil {
		return middleware.Claims{}, err
	}

	return middleware.Claims{
		Aud: validateResp.Aud,
		Exp: validateResp.Exp,
		Jti: validateResp.Jti,
		Iat: validateResp.Iat,
		Iss: validateResp.Iss,
		Nbf: validateResp.Nbf,
		Sub: validateResp.Sub,
	}, nil
}

func authRequired() gin.HandlerFunc {
	return middleware.AuthRequired(&jwtValidator{}, logger)
}

func buildRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	api := r.Group("/api", middleware.ZapLog(logger))

	api.POST("/login", loginHandler)
	api.POST("/register", registerHandler)

	todos := api.Group("/todos", authRequired())
	todos.POST("/", createTodoHandler)
	todos.DELETE("/:id", deleteTodoHandler)
	todos.GET("/", listTodosHandler)

	return r
}
