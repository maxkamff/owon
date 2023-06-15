package services

import (
	"fmt"

	"gitlab.com/api-gateway/config"
	c "gitlab.com/api-gateway/genproto/comment"
	p "gitlab.com/api-gateway/genproto/post"
	u "gitlab.com/api-gateway/genproto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type IServiceManager interface {
	UserService() u.UserServiceClient
	PostService() p.PostServiceClient
	CommentService() c.CommentServiceClient
}

type serviceManager struct {
	userService    u.UserServiceClient
	commentService c.CommentServiceClient
	postService    p.PostServiceClient
}

func (s *serviceManager) UserService() u.UserServiceClient {
	return s.userService
}

func (s *serviceManager) PostService() p.PostServiceClient {
	return s.postService
}

func (s *serviceManager) CommentService() c.CommentServiceClient {
	return s.commentService
}

func NewServiceManager(conf *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%s", conf.UserServiceHost, conf.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%s", conf.PostServiceHost, conf.PostServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	connComment, err := grpc.Dial(
		fmt.Sprintf("%s:%s", conf.CommentServiceHost, conf.CommentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	serviceManager := &serviceManager{
		userService:    u.NewUserServiceClient(connUser),
		postService:    p.NewPostServiceClient(connPost),
		commentService: c.NewCommentServiceClient(connComment),
	}

	return serviceManager, nil
}
