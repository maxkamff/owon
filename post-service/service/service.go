package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	p "gitlab.com/post-service/genproto/post"
	"gitlab.com/post-service/pkg/logger"
	gp "gitlab.com/post-service/service/grpc_client"
	"gitlab.com/post-service/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostService struct {
	storage storage.IStorage
	Logger  logger.Logger
	Client  gp.GrpcClient
}

func NewPostService(db *sqlx.DB, log logger.Logger, client gp.GrpcClient) *PostService {
	return &PostService{
		storage: storage.NewStoragePg(db),
		Logger:  log,
		Client:  client,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *p.CreatePostRequest) (*p.CreatePostResponse, error) {
	res, err := s.storage.Post().CreatePost(ctx, req)
	if err != nil {
		s.Logger.Error("Error insert posts", logger.Any("error insert posts", err))
		return &p.CreatePostResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}

	comment := req.Comments

	if comment != nil{
		uid := res.UserId
		pid := res.Id
		comment.UserId = uid
		comment.PostId = pid
		if err != nil {
			fmt.Println(err)
		}
		return res, nil
	} else {
		return res, nil
	}
}

func (s *PostService) GetPostById(ctx context.Context, req *p.GetPostRequest) (*p.GetPostResponse, error) {
	res, err := s.storage.Post().GetPostById(req)
	if err != nil {
		s.Logger.Error("Error get post", logger.Any("error get post", err))
		return &p.GetPostResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}
	return res, nil
}

// func (s *PostService) GetPostsByUserId(ctx context.Context, req *p.GetUserPostsRequest) (*p.GetUserPostsResponse, error) {
// 	res, err := s.storage.Post().GetPostsByUserId(req)
// 	if err != nil {
// 		s.Logger.Error("Error gett users", logger.Any("Error get users", err))
// 		return &p.GetUserPostsResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
// 	}
// 	return res, nil
// }

func (s *PostService) UpdatePost(ctx context.Context, req *p.UpdatePostRequest) (*p.UpdatePostRespoonse, error) {
	res, err := s.storage.Post().UpdatePost(req)
	if err != nil {
		s.Logger.Error("Error updating users", logger.Any("Error updating users", err))
		return &p.UpdatePostRespoonse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}
	return res, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *p.DeletePostRequest) (*p.Empty, error) {
	res, err := s.storage.Post().DeletePost(req)
	if err != nil {
		s.Logger.Error("Error deleting methods", logger.Any("Error deleting methods", err))
		return &p.Empty{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}
	return res, nil
}
