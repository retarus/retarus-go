package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/retarus/retarus-go/common"
	"net/http"
	"net/url"
)

// Client is responsible for sending requests and handling transportation for an SMS service.
// It needs to be properly configured to interact with the service.
// Note: To create a new instance of Client, use the NewClient function.
type Client struct {
	// Config holds the configuration settings for the client.
	Config Config

	// Transporter is responsible for the actual HTTP requests and responses.
	Transporter common.Transporter
}

// NewClient creates and returns a new Client instance.
// The client is configured based on the provided Config object.
// The Transporter is initialized with a default timeout of 5 seconds.
//
// This is the preferred way to create a new Client instance.
//
// Parameters:
//   - config: A Config object containing settings like API keys, base URLs, etc.
//
// Returns:
//   - A configured Client object ready to send requests to the SMS service.
func NewClient(config Config) Client {
	return Client{
		Config:      config,
		Transporter: common.NewTransporter(5), // Initialize transporter with a timeout of 5 seconds
	}
}

// Send sends a sms job to the specified numbers in the job.
func (c *Client) Send(job Job) (jobID string, err error) {
	jobBytes, err := json.Marshal(job)
	if err != nil {
		return "", err
	}

	u, err := url.JoinPath(c.Config.Region.HAAddr, "/jobs")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(jobBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.Config.User, c.Config.Password)

	resp, err := c.Transporter.HTTPClient.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := statusToError(resp.StatusCode, resp.Body); err != nil {
		return "", err
	}

	type jobResp struct {
		JobID string `json:"jobId"`
	}

	var jobResponse jobResp
	if err := json.NewDecoder(resp.Body).Decode(&jobResponse); err != nil {
		return "", err
	}

	return jobResponse.JobID, nil
}

// GetReport retrieves the status and list of SMS IDs for a specific job by its job ID.
// It uses the configured Transporter to make an HTTP GET request to the service.
//
// This method is designed to get the job status and the list of SMS IDs for the job.
// To get the individual SMS status for all SMS's of a job, use the SMS status endpoint.
// Parameters:
//   - jobID: A string representing the job ID for which the report is being requested.
//
// Returns:
//   - A pointer to a Report object containing details about the job's SMS statuses and IDs.
//   - An error object if an error occurs during the fetch operation or if no report is found.
func (c *Client) GetReport(jobID string) (*Report, error) {
	var smsReport Report

	resp := c.Transporter.DoDatacenterFetch(c.Config.Region.Servers, c.Config.User, c.Config.Password, []byte{}, "/jobs/"+jobID, http.MethodGet)
	if len(resp) == 0 {
		return nil, errors.New("Error occured during fetch of responses.")
	}
	for x := range resp {

		//fmt.Println("Server: ", resp[x].Request.URL, "Response-Yeah: ", resp[x].StatusCode)
		if resp[x].StatusCode == 404 {
			continue
		}
		defer resp[x].Body.Close()

		if err := statusToError(resp[x].StatusCode, resp[x].Body); err != nil {
			return nil, err
		}

		if err := json.NewDecoder(resp[x].Body).Decode(&smsReport); err != nil {
			return nil, err
		}
	}
	if smsReport.IsZero() == true {
		return nil, errors.New("no reports found, try again later or contact customer service")
	}
	return &smsReport, nil
}

// GetSmsStatus retrieves the status of individual SMS messages within a given job.
// It uses the configured Transporter to make an HTTP GET request to the service.
//
// This method is designed to get the individual SMS status for all SMS's of a job.
//
// Parameters:
//   - jobID: A string representing the job ID for which the statuses are being requested.
//
// Returns:
//   - A pointer to an SmsStatus object containing details about the individual SMS statuses within the job.
//   - An error object if an error occurs during the fetch operation or if no statuses are found.
func (c *Client) GetSmsStatus(jobID string) (*[]SmsStatus, error) {
	var status []SmsStatus

	parms := common.KvParams{Key: "jobId", Value: jobID}
	resp := c.Transporter.DoDatacenterFetch(c.Config.Region.Servers, c.Config.User, c.Config.Password, []byte{}, "/sms", http.MethodGet, parms)
	for x := range resp {
		if resp[x].StatusCode == 404 {
			continue
		}
		defer resp[x].Body.Close()

		if err := statusToError(resp[x].StatusCode, resp[x].Body); err != nil {
			return nil, err
		}

		if err := json.NewDecoder(resp[x].Body).Decode(&status); err != nil {
			return nil, err
		}
	}
	if len(status) == 0 {
		return nil, errors.New("no reports found, try again later or contact customer service")
	}
	return &status, nil
}
