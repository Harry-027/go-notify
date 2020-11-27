package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Harry-027/go-notify/api-server/config"
	"github.com/Harry-027/go-notify/api-server/database"
	"github.com/Harry-027/go-notify/api-server/models"
	"github.com/Harry-027/go-notify/api-server/repository"
	"github.com/Harry-027/go-notify/api-server/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
	"log"
	"strconv"
	"strings"
)

var KafkaConn *kafka.Conn

// SendMail godoc
// @Summary Send mail to client
// @Description Sends the mail to the registered client
// @Tags Send Mail
// @Accept json
// @Produce json
// @Param Body body []models.SendMailInput true "Send Mail Input"
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 204 {object} config.ApiResponse{}
// @Failure 417 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/sendMail [post]
func SendMail(ctx *fiber.Ctx) error {
	var user models.User
	var input []models.SendMailInput
	var processedIds []string

	userIdentifier := int(ctx.Locals("user_id").(float64))
	log.Println("user id :: ", userIdentifier)

	err := ctx.BodyParser(&input)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	user, err = repository.GetUserById(userIdentifier)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}

	if (user.NotificationCounter == 0) || (len(input) > user.NotificationCounter) {
		return utils.ApiResponse(ctx, fiber.StatusExpectationFailed, config.ApiConst[config.SUBS_EXPIRED])
	}

	for _, detail := range input {

		var template models.Template
		var clientRecord models.Client
		var clientDetails []models.Field
		var err error

		log.Println("User Id :: ", userIdentifier)
		log.Println("Template Id :: ", detail.TemplateId)
		template, err = repository.GetTemplateById(detail.TemplateId, userIdentifier)
		if err != nil {
			msg := fmt.Sprintf("Template with templateId %s not found !!", strconv.FormatUint(uint64(detail.TemplateId), 10))
			resp := fiber.Map{"status": "Ok", "message": msg}
			return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusNoContent, resp)
		}
		log.Println("template details :: ", template)

		clientRecord, err = repository.GetClientById(detail.ClientId)
		if err != nil {
			msg := fmt.Sprintf("Client with clientId %s not found !!", strconv.FormatUint(uint64(detail.ClientId), 10))
			resp := fiber.Map{"status": "Ok", "message": msg}
			return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusNoContent, resp)
		}
		log.Println("ClientDetails :: ", clientRecord)

		clientDetails, _ = repository.GetClientFields(detail.ClientId)
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
			To:      clientRecord.MailId,
			From:    user.Name,
			Subject: subject,
			Text:    bodyContent,
		}

		log.Println("kafka Payload :: ", kafkaPayload)
		payload, err := json.Marshal(kafkaPayload)
		if err != nil {
			return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
		}
		_, err = KafkaConn.WriteMessages(
			kafka.Message{Value: payload},
		)
		if err != nil {
			log.Println("failed to write messages:", err)
		}
		processedIds = append(processedIds, strconv.FormatUint(uint64(detail.ClientId), 10))
	}
	db := database.DB
	count := len(processedIds)
	user.NotificationCounter = user.NotificationCounter - count
	db.Save(&user)
	msg := fmt.Sprintf("Processed clientIds: %s", strings.Join(processedIds, ","))
	resp := fiber.Map{"status": "OK", "message": msg}
	return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusOK, resp)
}

// ScheduleMail godoc
// @Summary schedule the mail for clients
// @Description schedule the mail for clients
// @Tags Schedule Mail
// @Accept json
// @Produce json
// @Param Body body []models.SendMailInput true "Send Mail Input"
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 204 {object} config.ApiResponse{}
// @Failure 417 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/scheduleMail [post]
func ScheduleMail(ctx *fiber.Ctx) error {
	var input []models.SendMailInput
	var processedIds []string
	userIdentifier := int(ctx.Locals("user_id").(float64))
	log.Println("user id :: ", userIdentifier)

	err := ctx.BodyParser(&input)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	user, err := repository.GetUserById(userIdentifier)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}

	if (user.NotificationCounter == 0) || (len(input) > user.NotificationCounter) {
		return utils.ApiResponse(ctx, fiber.StatusExpectationFailed, config.ApiConst[config.SUBS_EXPIRED])
	}

	for _, detail := range input {

		var template models.Template
		var clientRecord models.Client
		var err error
		log.Println("User Id :: ", userIdentifier)
		log.Println("Template Id :: ", detail.TemplateId)
		template, err = repository.GetTemplateById(detail.TemplateId, userIdentifier)
		if err != nil {
			msg := fmt.Sprintf("Template with templateId %s not found !!", strconv.FormatUint(uint64(detail.TemplateId), 10))
			resp := fiber.Map{"status": "Ok", "message": msg}
			return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusNoContent, resp)
		}
		log.Println("template details :: ", template)

		clientRecord, err = repository.GetClientByIdUserId(strconv.FormatUint(uint64(detail.ClientId), 10), userIdentifier)
		if err != nil {
			msg := fmt.Sprintf("Client with clientId %s not found !!", strconv.FormatUint(uint64(detail.ClientId), 10))
			resp := fiber.Map{"status": "Ok", "message": msg}
			return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusNoContent, resp)
		}
		log.Println("ClientDetails :: ", clientRecord)
		prefType := utils.GetType(clientRecord.Preference)
		job := models.Job{
			Type:       prefType,
			Status:     "PENDING",
			TemplateID: template.ID,
			To:         clientRecord.MailId,
			From:       user.Name,
			ClientID:   clientRecord.ID,
		}
		err = repository.ScheduleJob(job)
		if err != nil {
			return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
		}
		processedIds = append(processedIds, strconv.FormatUint(uint64(detail.ClientId), 10))
	}
	msg := fmt.Sprintf("Processed clientIds for scheduled mail: %s", strings.Join(processedIds, ","))
	resp := fiber.Map{"status": "OK", "message": msg}
	return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusOK, resp)
}

// DeleteScheduleMail godoc
// @Summary Deletes the scheduled job
// @Description Deletes the scheduled job for sending mail
// @Tags Delete Scheduled Job
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Param Body body models.DeleteJob true "Delete job Input"
// @Success 200 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 401 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/deleteScheduleMail [post]
func DeleteScheduleMail(ctx *fiber.Ctx) error {
	var input models.DeleteJob
	userIdentifier := int(ctx.Locals("user_id").(float64))

	err := ctx.BodyParser(&input)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	job, err := repository.GetJob(input.JobID)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}

	client, err := repository.GetClientById(job.ClientID)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}

	user, err := repository.GetUserById(int(client.UserID))
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}

	if userIdentifier != int(user.ID) {
		return utils.ApiResponse(ctx, fiber.StatusUnauthorized, config.ApiConst[config.UNAUTHORIZED_ACTION])
	}

	err = repository.DeleteJob(job.ID)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}
	return utils.ApiResponse(ctx, fiber.StatusOK, config.ApiConst[config.JOB_DEL])
}
