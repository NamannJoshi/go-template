package dtos

import "time"

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDto struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    time.Duration `json:"expires_in"`
}

type RegisterDto struct {
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `db:"birth_of_date" json:"birth_of_date"`
	PhoneNumber int       `db:"phone_number" json:"phone_number"`
	Password    string    `json:"password"`
}

type RefreshTokenDto struct {
	RefreshTokenDto string `json:"refresh_token"`
}
