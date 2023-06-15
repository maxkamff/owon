package repo

import (
	"context"

	p "gitlab.com/post-service/genproto/post"
)

type PostStorageI interface {
	CreatePost(context.Context, *p.CreatePostRequest) (*p.CreatePostResponse, error)
	GetPostById(*p.GetPostRequest) (*p.GetPostResponse, error)
	//GetPostsByUserId(*p.GetUserPostsRequest) (*p.GetUserPostsResponse, error)
	UpdatePost(*p.UpdatePostRequest) (*p.UpdatePostRespoonse, error)
	DeletePost(*p.DeletePostRequest) (*p.Empty, error)
}
