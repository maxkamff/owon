package v1

import (
	"errors"
	"net/http"

	//"strings"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	models "gitlab.com/api-gateway/api/handlers/models"
	"gitlab.com/api-gateway/api/token"
	"gitlab.com/api-gateway/config"
	"gitlab.com/api-gateway/pkg/logger"
	"gitlab.com/api-gateway/services"
	"gitlab.com/api-gateway/storage/repo"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redis          repo.RedisRepo
	jwtHandler     token.JWTHandler
	enforcer       *casbin.Enforcer
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.RedisRepo
	Jwthandler     token.JWTHandler
	Enforcer       *casbin.Enforcer
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redis:          c.Redis,
		jwtHandler:     c.Jwthandler,
		enforcer:       c.Enforcer,
	}
}

func GetClaims(h *handlerV1, c *gin.Context) jwt.MapClaims {
	var (
		ErrUnAuthorized = errors.New("UnAuthorized")
		authorization   models.GetProfileByJwtRequestModel
		claims          jwt.MapClaims
		err             error
	)
	authorization.Token = c.GetHeader("Authorization")
	if c.Request.Header.Get("Authorization") == "" {
		ErrorCodeUnauthirized := http.StatusUnauthorized
		c.JSON(http.StatusUnauthorized, models.Error{
			Error:       err,
			Code:        ErrorCodeUnauthirized,
			Description: "UnAuthorized",
		})
		h.log.Error("UnAuthorized request ", logger.Error(ErrUnAuthorized))
		return nil
	}


	//authorization.Token = strings.TrimSpace(strings.Trim(authorization.Token, "Bearer"))

	h.jwtHandler.Token = authorization.Token
	claims, err = h.jwtHandler.ExtractClaims()
	if err != nil {
		ErrorCodeUnauthirized := http.StatusUnauthorized
		c.JSON(http.StatusUnauthorized, models.Error{
			Error:       err,
			Code:        ErrorCodeUnauthirized,
			Description: "UnAuthorized",
		})
		h.log.Error("UnAuthorized request ", logger.Error(err))
		return nil
	}

	return claims

}
