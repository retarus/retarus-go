package sms_test

import (
	"log"
	"net/http"
	"time"

	"github.com/retarus/retarus-go/sms"
)

func Example_SendAndGetReport() {
	client := sms.Client{
		Config: sms.Config{
			User:     "test@example.com",
			Password: "",
			Endpoint: sms.DE1,
		},
		HTTPClient: http.Client{Timeout: 5 * time.Second},
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

func Example_GetAndPrintReports() {
	client := sms.Client{
		Config: sms.Config{
			User:     "test@example.com",
			Password: "",
			Endpoint: sms.DE1,
		},
		HTTPClient: http.Client{Timeout: 5 * time.Second},
	}

	jobids, err := client.GetJobIDs(sms.Params{
		JobIDsOnly: true,
	})
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	for _, jobid := range jobids {
		log.Printf("Report found for jobID %s", jobid)
	}
}
