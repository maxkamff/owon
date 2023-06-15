package models

type UpdateUserRequest struct {
	Name      string `json:"name"`
	Last_name string `json:"last_name"`
	Email     string `json:"email"`
	UserName string `json:"username"`
}

type UpdateUserResponse struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	UserName string `json:"username"`
	Refresh_token string `json:"refresh_token"`
}

type DeleteUser struct {
	Id int64 `json:"id"`
}

type Empty struct {
}

type GetAllUsersResponse struct {
	User []*GetAllUserResponse `json:"users"`
}

type GetAllUserResponse struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	UserName string `json:"username"`
	Refresh_token string `json:"refresh_token"`
}

type GetUserRequest struct {
	Id int64 `json:"id"`
}

type GetUserResponse struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	UserName string `json:"username"`
	Refresh_token string `json:"refresh_token"`
}

type CreateRegister struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Last_name string `json:"last_name"`
	Email     string `json:"email"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Code      string `json:"code"`
}

type GetAllUserRequest struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type VerifyResponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Last_name     string `json:"last_name"`
	Email         string `json:"email"`
	Created_at    string `json:"created_at"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Refresh_token string `json:"refresh_token"`
	Access_token  string `json:"access_token"`
}

type LoginReponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Last_name     string `json:"last_name"`
	Email         string `json:"email"`
	Created_at    string `json:"created_at"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Refresh_token string `json:"refresh_token"`
	Access_token  string `json:"access_token"`
}

type GetProfileByJwtRequestModel struct {
	Token string `header:"Authorization"`
}

type CreateCommentRequest struct{
	User_id string `json:"user_id"`
	Post_id string `json:"post_id"`
	Description string `json:"description"`
	Liked bool `json:"liked"`
}

type CreateCommentResponse struct{
	User_id string `json:"user_id"`
	Post_id string `json:"post_id"`
	Description string `json:"description"`
	Liked bool `json:"liked"`
	Id string `json:"id"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type CreateUserRequest struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Last_name string `json:"last_name"`
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Refresh_token string `json:"refresh_token"`
	Code string `json:"code"`
	Post *CreatePostRequest `json:"post"`
}

type CreateUserResponse struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Last_name string `json:"last_name"`
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Refresh_token string `json:"refresh_token"`
	Code string `json:"code"`
	Access_token string `json:"access_token"`
	Post *CreatePostResponse `json:"post"`
}