package repository

import (
	"github.com/Harry-027/go-notify/api-server/models"
	"github.com/jinzhu/gorm"
	"log"
)

var DB *gorm.DB

func AddUser(user models.User) error {
	dbc := DB.Create(&user)
	if dbc.Error != nil {
		log.Println("An error occurred ::", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func GetUserByName(name string) (models.User, error) {
	var user models.User
	dbc := DB.Where("name = ?", name).Find(&user)
	if dbc.Error != nil {
		log.Println("An error occurred :: ", dbc.Error.Error())
		return models.User{}, dbc.Error
	}
	return user, nil
}

func GetUserById(id int) (models.User, error) {
	var user models.User
	dbc := DB.Where("id = ?", id).Find(&user)
	if dbc.Error != nil {
		log.Println("An error occurred :: ", dbc.Error.Error())
		return models.User{}, dbc.Error
	}
	return user, nil
}

func UpdateUserCounter(user models.User, count int) error {
	user.NotificationCounter = user.NotificationCounter - count
	dbc := DB.Save(&user)
	if dbc.Error != nil {
		log.Println("An error occurred while updating user counter ::", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func UpdateUserDetails(user models.User) error {
	dbc := DB.Save(&user)
	if dbc.Error != nil {
		log.Println("An error occurred while update transaction", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func UpdateUserPassword(userId int, hashPwd string) error {
	dbc := DB.Model(&models.User{}).Where("id = ?", userId).Update("password", hashPwd)
	if dbc.Error != nil {
		log.Println("An error occurred while updating password : ", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	dbc := DB.Find(&users)
	if dbc.Error != nil {
		log.Println("An error occurred while fetching all user details :: ", dbc.Error.Error())
		return []models.User{}, dbc.Error
	}
	return users, nil
}

func DeleteUser(userId int) error {
	dbc := DB.Where("id = ?", userId).Unscoped().Delete(&models.User{})
	if dbc.Error != nil {
		log.Println("An error occurred :: ", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func SaveUserAuth(auth models.Auth) error {
	dbc := DB.Create(&auth)
	if dbc.Error != nil {
		log.Println("An error occurred while saving login details", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func IfUsersCurrentSession(uuid interface{}) bool {
	var user models.User
	dbc := DB.Where("uuid = ?", uuid).Find(&user)
	if dbc.Error != nil {
		log.Println("An error occurred :: ", dbc.Error.Error())
		return false
	}
	if user.ID == 0 {
		return false
	}
	return true
}

func InvalidateCurrentSession(userID int, uuid interface{}) {
	dbc := DB.Where("user_id = ? AND uuid = ?", userID, uuid).Unscoped().Delete(&models.User{})
	if dbc.Error != nil {
		log.Println("An error occurred while deleting user session", dbc.Error.Error())
	}
}
