package database

import (
	"context"
	"log"
	"vk/app"
	"vk/model"
)


func (db DB) SignIn(req model.Credentials) (*model.User, error) {

	query := "SELECT * FROM users WHERE username = $1"

	ctx := context.Background()
	row := db.pool.QueryRow(ctx, query, req.Username)
	var u model.User

	if err := row.Scan(&u.ID, &u.Name, &u.Password, &u.Role); err != nil {
		log.Println("[-]err:SignIn:Scan:", err)
		return nil, err
	}

	if req.Password != u.Password {
		log.Println("[-]err:SignIn:req.Password != u.Password:", req.Password, u.Password)
		return nil, app.ErrBadPassword
	}

	log.Println("[+]sussess:SignIn", req.Username)
	return &u, nil
}

func (db DB) SignUp(req model.Credentials) error {

	exist, err := db.UserExist(req.Username)
	if err != nil {
		log.Println("[-]err:SignUp:UserExist:", err)
		return err
	}

	if exist {
		log.Println("[-]err:SignUp:UserConflict:", req.Username)
		return app.ErrUserConflict
	}

	log.Println(req.Username, req.Password)
	query := "INSERT INTO users(username, password, role) VALUES ($1, $2, 'User')"

	ctx := context.Background()
	_, err = db.pool.Exec(ctx, query, req.Username, req.Password)
	if err != nil {
		log.Println("[-]err:SignUp:Exec:", err)
		return err
	}

	log.Println("[+]sussess:SignUp", req.Username)
	return nil
}

func (db DB) SignUpAdministrator(req model.Credentials) error {

	exist, err := db.UserExist(req.Username)
	if err != nil {
		log.Println("[-]err:SignUpAdministrator:UserExist:", err)
		return err
	}

	if exist {
		log.Println("[-]err:SignUpAdministrator:UserConflict:", req.Username)
		return app.ErrUserConflict
	}

	query := "INSERT INTO users(username, password, role) VALUES ($1, $2, 'admin')"

	ctx := context.Background()
	_, err = db.pool.Exec(ctx, query, req.Username, req.Password)
	if err != nil {
		log.Println("[-]err:SignUpAdministrator:Exec:", err)
		return err
	}

	log.Println("[+]sussess:SignUpAdministrator", req.Username)
	return app.ErrCreated
}

func (db DB) User(user string) (*model.User, error) {

	query := "SELECT * FROM users WHERE username = $1"

	ctx := context.Background()
	row := db.pool.QueryRow(ctx, query, user)
	var u model.User

	if err := row.Scan(&u.ID, &u.Name, &u.Password, &u.Role); err != nil {
		log.Println("[-]err:User:Scan:", err)
		return nil, err
	}
	log.Println("[+]sussess:User", user)
	return &u, nil
}

func (db DB) UserExist(user string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = $1"
	ctx := context.Background()

	var count int
	err := db.pool.QueryRow(ctx, query, user).Scan(&count)
	if err != nil {
		log.Println("[-]err:UserExist:QueryRow:", err)
		return false, err
	}
	log.Println("[+]sussess:UserExist", user)
	return count > 0, nil
}

func (db DB) DeleteUser(req model.Credentials) error {

	query := "DELETE * FROM users WHERE username = $1"

	ctx := context.Background()
	_, err := db.pool.Exec(ctx, query, req.Username)

	if err != nil {
		log.Println("[-]err:DeleteUser:Exec", req.Username, err)
	}
	log.Println("[+]sussess:DeleteUser", req.Username)
	return nil
}
