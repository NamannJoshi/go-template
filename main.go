package main

import (
	"fmt"
	"template/db"
	"template/repositories"
	"template/router"
	"template/services"
)

func main() {
	fmt.Println("start adja")

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("PANIC: gg :", r)
	// 	}
	// }()

	conn, client := db.InitDatabase()
	db.CreateSchemas(conn)

	fmt.Println("adja", conn)

	repo := repositories.RepositoryBundler(conn)
	svc := services.ServiceBundler(repo)
	fmt.Println("ggagaglasglgsd", svc.CredServices)

	router.InitRouter(svc, client)
}
