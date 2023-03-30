package main

import (
	"log"
	"os"

	"github.com/arekbor/oauth/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal(err)
		return
	}

	a := api.New(os.Getenv("API_ADDR"))
	err = a.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
}
