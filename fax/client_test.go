package fax

import (
	"fmt"
	"github.com/retarus/retarus-go/common"
	"testing"
)

func faxClientProvider() Client {
	config := NewConfigFromEnv(common.Europe)
	client := NewClient(config)
	return client
}

func faxGeneratorCreator(amount int) []string {
	job := Job{
		Recipients: []Recipient{
			{
				Number: "+4989000000000", // number to send to
			},
		},
		Documents: []Document{
			{
				Name: "test.txt",     // local document to send
				Data: "dGVzdGZheAo=", // testfax
			},
		},
	}
	jobIds := []string{}
	client := faxClientProvider()
	for i := 1; i < amount+1; i++ {
		res, _ := client.Send(job)
		jobIds = append(jobIds, res)
	}
	return jobIds
}

// send normal fax
func TestNormalFaxSend(t *testing.T) {
	job := Job{
		Recipients: []Recipient{
			{
				Number: "+4989000000000", // number to send to
			},
		},
		Documents: []Document{
			{
				Name: "test.txt",     // local document to send
				Data: "dGVzdGZheAo=", // testfax
			},
		},
	}
	client := faxClientProvider()
	res, err := client.Send(job)

	if err != nil {
		t.Errorf("Error shouldn't happen here: %s", err)
	}
	if res == "" {
		t.Errorf("Error, empty job id, should contain id")
	}
}

func TestGetNormalFaxReport(t *testing.T) {
	client := faxClientProvider()
	jobIds := faxGeneratorCreator(1)

	first := jobIds[0]
	res, err := client.GetReport(first)
	if err != nil {
		t.Errorf("Error shouldn't happen here: %s", err)
	}
	fmt.Println(res)
}

func TestGetBulkReport(t *testing.T) {
	client := faxClientProvider()
	jobIds := faxGeneratorCreator(5)
	res, err := client.GetBulkReports(jobIds)
	if err != nil {
		t.Errorf("Error shouldn't happen here: %s", err)
	}
	if len(res) != 5 {
		t.Errorf("Missing some reports: %s", err)

	}
}
func TestDeleteReports(t *testing.T) {
	client := faxClientProvider()
	jobIds := faxGeneratorCreator(5)
	res, err := client.DeleteBulkReports(jobIds)
	if err != nil {
		t.Errorf("Error shouldn't happen here: %s", err)
	}
	fmt.Println(res)
	if len(res) != 5 {
		t.Errorf("Amount of delted report: %d", len(res))
	}
}

func TestGetBulkFaxReport(t *testing.T) {
	client := faxClientProvider()
	jobIds := faxGeneratorCreator(5)
	res, err := client.GetBulkReports(jobIds)
	if err != nil {
		t.Errorf("Error shouldn't happen here: %s", err)
	}
	fmt.Println(res)
	if len(res) != 5 {
		t.Errorf("Amount of got report: %d", len(res))
	}
}
