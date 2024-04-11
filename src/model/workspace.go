package model

import (
	"fmt"
	"time"
)

type Workspace struct {
	ID string `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
	Permissions []WorkspacePermission `json:"permissions" firestore:"permissions"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWorkspace(name string, perms []WorkspacePermission) (Workspace, error) {

	workspace := Workspace{ "", name, perms, time.Now(), time.Now()}

	// Validate permissions
	for idx, perm := range perms {
		err := perm.Validate()
		if err != nil {
			return workspace, fmt.Errorf("Invalid permission at index %d: %v", idx, err)
		}
	}

	return workspace, nil
}

func (w Workspace) ViewableBy(principal string) bool {

	// If principal has any valid role for workspace, they can view
	for _, perm := range w.Permissions {
		if perm.Principal == principal {
			err := perm.Validate()
			if err == nil {
				return true
			}
		}
	}
	return false
}
func (w Workspace) EditableBy(principal string) bool {

	// If principal has any valid role for workspace, they can view
	for _, perm := range w.Permissions {
		if perm.Principal == principal {
			err := perm.Validate()
			if err == nil {
				if perm.Role == "admin" {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}
	return false
}

type WorkspacePermission struct {
	Principal string `json:"user" firestore:"user"`
	Role string `json:"role" firestore:"role"` // "admin",  "editor", "viewer"
}

func NewWorkspacePermission(userEmail string, role string) (WorkspacePermission, error) {
	
	perm := WorkspacePermission{ userEmail, role }
	err := perm.Validate(); if err != nil {
		return perm, fmt.Errorf("Invalid permission: %v", err)
	}

	return perm, nil
}

func (p WorkspacePermission) Validate() error {

	// This validation methods exists because Workspace permissions won't always be created by
	// the constructor method. Sometimes, they will be marshalled from JSON.
	// In that case, we'll need to parse them first and then validate it
	
	// Checks if the role is valid
	if p.Role != "viewer" && p.Role != "editor" && p.Role != "owner" {
		return fmt.Errorf("Invalid role %s. Valid roles are \"viewer\", \"editor\" and  \"owner\"", p.Role)
	}

	// NOTE: I thought about addind a validation to check if user principal is existing BUT
	// this method should be database agnostic. So this check should be done at db insert level

	return nil
}
