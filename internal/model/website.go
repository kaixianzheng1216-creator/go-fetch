package model

import "time"

type Website struct {
	ID        string
	Name      string
	Domain    string
	CreatedAt time.Time
}
