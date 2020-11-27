package notifier

import (
	"context"
	"encoding/json"
	"github.com/mailgun/mailgun-go/v4"
	"log"
	"time"
)

// Notifier ...
type MailgunNotifier struct {
	ApiKey string
	Domain string
	From   string
}

// Payload ...
type KafkaPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Html    string `json:"text"`
	Subject string `json:"subject"`
}

// Send the mail ...
func (notifier MailgunNotifier) Send(payload []byte) (string, error) {

	var details KafkaPayload
	err := json.Unmarshal(payload, &details)

	if err != nil {
		log.Println("An error occurred while unmarshalling ", err)
		return "", err
	}

	mg := mailgun.NewMailgun(notifier.Domain, notifier.ApiKey)
	m := mg.NewMessage(
		details.From,
		details.Subject,
		"",
		details.To,
	)
	m.SetHtml(details.Html)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	msg, id, err := mg.Send(ctx, m)
	if err != nil {
		log.Println("An error occurred while sending mail ::", err)
	}
	log.Println("Mail sent successfully !!", msg)
	return id, err
}
