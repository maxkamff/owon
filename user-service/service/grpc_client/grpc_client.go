package grpcclient

import (
	"fmt"

	"gitlab.com/user-service/config"
	c "gitlab.com/user-service/genproto/comment"
	p "gitlab.com/user-service/genproto/post"
	"google.golang.org/grpc"
)

type GrpcClientI interface {
	Post() p.PostServiceClient
	Comment() c.CommentServiceClient
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

	connComment, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.CommentServiceHost, cfg.CommentServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Comment Service dial host:%s, port:%s", cfg.CommentServiceHost, cfg.CommentServicePort)
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"post-service":    p.NewPostServiceClient(connPost),
			"comment-service": c.NewCommentServiceClient(connComment),
		},
	}, nil
}
func (s *GrpcClient) Post() p.PostServiceClient {
	return s.connections["post-service"].(p.PostServiceClient)
}

func (s *GrpcClient) Comment() c.CommentServiceClient {
	return s.connections["comment-service"].(c.CommentServiceClient)
}
