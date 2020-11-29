package config

type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// map for Api response ...
var ApiConst = map[string]ApiResponse{
	SIGNUP_SUCCESS: {
		Status:  "Success",
		Message: "Congrats! Signup completed successfully",
	},
	SIGNUP_FAILED: {
		Status:  "Error",
		Message: "SignUp failed !!",
	},
	INCORRECT_PASSWORD_CONRIRM: {
		Status:  "Error",
		Message: "Password confirmation is incorrect",
	},
	VALIDATION_ERROR: {
		Status:  "Error",
		Message: "Payload validation failed",
	},
	BAD_REQUEST: {
		Status:  "Bad Request",
		Message: "Bad Payload",
	},
	NO_USER_DETAILS: {
		Status:  "Error",
		Message: "User details not found !!",
	},
	INTERNAL_SERVER_ERROR: {
		Status:  "Internal Server Error",
		Message: "An internal server error occurred !!",
	},
	SUCCESS_PWD_CHANGE: {
		Status:  "Success",
		Message: "Password changed successfully !!",
	},
	SUCCESS_ACC_DEL: {
		Status:  "Success",
		Message: "Account deleted successfully !!",
	},
	UNAUTHORIZED_ACTION: {
		Status:  "UnAuthorized",
		Message: "UnAuthorized Action",
	},
	FORBIDDEN: {
		Status:  "Forbidden",
		Message: "Not allowed to perform this action !!",
	},
	SUBS_EXPIRED: {
		Status:  "Error",
		Message: "Either subscription expired or more number of mails were expected !!..",
	},
	JOB_DEL: {
		Status:  "Ok",
		Message: "Job deleted successfully !!",
	},
	SUBS_SUCCESS: {
		Status:  "Success",
		Message: "Subscription Successful !!",
	},
	NO_CONTENT: {
		Status:  "Not Found",
		Message: "Details not found",
	},
}

// constants ...
const (
	SERVER_PORT                = "SERVER_PORT"
	ADMIN_ROLE                 = "admin"
	KAFKA_PROTOCOL             = "KAFKA_PROTOCOL"
	BROKER                     = "BROKER"
	TOPIC                      = "TOPIC"
	FREEPLAN                   = "sandbox"
	SIGNUP_FAILED              = "signUpFailed"
	SIGNUP_SUCCESS             = "signUpSuccess"
	INCORRECT_PASSWORD_CONRIRM = "incorrectPasswordConfirm"
	VALIDATION_ERROR           = "validationError"
	BAD_REQUEST                = "badRequest"
	NO_USER_DETAILS            = "noUserDetails"
	INTERNAL_SERVER_ERROR      = "internalServerError"
	SUCCESS_PWD_CHANGE         = "successPwdChange"
	SUCCESS_ACC_DEL            = "successAccDel"
	UNAUTHORIZED_ACTION        = "unauthorized"
	SECRET                     = "secret"
	SUBS_EXPIRED               = "subscriptionExpired"
	JOB_DEL                    = "jobDel"
	SUBS_SUCCESS               = "subsSuccess"
	NO_CONTENT                 = "nocontent"
	FORBIDDEN                  = "forbidden"
	DB_PORT                    = "DB_PORT"
	SERVER_URL                 = "SERVER_URL"
)
