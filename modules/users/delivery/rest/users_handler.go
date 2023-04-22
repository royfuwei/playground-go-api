package deliveryRest

import (
	"playground-go-api/domain"

	"github.com/gin-gonic/gin"
)

type usersHandler struct {
	usersDelivery
	usersUsecase domain.UsersUsecase
	e            *gin.Engine
}

type usersDelivery interface {
	RegisterSendCaptchaSms(c *gin.Context)
	RegisterValidateCaptchaSms(c *gin.Context)
}

func NewUsersHandler(e *gin.Engine, usc domain.UsersUsecase) {
	handler := &usersHandler{
		usersUsecase: usc,
		e:            e,
	}
	path := e.Group("/users")
	path.POST("/register/captcha/sms", handler.RegisterSendCaptchaSms)
}

func (h *usersHandler) RegisterSendCaptchaSms(c *gin.Context) {
	/* var body domain.ReqRegisterTelDTO
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, tools.NewResponseErr(err, errcode.Default, c.Request.URL.Path, tools.Int(http.StatusBadRequest)))
		return
	}
	result, uCaseErr := h.usersUsecase.RegisterByTelephone(&body)
	if uCaseErr != nil {
		c.JSON(tools.ParseErrStatus(uCaseErr.Err), tools.NewResponseErr(uCaseErr.Err, uCaseErr.ErrorCode, c.Request.URL.Path, nil))
		return
	}
	c.JSON(http.StatusOK, result) */
}
