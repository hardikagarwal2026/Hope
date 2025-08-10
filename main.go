package main

import (
	"fmt"
	"log"
	"hope/config"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hope Main")
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found")
	}

	dbconfig := config.GetDatabaseConfig()
	db, err := config.InitDatabase(dbconfig)
	if err != nil {
		log.Fatalf("Failed to Initialize Database: %v", err)
	}else{
		fmt.Println("Database Initialized Successfully")
	}

	fmt.Println(db)
	


	
}