package controller

import "smartgrowth-connectors/configapi/model"


func (ctr *Controller) CreateUser(name string, email string, subject string, appRole string) (model.User, error) {

	var createdUser model.User

	// Check if user has permissions
	if !(ctr.User.AppRole == "??" || ctr.User.AppRole == "??" {
		return createdUser, errors.New("User doesn't have permission to perform this action")
	}

	return createdUser, nil
}

func (cont *Controller) ListUsers(page int, limit int) ([]model.User, error) {

	var usersPage []model.User 
	
	// TODO: Implement

	return usersPage, nil
}
func (cont *Controller) GetUser(userId string) (model.User, error) {

	var user model.User

	return user, nil
}
func (cont *Controller) UpdateUser(userId string, name string, email string, subject string, appRoles [] string) (model.User, error) {

	var createdUser model.User

	// TODO: Implement

	return createdUser, nil
}
func (cont *Controller) DeleteUser(userId string) (model.User, error) {

	var deletedUser model.User

	// TODO: Implement

	return deletedUser, nil
}
