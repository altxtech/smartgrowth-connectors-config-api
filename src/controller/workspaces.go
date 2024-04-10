package controller

import (
	"fmt"
	"time"
	"smartgrowth-connectors/configapi/model"
)

func (ctr *Controller) CreateWorkspace(name string, permissions []model.WorkspacePermission)  (model.Workspace, error) {
	
	var workspace model.Workspace

	// User is owner
	basePermission, err  := model.NewWorkspacePermission(ctr.User.Email, "owner")
	if err != nil {
		return workspace, fmt.Errorf("Error creating base permissions: %v", err)
	}


	permissions = append(permissions, basePermission)
	permissions = dedupePermissions(permissions)

	// Create the workspace and insert it into the database
	workspace, err = model.NewWorkspace(name, permissions)
	if err != nil {
		return workspace, fmt.Errorf("Error creating workspace: %v", err) 
	}

	workspace, err = ctr.db.InsertWorkspace(workspace)
	if err != nil {
		fmt.Errorf("Error inserting workspace into database: %v", err)
	}

	return workspace, nil
}


func dedupePermissions(perms []model.WorkspacePermission) ([]model.WorkspacePermission) {
	
	uniquePerms := map[string]model.WorkspacePermission{}

	for _, val := range perms {
		uniquePerms[val.Principal + val.Role] = val
	}

	v := []model.WorkspacePermission{}

	for _, val := range uniquePerms {
		v = append(v, val)
	}

	return v
}

func (ctr *Controller) ListWorkspaces(offset int, limit int) ([]model.Workspace, error){

	// TODO: Super Admins should be able to read all workspaces
	
	workspaces, err := ctr.db.ListWorkspacesForPrincipal(ctr.User.Email)
	if err != nil {
		return workspaces, fmt.Errorf("Error reading workspaces from database: %v")
	}

	return workspaces,  nil
}

func (ctr *Controller) ReadWorkspace(id string) (model.Workspace, error) {
	
	var result model.Workspace
	workspace, err := ctr.db.GetWorkspaceByID(id)
	if err != nil {
		return result, fmt.Errorf("Error reading workspace from database: %v", err)
	}

	// Check dedupePermissions
	if !workspace.ViewableBy(ctr.User.Email) {
		return result, fmt.Errorf("User does not have permission to view workspace")
	}

	return workspace, nil
}

func (ctr *Controller) UpdateWorkspace(id string, name string, permissions []model.WorkspacePermission)  (model.Workspace, error) {
	
	var workspace model.Workspace

	// Check permissons
	// Since the only thing that really matters about editing in workspace are it's permissions, then we can define that 
	// only workspace admins can use this method
	workspace, err := ctr.db.GetWorkspaceByID(id)
	if err != nil {
		return workspace, fmt.Errorf("Error reading workspace from database: %v", err)
	}
	if !workspace.EditableBy(ctr.User.Email) {
		return workspace, fmt.Errorf("User can't edit this workspace")
	}

	// User is owner
	basePermission, err  := model.NewWorkspacePermission(ctr.User.Email, "owner")
	if err != nil {
		return workspace, fmt.Errorf("Error creating base permissions: %v", err)
	}


	permissions = append(permissions, basePermission)
	permissions = dedupePermissions(permissions)

	// Create the workspace and insert it into the database
	workspace.Permissions = permissions
	workspace.UpdatedAt = time.Now()

	workspace, err = ctr.db.UpdateWorkspace(workspace)
	if err != nil {
		fmt.Errorf("Error inserting workspace into database: %v", err)
	}

	return workspace, nil
}

func (ctr *Controller) DeleteWorkspace(id string) (model.Workspace, error) {
	
	// Read the workspace, check if it is existing
	workspace, err := ctr.db.GetWorkspaceByID(id)
	if err != nil {
		return workspace, fmt.Errorf("Error reading workspace from database: %v", err)
	}

	// Check permissions
	if !workspace.EditableBy(ctr.User.Email) {
		return workspace, fmt.Errorf("User does not have permission to delete workspace")
	}

	deletedWorkspace, err := ctr.db.DeleteWorkspaceByID(id)
	if err != nil {
		return deletedWorkspace, fmt.Errof("Error deleting workspace from database: %v")
	}

	return deletedWorkspace, nil
}
