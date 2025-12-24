package repositories

import (
	"context"
	"fmt"
	"template/data/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ICredentialRepository interface {
	GetAllCredentials(context.Context) ([]entities.Credential, error)
	CreateCredential(context.Context, entities.Credential) error
	GetCredentialById(context.Context, int) (entities.Credential, error)
	UpdateCredential(context.Context, int, entities.Credential) error
	DeleteCredential(context.Context, int) error
}

type CredentialRepoStruct struct {
	DB *pgxpool.Pool
}

func CredentialRepoInit(db *pgxpool.Pool) *CredentialRepoStruct {
	return &CredentialRepoStruct{DB: db}
}

func (u *CredentialRepoStruct) GetAllCredentials(context context.Context) ([]entities.Credential, error) {
	query := "select * from creds;"

	rows, err := u.DB.Query(context, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	creds, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.Credential])

	return creds, nil
}

func (u *CredentialRepoStruct) CreateCredential(context context.Context, creds entities.Credential) error {
	query := `insert into creds (user_id, access_token, refresh_token, expires_in) values (@userId, @accessToken, @refreshToken, @expiresIn);`

	args := pgx.NamedArgs{
		"user_id":       creds.UserId,
		"access_token":  creds.AccessToken,
		"refresh_token": creds.RefreshToken,
		"expires_in":    creds.ExpiresIn,
	}

	_, err := u.DB.Exec(context, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (u *CredentialRepoStruct) GetCredentialById(context context.Context, id int) (entities.Credential, error) {
	query := `select * from creds where id=@id;`

	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := u.DB.Query(context, query, args)
	if err != nil {
		return entities.Credential{}, err
	}
	fmt.Println("entered successfully", rows)
	cred, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.Credential])
	if err != nil {
		return entities.Credential{}, err
	}

	return cred, nil
}

func (u *CredentialRepoStruct) UpdateCredential(context context.Context, id int, updatedUser entities.Credential) error {
	query := `update users set username=@username, email=@email, phone_number=@phoneNumber, birth_of_date=@dateOfBirth, last_modified_at=now() where id=@id;`

	args := pgx.NamedArgs{
		// "id":          id,
		// "username":    updatedUser.Username,
		// "email":       updatedUser.Email,
		// "phoneNumber": updatedUser.PhoneNumber,
		// "dateOfBirth": updatedUser.DateOfBirth,
	}

	_, err := u.DB.Exec(context, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (u *CredentialRepoStruct) DeleteCredential(context context.Context, userId int) error {
	query := `delete from creds where user_id=@user_id;`

	args := pgx.NamedArgs{
		"user_id": userId,
	}

	_, err := u.DB.Exec(context, query, args)
	if err != nil {
		return err
	}

	return nil
}
