package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/api-gateway/api/handlers/models"
	"gitlab.com/api-gateway/api/token"
	"gitlab.com/api-gateway/config"
)

type JWTRoleAuthorizer struct {
	enforcer *casbin.Enforcer
	cfg      config.Config
	jwtHandler token.JWTHandler
}

// NewAuthorizer is a middleware for gin get role and
// allow or deny access to endpoints
func NewAuthorizer(e *casbin.Enforcer, jwtHandler token.JWTHandler, cfg config.Config) gin.HandlerFunc{
	a := &JWTRoleAuthorizer{
		enforcer: e,
		cfg: cfg,
		jwtHandler: jwtHandler,
	}

	return func(c *gin.Context){
		allow, err := a.CheckPermission(c.Request)
		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			if v.Errors == jwt.ValidationErrorExpired {
				a.RequireRefresh(c)
			}else {
				a.RequirePermisson(c)
			}
		} else if !allow{
			a.RequirePermisson(c)
		}
	}
}

// GetRole Gets role from Authorization header if there is a token then it is
// parsed and in role got from role claim. If there is no token then role is
// Unauthorized
func (a *JWTRoleAuthorizer) GetRole(r *http.Request) (string, error) {
	var(
		role string 
		claims jwt.MapClaims
		err error
	)

	jwtToken := r.Header.Get("Authorization")
	if jwtToken == ""{
		return "unauthorized", nil
	}else if strings.Contains(jwtToken, "Basic"){
		return "unauthorized", nil
	}

	a.jwtHandler.Token =  jwtToken
	claims, err = a.jwtHandler.ExtractClaims()
	if err != nil{
		return "", err
	}

	if claims["role"].(string) == "user"{
		role = "user"
	} else if claims["role"].(string) == "admin"{
		role = "admin"
	} else {
		role = "unknown"
	}
	return role, nil
}

//CheckPermission checks whether user is allowed to use certain endpoint
func (a *JWTRoleAuthorizer) CheckPermission(r *http.Request) (bool, error){
	user, err := a.GetRole(r)
	if err != nil{
		return false, err
	}
	method := r.Method
	path := r.URL.Path

	fmt.Println("User: ", user)
	fmt.Println("Method: ", method)
	fmt.Println("Path: ", path)

	allowed, err := a.enforcer.Enforce(user, path, method)
	if err != nil {
		panic(err)
	}
	 return allowed, nil
}

func (a *JWTRoleAuthorizer) RequirePermisson(c *gin.Context){
	c.AbortWithStatus(403)
}

func (a *JWTRoleAuthorizer) RequireRefresh(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, models.Error{
		Code: http.StatusUnauthorized,
		Description: "Token Is Expired",
	})
	c.AbortWithStatus(401)
}