// Example 02.delete_fax_report.go shows, how easily you can delete the fax report which are no longer needed, but they will be deleted after 30 days anyway.

package main

import (
	"fmt"
	"github.com/retarus/retarus-go/common"
	"github.com/retarus/retarus-go/fax"
	"log"
	"os"
)

func main() {
	// set your fax job id
	faxJobId := ""
	// Configure the client which is beeing used to send the
	username := os.Getenv("retarus_fax_username")
	password := os.Getenv("retarus_fax_password")
	customerNumber := os.Getenv("retarus_cuno")

	config := fax.NewConfig(username, password, customerNumber, common.Europe)
	client := fax.NewClient(config)
	res, err := client.DeleteReport(faxJobId)
	if err != nil {
		log.Fatalf("Error occured during delete request of report: %w", err)
	}
	fmt.Println("Response: ", res)
	fmt.Println("Job was successfully deleted")
}
