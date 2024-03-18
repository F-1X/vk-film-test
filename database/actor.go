package database

import (
	"context"
	"errors"
	"log"
	"strconv"
	"vk/app"
	"vk/model"

	"github.com/jackc/pgx/v5"
)

func (db DB) PostActor(req model.Actor) error {

	actor, err := db.GetActor(req)
	if err != nil && !errors.Is(err, app.ErrActorNotFound) {
		log.Println("err:PostActor:GetActor:", err)
		return err
	}

	if actor != nil {
		log.Println("err:PostActor:ActorConflict:", actor)
		return app.ErrActorConflict
	}

	query := "INSERT INTO actors (name, gender, birthday) VALUES ($1, $2, $3)"

	ctx := context.Background()
	_, err = db.pool.Exec(ctx, query, req.Name, req.Gender, req.Birthday)
	if err != nil {
		log.Println("err:PostActor:Exec:", err)
		return err
	}

	log.Println("sussess:PostActor", req)
	log.Println(query)
	log.Println(req)
	return nil
}

func (db DB) GetActor(req model.Actor) (*model.Actor, error) {
	var actor model.Actor
	query := "SELECT name, gender, birthday FROM actors WHERE 1=1"
	values := []interface{}{}
	if req.Name != "" {
		query += " AND name = $1"
		values = append(values, req.Name)
	}
	if req.Gender != "" {
		query += " AND gender = $2"
		values = append(values, req.Gender)
	}
	if !req.Birthday.IsZero() {
		query += " AND birthday = $3"
		values = append(values, req.Birthday)
	}
	query += " LIMIT 1"
	ctx := context.Background()
	row := db.pool.QueryRow(ctx, query, values...)
	err := row.Scan(&actor.Name, &actor.Gender, &actor.Birthday)
	if err != nil {
		log.Println("err:GetActor:", err)
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, app.ErrActorNotFound
		}
		return nil, err
	}

	log.Println("sussess:GetActor", actor)
	log.Println(query)
	log.Println(values...)
	return &actor, nil
}

func (db DB) UpdateActor(req model.ActorUpdate) error {
	var actor model.Actor

	actor.Name = req.Name
	actor.Gender = req.Gender
	actor.Birthday = req.Birthday

	_, err := db.GetActor(actor)
	if err != nil {
		log.Println("err:UpdateActor:GetActor:", err)
		return err
	}

	query := "UPDATE actors SET"
	values := []interface{}{}
	var i int = 1

	if req.Update.Gender != "" {
		query += " gender = $" + strconv.Itoa(i)
		values = append(values, req.Update.Gender)
		i++
	}

	if !req.Update.Birthday.IsZero() {
		if i != 1 {
			query += ", birthday = $" + strconv.Itoa(i)
		} else {
			query += " birthday = $" + strconv.Itoa(i)
		}
		values = append(values, req.Update.Birthday)
		i++
	}

	if req.Update.Name != "" {
		if i != 1 {
			query += ", name = $" + strconv.Itoa(i)
		} else {
			query += " name = $" + strconv.Itoa(i)
		}
		values = append(values, req.Update.Name)
		i++
	}

	query += " WHERE name = $" + strconv.Itoa(i)
	values = append(values, req.Name)

	if i == 1 {
		log.Println("err:UpdateActor:i == 1:", err)
		return app.ErrBadRequest
	}

	ctx := context.Background()
	_, err = db.pool.Exec(ctx, query, values...)
	if err != nil {
		log.Println("err:UpdateActor:Exec:", err)
		return err
	}
	log.Println("sussess query UpdateActor:", req.Name)
	log.Println(query)
	log.Println(values...)
	return nil
}

func (db DB) DeleteActor(req model.Actor) error {
	_, err := db.GetActor(req)
	if err != nil && errors.Is(err, app.ErrActorNotFound) {
		log.Println("err:DeleteActor:GetActor:", err)
		return err
	}
	if errors.Is(err, app.ErrActorNotFound) {
		log.Println("err:DeleteActor:ErrActorNotFound:", err)
		return app.ErrActorNotFound
	}

	query := "DELETE FROM actors WHERE 1=1"
	values := []interface{}{}
	var i int = 1

	if req.Name != "" {
		query += " AND name = $1"
		values = append(values, req.Name)
		i++
	}
	if req.Gender != "" {
		query += " AND gender = $2"
		values = append(values, req.Gender)
	}
	if !req.Birthday.IsZero() {
		query += " AND birthday = $3"
		values = append(values, req.Birthday)
	}

	ctx := context.Background()
	_, err = db.pool.Exec(ctx, query, values...)
	if err != nil {
		log.Println("err:DeleteActor:Exec:", err)
		return err
	}
	log.Println("sussess query DeleteActor:", req.Name)
	log.Println(query)
	log.Println(values...)
	return nil
}

func (db DB) GetActors(req model.Actor, offset string) ([]model.Actor, error) {

	var actors []model.Actor

	query := "SELECT name, gender, birthday FROM actors WHERE 1=1"

	values := make([]interface{}, 0)
	var i int = 1

	if req.Name != "" {
		query += " AND name ILIKE '%' || $" + strconv.Itoa(i) + " || '%'"
		values = append(values, req.Name)
		i++
	}

	if req.Gender != "" {
		query += " AND gender = LOWER($2)"
		values = append(values, req.Gender)
		i++
	}

	if !req.Birthday.IsZero() {
		query += " AND birthday LIKE $3"
		values = append(values, req.Birthday)
		i++
	}

	query += " LIMIT " + offset

	ctx := context.Background()
	rows, err := db.pool.Query(ctx, query, values...)
	if err != nil {
		log.Println("err:GetActors:Query:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var actor model.Actor
		err = rows.Scan(&actor.Name, &actor.Gender, &actor.Birthday)
		if err != nil {
			log.Println("err:GetActors:Scan:", err)
			return nil, err
		}
		actors = append(actors, actor)

	}
	if err := rows.Err(); err != nil {
		log.Println("err:GetActors:Rows:", err)
		return nil, err
	}

	if len(actors) == 0 {
		log.Println("err:GetActors:len(actors):", app.ErrActorNotFound)
		return nil, app.ErrActorNotFound
	}

	log.Println("sussess:GetActors")
	log.Println(query)
	log.Println(values...)
	return actors, nil
}
