package middleware

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const zapLogKey = "middleware-zap-log-key"

func ZapLog(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(zapLogKey, logger)
		c.Next()
	}
}

func MustGetZapLog(c *gin.Context) *zap.Logger {
	v, exist := c.Get(zapLogKey)
	if !exist {
		panic(errors.New("there is no zap logger"))
	}
	logger, ok := v.(*zap.Logger)
	if !ok {
		panic(fmt.Errorf("type error, expect  *zap.logger, actual is %v", reflect.TypeOf(v).Kind()))
	}
	return logger
}
