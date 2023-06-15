package repo

import (
	"context"

	u "gitlab.com/user-service/genproto/user"
)

type UserStorageI interface {
	CreateUser(context.Context, *u.CreateUserRequest) (*u.CreateUserResponse, error)
	GetUserById(*u.GetUserRequest) (*u.GetUserResponse, error)
	DeleteUser(*u.DeleteUserRequest) (*u.Empty, error)
	SearchUser(*u.SearchUserRequest) (*u.SearchUsersResponse, error)
	GetAllUsers(*u.GetAllUserRequest) (*u.GetAllUsersResponse, error)
	UpdateUser(*u.UpdateUserRequest) (*u.UpdateUserRespoonse, error)
	CheckField(feild, value string) (*u.CheckFeildResponse, error)
	GetByEmail(email string) (*u.LoginResponse, error)
	UpdateTokens(*u.TokenRequest) (*u.TokenResponse, error)
	CheckRole(*u.CheckRoleRequest) (*u.CheckRoleResponse, error)
}
