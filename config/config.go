package config

import (
	"fmt"
	"log"
	"os"

	"github.com/aysf/gojwt/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err2 := godotenv.Load("./.env")
	if err2 != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	var err error

	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	InitialMigration()
}

func InitialMigration() {
	DB.AutoMigrate(&models.Book{})
	DB.AutoMigrate(&models.User{})
}

func InitDBTest() {

	// err2 := godotenv.Load("./.env")
	// if err2 != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	os.Getenv("DB_USERNAME_TEST"),
	// 	os.Getenv("DB_PASSWORD_TEST"),
	// 	os.Getenv("DB_HOST_TEST"),
	// 	os.Getenv("DB_PORT_TEST"),
	// 	os.Getenv("DB_NAME_TEST"))

	dsn := "root:0123@tcp(127.0.0.1:3306)/userbook?charset=utf8&parseTime=True&loc=Local"

	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	InitMigrateTest()

}

func InitMigrateTest() {
	DB.Migrator().DropTable(&models.User{})
	DB.AutoMigrate(&models.User{})
}
