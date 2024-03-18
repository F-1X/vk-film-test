package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"vk/model"
)

var testCasesSignIn = []struct {
	name         string
	give         model.Credentials
	wantCode     int
	wantResponce string
}{
	{
		name: "signin success",
		give: model.Credentials{
			Username: "Gleb",
			Password: "qwerty1234",
		},
		wantCode:     409,
		wantResponce: `{"status":409,"message":"authorization success"}`,
	},
	{
		name: "password failed",
		give: model.Credentials{
			Username: "Gleb",
			Password: "qwerty12345",
		},
		wantCode:     401,
		wantResponce: `{"status":401,"message":"failed singin"}`,
	},
	{
		name: "username failed",
		give: model.Credentials{
			Username: "Gleb1",
			Password: "qwerty1234",
		},
		wantCode:     404,
		wantResponce: `{"status":404,"message":"not found"}`,
	},
	{
		name: "admin success",
		give: model.Credentials{
			Username: "admin",
			Password: "admin123",
		},
		wantCode:     409,
		wantResponce: `{"status":409,"message":"authorization success"}`,
	},
}

func TestAPIServer_SignIn(t *testing.T) {
	for _, tt := range testCasesSignIn {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.give)

			req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewReader(reqBody))

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(wrapHandleFunc(NewAPIServer(testServer.cfg, *testDB).SignIn))

			handler.ServeHTTP(rr, req)

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
		})
	}
}

var testCasesSignUp = []struct {
	name         string
	give         model.Credentials
	wantCode     int
	wantResponce string
}{
	{
		name: "signup success",
		give: model.Credentials{
			Username: "Glebka",
			Password: "qwertyzxc",
		},
		wantCode:     409,
		wantResponce: `{"status":409,"message":"authorization success"}`,
	},
	{
		name: "already exist",
		give: model.Credentials{
			Username: "Glebka",
			Password: "qwerty",
		},
		wantCode:     409,
		wantResponce: `{"status":409,"message":"bad password sorry"}`,
	},
	{
		name: "bad password",
		give: model.Credentials{
			Username: "Glebka",
			Password: "qwertyzxc",
		},
		wantCode:     409,
		wantResponce: `{"status":409,"message":"already exist"}`,
	},
}

func TestAPIServer_SignUp(t *testing.T) {
	for _, tt := range testCasesSignUp {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.give)

			req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewReader(reqBody))

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(wrapHandleFunc(NewAPIServer(testServer.cfg, *testDB).SignUp))

			handler.ServeHTTP(rr, req)

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
		})
	}
}

var testCasesAdminHandler = []struct {
	name                  string
	giveCreds             model.Credentials
	giveActor             model.Actor
	wantCodeSignIn        int
	wantResponceSignIn    string
	wantCodePostActor     int
	wantResponcePostActor string
}{
	{
		name: "admin signin success",
		giveCreds: model.Credentials{
			Username: "admin",
			Password: "admin123",
		},
		giveActor: model.Actor{
			Name:     "mikle",
			Gender:   "male",
			Birthday: time.Date(2024, time.March, 16, 0, 0, 0, 0, time.UTC),
		},
		wantCodeSignIn:        409,
		wantResponceSignIn:    `{"status":409,"message":"authorization success"}`,
		wantCodePostActor:     201,
		wantResponcePostActor: `{"status":201,"message":"success created"}`,
	},
}

func TestAPIServer_PostActor(t *testing.T) {
	for _, tt := range testCasesAdminHandler {
		t.Run(tt.name, func(t *testing.T) {

			adminCredsJSON, _ := json.Marshal(tt.giveCreds)
			req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(adminCredsJSON))

			signInRR := httptest.NewRecorder()

			handler := http.HandlerFunc(wrapHandleFunc(NewAPIServer(testServer.cfg, *testDB).SignIn))

			handler.ServeHTTP(signInRR, req)

			var sessionToken string
			for _, cookie := range signInRR.Result().Cookies() {
				if cookie.Name == "session_token" {
					sessionToken = cookie.Value
					break
				}
			}

			t.Log("case:", tt.name)
			if status := signInRR.Code; status != tt.wantCodeSignIn {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCodeSignIn)
				t.FailNow()
			}

			if status := strings.TrimSpace(signInRR.Body.String()); status != tt.wantResponceSignIn {
				t.Errorf("handler returned wrong response: got %v want %v",
					strings.TrimSpace(signInRR.Body.String()), tt.wantResponceSignIn)
				t.FailNow()
			}

			reqBody, _ := json.Marshal(tt.giveActor)
			postActorReq, _ := http.NewRequest(http.MethodPost, "/api/v1/actor", bytes.NewBuffer(reqBody))
			postActorReq.AddCookie(&http.Cookie{Name: "session_token", Value: sessionToken})

			postActorRR := httptest.NewRecorder()
			handler = http.HandlerFunc(wrapHandleFunc(NewAPIServer(testServer.cfg, *testDB).PostActorHandler))

			handler.ServeHTTP(postActorRR, postActorReq)

			t.Log("case:", tt.name)
			if status := postActorRR.Code; status != tt.wantCodePostActor {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantCodePostActor)
					t.FailNow()
			}
			if strings.TrimSpace(postActorRR.Body.String()) != tt.wantResponcePostActor {
				t.Errorf("handler returned wrong response: got %v want %v",
					strings.TrimSpace(postActorRR.Body.String()), tt.wantResponcePostActor)
				t.FailNow()
			}

		})
	}
}

// var testCasesAdminHandler = []struct {
// 	name         string
// 	give         model.Credentials
// 	wantCode     int
// 	wantResponce string
// }{
// 	{
// 		name: "delete failed",
// 		give: model.Credentials{
// 			Username: "Glebka",
// 			Password: "qwerty",
// 		},
// 		wantCode:     401,
// 		wantResponce: `{"status":401,"message":"bad password sorry"}`,
// 	},
// }

// func TestAPIServer_DeleteUser(t *testing.T) {
// 	for _, tt := range testCasesAdminHandler {
// 		t.Run(tt.name, func(t *testing.T) {

// 			adminCreds := model.Credentials{Username: "admin", Password: "adminpassword"}
// 			adminCredsJSON, _ := json.Marshal(adminCreds)
// 			signInReq, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(adminCredsJSON))

// 			reqBody, _ := json.Marshal(tt.give)

// 			req, _ := http.NewRequest(http.MethodPost, "/admin/user", bytes.NewReader(reqBody))

// 			rr := httptest.NewRecorder()

// 			handler := http.HandlerFunc(wrapHandleFunc(NewAPIServer(testServer.cfg, *testDB).DeleteUser))

// 			handler.ServeHTTP(rr, req)

// 			if status := rr.Code; status != tt.wantCode {
// 				t.Errorf("handler returned wrong status code: got %v want %v",
// 					status, tt.wantCode)
// 				t.FailNow()
// 			}
// 			if strings.TrimSpace(rr.Body.String()) != tt.wantResponce {
// 				t.Errorf("handler returned wrong response: got %v want %v",
// 					strings.TrimSpace(rr.Body.String()), tt.wantResponce)
// 				t.FailNow()
// 			}
// 		})
// 	}
// }
