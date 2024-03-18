package model

import (
	"encoding/json"
	"strings"
	"time"
)

type Actor struct {
	Name     string    `json:"name"`
	Gender   string    `json:"gender,omitempty"`
	Birthday time.Time `json:"birthday,omitempty"`
}

func (a *Actor) UnmarshalJSON(b []byte) error {
	type alias Actor
	d := struct {
		*alias
	}{
		alias: (*alias)(a),
	}

	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}

	a.Gender = strings.ToLower(a.Gender)

	return nil
}

type ActorUpdate struct {
	Name     string    `json:"name"`
	Gender   string    `json:"gender,omitempty"`
	Birthday time.Time `json:"birthday,omitempty"`
	Update   Actor     `json:"updateq"`
}
