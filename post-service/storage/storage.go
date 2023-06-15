package storage

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/post-service/storage/postgres"
	"gitlab.com/post-service/storage/repo"
)

type IStorage interface {
	Post() repo.PostStorageI
}

type storagePg struct {
	db       *sqlx.DB
	PostRepo repo.PostStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		PostRepo: postgres.NewPostRepo(db),
	}
}

func (s storagePg) Post() repo.PostStorageI {
	return s.PostRepo
}
