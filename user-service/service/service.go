package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	// "gitlab.com/user-service/genproto/comment"
	// "gitlab.com/user-service/genproto/post"
	u "gitlab.com/user-service/genproto/user"
	"gitlab.com/user-service/pkg/logger"
	gp "gitlab.com/user-service/service/grpc_client"
	"gitlab.com/user-service/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	storage storage.IStorage
	Logger  logger.Logger
	Client  gp.GrpcClient
}

func NewUserService(db *sqlx.DB, log logger.Logger, client gp.GrpcClient) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		Logger:  log,
		Client:  client,
	}
}


func (s *UserService) CreateUser(ctx context.Context, req *u.CreateUserRequest) (*u.CreateUserResponse, error) {
	user, err := s.storage.User().CreateUser(ctx, req)
	if err != nil {
		s.Logger.Error("Error insert user", logger.Any("Error insert user", err))
		return &u.CreateUserResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}

	post := req.Posts

	if post != nil {
		id := user.Id
		post.UserId = id
		if err != nil {
			fmt.Println(err)
		}
		return user, nil
	} else {
		return user, nil
	}
}

func (s *UserService) UpdateUser(ctx context.Context, req *u.UpdateUserRequest) (*u.UpdateUserRespoonse, error) {
	res, err := s.storage.User().UpdateUser(req)
	if err != nil {
		s.Logger.Error("Error update user", logger.Any("Error update user", err))
		return &u.UpdateUserRespoonse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}
	return res, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *u.GetUserRequest) (*u.GetUserResponse, error) {
	res, err := s.storage.User().GetUserById(req)
	if err != nil {
		s.Logger.Error("Error get user", logger.Any("Error get user", err))
		return &u.GetUserResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}

	return res, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *u.DeleteUserRequest) (*u.Empty, error) {
	res, err := s.storage.User().DeleteUser(req)
	if err != nil {
		s.Logger.Error("Error delete user", logger.Any("Error delete user", err))
		return &u.Empty{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}
	return res, nil
}

func (s *UserService) SearchUser(ctx context.Context, req *u.SearchUserRequest) (*u.SearchUsersResponse, error) {
	res, err := s.storage.User().SearchUser(req)
	if err != nil {
		s.Logger.Error("Error search user", logger.Any("Error search user", err))
		return &u.SearchUsersResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}
	return res, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, req *u.GetAllUserRequest) (*u.GetAllUsersResponse, error) {
	res, err := s.storage.User().GetAllUsers(req)
	if err != nil {
		s.Logger.Error("Error get all users", logger.Any("Error get all user", err))
		return &u.GetAllUsersResponse{}, status.Error(codes.Internal, "Somthing went wrong, please check users info")
	}
	return res, nil
}

func (s *UserService) CheckField(ctx context.Context, req *u.CheckFeildRequest) (*u.CheckFeildResponse, error) {
	boolean, err := s.storage.User().CheckField(req.Feild, req.Value)
	if err != nil {
		s.Logger.Error("Error check user", logger.Any("Error check user", err))
		return &u.CheckFeildResponse{}, status.Error(codes.Internal, "Somthing went wrong, please check users info")
	}
	return &u.CheckFeildResponse{
		Exists: boolean.Exists,
	}, nil
}

func (s *UserService) GetByEmail(ctx context.Context, req *u.Email) (*u.LoginResponse, error){
	res, err := s.storage.User().GetByEmail(req.Email)
	if err != nil {
		s.Logger.Error("Error get user", logger.Any("Error get user", err))
		return &u.LoginResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}
		return res, nil
}

func (s *UserService) UpdateTokens(ctx context.Context, req *u.TokenRequest) (*u.TokenResponse, error){
	res, err := s.storage.User().UpdateTokens(req)
	if err != nil {
		s.Logger.Error("Error update token", logger.Any("Error update token", err))
		return &u.TokenResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}
		return res, nil
}

func (s *UserService) GetUserRole(ctx context.Context, req *u.CheckRoleRequest) (*u.CheckRoleResponse, error){
	res, err := s.storage.User().CheckRole(req)
	if err != nil {
		s.Logger.Error("Error get user", logger.Any("Error get user", err))
		return &u.CheckRoleResponse{}, status.Error(codes.Internal, "Something went wrong, please check users info")
	}
		return res, nil
}