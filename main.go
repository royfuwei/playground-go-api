package main

import (
	"playground-go-api/config"
	"playground-go-api/infrastructures/ginrest"
	"playground-go-api/infrastructures/mongodb"
	_ "playground-go-api/swagger"
	"runtime"
	"time"

	"github.com/golang/glog"
)

// @title Seeks Auth Service
// @version 1.0
// @description Seeks Auth API

// @host
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() + 1)
	/* new config to global value */
	config.NewConfig()

	glog.Infof("goroutine number: %d", config.Cfgs.Goroutine)
	/* mongodb connection */
	mongoClient, ctx := mongodb.NewMongoClient(config.Cfgs.MongoAddr)
	defer mongoClient.Disconnect(ctx)
	/* gin rest api service */
	forever := make(chan bool)
	go func() {
		apiService := ginrest.NewAPIService()
		apiService.Start(mongoClient)
		time.Sleep(0)
	}()
	<-forever
}
