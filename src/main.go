package main

import (
	"log"
	"smartgrowth-connectors/configapi/controller"
	"smartgrowth-connectors/configapi/database"
	"smartgrowth-connectors/configapi/server"
)

/*
Main package should be responsible for initializing the app.

Using enviroment variables to initialize the database, configure the controller, start the app, etc..
*/

func main(){

	// Initialize database
	db, err := database.NewInMemoryDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}


	// Configure controller
	controller, err := controller.NewController(db, nil)
	if err != nil {
		log.Fatalf("Failed to initialize controller: %v", err)
	}
	
	server, err := server.NewServer(controller, "dev-sn5f570cadx2zciy", "https://api.connectors.smartgrowth.consulting")
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	server.Run()
}
