package repository

import (
	"github.com/Harry-027/go-notify/api-server/models"
	"log"
)

func ScheduleJob(job models.Job) error {
	dbc := DB.Create(&job)
	if dbc.Error != nil {
		log.Println("An error occurred while scheduling job :: ", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func UpdateJobStatus(id uint, status string) error {
	dbc := DB.Model(&models.Job{}).Where("id = ?", id).Update("status", status)
	if dbc.Error != nil {
		log.Println("An error occurred while updating job :: ", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func GetJob(id uint) (models.Job, error) {
	var job models.Job
	dbc := DB.Where("id = ?", id).Find(&job)
	if dbc.Error != nil {
		log.Println("An error occurred while fetching clients for a given userId :: ", dbc.Error.Error())
		return models.Job{}, dbc.Error
	}
	return job, nil
}

func DeleteJob(id uint) error {
	dbc := DB.Model(&models.Job{}).Where("id = ?", id).Update("Status", "INACTIVE")
	if dbc.Error != nil {
		log.Println("An error occurred while transaction ::", dbc.Error.Error())
		return dbc.Error
	}
	return nil
}

func GetActiveJobs() ([]models.Job, error) {
	var jobs []models.Job
	dbc := DB.Where("status = ?", "ACTIVE").Find(&jobs)
	if dbc.Error != nil {
		log.Println("An error occurred while transaction ::", dbc.Error.Error())
		return []models.Job{}, dbc.Error
	}
	return jobs, nil
}

func GetPendingJobs() ([]models.Job, error) {
	var jobs []models.Job
	dbc := DB.Where("status = ?", "PENDING").Find(&jobs)
	if dbc.Error != nil {
		log.Println("An error occurred while transaction ::", dbc.Error.Error())
		return []models.Job{}, dbc.Error
	}
	return jobs, nil
}
