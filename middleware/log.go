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
	l, err := GetZapLog(c)
	if err != nil {
		panic(err)
	}
	return l
}

func GetZapLog(c *gin.Context) (*zap.Logger, error) {
	v, exist := c.Get(zapLogKey)
	if !exist {
		return nil, errors.New("there is no zap logger")
	}
	logger, ok := v.(*zap.Logger)
	if !ok {
		return nil, fmt.Errorf("type error, expect  *zap.logger, actual is %v", reflect.TypeOf(v).Kind())
	}
	return logger, nil
}
