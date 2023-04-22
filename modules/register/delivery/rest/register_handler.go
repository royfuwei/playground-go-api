package deliveryRest

import (
	"net/http"
	"playground-go-api/domain"
	"playground-go-api/domain/errcode"
	ginguard "playground-go-api/infrastructures/ginrest/guard"
	"playground-go-api/infrastructures/tools"

	"github.com/gin-gonic/gin"
)

type registerHandler struct {
	// registerDelivery
	registerUcase domain.RegisterUsecase
	e             *gin.Engine
}

// type registerDelivery interface{}

func NewRegisterHandler(e *gin.Engine, usc domain.RegisterUsecase) {
	handler := &registerHandler{
		registerUcase: usc,
		e:             e,
	}
	root := e.Group("/register")
	root.POST("/captcha/sms/send", handler.RegisterSmsCaptchaSend)
	root.POST("/telephone/user", ginguard.AuthGuard, handler.RegisterUserByTelephone)
}

// @Summary send sms Register captcha
// @Description 發送註冊的sms captcha
// @Tags register
// @Produce json
// @Param default body domain.ReqSmsCaptchaSendDTO true "發送的手機號碼"
// @Success 200 {object} domain.ResSmsCaptchaSendDTO "success"
// @Failure 400 {object} domain.ResponseError "請求的body、header驗證失敗"
// @Router /register/captcha/sms/send [post]
func (h *registerHandler) RegisterSmsCaptchaSend(c *gin.Context) {
	var body domain.ReqSmsCaptchaSendDTO
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, tools.NewResponseErr(err, errcode.Default, c.Request.URL.Path, tools.Int(http.StatusBadRequest)))
		return
	}
	result, uCaseErr := h.registerUcase.RegisterSmsCaptchaSend(&body)
	if uCaseErr != nil {
		c.JSON(tools.ParseErrStatus(uCaseErr.Err), tools.NewResponseErr(uCaseErr.Err, uCaseErr.ErrorCode, c.Request.URL.Path, nil))
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary Register User By telephone
// @Description Register User By telephone
// @Tags register
// @Security BearerAuth
// @Produce json
// @Param default body domain.ReqCreateUserByTelephone true "建立"
// @Success 200 {object} domain.ResLoginTokenDTO "success"
// @Failure 400 {object} domain.ResponseError "請求的body、header驗證失敗"
// @Router /register/telephone/user [post]
func (h *registerHandler) RegisterUserByTelephone(c *gin.Context) {
	claims, ok := c.MustGet("claims").(*domain.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, tools.NewResponseErr(domain.ErrUnauthorized, errcode.JwtVerifyError, c.Request.URL.Path, tools.Int(http.StatusUnauthorized)))
		return
	}
	var body domain.ReqCreateUserByTelephone
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, tools.NewResponseErr(err, errcode.Default, c.Request.URL.Path, tools.Int(http.StatusBadRequest)))
		return
	}
	result, uCaseErr := h.registerUcase.RegisterUserByTelephone(claims, &body)
	if uCaseErr != nil {
		c.JSON(tools.ParseErrStatus(uCaseErr.Err), tools.NewResponseErr(uCaseErr.Err, uCaseErr.ErrorCode, c.Request.URL.Path, nil))
		return
	}
	c.JSON(http.StatusOK, result)
}
