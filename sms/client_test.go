package sms

import (
	"github.com/retarus/retarus-go/common"
	"testing"
	"time"
)

func smsClientProvider() Client {
	config := NewConfigFromEnv(common.Europe)
	return NewClient(config)
}

func generateSms(amount int) []string {
	job := Job{
		Messages: []Message{
			{
				Text: "this is a test message",
				Recipients: []Recipient{
					{
						Dst:         "0049176000000000",
						CustomerRef: "retarus",
					},
				},
			},
		},
	}
	jobIds := []string{}
	client := smsClientProvider()
	for i := 1; i < amount+1; i++ {
		res, _ := client.Send(job)
		jobIds = append(jobIds, res)
	}
	return jobIds
}

func TestNormalSend(t *testing.T) {
	job := Job{
		Messages: []Message{
			{
				Text: "this is a test message",
				Recipients: []Recipient{
					{
						Dst:         "0049176000000000",
						CustomerRef: "retarus",
					},
				},
			},
		},
	}
	client := smsClientProvider()
	_, err := client.Send(job)
	if err != nil {
		t.Errorf("Error shouldn't happen here: %s", err)
	}
}
func TestGetReport(t *testing.T) {
	jobId := generateSms(1)
	client := smsClientProvider()
	time.Sleep(time.Duration(8 * time.Second))
	res, err := client.GetReport(jobId[0])
	if err != nil {
		t.Errorf("Error shouldn't happen here: %s", err)
	}
	if res == nil {
		t.Fatal("Response shouldn't be null")
	}
	if res.JobID != jobId[0] {
		t.Errorf("Returned Job Id isn't matching with requested one.")
	}
}

func TestGetUnknownReport(t *testing.T) {
	unkownJobId := "123412341234124"
	client := smsClientProvider()
	_, err := client.GetReport(unkownJobId)
	if err == nil {
		t.Errorf("Error should happen here")
	}
	//fmt.Println("response", err)
}
