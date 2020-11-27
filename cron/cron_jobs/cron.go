package cron_jobs

import (
	"encoding/json"
	"fmt"
	"github.com/Harry-027/go-notify/api-server/models"
	"github.com/Harry-027/go-notify/api-server/repository"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
)

// Send the scheduled mails ...
func SendScheduledMail(job models.Job) error {
	log.Println("Sending scheduled mail ...")
	template, err := repository.GetTemplate(job.TemplateID)
	if err != nil {
		log.Println("An error occurred while fetching template: ", err.Error())
		return err
	}
	clientRecord, err := repository.GetClientById(job.ClientID)
	if err != nil {
		log.Println("An error occurred while fetching Client: ", err.Error())
		return err
	}
	log.Println("ClientDetails :: ", clientRecord)

	clientDetails, _ := repository.GetClientFields(job.ClientID)
	log.Println("Client Fields :: ", clientDetails)

	subject := template.Subject
	bodyContent := template.Body

	if len(clientDetails) > 0 {
		for _, field := range clientDetails {
			oldVal := fmt.Sprintf("{{ %s }}", field.Key)
			newVal := field.Value
			subject = strings.Replace(subject, oldVal, newVal, -1)
			bodyContent = strings.Replace(bodyContent, oldVal, newVal, -1)
		}
	}

	kafkaPayload := models.KafkaPayload{
		To:      job.To,
		From:    job.From,
		Subject: subject,
		Text:    bodyContent,
	}

	log.Println("kafka Payload :: ", kafkaPayload)
	payload, err := json.Marshal(kafkaPayload)
	if err != nil {
		log.Println("An error occurred while marshalling payload: ", err.Error())
		return err
	}
	_, err = KafkaConn.WriteMessages(
		kafka.Message{Value: payload},
	)
	if err != nil {
		log.Println("failed to write messages:", err)
		return err
	}
	user, _ := repository.GetUserByName(job.From)
	err = repository.UpdateUserCounter(user, 1)
	return err
}

// Job scheduler ...
func JobScheduler(jobs []models.Job) {
	log.Println("Job scheduler started ...")
	cronScheduler := GetNewCron()
	for _, job := range jobs {
		SetCronJob(cronScheduler, job, SendScheduledMail)
	}
	go cronScheduler.c.Start()
}

// Server startup job scheduler ...
func ScheduleJobOnServerStart() {
	log.Println("Server startup job scheduler initiated ...")
	jobs, err := repository.GetActiveJobs()
	if err != nil {
		log.Println("An error occurred while fetching jobs: ", err)
		return
	}
	log.Println("Server startup job scheduler started ...")
	JobScheduler(jobs)
}

// Schedule the daily active jobs ...
func ScheduleDailyJob() {
	log.Println("Daily job scheduler initiated ...")
	cronScheduler := GetNewCron()
	entryId, err := cronScheduler.c.AddFunc("@daily", func() {
		DailyJobScheduler()
	})
	if err != nil {
		log.Println("An error occurred while scheduling cron job", err)
	}
	log.Println("EntryId for Daily Cron Job :: ", entryId)
	go cronScheduler.c.Start()
}

// Daily job scheduler ...
func DailyJobScheduler() {
	log.Println("Daily Job Scheduler started ...")
	jobs, err := repository.GetPendingJobs()
	if err != nil {
		return
	}
	for _, job := range jobs {
		_ = repository.UpdateJobStatus(job.ID, "ACTIVE")
	}
	JobScheduler(jobs)
}
