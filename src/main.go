package main

import (
	"os"
	"log"
	"smartgrowth-connectors/configapi/controller"
	"smartgrowth-connectors/configapi/database"
	"smartgrowth-connectors/configapi/server"
	"smartgrowth-connectors/configapi/scripts"
)

/*
Main package should be responsible for initializing the app.

Using enviroment variables to initialize the database, configure the controller, start the app, etc..
*/

func main(){

	// Initialize database
	// If in testing mode, use in memory database
	var db database.Database
	if os.Getenv("ENV") == "LOCAL" {
		var err error
		db, err = database.NewInMemoryDB()

		// Seed database with inital database
		err = scripts.SeedDatabase(db)
		if err != nil {
			log.Fatalf("Failed to seed in memory database: %v", err)
		}

		if err != nil {
			log.Fatalf("Failed to initialize in memory database: %v", err)
		}
	} else {
		log.Fatalf("Production features not implemented yet")
	}
		


	// Configure controller
	controller, err := controller.NewController(db, nil)
	if err != nil {
		log.Fatalf("Failed to initialize controller: %v", err)
	}
	
	server, err := server.NewServer(controller, os.Getenv("AUTH0_DOMAIN"), os.Getenv("AUTH0_IDENTIFIER"))
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	server.Run()
}
