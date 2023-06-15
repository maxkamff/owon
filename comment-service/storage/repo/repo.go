package repo

import (
	"context"

	c "gitlab.com/comment-service/genproto/comment"
)

type CommentStorageI interface {
	CreateComment(context.Context, *c.CreateCommentRequest) (*c.CreateCommentResponse, error)
	DeleteComment(*c.DeleteCommentRequest) (*c.DeleteCommentEmpty, error)
	UpdateComment(*c.UpdateCommentRequest) (*c.UpdateCommentResponse, error)
	GetComment(*c.GetACommentRequest) (*c.GetACommentResponse, error)
	GetAllCommentsByPostId(*c.GetCommentsPostRequest) (*c.GetCommentsPostResponse, error)
}
