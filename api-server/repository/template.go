package repository

import (
	"errors"
	"github.com/Harry-027/go-notify/api-server/models"
	"log"
)

func GetTemplate(id uint) (models.Template, error) {
	var template models.Template
	DB.Where("id = ?", id).Find(&template)
	if template.ID == 0 {
		return models.Template{}, errors.New("An error occurred - template details not found")
	}
	return template, nil
}

func GetTemplateByUserId(userId int) ([]models.Template, error) {
	var templates []models.Template
	dbc := DB.Where("user_id = ?", userId).Find(&templates)
	if dbc.Error != nil {
		log.Println("An error occurred: ", dbc.Error.Error())
		return []models.Template{}, dbc.Error
	}
	return templates, nil
}

func GetTemplateById(id uint, userId int) (models.Template, error) {
	var template models.Template
	dbc := DB.Where("id = ? AND user_id = ?", id, userId).Find(&template)

	if dbc.Error != nil {
		log.Println("An error occurred: ", dbc.Error.Error())
		return models.Template{}, dbc.Error
	}

	if template.ID == 0 {
		return models.Template{}, errors.New("an error occurred - template details not found")
	}
	return template, nil
}

func AddTemplate(template models.Template) error {
	dbc := DB.Save(&template)
	if dbc.Error != nil {
		log.Println("An error occurred while update transaction", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func UpdateTemplate(id uint, template models.Template) error {
	dbc := DB.Model(&models.Template{}).Where("id = ?", id).Update(template)
	if dbc.Error != nil {
		log.Println("An error occurred :: ", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func DeleteTemplate(id uint) error {
	dbc := DB.Where("id = ?", id).Delete(&models.Template{})
	if dbc.Error != nil {
		log.Println("An error occurred while delete template transaction :: ", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}
