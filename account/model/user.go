package model

import "github.com/google/uuid"

// User model.
type User struct {
	UID      uuid.UUID `db:"uid" json:"uid"`
	Email    string    `db:"email" json:"email"`
	Password string    `db:"password" json:"-"`
	Username string    `db:"username" json:"username"`
	ImageURL string    `db:"image_url" json:"imageURL"`
	Website  string    `db:"website" json:"website"`
}
