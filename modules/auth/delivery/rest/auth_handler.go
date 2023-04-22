package deliveryRest

import (
	"net/http"
	"playground-go-api/domain"
	"playground-go-api/domain/errcode"
	ginguard "playground-go-api/infrastructures/ginrest/guard"
	"playground-go-api/infrastructures/tools"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authUsecase domain.AuthUsecase
	e           *gin.Engine
}

func NewAuthHandler(e *gin.Engine, usc domain.AuthUsecase) {
	handler := &authHandler{
		authUsecase: usc,
		e:           e,
	}
	root := e.Group("/auth")
	root.POST("/login/account", handler.LoginAccount)
	root.POST("/login/telephone", handler.LoginTelephone)
	root.POST("/jwt/decode", handler.JwtDecode)
	root.GET("/jwt/verify", ginguard.AuthGuard, handler.JwtVerify)
	root.GET("/jwt/verify-expired", ginguard.AuthGuardExpired, handler.JwtVerifyExpired)
	root.POST("/refresh-access-token", ginguard.AuthGuard, handler.RefreshAccessToken)
	root.POST("/logout", ginguard.AuthGuardExpired, handler.Logout)
}

// Auth Account Login
// @Summary Auth Account Login
// @Description Auth Account Login
// @Tags auth
// @Produce json
// @Param default body domain.ReqLoginAccountDTO true "account login 內容"
// @Success 200 {object} domain.ResLoginTokenDTO "success login response"
// @Failure 400 {object} domain.ResponseError "請求的body、header驗證失敗"
// @Router /auth/login/account [post]
func (h *authHandler) LoginAccount(c *gin.Context) {
	var body domain.ReqLoginAccountDTO
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, tools.NewResponseErr(err, errcode.Default, c.Request.URL.Path, tools.Int(http.StatusBadRequest)))
		return
	}
	result, uCaseErr := h.authUsecase.LoginAccount(&body)
	if uCaseErr != nil {
		status := tools.ParseErrStatus(uCaseErr.Err)
		c.JSON(status, tools.NewResponseErr(uCaseErr.Err, uCaseErr.ErrorCode, c.Request.URL.Path, tools.Int(status)))
		return
	}
	c.JSON(http.StatusOK, result)
}

// Auth Telephone Login
// @Summary Auth Telephone Login
// @Description Auth Telephone Login
// @Tags auth
// @Produce json
// @Param default body domain.ReqLoginTelephoneDTO true "telephone login 內容"
// @Success 200 {object} domain.ResLoginTokenDTO "success login response"
// @Failure 400 {object} domain.ResponseError "請求的body、header驗證失敗"
// @Router /auth/login/telephone [post]
func (h *authHandler) LoginTelephone(c *gin.Context) {
	var body domain.ReqLoginTelephoneDTO
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, tools.NewResponseErr(err, errcode.Default, c.Request.URL.Path, tools.Int(http.StatusBadRequest)))
		return
	}
	result, uCaseErr := h.authUsecase.LoginTelephone(&body)
	if uCaseErr != nil {
		status := tools.ParseErrStatus(uCaseErr.Err)
		c.JSON(status, tools.NewResponseErr(uCaseErr.Err, uCaseErr.ErrorCode, c.Request.URL.Path, tools.Int(status)))
		return
	}
	c.JSON(http.StatusOK, result)
}

// 當access token失效時，拿refresh token來重新產生一組
// @Summary refresh token 產生token
// @Description 當access token失效時，拿refresh token來重新產生一組
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Param default body domain.ReqRefreshAccessTokenDTO true "telephone login 內容"
// @Success 200 {object} domain.ResLoginTokenDTO "success login response"
// @Failure 400 {object} domain.ResponseError "請求的body、header驗證失敗"
// @Router /auth/refresh-access-token [post]
func (h *authHandler) RefreshAccessToken(c *gin.Context) {
	claims, ok := c.MustGet("claims").(*domain.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, tools.NewResponseErr(domain.ErrUnauthorized, errcode.JwtVerifyError, c.Request.URL.Path, tools.Int(http.StatusUnauthorized)))
		return
	}
	var body domain.ReqRefreshAccessTokenDTO
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, tools.NewResponseErr(err, errcode.Default, c.Request.URL.Path, tools.Int(http.StatusBadRequest)))
		return
	}
	result, uCaseErr := h.authUsecase.RefreshAccessToken(claims, &body)
	if uCaseErr != nil {
		status := tools.ParseErrStatus(uCaseErr.Err)
		c.JSON(status, tools.NewResponseErr(uCaseErr.Err, uCaseErr.ErrorCode, c.Request.URL.Path, tools.Int(status)))
		return
	}
	c.JSON(http.StatusOK, result)
}

// Logout Auth 登出
// @Summary Auth 登出
// @Description Auth 登出
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.ResLogoutDTO "User has been successfully logout."
// @Failure 400 {object} domain.ResponseError "請求的body、header驗證失敗"
// @Router /auth/logout [post]
func (h *authHandler) Logout(c *gin.Context) {
	claims, ok := c.MustGet("claims").(*domain.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, tools.NewResponseErr(domain.ErrUnauthorized, errcode.JwtVerifyError, c.Request.URL.Path, tools.Int(http.StatusUnauthorized)))
		return
	}
	var body domain.ReqRefreshAccessTokenDTO
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, tools.NewResponseErr(err, errcode.Default, c.Request.URL.Path, tools.Int(http.StatusBadRequest)))
		return
	}
	result, uCaseErr := h.authUsecase.Logout(claims, &body)
	if uCaseErr != nil {
		status := tools.ParseErrStatus(uCaseErr.Err)
		c.JSON(http.StatusBadRequest, tools.NewResponseErr(uCaseErr.Err, uCaseErr.ErrorCode, c.Request.URL.Path, tools.Int(status)))
		return
	}
	c.JSON(http.StatusOK, result)
}

// JwtVerify
// @Summary JwtVerify
// @Description JwtVerify 驗證jwt，成功返回內容
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.ResponseError 驗證jwt，成功返回內容
// @Failure 400 {object} domain.ResponseError "請求的body、header驗證失敗"
// @Router /auth/jwt/verify [get]
func (h *authHandler) JwtVerify(c *gin.Context) {
	// c.MustBindWith()
	result, ok := c.MustGet("claims").(*domain.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, tools.NewResponseErr(domain.ErrUnauthorized, errcode.JwtVerifyError, c.Request.URL.Path, tools.Int(http.StatusUnauthorized)))
		return
	}
	c.JSON(http.StatusOK, result)
}

// JwtVerifyExpired
// @Summary JwtVerifyExpired
// @Description JwtVerifyExpired 驗證過期jwt，成功返回內容
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.ResponseError 驗證過期jwt，成功返回內容
// @Failure 400 {object} domain.ResponseError "請求的body、header驗證失敗"
// @Router /auth/jwt/verify-expired [get]
func (h *authHandler) JwtVerifyExpired(c *gin.Context) {
	result, ok := c.MustGet("claims").(*domain.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, tools.NewResponseErr(domain.ErrUnauthorized, errcode.JwtVerifyError, c.Request.URL.Path, tools.Int(http.StatusUnauthorized)))
		return
	}
	c.JSON(http.StatusOK, result)
}

// JwtDecode
// @Summary JwtDecode
// @Description JwtDecode 解析jwt 內容
// @Tags auth
// @Produce json
// @Param default body domain.ReqAuthJwtDecodeDTO true "jwt"
// @Success 200 {object} domain.ResponseError 解析jwt 內容
// @Failure 400 {object} domain.ResponseError "請求的body、header驗證失敗"
// @Router /auth/jwt/decode [post]
func (h *authHandler) JwtDecode(c *gin.Context) {
	var body domain.ReqAuthJwtDecodeDTO
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, tools.NewResponseErr(err, errcode.Default, c.Request.URL.Path, tools.Int(http.StatusBadRequest)))
		return
	}
	result, uCaseErr := h.authUsecase.AuthJwtDecode(&body)
	if uCaseErr != nil {
		status := tools.ParseErrStatus(uCaseErr.Err)
		c.JSON(status, tools.NewResponseErr(uCaseErr.Err, uCaseErr.ErrorCode, c.Request.URL.Path, tools.Int(status)))
		return
	}
	c.JSON(http.StatusOK, result)
}
