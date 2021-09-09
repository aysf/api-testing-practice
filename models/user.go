package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `json:"firstName" form:"firstName"`
	LastName  string `json:"lastName" form:"lastName"`
	Email     string `json:"email" form:"email"`
	Password  string `json:"password" form:"password"`
	Token     string `json:"token" form:"token"`
}
