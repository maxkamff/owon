package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	u "gitlab.com/user-service/genproto/user"
)

var ()

// This Method saves a user to db
func (r *UserRepo) CreateUser(ctx context.Context ,user *u.CreateUserRequest) (*u.CreateUserResponse, error) {

	trace, ctx := opentracing.StartSpanFromContext(ctx, "CreateUser")
	defer trace.Finish()

	var res u.CreateUserResponse

	err := r.db.QueryRow(`
	INSERT INTO 
		users(id, name, last_name, email, username, password, refresh_token) 
	VALUES(
		$1, $2, $3, $4, $5, $6, $7)
	RETURNING 
		id, 
		name, 
		last_name,
		email,
		username, 
		refresh_token, 
		created_at`,
		user.Id, 
		user.Name, 
		user.LastName, 
		user.Email, 
		user.Username,
		user.Password,
		user.RefreshToken).
		Scan(
			&res.Id, 
			&res.Name, 
			&res.LastName, 
			&res.Email,  
			&res.Username, 
			&res.RefreshToken, 
			&res.CreatedAt)
	if err != nil {
		log.Println("Error while inserting user info to db")
		return &u.CreateUserResponse{}, err
	}
	return &res, nil
}

// This Method gets a user from db
func (r *UserRepo) GetUserById(user *u.GetUserRequest) (*u.GetUserResponse, error) {
	var res u.GetUserResponse
	err := r.db.QueryRow(`
	SELECT id, name, last_name, email, username, refresh_token, created_at, updated_at
	FROM
		users
	WHERE
		id = $1 and deleted_at is null
	`, user.Id).
		Scan(&res.Id, &res.Name, &res.LastName, &res.Email, &res.Username, &res.RefreshToken, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		log.Println("Error while getting user info from db")
		return &u.GetUserResponse{}, err
	}
	return &res, nil
}

// This Method updates user informations in db
func (r *UserRepo) UpdateUser(user *u.UpdateUserRequest) (*u.UpdateUserRespoonse, error) {
	var res u.UpdateUserRespoonse
	err := r.db.QueryRow(`
	UPDATE 
		users
	SET
		name = $1, last_name = $2, email = $3, updated_at = $4, username = $5
	WHERE
		id = $6 and deleted_at is null
	RETURNING
		id, name, last_name, email, created_at, updated_at, refresh_token, username
	`, user.Name, user.LastName, user.Email, time.Now(), user.Username, user.Id).
		Scan(&res.Id, &res.Name, &res.LastName, &res.Email, &res.CreatedAt, &res.UpdatedAt, &res.RefreshToken, &res.Username)
	if err != nil {
		log.Println("Error while Updating user info in db")
		return &u.UpdateUserRespoonse{}, err
	}
	return &res, nil
}

// This Method deletes a user from db
func (r *UserRepo) DeleteUser(user *u.DeleteUserRequest) (*u.Empty, error) {
	_, err := r.db.Exec(
		`UPDATE
			users
		SET
			deleted_at = $1
		WHERE
			id = $2
		`, time.Now(), user.Id)
	if err != nil {
		return &u.Empty{}, err
	}
	return &u.Empty{}, nil
}

// This Method gets users from db by name 
func (r *UserRepo) SearchUser(user *u.SearchUserRequest) (*u.SearchUsersResponse, error) {
	var res u.SearchUsersResponse
	query := `
	SELECT id, name, last_name, email, created_at, updated_at, username, refresh_token
	FROM
		users
	WHERE
		name LIKE $1 and delete_at is null
	`
	rows, err := r.db.Query(query, "%"+user.Name+"%")
	if err != nil {
		log.Println("Error while finding user info from db")
		return &u.SearchUsersResponse{}, err
	}

	for rows.Next() {
		temp := u.SearchUserResponse{}
		err = rows.Scan(
			&temp.Id, &temp.Name, &temp.LastName, &temp.Email, &temp.CreatedAt, &temp.UpdatedAt, &temp.Username, &temp.RefreshToken,
		)

		if err != nil {
			return &u.SearchUsersResponse{}, err
		}
		res.User = append(res.User, &temp)
	}
	return &res, nil
}

// This Method gets users from db by limit and page
func (r *UserRepo) GetAllUsers(user *u.GetAllUserRequest) (*u.GetAllUsersResponse, error) {
	var res u.GetAllUsersResponse
	query := `
	SELECT 
		id, name, last_name, email, created_at, updated_at, username
	FROM 	
		users
	WHERE deleted_at is null
		LIMIT $1
	`
	rows, err := r.db.Query(query, user.Limit)
	if err != nil {
		return &u.GetAllUsersResponse{}, err
	}
	for rows.Next() {
		temp := u.GetAllUserResponse{}
		err = rows.Scan(
			&temp.Id, &temp.Name, &temp.LastName, &temp.Email, &temp.CreatedAt, &temp.UpdatedAt, &temp.Username,
		)
		if err != nil {
			return &u.GetAllUsersResponse{}, err
		}
		res.User = append(res.User, &temp)
	}
	return &res, nil
}

// This Method checks user email or username is uniqe
func (r *UserRepo) CheckField(feild, value string) (*u.CheckFeildResponse, error) {
	query := fmt.Sprintf("SELECT 1 FROM users WHERE %s = $1", feild)
	var exists int
	err := r.db.QueryRow(query, value).Scan(exists)
	if err == sql.ErrNoRows {
		return &u.CheckFeildResponse{Exists: false}, nil
	}

	if err != nil {
		return &u.CheckFeildResponse{}, err
	}

	if exists == 0 {
		return &u.CheckFeildResponse{Exists: false}, nil
	}
	return &u.CheckFeildResponse{Exists: true}, nil
}

// This Method gets a user informations by user name
func (r *UserRepo) GetByEmail(email string) (*u.LoginResponse, error){
	var res u.LoginResponse
	err := r.db.QueryRow(`
	SELECT id, name, last_name, email, username, password, refresh_token, created_at, updated_at
	FROM
		users
	WHERE
		email =  $1 and deleted_at is null
	`, email).
		Scan(&res.Id, &res.Name, &res.LastName, &res.Email, &res.Username, &res.Password, &res.RefreshToken, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		log.Println("Error while getting user info by name")
		return &u.LoginResponse{}, err
	}
	fmt.Println(res)
	return &res, nil
}

// This Method updates user Tokens in db
func (r *UserRepo) UpdateTokens(user *u.TokenRequest) (*u.TokenResponse, error){
	var res u.TokenResponse
	err := r.db.QueryRow(`
	UPDATE
		users
	SET 
		refresh_token = $1, updated_at = $2
	WHERE
		id = $3 and deleted_at is null
	RETURNING refresh_token, updated_at
	`, user.RefreshToken, time.Now(), user.Id).
	Scan(
		&res.RefreshToken,
		&res.UpdatedAt,
	)
	if err != nil {
		log.Println("Error while Updating tokens")
		return &u.TokenResponse{}, err
	}
	return &res, nil
}

// This method checks user role
func (r *UserRepo) CheckRole(user *u.CheckRoleRequest) (*u.CheckRoleResponse, error){
	var res u.CheckRoleResponse
	err := r.db.QueryRow(`
	SELECT v0, v1
	FROM
		casbin_rule
	WHERE 	
		id = v0
	LIMIT 1
	`, user.Id).
	Scan(&res.Id, &res.Role)
	if err != nil {
		log.Println("Error while getting user info for check")
		return &u.CheckRoleResponse{}, err
	}
	return &res, nil
}