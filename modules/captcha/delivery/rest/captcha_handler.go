package deliveryRest

import (
	"net/http"
	"playground-go-api/domain"
	"playground-go-api/domain/errcode"
	"playground-go-api/infrastructures/tools"

	"github.com/gin-gonic/gin"
)

type captchaHandler struct {
	captchaUsecase domain.CaptchaUsecase
	e              *gin.Engine
	CaptchaDelivery
}

type CaptchaDelivery interface {
	SmsCaptchaSend(c *gin.Context)
	SmsCaptchaValidate(c *gin.Context)
}

func NewCaptchaHandler(e *gin.Engine, usc domain.CaptchaUsecase) {
	handler := &captchaHandler{
		captchaUsecase: usc,
		e:              e,
	}
	route := e.Group("/captcha")
	route.POST("/sms/send", handler.SmsCaptchaSend)
	route.POST("/sms/validate", handler.SmsCaptchaValidate)
}

// Send Captcha SMS
// @Summary Send Captcha SMS
// @Description Send Captcha SMS
// @Tags captcha
// @Produce json
// @Param default body domain.ReqSmsCaptchaSendDTO true "Send Captcha SMS 內容"
// @Success 200 {object} domain.ResSmsCaptchaSendDTO "Send Captcha SMS 成功"
// @Failure 400 {object} domain.ResponseError "Send Captcha SMS 失敗"
// @Router /captcha/sms/send [post]
func (h *captchaHandler) SmsCaptchaSend(c *gin.Context) {
	var body domain.ReqSmsCaptchaSendDTO
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, tools.NewResponseErr(err, errcode.Default, c.Request.URL.Path, tools.Int(http.StatusBadRequest)))
		return
	}
	result, uCaseErr := h.captchaUsecase.SmsCaptchaSend(&body)
	if uCaseErr != nil {
		c.JSON(tools.ParseErrStatus(uCaseErr.Err), tools.NewResponseErr(uCaseErr.Err, uCaseErr.ErrorCode, c.Request.URL.Path, nil))
		return
	}
	c.JSON(http.StatusOK, result)
}

// 驗證 Captcha SMS
// @Summary 驗證 Captcha SMS
// @Description 驗證 Captcha SMS，不會檢查是否已經有使用者及電話號碼
// @Tags captcha
// @Produce json
// @Param default body domain.ReqSmsCaptchaValidateDTO true "驗證 Captcha SMS 內容"
// @Success 200 {object} domain.ResSmsCaptchaValidateDTO "驗證 Captcha SMS 成功"
// @Success 401 {object} domain.ResSmsCaptchaValidateDTO "驗證 Captcha SMS 失敗"
// @Failure 400 {object} domain.ResponseError "驗證 Captcha SMS 失敗"
// @Router /captcha/sms/validate [post]
func (h *captchaHandler) SmsCaptchaValidate(c *gin.Context) {
	var body domain.ReqSmsCaptchaValidateDTO
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, tools.NewResponseErr(err, errcode.Default, c.Request.URL.Path, tools.Int(http.StatusBadRequest)))
		return
	}
	result, uCaseErr := h.captchaUsecase.SmsCaptchaValidate(&body)
	if uCaseErr != nil {
		if result != nil {
			c.JSON(tools.ParseErrStatus(uCaseErr.Err), result)
			return
		}
		c.JSON(tools.ParseErrStatus(uCaseErr.Err), tools.NewResponseErr(uCaseErr.Err, uCaseErr.ErrorCode, c.Request.URL.Path, nil))
		return
	}
	c.JSON(http.StatusOK, result)
}
