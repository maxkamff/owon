package postgres

import (
	"context"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	p "gitlab.com/post-service/genproto/post"
)

// This method creates new post
func (r *PostRepo) CreatePost(ctx context.Context, post *p.CreatePostRequest) (*p.CreatePostResponse, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	var res p.CreatePostResponse
	err := r.db.DB.QueryRow(`
	INSERT INTO
		posts(photo, id, title, user_id)
	VALUES(
		$1, $2, $3, $4)
	RETURNING
		id, title, photo, user_id, created_at`, post.Id, post.Title, post.Photo, post.UserId).
		Scan(&res.Id, &res.Title, &res.Photo, &res.UserId, &res.CreatedAt)
	if err != nil {
		log.Println("error while inserting post info to db")
		return &p.CreatePostResponse{}, err
	}

	return &res, nil
}

// This method gets post by id
func (r *PostRepo) GetPostById(post *p.GetPostRequest) (*p.GetPostResponse, error) {
	var res p.GetPostResponse
	err := r.db.DB.QueryRow(`
	SELECT
		id, title, photo, user_id, created_at, updated_at
	FROM 
		posts
	WHERE
		id = $1 and deleted_at is null
	`, post.Id).
		Scan(&res.Id, &res.Title, &res.Photo, &res.UserId, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		log.Println("error while getting post info from db")
		return &p.GetPostResponse{}, err
	}

	return &res, nil
}

// This method gets posts bu user_id
// func (r *PostRepo) GetPostsByUserId(post *p.GetUserPostsRequest) (*p.GetUserPostsResponse, error) {
// 	var res p.GetUserPostsResponse
// 	query := `
// 	SELECT 
// 		id, title, photo, user_id, created_at, updated_at
// 	FROM 
// 		posts
// 	WHERE
// 		user_id = $1 and deleted_at is null
// 	LIMIT $2`
// 	rows, err := r.db.Query(query, post.UserId)
// 	if err != nil {
// 		return &p.GetUserPostsResponse{}, err
// 	}
// 	for rows.Next() {
// 		temp := p.GetUserPostResponse{}
// 		err = rows.Scan(
// 			&temp.Id, &temp.Title, &temp.Photo, &temp.UserId, &temp.CreatedAt, &temp.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		res.Post = append(res.Post, &temp)
// 	}
// 	return &res, nil
// }

// This method update post info
func (r *PostRepo) UpdatePost(post *p.UpdatePostRequest) (*p.UpdatePostRespoonse, error) {
	var res p.UpdatePostRespoonse
	err := r.db.QueryRow(`
	UPDATE 
		posts
	SET
		title = $1
	WHERE
		id = $2 and deleted_at is null
	RETURNING
		id, title, photo, user_id, created_at, updated_at
	`, post.Title, post.Id).
		Scan(&res.Id, &res.Title, &res.Photo, &res.UserId, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		log.Println("Error while updating user info")
		return &p.UpdatePostRespoonse{}, err
	}
	return &res, nil
}

// This method delete post info from db
func (r *PostRepo) DeletePost(com *p.DeletePostRequest) (*p.Empty, error) {
	_, err := r.db.Exec(
		`UPDATE
			posts
		SET
			deleted_at = $1
		WHERE
			id = $2
		`, time.Now(), com.Id)
	if err != nil {
		return &p.Empty{}, err
	}
	return &p.Empty{}, nil
}