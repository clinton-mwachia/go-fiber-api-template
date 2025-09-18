package main

import (
	"fmt"

	"github.com/clinton-mwachia/go-fiber-api-template/config"
)

func main() {
	// load env
	config.Load()
	// connect DB
	config.ConnectDB()
	fmt.Println("Hello go fiver template")
}
