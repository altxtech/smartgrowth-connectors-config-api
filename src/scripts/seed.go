package scripts

import (
	"os"
	"smartgrowth-connectors/configapi/database"
	"smartgrowth-connectors/configapi/model"
)

func SeedDatabase(db database.Database) error {
	
	// Just a single super admin user so we can test the  API manually
	superAdmin := model.NewUser("Super Admin", "", os.Getenv("SUPER_ADMIN_EMAIL"), "Super Admin")
	
	_, err := db.InsertUser(superAdmin)
	if err != nil {
		return err
	}

	return nil
}

