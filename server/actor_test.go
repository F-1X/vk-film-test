package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"vk/config"
	"vk/database"
	"vk/model"
)

var testDB *database.DB
var testServer *APIServer

func TestMain(m *testing.M) {
	cfg, err := config.Read("../config.yml")
	if err != nil {
		log.Fatalf("config read failed: %v", err)
	}

	testDB, err = database.New(cfg.Database)
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	testServer = NewAPIServer(cfg.Server, *testDB)

	m.Run()
}

var testCasesPostActorHandler = []struct {
	name         string
	give         model.Actor
	wantCode     int
	wantResponce string
}{
	{
		name: "Succees add actor1",
		give: model.Actor{
			Name:     "Tom Hanks1",
			Gender:   "Male",
			Birthday: time.Date(1956, time.July, 9, 0, 0, 0, 0, time.UTC),
		},
		wantCode:     201,
		wantResponce: `{"status":201,"message":"success created"}`,
	},
	{
		name: "Succees add actor2",
		give: model.Actor{
			Name:     "Tom Hanks2",
			Gender:   "Male",
			Birthday: time.Date(1956, time.July, 9, 0, 0, 0, 0, time.UTC),
		},
		wantCode:     201,
		wantResponce: `{"status":201,"message":"success created"}`,
	},
	{
		name: "Duplicate add actor1",
		give: model.Actor{
			Name:     "Tom Hanks1",
			Gender:   "Male",
			Birthday: time.Date(1956, time.July, 9, 0, 0, 0, 0, time.UTC),
		},
		wantCode:     409,
		wantResponce: `{"status":409,"message":"actor already exist"}`,
	},
}

func TestAPIServer_PostActorHandler(t *testing.T) {

	for _, tt := range testCasesPostActorHandler {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.give)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/actor", bytes.NewReader(reqBody))

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(wrapHandleFunc(NewAPIServer(testServer.cfg, *testDB).PostActorHandler))

			handler.ServeHTTP(rr, req)
			t.Log("case:", tt.name)
			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			if strings.TrimSpace(rr.Body.String()) != tt.wantResponce {
				t.Errorf("handler returned wrong response: got %v want %v",
					strings.TrimSpace(rr.Body.String()), tt.wantResponce)
			}
		})
	}
}

var testCasesPatchActorHandler = []struct {
	name         string
	give         model.ActorUpdate
	wantCode     int
	wantResponce string
}{
	{
		name: "Succees patch actor",
		give: model.ActorUpdate{
			Name:     "Tom Hanks2",
			Gender:   "Male",
			Birthday: time.Date(1956, time.July, 9, 0, 0, 0, 0, time.UTC),
			Update: model.Actor{
				Name: "Tom Hanks2 updated name",
			},
		},
		wantCode:     201,
		wantResponce: `{"status":201,"message":"update success"}`,
	},
	{
		name: "Repeat last patch",
		give: model.ActorUpdate{
			Name:     "Tom Hanks2",
			Gender:   "Male",
			Birthday: time.Date(1956, time.July, 9, 0, 0, 0, 0, time.UTC),
			Update: model.Actor{
				Name: "Tom Hanks2 updated name",
			},
		},
		wantCode:     404,
		wantResponce: `{"status":404,"message":"actor not found"}`,
	},
	{
		name: "Return it back last patch",
		give: model.ActorUpdate{
			Name:     "Tom Hanks2 updated name",
			Gender:   "Male",
			Birthday: time.Date(1956, time.July, 9, 0, 0, 0, 0, time.UTC),
			Update: model.Actor{
				Name: "Tom Hanks2",
			},
		},
		wantCode:     201,
		wantResponce: `{"status":201,"message":"update success"}`,
	},
}

func TestAPIServer_PatchActorHandler(t *testing.T) {

	for _, tt := range testCasesPatchActorHandler {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.give)
			req, _ := http.NewRequest(http.MethodPatch, "/api/v1/actor", bytes.NewReader(reqBody))

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(wrapHandleFunc(NewAPIServer(testServer.cfg, *testDB).PatchActorHandler))

			handler.ServeHTTP(rr, req)

			t.Log("case:", tt.name)
			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
				t.FailNow()
			}
			if strings.TrimSpace(rr.Body.String()) != tt.wantResponce {
				t.Errorf("handler returned wrong status code: got %v want %v",
					strings.TrimSpace(rr.Body.String()), tt.wantResponce)
				t.FailNow()
			}
		})
	}
}

var testCasesDeleteActorHandler = []struct {
	name         string
	give         model.Actor
	wantCode     int
	wantResponce string
}{
	{
		name: "Succees delete actor1",
		give: model.Actor{
			Name:     "Tom Hanks1",
			Gender:   "Male",
			Birthday: time.Date(1956, time.July, 9, 0, 0, 0, 0, time.UTC),
		},
		wantCode:     204,
		wantResponce: `{"status":204,"message":"delete success"}`,
	},
	{
		name: "Duplicate delete actor1",
		give: model.Actor{
			Name:     "Tom Hanks1",
			Gender:   "Male",
			Birthday: time.Date(1956, time.July, 9, 0, 0, 0, 0, time.UTC),
		},
		wantCode:     404,
		wantResponce: `{"status":404,"message":"actor not found"}`,
	},
	{
		name: "Duplicate delete actor2",
		give: model.Actor{
			Name:     "Tom Hanks2",
			Gender:   "Male",
			Birthday: time.Date(1956, time.July, 9, 0, 0, 0, 0, time.UTC),
		},
		wantCode:     204,
		wantResponce: `{"status":204,"message":"delete success"}`,
	},
	{
		name: "Not exists actor delete",
		give: model.Actor{
			Name:     "Not exist actor",
			Gender:   "Male",
			Birthday: time.Date(1956, time.July, 9, 0, 0, 0, 0, time.UTC),
		},
		wantCode:     404,
		wantResponce: `{"status":404,"message":"actor not found"}`,
	},
	{
		name: "Breaked model request",
		give: model.Actor{
			Gender:   "Male",
			Birthday: time.Date(1956, time.July, 9, 0, 0, 0, 0, time.UTC),
		},
		wantCode:     400,
		wantResponce: `{"status":400,"message":"bad request"}`,
	},
}

func TestAPIServer_DeleteActorHandler(t *testing.T) {

	for _, tt := range testCasesDeleteActorHandler {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.give)
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/actor", bytes.NewReader(reqBody))

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(wrapHandleFunc(NewAPIServer(testServer.cfg, *testDB).DeleteActorHandler))

			handler.ServeHTTP(rr, req)

			t.Log("case:", tt.name)
			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
				t.FailNow()
			}
			if strings.TrimSpace(rr.Body.String()) != tt.wantResponce {
				t.Errorf("handler returned wrong status code: got %v want %v",
					strings.TrimSpace(rr.Body.String()), tt.wantResponce)
				t.FailNow()
			}
		})
	}
}

var testCasesGetActors = []struct {
	name            string
	giveCreds       model.Credentials
	giveActors      model.Actor
	wantCode        int
	wantResponce    string
	wantCodeAct     int
	wantResponceAct string
}{
	{
		name: "signin success",
		giveCreds: model.Credentials{
			Username: "Gleb",
			Password: "qwerty1234",
		},
		giveActors: model.Actor{
			Name: "tom",
		},
		wantCode:        409,
		wantResponce:    `{"status":409,"message":"authorization success"}`,
		wantCodeAct:     409,
		wantResponceAct: `{"status":409,"message":"authorization success"}`,
	},
}

func TestAPIServer_GetActors(t *testing.T) {
	for _, tt := range testCasesGetActors {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.giveCreds)

			req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewReader(reqBody))

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(wrapHandleFunc(NewAPIServer(testServer.cfg, *testDB).SignIn))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
				t.FailNow()
			}
			if strings.TrimSpace(rr.Body.String()) != tt.wantResponce {
				t.Errorf("handler returned wrong response: got %v want %v",
					strings.TrimSpace(rr.Body.String()), tt.wantResponce)
				t.FailNow()
			}

			reqBodyAct, _ := json.Marshal(tt.giveActors)

			reqAct, _ := http.NewRequest(http.MethodPost, "/api/v1/actors?offset=2", bytes.NewReader(reqBodyAct))

			rr2 := httptest.NewRecorder()

			handler = http.HandlerFunc(wrapHandleFunc(NewAPIServer(testServer.cfg, *testDB).GetActorsHandler))

			handler.ServeHTTP(rr2, reqAct)

			t.Log("case:", tt.name)
			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCode)
				t.FailNow()
			}
			if strings.TrimSpace(rr.Body.String()) != tt.wantResponce {
				t.Errorf("handler returned wrong response: got %v want %v",
					strings.TrimSpace(rr.Body.String()), tt.wantResponce)
				t.FailNow()
			}

			t.Log(rr2.Body.String())

		})
	}
}
