package server

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"vk/app"
	"vk/model"
)

type Logger struct {
	handler http.Handler
}

func wrapLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading request body:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(strings.NewReader(string(body)))
	log.Printf("[request] %s %s %s", r.Method, r.URL.Path, body)
	l.handler.ServeHTTP(w, r)
	log.Printf("[metric] %s %s %v", r.Method, r.URL.Path, time.Since(start))
}

type EnsureAuth struct {
	handler http.Handler
}

func GetAuthenticatedUser(r *http.Request) (*model.Session, error) {
	body, _ := io.ReadAll(r.Body)

	c, err := r.Cookie("session_token")
	if err != nil {

		if err == http.ErrNoCookie {
			log.Printf("[-][no cookie] %s %s %s", r.Method, r.URL.Path, body)
			return nil, app.ErrCookies
		}
		log.Printf("[-][auth] %s %s %s", r.Method, r.URL.Path, body)
		return nil, app.ErrBadRequest
	}

	sessionToken := c.Value
	userSession, exists := SessionsCache[sessionToken]

	if !exists {
		log.Printf("[-][cookie][exists] %s %s %s", r.Method, r.URL.Path, body)
		return nil, app.ErrCookies
	}

	if userSession.IsExpired() {
		delete(SessionsCache, sessionToken)
		return nil, app.ErrCookies
	}

	log.Println("succes auth checked token")
	return &userSession, nil
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := GetAuthenticatedUser(r)
	if err != nil {
		log.Println("[session] bad cookie", r.Method, r.URL.Path)
		http.Error(w, "not authorized, please signin", http.StatusUnauthorized)
		return
	}

	if session.Role == "admin" {
		log.Println("[admin] continue admin role", r.Method, r.URL.Path)
		ea.handler.ServeHTTP(w, r)
	}

	if session.Role == "user" {
		if r.Method != "GET" {
			log.Println("[user] not permissioned method", r.Method, r.URL.Path)
			http.Error(w, "no permission", http.StatusForbidden)
			return
		}

		ea.handler.ServeHTTP(w, r)
	}

}

func NewEnsureAuth(handlerToWrap http.Handler) *EnsureAuth {
	return &EnsureAuth{handlerToWrap}
}

type AdminMiddleware struct {
	handler http.Handler
}

func NewAdminMiddleware(handlerToWrap http.Handler) *AdminMiddleware {
	return &AdminMiddleware{handlerToWrap}
}

func (ea *AdminMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "please sign-in", http.StatusUnauthorized)
		return
	}

	if session.Role != "admin" {
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	ea.handler.ServeHTTP(w, r)

}
