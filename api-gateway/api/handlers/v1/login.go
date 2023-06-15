package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/api-gateway/api/handlers/models"
	"gitlab.com/api-gateway/genproto/user"
	"gitlab.com/api-gateway/pkg/etc"
	"gitlab.com/api-gateway/pkg/logger"
)

// User             Login
// @Summary         Login User
// @Description     This Function get login User
// @Tags            Login
// @Accept 			json
// @Produce			json
// @Param 			email 		path string true "email"
// @Param 			password 	path string true "password"
// @Succes 			201 {object}    models.LoginResponse
// @Failure			500 {object} 	models.Error
// @Failure			400 {object} 	models.Error
// @Router			/v1/login/{email}/{password} [get]
func (h *handlerV1) Login(c *gin.Context) {
	var (
		email = c.Param("email")
		password = c.Param("password")
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.UserService().GetByEmail(ctx, &user.Email{Email: email})
	fmt.Println(res)
	if err != nil{
		c.JSON(http.StatusNotFound, models.Error{
			Error:       err,
			Description: "Couln't find matching information, Have you registered before?",
		})
		h.log.Error("Error while getting customer by email", logger.Any("post", err))
		return
	}

	if !etc.CheckPasswordHash(password, res.Password){
		c.JSON(http.StatusNotFound, models.Error{
			Description: "Password or Email error",
			Code:        http.StatusBadRequest,
		})
		return
	}
	
	ans, err := h.serviceManager.UserService().GetUserRole(ctx, &user.CheckRoleRequest{Id: res.Id})
	if err != nil{
		c.JSON(http.StatusNotFound, models.Error{
			Error:       err,
			Description: "You don't have permisson to this function",
		})
		h.log.Error("Error while getting customer by id", logger.Any("post", err))
		return
	}
	h.jwtHandler.Iss = "User"
	h.jwtHandler.Sub = res.Id
	h.jwtHandler.Role = ans.Role
	h.jwtHandler.Aud = []string{"user-frontend"}
	h.jwtHandler.SigninKey = h.cfg.SignInKey
	h.jwtHandler.Log = h.log
	tokens, err := h.jwtHandler.GenerateAuthJWT()
	accesstoken := tokens[0]
	refreshToken := tokens[1]
	if err != nil {
		h.log.Error("error occured while generating tokens")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong,please try again",
		})
		return
	}

	newRes, err := h.serviceManager.UserService().UpdateTokens(ctx, &user.TokenRequest{
		Id: res.Id,
		RefreshToken: refreshToken,
	})

	res.AccessToken = accesstoken
	res.RefreshToken = newRes.RefreshToken
	res.UpdatedAt = newRes.UpdatedAt

	response := models.LoginReponse{
		Id: res.Id,
		Name: res.Name,
		Last_name: res.LastName,
		Username: res.Username,
		Email: res.Email,
		Refresh_token: res.RefreshToken,
		Access_token: res.AccessToken,
		Created_at: res.CreatedAt,
	}
	c.JSON(http.StatusOK, response)
}