package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"gitlab.com/api-gateway/api/handlers/models"
	"gitlab.com/api-gateway/email"
	pb "gitlab.com/api-gateway/genproto/user"
	"gitlab.com/api-gateway/pkg/etc"
	"gitlab.com/api-gateway/pkg/logger"
	l "gitlab.com/api-gateway/pkg/logger"
	"google.golang.org/protobuf/encoding/protojson"
)

// register customer
// @Summary		register customer
// @Description	this registers customer
// @Tags		register
// @Accept		json
// @Produce 	json
// @Param 		body	body  	 models.CreateRegister true "Register customer"
// @Success		201 	{object} models.Error
// @Failure		500 	{object} models.Error
// @Router		/v1/register 	[post]
func (h *handlerV1) Register(c *gin.Context) {
	var body models.CreateRegister

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
			"Hint":  "Check your data",
		})
		h.log.Error("Error while binding json", l.Any("json", err))
		return
	}
	body.Email = strings.TrimSpace(body.Email)
	body.UserName = strings.TrimSpace(body.Name)
	body.Email = strings.ToLower(body.Email)
	body.UserName = strings.ToLower(body.UserName)
	body.Password, err = etc.GeneratePasswordHash(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error: err,
		})
		h.log.Error("couldn't hash the password")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	emailExists, err := h.serviceManager.UserService().CheckField(ctx, &pb.CheckFeildRequest{
		Feild: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		h.log.Error("Error while cheking email uniqeness", l.Any("check", err))
		return
	}

	if emailExists.Exists {
		c.JSON(http.StatusConflict, gin.H{
			"info": "Email is already used",
		})
		return
	}

	body.Code = etc.GenerateCode(6)
	bodyByte, err := json.Marshal(body)
	if err != nil {
		h.log.Error("Error while marshaling to json", l.Any("json", err))
		return
	}

	msg := "Subject: User email verification\n Your verification code: " + body.Code
	err = email.SendEmail([]string{body.Email}, []byte(msg))

	err = h.redis.SetWithTTL(body.Email, string(bodyByte), 300)
	if err != nil {
		h.log.Error("Error while marshaling to json", l.Any("json", err))
		return
	}
	c.JSON(http.StatusAccepted,body.Code)

}

// register customer
// @Summary		register customer
// @Description	this registers customer
// @Tags		register
// @Accept		json
// @Produce 	json
// @Param 		email path string true "email"
// @Param 		code  path string true "code"
// @Success		201 	{object} models.VerifyResponse
// @Failure		500 	{object} models.Error
// @Router		/v1/verify/{email}/{code} 	[get]
func (h *handlerV1) Verify(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true 

	var (
		email = c.Param("email")
		code  = c.Param("code")
	)

	userBody, err := h.redis.Get(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error get from redis user body", l.Error(err))
		return
	}

	if userBody == nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"info": "Your time has expired",
		})
		return
	}

	userBodys := cast.ToString(userBody)
	body := pb.CreateUserRequest{}

	err = json.Unmarshal([]byte(userBodys), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error while unmarshaling from json to customer body", logger.Any("json", err))
		return
	}
	fmt.Println(">>>>",body.Code)
	if body.Code != code {
		fmt.Println(body.Code)
		c.JSON(http.StatusConflict, gin.H{
			"info": "Wrong code",
		})
		return
	}

	id, err := uuid.NewRandom()
	if err != nil{
		panic("Can't generate uuid")
	}
	body.Id = id.String()
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	h.jwtHandler.Iss = "user"
	h.jwtHandler.Sub = body.Id
	h.jwtHandler.Role = "user"
	h.jwtHandler.Aud = []string{"user-fronted"}
	h.jwtHandler.SigninKey = h.cfg.SignInKey
	h.jwtHandler.Log = h.log
	tokens, err := h.jwtHandler.GenerateAuthJWT()
	accessToken := tokens[0]
	refreshToken := tokens[1]
	if err != nil {
		h.log.Error("error occured while generating tokens")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong,please try again",
		})
		return
	}

	body.RefreshToken = refreshToken

	user, err := h.serviceManager.UserService().CreateUser(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error while creating user to db", l.Error(err))
		return
	}
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	c.JSON(http.StatusAccepted, user)
}

func handleInternalWithMessage(c *gin.Context, logger l.Logger, err error, s string) {
	panic("unimplemented")
}
