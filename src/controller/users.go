package controller

import (
	"errors"
	"fmt"
	"smartgrowth-connectors/configapi/model"
)


func (ctr *Controller) CreateUser(name string, email string, subject string, appRole string) (model.User, error) {

	var idUser model.User

	// Authorization
	switch ctr.User.AppRole {
	case "Super Admin":
		// Just move on
	case "Client App":
		if appRole != "Customer"{
			return idUser, errors.New("Client Apps can't create users with roles other than \"Customer\"")
		}
	case "Customer":
		return idUser, errors.New("\"Customer\" users can't permform this action")
	default:
		return idUser, errors.New("Invalid AppRole")
	}

	// Create new User and insert in the database
	newUser := model.NewUser( name, email, subject, appRole)
	idUser, err := ctr.db.InsertUser( newUser )
	if err != nil {
		return idUser, fmt.Errorf("Error inserting user to database: %v", err)
	}

	return idUser, nil
}

func (cont *Controller) ListUsers(page int, limit int) ([]model.User, error) {

	var usersPage []model.User 
	
	// Authorization. Super Admins and Client Apps can list all users. Customers can only list their own user.
	switch cont.User.AppRole {
	case "Super Admin":
		// Just move on
	case "Client App":
		// Just move on
	case "Customer":
		user, err := cont.db.GetUserById(cont.User.ID)
		if err != nil {
			return usersPage, fmt.Errorf("Error getting user from database: %v", err)
		}
		usersPage = append(usersPage, user)
		return usersPage, nil
	default:
		return usersPage, errors.New("Invalid AppRole")
	}

	// Get users from database
	usersPage, err := cont.db.ListUsers(page, limit)
	if err != nil {
		return usersPage, fmt.Errorf("Error getting users from database: %v", err)
	}

	return usersPage, nil
}
func (cont *Controller) GetUser(userId string) (model.User, error) {

	var user model.User

	// Authorization.  Customer can only get their own user.
	if cont.User.AppRole == "Customer" && cont.User.ID != userId {
		return user, errors.New("Customers can only get their own user")
	}

	// Get user from database
	user, err := cont.db.GetUserById(userId)
	if err != nil {
		return user, fmt.Errorf("Error getting user from database: %v", err)
	}

	return user, nil
}

func (cont *Controller) UpdateUser(userId string, name string, email string, subject string, appRole string) (model.User, error) {

	var createdUser model.User

	// Authorization.  Customer can only update their own user.
	switch cont.User.AppRole {
	case "Super Admin":
		// Just move on
	case "Client App":
		// Just move on
	case "Customer":
		if cont.User.ID != userId {
			return createdUser, errors.New("Customers can only update their own user")
		}
	default:
		return createdUser, errors.New("Invalid AppRole")
	}

	// Update user in database
	updatedUser := model.NewUser(name, email, subject, appRole)
	createdUser, err := cont.db.UpdateUser(userId, updatedUser)
	if err != nil {
		return createdUser, fmt.Errorf("Error updating user in database: %v", err)
	}

	return createdUser, nil
}

func (cont *Controller) DeleteUser(userId string) (model.User, error) {

	var deletedUser model.User

	// Authorization.  Customer can only delete their own user.
	switch cont.User.AppRole {
	case "Super Admin":
		// Just move on
	case "Client App":
		// Just move on
	case "Customer":
		if cont.User.ID != userId {
			return deletedUser, errors.New("Customers can only delete their own user")
		}
	default:
		return deletedUser, errors.New("Invalid AppRole")
	}
	
	// Delete user from database
	deletedUser, err := cont.db.DeleteUserById(userId)
	if err != nil {
		return deletedUser, fmt.Errorf("Error deleting user from database: %v", err)
	}

	return deletedUser, nil
}
