package services

import (
	"fmt"
	"template/repositories"
)

type Services struct {
	AuthServices IAuthService
	CredServices ICredentialService
	UserServices IUserService
}

func ServiceBundler(repo *repositories.Repositories) Services {
	fmt.Println("repo: ", repo)
	// credSvc := CredentialServiceInit(repo.CredentialRepository)
	userSvc := UserServiceInit(repo.UserRepository)

	return Services{
		AuthServices: AuthServiceInit(userSvc),
		// CredServices: credSvc,
		UserServices: userSvc,
	}
}
