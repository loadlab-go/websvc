package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/loadlab-go/websvc/resp"
	"go.uber.org/zap"
)

type validator interface {
	Validate(ctx context.Context, token string) error
}

func getJWTFromContext(c *gin.Context) (string, error) {
	v := c.Request.Header.Get("Authorization")
	if v == "" {
		return "", errors.New("missing authorization header")
	}
	if !strings.HasPrefix(v, "Bearer ") {
		return "", fmt.Errorf("authorization header invalid: %q", v)
	}
	jwt := strings.TrimSpace(strings.TrimPrefix(v, "Bearer "))
	if jwt == "" {
		return "", fmt.Errorf("authorization header invalid: %q", v)
	}
	return jwt, nil
}

func AuthRequired(v validator, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt, err := getJWTFromContext(c)
		if err != nil {
			logger.Warn("get jwt failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp.ErrResponse(err))
			return
		}
		err = v.Validate(c.Request.Context(), jwt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp.ErrResponse(err))
			return
		}
		c.Next()
	}
}
