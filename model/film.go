package model

import "time"

type Film struct {
	ID     int
	Name   string     `json:"name"`
	Desc   string     `json:"desc"`
	Date   CustomDate `json:"date"`
	Rate   float32    `json:"rate"`
	Actors []Actor    `json:"actors"`
}

type Films struct {
	Films  []Film
	Offset int
}

type CustomDate struct {
	time.Time
}

func (d *CustomDate) UnmarshalJSON(data []byte) error {

	trimmed := string(data[1 : len(data)-1])

	parsedDate, err := time.Parse("2006-01-02", trimmed)
	if err != nil {
		return err
	}

	d.Time = parsedDate

	return nil
}

func NewCustomDate(t time.Time) CustomDate {
	return CustomDate{t}
}
