package verify

import (
	"log"
	"regexp"
	"strconv"
	"time"
	"vk/model"
)

func Actor(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Zа-яА-Я ]+$`)
	log.Println(re.MatchString(s))
	return re.MatchString(s)
}

func Gender(s string) bool {
	return s == "male" || s == "female" || s == ""
}

func Birthday(t time.Time) time.Time {
	layout := "2006-01-02"
	dateString := t.Format(layout)
	newTime, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}
	}
	return newTime
}

func ActorsModel(actors []model.Actor) bool {
	for _, actor := range actors {
		if !Actor(actor.Name) {
			return false
		}
		if !Gender(actor.Gender) {
			return false
		}
		actor.Birthday = Birthday(actor.Birthday)

	}
	return true
}

func ActorModel(actor *model.Actor) bool {
	if !Actor(actor.Name) {
		log.Println("actor failed here", actor.Name)
		return false
	}
	if !Gender(actor.Gender) {
		log.Println("gender failed")
		return false
	}
	actor.Birthday = Birthday(actor.Birthday)
	return true
}

func Film(s string) bool {
	l := len(s)
	if l == 0 {
		return true
	}
	if l > 150 {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9 ]+$`)
	return re.MatchString(s)
}

func Desc(s string) bool {
	return len(s) < 1001
}

func Rate(rate float32) bool {
	return rate >= 0 && rate <= 10
}

func FilmModel(f *model.Film) bool {

	if !Film(f.Name) {
		log.Println("film failed")
		return false
	}

	if !Desc(f.Desc) {
		log.Println("desc failed")
		return false
	}

	if !Rate(f.Rate) {
		log.Println("rate failed")
		return false
	}

	if !ActorsModel(f.Actors) {
		log.Println("actors failed")
		return false
	}

	return true
}

func Offset(f string) string {
	if f == "" {
		return "100"
	}
	offset, _ := strconv.Atoi(f)
	if offset >= 100 {
		return "100"
	}
	if offset <= 0 {
		return "100"
	}
	return f
}

func ActorUpdateModel(f *model.ActorUpdate) bool {

	if !Actor(f.Name) {
		log.Println("actor failed")
		return false
	}

	if len(f.Update.Name) != 0 {
		if !Actor(f.Update.Name) {
			log.Println("new actor failed")

			return false
		}
	}

	return true
}

func Creds(creds model.Credentials) bool {
	if len(creds.Username) == 0 {
		return false
	}

	if len(creds.Password) < 8 {
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9]+$`)
	if !re.MatchString(creds.Username) {
		return false
	}

	re = regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9!@#$]+$`)

	return re.MatchString(creds.Password)
}
