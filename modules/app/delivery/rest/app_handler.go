package deliveryRest

import (
	"net/http"
	"playground-go-api/domain"

	"github.com/gin-gonic/gin"
)

type appHandler struct {
	AppUsecase domain.AppUsecase
	e          *gin.Engine
	AppDelivery
}

type AppDelivery interface {
	// GetApp 取得App 名稱
	GetApp(c *gin.Context)
	GetSwagger(c *gin.Context)
}

func NewAppHandler(e *gin.Engine, usc domain.AppUsecase) {
	handler := &appHandler{
		AppUsecase: usc,
		e:          e,
	}
	root := e.Group("/")
	root.GET("", handler.GetApp)
	root.GET("/docs", handler.GetSwagger)
	root.GET("/swagger", handler.GetSwagger)
}

// GetApp 取得App 名稱
// @Summary 取得App 名稱
// @Description 取得App 名稱
// @Tags app
// @Produce json
// @Success 200 {object} domain.App
// @Router / [get]
func (aH *appHandler) GetApp(c *gin.Context) {
	result := aH.AppUsecase.GetApp()
	c.JSON(http.StatusOK, result)
}

func (aH *appHandler) GetSwagger(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
}
