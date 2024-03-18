package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
	"vk/app"
	"vk/model"

	"github.com/jackc/pgx"
)

func (db DB) CreateFilm(film *model.Film) error {

	existingFilm, err := db.GetFilmByFilmName(film)
	if err != nil {
		if !errors.Is(err, app.ErrFilmNotFound) {
			log.Println("err:CreateFilm:DB:GetFilm:", err)
			return err

		}
	}

	if existingFilm == nil {
		query := `INSERT INTO films (name, description, date, rate) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING id`
		err := db.pool.QueryRow(context.Background(), query, film.Name, film.Desc, film.Date.Format("2006-01-02"), film.Rate).Scan(&film.ID)
		if err != nil {
			log.Println("err:CreateFilm:DB:QueryRow:", err)
			return err
		}
	} else {
		return app.ErrFilmConflict
	}

	for _, actorName := range film.Actors {
		var actorID int
		err := db.pool.QueryRow(context.Background(), `SELECT id FROM actors WHERE name = $1`, actorName.Name).Scan(&actorID)
		if err != nil {
			log.Println("err:CreateFilm:DB:QueryRow:", err)
			return err
		}

		_, err = db.pool.Exec(context.Background(), `INSERT INTO actor_film (actor_id, film_id) VALUES ($1, $2)`, actorID, film.ID)
		if err != nil {
			log.Println("err:CreateFilm:DB:Exec:", err)
			return err
		}
	}

	log.Println("success:CreateFilm", film)
	return app.ErrCreated
}

func (db DB) GetFilmByFilmName(film *model.Film) (*model.Film, error) {
	log.Println("starting GetFilm")
	query := `
		SELECT f.id, f.name, f.description, f.date, f.rate, array_agg(a.name) as actors
		FROM films f
		LEFT JOIN actor_film af ON f.id = af.film_id
		LEFT JOIN actors a ON af.actor_id = a.id
		WHERE f.name=$1
		GROUP BY f.id, f.name, f.description, f.date, f.rate;
	`

	row := db.pool.QueryRow(context.Background(), query, film.Name)
	var foundFilm model.Film
	var actors []string
	var date sql.NullTime
	if err := row.Scan(&foundFilm.ID, &foundFilm.Name, &foundFilm.Desc, &date, &foundFilm.Rate, &actors); err != nil {

		if err.Error() == pgx.ErrNoRows.Error() {
			log.Println("err:GetFilm:NoRows:", err)
			return nil, app.ErrFilmNotFound
		}

		log.Println("err:GetFilm:Scan:")
		return nil, err
	}

	log.Println(actors)
	if date.Valid {
		foundFilm.Date = model.NewCustomDate(date.Time)
	}

	foundFilm.Actors = []model.Actor{}

	for _, actorName := range actors {

		actor := model.Actor{Name: actorName}
		foundFilm.Actors = append(foundFilm.Actors, actor)
	}

	log.Println("success:GetFilm", foundFilm)
	return &foundFilm, nil
}

func (db DB) GetFilmsByActor(actor *model.Actor) ([]model.Film, error) {
	log.Println("starting GetFilmsByActor")
	query := `
        SELECT f.id, f.name, f.description, f.date, f.rate, a.name as actor_name
        FROM films f
        INNER JOIN actor_film af ON f.id = af.film_id
        INNER JOIN actors a ON af.actor_id = a.id
        WHERE a.name ILIKE %($1)%
    `

	log.Println("GetFilmsByActor:query,values", query, actor.Name)

	rows, err := db.pool.Query(context.Background(), query, actor.Name)
	if err != nil {
		log.Println("err:GetFilmsByActor:Query:", err)
		return nil, err
	}
	defer rows.Close()

	films := make([]model.Film, 0)
	actorSet := make(map[int]bool)

	for rows.Next() {
		var film model.Film
		var date sql.NullTime
		var actorName string
		err := rows.Scan(&film.ID, &film.Name, &film.Desc, &date, &film.Rate, &actorName)
		if err != nil {
			log.Println("err:GetFilmsByActor:Scan:", err)
			return nil, err
		}

		if date.Valid {
			film.Date = model.NewCustomDate(date.Time)
		}

		if !actorSet[film.ID] {
			films = append(films, film)
			actorSet[film.ID] = true
		}
	}

	if err := rows.Err(); err != nil {
		log.Println("err:GetFilmsByActor:Rows:", err)
		return nil, err
	}

	log.Println("success:GetFilmsByActor", len(films), "films found for actor", actor)
	return films, nil
}

// func (db DB) GetFilm(film *model.Film) (*model.Film, error) {

// 	query := "SELECT * FROM films WHERE 1=1"

// 	values := make([]interface{}, 0)
// 	var i int = 1

// 	if film.Name != "" {
// 		query += " AND name ILIKE '%' || $" + strconv.Itoa(i) + " || '%'"
// 		values = append(values, film.Name)
// 		i++
// 	}

// 	if film.Desc != "" {
// 		query += " AND description ILIKE $" + strconv.Itoa(i)
// 		values = append(values, film.Desc)
// 		i++
// 	}

// 	if !film.Date.IsZero() {
// 		log.Println(film.Date)
// 		query += " AND date >= $" + strconv.Itoa(i)
// 		values = append(values, film.Date)
// 		i++
// 	}

// 	if film.Rate != 0 {
// 		query += " AND rate >= $" + strconv.Itoa(i)
// 		values = append(values, film.Rate)
// 		i++
// 	}

// 	if len(film.Actors) > 0 {
// 		actorsConditions := []string{}
// 		for _, actor := range film.Actors {
// 			actorsConditions = append(actorsConditions, "$"+strconv.Itoa(i)+" = ANY(actors)")
// 			values = append(values, actor)
// 			i++
// 		}

// 		actorsQuery := "(" + strings.Join(actorsConditions, " OR ") + ")"

// 		query += " AND " + actorsQuery
// 	}

// 	query += " LIMIT 1"

// 	args := []interface{}{film.Name, film.Desc, film.Date, film.Rate}
// 	for _, actor := range film.Actors {
// 		args = append(args, actor)
// 	}

// 	var foundFilm model.Film

// 	row := db.pool.QueryRow(context.Background(), query, values...)

// 	var date sql.NullTime
// 	if err := row.Scan(new(interface{}), &foundFilm.Name, &foundFilm.Desc, &date, &foundFilm.Rate, &foundFilm.Actors); err != nil {
// 		if err == pgx.ErrNoRows {
// 			log.Println("err:GetFilm:NoRows:", err)
// 			return nil, app.ErrFilmNotFound
// 		}

// 		log.Println("err:GetFilm:Scan:", err)
// 		return nil, err
// 	}

// 	if date.Valid {
// 		foundFilm.Date = model.NewCustomDate(date.Time)
// 	}

// 	log.Println("sussess:GetFilm", foundFilm)
// 	log.Println(query)
// 	log.Println(args)
// 	return &foundFilm, nil
// }

func (db DB) DeleteFilm(film *model.Film) error {
	film, err := db.GetFilmByFilmName(film)
	if err != nil {

		if errors.Is(err, app.ErrFilmNotFound) {
			log.Println("err:DeleteFilm:FilmNotFound:", err)
			return app.ErrFilmNotFound
		}

		log.Println("err:DeleteFilm:GetFilm:", err)
		return err
	}

	query := "DELETE FROM films WHERE 1=1"

	values := make([]interface{}, 0)
	var i int = 1

	if film.Name != "" {
		query += " AND name = '%' || $" + strconv.Itoa(i) + " || '%'"
		values = append(values, film.Name)
		i++
	}

	if film.Desc != "" {
		query += " AND description ILIKE $" + strconv.Itoa(i)
		values = append(values, film.Desc)
		i++
	}

	if !film.Date.IsZero() {
		log.Println(film.Date)
		query += " AND date = $" + strconv.Itoa(i)
		values = append(values, film.Date)
		i++
	}

	if film.Rate != 0 {
		query += " AND rate = $" + strconv.Itoa(i)
		values = append(values, film.Rate)
		i++
	}

	ctx := context.Background()
	_, err = db.pool.Exec(ctx, query, values...)
	if err != nil {
		log.Println("err:DeleteFilm:Exec:", err)
		return err
	}

	log.Println("sussess:DeleteFilm:", film.Name)
	log.Println(query)
	log.Println(values...)
	return nil
}

func (db DB) GetFilms(film *model.Film, offset string, sort string, order bool) ([]model.Film, error) {

	query := "SELECT * FROM films WHERE 1=1"

	values := make([]interface{}, 0)
	var i int = 1

	if film.Name != "" {
		query += " AND name ILIKE '%' || $" + strconv.Itoa(i) + " || '%'"
		values = append(values, film.Name)
		i++
	}

	if film.Desc != "" {
		query += " AND description ILIKE $" + strconv.Itoa(i)
		values = append(values, film.Desc)
		i++
	}

	if !film.Date.IsZero() {
		log.Println(film.Date)
		query += " AND date >= $" + strconv.Itoa(i)
		values = append(values, film.Date)
		i++
	}

	if film.Rate != 0 {
		query += " AND rate >= $" + strconv.Itoa(i)
		values = append(values, film.Rate)
		i++
	}

	if len(film.Actors) > 0 {
		actorsConditions := []string{}
		for _, actor := range film.Actors {
			actorsConditions = append(actorsConditions, "$"+strconv.Itoa(i)+" = ANY(actors)")
			values = append(values, actor)
			i++
		}

		actorsQuery := "(" + strings.Join(actorsConditions, " OR ") + ")"

		query += " AND " + actorsQuery
	}

	if sort == "" {
		query += " ORDER BY rate DESC"
	} else {
		query += " ORDER BY " + sort
		if order {
			query += " DESC"
		} else {
			query += " ASC"
		}
	}

	query += " LIMIT " + offset

	args := []interface{}{film.Name, film.Desc, film.Date, film.Rate}
	for _, actor := range film.Actors {
		args = append(args, actor)
	}

	rows, err := db.pool.Query(context.Background(), query, values...)
	if err != nil {
		log.Println("err:GetFilms:Query:", err)
		return nil, err
	}
	defer rows.Close()

	var foundFilms []model.Film

	for rows.Next() {
		var foundFilm model.Film
		var date sql.NullTime
		if err := rows.Scan(new(interface{}), &foundFilm.Name, &foundFilm.Desc, &date, &foundFilm.Rate, &foundFilm.Actors); err != nil {
			log.Println("err:GetFilms:Scan:", err)
			return nil, err
		}
		if date.Valid {

			foundFilm.Date = model.NewCustomDate(date.Time)
		}
		foundFilms = append(foundFilms, foundFilm)
	}

	if len(foundFilms) == 0 {
		log.Println("err:GetFilms:len(foundFilms) == 0:", err)
		return nil, app.ErrFilmNotFound
	}

	log.Println("sussess:GetFilms", foundFilms)
	log.Println(query)
	log.Println(args)
	return foundFilms, nil
}
