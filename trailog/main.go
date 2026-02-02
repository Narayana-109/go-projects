package main

import (
	"fmt"
	"trailog/api"
	"trailog/database"
)

func main() {
	fmt.Print("Started execution...")
	database.Connect()
	api.LoadRouter()
	defer database.DB.Close()

}
