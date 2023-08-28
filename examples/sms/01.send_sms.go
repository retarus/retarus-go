package main

import (
	"github.com/retarus/retarus-go/common"
	"github.com/retarus/retarus-go/sms"
	"log"
	"time"
)

func main() {
	config := sms.NewConfigFromEnv(common.Europe)
	client := sms.Client{
		Config:      config,
		Transporter: common.NewTransporter(5),
	}

	job := sms.Job{
		Messages: []sms.Message{
			{
				Text: "this is a test message",
				Recipients: []sms.Recipient{
					{
						Dst:         "0049176000000000",
						CustomerRef: "retarus",
					},
				},
			},
		},
	}
	jobID, err := client.Send(job)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	log.Printf("SMS send with JobID %s", jobID)

	time.Sleep(20 * time.Second)

	rep, err := client.GetReport(jobID)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	log.Printf("Report found: %+v", rep)
}
