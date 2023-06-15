package api

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	_ "gitlab.com/api-gateway/api/docs"
	v1 "gitlab.com/api-gateway/api/handlers/v1"
	"gitlab.com/api-gateway/api/middleware"
	"gitlab.com/api-gateway/api/token"
	"gitlab.com/api-gateway/config"
	"gitlab.com/api-gateway/pkg/logger"
	"gitlab.com/api-gateway/services"
	"gitlab.com/api-gateway/storage/repo"

	swaggerfile "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Option ...
type Option struct {
	Conf            config.Config
	Logger          logger.Logger
	ServiceManager  services.IServiceManager
	InMemoryStorage repo.RedisRepo
	CasbinEnforcer *casbin.Enforcer
}

// New ...

// @title Go-BootCamp N7
// @version 1.0
// @host localhost:7000

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	jwtHandler := token.JWTHandler{
		SigninKey: option.Conf.SignInKey,
		Log:       option.Logger,
	}
	router.Use(middleware.NewAuthorizer(option.CasbinEnforcer, jwtHandler, option.Conf))
	
	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Redis:          option.InMemoryStorage,
		Jwthandler:     jwtHandler,
		Enforcer: *&option.CasbinEnforcer,
	})
	
	api := router.Group("/v1")

	// users
	api.POST("/users", handlerV1.CreateUser)
	api.GET("/user", handlerV1.GetUserById)
	api.GET("/users", handlerV1.GetAllUsers)
	api.DELETE("/user", handlerV1.DeleteUser)
	api.PUT("/user", handlerV1.UpdateUser)

	// posts
	api.POST("/post", handlerV1.CreatePost)
	api.PUT("/post", handlerV1.UpdatePost)
	api.GET("/post", handlerV1.GetPostById)
	api.GET("/posts", handlerV1.GetPostsByUserId)
	api.DELETE("/post", handlerV1.DeletePost)

	// comments
	api.POST("/comments", handlerV1.CreateComment)
	api.DELETE("/comments/:id", handlerV1.DeleteComment)
	api.PUT("comments/:id", handlerV1.UpdateComment)
	api.GET("/comments/:id", handlerV1.GetComment)
	api.GET("/comments", handlerV1.GetAllCommentsByPostId)

	// register
	api.POST("/register", handlerV1.Register)
	api.GET("/verify/:email/:code", handlerV1.Verify)

	// login
	api.GET("/login/:email/:password", handlerV1.Login)

	// admin
	api.POST("/admin/add/policy", handlerV1.AddPolicy)
	api.POST("/admin/add/role", handlerV1.AddRoleToUser)
	api.DELETE("admin/delete/policy", handlerV1.RemovePolicy)
	api.POST("/admin/remove/role", handlerV1.RemoveRoleFromUser)
	api.GET("/login/admin/:email/:password", handlerV1.LoginFodrAdmins)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfile.Handler, url))

	// api.PUT("/users/:id", handlerV1.UpdateUser)
	// posts
	return router
}
