package model

import (
	"fmt"
	"time"
)

// User
type User struct {
	ID string `json:"id" firestore:"id"` // "" for unidentified
	Name string `json:"name" firestore:"name"`
	Email string `json:"email" firestore:"email"`
	Sub string `json:"sub" firestore:"sub"`
	AppRole string `json:"app_role" firestore:"app_role"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

func NewUser(name string, email string, sub string, appRole string) User {
	// Creaates new user without an identity (attributed at datadabase insertion)
	return User{ "", name, email, sub, appRole, time.Now(), time.Now() }
}

func (u User) HasIdentity() bool {
	return u.ID != ""
}

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
			return workspace, fmt.Errorf("Invalid permission at index %i: %v", idx, err)
		}
	}

	return nil
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
	if perm.Validate() != nil {
		return perm, fmt.Errorf("Invalid permission: %v", err)
	}

	return perm, nil
}

func (p WorkspacePermission) Validate() err {

	// This validation methods exists because Workspace permissions won't always be created by
	// the constructor method. Sometimes, they will be marshalled from JSON.
	// In that case, we'll need to parse them first and then validate it
	
	// Checks if the role is valid
	if role != "viewer" && role != "editor" && role != "owner" {
		return fmt.Errorf("Invalid role %s. Valid roles are \"viewer\", \"editor\" and  \"owner\""))
	}

	// NOTE: I thought about addind a validation to check if user principal is existing BUT
	// this method should be database agnostic. So this check should be done at db insert level

	return nil
}

type Source struct {
	ID string `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"` 
}

type Destination struct {
}
