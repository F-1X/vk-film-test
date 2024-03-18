package database

import (
	"context"
	"log"
)

func (db DB) InitDb() {

	ctx := context.Background()
	_, err := db.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS actors (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			gender VARCHAR(10) NOT NULL,
			birthday DATE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS films (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			date TIMESTAMP,
			rate NUMERIC
		);

		CREATE TABLE IF NOT EXISTS actor_film (
			actor_id INT REFERENCES actors(id),
			film_id INT REFERENCES films(id),
			PRIMARY KEY (actor_id, film_id)
		);

		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS sessions (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			token  TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at TIMESTAMP,
			expired_at TIMESTAMP
		);

		

	`)

	if err != nil {
		log.Fatal("failed init db",err)
	}

	var adminExists bool
	
	err = db.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = 'admin')").Scan(&adminExists)
	if err != nil {
		log.Fatalf("failed to check admin existence: %v", err)
	}

	if !adminExists {
		log.Println("creating a new admin")
		_, err := db.pool.Exec(ctx, "INSERT INTO users (username, password, role) VALUES ('admin', 'admin123', 'admin')")
		if err != nil {
			log.Fatalf("failed to create admin: %v", err)
		}
	}

	var Gleb bool
	err = db.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = 'Gleb')").Scan(&Gleb)
	if err != nil {
		log.Fatalf("failed to check admin existence: %v", err)
	}

	if !Gleb {
		log.Println("creating a new user")
		_, err := db.pool.Exec(ctx, "INSERT INTO users (username, password, role) VALUES ('Gleb', 'qwerty1234', 'user')")
		if err != nil {
			log.Fatalf("failed to create admin: %v", err)
		}
	}


}
