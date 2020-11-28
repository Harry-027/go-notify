package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Harry-027/go-notify/api-server/config"
	"github.com/Harry-027/go-notify/api-server/models"
	"github.com/Harry-027/go-notify/api-server/payment"
	"github.com/Harry-027/go-notify/api-server/repository"
	"github.com/Harry-027/go-notify/api-server/utils"
	"github.com/Harry-027/go-notify/api-server/validator"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"strconv"
	"strings"
)

// AddClients godoc
// @Summary Registers the client details
// @Description Adds the clients for a logged in user
// @Tags Add Client
// @Accept json
// @Produce json
// @Param Body body []models.SwaggerClient true "Add the client details"
// @Param Authorization header string true "Authentication header"
// @Success 201 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/clients [post]
func AddClients(ctx *fiber.Ctx) error {
	var clients []models.Client

	err := ctx.BodyParser(&clients)
	if err != nil {
		log.Println(" An error occurred while body parsing :: ", err)
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	isInvalid, err := validator.ValidateSliceOfStruct(clients)
	if err != nil || isInvalid {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	userIdentifier := uint(ctx.Locals("user_id").(float64))
	var dbSaveError []string
	var msg string
	var status string
	for _, client := range clients {
		isValidPref := validator.ValidatePreference(client.Preference)
		if isValidPref == false {
			err := fmt.Sprintf("Client preference is incorrect for: %s", client.Name)
			dbSaveError = append(dbSaveError, err)
			continue
		}
		client.UserID = userIdentifier
		err := repository.AddClient(client)
		if err != nil {
			log.Println("An error occurred while saving client :: ", err.Error())
			errMsg := fmt.Sprintf("An error occurred while saving client :: %s", client.Name)
			dbSaveError = append(dbSaveError, errMsg)
		}
		addedClient, _ := repository.GetClientByMailId(client.MailId)
		nameField := models.Field{
			Key:      "Name",
			Value:    client.Name,
			ClientID: addedClient.ID,
		}
		mailField := models.Field{
			Key:      "Mail",
			Value:    client.MailId,
			ClientID: addedClient.ID,
		}
		err = repository.AddClientFields(nameField)
		err = repository.AddClientFields(mailField)
		if err != nil {
			errMsg := fmt.Sprintf("An error occurred while saving client Field:: %s", client.Name)
			log.Println("An error occurred while saving client fields :: ", err.Error())
			dbSaveError = append(dbSaveError, errMsg)
		}
	}
	key := fmt.Sprintf("Clients:%d", userIdentifier)
	Cache.deleteKey(key)
	if len(dbSaveError) > 0 {
		msg = strings.Join(dbSaveError, ",")
		status = "Error"
	} else {
		status = "Ok"
		msg = "Clients saved successfully ..."
	}
	resp := fiber.Map{"status": status, "message": msg}
	return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusCreated, resp)
}

// AddTemplate godoc
// @Summary Add the mail template
// @Description Add the mail template
// @Tags Add template
// @Accept json
// @Produce json
// @Param Body body models.SwaggerTemplate true "Add the template details"
// @Param Authorization header string true "Authentication header"
// @Success 201 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 412 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/addTemplate [post]
func AddTemplate(ctx *fiber.Ctx) error {
	var template models.Template
	var user models.User

	err := ctx.BodyParser(&template)
	if err != nil {
		log.Println(" An error occurred while body parsing :: ", err)
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	userIdentifier := int(ctx.Locals("user_id").(float64))
	user, _ = repository.GetUserById(userIdentifier)
	if user.ID == 0 {
		log.Println("User doesn't exists !!")
		resp := fiber.Map{"status": "Error", "message": "Signup required before adding Templates!!"}
		return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusPreconditionFailed, resp)
	}

	template.UserID = user.ID
	err = repository.AddTemplate(template)
	if err != nil {
		log.Println("Error while saving template: ", err.Error())
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}
	key := fmt.Sprintf("Templates:%d", userIdentifier)
	Cache.deleteKey(key)
	resp := fiber.Map{"status": "Created", "message": "Template added successfully !!"}
	return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusCreated, resp)
}

// GetTemplates godoc
// @Summary Fetch the templates
// @Description Fetch all the registered mail templates
// @Tags Template details
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/getTemplates [get]
func GetTemplates(ctx *fiber.Ctx) error {
	var err error
	var result []fiber.Map
	var templates []models.Template
	userIdentifier := int(ctx.Locals("user_id").(float64))
	cacheKey := fmt.Sprintf("Templates:%d", userIdentifier)
	ifExists, _ := Cache.ifExistsInCache(cacheKey)
	if ifExists {
		result, err := Cache.getDetails(cacheKey)
		if err == nil {
			unmarshalErr := json.Unmarshal([]byte(result), &templates)
			if unmarshalErr != nil {
				log.Println("An error occurred while fetching template: ", unmarshalErr.Error())
				return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
			}
			log.Println("Serving templates from cache ...")
		}
	} else {
		templates, err = repository.GetTemplateByUserId(userIdentifier)
		if err != nil {
			log.Println("An error occurred while fetching template: ", err.Error())
			return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
		}
		byteTemplates, marshalErr := json.Marshal(templates)
		if marshalErr == nil {
			Cache.setDetails(cacheKey, string(byteTemplates))
		}
	}
	for _, template := range templates {
		result = append(result, fiber.Map{"name": template.Name, "subject": template.Subject, "body": template.Body})
	}
	return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusOK, result)
}

// AddUserVariable godoc
// @Summary Add the user variable
// @Description Add the user variable
// @Tags Add user variable
// @Accept json
// @Produce json
// @Param Body body models.VariableTemplateInput true "Add the template details"
// @Param Authorization header string true "Authentication header"
// @Success 201 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 424 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/clientDetails [post]
func AddUserVariable(ctx *fiber.Ctx) error {
	var tempInput models.VariableTemplateInput

	err := ctx.BodyParser(&tempInput)
	if err != nil {
		log.Println(" An error occurred while body parsing :: ", err)
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	validationErr, result := validator.Validate(tempInput)
	if validationErr != nil || result == false {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	var client models.Client
	client, err = repository.GetClientByMailId(tempInput.ClientMailId)
	if err != nil {
		resp := fiber.Map{"status": "Failed Dependency", "message": "Given Client not yet registered !!"}
		return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusFailedDependency, resp)
	}

	field := models.Field{
		Key:      tempInput.Key,
		Value:    tempInput.Value,
		ClientID: client.ID,
	}
	err = repository.AddClientFields(field)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}
	resp := fiber.Map{"status": "Created", "message": "Template variable added successfully !!"}
	return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusCreated, resp)
}

// Subscribe godoc
// @Summary User subscription
// @Description user can subscribe to a specific plan
// @Tags Subscribe
// @Accept json
// @Produce json
// @Param Body body models.SubscriptionInput true "subscription details"
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 424 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/subscribe [post]
func Subscribe(ctx *fiber.Ctx) error {
	var (
		subs models.SubscriptionInput
		user models.User
	)

	err := ctx.BodyParser(&subs)
	if err != nil {
		log.Println(" An error occurred while body parsing :: ", err)
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	userIdentifier := int(ctx.Locals("user_id").(float64))
	user, _ = repository.GetUserById(userIdentifier)
	if user.ID == 0 {
		resp := fiber.Map{"status": "Error", "message": "User doesn't exist"}
		return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusFailedDependency, resp)
	}

	tokens, err := payment.MakePayment(subs)
	if err != nil {
		log.Println("An error occurred while payment transaction", err)
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}

	user.NotificationCounter = user.NotificationCounter + tokens
	user.Subscription = subs.SubscriptionType
	err = repository.UpdateUserDetails(user)
	if err != nil {
		log.Println("An error occurred while updating user details: ", err.Error())
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}
	return utils.ApiResponse(ctx, fiber.StatusOK, config.ApiConst[config.SUBS_SUCCESS])
}

// GetSubsDetail godoc
// @Summary Current subscription details
// @Description Get the current subscription details
// @Tags Subscription details
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/subscriptionDetails [get]
func GetSubsDetail(ctx *fiber.Ctx) error {
	userIdentifier := int(ctx.Locals("user_id").(float64))
	usr, err := repository.GetUserById(userIdentifier)
	if err != nil {
		log.Println("An error occurred while fetching user details: ", err.Error())
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}
	msg := fmt.Sprintf("Mail Count left: %d", usr.NotificationCounter)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Ok", "message": msg})
}

// UpdateTemplate godoc
// @Summary Update the template
// @Description Update the template details
// @Tags Update Template
// @Accept json
// @Produce json
// @Param id path string true "template id"
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 204 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/updateTemplate [put]
func UpdateTemplate(ctx *fiber.Ctx) error {
	var template models.Template
	templateID, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	userIdentifier := int(ctx.Locals("user_id").(float64))

	template, _ = repository.GetTemplateById(uint(templateID), userIdentifier)
	if template.ID != 0 {
		return utils.ApiResponse(ctx, fiber.StatusNoContent, config.ApiConst[config.NO_CONTENT])
	}

	var templateDetails models.Template
	err = ctx.BodyParser(&templateDetails)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	err = repository.UpdateTemplate(template.ID, templateDetails)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Ok", "message": "Template details updated successfully"})
}

// UpdateClient godoc
// @Summary Update the client details
// @Description Update the client details
// @Tags Update Client
// @Accept json
// @Produce json
// @Param Body body models.SwaggerClient true "client details"
// @Param mailId path string true "mail id"
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 204 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/updateClient [put]
func UpdateClient(ctx *fiber.Ctx) error {
	var client models.Client
	clientMailID := ctx.Params("mailId")
	userIdentifier := int(ctx.Locals("user_id").(float64))
	client, _ = repository.GetClientByUserIdMailId(clientMailID, userIdentifier)
	if client.ID != 0 {
		var clientDetails models.Client
		err := ctx.BodyParser(&clientDetails)
		if err != nil {
			return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
		}

		err = repository.UpdateClientById(client.ID, clientDetails)
		if err != nil {
			return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Ok", "message": "Client details updated successfully"})
	}
	return utils.ApiResponse(ctx, fiber.StatusNoContent, config.ApiConst[config.NO_CONTENT])
}

// DeleteClient godoc
// @Summary Delete the client
// @Description Delete the client
// @Tags Delete Client
// @Accept json
// @Produce json
// @Param mailId path string true "mail id"
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 204 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/deleteClient [delete]
func DeleteClient(ctx *fiber.Ctx) error {
	var client models.Client
	clientMailID := ctx.Params("mailId")
	userIdentifier := int(ctx.Locals("user_id").(float64))
	client, _ = repository.GetClientByUserIdMailId(clientMailID, userIdentifier)
	if client.ID == 0 {
		return utils.ApiResponse(ctx, fiber.StatusNoContent, config.ApiConst[config.NO_CONTENT])
	}
	err := repository.DeleteClientById(client.ID)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Ok", "message": "Client removed successfully !!"})
}

// DeleteTemplate godoc
// @Summary Delete the template
// @Description Delete the template
// @Tags Delete Template
// @Accept json
// @Produce json
// @Param id path string true "mail id"
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 204 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/deleteTemplate [delete]
func DeleteTemplate(ctx *fiber.Ctx) error {
	var template models.Template
	name, _ := strconv.ParseUint(ctx.Params("id"), 10, 64)
	userIdentifier := int(ctx.Locals("user_id").(float64))
	template, _ = repository.GetTemplateById(uint(name), userIdentifier)
	if template.ID == 0 {
		return utils.ApiResponse(ctx, fiber.StatusNoContent, config.ApiConst[config.NO_CONTENT])
	}
	err := repository.DeleteTemplate(template.ID)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Ok", "message": "Template removed successfully !!"})
}

// DeleteUserVariable godoc
// @Summary Delete the user variable
// @Description Delete the user variable
// @Tags Delete user variable
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Param Body body models.DeleteTemplateInput true "delete user variable"
// @Success 200 {object} config.ApiResponse{}
// @Failure 204 {object} config.ApiResponse{}
// @Failure 403 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/clientDetails [delete]
func DeleteUserVariable(ctx *fiber.Ctx) error {
	var input models.DeleteTemplateInput
	var client models.Client
	var user models.User

	userIdentifier := uint(ctx.Locals("user_id").(float64))
	err := ctx.BodyParser(&input)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	validationErr, result := validator.Validate(input)
	if validationErr != nil || result == false {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	client, err = repository.GetClientByMailId(input.ClientMailId)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusNoContent, config.ApiConst[config.NO_CONTENT])
	}

	user, err = repository.GetUserById(int(client.UserID))
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusNoContent, config.ApiConst[config.NO_CONTENT])
	}

	if userIdentifier != user.ID {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "Forbidden", "message": "Client does not belong to you !!"})
	}

	err = repository.DeleteClientField(input.Key)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "Error", "message": "Error deleting given key !!"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Ok", "message": "Given key deleted successfully !!"})
}

// GetUsers godoc
// @Summary Fetch the users
// @Description Fetch the users registered in system
// @Tags Get User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 403 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/users [get]
func GetUsers(ctx *fiber.Ctx) error {
	userRole := ctx.Locals("user_role")
	if userRole == config.ADMIN_ROLE {
		var userDetails []fiber.Map
		var users []models.User
		var err error
		cacheKey := "USERS"
		ifExists, _ := Cache.ifExistsInCache(cacheKey)
		if ifExists {
			cachedUsers, err := Cache.getDetails(cacheKey)
			if err == nil {
				unmarshalErr := json.Unmarshal([]byte(cachedUsers), &users)
				if unmarshalErr != nil {
					log.Println("An error occurred while fetching details: ", unmarshalErr.Error())
					return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
				}
				log.Println("Serving users from cache ...")
			}
		} else {
			users, err = repository.GetAllUsers()
			if err != nil {
				return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
			}
			byteUsers, marshalErr := json.Marshal(users)
			if marshalErr == nil {
				Cache.setDetails(cacheKey, string(byteUsers))
			}
		}
		for _, user := range users {
			details := fiber.Map{"Email": user.Name, "Role": user.Role, "Subscription": user.Subscription, "NotificationCounter": user.NotificationCounter}
			userDetails = append(userDetails, details)
		}
		return ctx.Status(fiber.StatusOK).JSON(userDetails)
	}
	return utils.ApiResponse(ctx, fiber.StatusForbidden, config.ApiConst[config.FORBIDDEN])
}

// GetClientVariable godoc
// @Summary Fetch the client variable (variable that can be replaced in a given email template)
// @Description Fetch the client variable
// @Tags Fetch Client variable
// @Accept json
// @Produce json
// @Param clientID path string true "client ID"
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 403 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/clientDetails [get]
func GetClientVariable(ctx *fiber.Ctx) error {
	var clientVariable []models.Field

	clientID, _ := strconv.ParseUint(ctx.Params("clientID"), 10, 64)
	userIdentifier := int(ctx.Locals("user_id").(float64))

	client, _ := repository.GetClientByIdUserId(strconv.FormatUint(clientID, 10), userIdentifier)
	if client.ID == 0 {
		return utils.ApiResponse(ctx, fiber.StatusNoContent, config.ApiConst[config.NO_CONTENT])
	}
	clientVariable, _ = repository.GetClientFields(uint(clientID))
	if len(clientVariable) > 0 {
		var result []fiber.Map
		for _, variable := range clientVariable {
			data := fiber.Map{"Key": variable.Key, "Value": variable.Value}
			result = append(result, data)
		}
		return ctx.Status(fiber.StatusOK).JSON(result)
	}
	return utils.ApiResponse(ctx, fiber.StatusNoContent, config.ApiConst[config.NO_CONTENT])
}

// GetClients godoc
// @Summary Fetch the clients
// @Description Fetch the clients
// @Tags Fetch Client
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} []models.SwaggerClient{}
// @Failure 204 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /api/clients [get]
func GetClients(ctx *fiber.Ctx) error {
	var clientDetails []fiber.Map
	var clients []models.Client
	var err error
	userIdentifier := int(ctx.Locals("user_id").(float64))
	cacheKey := fmt.Sprintf("Clients:%d", userIdentifier)
	log.Println("cacheKey: ", cacheKey)
	ifExists, _ := Cache.ifExistsInCache(cacheKey)
	if ifExists {
		cachedClients, err := Cache.getDetails(cacheKey)
		log.Println("Cached clients: ", cachedClients)
		if err == nil {
			unmarshalErr := json.Unmarshal([]byte(cachedClients), &clients)
			if unmarshalErr != nil {
				log.Println("An error occurred while fetching clients: ", unmarshalErr.Error())
				return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
			}
			log.Println("Serving clients from cache ...")
		}
	} else {
		clients, err = repository.GetClientsByUserId(userIdentifier)
		if err != nil {
			return utils.ApiResponse(ctx, fiber.StatusNoContent, config.ApiConst[config.NO_CONTENT])
		}
		byteClients, marshalErr := json.Marshal(clients)
		if marshalErr == nil {
			Cache.setDetails(cacheKey, string(byteClients))
		}
	}
	for _, client := range clients {
		details := fiber.Map{"Name": client.Name, "Email": client.MailId, "Phone": client.Phone, "Preference": client.Preference}
		clientDetails = append(clientDetails, details)
	}
	return ctx.Status(fiber.StatusOK).JSON(clientDetails)
}

// Logout godoc
// @Summary Logout
// @Description Logout
// @Tags Logout
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /privacy/logout [post]
func Logout(ctx *fiber.Ctx) error {
	userIdentifier := int(ctx.Locals("user_id").(float64))
	uuidDetail := ctx.Locals("uuid")
	repository.InvalidateCurrentSession(userIdentifier, uuidDetail)
	return ctx.SendStatus(fiber.StatusOK)
}

// ForgotPassword godoc
// @Summary ForgotPassword
// @Description ForgotPassword
// @Tags ForgotPassword
// @Accept json
// @Produce json
// @Param Body body models.ForgotPassword true "forgot password"
// @Success 200 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Router /auth/forgotPassword [post]
func ForgotPassword(ctx *fiber.Ctx) error {
	var forgotPasswd models.ForgotPassword
	err := ctx.BodyParser(&forgotPasswd)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	validationErr, _ := validator.Validate(forgotPasswd)
	if validationErr != nil {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	user, err := repository.GetUserByName(forgotPasswd.Email)
	if err != nil {
		log.Println("Invalid User .. User doesn't exists")
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	cacheKey := uuid.NewV4().String()
	cacheValue := user.Name
	Cache.setDetails(cacheKey, cacheValue)
	Cache.setExpiry(cacheKey, 3000)

	content := fmt.Sprintf("Post new password on given link: http://localhost:3001/auth/getNewPassword/%s to update your password. Request is valid for 5 minutes.", cacheKey)
	kafkaPayload := models.ForgotPasswordPayload{
		To:      user.Name,
		From:    config.GetConfig("MAILGUN_FROM"),
		Text: content,
		Subject: "Password Reset",
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
	msg := "Please check your mailbox for update password link. Request valid for 5 minutes"
	resp := fiber.Map{"message": msg}
	return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusOK, resp)
}

// GetNewPassword godoc
// @Summary GetNewPassword
// @Description GetNewPassword
// @Tags GetNewPassword
// @Accept json
// @Produce json
// @Param id path string true "uuid"
// @Success 200 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 412 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Router /auth/getNewPassword [post]
func GetNewPassword(ctx *fiber.Ctx) error {
	var pwd models.NewPassword
	cacheKey := ctx.Params("id")
	mail, _ := Cache.getDetails(cacheKey)
	if mail != "" {
		err := ctx.BodyParser(&pwd)
		if err != nil {
			return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
		}
		user, err := repository.GetUserByName(mail)
		if err != nil {
			return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
		}
		hashPassword, _ := utils.HashPassword(pwd.Password)
		err = repository.UpdateUserPassword(int(user.ID), hashPassword)
		if err != nil {
			return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
		}
		repository.InvalidateAllSessions(int(user.ID))
		return utils.ApiResponse(ctx, fiber.StatusOK, config.ApiConst[config.SUCCESS_PWD_CHANGE])
	}
	resp := fiber.Map{"status": "Error", "message": "Request to change password has been expired"}
	return utils.ApiResponseWithCustomMsg(ctx, fiber.StatusPreconditionFailed, resp)
}

// CheckAuditLog godoc
// @Summary CheckAuditLog
// @Description CheckAuditLog
// @Tags CheckAuditLog
// @Accept json
// @Produce json
// @Success 200 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Router /api/checkAuditLog [get]
func CheckAuditLog(ctx *fiber.Ctx) error {
	var auditDetails []fiber.Map
	userIdentifier := int(ctx.Locals("user_id").(float64))
	user, _ := repository.GetUserById(userIdentifier)
	if user.ID == 0 {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}
	audits, err := repository.GetAuditLog(user.ID)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}
	for _, audit := range audits {
		details := fiber.Map{"To": audit.To, "template": audit.TemplateName, "SentOn": audit.CreatedAt}
		auditDetails = append(auditDetails, details)
	}
	return ctx.Status(fiber.StatusOK).JSON(auditDetails)
}
