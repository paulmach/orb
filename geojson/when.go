package geojson

import "time"

// "When" is a datetime info bound to Features objects
// Geojson spec at https://github.com/geojson/geojson-ld

type When struct {
	Type         string           `json:"@type,omitempty"`
	Datetime     *time.Time        `json:"datetime,omitempty"`
}

// NewWhen creates a when clause
func NewWhen(Type string, Datetime *time.Time) When {
	return When{Type, Datetime}
}

func (w *When) Valid() bool {
	if w.Type == "Instant" || w.Type=="Interval" {
		return w.Datetime.Year() != 1
	}
	return false
}

