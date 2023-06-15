package postgres

// import (
// 	"github.com/user-service/config"
// 	pb "github.com/user-service/genproto/user"
// 	"github.com/user-service/pkg/db"
// 	"github.com/user-service/storage/repo"

// 	"testing"

// 	"github.com/stretchr/testify/suite"
// )

// type UserSuiteTest struct {
// 	suite.Suite
// 	CleanUpfunc func()
// 	Repo        repo.UserStorageI
// }

// func (s *UserSuiteTest) SetupSuite() {
// 	pgPool, cleanUp := db.ConnectToDBForSuite(config.Load())
// 	s.Repo = NewUserRepo(pgPool)
// 	s.CleanUpfunc = cleanUp
// }

// func (s *UserSuiteTest) TestUserCrud() {
// 	// CREATE
// 	user := &pb.CreateUserRequest{
// 		Name: "Ahmadjon",
// 		LastName:  "Darxonov",
// 		Email:     "darxonovahmad@gmail.com",
// 	}

// 	createdUserResp, err := s.Repo.CreateUser(user)
// 	s.Nil(err)
// 	s.NotNil(createdUserResp)
// 	s.Equal(user.Name, createdUserResp.Name)
// 	s.Equal(user.LastName, createdUserResp.LastName)
// 	s.Equal(user.Email, createdUserResp.Email)

// 	// GET
// 	getUserResp, err := s.Repo.GetUserById(&pb.GetUserRequest{
// 		Id: createdUserResp.Id,
// 	})
// 	s.Nil(err)
// 	s.NotNil(getUserResp)
// 	s.Equal(createdUserResp.Name, getUserResp.Name)
// 	s.Equal(createdUserResp.LastName, getUserResp.LastName)

// 	// GET ALL
// 	getAllUsersResp, err := s.Repo.GetAllUsers(&pb.Empty{})
// 	s.Nil(err)
// 	s.NotNil(getAllUsersResp)

// 	// SEARCH
// 	searchUserResp, err := s.Repo.SearchUser(&pb.SearchUserRequest{
// 		Name: createdUserResp.Name,
// 	})
// 	s.Nil(err)
// 	s.NotNil(searchUserResp)

// 	// UPDATE
// 	updateUser := &pb.UpdateUserRequest{
// 		Id:        createdUserResp.Id,
// 		Name: "UpdatedName",
// 		LastName:  "UpdatedLastName",
// 		Email:     "updated-email@example.com",
// 	}

// 	UpdateUserResp, err := s.Repo.UpdateUser(updateUser)
// 	s.Nil(err)
// 	s.NotNil(UpdateUserResp)
// 	s.Equal(updateUser.Name, UpdateUserResp.Name)
// 	s.Equal(updateUser.LastName, UpdateUserResp.LastName)
// 	s.Equal(updateUser.Email, UpdateUserResp.Email)

// 	// DELETE
// 	_, err = s.Repo.DeleteUser(&pb.DeleteUserRequest{Id: createdUserResp.Id})
// 	s.Nil(err)

// 	getDeleteResp, err := s.Repo.GetUserById(&pb.GetUserRequest{Id: createdUserResp.Id})
// 	s.NotNil(err)
// 	s.Nil(getDeleteResp)
// }

// func (suite *UserSuiteTest) TearDownSuite() {
// 	suite.CleanUpfunc()
// }

// func TestUserRepositoryTestSuite(t *testing.T) {
// 	suite.Run(t, new(UserSuiteTest))
// }
