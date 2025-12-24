package repositories

import (
	"context"
	"fmt"
	"template/data/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IUserRepository interface {
	GetAllUsers(context.Context) ([]entities.User, error)
	CreateUser(context.Context, entities.User) error
	GetUserById(context.Context, int) (entities.User, error)
	GetUserByUsername(context.Context, string) (entities.User, error)
	UpdateUser(context.Context, int, entities.User) error
	DeleteUser(context.Context, int) error
}

type UserRepoStruct struct {
	DB *pgxpool.Pool
}

func UserRepoInit(db *pgxpool.Pool) *UserRepoStruct {
	return &UserRepoStruct{DB: db}
}

func (u *UserRepoStruct) GetAllUsers(context context.Context) ([]entities.User, error) {
	query := "select * from users;"

	rows, err := u.DB.Query(context, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.User])
	fmt.Println("usersrzz", users)

	return users, nil
}

func (u *UserRepoStruct) CreateUser(context context.Context, user entities.User) error {
	fmt.Println(user.PhoneNumber)
	query := `insert into users (username, password, email, roles, phone_number, birth_of_date, last_modified_at) values (@username, @password, @email, @roles, @phoneNumber, @dateOfBirth, now());`

	args := pgx.NamedArgs{
		"username":    user.Username,
		"password":    user.Password,
		"roles":       user.Roles,
		"email":       user.Email,
		"phoneNumber": user.PhoneNumber,
		"dateOfBirth": user.DateOfBirth,
	}

	_, err := u.DB.Exec(context, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepoStruct) GetUserById(context context.Context, id int) (entities.User, error) {
	query := `select * from users where id=@id;`

	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := u.DB.Query(context, query, args)
	if err != nil {
		return entities.User{}, err
	}
	fmt.Println("entered successfully", rows)
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.User])
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (u *UserRepoStruct) GetUserByUsername(context context.Context, username string) (entities.User, error) {
	query := `select * from users where username=@username;`

	args := pgx.NamedArgs{
		"username": username,
	}

	rows, err := u.DB.Query(context, query, args)
	if err != nil {
		return entities.User{}, err
	}
	fmt.Println("entered successfully", rows)
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entities.User])
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (u *UserRepoStruct) UpdateUser(context context.Context, id int, updatedUser entities.User) error {
	query := `update users set username=@username, password=@password, email=@email, roles=@roles, phone_number=@phoneNumber, birth_of_date=@dateOfBirth, last_modified_at=now() where id=@id;`

	args := pgx.NamedArgs{
		"id":          id,
		"username":    updatedUser.Username,
		"password":    updatedUser.Password,
		"roles":       updatedUser.Roles,
		"email":       updatedUser.Email,
		"phoneNumber": updatedUser.PhoneNumber,
		"dateOfBirth": updatedUser.DateOfBirth,
	}

	_, err := u.DB.Exec(context, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepoStruct) DeleteUser(context context.Context, id int) error {
	query := `delete from users where id=@id;`

	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := u.DB.Exec(context, query, args)
	if err != nil {
		return err
	}

	return nil
}
