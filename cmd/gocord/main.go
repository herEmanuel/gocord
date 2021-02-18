package main

import (
	"github.com/herEmanuel/gocord/pkg/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error: could not load the environment variables, " + err.Error())
	}

	api.Initialize()
}
