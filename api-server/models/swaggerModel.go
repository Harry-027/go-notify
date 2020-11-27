package models

// models for swagger definition
type SwaggerClient struct {
	Name       string `json:"name"`
	MailId     string `json:"mailID"`
	Phone      int    `json:"phone"`
	Preference string `json:"preference"`
}

type SwaggerTemplate struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
