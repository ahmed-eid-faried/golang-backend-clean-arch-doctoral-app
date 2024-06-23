package utils

import (
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SendSMS sends an SMS using Twilio
func SendSMS(toPhoneNumber, message string) error {
	// client := twilio.NewRestClient()
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})
	params := &openapi.CreateMessageParams{}
	params.SetTo(toPhoneNumber)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(message)

	_, err := client.Api.CreateMessage(params)
	return err

}
