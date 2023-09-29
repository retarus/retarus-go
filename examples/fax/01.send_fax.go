// "01_send_fax.go" elegantly demonstrates how to dispatch a fax (in this case a really simple config and options) using the SDKS's minimalist and intuitive settings.

package main

import (
	"fmt"
	"github.com/retarus/retarus-go/common"
	"github.com/retarus/retarus-go/fax"
	"os"
)

func main() {
	// get environment variables
	username := os.Getenv("retarus_fax_username")
	password := os.Getenv("retarus_fax_password")
	customerNumber := os.Getenv("retarus_cuno")

	// setup client
	config := fax.NewConfig(username, password, customerNumber, common.Europe)
	client := fax.NewClient(config)

	// create job
	job := fax.Job{
		Recipients: []fax.Recipient{
			{
				Number: "004989000000000", // number to send to
			},
		},
		Documents: []fax.Document{
			{
				Name: "test.txt", // local document to send
				// both options under here are not really working, backend returns 500
				Charset: fax.UTF_8,
				Data:    "dGVzdGZheAo=", // test fax
			},
		},
	}
	// dispatch and wait for report
	jobID, err := client.Send(job)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("JobId: ", jobID)
}
