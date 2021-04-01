package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/loadlab-go/websvc/resp"
	"go.uber.org/zap"
)

const claimsKey = "middleware-claims-key"

type validator interface {
	Validate(ctx context.Context, token string) (Claims, error)
}

type Claims struct {
	Aud []string
	Exp int64
	Jti string
	Iat int64
	Iss string
	Nbf int64
	Sub string
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

		claims, err := v.Validate(c.Request.Context(), jwt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp.ErrResponse(err))
			return
		}
		c.Set(claimsKey, claims)
		c.Next()
	}
}

func GetClaims(c *gin.Context) (Claims, error) {
	v, exist := c.Get(claimsKey)
	if !exist {
		return Claims{}, errors.New("there is no claims")
	}
	claims, ok := v.(Claims)
	if !ok {
		return Claims{}, fmt.Errorf("type error, expect Claims, actual is %v", reflect.TypeOf(v).Kind())
	}
	return claims, nil
}
