package models

import "gorm.io/gorm"

type File struct {
	gorm.Model

	// ID       string `json:"id" gorm:"primaryKey"`
	Title    string `json:"title" gorm:"not null"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photourl"`
	UserID   uint   `json:"userid"`
}
