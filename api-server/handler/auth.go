package handler

import (
	"github.com/Harry-027/go-notify/api-server/config"
	"github.com/Harry-027/go-notify/api-server/models"
	"github.com/Harry-027/go-notify/api-server/repository"
	"github.com/Harry-027/go-notify/api-server/utils"
	"github.com/Harry-027/go-notify/api-server/validator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

// Login godoc
// @Summary Authenticate the user
// @Description Authenticates a user and provides jwt token
// @Tags login
// @Accept json
// @Produce json
// @Param Body body models.LoginInput true "User login"
// @Success 200 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Router /auth/login [post]
func Login(ctx *fiber.Ctx) error {
	var input models.LoginInput

	err := ctx.BodyParser(&input)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}
	validationErr, result := validator.Validate(input)
	log.Println("Validation result: ", result)

	if validationErr != nil {
		log.Println("An error occurred :: ", validationErr.Error())
		return utils.ApiResponse(ctx, fiber.StatusPreconditionFailed, config.ApiConst[config.VALIDATION_ERROR])
	}

	user, err := repository.GetUserByName(input.Email)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusNoContent, config.ApiConst[config.NO_USER_DETAILS])
	}

	isValid := utils.CheckPasswordHash(input.Password, user.Password)
	if isValid == false {
		return utils.ApiResponse(ctx, fiber.StatusUnauthorized, config.ApiConst[config.UNAUTHORIZED_ACTION])
	}
	newUuid := uuid.NewV4().String()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Name
	claims["user_id"] = user.ID
	claims["user_role"] = user.Role
	claims["issuer"] = "Harish Bhawnani"
	claims["uuid"] = newUuid
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(config.GetConfig(config.JWT_SECRET)))
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}
	auth := models.Auth{
		Uuid:   newUuid,
		UserID: user.ID,
	}
	repository.SaveUserAuth(auth)
	return ctx.JSON(fiber.Map{"status": "Success", "message": "Success login", "token": t})
}

// Signup godoc
// @Summary The user account registration
// @Description Registers the user account
// @Tags signup
// @Accept json
// @Produce json
// @Param Body body models.SignupInput true "User SignupDetails"
// @Success 200 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Router /auth/signup [post]
func Signup(ctx *fiber.Ctx) error {
	var input models.SignupInput

	err := ctx.BodyParser(&input)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}
	log.Println("Signup Payload: ", input)
	validationErr, result := validator.Validate(input)
	log.Println("Validation result: ", result)

	if validationErr != nil {
		log.Println("Validation Error: ", validationErr.Error())
		return utils.ApiResponse(ctx, fiber.StatusPreconditionFailed, config.ApiConst[config.VALIDATION_ERROR])
	}

	if input.Password != input.ConfirmPassword {
		return utils.ApiResponse(ctx, fiber.StatusNotAcceptable, config.ApiConst[config.INCORRECT_PASSWORD_CONRIRM])
	}

	hashPassword, _ := utils.HashPassword(input.Password)
	user := models.User{}
	user.Name = input.Email
	user.Password = hashPassword
	user.Role = input.Role
	user.Subscription = config.FREEPLAN
	user.NotificationCounter = 100
	err = repository.AddUser(user)
	Cache.deleteKey("USERS")
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusExpectationFailed, config.ApiConst[config.SIGNUP_FAILED])
	}
	return utils.ApiResponse(ctx, fiber.StatusOK, config.ApiConst[config.SIGNUP_SUCCESS])
}

// UpdatePassword godoc
// @Summary The user password update
// @Description Helps the user to update the password
// @Tags UpdatePassword
// @Accept json
// @Produce json
// @Param Body body models.UpdatePassword true "Update Password"
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 400 {object} config.ApiResponse{}
// @Failure 412 {object} config.ApiResponse{}
// @Failure 204 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /privacy/updatePassword [post]
func UpdatePassword(ctx *fiber.Ctx) error {
	var input models.UpdatePassword

	err := ctx.BodyParser(&input)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusBadRequest, config.ApiConst[config.BAD_REQUEST])
	}

	validationErr, result := validator.Validate(input)
	log.Println("Validation result: ", result)

	if validationErr != nil {
		log.Println("Validation Error: ", validationErr.Error())
		return utils.ApiResponse(ctx, fiber.StatusPreconditionFailed, config.ApiConst[config.VALIDATION_ERROR])
	}

	userIdentifier := int(ctx.Locals("user_id").(float64))
	if userIdentifier == 0 {
		return utils.ApiResponse(ctx, fiber.StatusNoContent, config.ApiConst[config.NO_USER_DETAILS])
	}

	hashPassword, _ := utils.HashPassword(input.Password)
	err = repository.UpdateUserPassword(userIdentifier, hashPassword)
	if err != nil {
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}
	repository.InvalidateAllSessions(userIdentifier)
	return utils.ApiResponse(ctx, fiber.StatusOK, config.ApiConst[config.SUCCESS_PWD_CHANGE])
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
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "message": "Logged out"})
}

// DeleteAccount godoc
// @Summary Delete the User Account
// @Description Deletes the user account permanently. All the data gets lost permanently
// @Tags Delete the Account
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {object} config.ApiResponse{}
// @Failure 500 {object} config.ApiResponse{}
// @Security ApiKeyAuth
// @Router /privacy/deleteAccount [delete]
func DeleteAccount(ctx *fiber.Ctx) error {
	userIdentifier := int(ctx.Locals("user_id").(float64))
	err := repository.DeleteUser(userIdentifier)
	if err != nil {
		log.Println("Error deleting Account for User: ", userIdentifier)
		return utils.ApiResponse(ctx, fiber.StatusInternalServerError, config.ApiConst[config.INTERNAL_SERVER_ERROR])
	}
	return utils.ApiResponse(ctx, fiber.StatusOK, config.ApiConst[config.SUCCESS_ACC_DEL])
}
