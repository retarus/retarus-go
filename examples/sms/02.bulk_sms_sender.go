// Retarus GmbH ©2023
package main

import (
	"encoding/csv"
	"fmt"
	"github.com/noirbizarre/gonja"
	"github.com/retarus/retarus-go/common"
	"github.com/retarus/retarus-go/sms"
	"log"
	"os"
	"time"
)

func main() {
	// load needed credentials
	username := os.Getenv("retarus_username")
	password := os.Getenv("retarus_password")

	// create client and needed config for Retarus SMS-Client
	config := sms.NewConfig(username, password, common.Europe)
	client := sms.Client{
		Config:      config,
		Transporter: common.NewTransporter(5),
	}

	// load csv file containing the campaign contact which will receive the sms
	file, err := os.Open("assets/sms_data.csv")
	if err != nil {
		fmt.Println("Could not find load sms csv list")
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("failed to read the CSV file: %s", err)
	}
	messages := []sms.Message{}

	// This script will only send one sms job which contains multiple messages (multiple sms) which each has different receivers.
	for index, record := range records {
		if index == 0 {
			continue
		}
		// load campainge message
		b, err := os.ReadFile("assets/advertisement.txt") // just pass the file name
		if err != nil {
			panic(err)
		}
		// render the template with the surname of the target.
		tpl, err := gonja.FromString(string(b))
		if err != nil {
			panic(err)
		}
		out, err := tpl.Execute(gonja.Context{"firstname": "axel"})
		if err != nil {
			panic(err)
		}

		fmt.Println("index: ", index, "Record", record)

		// parse values into retarus data structure
		recipient := []sms.Recipient{sms.NewRecipient(record[2], "example_02_go_sdk", []time.Time{})} // Assuming sms.Receipt exists
		message := sms.NewMessage(out, recipient)
		messages = append(messages, message) // Assign the result back to the messages slice.
	}
	fmt.Println("Messages: ", messages)
	res, err := client.Send(sms.NewJob(messages, &sms.Options{}))

	if err != nil {
		log.Println("Error occurred during sending of sms messages: ", err)
	}
	fmt.Println("Successfully send sms with JobId: ", res)
}
