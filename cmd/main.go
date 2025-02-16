package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"merchio/internal/api/handler"
	"merchio/internal/api/router"
	repo "merchio/internal/repository/user"
	serv "merchio/internal/service/user"
	"net/http"
	"os"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	err := Load("/app/.env")
	if err != nil {
		log.Fatalf("Error loading .env file", err)
	}
	pgDSN := os.Getenv("POSTGRES_CONN")
	fmt.Println(pgDSN)

	repository := repo.NewRepository(ctx, pgDSN)
	service := serv.NewService(repository)
	h := handler.NewImplementation(service)
	r := router.NewRouter(h)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}

}
