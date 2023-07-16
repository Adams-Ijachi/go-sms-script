package main

import (
	"fmt"

	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"

	"context"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	lambda.Start(HandleLambdaEvent)
}

// get word from json
type Data struct {
	Entries []string `json:"entries"`
}

func getWord() string {
	// date format should be 02- January - 2006

	date := time.Now().Format("2 January 2006")

	message := fmt.Sprintf("It is %s", date)

	return message

}

func HandleLambdaEvent(ctx context.Context, event map[string]interface{}) (string, error) {
	godotenv.Load()

	word := getWord()

	client := twilio.NewRestClient()

	params := &api.CreateMessageParams{}

	params.SetBody(word)

	params.SetFrom(os.Getenv("TWILIO_NUMBER"))
	params.SetTo(os.Getenv("USER_NUMBER"))

	resp, err := client.Api.CreateMessage(params)

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	} else {
		if resp.Sid != nil {

			log.Println("Message sent successfully")

			event["message"] = "Message sent successfully"
		} else {
			fmt.Println(resp)
			log.Fatal(resp)
		}
	}

	return "Message sent successfully", nil

}
