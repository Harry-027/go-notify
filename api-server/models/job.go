package models

import "github.com/jinzhu/gorm"

type DeleteJob struct {
	JobID uint `json:"jobId"`
}

type Job struct {
	gorm.Model
	Type       string // daily, weekly, monthly
	Status     string
	TemplateID uint
	To         string
	From       string
	ClientID   uint
	Client     Client `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
