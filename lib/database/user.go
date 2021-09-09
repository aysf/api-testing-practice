package database

import (
	"fmt"

	"github.com/aysf/gojwt/config"
	"github.com/aysf/gojwt/middlewares"
	"github.com/aysf/gojwt/models"
)

func GetUsers() (interface{}, error) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserById(userId int) (interface{}, error) {
	var user models.User

	if err := config.DB.Find(&user, userId).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(user *models.User) (interface{}, error) {
	tx := config.DB.Create(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func UpdateUser(id int, userUpdate *models.User) (interface{}, error) {
	user := models.User{}
	tx := config.DB.First(&user, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected > 0 {
		user.FirstName = userUpdate.FirstName
		user.LastName = userUpdate.LastName
		user.Email = userUpdate.Email
		user.Password = userUpdate.Password
	}
	config.DB.Save(&user)
	return user, nil
}

func DeleteUser(id int) (interface{}, error) {
	users := models.User{}
	tx := config.DB.Delete(&users, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}

func LoginUsers(user *models.User) (interface{}, error) {
	var err error
	if err := config.DB.Where("email = ? AND password = ?", user.Email, user.Password).First(user).Error; err != nil {
		fmt.Println("err 1")
		return nil, err
	}
	user.Token, err = middlewares.CreateToken(int(user.ID))
	if err != nil {
		fmt.Println("err 2")
		return nil, err
	}
	if err := config.DB.Save(user).Error; err != nil {
		fmt.Println("err 3")
		return nil, err
	}
	return user, nil
}
