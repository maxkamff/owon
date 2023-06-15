package storage

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/user-service/storage/postgres"
	"gitlab.com/user-service/storage/repo"
)

type IStorage interface {
	User() repo.UserStorageI
}

type storagePg struct {
	db       *sqlx.DB
	UserRepo repo.UserStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		UserRepo: postgres.NewUserRepo(db),
	}
}

func (s storagePg) User() repo.UserStorageI {
	return s.UserRepo
}
