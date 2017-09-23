package main

import (
	"log"
	"os"
)

func main() {
	database, e := NewDB(
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
	)
	if e != nil {
		log.Fatal("Could not access database; ", e)
	}

	service, err := NewService(database, "8080")
	if err != nil {
		log.Fatal("Could not start service; ", err)
	}

	err = service.Start()
	if err != nil {
		log.Fatal("Service stopped; ", err)
	}
}
