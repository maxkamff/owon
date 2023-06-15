package grpcclient

import (
	"fmt"

	"gitlab.com/comment-service/config"
	p "gitlab.com/comment-service/genproto/post"
	u "gitlab.com/comment-service/genproto/user"
	"google.golang.org/grpc"
)

type GrpcClientI interface {
	Post() p.PostServiceClient
	User() u.UserServiceClient
}
type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

func New(cfg config.Config) (*GrpcClient, error) {
	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Post Service dial host:%s, port:%s", cfg.PostServiceHost, cfg.PostServicePort)
	}

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("User Service dial host:%s, port:%s", cfg.UserServiceHost, cfg.UserServicePort)
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"post-service": p.NewPostServiceClient(connPost),
			"user-service": u.NewUserServiceClient(connUser),
		},
	}, nil
}
func (s *GrpcClient) Post() p.PostServiceClient {
	return s.connections["post-service"].(p.PostServiceClient)
}

func (s *GrpcClient) User() u.UserServiceClient {
	return s.connections["user-service"].(u.UserServiceClient)
}
