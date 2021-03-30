package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loadlab/websvc/middleware"
)

func buildRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	api := r.Group("/api", middleware.ZapLog(logger))

	api.POST("/login", loginHandler)
	api.POST("/register", registerHandler)

	todos := api.Group("/todos")
	todos.POST("/", createTodoHandler)
	todos.DELETE("/:id", deleteTodoHandler)
	todos.GET("/", listTodosHandler)

	return r
}
