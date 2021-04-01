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
	"github.com/loadlab-go/websvc/resp"
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
		c.JSON(http.StatusBadRequest, resp.ErrResponse(fmt.Errorf("parse payload failed: %v", err)))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()
	validateResp, err := userClient.ValidatePassword(ctx, &userpb.ValidatePasswordRequest{Username: loginReq.Username, Password: loginReq.Password})
	if err != nil {
		middleware.MustGetZapLog(c).Warn("validate password failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, resp.ErrResponse(fmt.Errorf("validate password failed: %v", err)))
		return
	}
	tokenResp, err := jwtClient.GenerateJWT(ctx, &authpb.GenerateJWTRequest{Id: validateResp.Id, Username: loginReq.Password})
	if err != nil {
		middleware.MustGetZapLog(c).Warn("generate auth token failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, resp.ErrResponse(fmt.Errorf("generate auth token failed: %v", err)))
		return
	}

	var loginResp struct {
		Token string `json:"token"`
	}
	loginResp.Token = tokenResp.Token

	c.JSON(http.StatusOK, resp.OkResponse(loginResp))
}

func registerHandler(c *gin.Context) {
	var registerReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&registerReq)
	if err != nil {
		middleware.MustGetZapLog(c).Warn("parse payload failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, resp.ErrResponse(fmt.Errorf("parse payload failed: %v", err)))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*3)
	defer cancel()
	createResp, err := userClient.Create(ctx, &userpb.CreateRequest{Username: registerReq.Username, Password: registerReq.Password})
	if err != nil {
		middleware.MustGetZapLog(c).Warn("register failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, resp.ErrResponse(fmt.Errorf("register failed: %v", err)))
		return
	}
	c.JSON(http.StatusOK, resp.OkResponse(createResp.Id))
}

func createTodoHandler(c *gin.Context) {

}

func deleteTodoHandler(c *gin.Context) {

}

func listTodosHandler(c *gin.Context) {

}
