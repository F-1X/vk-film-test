package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"vk/app"
	"vk/model"
	"vk/verify"
)

func (s APIServer) GetActorHandler(w http.ResponseWriter, r *http.Request) interface{} {
	var req model.Actor
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	if !verify.ActorModel(&req) {
		return err
	}

	actor, err := s.repo.GetActor(req)
	if err != nil {
		if errors.Is(err, app.ErrActorNotFound) {
			return err
		} else {
			return err
		}
	}
	log.Println("[+][GetActorHandler]", req)
	return actor
}

// @Summary Create an Actor
// @Description Получает информацию о фильме из базы данных по его названию
// @Tags films
// @Accept json
// @Produce json
// @Param name, gender, birthday
// @Success 200 {object} FilmResponse
// @Failure 400 {object} ErrorResponse
// @Router /films/{name} [get]

func (s APIServer) PostActorHandler(w http.ResponseWriter, r *http.Request) interface{} {
	var req model.Actor
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return app.ErrBadRequest
	}
	if !verify.ActorModel(&req) {
		return app.ErrBadRequest
	}

	err = s.repo.PostActor(req)
	if err != nil {
		log.Println("[-][PostActorHandler][PostActor]", err)
		return err
	}
	log.Println("[+][PostActorHandler]", req)
	return app.ErrCreated
}

func (s APIServer) PatchActorHandler(w http.ResponseWriter, r *http.Request) interface{} {
	var req model.ActorUpdate
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("decode != nil", err)
		return app.ErrBadRequest
	}

	if !verify.ActorUpdateModel(&req) {
		return app.ErrBadRequest
	}

	err = s.repo.UpdateActor(req)
	if err != nil {
		log.Println("[-][PatchActorHandler][UpdateActor]", err)
		return err
	}
	log.Println("[+][PatchActorHandler]", req)
	return app.ErrSuccessUpdate
}

func (s APIServer) DeleteActorHandler(w http.ResponseWriter, r *http.Request) interface{} {
	var req model.Actor
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return app.ErrBadRequest
	}
	if !verify.ActorModel(&req) {
		return app.ErrBadRequest
	}
	err = s.repo.DeleteActor(req)
	if err != nil {
		log.Println("[-][PatchActorHandler][DeleteActor]", err)
		return err
	}
	log.Println("[+][DeleteActorHandler]", req)
	return app.ErrActorDelete
}

func (s APIServer) GetActorsHandler(w http.ResponseWriter, r *http.Request) interface{} {
	var req model.Actor
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	if !verify.ActorModel(&req) {
		return app.ErrBadRequest
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := verify.Offset(offsetStr)

	sort := r.URL.Query().Get("sort")
	if sort != "rate" && sort != "date" && sort != "name" && sort != "" {
		return err
	}

	actors, err := s.repo.GetActors(req, offset)
	if err != nil {
		if errors.Is(err, app.ErrActorNotFound) {
			log.Println("[-][PatchActorHandler][GetActors][ErrActorNotFound]", err)
			return app.ErrActorNotFound
		}
		log.Println("[-][PatchActorHandler][GetActors][ErrBadRequest]", err)
		return app.ErrBadRequest
	}
	log.Println("[+][PatchActorHandler]", req)
	return actors
}
