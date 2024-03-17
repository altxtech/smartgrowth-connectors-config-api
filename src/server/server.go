package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"smartgrowth-connectors/configapi/controller"
	"smartgrowth-connectors/configapi/middleware"
)

/*
	What do I need, and how should I organize them?

	- API (Server)
	The gin router, that is responsible for handle the network specific aspects of the application.
	For example, if I only want to offer a JSON base REST HTTP API, this component should be responsible
	for decoding the HTTP requests and calling the controller actions

	- Controller
	Contains the business logic of the application. Modifies the data model based on the called actions

	**Problem**
	If the API uses the controller, to create the API router we need, at some point, pass the controller
	as an argument for a function that returns gin.Handler functions that use the controller.

	This means that there are two approached that I could take
	1. The API is a separate package that offers an API type with a NewAPI() constructor that takes a Controller as an argument
	2. There is no API package. Instead the controller has some sort of .MakeAPI() method that either return a *gin.Router() or outright starts the API server

	- Database
	Database inferface for the App. Should be passed as an argument for the controller constructions. Should be used exclusivelly for the constructor.

	- Types / Model
	Has the types definitions to be used in the application. The database package should exclusively use these types. Controller and `API` (if existant)
	couls have their own private types, if necessary.

	For example, the API could have private types for the purpose of ingesting incoming HTTP requests and returning HTTP reponses (again, always within
	the intended scope for the package of abstracting away the networking elements).

	Ohter example is that the types defined by the models package could have functionality
	methods that could be used in the controller. However the required functionality may required data that exceeds what is presented in the data model.
	In that case we may need create other types within the contoller packages that are composed by and extend the model types.
*/

type Server struct {
	router *gin.Engine
	controller *controller.Controller
}

func NewServer(controller *controller.Controller, authDomain string, identifier string) (*Server, error){

	server := &Server{
		router: gin.Default(),
		controller: controller,
	}


	// Add authentication middleware
	authMiddleware := middleware.EnsureValidToken(authDomain, identifier)
	server.router.Use(authMiddleware)


	server.router.GET("/sources", server.getSources())


	return server, nil
}

// Handlers
// Configuration
func (s *Server) getSources() gin.HandlerFunc {
	return func (c *gin.Context) {

		// List all configurations
		configs, err := s.controller.ListSources()
		if err != nil {
			errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Failed to list configs: %v", err))
			return
		}

		c.IndentedJSON(http.StatusOK, configs)
		return
	}
}


func (s *Server) Run() {
	fmt.Println("Starting server...")	
	s.router.Run()
}


// API Errors
type apiError struct {
	Error string `json:"error"`
}

func errorResponse(c *gin.Context, status int, message string) {
	response := apiError{ message }
	c.IndentedJSON(status, response) 
}
