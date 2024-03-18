package database

import (
	"testing"
	"time"
	"vk/app"
	"vk/model"
)

func TestCreateFilmWithExistingActors(t *testing.T) {

	testFilm := model.Film{
		Name: "Test Film",
		Desc: "Test description",
		Date: model.NewCustomDate(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)),
		Rate: 8.0,
		Actors: []model.Actor{
			{
				Name:     "John Doe",
				Gender:   "Male",
				Birthday: time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC),
			},
			{
				Name:     "Doe John",
				Gender:   "Male",
				Birthday: time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	err := testDB.CreateFilm(&testFilm)
	if err != nil {
		if err != app.ErrFilmConflict {
			t.Errorf("unexpected error: %v", err)
		}
		
	}

}

func TestGetFilmByFilmName(t *testing.T) {

	wantSuccessCase := model.Film{
		Name: "Test Film",
		Desc: "Test description",
	}

	film, err := testDB.GetFilmByFilmName(&wantSuccessCase)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}

	if film.Name != wantSuccessCase.Name || film.Desc != wantSuccessCase.Desc {
		t.Errorf("unwanted values: film have %s got %s", wantSuccessCase.Name, film.Name)
		t.Errorf("unwanted values: film have %s got %s", wantSuccessCase.Desc, film.Desc)
	}

	wantFailedCase := model.Film{
		Name: "Unexisted film",
		Desc: "Test description",
	}

	film, err = testDB.GetFilmByFilmName(&wantFailedCase)
	if err != app.ErrFilmNotFound {
		t.Errorf("unexpected error: %v %+v", err, film)
		t.FailNow()
	}
	if film != nil {
		t.Error("unexpected value: ", film)
		t.FailNow()
	}

}

func TestGetFilmByActorName(t *testing.T) {

	testFilm := model.Actor{
		Name: "John",
	}

	films, err := testDB.GetFilmsByActor(&testFilm)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Log(films)
	t.Log(testFilm)

}
