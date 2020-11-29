package models

import "github.com/jinzhu/gorm"

type DeleteJob struct {
	JobID uint `json:"jobId"`
}

type Job struct {
	gorm.Model
	Type       string // daily, weekly, monthly
	Status     string
	TemplateID uint   `gorm:"unique_index:idx_job"`
	To         string `gorm:"unique_index:idx_job"`
	From       string
	ClientID   uint
	Client     Client `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Audit struct {
	gorm.Model
	To           string
	FromUser     uint `gorm:"not null"`
	TemplateName string
	TemplateID   uint
}
