package urlalias

import "time"

type URL struct {
	Alias     string
	Original  string
	ExpiresAt time.Time
	Visits    uint
}

func New(alias, original string) URL {
	return URL{
		Alias:     alias,
		Original:  original,
		ExpiresAt: time.Now().AddDate(0, 1, 0),
		Visits:    0,
	}
}
