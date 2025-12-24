package services

import (
	"context"
	"fmt"
	"template/data/dtos"
	"template/data/entities"
	"template/utils/auth"
	"time"
)

type IAuthService interface {
	Login(context.Context, dtos.LoginDto) (*entities.Credential, error)
	Register(context.Context, dtos.RegisterDto) error
	RefreshAccessToken(dtos.RefreshTokenDto) (*dtos.LoginResponseDto, error)
}

type AuthServiceStruct struct {
	// credSvc ICredentialService
	userSvc IUserService
}

func AuthServiceInit(userService IUserService) *AuthServiceStruct {
	return &AuthServiceStruct{userSvc: userService}
}

func (a *AuthServiceStruct) Login(context context.Context, registerDto dtos.LoginDto) (*entities.Credential, error) {
	user, err := a.userSvc.GetUserByUsername(context, registerDto.Username)
	if err != nil {
		return nil, fmt.Errorf("user does not exists.")
	}

	if user.Password != registerDto.Password {
		return nil, fmt.Errorf("credentials does not match.")
	}

	accessToken, err := auth.CreateToken(user.Id, user.Username, user.Roles, 3600)
	if err != nil {
		return nil, err
	}
	refreshToken, err := auth.CreateToken(user.Id, user.Username, user.Roles, 86400)
	if err != nil {
		return nil, err
	}

	creds := entities.Credential{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    time.Duration(time.Hour * 2400),
	}

	fmt.Println(creds)

	return &creds, nil
}

func (a *AuthServiceStruct) Register(context context.Context, registerDto dtos.RegisterDto) error {
	userExists, err := a.userSvc.GetUserByUsername(context, registerDto.Username)
	if userExists != nil || err != nil {
		fmt.Errorf("user already exists")
	}

	user := entities.User{
		Username:     registerDto.Username,
		Password:     registerDto.Password,
		Email:        registerDto.Email,
		PhoneNumber:  registerDto.PhoneNumber,
		Roles:        "User",
		DateOfBirth:  registerDto.DateOfBirth,
		CreatedAt:    time.Now(),
		LastModified: time.Now(),
	}

	err = a.userSvc.CreateUser(context, user)
	if err != nil {
		return err
	}

	// creds := entities.Credential{
	// 	UserId:       createdUser.Id,
	// 	AccessToken:  accessToken,
	// 	RefreshToken: refreshToken,
	// 	ExpiresIn:    3600,
	// }

	// err = a.credSvc.CreateCredential(context, creds)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (a *AuthServiceStruct) RefreshAccessToken(refreshToken dtos.RefreshTokenDto) (*dtos.LoginResponseDto, error) {
	fmt.Println("coming into refresh access token")
	claims, err := auth.VerfiyAccessToken(refreshToken.RefreshTokenDto)
	if err != nil {
		return nil, err
	}

	// fmt.Println(claims["exp"].(int64))

	isExpired := int64(claims["exp"].(float64)) > time.Now().Unix()
	userId := int(claims["user_id"].(float64))
	userName := claims["username"].(string)
	roles := claims["roles"].(string)

	fmt.Println(isExpired)
	if isExpired {
		return nil, fmt.Errorf("token expired")
	}

	accessToken, err := auth.CreateToken(userId, userName, roles, 3600)
	if err != nil {
		return nil, err
	}

	return &dtos.LoginResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.RefreshTokenDto,
		ExpiresIn:    3600,
	}, nil
}
