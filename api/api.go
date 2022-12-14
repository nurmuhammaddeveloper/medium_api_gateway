package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/nurmuhammaddeveloper/medium_api_gateway/api/v1"
	"github.com/nurmuhammaddeveloper/medium_api_gateway/config"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/nurmuhammaddeveloper/medium_api_gateway/api/docs" // for swagger

	grpcPkg "github.com/nurmuhammaddeveloper/medium_api_gateway/pkg/grpc_client"
)

type RouterOptions struct {
	Cfg        *config.Config
	GrpcClient grpcPkg.GrpcClientI
	Logger     *logrus.Logger
}

// @title           Swagger for blog api
// @version         1.0
// @description     This is a blog service api.
// @host      localhost:8000
// @BasePath  /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Security ApiKeyAuth
func New(opt *RouterOptions) *gin.Engine {
	router := gin.Default()

	handlerV1 := v1.New(&v1.HandlerV1Options{
		Cfg:        opt.Cfg,
		GrpcClient: opt.GrpcClient,
		Logger:     opt.Logger,
	})

	apiV1 := router.Group("/v1")

	apiV1.POST("/auth/register", handlerV1.Register)
	apiV1.POST("/auth/login", handlerV1.Login)
	apiV1.POST("/auth/verify", handlerV1.Verify)
	apiV1.POST("/auth/forgot-password", handlerV1.ForgotPassword)
	apiV1.POST("/auth/update-password",handlerV1.UpdatePassword)
	apiV1.POST("/auth/verify-forgot-password", handlerV1.VerifyForgotPassword)

	apiV1.POST("/users",handlerV1.AuthMiddleWare("users", "create"), handlerV1.CreateUser)
	apiV1.GET("/users/:id", handlerV1.GetUser)
	apiV1.PUT("/users/:id", handlerV1.UpdateUser)
	apiV1.GET("/users", handlerV1.GetAllUsers)
	apiV1.DELETE("/users/:id", handlerV1.DeleteUser)
	apiV1.GET("/users/email/:email", handlerV1.GetUserByEmail)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
