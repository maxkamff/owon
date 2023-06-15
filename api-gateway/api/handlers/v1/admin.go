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
	"google.golang.org/protobuf/encoding/protojson"
)

// This function can add policy
// @Summary Add policy
// @Description This API uses to add policy
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param policy body models.Policy true "Policy"
// @Success 201 {object} models.Empty
// @Router /v1/admin/add/policy [POST]
func(h *handlerV1) AddPolicy(c *gin.Context){
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	body := models.Policy{}


	err := c.ShouldBindJSON(&body)
	if err != nil{
		fmt.Println(err)
		return
	}
	ok, err := h.enforcer.AddPolicy(body.User, body.Domain, body.Action)
	if err != nil{
		fmt.Println(">>>", err)
	}
	h.enforcer.SavePolicy()
	fmt.Println(ok)

	c.JSON(http.StatusOK, models.Empty{})
}

// This function can delete policy
// @Summary Delete policy
// @Description This API uses to delete policy
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param policy body models.Policy true "Policy"
// @Success 201 {object} models.Empty
// @Router /v1/admin/delete/policy [DELETE]
func (h *handlerV1) RemovePolicy(c *gin.Context){
	body := models.Policy{}

	err := c.ShouldBindJSON(&body)
	if err != nil{
		fmt.Println(err)
	}
	ok, err := h.enforcer.RemovePolicy(body.User, body.Domain, body.Action)
	if err != nil{
		fmt.Println(err)
	}
	h.enforcer.SavePolicy()
	fmt.Println(ok)
	c.JSON(http.StatusOK, models.Empty{})
}

// This function can set role to user
// @Summary Add role to user
// @Description This API uses to give role to user
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param policy body models.Role true "Role"
// @Success 201 {object} models.Empty
// @Router /v1/admin/add/role [POST]
func (h *handlerV1) AddRoleToUser(c *gin.Context){
	body := models.Role{}

	err := c.ShouldBindJSON(&body)
	if err != nil{
		fmt.Println(err)
	}

	ok, err := h.enforcer.AddRoleForUser(body.Id, body.Type)
	if err != nil{
		fmt.Println(">>>", err)
	}
	fmt.Println(ok)

	c.JSON(http.StatusOK, models.Empty{})
}

// This function can remove role from user
// @Summary Remove role from user
// @Description This API uses to remove role from user
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param policy body models.Role true "Role"
// @Success 201 {object} models.Empty
// @Router /v1/admin/remove/role [POST]
func (h *handlerV1) RemoveRoleFromUser(c *gin.Context){
	body := models.Role{}

	err := c.ShouldBindJSON(&body)
	if err != nil{
		fmt.Println(err)
	}
	ok, err := h.enforcer.RemoveGroupingPolicy(body.Id, body.Type)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(ok)

	c.JSON(http.StatusOK, models.Empty{})
}

// Admin             Login
// @Summary         Login Admin
// @Description     This Function get login admin
// @Tags            admin
// @Accept 			json
// @Produce			json
// @Param 			email 		path string true "email"
// @Param 			password 	path string true "password"
// @Succes 			201 {object}    models.LoginResponse
// @Failure			500 {object} 	models.Error
// @Failure			400 {object} 	models.Error
// @Router			/v1/login/admin/{email}/{password} [get]
func (h *handlerV1) LoginFodrAdmins(c *gin.Context) {
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

	h.jwtHandler.Iss = "Admin"
	h.jwtHandler.Sub = res.Id
	h.jwtHandler.Role = "admin"
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