package models

type UpdatePostRequest struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Photo string `json:"photo"`
}

type UpdatedPostResponse struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Photo      string `json:"photo"`
	User_Id    string `json:"user_id"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}
type CreatePostRequest struct {
	Title   string               `json:"title"`
	Photo   string               `json:"photo"`
	User_id string               `json:"user_id"`
	Id      string               `json:"id"`
	Comment CreateCommentRequest `json:"comment"`
}

type CreatePostResponse struct {
	Id         string                 `json:"id"`
	Title      string                 `json:"title"`
	Photo      string                 `json:"photo"`
	User_Id    string                 `json:"user_id"`
	Created_at string                 `json:"created_at"`
	Comment    *CreateCommentResponse `json:"comment"`
}

type GetPostResponse struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Photo      string `json:"photo"`
	User_Id    string `json:"user_id"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type GetUserPostResponse struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Photo      string `json:"photo"`
	User_Id    string `json:"user_id"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type GetUserPostsResponse struct {
	Post []*GetUserPostResponse `json:"posts"`
}
