package ginguard

import (
	"net/http"
	"playground-go-api/domain"
	"playground-go-api/domain/errcode"
	"playground-go-api/infrastructures/tools"
	jwtSvc "playground-go-api/modules/auth/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthGuard(c *gin.Context) {
	jwtService := jwtSvc.NewJwtService()
	auth := c.GetHeader("Authorization")
	if len(auth) == 0 {
		c.Abort()
		c.JSON(http.StatusUnauthorized, tools.NewResponseErr(domain.ErrAuthorizationEmpty, errcode.ErrAuthorizationEmpty, c.Request.URL.Path, tools.Int(http.StatusUnauthorized)))
		return
	}
	auths := strings.Split(auth, " ")
	if len(auths) < 2 || strings.ToLower(auths[0]) != "bearer" {
		c.Abort()
		c.JSON(http.StatusUnauthorized, tools.NewResponseErr(domain.ErrBearerNotValid, errcode.ErrBearerNotValid, c.Request.URL.Path, tools.Int(http.StatusUnauthorized)))
		return
	}
	claims, err := jwtService.JwtVerify(auths[1])
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, tools.NewResponseErr(domain.ErrUnauthorized, errcode.JwtVerifyError, c.Request.URL.Path, tools.Int(http.StatusUnauthorized)))
		return
	}
	c.Set("claims", claims)
	c.Next()
}

func AuthGuardExpired(c *gin.Context) {
	jwtService := jwtSvc.NewJwtService()
	auth := c.GetHeader("Authorization")
	if len(auth) == 0 {
		c.Abort()
		c.JSON(http.StatusUnauthorized, tools.NewResponseErr(domain.ErrAuthorizationEmpty, errcode.ErrAuthorizationEmpty, c.Request.URL.Path, tools.Int(http.StatusUnauthorized)))
		return
	}
	auths := strings.Split(auth, " ")
	if len(auths) < 2 || strings.ToLower(auths[0]) != "bearer" {
		c.Abort()
		c.JSON(http.StatusUnauthorized, tools.NewResponseErr(domain.ErrBearerNotValid, errcode.ErrBearerNotValid, c.Request.URL.Path, tools.Int(http.StatusUnauthorized)))
		return
	}
	claims, err := jwtService.JwtVerifyExpired(auths[1])
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, tools.NewResponseErr(domain.ErrUnauthorized, errcode.JwtVerifyError, c.Request.URL.Path, tools.Int(http.StatusUnauthorized)))
		return
	}
	c.Set("claims", claims)
	c.Next()
}
