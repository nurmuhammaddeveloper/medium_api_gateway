package v1

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nurmuhammaddeveloper/medium_api_gateway/api/models"
	"github.com/nurmuhammaddeveloper/medium_api_gateway/config"
	grpcPkg "github.com/nurmuhammaddeveloper/medium_api_gateway/pkg/grpc_client"

	"github.com/sirupsen/logrus"
)

var (
	ErrWrongEmailOrPass = errors.New("wrong email or password")
	ErrEmailExists      = errors.New("email already exists")
	ErrUserNotVerified  = errors.New("user not verified")
	ErrIncorrectCode    = errors.New("incorrect verification code")
	ErrCodeExpired      = errors.New("verification code has been expired")
	ErrForbidden        = errors.New("forbidden")
	ErrNotAllowed       = errors.New("method not allowed")
)

type handlerV1 struct {
	cfg        *config.Config
	grpcClient grpcPkg.GrpcClientI
	logger     *logrus.Logger
}

type HandlerV1Options struct {
	Cfg        *config.Config
	GrpcClient grpcPkg.GrpcClientI
	Logger     *logrus.Logger
}

func New(options *HandlerV1Options) *handlerV1 {
	return &handlerV1{
		cfg:        options.Cfg,
		grpcClient: options.GrpcClient,
		logger:     options.Logger,
	}
}

func errorResponse(err error) *models.ErrorResponse {
	return &models.ErrorResponse{
		Error: err.Error(),
	}
}

func validateGetAllParams(c *gin.Context) (*models.GetAllParams, error) {
	var (
		limit int = 10
		page  int = 1
		err   error
	)

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			return nil, err
		}
	}

	return &models.GetAllParams{
		Limit:  int32(limit),
		Page:   int32(page),
		Search: c.Query("search"),
	}, nil
}
