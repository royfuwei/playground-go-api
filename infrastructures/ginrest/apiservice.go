package ginrest

import (
	"net/http"
	"os"
	"os/signal"
	"playground-go-api/config"
	appRest "playground-go-api/modules/app/delivery/rest"
	appUcase "playground-go-api/modules/app/usecases"
	authRest "playground-go-api/modules/auth/delivery/rest"
	refreshTokensMgo "playground-go-api/modules/auth/repository/mgo"
	jwtSvc "playground-go-api/modules/auth/service"
	authUcase "playground-go-api/modules/auth/usecases"
	captchaRest "playground-go-api/modules/captcha/delivery/rest"
	captchaMgo "playground-go-api/modules/captcha/repository/mgo"
	captchaSvc "playground-go-api/modules/captcha/service"
	captchaUcase "playground-go-api/modules/captcha/usecase"
	configsMgo "playground-go-api/modules/configs/repository/mgo"
	configsUcase "playground-go-api/modules/configs/usecases"
	encryptSvc "playground-go-api/modules/encrypt/service"
	registerRest "playground-go-api/modules/register/delivery/rest"
	registerUcase "playground-go-api/modules/register/usecase"
	smsSvc "playground-go-api/modules/sms/service"

	twsmsSmsUcase "playground-go-api/modules/sms/usecase"

	// msgSmsUcase "playground-go-api/modules/sms/usecase"
	usersRest "playground-go-api/modules/users/delivery/rest"
	usersMgo "playground-go-api/modules/users/repository/mgo"
	usersSvc "playground-go-api/modules/users/service"
	usersUcase "playground-go-api/modules/users/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIService struct{}

func NewAPIService() *APIService {
	return &APIService{}
}

// Start api service init and start
func (api *APIService) Start(mongoClient *mongo.Client) {
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	// Force log's color
	gin.ForceConsoleColor()
	// Disable log's color
	// gin.DisableConsoleColor()

	encryptService := encryptSvc.NewEncryptService()
	usersService := usersSvc.NewUsersService()
	jwtService := jwtSvc.NewJwtService()
	captchaService := captchaSvc.NewCaptchaService()
	smsService := smsSvc.NewSmsService()
	/* usecase, delivery 注入router */
	configsRepo := configsMgo.NewMgoConfigsRepository(mongoClient)
	usersRepo := usersMgo.NewMgoUsersRepository(mongoClient)
	refreshTokensRepo := refreshTokensMgo.NewMgoRefreshTokensRepository(mongoClient)
	captchaRepo := captchaMgo.NewMgoCaptchaRepository(mongoClient)

	// smsUsecase := msgSmsUcase.NewMessageSmsUsecase()
	smsUsecase := twsmsSmsUcase.NewTwsmsSmsUsecase()
	appUsecase := appUcase.NewAppUsecase()
	configsUsecase := configsUcase.NewConfigsUsecase(configsRepo, encryptService)
	usersUsecase := usersUcase.NewUsersUsecase(usersRepo, configsRepo, usersService, encryptService, smsUsecase)
	captchaUsecase := captchaUcase.NewCaptchaUsecase(captchaService, captchaRepo, usersUsecase, smsService, jwtService, smsUsecase)
	authUsecase := authUcase.NewAuthUsecase(refreshTokensRepo, jwtService, usersRepo, usersService)
	registerUsecase := registerUcase.NewRegisterUsecase(usersUsecase, authUsecase, usersRepo, captchaUsecase)

	// init initial/configs.json
	configsUsecase.InitConfigs()
	configsUsecase.GetAccountExpireAt()
	configsUsecase.GetCreateUserDefPwd()
	// init initial/users.json
	usersUsecase.InitUsers()

	usersRest.NewUsersHandler(r, usersUsecase)
	appRest.NewAppHandler(r, appUsecase)
	authRest.NewAuthHandler(r, authUsecase)
	captchaRest.NewCaptchaHandler(r, captchaUsecase)
	registerRest.NewRegisterHandler(r, registerUsecase)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := &http.Server{
		Addr:    config.Cfgs.Port,
		Handler: r,
	}

	api.gracefulShutdown(server)
	glog.Infof("Start API Service: 127.0.0.1%s", config.Cfgs.Port)
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			glog.Info("Server closed under request")
		} else {
			glog.Fatalf("Server closed unexpect: %v", err)
		}
	}
	glog.Info("Server exiting")
}

func (a *APIService) gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		glog.Info("receive interrupt signal")
		if err := server.Close(); err != nil {
			glog.Fatal("Server Close:", err)
		}
	}()
}

// func (a *APIService) validateMgoID() validator.Func {
// 	return func(field validator.FieldLevel) bool {
// 		switch field.Field().Kind() {
// 		case reflect.Invalid:
// 			return false
// 		case reflect.String:
// 			return bson.IsObjectIdHex(field.Field().String())
// 		case reflect.Slice:
// 			for i := 0; i < field.Field().Len(); i++ {
// 				e := field.Field().Index(i)
// 				if !bson.IsObjectIdHex(e.String()) {
// 					return false
// 				}
// 			}
// 			return true
// 		}
// 		return true
// 	}
// }
