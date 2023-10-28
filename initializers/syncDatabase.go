package initializers

import "example/goAPI/models"

func SyncDatabase() {
	// DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.User{}, &models.File{})
	// DB.AutoMigrate(&models.File{})
}
