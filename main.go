package main

import (
	"fmt"
	"hope/db"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hope Main")
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found")
	}

	dbconfig := db.GetDatabaseConfig()
	db, err := db.InitDatabase(dbconfig)
	if err != nil {
		log.Fatalf("Failed to Initialize Database: %v", err)
	}else{
		fmt.Println("Database Initialized Successfully")
	}

	fmt.Println(db)
	


	
}