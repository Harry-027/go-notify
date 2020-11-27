package validator

import (
	"fmt"
	"github.com/Harry-027/go-notify/api-server/models"
	"github.com/asaskevich/govalidator"
	"log"
	"reflect"
)

// Client preference type ...
var preferenceType = map[string]bool{
	"daily":   true,
	"weekly":  true,
	"monthly": true,
}

// Validate the given struct ...
func Validate(v interface{}) (error, bool) {
	result, err := govalidator.ValidateStruct(v)
	if err != nil {
		log.Println("Invalid data", err)
	}
	return err, result
}

// Validate the slice of struct ...
func ValidateSliceOfStruct(v interface{}) (bool, error) {
	isBadPayload := false
	interfaceType := reflect.TypeOf(v)
	if interfaceType.Kind() == reflect.Slice {
		clients := v.([]models.Client)
		for _, client := range clients {
			_, err := govalidator.ValidateStruct(client)
			if err != nil {
				log.Println("Invalid data", err)
				isBadPayload = true
			}
		}
		return isBadPayload, nil
	}
	isBadPayload = true
	err := fmt.Errorf("given %v is not a valid slice", v)
	return isBadPayload, err
}

// Validate the client preference ...
func ValidatePreference(pref string) bool {
	_, ok := preferenceType[pref]
	return ok
}
