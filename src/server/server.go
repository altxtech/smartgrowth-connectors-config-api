package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"smartgrowth-connectors/configapi/controller"
	"smartgrowth-connectors/configapi/middleware"
)

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
	server.router.Use(server.setUser)

	return server, nil
} 

func (s *Server) setUser(c *gin.Context) {

	sub := c.GetString("sub")
	userController, err := s.controller.AsUser(sub)
	if err != nil {
		message := fmt.Sprintf("Failed to set controller for subscription %s: %v", sub, err)
		response := apiError{message}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	c.Set("ctr", userController)
	c.Next()
}

func getController(c *gin.Context) (*controller.Controller, error) {

	ctr, ok := c.Get("ctr")
	if !ok {
		return nil, errors.New("Missing user controller from requests context")
	}
	controller, ok := ctr.(*controller.Controller)
	if !ok {
		return nil, errors.New("Controller is of wrong type")
	}
	
	return controller, nil
}

// Wrapper for handlers?

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
