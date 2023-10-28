package controllers

import (
	"example/goAPI/initializers"
	"example/goAPI/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	//get Email/pass from req body

	var body struct {
		Username string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	if len(body.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be 6 or more characters",
		})

		return
	}

	// hash password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	//create user

	user := models.User{Username: body.Username, Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {

	//get email and password

	var body struct {
		Username string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fail to create token",
		})

		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	fmt.Print(user.(models.User).ID)
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})

}

func UploadImage(c *gin.Context) {
	user, _ := c.Get("user")
	userID := user.(models.User).ID

	fmt.Println("INI UPLOAD FILE")
	fmt.Print(user.(models.User).ID)

	var body struct {
		Caption  string
		PhotoUrl string
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filePath := filepath.Join("uploads", file.Filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	fileMetadata := models.File{

		Title:    file.Filename,
		Caption:  body.Caption,
		PhotoUrl: body.PhotoUrl,
		UserID:   userID,
	}

	if err := initializers.DB.Create(&fileMetadata).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}
}

func DeleteImage(c *gin.Context) {
	var file models.File

	user, _ := c.Get("user")
	userID := user.(models.User).ID

	id := c.Param("id")
	err := initializers.DB.Where("id = ?", id).First(&file).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	if userID != file.UserID {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	filePath := filepath.Join("uploads", file.Title)

	err = os.Remove(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from upload folder"})
		return
	}

	err = initializers.DB.Delete(&file).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File " + file.Title + " deleted successfully",
	})

}

func GetImage(c *gin.Context) {
	var file models.File

	id := c.Param("id")
	err := initializers.DB.Where("id = ?", id).First(&file).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	filePath := filepath.Join("uploads", file.Title)
	// Open the file
	fileData, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileData.Close()

	fileHeader := make([]byte, 512)
	_, err = fileData.Read(fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	fileContentType := http.DetectContentType(fileHeader)
	// Get the file info
	fileInfo, err := fileData.Stat()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Title))
	c.Header("Content-Type", fileContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.File(filePath)
}
