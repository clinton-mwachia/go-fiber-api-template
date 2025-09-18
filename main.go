package main

import (
	"fmt"
	"log"

	"github.com/clinton-mwachia/go-fiber-api-template/config"
	"github.com/joho/godotenv"
)

func main() {
	// load env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system env")
	}
	config.ConnectDB()
	fmt.Println("Hello go fiver template")
}
