package main

import (
	"log"
	"tic-tac-toe/internal/di"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system env")
	}

	fx.New(
		di.Module,
	).Run()
}
