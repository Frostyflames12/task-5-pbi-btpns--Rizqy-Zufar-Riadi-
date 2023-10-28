package main

import (
	"example/goAPI/controllers"
	"example/goAPI/initializers"
	"example/goAPI/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.Authentication, controllers.Validate)
	r.POST("/photos", middleware.Authentication, controllers.UploadImage)
	r.DELETE("/photos/:id", middleware.Authentication, controllers.DeleteImage)
	r.GET("/photos/:id", middleware.Authentication, controllers.GetImage)

	r.Run()
}
