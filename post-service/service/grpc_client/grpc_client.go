package grpcclient

import (
	"fmt"

	"gitlab.com/post-service/config"
	c "gitlab.com/post-service/genproto/comment"
	u "gitlab.com/post-service/genproto/user"
	"google.golang.org/grpc"
)

type GrpcClientI interface {
	User() u.UserServiceClient
	Comment() c.CommentServiceClient
}
type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

func New(cfg config.Config) (*GrpcClient, error) {
	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("User Service dial host:%s, port:%s", cfg.UserServiceHost, cfg.UserServicePort)
	}

	connComment, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.CommentServiceHost, cfg.CommentServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Comment Service dial host:%s, port:%s", cfg.CommentServiceHost, cfg.CommentServicePort)
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"user-service":    u.NewUserServiceClient(connUser),
			"comment-service": c.NewCommentServiceClient(connComment),
		},
	}, nil
}

func (s *GrpcClient) User() u.UserServiceClient {
	return s.connections["user-service"].(u.UserServiceClient)
}

func (s *GrpcClient) Commment() c.CommentServiceClient {
	return s.connections["comment-service"].(c.CommentServiceClient)
}
