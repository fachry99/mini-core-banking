package domain

import "time"

type User struct {
	ID        string    `db:"id"`
	FullName  string    `db:"full_name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}
