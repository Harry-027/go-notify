package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/Harry-027/go-notify/api-server/models"
)

func GetClientById(id uint) (models.Client, error) {
	var client models.Client
	DB.Where("id = ?", id).Find(&client)
	if client.ID == 0 {
		return models.Client{}, errors.New("an error occurred - client details not found")
	}
	return client, nil
}

func GetClientsByUserId(id int) ([]models.Client, error) {
	var clients []models.Client
	dbc := DB.Where("user_id = ?", id).Find(&clients)
	if dbc.Error != nil {
		log.Println("An error occurred while fetching clients for a given userId :: ", dbc.Error.Error())
		return []models.Client{}, dbc.Error
	}
	return clients, nil
}

func GetClientByIdUserId(clientId string, userId int) (models.Client, error) {
	var client models.Client
	dbc := DB.Where("id = ? AND user_id = ?", clientId, userId).Find(&client)
	if dbc.Error != nil {
		log.Println("An error occurred while fetching client details :: ", dbc.Error.Error())
		return models.Client{}, dbc.Error
	}
	return client, nil
}

func GetClientByMailId(mailId string) (models.Client, error) {
	var client models.Client
	DB.Where("mail_id = ?", mailId).Find(&client)
	if client.ID == 0 {
		return models.Client{}, errors.New("an error occurred - client details not found")
	}
	return client, nil
}

func GetClientByUserIdMailId(clientMailId string, userId int) (models.Client, error) {
	var client models.Client
	dbc := DB.Where("mail_id = ? AND user_id = ?", clientMailId, userId).Find(&client)
	if dbc.Error != nil {
		log.Println("An error occurred :: ", dbc.Error.Error())
		return models.Client{}, dbc.Error
	}
	return client, nil
}

func DeleteClientById(id uint) error {
	dbc := DB.Where("id = ?", id).Unscoped().Delete(&models.Client{})
	if dbc.Error != nil {
		log.Println("An error occurred while transaction ::", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func UpdateClientById(id uint, clientDetails models.Client) error {
	dbc := DB.Model(&models.Client{}).Where("id = ?", id).Updates(clientDetails)
	if dbc.Error != nil {
		log.Println("An error occurred :: ", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func GetClientFields(id uint) ([]models.Field, error) {
	var fieldDetails []models.Field
	DB.Where("client_id = ?", id).Find(&fieldDetails)
	if len(fieldDetails) == 0 {
		return []models.Field{}, errors.New("an error occurred - client fields not found")
	}
	return fieldDetails, nil
}

func DeleteClientField(key string) error {
	dbc := DB.Where("key = ?", key).Delete(&models.Field{})
	if dbc.Error != nil {
		log.Println("An error occurred while field delete transaction :: ", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func AddClient(client models.Client) error {
	dbc := DB.Create(&client)
	if dbc.Error != nil {
		log.Println("An error occurred while saving client :: ", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func AddClientFields(field models.Field) error {
	dbc := DB.Create(&field)
	if dbc.Error != nil {
		fmt.Println("An error occurred while save transaction ::", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}
