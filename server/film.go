package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"vk/app"
	"vk/model"
	"vk/verify"
)
type ErrBadRequest struct {
    Message string `json:"message"`
}
// GetFilm godoc
// @Summary      Get a film by name, description or by actors
// @Description  Get a film by providing its name, description or actors.
// @Tags         films
// @Accept       json
// @Produce      json
// @Param        name      query    string  false  "Film name"
// @Param        description  query    string  false  "Film description"
// @Param        actors  query    string  false  "Comma-separated list of actor names"
// @Success      200  {object}  model.Film
// @Failure      400  {object}  ErrBadRequest
// @Router       /api/v1/films [get]
func (s APIServer) GetFilm(w http.ResponseWriter, r *http.Request) interface{} {

	body, err := checkRequestBodyType(r)
	if err != nil {
		return app.ErrBadRequest
	}

	switch t := body.(type) {
	case model.Film:
		res, err := s.repo.GetFilmByFilmName(&t)
		if err != nil {
			log.Println("APIServer:GetFilmByFilmName:err:", err)
			return err
		}
		if !verify.FilmModel(&t) {
			return app.ErrBadRequest
		}
		log.Println("APIServer:GetFilmByFilmName:success:", t.Name)
		return res
	case model.Actor:
		res, err := s.repo.GetFilmsByActor(&t)
		if err != nil {
			log.Println("APIServer:GetFilmsByActor:err:", err)
			return err
		}
		if !verify.Actor(t.Name) {
			return app.ErrBadRequest
		}
		log.Println("APIServer:GetFilmsByActor:success:", t.Name)
		return res
	}

	return nil

}

func checkRequestBodyType(r *http.Request) (interface{}, error) {
	var film model.Film
	var actor model.Actor

	err := json.NewDecoder(r.Body).Decode(&film)
	if err == nil {
		return film, nil
	}

	err = json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		return nil, errors.New("bad type of request")
	}

	return actor, nil
}

func (s APIServer) PostFilm(w http.ResponseWriter, r *http.Request) interface{} {
	var req model.Film
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	if !verify.FilmModel(&req) {
		return err
	}
	err = s.repo.CreateFilm(&req)
	if err != nil {
		log.Println("APIServer:PostFilm:err:", err)
		return err
	}
	log.Println("APIServer:PostFilm:success:", req.Name)
	return err
}

func (s APIServer) DeleteFilm(w http.ResponseWriter, r *http.Request) interface{} {
	var req model.Film
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	if !verify.FilmModel(&req) {
		return err
	}
	err = s.repo.DeleteFilm(&req)
	if err != nil {
		log.Println("APIServer:DeleteFilm:success:")
		return err
	}
	log.Println("APIServer:DeleteFilm:NoContent:")
	return http.StatusNoContent

}

func (s APIServer) GetFilms(w http.ResponseWriter, r *http.Request) interface{} {
	var req model.Film

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

		return http.StatusBadRequest
	}
	if !verify.FilmModel(&req) {

		return err
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := verify.Offset(offsetStr)

	sort := r.URL.Query().Get("sort")
	if sort != "rate" && sort != "date" && sort != "name" && sort != "" {
		return err
	}

	order := r.URL.Query().Get("order")
	asc := true
	if order != "" {
		if order != "1" && order != "0" && order != "true" && order != "false" {
			return http.StatusBadRequest
		}

		asc, err = strconv.ParseBool(order)
		if err != nil {
			return http.StatusBadRequest
		}
	}

	res, err := s.repo.GetFilms(&req, offset, sort, asc)
	if err != nil {
		log.Println("APIServer:GetFilms:err:", err)
		return err
	}
	log.Println("APIServer:GetFilms:success:", res)
	return res
}
