package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func InitDatabase() (*pgxpool.Pool, *redis.Client) {
	godotenv.Load()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_DB_ADDRESS"),
		Username: "default",
		Password: os.Getenv("REDIS_DB_PASS"),
		DB:       0,
	})

	return conn, client
}

func CreateSchemas(db *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	createUserSchema(ctx, db)
	// createCredsSchema(ctx, db)
}

func createUserSchema(ctx context.Context, db *pgxpool.Pool) {
	query := `
		create table if not exists users (
			id bigserial primary key,
			username varchar(200) not null unique,
			password varchar(200) not null,
			roles text not null,
			email varchar(200) not null,
			birth_of_date timestamp null,
			phone_number int not null,
			created_at timestamptz not null default now(),
			last_modified_at timestamp not null default now()
		);
	`

	if _, err := db.Exec(ctx, query); err != nil {
		fmt.Println("error while creating user schema %v", err)
		return
	}
}

func createCredsSchema(ctx context.Context, db *pgxpool.Pool) {
	query := `
			create table if not exists creds (
			id bigserial primary key,
			user_id bigint not null,
			access_token text not null,
			refresh_token text not null,
			expires_at timestamptz not null,
			created_at timestamptz not null default now(),

			constraint fk_user_cred
				foreign key (user_id)
				references users(id)
				on delete cascade,

			constraint uq_user_cred unique (user_id)
		);
	`

	if _, err := db.Exec(ctx, query); err != nil {
		fmt.Println("error while creating creds schema %v", err)
		return
	}
}
