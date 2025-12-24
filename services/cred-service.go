package services

import (
	"context"
	"template/data/dtos"
	"template/data/entities"
	"template/repositories"
	"time"
)

type ICredentialService interface {
	GetAllCredentials(context.Context) ([]entities.Credential, error)
	CreateCredential(context.Context, entities.Credential) error
	GetCredentialById(context.Context, int) (entities.Credential, error)
	UpdateCredential(context.Context, int, dtos.UpdatedCredDto) error
	DeleteCredential(context.Context, int) error
}

type CredentialServiceStruct struct {
	credentialRepo repositories.ICredentialRepository
}

func CredentialServiceInit(repo repositories.ICredentialRepository) *CredentialServiceStruct {
	return &CredentialServiceStruct{credentialRepo: repo}
}

func (c *CredentialServiceStruct) GetAllCredentials(context context.Context) ([]entities.Credential, error) {
	users, err := c.credentialRepo.GetAllCredentials(context)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (c *CredentialServiceStruct) CreateCredential(context context.Context, cred entities.Credential) error {
	if err := c.credentialRepo.CreateCredential(context, cred); err != nil {
		return err
	}

	return nil
}

func (c *CredentialServiceStruct) GetCredentialById(context context.Context, id int) (entities.Credential, error) {
	cred, err := c.credentialRepo.GetCredentialById(context, id)
	if err != nil {
		return entities.Credential{}, err
	}

	return cred, nil
}

func (c *CredentialServiceStruct) UpdateCredential(context context.Context, userId int, cred dtos.UpdatedCredDto) error {
	currentCred, err := c.credentialRepo.GetCredentialById(context, userId)
	if err != nil {
		return err
	}

	updated := entities.Credential{
		AccessToken: func() string {
			if cred.AccessToken != nil {
				return *cred.AccessToken
			}
			return currentCred.AccessToken
		}(),

		RefreshToken: func() string {
			if cred.RefreshToken != nil {
				return *cred.RefreshToken
			}
			return currentCred.RefreshToken
		}(),

		ExpiresIn: func() time.Duration {
			if cred.ExpiresIn != nil {
				return *cred.ExpiresIn
			}
			return currentCred.ExpiresIn
		}(),

		UserId: func() int {
			if cred.UserId != nil {
				return *cred.UserId
			}
			return currentCred.UserId
		}(),
	}

	if err := c.credentialRepo.UpdateCredential(context, userId, updated); err != nil {
		return err
	}

	return nil
}

func (c *CredentialServiceStruct) DeleteCredential(context context.Context, id int) error {
	if err := c.credentialRepo.DeleteCredential(context, id); err != nil {
		return err
	}

	return nil
}
