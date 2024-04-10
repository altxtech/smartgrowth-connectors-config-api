package server

import (
	"fmt"
	"net/http"
	"strconv"
	"smartgrowth-connectors/configapi/model"

	"github.com/gin-gonic/gin"
)

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
	Permissions []model.WorkspacePermission `json:"permissions"`
}

type UpdateWorkspaceRequest struct {
	Name string `json:"name"`

	// I'm using the model struct because these will be passed as arguments downstream to a controller function
	// and I'm ok with having controller methods depending on types defines in the model package, but not with
	// types defined in the server package, which is supposed to be more of an "upstream" package

	Permissions []model.WorkspacePermission `json:"permissions"`
}

func CreateWorkspace(c *gin.Context) {
	ctr, err := getController(c)
	if err != nil {
		error := apiError{"Failure fetching request controller"}
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	var request CreateWorkspaceRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	createdWorkspace, err := ctr.CreateWorkspace(request.Name, request.Permissions)
	if err != nil {
		message := fmt.Sprintf("Error creating workspace: %v", err)
		response := apiError{ message }
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, createdWorkspace)
	return
}

func ListWorkspaces(c *gin.Context) { 
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

	workspaces, err := ctr.ListWorkspaces(int(pageNumber), int(limitNumber))
	if err != nil {
		message := fmt.Sprintf("Error listing workspaces: %v", err)
		response := apiError{ message } 
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, workspaces)
	return
}

func GetWorkspace(c *gin.Context) {
	ctr, err := getController(c)
	if err != nil {
		error := apiError{"Failure fetching request controller"}
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	id := c.Param("id")
	workspace, err := ctr.ReadWorkspace(id)
	if err != nil {
		message := fmt.Sprintf("Error reading workspace: %v", err)
		c.JSON(http.StatusBadRequest, apiError{ message })
		return
	}

	c.JSON(http.StatusOK, workspace)
	return
}

func UpdateWorkspace(c *gin.Context) {
	ctr, err := getController(c)
	if err != nil {
		error := apiError{"Failure fetching request controller"}
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	var request UpdateWorkspaceRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	id := c.Param("id")
	updatedWorkspace, err := ctr.UpdateWorkspace(id, request.Name, request.Permissions)
	if err != nil {
		message := fmt.Sprintf("Error updating workspace: %v", err)
		response := apiError{ message }
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, updatedWorkspace)
	return
}
func DeleteWorkspace(c *gin.Context) {
	ctr, err := getController(c)
	if err != nil {
		error := apiError{"Failure fetching request controller"}
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	id := c.Param("id")
	workspace, err := ctr.DeleteWorkspace(id)
	if err != nil {
		message := fmt.Sprintf("Error deleting workspace: %v", err)
		c.JSON(http.StatusBadRequest, apiError{ message })
		return
	}

	c.JSON(http.StatusOK, workspace)
	return
}
