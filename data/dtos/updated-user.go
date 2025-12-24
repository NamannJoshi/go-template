package dtos

import "time"

type UpdatedUserDto struct {
	Username    *string    `json:"username"`
	Password    *string    `json:"password"`
	Email       *string    `json:"email"`
	Roles       *string    `json:"roles"`
	PhoneNumber *int       `json:"phone_number"`
	DateOfBirth *time.Time `json:"birth_of_date"`
}
