package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	authpb "github.com/loadlab-go/pkg/proto/auth"
	userpb "github.com/loadlab-go/pkg/proto/user"
	"github.com/loadlab-go/websvc/middleware"
	"go.uber.org/zap"
)

func loginHandler(c *gin.Context) {
	var loginReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&loginReq)
	if err != nil {
		middleware.MustGetZapLog(c).Warn("parse payload failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, errResponse(fmt.Errorf("parse payload failed: %v", err)))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*3)
	defer cancel()
	authResp, err := authClient.Authenticate(ctx, &authpb.AuthenticateRequest{Username: loginReq.Password, Password: loginReq.Password})
	if err != nil {
		middleware.MustGetZapLog(c).Warn("authenticate failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errResponse(fmt.Errorf("authenticate failed: %v", err)))
		return
	}

	var loginResp struct {
		JWT string `json:"jwt"`
	}
	loginResp.JWT = authResp.Jwt

	c.JSON(http.StatusOK, okResponse(loginResp))
}

func registerHandler(c *gin.Context) {
	var registerReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&registerReq)
	if err != nil {
		middleware.MustGetZapLog(c).Warn("parse payload failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, errResponse(fmt.Errorf("parse payload failed: %v", err)))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*3)
	defer cancel()
	createResp, err := userClient.Create(ctx, &userpb.CreateRequest{Username: registerReq.Username, Password: registerReq.Password})
	if err != nil {
		middleware.MustGetZapLog(c).Warn("register failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errResponse(fmt.Errorf("register failed: %v", err)))
		return
	}
	c.JSON(http.StatusOK, okResponse(createResp.Id))
}

func createTodoHandler(c *gin.Context) {

}

func deleteTodoHandler(c *gin.Context) {

}

func listTodosHandler(c *gin.Context) {

}
