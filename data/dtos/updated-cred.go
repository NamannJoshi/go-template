package dtos

import "time"

type UpdatedCredDto struct {
	UserId       *int           `db:"user_id" json:"userId"`
	AccessToken  *string        `db:"access_token" json:"accessToken"`
	RefreshToken *string        `db:"refresh_token" json:"refreshToken"`
	ExpiresIn    *time.Duration `db:"expires_in" json:"expiresIn"`
}
