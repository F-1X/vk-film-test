package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"vk/config"
	"vk/model"

	_ "database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ActorRepository interface {
	PostActor(req model.Actor) error
	GetActor(req model.Actor) (*model.Actor, error)
	UpdateActor(req model.ActorUpdate) error
}

type FilmRepository interface {
	GetFilm(req model.Film) (*model.Film, error)
	GetFilms(req model.Film) []model.Film
	CreateFilm(req model.Film) error
	DeleteFilm(req model.Film) error
}

type AuthRepository interface {
	SignIn(req model.Credentials) (*model.User, error)
	SignUp(req model.Credentials) error
	SignUpAdministrator(req model.Credentials) error
	User(req model.Credentials) (*model.User, error)
	UserExist(req model.Credentials) (bool, error)
	DeleteUser(req model.Credentials) error
}

type SessionsRepository interface {
	CreateSession(username string, token string, role string, createdAt time.Time, expiresAt time.Time) error
	GetSessionByUsername(username string) (*model.Session, error)
	DeleteSession(sessionToken string) error
}

type DB struct {
	pool *pgxpool.Pool
	cfg  config.Database
}

func New(cfg config.Database) (*DB, error) {

	DATABASE_URL := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)

	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, DATABASE_URL)
	if err != nil {
		log.Fatal("err:", err)
		return nil, err
	}

	if err := dbpool.Ping(ctx); err != nil {
		log.Fatal("err:", err)
		return nil, err
	}

	return &DB{
		pool: dbpool,
		cfg:  cfg,
	}, nil

}
