package v1

import (
	"context"
	"net/http"

	"time"

	_ "gitlab.com/api-gateway/api/handlers/models"
	pb "gitlab.com/api-gateway/genproto/user"
	l "gitlab.com/api-gateway/pkg/logger"

	"github.com/opentracing/opentracing-go"
	"github.com/gin-gonic/gin"
	//"gitlab.com/api-gateway/pkg/utils"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Summary CrateUser
// @Description This API Can Create User
// @Tags user
// @Accept json
// @Produce json
// @Param body body models.CreateUserRequest true "Get Users By LImit Page"
// @Success 201 {object} models.CreateUserResponse
// @Failure 400 string Error response
// @Router /v1/user [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var body pb.CreateUserRequest
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()

	response, err := h.serviceManager.UserService().CreateUser(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list users", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get User
// @Description This API Can Get User Info
// @Tags user
// @Security    BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} models.GetUserResponse
// @Failure 400 string Error response
// @Router /v1/user [get]
func (h *handlerV1) GetUserById(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	claims := GetClaims(h, c)
	id = claims["sub"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetUserById(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorf": err.Error(),
		})
		h.log.Error("failed to get user from db", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary GetUser
// @Description This API Can Get Users By Limit
// @Tags user
// @Accept json
// @Produce json
// @Param body body models.GetAllUserRequest true "Get Users By LImit Page"
// @Success 201 {object} models.GetAllUsersResponse
// @Failure 400 string Error response
// @Router /v1/users [get]
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	var body pb.GetAllUserRequest
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetAllUsers(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list users", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Delete User
// @Description This API Can Delete User
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} models.Empty
// @Failure 400 string Error response
// @Router /v1/user [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	claims := GetClaims(h, c)
	id := claims["sub"].(string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	_, err := h.serviceManager.UserService().DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Update User
// @Description This API Can Update All Informations
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body models.UpdateUserRequest true "UpdateUser"
// @Success 201 {object} models.UpdateUserResponse
// @Failure 400 string Error response
// @Router /v1/user [put]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		body        pb.UpdateUserRequest
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	claims := GetClaims(h, c)
	id := claims["sub"].(string)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Id = id
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().UpdateUser(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}