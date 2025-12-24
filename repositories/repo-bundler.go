package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type Repositories struct {
	// CredentialRepository ICredentialRepository
	UserRepository IUserRepository
}

func RepositoryBundler(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		// CredentialRepository: CredentialRepoInit(db),
		UserRepository: UserRepoInit(db),
	}
}
