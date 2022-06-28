package models

import "gorm.io/gorm"

type Template struct {
	gorm.Model
	Name    string `json:"name" gorm:"type:varchar(50);unique;not null"`
	Subject string `json:"subject" gorm:"type:varchar(50)"`
	Body    string `json:"body" gorm:"type:varchar(10000)"`
	UserID  uint
	User    User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
