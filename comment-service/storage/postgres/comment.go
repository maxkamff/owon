package postgres

import (
	"context"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	c "gitlab.com/comment-service/genproto/comment"
)

func (r *CommentRepo) CreateComment(ctx context.Context, com *c.CreateCommentRequest) (*c.CreateCommentResponse, error) {

	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()

	var res c.CreateCommentResponse
	err := r.db.QueryRow(`
	INSERT INTO
		comments(user_id, post_id, description, liked)
	VALUES(
		$1, $2, $3, $4)
	RETURNING
		id, user_id, post_id, description, liked, created_at, updated_at`, com.UserId, com.PostId, com.Description, com.Liked).
		Scan(&res.Id, &res.UserId, &res.PostId, &res.Description, &res.Liked, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		log.Println("error Inserting comment info")
		return &c.CreateCommentResponse{}, err
	}

	return &res, nil
}

func (r *CommentRepo) DeleteComment(com *c.DeleteCommentRequest) (*c.DeleteCommentEmpty, error) {
	_, err := r.db.Exec(
		`UPDATE
			comments
		SET
			deleted_at = $1
		WHERE
			id = $2
		`, time.Now(), com.PostId)
	if err != nil {
		return &c.DeleteCommentEmpty{}, err
	}
	return &c.DeleteCommentEmpty{}, nil
}

func (r *CommentRepo) UpdateComment(com *c.UpdateCommentRequest) (*c.UpdateCommentResponse, error) {
	var res c.UpdateCommentResponse
	err := r.db.QueryRow(`
	UPDATE
		comments
	SET
		user_id = $1, post_id = $2, description = $3, liked = $4
	WHERE
		id = $5
	RETURNING
		id, user_id, post_id, description, liked, created_at, updated_at
	`, com.UserId, com.PostId, com.Description, com.Liked, com.Id).
		Scan(
			&res.Id, &res.UserId, &res.PostId, &res.Description, &res.Liked, &res.CreatedAt, &res.UpdatedAt,
		)
	if err != nil {
		log.Println("Error Updating comment info")
		return &c.UpdateCommentResponse{}, err
	}
	return &res, nil
}

func (r *CommentRepo) GetComment(com *c.GetACommentRequest) (*c.GetACommentResponse, error) {
	var res c.GetACommentResponse
	err := r.db.QueryRow(`
	SELECT 
		id, user_id, post_id, description, liked, created_at, updated_at
	FROM
		comments
	WHERE
		id = $1 and deleted_at is null
	`, com.Id).
		Scan(
			&res.Id, &res.UserId, &res.PostId, &res.Description, &res.Liked, &res.CreatedAt, &res.UpdatedAt,
		)
	if err != nil {
		log.Println("Error getting comment info")
		return &c.GetACommentResponse{}, err
	}
	return &res, nil
}

func (r *CommentRepo) GetAllCommentsByPostId(com *c.GetCommentsPostRequest) (*c.GetCommentsPostResponse, error) {
	var res c.GetCommentsPostResponse
	query := `
	SELECT 
		id, user_id, post_id, description, liked, created_at, updated_at
	FROM 
		comments
	WHERE
		post_id = $1`
	rows, err := r.db.Query(query, com.PostId)
	if err != nil {
		return &c.GetCommentsPostResponse{}, err
	}
	for rows.Next() {
		temp := c.GetACommentResponse{}
		err = rows.Scan(
			&temp.Id, &temp.UserId, &temp.PostId, &temp.Description, &temp.Liked, &temp.CreatedAt, &temp.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		res.Comment = append(res.Comment, &temp)
	}
	return &res, nil
}
