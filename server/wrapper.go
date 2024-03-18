package server

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"vk/app"
	"vk/model"
)

type apiFunc func(http.ResponseWriter, *http.Request) interface{}

func wrapHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteJSON(w, err)
		}

	}
}

func WriteJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	log.Println("WriteJSON", v, reflect.TypeOf(v))
	switch t := v.(type) {

	case model.Actor:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(t)

	case *model.Actor:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(t)

	case []model.Actor:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(t)

	case *model.Film:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(t)

	case []model.Film:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(t)

	case model.Credentials:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(t)

	case *app.APIError:
		w.WriteHeader(t.Status)
		json.NewEncoder(w).Encode(t)

	case *app.APIMessage:
		log.Println("*app.APIMessage", t, t.Status)
		w.WriteHeader(t.Status)
		json.NewEncoder(w).Encode(t)

	case error:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("some internal error")

	default:
		log.Println("unknown type", reflect.TypeOf(v))
		json.NewEncoder(w).Encode("unknown error")
	}

}
