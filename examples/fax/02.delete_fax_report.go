// Example 02.delete_fax_report.go shows, how easily you can delete the fax report which are no longer needed, but they will be deleted after 30 days anyway.

package main

import (
	"fmt"
	"github.com/retarus/retarus-go/common"
	"github.com/retarus/retarus-go/fax"
	"os"
)

func main() {
	// set your fax job id
	faxJobId := ""
	// Configure the client which is beeing used to send the
	username := os.Getenv("retarus_username")
	password := os.Getenv("retarus_password")
	customerNumber := os.Getenv("retarus_cuno")

	config := fax.NewConfig(username, password, customerNumber, common.Europe)
	client := fax.NewClient(config)
	_, err := client.DeleteReport(faxJobId)
	if err == nil {
		fmt.Errorf("failed to execute someFunction: %w", err)
		panic("")
	}
	fmt.Println("Job was successfully deleted")
}
