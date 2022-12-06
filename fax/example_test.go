package fax_test

import (
	"log"
	"net/http"
	"time"

	"github.com/retarus/retarus-go/fax"
)

func Example_Send() {
	client := fax.Client{
		Config: fax.Config{
			User:           "test@example.com",
			Password:       "",
			CustomerNumber: "99999TE",
			Endpoint:       fax.DE,
		},
		HTTPClient: http.Client{Timeout: 5 * time.Second},
	}

	job := fax.Job{
		Recipients: []fax.Recipient{
			{
				Number: "004989000000000",
			},
		},
		Documents: []fax.Document{
			{
				Name:      "test.txt",
				Charset:   fax.UTF_8,
				Reference: "testJob",
				Data:      "dGVzdGZheAo=", // testfax
			},
		},
	}

	jobID, err := client.Send(job)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	log.Printf("Fax send with JobID %s", jobID)
}

func Example_GetAndDeleteReports() {
	client := fax.Client{
		Config: fax.Config{
			User:           "test@example.com",
			Password:       "",
			CustomerNumber: "99999TE",
			Endpoint:       fax.DE,
		},
		HTTPClient: http.Client{Timeout: 5 * time.Second},
	}

	reports, err := client.GetReports()
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	for _, report := range reports {
		// single report delete, we also could use bulk
		delReport, err := client.DeleteReport(report.JobID)
		if err != nil {
			log.Printf("error deleting report: %s", err.Error())
			continue
		}

		log.Printf("Report for jobID %s deleted: %v", delReport.JobID, delReport.Deleted)
	}
}
