package repo

import "time"

type Note struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
}
