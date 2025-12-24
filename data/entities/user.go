package entities

import (
	"time"
)

type User struct {
	Id           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Password     string    `db:"password" json:"password"`
	Email        string    `db:"email" json:"email"`
	PhoneNumber  int       `db:"phone_number" json:"phone_number"`
	Roles        string    `db:"roles" json:"roles"`
	DateOfBirth  time.Time `db:"birth_of_date" json:"birth_of_date"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	LastModified time.Time `db:"last_modified_at" json:"last_modified_at"`
}
