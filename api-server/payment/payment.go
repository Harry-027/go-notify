package payment

import (
	"errors"
	"github.com/Harry-027/go-notify/api-server/models"
	"log"
	"strings"
)

// monthly plan subscription enum
var SubsEnum = map[string]int{
	"sandbox":  100,
	"gold":     10000,
	"silver":   5000,
	"platinum": 1000,
}

// Make payment (stubbed for now)...
func MakePayment(subs models.SubscriptionInput) (int, error) {
	count, ok := SubsEnum[strings.ToLower(subs.SubscriptionType)]
	if !ok {
		return 0, errors.New("invalid subscription plan")
	}
	// Integrate payment gateway here
	log.Println("Payment Type ::", subs.PaymentType) // Credit,Debit,Net Banking
	return count, nil
}
