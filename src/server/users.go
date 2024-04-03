package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name string `json:"name" firestore:"name"`
	Email string `json:"email" firestore:"email"`
	Sub string `json:"sub" firestore:"sub"`
	AppRoles []string `json:"app_roles" firestore:"app_roles"`
}
type UpdateUserRequest struct {
	Name string `json:"name" firestore:"name"`
	Email string `json:"email" firestore:"email"`
	Sub string `json:"sub" firestore:"sub"`
	AppRoles []string `json:"app_roles" firestore:"app_roles"`
}

func CreateUser(c *gin.Context) {

	ctr, err := getController(c)
	if err != nil {
		error := apiError{"Failure fetching request controller"}
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	var request CreateUserRequest
	err = c.ShouldBindJSON(request)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	createdUser, err := ctr.CreateUser(request.Name, request.Email, request.Sub, request.AppRoles)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error creating user: %v", err))
	}

	c.IndentedJSON(http.StatusOK, createdUser)
	return
}

func ListUsers(c *gin.Context) {

	ctr, err := getController(c)
	if err != nil {
		error := apiError{"Failure fetching request controller"}
		c.JSON(http.StatusInternalServerError, error)
		return
	}
	
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("page", "100")

	pageNumber, err := strconv.ParseInt(page, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid page number")
		return
	}
	limitNumber, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid limit number")
		return
	}

	
	users, err := ctr.ListUsers(int(pageNumber), int(limitNumber))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error listing users: %v", err))
	}

	c.IndentedJSON(http.StatusOK, users)
	return

}

func  GetUser(c *gin.Context) {

	ctr, err := getController(c)
	if err != nil {
		error := apiError{"Failure fetching request controller"}
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	userId := c.Param("id")

	user, err := ctr.GetUser(userId)
	if err != nil {
		message := fmt.Sprintf("Error retrieving user with id %s: %v", userId, err)
		errorResponse(c, http.StatusBadRequest, message)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
	return
}

func  UpdateUser(c *gin.Context) {

	ctr, err := getController(c)
	if err != nil {
		error := apiError{"Failure fetching request controller"}
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	var request UpdateUserRequest
	c.ShouldBindJSON(request)

	userId := c.Param("id")

	updatedUser, err :=  ctr.UpdateUser(userId, request.Name, request.Email, request.Sub, request.AppRoles)
	if err != nil {
		message := fmt.Sprintf("Error updating user with id %s: %v", userId, err) 
		errorResponse(c, http.StatusBadRequest, message)
		return
	}

	c.IndentedJSON(http.StatusOK, updatedUser)
	return 
}

func  DeleteUser(c *gin.Context) {
	
	ctr, err := getController(c)
	if err != nil {
		error := apiError{"Failure fetching request controller"}
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	userId := c.Param("id")

	deletedUser, err := ctr.DeleteUser(userId)
	if err != nil {
		message := fmt.Sprintf("Error deleting user with id %s: %v", userId, err)
		errorResponse := apiError{message}
		c.IndentedJSON(http.StatusBadRequest, errorResponse)
		return
	}

	c.IndentedJSON(http.StatusOK, deletedUser)
	return 
}
