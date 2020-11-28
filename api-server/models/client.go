package models

import (
	"github.com/jinzhu/gorm"
)

type VariableTemplateInput struct {
	Key          string `json:"key" valid:"required"`
	Value        string `json:"value" valid:"required"`
	ClientMailId string `json:"clientMailID" valid:"email,required"`
}

type DeleteTemplateInput struct {
	Key          string `json:"key" valid:"required"`
	ClientMailId string `json:"clientMailID" valid:"email,required"`
}

type SendMailInput struct {
	TemplateId uint `json:"templateId"`
	ClientId   uint `json:"clientId"`
}

type Client struct {
	gorm.Model
	Name       string `json:"name" gorm:"type:varchar(50)"`
	MailId     string `json:"mailID" valid:"email,required" gorm:"type:varchar(50);not null;unique_index:idx_first_second"`
	Phone      int    `json:"phone" valid:"required" gorm:"type:bigint;not null"`
	Preference string `json:"preference" gorm:"type:varchar(50)"`
	UserID     uint   `gorm:"unique_index:idx_first_second"`
	User       User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Field struct {
	gorm.Model
	Key      string `json:"key" valid:"required" gorm:"type:varchar(50);not null;unique_index:idx_field"`
	Value    string `json:"value" valid:"required" gorm:"type:varchar(50);not null"`
	ClientID uint   `gorm:"unique_index:idx_field"`
	Client   Client `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type KafkaPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Text    string `json:"text"`
	Subject string `json:"subject"`
}
