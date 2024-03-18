package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"vk/app"
	"vk/model"
	"vk/verify"

	"github.com/google/uuid"
)

func (s APIServer) SignUp(w http.ResponseWriter, r *http.Request) interface{} {

	var cred model.Credentials

	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		return err
	}
	if !verify.Creds(cred) {
		return app.ErrBadPassword
	}
	if err = s.repo.SignUp(cred); err != nil {
		return app.AlreadExist
	}

	return app.SuccessCreated
}

func (s APIServer) SignUpAdmin(w http.ResponseWriter, r *http.Request) interface{} {

	var cred model.Credentials

	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		return err
	}
	log.Println(cred)

	if !verify.Creds(cred) {
		return app.ErrBadPassword
	}

	if err = s.repo.SignUp(cred); err != nil {
		return app.ErrSignUp
	}

	log.Println(cred)

	return nil
}

func (s APIServer) SignIn(w http.ResponseWriter, r *http.Request) interface{} {
	log.Println("income creed", r.Body)
	var creds model.Credentials

	// body, _ := io.ReadAll(r.Body)
	// log.Printf("[ReadAll] %s ", body)

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Println("try sign in decode", err, creds)
		return err
	}

	if !verify.Creds(creds) {
		log.Println("try sign in verify creds", creds)
		return app.ErrBadPassword
	}
	log.Println("try sign in")
	user, err := s.repo.SignIn(creds)

	if err != nil {
		if err == app.ErrBadPassword {
			log.Println("err:app.ErrBadPassword:", err)
			return app.ErrUnauthorized
		}
		if user == nil {
			log.Println("err:user == nil:", err)
			return app.ErrNotFound
		}

		log.Println("err:s.repo.SignIn(creds):", err)
		return app.ErrBadRequest

	}

	sessionToken := uuid.NewString()
	expiredAt := time.Now().Add(15 * time.Minute)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiredAt,
	})

	SessionsCache[sessionToken] = model.Session{
		User:      user.Name,
		Token:     sessionToken,
		Role:      user.Role,
		CreatedAt: time.Now(),
		ExpiresAt: expiredAt,
	}
	log.Println("succes created token", sessionToken)
	return app.SucceesSignIn
}

func (s APIServer) DeleteUser(w http.ResponseWriter, r *http.Request) interface{} {

	var cred model.Credentials

	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		return err
	}
	if !verify.Creds(cred) {
		return app.ErrBadPassword
	}
	if err = s.repo.DeleteUser(cred); err != nil {
		return app.ErrNotFound
	}

	return app.SussessDelete
}
