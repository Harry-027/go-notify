package models

import (
	"github.com/jinzhu/gorm"
)

type LoginInput struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"required"`
}

type SignupInput struct {
	Email           string `json:"email" valid:"email"`
	Role            string `json:"role" valid:"required"`
	Password        string `json:"password" valid:"required"`
	ConfirmPassword string `json:"confirm_password" valid:"required"`
}

type UpdatePassword struct {
	Password string `json:"password" valid:"required"`
}

type ForgotPassword struct {
	Email string `json:"email" valid:"email"`
}

type NewPassword struct {
	Password string `json:"password"`
}

type ForgotPasswordPayload struct {
	To      string
	From    string
	Text    string
	Subject string
}

type SubscriptionInput struct {
	SubscriptionType string `json:"subscriptionType"`
	PaymentType      string `json:"paymentType"`
}

type User struct {
	gorm.Model
	Name                string `json:"email" gorm:"unique;not null"`
	Password            string `json:"password" gorm:"not null"`
	Role                string `json:"role" gorm:"not null"`
	Subscription        string `gorm:"type:text"`
	NotificationCounter int    `gorm:"type:int"`
}

type Auth struct {
	gorm.Model
	Uuid   string `gorm:"unique;not null"`
	UserID uint
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
