package v1

import (
	"context"
	"net/http"
	"strconv"

	// "strconv"

	"time"

	//"gitlab.com/api-gateway/api/handlers/models"
	pb "gitlab.com/api-gateway/genproto/post"
	l "gitlab.com/api-gateway/pkg/logger"

	//"gitlab.com/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Summary Create Post
// @Description This API Can Create Post
// @Tags post
// @Security    BearerAuth
// @Accept json
// @Produce json
// @Param body body models.CreatePostRequest true "Create Post"
// @Success 201 {object} models.CreatePostResponse
// @Failure 400 string Error response
// @Router /v1/post [post]
func (h *handlerV1) CreatePost(c *gin.Context) {
	var (
		body        pb.CreatePostRequest
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	claims := GetClaims(h, c)
	user_id := claims["sub"].(string)
	body.UserId = user_id

	id, err := uuid.NewRandom()
	if err != nil{
		panic("Can't generate uuid")
	}
	body.Id = id.String()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.PostService().CreatePost(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create post", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Update Post
// @Description This API Can Update Post
// @Tags post
// @Security    BearerAuth
// @Accept json
// @Produce json
// @Param body body models.UpdatePostRequest true "Update Post"
// @Success 201 {object} models.UpdatedPostResponse
// @Failure 400 string Error response
// @Router /post [put]
func (h *handlerV1) UpdatePost(c *gin.Context) {
	var (
		body        pb.UpdatePostRequest
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	claims := GetClaims(h, c)
	user_id := claims["sub"].(string)
	body.UserId = user_id


	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.PostService().UpdatePost(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update post", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get Post
// @Description This API Can Get Post Info
// @Tags post
// @Security    BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "get Post"
// @Success 201 {object} models.GetPostResponse
// @Failure 400 string Error response
// @Router /post [get]
func (h *handlerV1) GetPostById(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	claims := GetClaims(h, c)
	_ = claims["sub"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.PostService().GetPostById(ctx, &pb.GetPostRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get post by id", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}


// @Summary Get Post
// @Description This API Can Get Post Info
// @Tags post
// @Security    BearerAuth
// @Accept json
// @Produce json
// @Param limit path string true "get post by limit"
// @Param page path string true "get post by page"
// @Success 201 {object} models.GetUserPostsResponse
// @Failure 400 string Error response
// @Router /posts [get]
func (h *handlerV1) GetPostsByUserId(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	limit := c.Param("limit")
	page := c.Param("page")

	lm, err:= strconv.Atoi(limit)
	if err != nil {
		panic(err)
	}
	pg, err := strconv.Atoi(page)
	if err != nil {
		panic(err)
	}

	claims := GetClaims(h, c)
	_  = claims["sub"].(string)

	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.PostService().GetPostsByUserId(ctx, &pb.GetUserPostsRequest{Limit: int64(lm), Page: int64(pg)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get posts", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get Post
// @Description This API Can Get Post Info
// @Tags post
// @Security    BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "delete post"
// @Success 201 {object} models.Empty
// @Failure 400 string Error response
// @Router /post [delete]
func (h *handlerV1) DeletePost(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	claims := GetClaims(h, c)
	_ = claims["sub"].(string)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	_, err := h.serviceManager.PostService().DeletePost(ctx, &pb.DeletePostRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete post", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}
