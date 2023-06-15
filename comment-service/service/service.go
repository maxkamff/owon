package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	c "gitlab.com/comment-service/genproto/comment"
	"gitlab.com/comment-service/pkg/logger"
	gp "gitlab.com/comment-service/service/grpc_client"
	"gitlab.com/comment-service/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommentService struct {
	storage storage.IStorage
	Logger  logger.Logger
	Client  gp.GrpcClient
}

func NewCommentService(db *sqlx.DB, log logger.Logger, client gp.GrpcClient) *CommentService {
	return &CommentService{
		storage: storage.NewStoragePg(db),
		Logger:  log,
		Client:  client,
	}
}

func (s *CommentService) CreateComment(ctx context.Context, req *c.CreateCommentRequest) (*c.CreateCommentResponse, error) {
	res, err := s.storage.Comment().CreateComment(ctx, req)
	if err != nil {
		s.Logger.Error("Error insert comments", logger.Any("error insert comments", err))
		return &c.CreateCommentResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}
	return res, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, req *c.DeleteCommentRequest) (*c.DeleteCommentEmpty, error) {
	res, err := s.storage.Comment().DeleteComment(req)
	if err != nil {
		s.Logger.Error("Error deleting methods", logger.Any("Error deleting methods", err))
		return &c.DeleteCommentEmpty{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}
	return res, nil
}

func (s *CommentService) UpdateComment(ctx context.Context, req *c.UpdateCommentRequest) (*c.UpdateCommentResponse, error) {
	res, err := s.storage.Comment().UpdateComment(req)
	if err != nil {
		s.Logger.Error("Error updating comments", logger.Any("Error updating comments", err))
		return &c.UpdateCommentResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}
	return res, nil
}

func (s *CommentService) GetComment(ctx context.Context, req *c.GetACommentRequest) (*c.GetACommentResponse, error) {
	res, err := s.storage.Comment().GetComment(req)
	if err != nil {
		s.Logger.Error("Error getting comments", logger.Any("Error getting comments", err))
		return &c.GetACommentResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}

	return res, nil
}

func (s *CommentService) GetAllCommentsByPostId(ctx context.Context, req *c.GetCommentsPostRequest) (*c.GetCommentsPostResponse, error) {
	res, err := s.storage.Comment().GetAllCommentsByPostId(req)
	if err != nil {
		s.Logger.Error("Error getting comments", logger.Any("Error getting comments", err))
		return &c.GetCommentsPostResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}

	return res, nil
}
