package services

import (
	"context"
	"template/data/dtos"
	"template/data/entities"
	"template/repositories"
	"time"
)

type IUserService interface {
	GetAllUsers(context.Context) ([]entities.User, error)
	CreateUser(context.Context, entities.User) error
	GetUserById(context.Context, int) (*entities.User, error)
	GetUserByUsername(context.Context, string) (*entities.User, error)
	UpdateUser(context.Context, int, dtos.UpdatedUserDto) (entities.User, error)
	DeleteUser(context.Context, int) error
}

type UserServiceStruct struct {
	userRepo repositories.IUserRepository
}

func UserServiceInit(repo repositories.IUserRepository) *UserServiceStruct {
	return &UserServiceStruct{userRepo: repo}
}

func (u *UserServiceStruct) GetAllUsers(context context.Context) ([]entities.User, error) {
	users, err := u.userRepo.GetAllUsers(context)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserServiceStruct) CreateUser(context context.Context, user entities.User) error {
	if err := u.userRepo.CreateUser(context, user); err != nil {
		return err
	}

	// user, err := u.userRepo.GetUserByUsername(context, user.Username)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (u *UserServiceStruct) GetUserById(context context.Context, id int) (*entities.User, error) {
	user, err := u.userRepo.GetUserById(context, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserServiceStruct) GetUserByUsername(context context.Context, username string) (*entities.User, error) {
	user, err := u.userRepo.GetUserByUsername(context, username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserServiceStruct) UpdateUser(context context.Context, id int, user dtos.UpdatedUserDto) (entities.User, error) {
	currentUser, err := u.userRepo.GetUserById(context, id)
	if err != nil {
		return entities.User{}, err
	}

	updated := entities.User{
		Id: id,

		Username: func() string {
			if user.Username != nil {
				return *user.Username
			}
			return currentUser.Username
		}(),

		Roles: func() string {
			if user.Roles != nil {
				return *user.Roles
			}
			return currentUser.Username
		}(),

		Password: func() string {
			if user.Password != nil {
				return *user.Password
			}
			return currentUser.Password
		}(),

		Email: func() string {
			if user.Email != nil {
				return *user.Email
			}
			return currentUser.Email
		}(),

		PhoneNumber: func() int {
			if user.PhoneNumber != nil {
				return *user.PhoneNumber
			}
			return currentUser.PhoneNumber
		}(),

		DateOfBirth: func() time.Time {
			if user.DateOfBirth != nil {
				return *user.DateOfBirth
			}
			return currentUser.DateOfBirth
		}(),
	}

	if err := u.userRepo.UpdateUser(context, id, updated); err != nil {
		return entities.User{}, err
	}

	updatedUser, err := u.userRepo.GetUserByUsername(context, updated.Username)
	if err != nil {
		return entities.User{}, err
	}

	return updatedUser, nil
}

func (u *UserServiceStruct) DeleteUser(context context.Context, id int) error {
	if err := u.userRepo.DeleteUser(context, id); err != nil {
		return err
	}

	return nil
}
