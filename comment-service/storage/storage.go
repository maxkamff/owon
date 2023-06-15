package storage

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/comment-service/storage/postgres"
	"gitlab.com/comment-service/storage/repo"
)

type IStorage interface {
	Comment() repo.CommentStorageI
}

type storagePg struct {
	db          *sqlx.DB
	CommentRepo repo.CommentStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:          db,
		CommentRepo: postgres.NewCommentRepo(db),
	}
}

func (s storagePg) Comment() repo.CommentStorageI {
	return s.CommentRepo
}
