package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loadlab/web/middleware"
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

	var loginResp struct {
		Token string `json:"token"`
	}
	loginResp.Token = "xxxx"
	c.JSON(http.StatusOK, okResponse(loginResp))
}

func registerHandler(c *gin.Context) {

}

func createTodoHandler(c *gin.Context) {

}

func deleteTodoHandler(c *gin.Context) {

}

func listTodosHandler(c *gin.Context) {

}
