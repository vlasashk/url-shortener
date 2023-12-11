package urlalias

import "time"

type URL struct {
	Alias     string    `json:"-"`
	Original  string    `json:"original" validate:"required,url"`
	ExpiresAt time.Time `json:"-"`
	Visits    uint      `json:"-"`
}

func New(alias, original string) URL {
	return URL{
		Alias:     alias,
		Original:  original,
		ExpiresAt: time.Now().AddDate(0, 1, 0),
		Visits:    0,
	}
}
