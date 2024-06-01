package entities

import "time"

type Todo struct {
	ID        string
	Title     string
	StartDate time.Time
	EndDate   time.Time
	UserID    string
}
