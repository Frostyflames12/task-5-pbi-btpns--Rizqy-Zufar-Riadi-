package initializers

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var test string

func ConnectToDb() {

	var err error

	test = "what"

	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Fail to connect to Database")
	}
}
