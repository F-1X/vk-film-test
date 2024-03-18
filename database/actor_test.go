package database

import (
	"log"
	"testing"
	"time"
	"vk/app"
	"vk/config"
	"vk/model"
)

var testDB *DB

func TestMain(m *testing.M) {
	cfg, err := config.Read("../config.yml")
	if err != nil {
		log.Fatalf("config read failed: %v", err)
	}

	testDB, err = New(cfg.Database)
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	m.Run()
}

func TestPostActor(t *testing.T) {

	testActor1 := model.Actor{Name: "John Doe", Gender: "Male", Birthday: time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC)}
	testActor2 := model.Actor{Name: "Doe John", Gender: "Male", Birthday: time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC)}

	err := testDB.PostActor(testActor1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = testDB.PostActor(testActor2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

}

func TestGetActor(t *testing.T) {

	testActor := model.Actor{Name: "John Doe", Gender: "Male", Birthday: time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC)}

	findActor, err := testDB.GetActor(testActor)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}

	if findActor.Name != testActor.Name || findActor.Birthday != testActor.Birthday || findActor.Gender != testActor.Gender {
		t.Error("unexpected actor", findActor, testActor)
	}

}

func TestFailedGetActor(t *testing.T) {

	testUnexistedActor := model.Actor{Name: "John Doe2", Gender: "Male", Birthday: time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC)}

	_, err := testDB.GetActor(testUnexistedActor)
	if err != app.ErrActorNotFound {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}
}
