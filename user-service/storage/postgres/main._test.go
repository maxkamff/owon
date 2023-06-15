package postgres

import (
	"log"
	"os"
	"testing"

	"gitlab.com/user-service/config"
	"gitlab.com/user-service/pkg/db"
	"gitlab.com/user-service/pkg/logger"
)

var pgRepo *UserRepo

func TestMain(m *testing.M) {
	conf := config.Load()
	connDb, err := db.ConnectToDB(conf)
	if err != nil {
		log.Fatal("Error while connecting to db", logger.Error(err))
	}

	pgRepo = NewUserRepo(connDb)
	os.Exit(m.Run())
}
