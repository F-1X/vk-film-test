package server

import (
	"context"
	"log"
	"net/http"
	"time"
	"vk/config"
	"vk/database"
	"vk/model"
)

type Server interface {
	Shutdown(ctx context.Context) error
}

type UserService interface {
	GetActorHandler(http.ResponseWriter, *http.Request) interface{}
}

// AdminService определяет интерфейс для сервиса, доступного администраторам.
type AdminService interface {
	GetActorHandler(http.ResponseWriter, *http.Request) interface{}
	PostActorHandler(http.ResponseWriter, *http.Request) interface{}
	PatchActorHandler(http.ResponseWriter, *http.Request) interface{}
	DeleteActorHandler(http.ResponseWriter, *http.Request) interface{}
}

type APIServer struct {
	cfg          config.Server
	repo         database.DB
	router       *http.ServeMux
	server       *http.Server
	sessionsCache map[string]model.Session
}

func NewAPIServer(cfg config.Server, repo database.DB) *APIServer {

	apiServer := &APIServer{
		cfg:         cfg,
		repo:        repo,
		router:      http.NewServeMux(),
		sessionsCache: *NewSessionsCache(),
	}

	apiServer.router.Handle("/api/v1/", NewEnsureAuth(apiServer.NewAPIRouter()))
	apiServer.router.Handle("/admin/", NewAdminMiddleware(apiServer.NewAdminAuthRouter()))
	apiServer.router.Handle("/", apiServer.NewAuthRouter())

	return apiServer
}

func (s *APIServer) Run() error {
	s.server = &http.Server{
		Addr:    ":" + s.cfg.Port,
		Handler: wrapLogger(s.router),
	}

	log.Println("[!] server started at", time.Now().Format("2006/01/02 15:04:05"), "on port:", s.cfg.Port)
	return http.ListenAndServe(":"+s.cfg.Port, wrapLogger(s.router))
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s APIServer) NewAuthRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			wrapHandleFunc(s.SignUp)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			wrapHandleFunc(s.SignIn)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return r
}

func (s APIServer) NewAdminAuthRouter() *http.ServeMux {
	r := http.NewServeMux()
	// админский хендлер умеет удалять пользователей и регистрировать новых админов
	r.HandleFunc("/admin/signup", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			wrapHandleFunc(s.SignUpAdmin)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.HandleFunc("/admin/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			wrapHandleFunc(s.DeleteUser)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return r
}

func (s APIServer) NewAPIRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/api/v1/actors", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			wrapHandleFunc(s.GetActorsHandler)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.HandleFunc("/api/v1/actor", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			wrapHandleFunc(s.GetActorHandler)(w, r)
		case http.MethodPost:
			wrapHandleFunc(s.PostActorHandler)(w, r)
		case http.MethodPatch:
			wrapHandleFunc(s.PatchActorHandler)(w, r)
		case http.MethodDelete:
			wrapHandleFunc(s.DeleteActorHandler)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.HandleFunc("/api/v1/film", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			wrapHandleFunc(s.GetFilm)(w, r)
		case http.MethodPost:
			wrapHandleFunc(s.PostFilm)(w, r)
		case http.MethodDelete:
			wrapHandleFunc(s.DeleteFilm)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.HandleFunc("/api/v1/films", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			wrapHandleFunc(s.GetFilms)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return r
}
