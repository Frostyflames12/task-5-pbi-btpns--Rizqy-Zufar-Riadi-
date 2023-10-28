package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `json:"id" binding:"required" gorm:"primaryKey"`
	Username string `json:"username" binding:"required" gorm:"not null;default:null"`
	Email    string `json:"email" binding:"required" gorm:"unique;not null;default:null"`
	Password string `json:"password" binding:"required" gorm:"not null;default:null"`
}
